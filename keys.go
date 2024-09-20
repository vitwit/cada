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
	// PendingBlocksToTimeouts = collections.NewPrefix(3)

	// TimeoutsToPendingBlocks maps timeouts to a set of pending blocks
	// TimeoutsToPendingBlocks = collections.NewPrefix(4)

	// light client store key
	// ClientStoreKey = []byte("client_store/")

	PendingBlobsKey = collections.NewPrefix(3)

	BlobStatusKey = collections.NewPrefix(4)

	PrevHeightKey = collections.NewPrefix(5)

	NextHeightKey = collections.NewPrefix(6)

	VotingEndHeightKey = collections.NewPrefix(7)

	LastVotingEndHeightKey = collections.NewPrefix(8)

	AvailHeightKey = collections.NewPrefix(9)
)

const (
	// ModuleName is the name of the module
	ModuleName = "cada"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName
)

// PendingBlobsStoreKey generates a store key for pending blobs based on the given block range.
// The key is constructed by appending the byte-encoded 'From' and 'To' values from the `blocksRange`
// to a base key (`PendingBlobsKey`). This unique key is used to store and retrieve pending blob data
// in a key-value store.
func PendingBlobsStoreKey(blocksRange types.Range) []byte {
	fromBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(fromBytes, blocksRange.From)

	toBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(toBytes, blocksRange.To)

	key := PendingBlobsKey
	key = append(key, fromBytes...)
	key = append(key, toBytes...)
	return key
}
