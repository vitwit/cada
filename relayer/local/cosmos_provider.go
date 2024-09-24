package local

import (
	"context"

	cometrpc "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

// CosmosProvider defines the methods required for interacting with a Cosmos blockchain.
type CosmosProvider interface {
	QueryABCI(ctx context.Context, path string, data []byte) ([]byte, error)
	GetCodec() codec.BinaryCodec
	GetRPCClient() RPCClient
	GetBlockAtHeight(ctx context.Context, height int64) (*coretypes.ResultBlock, error)
}
type DefaultCosmosProvider struct {
	Cdc       codec.BinaryCodec
	RPCClient RPCClient
}

// QueryABCI queries the ABCI for the given path and data.
func (cp *DefaultCosmosProvider) QueryABCI(ctx context.Context, path string, data []byte) ([]byte, error) {
	res, err := cp.RPCClient.ABCIQuery(ctx, path, data)
	if err != nil {
		return nil, err
	}
	return res.Response.Value, nil
}

// GetCodec returns the codec used by the CosmosProvider.
func (cp *DefaultCosmosProvider) GetCodec() codec.BinaryCodec {
	return cp.Cdc
}

// GetRPCClient returns the RPC client used by the CosmosProvider.
func (cp *DefaultCosmosProvider) GetRPCClient() RPCClient {
	return cp.RPCClient
}

// NewProvider validates the CosmosProviderConfig, instantiates a ChainClient and then instantiates a CosmosProvider
func NewDefaultCosmosProvider(cdc codec.BinaryCodec, rpc string) (CosmosProvider, error) {
	rpcClient, err := cometrpc.NewWithTimeout(rpc, "/websocket", uint(3))
	if err != nil {
		return nil, err
	}

	cp := &DefaultCosmosProvider{
		Cdc:       cdc,
		RPCClient: rpcClient,
	}

	return cp, nil
}
