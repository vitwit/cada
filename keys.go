package availblob

import (
	"encoding/binary"

	"cosmossdk.io/collections"
	"github.com/vitwit/avail-da-module/types"
)

var (

	// ValidatorsKey saves the current validators.
	ValidatorsKey = collections.NewPrefix(0)

	// ClientIDKey saves the current clientID.
	ClientIDKey = collections.NewPrefix(1)

	// ProvenHeightKey saves the current proven height.
	ProvenHeightKey = collections.NewPrefix(2)

	// PendingBlocksToTimeouts maps pending blocks to their timeout
	PendingBlocksToTimeouts = collections.NewPrefix(3)

	// TimeoutsToPendingBlocks maps timeouts to a set of pending blocks
	TimeoutsToPendingBlocks = collections.NewPrefix(4)

	// light client store key
	ClientStoreKey = []byte("client_store/")

	PendingBlobsKey = collections.NewPrefix(5)
)

const (
	// ModuleName is the name of the module
	ModuleName = "availdamodule"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName

	// TransientStoreKey defines the transient store key
	TransientStoreKey = "transient_" + ModuleName
)

func PendingBlobsStoreKey(blocksRange types.Range) []byte {
	fromBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(fromBytes, uint64(blocksRange.From))

	toBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(toBytes, uint64(blocksRange.To))

	key := PendingBlobsKey
	key = append(key, fromBytes...)
	key = append(key, toBytes...)
	return key
}
