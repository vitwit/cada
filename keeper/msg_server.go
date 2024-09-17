package keeper

import (
	"context"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/types"
)

type msgServer struct {
	k *Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

func (s msgServer) UpdateBlobStatus(ctx context.Context, req *types.MsgUpdateBlobStatusRequest) (*types.MsgUpdateBlobStatusResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// status should be changed to Voting or Ready, depending on the request
	store := sdkCtx.KVStore(s.k.storeKey)
	provenHeight := s.k.GetProvenHeightFromStore(sdkCtx)
	endHeight := s.k.GetEndHeightFromStore(sdkCtx)
	status := GetStatusFromStore(store)

	if req.BlocksRange.From != provenHeight+1 || req.BlocksRange.To != endHeight {
		return nil, fmt.Errorf("invalid blocks range request: expected range [%d -> %d], got [%d -> %d]",
			provenHeight+1, endHeight, req.BlocksRange.From, req.BlocksRange.To)
	}

	if status != PendingState {
		return nil, errors.New("can't update the status if it is not pending")
	}

	newStatus := InVotingState
	if !req.IsSuccess {
		newStatus = FailureState
	} else {
		currentHeight := sdkCtx.BlockHeight()
		UpdateAvailHeight(sdkCtx, store, req.AvailHeight) // updates avail height at which the blocks got submitted to DA
		UpdateVotingEndHeight(sdkCtx, store, uint64(currentHeight)+s.k.VotingInterval)
	}

	UpdateBlobStatus(sdkCtx, store, newStatus)

	return &types.MsgUpdateBlobStatusResponse{}, nil
}
