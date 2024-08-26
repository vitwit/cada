package client

import (
	"bytes"
	"os"

	cometrpc "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

// "github.com/emerishq/demeris-backend-models/cns"

const (
	StagingEnvKey      = "staging"
	AkashMnemonicKey   = "AKASH_MNEMONIC"
	CosmosMnemonicKey  = "COSMOS_MNEMONIC"
	TerraMnemonicKey   = "TERRA_MNEMONIC"
	OsmosisMnemonicKey = "OSMOSIS_MNEMONIC"
)

func CreateChainClient(nodeAddress, keyringServiceName, chainID, homePath string, codec codec.Codec) (*ChainClient, error) {
	nodeAddress = "http://localhost:26657"
	kr, err := keyring.New(keyringServiceName, KeyringBackendTest, homePath, os.Stdin, codec)
	if err != nil {
		return nil, err
	}

	wsClient, err := cometrpc.New(nodeAddress, "/websocket")
	if err != nil {
		return nil, err
	}
	out := &bytes.Buffer{}
	clientCtx := NewClientCtx(kr, wsClient, chainID, codec, out).WithChainID(chainID).WithNodeURI(nodeAddress)

	factory := NewFactory(clientCtx)
	return &ChainClient{
		factory:   factory,
		clientCtx: clientCtx,
		out:       out,
	}, nil
}

// GetClient is to create client and imports mnemonic and returns created chain client
// func GetClient(env string, chainName string, cc ChainClient, dir string) (c *ChainClient, err error) {
// 	// get chain info
// 	// info, err := LoadSingleChainInfo(env, chainName)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// initSDKConfig(info.NodeInfo.Bech32Config)
// 	c, err = CreateChainClient(cc.RPC, cc.KeyringServiceName, info.NodeInfo.ChainID, dir)
// 	if err != nil {
// 		return nil, err
// 	}

// 	mnemonic := cc.Mnemonic
// 	if env == StagingEnvKey {
// 		mnemonic = GetMnemonic(chainName)
// 	}

// 	c.AddressPrefix = info.NodeInfo.Bech32Config.PrefixAccount
// 	c.HDPath = info.DerivationPath
// 	c.Enabled = info.Enabled
// 	c.ChainName = info.ChainName
// 	c.Mnemonic = mnemonic
// 	c.ChainName = chainName
// 	if len(info.Denoms) != 0 {
// 		c.Denom = info.Denoms[0].Name
// 	}

// 	_, err = c.ImportMnemonic(cc.Key, c.Mnemonic, c.HDPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return c, nil
// }
