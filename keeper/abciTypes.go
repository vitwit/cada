package keeper

import abci "github.com/cometbft/cometbft/abci/types"

// StakeWeightedVotes stores the stake-weighted votes from validators and the commit info
// for a consensus round.
type StakeWeightedVotes struct {
	// A map where the key is the validator's address , and the value is the
	// validator's voting power.
	Votes map[string]int64

	// ExtendedCommitInfo Contains additional information about the commit phase, including
	//  vote extensions and details about the current consensus round.
	ExtendedCommitInfo abci.ExtendedCommitInfo
}
