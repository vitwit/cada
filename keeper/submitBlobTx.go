package keeper

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	clitx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/spf13/pflag"
	"github.com/vitwit/avail-da-module/types"
)

func (k Keeper) SubmitBlobTx(ctx sdk.Context, msg types.MsgSubmitBlobRequest) error {
	cdc := k.cdc
	chainID := ctx.ChainID()
	address := k.proposerAddress
	fromAddress := sdk.AccAddress(address)
	// Assuming you have access to the keyring and broadcast mode
	broadcastMode := "block"

	homepath := "/home/vitwit/.availsdk/keyring-test"

	// keyring
	kr, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest,
		homepath, os.Stdin, cdc.(codec.Codec))

	if err != nil {
		fmt.Println("error while creating keyring..", err)
		return err
	}

	/*
		clientCtx := client.Context{}.
		WithCodec(cdc).
		WithTxConfig(authTx.NewTxConfig(cdc, authTx.DefaultSignModes)).
		WithChainID("demo").
		WithKeyringDir("~/.availsdk/keyring-test").
		WithHomeDir("~/.availsdk").
		WithInput(os.Stdin)
	*/

	// k.keyring.Backend()

	clientCtx := client.Context{}.
		WithCodec(cdc.(codec.Codec)).
		WithChainID(chainID).
		WithFromAddress(fromAddress).
		WithFromName("alice").
		// WithKeyringDir("~/.availsdk/keyring-test").
		WithBroadcastMode(broadcastMode).
		WithTxConfig(authTx.NewTxConfig(cdc.(codec.Codec), authTx.DefaultSignModes)).
		WithKeyring(kr)

	fmt.Println("coming upto hereeeee.........")
	msg.ValidatorAddress = fromAddress.String()

	fmt.Println("validator addressssssss............, ", msg.ValidatorAddress)

	flags := *pflag.NewFlagSet("my-flags", pflag.ContinueOnError)
	fmt.Println("new flagssssss.......", flags)
	// fmt.Println("account and sequence numberrr.......", flags.)
	// Set any additional flags here, like fees, gas, etc.

	fmt.Println("txxxxxxxxxxx........", clientCtx.ChainID)
	fmt.Println("txxxxxxxxxxx........", clientCtx.CmdContext)
	fmt.Println("txxxxxxxxxxx........", clientCtx.Codec)
	fmt.Println("txxxxxxxxxxx........", clientCtx.FromAddress.String())
	fmt.Println("txxxxxxxxxxx........", clientCtx.BroadcastMode)
	fmt.Println("aaaaaaaa.......", clientCtx.TxConfig)
	fmt.Println("aaaaaaaa.......", clientCtx.AccountRetriever)

	err = clitx.GenerateOrBroadcastTxCLI(clientCtx, &flags, &msg)
	if err != nil {
		fmt.Println("error insideeeeeeeeeeee............", err)
		return err
	}

	fmt.Println("heree.....")

	// handle the response, log, or return as needed
	return nil
}
