package types

import (
	"cosmossdk.io/collections"
)

var (

	// ValidatorsKey saves the current validators.
	ValidatorsKey = collections.NewPrefix(0)

	// ClientIDKey saves the current clientID.
	ClientIDKey = collections.NewPrefix(1)

	// ProvenHeightKey saves the current proven height.
	ProvenHeightKey = collections.NewPrefix(2)

	PendingBlobsKey = collections.NewPrefix(3)

	BlobStatusKey = collections.NewPrefix(4)

	PrevHeightKey = collections.NewPrefix(5)

	NextHeightKey = collections.NewPrefix(6)

	VotingEndHeightKey = collections.NewPrefix(7)

	LastVotingEndHeightKey = collections.NewPrefix(8)

	AvailHeightKey = collections.NewPrefix(9)

	// PendingHeightKey saves the height at which the state is changed to pending
	PendingHeightKey = collections.NewPrefix(10)
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
