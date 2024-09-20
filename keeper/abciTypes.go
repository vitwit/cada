package keeper

import abci "github.com/cometbft/cometbft/abci/types"

// StakeWeightedVotes represents the aggregated stake-weighted votes from validators,
// along with the associated commit information for a specific consensus round.
type StakeWeightedVotes struct {
	// A map where the key is the range of pending blocks(e.g. "1 10"), and the value is the
	// validator's voting power.
	Votes map[string]int64

	// ExtendedCommitInfo Contains additional information about the commit phase, including
	//  vote extensions and details about the current consensus round.
	ExtendedCommitInfo abci.ExtendedCommitInfo
}
