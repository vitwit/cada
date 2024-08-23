package keeper

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	clientCtx := client.Context{}.
		WithCodec(cdc.(codec.Codec)).
		WithChainID(chainID).
		WithFromAddress(fromAddress).
		WithFromName("alice").
		WithKeyringDir("~/.availsdk/keyring-test").
		WithBroadcastMode(broadcastMode)

	msg.ValidatorAddress = fromAddress.String()

	flags := pflag.NewFlagSet("my-flags", pflag.ContinueOnError)
	// Set any additional flags here, like fees, gas, etc.

	err := tx.GenerateOrBroadcastTxCLI(clientCtx, flags, &msg)
	if err != nil {
		return err
	}

	// handle the response, log, or return as needed
	return nil
}
