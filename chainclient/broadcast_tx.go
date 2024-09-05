package chainclient

import (
	"fmt"
	"log"
	"os"

	cometrpc "github.com/cometbft/cometbft/rpc/client/http"
	clitx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/types"

	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func GetBinPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	availdHomePath := filepath.Join(homeDir, ".availsdk")
	fmt.Println("availdHonmePath.......", availdHomePath)
	return availdHomePath
}

func ExecuteTX(ctx sdk.Context, msg types.MsgUpdateBlobStatusRequest, cdc codec.BinaryCodec) error {
	// Define keyring and RPC client configuration

	// homePath := "/home/vitwit/.availsdk"
	homePath := GetBinPath()
	key := os.Getenv("KEY")
	fmt.Println("get key namee.........", key)
	if key == "" { //TODO : remove this later
		key = "alice"
	}
	keyName := key
	rpcAddress := "http://localhost:26657"

	// Create a keyring
	kr, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, homePath, os.Stdin, cdc.(codec.Codec))
	if err != nil {
		return fmt.Errorf("error creating keyring: %w", err)
	}

	// List all keys in the keyring
	// keys, err := kr.List()
	// if err != nil {
	// 	fmt.Println("error listing keys:", err)
	// }

	info, err := kr.Key(keyName)
	valAddr, err := info.GetAddress()
	fmt.Println("after address................", valAddr)

	// valAddr, err := sdk.AccAddressFromBech32(addr.String())
	// fmt.Println("val addr, err..", valAddr, err, addr)

	// fmt.Println("keysss........", keys)

	// // Print out the keys
	// for _, keyInfo := range keys {
	// 	addr, err := keyInfo.GetAddress()
	// 	fmt.Println("err..", err)
	// 	fmt.Printf("Name: %s, Address: %s\n", keyInfo.Name, addr)
	// }

	// Create an RPC client
	rpcClient, err := cometrpc.NewWithTimeout(rpcAddress, "/websocket", 3)
	if err != nil {
		return fmt.Errorf("error creating RPC client: %w", err)
	}

	// Create a new client context
	clientCtx := NewClientCtx(kr, rpcClient, ctx.ChainID(), cdc, homePath, valAddr)

	// Retrieve the validator address (replace with actual logic to get the address)
	// valAddr, err = sdk.AccAddressFromBech32("cosmos1fhqer4tc50nut2evvnj6yegcah2yfu3s844n9a")
	// if err != nil {
	// 	return fmt.Errorf("error parsing validator address: %w", err)
	// }

	// Set the client context's from fields
	clientCtx.FromName = keyName
	clientCtx.FromAddress = valAddr

	// Fetch account number and sequence from the blockchain
	accountRetriever := authtypes.AccountRetriever{}
	account, err := accountRetriever.GetAccount(clientCtx, valAddr)
	if err != nil {
		return fmt.Errorf("error retrieving account: %w", err)
	}

	fmt.Println("account details......", account.GetAccountNumber(), account.GetSequence())

	// Set the correct account number and sequence
	factory := NewFactory(clientCtx).
		WithAccountNumber(account.GetAccountNumber()).
		WithSequence(account.GetSequence())

	// Create a transaction factory and set the validator address in the message
	// factory := NewFactory(clientCtx)
	msg.ValidatorAddress = valAddr.String()
	// time.Sleep(10 * time.Second)

	// Generate and broadcast the transaction
	if err := clitx.GenerateOrBroadcastTxWithFactory(clientCtx, factory, &msg); err != nil {
		return fmt.Errorf("error broadcasting transaction: %w", err)
	}

	return nil
}
