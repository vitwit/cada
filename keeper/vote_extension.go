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

// TODO: add required parameters like avail light client url, etc..
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

func (h *VoteExtHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, _ *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
		from := h.Keeper.GetStartHeightFromStore(ctx)
		end := h.Keeper.GetEndHeightFromStore(ctx)

		availHeight := h.Keeper.GetAvailHeightFromStore(ctx)

		pendingRangeKey := Key(from, end)

		blobStatus := h.Keeper.GetBlobStatus(ctx)
		currentHeight := ctx.BlockHeight()
		voteEndHeight := h.Keeper.GetVotingEndHeightFromStore(ctx)
		Votes := make(map[string]bool, 1)

		abciResponseVoteExt := &abci.ResponseExtendVote{}

		if currentHeight+1 != int64(voteEndHeight) || blobStatus != InVotingState {
			voteExt := VoteExtension{
				Votes: Votes,
			}

			// TODO: use better marshaling instead of json (eg: proto marshaling)
			votesBytes, err := json.Marshal(voteExt)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
			}
			abciResponseVoteExt.VoteExtension = votesBytes
			return abciResponseVoteExt, nil
		}

		ok, err := h.Keeper.relayer.IsDataAvailable(ctx, from, end, availHeight, "http://localhost:8000") // TODO: read light client url from config
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

		// TODO: use proto marshaling instead
		votesBytes, err := json.Marshal(voteExt)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
		}

		return &abci.ResponseExtendVote{
			VoteExtension: votesBytes,
		}, nil
	}
}

func (h *VoteExtHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(_ sdk.Context, _ *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
		// TODO: write proper validation for the votes
		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}
