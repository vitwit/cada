package keeper

import (
	"encoding/binary"
	"strconv"

	"cosmossdk.io/collections"
	storetypes2 "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/vitwit/avail-da-module/types"
)

const (
	ReadyState    uint32 = 0
	PendingState  uint32 = 1
	InVotingState uint32 = 2
	FailureState  uint32 = 3
)

func ParseStatus(status uint32, startHeight, endHeight uint64) string {
	if startHeight == 0 && endHeight == 0 {
		return ""
	}

	switch status {
	case ReadyState:
		return "SUCCESS"
	case PendingState:
		return "PENDING"
	case InVotingState:
		return "IN_VOTING"
	case FailureState:
		return "FAILURE"
	default:
		return "UNKNOWN"
	}
}

// ParseVotingEndHeight converts a given voting end height from uint64 to a string format.
func ParseVotingEndHeight(height uint64) string {
	if height == 0 {
		return ""
	}
	return strconv.FormatUint(height, 10)
}

// CanUpdateStatusToPending checks if the blob status can be updated to "pending".
// This function verifies whether the current status allows transitioning to the "pending" state.
func CanUpdateStatusToPending(store storetypes2.KVStore) bool {
	statusBytes := store.Get(types.BlobStatusKey)
	if len(statusBytes) == 0 {
		return true
	}

	status := binary.BigEndian.Uint32(statusBytes)

	return status == ReadyState || status == FailureState
}

// GetStatusFromStore retrieves the current status of the blob from the store.
func GetStatusFromStore(store storetypes2.KVStore) uint32 {
	statusBytes := store.Get(types.BlobStatusKey)

	if len(statusBytes) == 0 {
		return ReadyState
	}

	status := binary.BigEndian.Uint32(statusBytes)

	return status
}

// UpdateBlobStatus updates the blob status in the KV store.
func UpdateBlobStatus(_ sdk.Context, store storetypes2.KVStore, status uint32) error {
	statusBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(statusBytes, status)

	store.Set(types.BlobStatusKey, statusBytes)
	return nil
}

// UpdateStartHeight updates the start height in the KV store.
func UpdateStartHeight(_ sdk.Context, store storetypes2.KVStore, startHeight uint64) error {
	return updateHeight(store, types.PrevHeightKey, startHeight)
}

// UpdateEndHeight updates the end height in the KV store.
func UpdateEndHeight(_ sdk.Context, store storetypes2.KVStore, endHeight uint64) error {
	return updateHeight(store, types.NextHeightKey, endHeight)
}

// UpdateProvenHeight updates the proven height in the KV store.
func UpdateProvenHeight(_ sdk.Context, store storetypes2.KVStore, provenHeight uint64) error {
	return updateHeight(store, types.ProvenHeightKey, provenHeight)
}

// UpdateAvailHeight updates the avail height in the store
func UpdateAvailHeight(_ sdk.Context, store storetypes2.KVStore, availHeight uint64) error {
	return updateHeight(store, types.AvailHeightKey, availHeight)
}

// UpdateVotingEndHeight updates the voting end height in the KV store.
func UpdateVotingEndHeight(_ sdk.Context, store storetypes2.KVStore, votingEndHeight uint64, isLastVoting bool) error {
	key := types.VotingEndHeightKey
	if isLastVoting {
		key = types.LastVotingEndHeightKey
	}
	return updateHeight(store, key, votingEndHeight)
}

// updateHeight encodes and stores a height value in the KV store.
func updateHeight(store storetypes2.KVStore, key collections.Prefix, height uint64) error {
	heightBytes := make([]byte, 8)

	binary.BigEndian.PutUint64(heightBytes, height)

	store.Set(key, heightBytes)
	return nil
}

// GetProvenHeightFromStore retrieves the proven height from the KV store.
func (k *Keeper) GetProvenHeightFromStore(ctx sdk.Context) uint64 {
	return k.getHeight(ctx, types.ProvenHeightKey)
}

// GetAvailHeightFromStore retrieves the avail height from the KV store.
func (k *Keeper) GetAvailHeightFromStore(ctx sdk.Context) uint64 {
	return k.getHeight(ctx, types.AvailHeightKey)
}

// GetVotingEndHeightFromStore retrieves the ending vote height from store
func (k *Keeper) GetVotingEndHeightFromStore(ctx sdk.Context, isLastVoting bool) uint64 {
	key := types.VotingEndHeightKey
	if isLastVoting {
		key = types.LastVotingEndHeightKey
	}
	return k.getHeight(ctx, key)
}

// GetStartHeightFromStore retrieves the start height from store
func (k *Keeper) GetStartHeightFromStore(ctx sdk.Context) uint64 {
	return k.getHeight(ctx, types.PrevHeightKey)
}

// GetEndHeightFromStore retrieves the end height from store
func (k *Keeper) GetEndHeightFromStore(ctx sdk.Context) uint64 {
	return k.getHeight(ctx, types.NextHeightKey)
}

// getHeight retrieves and decodes a height value from the KV store.
func (k *Keeper) getHeight(ctx sdk.Context, key collections.Prefix) uint64 {
	store := ctx.KVStore(k.storeKey)
	heightBytes := store.Get(key)

	if len(heightBytes) == 0 {
		return 0
	}

	height := binary.BigEndian.Uint64(heightBytes)
	return height
}
