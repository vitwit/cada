package keeper

import (
	"fmt"
	"os"

	cometrpc "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	clitx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/spf13/pflag"
	"github.com/vitwit/avail-da-module/types"
)

func (k Keeper) SubmitBlobTx(ctx sdk.Context, msg types.MsgSubmitBlobRequest) error {
	// address := k.proposerAddress
	cdc := k.cdc
	homepath := "/home/vitwit/.availsdk/keyring-test"
	// keyring
	kr, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest,
		homepath, os.Stdin, cdc.(codec.Codec))

	if err != nil {
		fmt.Println("error while creating keyring..", err)
		return err
	}

	rpcClient, err := cometrpc.NewWithTimeout("http://localhost:26657", "/websocket", uint(3))
	if err != nil {
		return err
	}

	// create new client context
	clientCtx := NewClientCtx(kr, rpcClient, ctx.ChainID(), cdc)

	flags := *pflag.NewFlagSet("my-flags", pflag.ContinueOnError)
	fmt.Println("new flagssssss.......", flags)

	msg.ValidatorAddress = "cosmos1ux2hl3y42nz6vtdl8k7t7f05k9p3r2k62zfvtv"

	// txf, err := clitx.NewFactoryCLI(clientCtx, &flags)
	// fmt.Println("here the eroor with txf....", txf, err)
	// if err != nil {
	// 	return err
	// }

	factory := NewFactory(clientCtx)

	err = clitx.GenerateOrBroadcastTxWithFactory(clientCtx, factory, &msg)
	if err != nil {
		fmt.Println("error insideeeeeeeeeeee............", err)
		return err
	}

	return nil
}

func (k Keeper) SubmitBlobTx1(ctx sdk.Context, msg types.MsgSubmitBlobRequest) error {
	cdc := k.cdc
	chainID := ctx.ChainID()
	address := k.proposerAddress

	fmt.Println("address heree......", address)
	// fromAddress := sdk.AccAddress(address)
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

	rpcClient, err := cometrpc.NewWithTimeout("http://localhost:26657", "/websocket", uint(3))
	if err != nil {
		return err
	}

	addr, err := sdk.AccAddressFromBech32("cosmos1ux2hl3y42nz6vtdl8k7t7f05k9p3r2k62zfvtv")
	fmt.Println("address and errorr......", addr, err)

	// clientCtx := NewClientCtx(kr, rpcClient)

	// clientCtx := client.Context{}.
	// 	WithCodec(cdc).
	// 	WithTxConfig(authTx.NewTxConfig(cdc, authTx.DefaultSignModes)).
	// 	WithChainID("demo").
	// 	WithKeyringDir("~/.availsdk/keyring-test").
	// 	WithHomeDir("~/.availsdk").
	// 	WithInput(os.Stdin)

	// k.keyring.Backend()

	clientCtx := client.Context{}.
		WithCodec(cdc.(codec.Codec)).
		WithChainID(chainID).
		WithFromAddress(addr).
		WithFromName("alice").
		// WithKeyringDir("~/.availsdk/keyring-test").
		WithKeyringDir(homepath).
		WithBroadcastMode(broadcastMode).
		WithTxConfig(authTx.NewTxConfig(cdc.(codec.Codec), authTx.DefaultSignModes)).
		WithKeyring(kr).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithClient(rpcClient).
		WithSkipConfirmation(true)

	// a, b, c := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, addr)
	// fmt.Println("coming upto hereeeee.........", a, b, c)
	msg.ValidatorAddress = "cosmos1ux2hl3y42nz6vtdl8k7t7f05k9p3r2k62zfvtv"

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

	// err = clitx.GenerateOrBroadcastTxCLI(clientCtx, &flags, &msg)
	txf, err := clitx.NewFactoryCLI(clientCtx, &flags)
	fmt.Println("here the eroor with txf....", txf, err)
	if err != nil {
		return err
	}

	err = clitx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, &msg)
	if err != nil {
		fmt.Println("error insideeeeeeeeeeee............", err)
		return err
	}

	fmt.Println("heree.....")

	// handle the response, log, or return as needed
	return nil
}
