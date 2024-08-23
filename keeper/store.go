package keeper

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	storetypes2 "cosmossdk.io/store/types"
	availblob1 "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/types"
)

const (
	PENDING uint32 = 0
	SUCCESS uint32 = 1
	FAILURE uint32 = 2
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

func updateBlobStatus(ctx context.Context, store storetypes2.KVStore, blocksRange types.Range, status uint32) error {
	if status != PENDING && status != SUCCESS && status != FAILURE {
		return errors.New("unknown status")
	}
	pendingBlobStoreKey := availblob1.PendingBlobsStoreKey(blocksRange)

	statusBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(statusBytes, status)

	store.Set(pendingBlobStoreKey, statusBytes)
	return nil
}
