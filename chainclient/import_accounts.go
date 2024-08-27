package client

import (
	"bytes"
	"fmt"
	"io"

	cometrpc "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"

	// "github.com/tendermint/starport/starport/pkg/xfilepath"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/go-bip39"
	// "github.com/tendermint/starport/starport/pkg/spn"
)

// var availdHomePath = xfilepath.JoinFromHome(xfilepath.Path("availsdk"))

const (
	KeyringBackendTest = "test"
	ALICE_MNEMONIC     = ""
	// ALICE_MNEMONIC     = "all soap kiwi cushion federal skirt tip shock exist tragic verify lunar shine rely torch please view future lizard garbage humble medal leisure mimic"
)

// ChainClient is client to interact with SPN.
type ChainClient struct {
	factory            tx.Factory
	clientCtx          client.Context
	out                *bytes.Buffer
	Address            string `json:"address"`
	AddressPrefix      string `json:"account_address_prefix"`
	RPC                string `json:"rpc"`
	Key                string `json:"key"`
	Mnemonic           string `json:"mnemonic"`
	KeyringServiceName string `json:"keyring_service_name"`
	HDPath             string `json:"hd_path"`
	Enabled            bool   `json:"enabled"`
	ChainName          string `json:"chain_name"`
	Denom              string `json:"denom"`
}

// ImportMnemonic is to import existing account mnemonic in keyring
func (c ChainClient) ImportMnemonic(keyName, mnemonic, hdPath string) (err error) {
	err = c.AccountCreate(keyName, mnemonic, hdPath) // return account also
	fmt.Println("here the accc details.......", keyName, mnemonic, hdPath)
	if err != nil {
		return err
	}

	return nil
}

// AccountCreate creates an account by name and mnemonic (optional) in the keyring.
func (c *ChainClient) AccountCreate(accountName, mnemonic, hdPath string) error {
	if mnemonic == "" {
		entropySeed, err := bip39.NewEntropy(256)
		if err != nil {
			return err
		}
		mnemonic, err = bip39.NewMnemonic(entropySeed)
		fmt.Println("mnemoniccccc here.....", mnemonic)
		if err != nil {
			return err
		}
	}

	algos, _ := c.clientCtx.Keyring.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), algos)
	if err != nil {
		return err
	}

	info, err := c.clientCtx.Keyring.NewAccount(accountName, mnemonic, "", hdPath, algo)
	if err != nil {
		return err
	}
	pk, err := info.GetPubKey()
	if err != nil {
		return err
	}
	addr := sdk.AccAddress(pk.Address())
	fmt.Println("address hereee...", addr)
	// account := c.ToAccount(info)
	// account.Mnemonic = mnemonic
	return nil
}

// func initSDKConfig() {
// 	// sdkConfig := sdk.GetConfig()
// 	// bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()
// 	// sdkConfig.SetBech32PrefixForAccount(config.Bech32PrefixAccAddr(), config.Bech32PrefixAccPub())
// 	// sdkConfig.SetBech32PrefixForValidator(config.Bech32PrefixValAddr(), config.Bech32PrefixValPub())
// 	// sdkConfig.SetBech32PrefixForConsensusNode(config.Bech32PrefixConsAddr(), config.Bech32PrefixConsPub())
// }

func NewClientCtx(kr keyring.Keyring, c *cometrpc.HTTP, chainID string, cdc codec.BinaryCodec, out io.Writer, address string) client.Context {
	encodingConfig := MakeEncodingConfig()
	authtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	sdk.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	staking.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	// chainID := ctx.ChainID()

	// fmt.Println("address heree......", address)
	fromAddress := sdk.AccAddress(address)
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
		WithOutput(out).WithClient(c).WithInterfaceRegistry(encodingConfig.InterfaceRegistry)
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

// AccountList returns a list of accounts.
// func (c *ChainClient) AccountList() (accounts []sdk.Account, err error) {
// 	infos, err := c.clientCtx.Keyring.List()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, info := range infos {
// 		accounts = append(accounts, c.ToAccount(info))
// 	}
// 	return accounts, nil
// }
