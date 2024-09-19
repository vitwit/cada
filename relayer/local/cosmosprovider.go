package local

import (
	cometrpc "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/codec"
)

type CosmosProvider struct {
	Cdc       codec.BinaryCodec
	RPCClient RPCClient
}

// NewProvider validates the CosmosProviderConfig, instantiates a ChainClient and then instantiates a CosmosProvider
func NewProvider(cdc codec.BinaryCodec, rpc string) (*CosmosProvider, error) {
	rpcClient, err := cometrpc.NewWithTimeout(rpc, "/websocket", uint(3))
	if err != nil {
		return nil, err
	}

	cp := &CosmosProvider{
		Cdc:       cdc,
		RPCClient: rpcClient,
	}

	return cp, nil
}
