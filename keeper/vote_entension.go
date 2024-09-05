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

func (h *VoteExtHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {

	return func(ctx sdk.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {

		fmt.Println("coming to extend vote handler.........")
		// TODO: implement proper logic, this is for demo purpose only
		from := h.Keeper.GetStartHeightFromStore(ctx)
		end := h.Keeper.GetEndHeightFromStore(ctx)

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

func (h *VoteExtHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
		// TODO: write proper validation for the votes
		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}
