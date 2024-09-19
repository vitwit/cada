package relayer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/retriever"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubmitDataToAvailClient submits data to the Avail client using an HTTP POST request.
func (r *Relayer) SubmitDataToAvailClient(_ string, _ int, data []byte, blocks []int64, lightClientURL string) (BlockInfo, error) {
	var blockInfo BlockInfo

	handler := NewHTTPClientHandler()
	datab := base64.StdEncoding.EncodeToString(data)

	jsonData := []byte(fmt.Sprintf(`{"data":"%s"}`, datab))

	url := fmt.Sprintf("%s/v2/submit", lightClientURL)

	// Make the POST request
	responseBody, err := handler.Post(url, jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return blockInfo, err
	}

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(responseBody, &blockInfo)
	if err != nil {
		r.logger.Error("Error while posting block(s) to Avail DA",
			"height_start", blocks[0],
			"height_end", blocks[len(blocks)-1],
			"appID", strconv.Itoa(r.rpcClient.config.AppID),
		)
	}

	if err == nil {
		r.logger.Info("Successfully posted block(s) to Avail DA",
			"height_start", blocks[0],
			"height_end", blocks[len(blocks)-1],
			"appID", strconv.Itoa(r.rpcClient.config.AppID),
			"block_hash", blockInfo.BlockHash,
			"block_number", blockInfo.BlockNumber,
			"hash", blockInfo.Hash,
		)
	}

	return blockInfo, nil
}

// GetSubmittedData retrieves the data submitted to the Avail client for a specific block.
// It constructs a URL based on the given light client URL and block number, performs an HTTP GET request to fetch
// the data, and decodes the response into a BlockData struct.
func (r *Relayer) GetSubmittedData(lightClientURL string, blockNumber int) (BlockData, error) {
	handler := NewHTTPClientHandler()

	// Construct the URL with the block number
	url := fmt.Sprintf("%s/v2/blocks/%v/data?fields=data", lightClientURL, blockNumber)

	// Perform the GET request, returning the body directly
	body, err := handler.Get(url)
	if err != nil {
		return BlockData{}, fmt.Errorf("failed to fetch data from the avail: %w", err)
	}

	// Decode the response body into the BlockData struct
	var blockData BlockData
	err = json.Unmarshal(body, &blockData)
	if err != nil {
		return BlockData{}, fmt.Errorf("failed to decode block response: %w", err)
	}

	return blockData, nil
}

// IsDataAvailable is to query the avail light client and check if the data is made available at the given height
func (r *Relayer) IsDataAvailable(ctx sdk.Context, from, to, availHeight uint64, lightClientURL string) (bool, error) {
	availBlock, err := r.GetSubmittedData(lightClientURL, int(availHeight))
	if err != nil {
		return false, err
	}

	var blocks []int64
	for i := from; i <= to; i++ {
		blocks = append(blocks, int64(i))
	}

	cosmosBlocksData := r.GetBlocksDataFromLocal(ctx, blocks)
	base64CosmosBlockData := base64.StdEncoding.EncodeToString(cosmosBlocksData)

	// TODO: any better / optimized way to check if data is really available?
	return isDataIncludedInBlock(availBlock, base64CosmosBlockData), nil
}

// isDataIncludedInBlock checks whether the provided Cosmos block data is included in the given Avail block.
func isDataIncludedInBlock(availBlock BlockData, base64cosmosData string) bool {
	for _, data := range availBlock.Extrinsics {
		if data.Data == base64cosmosData {
			return true
		}
	}

	return false
}

// GetBlock represents the data structure for a block with its associated transactions.
type GetBlock struct {
	BlockNumber      int      `json:"block_number"`
	DataTransactions []string `json:"data_transactions"`
}

// SubmitDataToAvailDA submits data to the Avail light client using a provided API URL and seed.
// The function connects to a Substrate API, creates an extrinsic with the provided data and appID, signs it using
// the seed, and submits it to the chain. It then monitors the transaction status and checks for finalization.
func (r *Relayer) SubmitDataToAvailDA(apiURL, seed string, availAppID int, data []byte, blocks []int64) error {
	api, err := gsrpc.NewSubstrateAPI(apiURL)
	if err != nil {
		r.logger.Error("cannot create api:%w", err)
		return err
	}

	fmt.Println("appurl", apiURL, "seed", seed)
	fmt.Println("created substrate api")

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		r.logger.Error("cannot get metadata:%w", err)
		return err
	}

	// Set data and appID according to need
	// data, _ := RandToken(size)
	appID := 0

	// if app id is greater than 0 then it must be created before submitting data
	if availAppID != 0 {
		appID = availAppID
	}

	c, err := types.NewCall(meta, "DataAvailability.submit_data", types.NewBytes(data))
	if err != nil {
		r.logger.Error("cannot create new call:%w", err)
		return err
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		r.logger.Error("cannot get block hash:%w", err)
		return err
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		r.logger.Error("cannot get runtime version:%w", err)
		return err
	}

	keyringPair, err := signature.KeyringPairFromSecret(seed, 42)
	if err != nil {
		r.logger.Error("cannot create KeyPair:%w", err)
		return err
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyringPair.PublicKey)
	if err != nil {
		r.logger.Error("cannot create storage key:%w", err)
		return err
	}

	var accountInfo types.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	fmt.Println("ok", ok, err)
	if err != nil || !ok {
		r.logger.Error("cannot get latest storage:%w", err)
		return err
	}

	nonce := uint32(accountInfo.Nonce)

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		AppID:              types.NewUCompactFromUInt(uint64(0)),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction using Alice's default account
	err = ext.Sign(keyringPair, o)
	if err != nil {
		r.logger.Error("cannot sign:%w", err)
		return err
	}

	// Send the extrinsic
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		r.logger.Error("cannot submit extrinsic:%v", err)
		return err
	}

	fmt.Printf("** Data submitted to Avail DA using APPID: ** %v \n", appID)
	if err == nil {
		r.logger.Info("Posted block(s) to Avail DA",
			"height_start", blocks[0],
			"height_end", blocks[len(blocks)-1],
			"appID", strconv.Itoa(r.rpcClient.config.AppID),
		)
	}

	defer sub.Unsubscribe()
	timeout := time.After(100 * time.Second)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				fmt.Printf("Txn inside block %v\n", status.AsInBlock.Hex())
			} else if status.IsFinalized {
				retriever, err := retriever.NewDefaultEventRetriever(state.NewEventProvider(api.RPC.State), api.RPC.State)
				if err != nil {
					r.logger.Error("Couldn't create event retriever: %s", err)
					return err
				}

				h := status.AsFinalized
				events, err := retriever.GetEvents(h)
				if err != nil {
					r.logger.Error("Couldn't retrieve events")
				}

				for _, event := range events {
					if event.Name == "DataAvailability.DataSubmitted" {
						from, _ := registry.ProcessDecodedFieldValue[*types.AccountID](
							event.Fields,
							func(_ int, field *registry.DecodedField) bool {
								return field.Name == "sp_core.crypto.AccountId32.who"
							},
							func(value any) (*types.AccountID, error) {
								fields, ok := value.(registry.DecodedFields)

								if !ok {
									return nil, fmt.Errorf("unexpected value: %v", value)
								}

								accByteSlice, err := registry.GetDecodedFieldAsSliceOfType[types.U8](fields, func(fieldIndex int, _ *registry.DecodedField) bool {
									return fieldIndex == 0
								})
								if err != nil {
									return nil, err
								}

								var accBytes []byte

								for _, accByte := range accByteSlice {
									accBytes = append(accBytes, byte(accByte))
								}

								return types.NewAccountID(accBytes)
							},
						)
						a := from.ToHexString()

						// // add, _ := types.NewAddressFromHexAccountID(a)
						// fmt.Println(from)
						fmt.Printf("from address read from event: %s \n", a)
						dataHash, err := registry.ProcessDecodedFieldValue[*types.Hash](
							event.Fields,
							func(fieldIndex int, _ *registry.DecodedField) bool {
								return fieldIndex == 1
							},
							func(value any) (*types.Hash, error) {
								fields, ok := value.(registry.DecodedFields)
								if !ok {
									return nil, fmt.Errorf("unexpected value: %v", value)
								}

								hashByteSlice, err := registry.GetDecodedFieldAsSliceOfType[types.U8](fields, func(fieldIndex int, _ *registry.DecodedField) bool {
									return fieldIndex == 0
								})
								if err != nil {
									return nil, err
								}

								var hashBytes []byte
								for _, hashByte := range hashByteSlice {
									hashBytes = append(hashBytes, byte(hashByte))
								}

								hash := types.NewHash(hashBytes)
								return &hash, nil
							},
						)

						switch {
						case err != nil:
							fmt.Printf("DataHash parsing err: %s\n", err.Error())
						case dataHash == nil:
							fmt.Println("DataHash is nil")
						default:
							fmt.Printf("DataHash read from event: %s \n", dataHash.Hex())
						}

					}
				}
				fmt.Printf("Txn inside finalized block\n")
				hash := status.AsFinalized
				err = GetDataFromAvailDA(hash, api, string(data))
				if err != nil {
					r.logger.Error("cannot get data:%v", err)
					return err
				}
				return nil
			}
		case <-timeout:
			r.logger.Warn("timeout of 100 seconds reached without getting finalized status for extrinsic")
			return nil
		}
	}
}

// GetDataFromAvailDA retrieves and verifies submitted data from the Avail Data Availability (DA) chain
// by searching for the given data in the extrinsics of a block identified by the provided block hash.
func GetDataFromAvailDA(hash types.Hash, api *gsrpc.SubstrateAPI, data string) error {
	block, err := api.RPC.Chain.GetBlock(hash)
	if err != nil {
		return fmt.Errorf("cannot get block by hash:%w", err)
	}

	for _, ext := range block.Block.Extrinsics {
		// these values below are specific indexes only for data submission, differs with each extrinsic
		if ext.Method.CallIndex.SectionIndex == 29 && ext.Method.CallIndex.MethodIndex == 1 {
			arg := ext.Method.Args
			str := string(arg)
			slice := str[2:]

			fmt.Println("Data retrieved:")
			if slice == data {
				fmt.Println("Data found in block")
				return nil
			}
		}
	}
	return fmt.Errorf("data not found")
}
