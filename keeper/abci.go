package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ProofOfBlobProposalHandler struct {
	keeper *Keeper

	prepareProposalHandler sdk.PrepareProposalHandler
	processProposalHandler sdk.ProcessProposalHandler
	VoteExtHandler         VoteExtHandler
}

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

func (h *ProofOfBlobProposalHandler) PrepareProposal(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
	h.keeper.ProposerAddress = req.ProposerAddress
	proposalTxs := req.Txs

	votes, err := h.aggregateVotes(ctx, req.LocalLastCommit)
	if err != nil {
		fmt.Println("error while aggregating votes", err)
		return nil, err
	}

	injectedVoteExtTx := StakeWeightedVotes{
		Votes:              votes,
		ExtendedCommitInfo: req.LocalLastCommit,
	}
	bz, err := json.Marshal(injectedVoteExtTx)
	if err != nil {
		fmt.Println("failed to encode injected vote extension tx", "err", err)
	}

	proposalTxs = append(proposalTxs, bz)
	return &abci.ResponsePrepareProposal{
		Txs: proposalTxs,
	}, nil
}

func (h *ProofOfBlobProposalHandler) ProcessProposal(_ sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
	if len(req.Txs) == 0 {
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}

	var injectedVoteExtTx StakeWeightedVotes
	if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err != nil {
		fmt.Println("failed to decode injected vote extension tx", "err", err)
		// return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
	}

	// TODO: write some validations

	return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
}

func (k *Keeper) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) error {
	votingEndHeight := k.GetVotingEndHeightFromStore(ctx, false)
	blobStatus := k.GetBlobStatus(ctx)
	currentHeight := ctx.BlockHeight()

	if len(req.Txs) > 0 && currentHeight == int64(votingEndHeight) && blobStatus == InVotingState {
		var injectedVoteExtTx StakeWeightedVotes
		if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err != nil {
			fmt.Println("preblocker failed to decode injected vote extension tx", "err", err)
		} else {
			from := k.GetStartHeightFromStore(ctx)
			to := k.GetEndHeightFromStore(ctx)

			pendingRangeKey := Key(from, to)
			votingPower := injectedVoteExtTx.Votes[pendingRangeKey]

			state := FailureState

			if votingPower > 0 { // TODO: calculate voting power properly
				state = ReadyState
			}

			store := ctx.KVStore(k.storeKey)
			UpdateVotingEndHeight(ctx, store, 0, false)
			k.SetBlobStatus(ctx, state)
		}
	}

	currentBlockHeight := ctx.BlockHeight()
	if !k.IsValidBlockToPostTODA(uint64(currentBlockHeight)) {
		return nil
	}

	provenHeight := k.GetProvenHeightFromStore(ctx)
	fromHeight := provenHeight + 1
	endHeight := min(fromHeight+uint64(k.MaxBlocksForBlob), uint64(ctx.BlockHeight())) // exclusive i.e [fromHeight, endHeight)
	// Calculate pending range of blocks to post data

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

func (k *Keeper) IsValidBlockToPostTODA(height uint64) bool {
	if height <= uint64(1) {
		return false
	}

	if (height-1)%k.PublishToAvailBlockInterval != 0 {
		return false
	}

	return true
}

func (h *ProofOfBlobProposalHandler) aggregateVotes(ctx sdk.Context, ci abci.ExtendedCommitInfo) (map[string]int64, error) {
	from := h.keeper.GetStartHeightFromStore(ctx)
	to := h.keeper.GetEndHeightFromStore(ctx)

	pendingRangeKey := Key(from, to)
	votes := make(map[string]int64, 1)

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

		for voteRange, isVoted := range voteExt.Votes {
			if voteRange != pendingRangeKey || !isVoted {
				continue
			}

			votes[voteRange] += v.Validator.Power
		}

	}
	return votes, nil
}
