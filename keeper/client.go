package keeper

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
)

func InitClientCtx(cdc *codec.ProtoCodec) client.Context {
	// interfaceRegistry := codectypes.NewInterfaceRegistry()

	// // Register the custom type InjectedData
	// availblobTypes.RegisterInterfaces(interfaceRegistry)

	// cdc := codec.NewProtoCodec(interfaceRegistry)
	// cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

	clientCtx := client.Context{}.
		WithCodec(cdc).
		WithTxConfig(authTx.NewTxConfig(cdc, authTx.DefaultSignModes)).
		WithChainID("demo").
		WithKeyringDir("~/.availsdk/keyring-test").
		WithHomeDir("~/.availsdk").
		WithInput(os.Stdin)

	return clientCtx
}
