package relayer

// import (
// 	"crypto/rand"
// 	"encoding/hex"
// 	"fmt"
// 	"time"

// 	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
// 	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
// 	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
// )

// // submitData creates a transaction and makes a Avail data submission
// func (client *AvailClient) SubmitData(ApiURL string, Seed string, AppID int, data []byte) (types.Hash, error) {
// 	var hash types.Hash
// 	api, err := gsrpc.NewSubstrateAPI(ApiURL)
// 	if err != nil {
// 		return hash, fmt.Errorf("cannot create api:%w", err)
// 	}

// 	fmt.Println("appurl", ApiURL, "seed", Seed)
// 	fmt.Println("created substrate api")

// 	meta, err := api.RPC.State.GetMetadataLatest()
// 	if err != nil {
// 		return hash, fmt.Errorf("cannot get metadata:%w", err)
// 	}

// 	fmt.Println("meta data")

// 	// Set data and appID according to need
// 	// data, _ := RandToken(size)
// 	appID := 0

// 	fmt.Println("randome data")

// 	// if app id is greater than 0 then it must be created before submitting data
// 	if AppID != 0 {
// 		appID = AppID
// 	}

// 	c, err := types.NewCall(meta, "DataAvailability.submit_data", types.NewBytes(data))
// 	if err != nil {
// 		return hash, fmt.Errorf("cannot create new call:%w", err)
// 	}

// 	fmt.Println("made call")
// 	// Create the extrinsic
// 	ext := types.NewExtrinsic(c)

// 	fmt.Println("extrinsics")
// 	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
// 	if err != nil {
// 		return hash, fmt.Errorf("cannot get block hash:%w", err)
// 	}

// 	fmt.Println("genesisHash", genesisHash)

// 	rv, err := api.RPC.State.GetRuntimeVersionLatest()
// 	if err != nil {
// 		return hash, fmt.Errorf("cannot get runtime version:%w", err)
// 	}

// 	// fmt.Println("rv", rv)

// 	keyringPair, err := signature.KeyringPairFromSecret(Seed, 42)
// 	if err != nil {
// 		return hash, fmt.Errorf("cannot create KeyPair:%w", err)
// 	}

// 	// fmt.Println("keyring pair", keyringPair)

// 	key, err := types.CreateStorageKey(meta, "System", "Account", keyringPair.PublicKey)
// 	if err != nil {
// 		return hash, fmt.Errorf("cannot create storage key:%w", err)
// 	}

// 	fmt.Println("keyyyy", key)
// 	var accountInfo types.AccountInfo
// 	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
// 	fmt.Println("ok", ok, err)
// 	if err != nil || !ok {
// 		return hash, fmt.Errorf("cannot get latest storage:%w", err)
// 	}

// 	nonce := uint32(accountInfo.Nonce)

// 	// fmt.Println("storage", ok, accountInfo, nonce)

// 	o := types.SignatureOptions{
// 		BlockHash:          genesisHash,
// 		Era:                types.ExtrinsicEra{IsMortalEra: false},
// 		GenesisHash:        genesisHash,
// 		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
// 		SpecVersion:        rv.SpecVersion,
// 		Tip:                types.NewUCompactFromUInt(0),
// 		AppID:              types.NewUCompactFromUInt(uint64(0)),
// 		TransactionVersion: rv.TransactionVersion,
// 	}

// 	fmt.Println("before signing", err)
// 	// Sign the transaction using Alice's default account
// 	err = ext.Sign(keyringPair, o)

// 	fmt.Println("after signing", err)

// 	// fmt.Println("oooooo", o)
// 	if err != nil {
// 		return hash, fmt.Errorf("cannot sign:%w", err)
// 	}

// 	fmt.Println("before sumbmitting extrinsincs")
// 	// Submit the extrinsic
// 	hash, err = api.RPC.Author.SubmitExtrinsic(ext)
// 	if err != nil {
// 		return hash, fmt.Errorf("cannot submit extrinsic:%w", err)
// 	}

// 	fmt.Printf("Data submitted by Alice: against appID %v  sent with hash %#x\n", appID, hash)

// 	go getData(hash, api, string(types.NewBytes(data)))
// 	return hash, err
// }

// // RandToken generates a random hex value.
// func RandToken(n int) (string, error) {
// 	bytes := make([]byte, n)
// 	if _, err := rand.Read(bytes); err != nil {
// 		return "", err
// 	}
// 	return hex.EncodeToString(bytes), nil
// }

// func getData(hash types.Hash, api *gsrpc.SubstrateAPI, data string) error {

// 	time.Sleep(50 * time.Second)
// 	fmt.Println("checking get data.............")

// 	dataBytes := []byte(data)
// 	db := types.NewBytes(dataBytes)

// 	fmt.Printf("sent bytes %s", BytesToHex(db))

// 	block, err := api.RPC.Chain.GetBlock(hash)
// 	if err != nil {
// 		fmt.Println("can't get block data..................................")
// 		return fmt.Errorf("cannot get block by hash:%w", err)
// 	}

// 	fmt.Println("block", block)

// 	fmt.Println("length of extrinsics: ", len(block.Block.Extrinsics))
// 	for _, ext := range block.Block.Extrinsics {

// 		// these values below are specific indexes only for data submission, differs with each extrinsic
// 		if ext.Method.CallIndex.SectionIndex == 29 && ext.Method.CallIndex.MethodIndex == 1 {
// 			arg := ext.Method.Args
// 			str := string(arg)
// 			slice := str[2:]

// 			prefix := slice[0:2]
// 			fmt.Println("prefix........", prefix)

// 			fmt.Println("Data retrieved:")
// 			if slice == data {
// 				fmt.Println("retrieved data", BytesToHex([]byte(str)))
// 				fmt.Println("Data found in block")
// 			}
// 		}
// 	}

// 	fmt.Println("done get data.....................")
// 	return nil
// }

// func BytesToHex(data []byte) string {
// 	return hex.EncodeToString(data)
// }
