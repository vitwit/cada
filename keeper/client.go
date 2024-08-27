package keeper

import (
	"os"

	cometrpc "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// "github.com/tendermint/starport/starport/pkg/xfilepath"
)

// var availdHomePath = xfilepath.JoinFromHome(xfilepath.Path("availsdk"))

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

func NewClientCtx(kr keyring.Keyring, c *cometrpc.HTTP, chainID string, cdc codec.BinaryCodec) client.Context {
	encodingConfig := MakeEncodingConfig()
	// authtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	// cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	// sdk.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	// staking.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	// cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	// chainID := ctx.ChainID()

	// fmt.Println("address heree......", address)
	// sdk.AccAddressFromBech32()
	fromAddress, err := sdk.AccAddressFromBech32("cosmos1ux2hl3y42nz6vtdl8k7t7f05k9p3r2k62zfvtv")
	// fmt.Println("here errorr...", err)
	if err != nil {
		// return err
	}
	// fmt.Println("from addresss.........", fromAddress)
	// Assuming you have access to the keyring and broadcast mode
	broadcastMode := "block"

	homepath := "/home/vitwit/.availsdk/keyring-test"

	return client.Context{}.
		WithCodec(cdc.(codec.Codec)).
		WithChainID(chainID).
		WithFromAddress(fromAddress).
		WithFromName("alice").
		WithKeyringDir(homepath).
		WithBroadcastMode(broadcastMode).
		WithTxConfig(authTx.NewTxConfig(cdc.(codec.Codec), authTx.DefaultSignModes)).
		WithKeyring(kr).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithClient(c).WithInterfaceRegistry(encodingConfig.InterfaceRegistry)
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

// MakeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func MakeEncodingConfig() EncodingConfig {
	aminoCodec := codec.NewLegacyAmino()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)
	txCfg := authTx.NewTxConfig(codec, authTx.DefaultSignModes)

	encCfg := EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          txCfg,
		Amino:             aminoCodec,
	}

	std.RegisterLegacyAminoCodec(encCfg.Amino)
	std.RegisterInterfaces(encCfg.InterfaceRegistry)
	// mb.RegisterLegacyAminoCodec(encCfg.Amino)
	// mb.RegisterInterfaces(encCfg.InterfaceRegistry)

	return encCfg
}

// EncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}
