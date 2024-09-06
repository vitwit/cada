package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetBlobStatusPending sets the status of a blob to "pending" and updates its height range.
// This method ensures that the status can be updated to pending before making the changes
// and sets the start and end heights for the blob.
func (k *Keeper) SetBlobStatusPending(ctx sdk.Context, startHeight, endHeight uint64) bool {

	store := ctx.KVStore(k.storeKey)

	if !CanUpdateStatusToPending(store) { //TOodo: we should check for expiration too
		return false
	}

	UpdateBlobStatus(ctx, store, PENDING_STATE)
	UpdateStartHeight(ctx, store, startHeight)
	UpdateEndHeight(ctx, store, endHeight)
	return true
}

// setBlobStatusSuccess updates the status of a blob to "ready" and sets the proven height.
// This method marks the blob as successfully processed and updates the proven height to the blob's end height.
func (k *Keeper) setBlobStatusSuccess(ctx sdk.Context) error {
	store := ctx.KVStore(k.storeKey)
	endHeight := k.GetEndHeightFromStore(ctx)

	err := UpdateProvenHeight(ctx, store, endHeight)
	if err != nil {
		return err
	}
	return UpdateBlobStatus(ctx, store, READY_STATE)
}

// SetBlobStatusFailure sets the status of a blob to "failure".
// This method updates the blob status to indicate that the processing has failed.
func (k *Keeper) SetBlobStatusFailure(ctx sdk.Context) error {

	store := ctx.KVStore(k.storeKey)

	return UpdateBlobStatus(ctx, store, FAILURE_STATE)
}
