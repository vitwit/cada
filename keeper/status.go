package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

func (k *Keeper) setBlobStatusSuccess(ctx sdk.Context) error {
	store := ctx.KVStore(k.storeKey)
	endHeight := k.GetEndHeightFromStore(ctx)

	err := UpdateProvenHeight(ctx, store, endHeight)
	if err != nil {
		return err
	}
	return UpdateBlobStatus(ctx, store, READY_STATE)
}

func (k *Keeper) SetBlobStatusFailure(ctx sdk.Context) error {

	store := ctx.KVStore(k.storeKey)

	return UpdateBlobStatus(ctx, store, FAILURE_STATE)
}
