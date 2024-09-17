package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (qs queryServer) SubmittedBlobStatus(ctx context.Context, _ *types.QuerySubmittedBlobStatusRequest) (*types.QuerySubmittedBlobStatusResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(qs.k.storeKey)
	startHeight := qs.k.GetStartHeightFromStore(sdkCtx)
	endHeight := qs.k.GetEndHeightFromStore(sdkCtx)
	status := GetStatusFromStore(store)
	blobStatus := ParseStatus(status)
	provenHeight := qs.k.GetProvenHeightFromStore(sdkCtx)
	votingEndHeight := qs.k.GetVotingEndHeightFromStore(sdkCtx)

	return &types.QuerySubmittedBlobStatusResponse{
		Range:                &types.Range{From: startHeight, To: endHeight},
		Status:               blobStatus,
		ProvenHeight:         provenHeight,
		LastBlobVotingEndsAt: votingEndHeight,
	}, nil
}
