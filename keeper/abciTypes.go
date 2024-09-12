package keeper

import abci "github.com/cometbft/cometbft/abci/types"

type StakeWeightedVotes struct {
	Votes              map[string]int64
	ExtendedCommitInfo abci.ExtendedCommitInfo
}
