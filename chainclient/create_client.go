package client

import (
	"bytes"
	"fmt"
	"os"

	cometrpc "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// "github.com/emerishq/demeris-backend-models/cns"

const (
	StagingEnvKey      = "staging"
	AkashMnemonicKey   = "AKASH_MNEMONIC"
	CosmosMnemonicKey  = "COSMOS_MNEMONIC"
	TerraMnemonicKey   = "TERRA_MNEMONIC"
	OsmosisMnemonicKey = "OSMOSIS_MNEMONIC"
)

func CreateChainClient(keyringServiceName, chainID, homePath string, codec codec.Codec) (*ChainClient, error) {
	nodeAddress := "http://localhost:26657"
	kr, err := keyring.New(keyringServiceName, KeyringBackendTest, homePath, os.Stdin, codec)
	if err != nil {
		return nil, err
	}

	wsClient, err := cometrpc.New(nodeAddress, "/websocket")
	if err != nil {
		return nil, err
	}
	out := &bytes.Buffer{}

	address := "cosmos1ux2hl3y42nz6vtdl8k7t7f05k9p3r2k62zfvtv"
	clientCtx := NewClientCtx(kr, wsClient, chainID, codec, out, address).WithChainID(chainID).WithNodeURI(nodeAddress)
	fmt.Println("client ctxxx.......", clientCtx.FromName, clientCtx.FromAddress)

	factory := NewFactory(clientCtx)
	return &ChainClient{
		factory:   factory,
		clientCtx: clientCtx,
		out:       out,
	}, nil
}

// GetClient is to create client and imports mnemonic and returns created chain client
func GetClient(chainID string, cc ChainClient, homePath string, codec codec.Codec) (c *ChainClient, err error) {
	// get chain info
	// info, err := LoadSingleChainInfo(env, chainName)
	// if err != nil {
	// 	return nil, err
	// }

	// initSDKConfig(info.NodeInfo.Bech32Config)
	c, err = CreateChainClient(sdk.KeyringServiceName(), chainID, homePath, codec)
	if err != nil {
		return nil, err
	}

	// // mnemonic := cc.Mnemonic
	// // if env == StagingEnvKey {
	// // 	mnemonic = GetMnemonic(chainName)
	// // }

	// // c.AddressPrefix = info.NodeInfo.Bech32Config.PrefixAccount
	// // c.HDPath = info.DerivationPath
	// // c.Enabled = info.Enabled
	// // c.ChainName = info.ChainName
	// c.Mnemonic = ALICE_MNEMONIC
	// c.ChainName = chainID
	// // if len(info.Denoms) != 0 {
	// // 	c.Denom = info.Denoms[0].Name
	// // }

	// fmt.Println("mnemonic and chain id....", c.Mnemonic, c.ChainName, c.Key)

	// err = c.ImportMnemonic(c.Key, c.Mnemonic, c.HDPath)
	// if err != nil {
	// 	return nil, err
	// }

	return c, nil
}
