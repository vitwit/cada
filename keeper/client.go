package keeper

import (
	"os"

	cometrpc "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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

func NewClientCtx(kr keyring.Keyring, c *cometrpc.HTTP, ctx sdk.Context, cdc codec.BinaryCodec) client.Context {
	// encodingConfig := params.MakeEncodingConfig()
	// authtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	// cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	// sdk.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	// staking.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	// cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	chainID := ctx.ChainID()

	// fmt.Println("address heree......", address)
	// fromAddress := sdk.AccAddress(address)
	// Assuming you have access to the keyring and broadcast mode
	broadcastMode := "block"

	homepath := "/home/vitwit/.availsdk/keyring-test"

	return client.Context{}.
		WithCodec(cdc.(codec.Codec)).
		WithChainID(chainID).
		// WithFromAddress(fromAddress).
		WithFromName("alice").
		WithKeyringDir(homepath).
		WithBroadcastMode(broadcastMode).
		WithTxConfig(authTx.NewTxConfig(cdc.(codec.Codec), authTx.DefaultSignModes)).
		WithKeyring(kr).
		WithAccountRetriever(authtypes.AccountRetriever{})
}

// NewFactory creates a new Factory.
func NewFactory(clientCtx client.Context) tx.Factory {
	return tx.Factory{}.
		WithChainID(clientCtx.ChainID).
		WithKeybase(clientCtx.Keyring).
		// WithGas(defaultGasLimit).
		// WithGasAdjustment(defaultGasAdjustment).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithTxConfig(clientCtx.TxConfig)
}
