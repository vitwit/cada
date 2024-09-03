package keeper

import (
	"encoding/binary"

	"cosmossdk.io/collections"
	storetypes2 "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	availblob1 "github.com/vitwit/avail-da-module"
)

const (
	READY_STATE     uint32 = 0
	PENDING_STATE   uint32 = 1
	IN_VOTING_STATE uint32 = 2
	FAILURE_STATE   uint32 = 3
)

func ParseStatus(status uint32) string {
	switch status {
	case READY_STATE:
		return "SUCCESS"
	case PENDING_STATE:
		return "PENDING"
	case IN_VOTING_STATE:
		return "IN_VOTING"
	case FAILURE_STATE:
		return "FAILUTE"
	default:
		return "UNKNOWN"
	}
}

func CanUpdateStatusToPending(store storetypes2.KVStore) bool {
	statusBytes := store.Get(availblob1.BlobStatusKey)
	if statusBytes == nil || len(statusBytes) == 0 {
		return true
	}

	status := binary.BigEndian.Uint32(statusBytes)

	return status == READY_STATE || status == FAILURE_STATE
}

func GetStatusFromStore(store storetypes2.KVStore) uint32 {
	statusBytes := store.Get(availblob1.BlobStatusKey)

	if statusBytes == nil || len(statusBytes) == 0 {
		return READY_STATE
	}

	status := binary.BigEndian.Uint32(statusBytes)

	return status
}

func UpdateBlobStatus(ctx sdk.Context, store storetypes2.KVStore, status uint32) error {

	statusBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(statusBytes, status)

	store.Set(availblob1.BlobStatusKey, statusBytes)
	return nil
}

func UpdateStartHeight(ctx sdk.Context, store storetypes2.KVStore, startHeight uint64) error {
	return updateHeight(store, availblob1.PrevHeightKey, startHeight)
}

func UpdateEndHeight(ctx sdk.Context, store storetypes2.KVStore, endHeight uint64) error {
	return updateHeight(store, availblob1.NextHeightKey, endHeight)
}

func UpdateProvenHeight(ctx sdk.Context, store storetypes2.KVStore, provenHeight uint64) error {
	return updateHeight(store, availblob1.ProvenHeightKey, provenHeight)
}

func updateHeight(store storetypes2.KVStore, key collections.Prefix, height uint64) error {
	heightBytes := make([]byte, 8)

	binary.BigEndian.PutUint64(heightBytes, height)

	store.Set(key, heightBytes)
	return nil
}

func (k *Keeper) GetProvenHeightFromStore(ctx sdk.Context) uint64 {
	return k.getHeight(ctx, availblob1.ProvenHeightKey)
}

func (k *Keeper) GetStartHeightFromStore(ctx sdk.Context) uint64 {
	return k.getHeight(ctx, availblob1.PrevHeightKey)
}

func (k *Keeper) GetEndHeightFromStore(ctx sdk.Context) uint64 {

	return k.getHeight(ctx, availblob1.NextHeightKey)
}

func (k *Keeper) getHeight(ctx sdk.Context, key collections.Prefix) uint64 {
	store := ctx.KVStore(k.storeKey)
	heightBytes := store.Get(key)

	if heightBytes == nil || len(heightBytes) == 0 {
		return 0
	}

	height := binary.BigEndian.Uint64(heightBytes)
	return height
}
