package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dacli "github.com/vitwit/avail-da-module/chainclient"
	"github.com/vitwit/avail-da-module/types"
)

// SubmitBlobTx2 submits a blob transaction to the chain using a chain client.
// This function creates a chain client, sets the validator address from the client, and broadcasts the transaction.
func (k Keeper) SubmitBlobTx2(ctx sdk.Context, msg types.MsgSubmitBlobRequest) error {
	cdc := k.cdc
	homepath := "/home/vitwit/.availsdk" // TODO: get from config

	cc, err := dacli.CreateChainClient(sdk.KeyringServiceName(), ctx.ChainID(), homepath, cdc.(codec.Codec))
	if err != nil {
		return err
	}

	msg.ValidatorAddress = cc.Address
	err = cc.BroadcastTx(msg, cc.Key, sdk.AccAddress(cc.Address))
	if err != nil {
		fmt.Println("error while broadcastig the txxx.........", err)
		return err
	}

	return nil
}

// func SubmitBlobTx(ctx sdk.Context, msg types.MsgSubmitBlobRequest, cdc codec.Codec) error {
// 	// Define keyring and RPC client configuration

// 	homePath := "/home/vitwit/.availsdk"
// 	keyName := "alice"
// 	rpcAddress := "http://localhost:26657"

// 	// Create a keyring
// 	kr, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, homePath, os.Stdin, cdc)
// 	if err != nil {
// 		return fmt.Errorf("error creating keyring: %w", err)
// 	}

// 	// List all keys in the keyring
// 	// keys, err := kr.List()
// 	// if err != nil {
// 	// 	fmt.Println("error listing keys:", err)
// 	// }

// 	info, err := kr.Key(keyName)
// 	// log.Println("uuu....", info, err)

// 	valAddr, err := info.GetAddress()

// 	// valAddr, err := sdk.AccAddressFromBech32(addr.String())
// 	// fmt.Println("val addr, err..", valAddr, err, addr)

// 	// fmt.Println("keysss........", keys)

// 	// // Print out the keys
// 	// for _, keyInfo := range keys {
// 	// 	addr, err := keyInfo.GetAddress()
// 	// 	fmt.Println("err..", err)
// 	// 	fmt.Printf("Name: %s, Address: %s\n", keyInfo.Name, addr)
// 	// }

// 	// Create an RPC client
// 	rpcClient, err := cometrpc.NewWithTimeout(rpcAddress, "/websocket", 3)
// 	if err != nil {
// 		return fmt.Errorf("error creating RPC client: %w", err)
// 	}

// 	// Create a new client context
// 	clientCtx := NewClientCtx(kr, rpcClient, ctx.ChainID(), cdc, homePath, valAddr)

// 	// Retrieve the validator address (replace with actual logic to get the address)
// 	// valAddr, err = sdk.AccAddressFromBech32("cosmos1fhqer4tc50nut2evvnj6yegcah2yfu3s844n9a")
// 	// if err != nil {
// 	// 	return fmt.Errorf("error parsing validator address: %w", err)
// 	// }

// 	// Set the client context's from fields
// 	clientCtx.FromName = keyName
// 	clientCtx.FromAddress = valAddr

// 	// Fetch account number and sequence from the blockchain
// 	accountRetriever := authtypes.AccountRetriever{}
// 	account, err := accountRetriever.GetAccount(clientCtx, valAddr)
// 	if err != nil {
// 		return fmt.Errorf("error retrieving account: %w", err)
// 	}

// 	fmt.Println("account details......", account.GetAccountNumber(), account.GetSequence())

// 	// Set the correct account number and sequence
// 	factory := NewFactory(clientCtx).
// 		WithAccountNumber(account.GetAccountNumber()).
// 		WithSequence(account.GetSequence())

// 	// Create a transaction factory and set the validator address in the message
// 	// factory := NewFactory(clientCtx)
// 	msg.ValidatorAddress = valAddr.String()
// 	// time.Sleep(10 * time.Second)

// 	// Generate and broadcast the transaction
// 	if err := clitx.GenerateOrBroadcastTxWithFactory(clientCtx, factory, &msg); err != nil {
// 		return fmt.Errorf("error broadcasting transaction: %w", err)
// 	}

// 	return nil
// }
