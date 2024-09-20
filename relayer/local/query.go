package local

import (
	"context"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/bytes"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RPCClient defines a set of methods for interacting with a blockchain node via RPC.
// It includes methods for retrieving block information, querying node status, and
// performing ABCI queries.
type RPCClient interface {
	Block(ctx context.Context, height *int64) (*coretypes.ResultBlock, error)
	Status(ctx context.Context) (*coretypes.ResultStatus, error)
	ABCIQuery(ctx context.Context, path string, data bytes.HexBytes) (*coretypes.ResultABCIQuery, error)
}

// GetBlockAtHeight retrieves the block information at a specific height from the chain
// using the RPC client.
func (cc *CosmosProvider) GetBlockAtHeight(ctx context.Context, height int64) (*coretypes.ResultBlock, error) {
	block, err := cc.RPCClient.Block(ctx, &height)
	if err != nil {
		return nil, fmt.Errorf("error querying block at height %d: %w", height, err)
	}
	return block, nil
}

// Status queries the status of this node, can be used to check if it is catching up or a validator
func (cc *CosmosProvider) Status(ctx context.Context) (*coretypes.ResultStatus, error) {
	status, err := cc.RPCClient.Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("error querying status: %w", err)
	}
	return status, nil
}

// QueryABCI performs an ABCI query and returns the appropriate response and error sdk error code.
func (cc *CosmosProvider) QueryABCI(ctx context.Context, path string, data []byte) (abci.ResponseQuery, error) {
	result, err := cc.RPCClient.ABCIQuery(ctx, path, data)
	if err != nil {
		return abci.ResponseQuery{}, err
	}

	if !result.Response.IsOK() {
		return abci.ResponseQuery{}, sdkErrorToGRPCError(result.Response)
	}

	return result.Response, nil
}

func sdkErrorToGRPCError(resp abci.ResponseQuery) error {
	switch resp.Code {
	case sdkerrors.ErrInvalidRequest.ABCICode():
		return status.Error(codes.InvalidArgument, resp.Log)
	case sdkerrors.ErrUnauthorized.ABCICode():
		return status.Error(codes.Unauthenticated, resp.Log)
	case sdkerrors.ErrKeyNotFound.ABCICode():
		return status.Error(codes.NotFound, resp.Log)
	default:
		return status.Error(codes.Unknown, resp.Log)
	}
}
