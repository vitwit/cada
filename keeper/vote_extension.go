package keeper

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VoteExtHandler struct {
	logger log.Logger

	Keeper *Keeper
}

func NewVoteExtHandler(
	logger log.Logger,
	keeper *Keeper,
) *VoteExtHandler {
	return &VoteExtHandler{
		logger: logger,
		Keeper: keeper,
	}
}

func Key(from, to uint64) string {
	return fmt.Sprintln(from, " ", to)
}

type VoteExtension struct {
	Votes map[string]bool
}

// ExtendVoteHandler handles the extension of votes by providing a vote extension for the given block.
// This function is used to extend the voting information with the necessary vote extensions based on the current blockchain state.
func (h *VoteExtHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, _ *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
		from := h.Keeper.GetStartHeightFromStore(ctx)
		end := h.Keeper.GetEndHeightFromStore(ctx)

		availHeight := h.Keeper.GetAvailHeightFromStore(ctx)

		pendingRangeKey := Key(from, end)

		blobStatus := h.Keeper.GetBlobStatus(ctx)
		currentHeight := ctx.BlockHeight()
		voteEndHeight := h.Keeper.GetVotingEndHeightFromStore(ctx, false)
		Votes := make(map[string]bool, 1)

		abciResponseVoteExt := &abci.ResponseExtendVote{}

		if currentHeight+1 != int64(voteEndHeight) || blobStatus != InVotingState {
			voteExt := VoteExtension{
				Votes: Votes,
			}

			votesBytes, err := json.Marshal(voteExt)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
			}
			abciResponseVoteExt.VoteExtension = votesBytes
			return abciResponseVoteExt, nil
		}

		ok, err := h.Keeper.relayer.IsDataAvailable(ctx, from, end, availHeight)
		if ok {
			h.logger.Info("submitted data to Avail verified successfully at",
				"block_height", availHeight,
			)
		}

		if err != nil {
			fmt.Printf("error while checking for submitted data to avail %v", err)
			return abciResponseVoteExt, err
		}

		Votes[pendingRangeKey] = ok
		voteExt := VoteExtension{
			Votes: Votes,
		}

		votesBytes, err := json.Marshal(voteExt)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
		}

		return &abci.ResponseExtendVote{
			VoteExtension: votesBytes,
		}, nil
	}
}

// VerifyVoteExtensionHandler handles the verification of vote extensions by validating the provided vote extension data.
// This function is used to verify the correctness and validity of the vote extensions submitted during the voting process.
// func (h *VoteExtHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
// 	return func(_ sdk.Context, _ *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
// 		// TODO: add proper validation for the votes if any
// 		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
// 	}
// }

func (h *VoteExtHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
		if req == nil {
			return &abci.ResponseVerifyVoteExtension{
				Status: abci.ResponseVerifyVoteExtension_REJECT,
			}, fmt.Errorf("request is nil")
		}

		// Example: Validate vote height (assuming the vote has a height field)
		if req.Height <= 0 {
			return &abci.ResponseVerifyVoteExtension{
				Status: abci.ResponseVerifyVoteExtension_REJECT,
			}, fmt.Errorf("invalid vote height: %d", req.Height)
		}

		if len(req.VoteExtension) == 0 {
			return &abci.ResponseVerifyVoteExtension{
				Status: abci.ResponseVerifyVoteExtension_REJECT,
			}, fmt.Errorf("vote extension data is empty")
		}

		return &abci.ResponseVerifyVoteExtension{
			Status: abci.ResponseVerifyVoteExtension_ACCEPT,
		}, nil
	}
}
