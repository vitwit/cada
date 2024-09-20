package keeper

import (
	"encoding/binary"
	"fmt"

	"cosmossdk.io/collections"
	storetypes2 "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	availblob1 "github.com/vitwit/avail-da-module"
)

const (
	ReadyState    uint32 = 0
	PendingState  uint32 = 1
	InVotingState uint32 = 2
	FailureState  uint32 = 3
)

func ParseStatus(status uint32, start, end uint64) string {
	if start == 0 && end == 0 {
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

func ParseVotingEndHeight(height uint64) string {
	if height == 0 {
		return ""
	}
	return fmt.Sprintln(height)
}

func CanUpdateStatusToPending(store storetypes2.KVStore) bool {
	statusBytes := store.Get(availblob1.BlobStatusKey)
	if len(statusBytes) == 0 {
		return true
	}

	status := binary.BigEndian.Uint32(statusBytes)

	return status == ReadyState || status == FailureState
}

func GetStatusFromStore(store storetypes2.KVStore) uint32 {
	statusBytes := store.Get(availblob1.BlobStatusKey)

	if len(statusBytes) == 0 {
		return ReadyState
	}

	status := binary.BigEndian.Uint32(statusBytes)

	return status
}

func UpdateBlobStatus(_ sdk.Context, store storetypes2.KVStore, status uint32) error {
	statusBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(statusBytes, status)

	store.Set(availblob1.BlobStatusKey, statusBytes)
	return nil
}

func UpdateStartHeight(_ sdk.Context, store storetypes2.KVStore, startHeight uint64) error {
	return updateHeight(store, availblob1.PrevHeightKey, startHeight)
}

func UpdateEndHeight(_ sdk.Context, store storetypes2.KVStore, endHeight uint64) error {
	return updateHeight(store, availblob1.NextHeightKey, endHeight)
}

func UpdateProvenHeight(_ sdk.Context, store storetypes2.KVStore, provenHeight uint64) error {
	return updateHeight(store, availblob1.ProvenHeightKey, provenHeight)
}

func UpdateAvailHeight(_ sdk.Context, store storetypes2.KVStore, availHeight uint64) error {
	return updateHeight(store, availblob1.AvailHeightKey, availHeight)
}

func UpdateVotingEndHeight(_ sdk.Context, store storetypes2.KVStore, votingEndHeight uint64, isLastVoting bool) error {
	key := availblob1.VotingEndHeightKey
	if isLastVoting {
		key = availblob1.LastVotingEndHeightKey
	}
	return updateHeight(store, key, votingEndHeight)
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

func (k *Keeper) GetAvailHeightFromStore(ctx sdk.Context) uint64 {
	return k.getHeight(ctx, availblob1.AvailHeightKey)
}

func (k *Keeper) GetVotingEndHeightFromStore(ctx sdk.Context, isLastVoting bool) uint64 {
	key := availblob1.VotingEndHeightKey
	if isLastVoting {
		key = availblob1.LastVotingEndHeightKey
	}
	return k.getHeight(ctx, key)
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

	if len(heightBytes) == 0 {
		return 0
	}

	height := binary.BigEndian.Uint64(heightBytes)
	return height
}
