package keeper

import (
	"context"
	"encoding/binary"
	"fmt"

	storetypes2 "cosmossdk.io/store/types"
	availblob1 "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/types"
)

const (
	READY_STATE     uint32 = 0
	PENDING_STATE   uint32 = 1
	IN_VOTING_STATE uint32 = 2
	FAILURE_STATE   uint32 = 3
)

func IsAlreadyExist(ctx context.Context, store storetypes2.KVStore, blocksRange types.Range) bool {
	pendingBlobStoreKey := availblob1.PendingBlobsStoreKey(blocksRange)
	blobStatus := store.Get(pendingBlobStoreKey)
	fmt.Println("blob status:", blobStatus, blobStatus == nil)
	if blobStatus == nil {
		return false
	}
	return true
}

func IsStateReady(store storetypes2.KVStore) bool {
	fmt.Printf("availblob1.BlobStatusKey: %v\n", availblob1.BlobStatusKey)
	statusBytes := store.Get(availblob1.BlobStatusKey)
	fmt.Printf("statusBytes::::::::::::: %v\n", statusBytes)
	if statusBytes != nil || len(statusBytes) != 0 {
		return true
	}

	fmt.Printf("\"IsStateReadhyyyyyyyyyyyyyy\": %v\n", "IsStateReadhyyyyyyyyyyyyyy")
	fmt.Printf("binary.BigEndian.Uint32(statusBytes): %v\n", binary.BigEndian.Uint32(statusBytes))
	status := binary.BigEndian.Uint32(statusBytes)

	fmt.Printf("status********************: %v\n", status)

	return status == READY_STATE

}

func UpdateBlobStatus(ctx context.Context, store storetypes2.KVStore, status uint32) error {

	statusBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(statusBytes, status)

	store.Set(availblob1.BlobStatusKey, statusBytes)
	return nil
}

func UpdateEndHeight(ctx context.Context, store storetypes2.KVStore, endHeight uint64) error {

	heightBytes := make([]byte, 8)

	binary.BigEndian.PutUint64(heightBytes, endHeight)

	store.Set(availblob1.NextHeightKey, heightBytes)
	return nil
}

func UpdateProvenHeight(ctx context.Context, store storetypes2.KVStore, endHeight uint64) error {

	heightBytes := make([]byte, 8)

	binary.BigEndian.PutUint64(heightBytes, endHeight)

	store.Set(availblob1.ProvenHeightKey, heightBytes)

	fmt.Printf("store.Get(availblob1.ProvenHeightKey): %v\n", store.Get(availblob1.ProvenHeightKey))

	fmt.Printf("\"Update Proven Height\": %v\n", "Update Proven Height")
	return nil
}
