package relayer

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/retriever"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

func (r *Relayer) SubmitDataToClient(Seed string, AppID int, data []byte, blocks []int64, lightClientUrl string) error {
	if r.submittedBlocksCache[blocks[0]] {
		return nil
	}

	r.submittedBlocksCache[blocks[0]] = true
	delete(r.submittedBlocksCache, blocks[0]-int64(len(blocks)))

	handler := NewHTTPClientHandler()
	datab := base64.StdEncoding.EncodeToString(data)

	jsonData := []byte(fmt.Sprintf(`{"data":"%s"}`, datab))

	// Define the URL
	//url := "http://127.0.0.1:8000/v2/submit"

	url := fmt.Sprintf("%s/v2/submit", lightClientUrl)

	// Make the POST request
	responseBody, err := handler.Post(url, jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	// Create an instance of the struct
	var blockInfo BlockInfo

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(responseBody, &blockInfo)
	if err != nil {
		r.logger.Error("Error while posting block(s) to Avail DA",
			"height_start", blocks[0],
			"height_end", blocks[len(blocks)-1],
			"appID", string(r.rpcClient.config.AppID),
		)
	}

	if err == nil {
		r.logger.Info("Successfully posted block(s) to Avail DA",
			"height_start", blocks[0],
			"height_end", blocks[len(blocks)-1],
			"appID", string(r.rpcClient.config.AppID),
			"block_hash", blockInfo.BlockHash,
			"block_number", blockInfo.BlockNumber,
			"hash", blockInfo.Hash,
		)
	}

	return nil
}

func (r *Relayer) GetSubmittedData(lightClientUrl string, blockNumber int) {
	handler := NewHTTPClientHandler()
	// get submitted block data using light client api with avail block height
	// time.Sleep(20 * time.Second) // wait upto data to be submitted data to be included in the avail blocks

	url := fmt.Sprintf("%s/v2/blocks/%v/data?feilds=data", lightClientUrl)
	url = fmt.Sprintf(url, blockNumber)

	body, err := handler.Get(url)
	if err != nil {
		return
	}

	if body != nil {
		r.logger.Info("submitted data to Avail verfied successfully at",
			"block_height", blockNumber,
		)
	}
}

// Define the struct that matches the JSON structure
type GetBlock struct {
	BlockNumber      int      `json:"block_number"`
	DataTransactions []string `json:"data_transactions"`
}

// submitData creates a transaction and makes a Avail data submission
func (r *Relayer) SubmitData1(ApiURL string, Seed string, AppID int, data []byte, blocks []int64) error {
	api, err := gsrpc.NewSubstrateAPI(ApiURL)
	if err != nil {
		r.logger.Error("cannot create api:%w", err)
		return err
	}

	fmt.Println("appurl", ApiURL, "seed", Seed)
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
	if AppID != 0 {
		appID = AppID
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

	keyringPair, err := signature.KeyringPairFromSecret(Seed, 42)
	if err != nil {
		r.logger.Error("cannot create KeyPair:%w", err)
		return err
	}

	// fmt.Println("keyring pair", keyringPair)

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
			"appID", string(r.rpcClient.config.AppID),
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
							func(fieldIndex int, field *registry.DecodedField) bool {

								return field.Name == "sp_core.crypto.AccountId32.who"
							},
							func(value any) (*types.AccountID, error) {
								fields, ok := value.(registry.DecodedFields)

								if !ok {
									return nil, fmt.Errorf("unexpected value: %v", value)
								}

								accByteSlice, err := registry.GetDecodedFieldAsSliceOfType[types.U8](fields, func(fieldIndex int, field *registry.DecodedField) bool {
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
							func(fieldIndex int, field *registry.DecodedField) bool {
								return fieldIndex == 1
							},
							func(value any) (*types.Hash, error) {
								fields, ok := value.(registry.DecodedFields)
								if !ok {
									return nil, fmt.Errorf("unexpected value: %v", value)
								}

								hashByteSlice, err := registry.GetDecodedFieldAsSliceOfType[types.U8](fields, func(fieldIndex int, field *registry.DecodedField) bool {
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
						if err != nil {
							fmt.Printf("DataHash parsing err: %s\n", err.Error())
						} else if dataHash == nil {
							fmt.Println("DataHash is nil")
						} else {
							fmt.Printf("DataHash read from event: %s \n", dataHash.Hex())
						}
						fmt.Printf("DataHash read from event: %s \n", dataHash.Hex())

					}

				}
				fmt.Printf("Txn inside finalized block\n")
				hash := status.AsFinalized
				err = getData1(hash, api, string(data))
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

// RandToken generates a random hex value.
func RandToken1(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getData1(hash types.Hash, api *gsrpc.SubstrateAPI, data string) error {

	block, err := api.RPC.Chain.GetBlock(hash)
	if err != nil {
		return fmt.Errorf("cannot get block by hash:%w", err)
	}

	// Encode the struct to JSON
	// jsonData, err := json.Marshal(block.Block.Header)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("length of extrinsics: ", len(block.Block.Extrinsics))
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
