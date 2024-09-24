package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StakeWeightedVotes represents the aggregated stake-weighted votes from validators,
// along with the associated commit information for a specific consensus round.
type StakeWeightedVotes struct {
	// A map where the key is the range of pending blocks(e.g. "1 10"), and the value is the
	// validator's voting power.
	Votes map[string]int64

	// ExtendedCommitInfo Contains additional information about the commit phase, including
	//  vote extensions and details about the current consensus round.
	ExtendedCommitInfo abci.ExtendedCommitInfo

	// TotalVotingPower is the sum of all validators' voting power
	TotalVotingPower int64
}

// ProofOfBlobProposalHandler manages the proposal and vote extension logic related to
// blob transactions in the consensus process.
type ProofOfBlobProposalHandler struct {
	keeper *Keeper

	// prepareProposalHandler is responsible for preparing proposals. This is called
	// during the proposal preparation phase to adjust or validate the proposal before
	// it's broadcast to other validators.
	prepareProposalHandler sdk.PrepareProposalHandler

	// processProposalHandler processes proposals during the consensus. It verifies
	// the validity of the received proposal and determines whether it should be accepted or rejected.
	processProposalHandler sdk.ProcessProposalHandler

	// VoteExtHandler handles vote extensions, allowing validators to submit additional information
	// like avail height along with their votes during the consensus process.
	VoteExtHandler VoteExtHandler
}

// NewProofOfBlobProposalHandler creates a new instance of the ProofOfBlobProposalHandler.
func NewProofOfBlobProposalHandler(
	k *Keeper,
	prepareProposalHandler sdk.PrepareProposalHandler,
	processProposalHandler sdk.ProcessProposalHandler,
	voteExtHandler VoteExtHandler,
) *ProofOfBlobProposalHandler {
	return &ProofOfBlobProposalHandler{
		keeper:                 k,
		prepareProposalHandler: prepareProposalHandler,
		processProposalHandler: processProposalHandler,
		VoteExtHandler:         voteExtHandler,
	}
}

// PrepareProposal is responsible for preparing a proposal by aggregating vote extensions
// and injecting them into the list of transactions for the proposal.
func (h *ProofOfBlobProposalHandler) PrepareProposal(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
	h.keeper.ProposerAddress = req.ProposerAddress
	proposalTxs := req.Txs

	votes, totalVotingPower, err := h.aggregateVotes(ctx, req.LocalLastCommit)
	if err != nil {
		fmt.Println("error while aggregating votes", err)
		return nil, err
	}

	injectedVoteExtTx := StakeWeightedVotes{
		Votes:              votes,
		ExtendedCommitInfo: req.LocalLastCommit,
		TotalVotingPower:   totalVotingPower,
	}

	// if there is any another tx, it might give any marshelling error, so ignoring this err
	bz, _ := json.Marshal(injectedVoteExtTx)

	proposalTxs = append([][]byte{bz}, proposalTxs...)
	return &abci.ResponsePrepareProposal{
		Txs: proposalTxs,
	}, nil
}

// ProcessProposal handles the validation and processing of a proposed block during the consensus process.
// It checks if the proposal contains any transactions, attempts to decode the injected vote extension
// transaction, and performs necessary validations before deciding to accept or reject the proposal.
func (h *ProofOfBlobProposalHandler) ProcessProposal(_ sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
	if len(req.Txs) == 0 {
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}

	var injectedVoteExtTx StakeWeightedVotes
	if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err != nil {
		// if there is any another tx, it might give any unmarshelling error, so ignoring this err
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}

	// TODO: write some validations
	// if injectedVoteExtTx.ExtendedCommitInfo != nil {// }

	return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
}

// 66% of voting power is needed. Change the percentage if required
func isEnoughVoting(voting, totalVoting int64) bool {
	// division can be inaccurate due to decimal roundups
	// voting / totalVoting * 100 > 66
	// voting * 100 > 66 * totalVoting
	return voting*100 > 66*totalVoting
}

// PreBlocker runs before finalizing each block, responsible for handling vote extensions
// and managing the posting of blocks to the Avail light client.
func (k *Keeper) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) error {
	votingEndHeight := k.GetVotingEndHeightFromStore(ctx, false)
	blobStatus := k.GetBlobStatus(ctx)
	currentHeight := ctx.BlockHeight()

	if len(req.Txs) > 0 && currentHeight == int64(votingEndHeight) && blobStatus == InVotingState {
		var injectedVoteExtTx StakeWeightedVotes
		if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err == nil {
			from := k.GetStartHeightFromStore(ctx)
			to := k.GetEndHeightFromStore(ctx)

			pendingRangeKey := Key(from, to)
			votingPower := injectedVoteExtTx.Votes[pendingRangeKey]

			state := FailureState
			totalVotingPower := injectedVoteExtTx.TotalVotingPower

			if isEnoughVoting(votingPower, totalVotingPower) {
				state = ReadyState
			}

			store := ctx.KVStore(k.storeKey)
			UpdateVotingEndHeight(ctx, store, 0, false)
			k.SetBlobStatus(ctx, state)
		}
	}

	currentBlockHeight := ctx.BlockHeight()
	if !k.IsValidBlockToPostToDA(uint64(currentBlockHeight)) {
		return nil
	}

	// Calculate pending range of blocks to post data
	provenHeight := k.GetProvenHeightFromStore(ctx)
	fromHeight := provenHeight + 1
	endHeight := min(fromHeight+k.relayer.AvailConfig.MaxBlobBlocks, uint64(ctx.BlockHeight())) // exclusive i.e [fromHeight, endHeight)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	ok := k.SetBlobStatusPending(sdkCtx, fromHeight, endHeight-1)
	if !ok {
		return nil
	}

	var blocksToSumit []int64
	for i := fromHeight; i < endHeight; i++ {
		blocksToSumit = append(blocksToSumit, int64(i))
	}

	// only proposar should should run the this
	if bytes.Equal(req.ProposerAddress, k.ProposerAddress) {
		k.relayer.PostBlocks(ctx, blocksToSumit, k.cdc, req.ProposerAddress)
	}

	return nil
}

// IsValidBlockToPostToDA checks if the given block height is valid for posting data.
// The block is considered valid if it meets the defined interval for posting.
func (k *Keeper) IsValidBlockToPostToDA(height uint64) bool {
	if height <= uint64(1) {
		return false
	}

	if (height-1)%k.relayer.AvailConfig.PublishBlobInterval != 0 {
		return false
	}

	return true
}

// For each vote, it decodes the `VoteExtension`, which contains the vote information for
// specific block ranges.
// If the vote extension contains a vote for the pending range, it sums the voting power
// of validators.
func (h *ProofOfBlobProposalHandler) aggregateVotes(ctx sdk.Context, ci abci.ExtendedCommitInfo) (map[string]int64, int64, error) {
	from := h.keeper.GetStartHeightFromStore(ctx)
	to := h.keeper.GetEndHeightFromStore(ctx)

	pendingRangeKey := Key(from, to)
	votes := make(map[string]int64, 1)
	totalVoting := 0

	for _, v := range ci.Votes {
		// Process only votes with BlockIDFlagCommit, indicating the validator committed to the block.
		// Skip votes with other flags (e.g., BlockIDFlagUnknown, BlockIDFlagNil).
		if v.BlockIdFlag != cmtproto.BlockIDFlagCommit {
			continue
		}

		var voteExt VoteExtension
		if err := json.Unmarshal(v.VoteExtension, &voteExt); err != nil {
			h.VoteExtHandler.logger.Error("failed to decode vote extension", "err", err, "validator", fmt.Sprintf("%x", v.Validator.Address))
			continue
		}

		if voteExt.Votes == nil {
			continue
		}

		currentTotalVoting := 0
		for voteRange, isVoted := range voteExt.Votes {

			currentTotalVoting += int(v.Validator.Power)
			if voteRange != pendingRangeKey || !isVoted {
				continue
			}

			votes[voteRange] += v.Validator.Power
		}

		totalVoting = max(totalVoting, currentTotalVoting)

	}
	return votes, int64(totalVoting), nil
}
