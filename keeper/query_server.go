package keeper

import (
	"context"
	"time"

	"cosmossdk.io/collections"
	"github.com/vitwit/avail-da-module/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k *Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k *Keeper
}

// PendingValidators returns the pending validators.
func (qs queryServer) Validators(ctx context.Context, _ *types.QueryValidatorsRequest) (*types.QueryValidatorsResponse, error) {
	vals, err := qs.k.GetAllValidators(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryValidatorsResponse{Validators: vals.Validators}, nil
}

func (qs queryServer) AvailAddress(ctx context.Context, req *types.QueryAvailAddressRequest) (*types.QueryAvailAddressResponse, error) {
	addr, err := qs.k.GetValidatorAvailAddress(ctx, req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	return &types.QueryAvailAddressResponse{AvailAddress: addr}, nil
}

func (qs queryServer) ProvenHeight(ctx context.Context, _ *types.QueryProvenHeightRequest) (*types.QueryProvenHeightResponse, error) {
	provenHeight, err := qs.k.GetProvenHeight(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryProvenHeightResponse{
		ProvenHeight: provenHeight,
	}, nil
}

func (qs queryServer) PendingBlocks(ctx context.Context, _ *types.QueryPendingBlocksRequest) (*types.QueryPendingBlocksResponse, error) {
	pendingBlocks, err := qs.k.GetPendingBlocksWithExpiration(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryPendingBlocksResponse{
		PendingBlocks: pendingBlocks,
	}, nil
}

func (qs queryServer) ExpiredBlocks(ctx context.Context, _ *types.QueryExpiredBlocksRequest) (*types.QueryExpiredBlocksResponse, error) {
	currentTimeNs := time.Now().UnixNano()
	iterator, err := qs.k.TimeoutsToPendingBlocks.
		Iterate(ctx, (&collections.Range[int64]{}).StartInclusive(0).EndInclusive(currentTimeNs))
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var expiredBlocks []*types.BlockWithExpiration
	for ; iterator.Valid(); iterator.Next() {
		expiration, err := iterator.Key()
		if err != nil {
			return nil, err
		}
		blocks, err := iterator.Value()
		if err != nil {
			return nil, err
		}
		for _, block := range blocks.BlockHeights {
			expiredBlocks = append(expiredBlocks, &types.BlockWithExpiration{
				Height:     block,
				Expiration: time.Unix(0, expiration),
			})
		}
	}
	return &types.QueryExpiredBlocksResponse{
		CurrentTime:   time.Unix(0, currentTimeNs),
		ExpiredBlocks: expiredBlocks,
	}, nil
}
