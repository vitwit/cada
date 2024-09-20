package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetBlobStatusPending sets the status of a blob to "pending" and updates its height range.
// This method ensures that the status can be updated to pending before making the changes
// and sets the start and end heights for the blob.
func (k *Keeper) SetBlobStatusPending(ctx sdk.Context, startHeight, endHeight uint64) bool {
	store := ctx.KVStore(k.storeKey)

	if !CanUpdateStatusToPending(store) { // TOodo: we should check for expiration too
		return false
	}

	UpdateBlobStatus(ctx, store, PendingState)
	UpdateStartHeight(ctx, store, startHeight)
	UpdateEndHeight(ctx, store, endHeight)
	return true
}

// SetBlobStatus updates the status of the blob
func (k *Keeper) SetBlobStatus(ctx sdk.Context, state uint32) error {
	store := ctx.KVStore(k.storeKey)

	if state == ReadyState {
		endHeight := k.GetEndHeightFromStore(ctx)
		err := UpdateProvenHeight(ctx, store, endHeight)
		if err != nil {
			return err
		}
	}

	return UpdateBlobStatus(ctx, store, state)
}
