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

// UpdateBlobStatus updates the status of the blob.
// This method verifies the blocks range and updates the status to either Voting or Failure,
// depending on the request's success flag.
func (s msgServer) UpdateBlobStatus(ctx context.Context, req *types.MsgUpdateBlobStatusRequest) (*types.MsgUpdateBlobStatusResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// status should be changed to Voting or Ready, depending on the request
	store := sdkCtx.KVStore(s.k.storeKey)
	provenHeight := s.k.GetProvenHeightFromStore(sdkCtx)
	endHeight := s.k.GetEndHeightFromStore(sdkCtx)
	status := GetStatusFromStore(store)

	// Validate the block range
	if req.BlocksRange.From != provenHeight+1 || req.BlocksRange.To != endHeight {
		return nil, fmt.Errorf("invalid blocks range request: expected range [%d -> %d], got [%d -> %d]",
			provenHeight+1, endHeight, req.BlocksRange.From, req.BlocksRange.To)
	}

	// Ensure that the blob is in the PendingState before updating
	if status != PendingState {
		return nil, errors.New("can't update the status if it is not pending")
	}

	newStatus := InVotingState
	if !req.IsSuccess {
		// Mark as failure if the request indicates failure
		newStatus = FailureState
	} else {
		// If success, update the voting-related heights
		currentHeight := sdkCtx.BlockHeight()

		// Update the avail height at which blocks were submitted to DA
		UpdateAvailHeight(sdkCtx, store, req.AvailHeight)

		// Retrieve the last voting end height
		lastVotingEndHeight := s.k.GetVotingEndHeightFromStore(sdkCtx, false)

		// Set the new voting end height, with the configured interval
		newVotingEndHeight := uint64(currentHeight) + s.k.relayer.AvailConfig.VoteInterval
		UpdateVotingEndHeight(sdkCtx, store, newVotingEndHeight, false)

		// Update the previous voting end height to mark the end of the last round
		UpdateVotingEndHeight(sdkCtx, store, lastVotingEndHeight, true)
	}

	// Finally, update the blob status in the store
	UpdateBlobStatus(sdkCtx, store, newStatus)

	return &types.MsgUpdateBlobStatusResponse{}, nil
}
