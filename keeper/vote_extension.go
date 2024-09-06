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

// TODO: change the Vote Extension to be actually usable
type VoteExtension struct {
	Votes map[string]bool
}

// ExtendVoteHandler handles the extension of votes by providing a vote extension for the given block.
// This function is used to extend the voting information with the necessary vote extensions based on the current blockchain state.
func (h *VoteExtHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {

	return func(ctx sdk.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {

		// TODO: implement proper logic, this is for demo purpose only
		from := h.Keeper.GetStartHeightFromStore(ctx)
		end := h.Keeper.GetEndHeightFromStore(ctx)

		availHeight := h.Keeper.GetAvailHeightFromStore(ctx)

		pendingRangeKey := Key(from, end)

		blobStatus := h.Keeper.GetBlobStatus(ctx)
		currentHeight := ctx.BlockHeight()
		voteEndHeight := h.Keeper.GetVotingEndHeightFromStore(ctx)
		Votes := make(map[string]bool, 1)

		abciResponseVoteExt := &abci.ResponseExtendVote{}

		if currentHeight+1 != int64(voteEndHeight) || blobStatus != IN_VOTING_STATE {
			voteExt := VoteExtension{
				Votes: Votes,
			}

			//TODO: use better marshalling instead of json (eg: proto marshalling)
			votesBytes, err := json.Marshal(voteExt)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
			}
			abciResponseVoteExt.VoteExtension = votesBytes
			return abciResponseVoteExt, nil
		}

		ok, err := h.Keeper.relayer.IsDataAvailable(ctx, from, end, availHeight, "http://localhost:8000")
		fmt.Println("checking light client...", ok, err)

		// ok, checkLightClient()
		Votes[pendingRangeKey] = true
		voteExt := VoteExtension{
			Votes: Votes,
		}

		//TODO: use proto marshalling instead
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
func (h *VoteExtHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
		// TODO: write proper validation for the votes
		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}
