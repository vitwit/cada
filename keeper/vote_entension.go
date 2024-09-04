package keeper

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/types"
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

// TODO: change the Vote Extension to be actually usable
type VoteExtension struct {
	Votes map[types.Range]bool
}

func (h *VoteExtHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {

		fmt.Println("coming to extend vote handler.........")
		// TODO: implement proper logic, this is for demo purpose only
		from := h.Keeper.GetStartHeightFromStore(ctx)
		end := h.Keeper.GetEndHeightFromStore(ctx)

		pendingRange := types.Range{
			From: from,
			To:   end,
		}

		var Votes map[types.Range]bool
		Votes[pendingRange] = true
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
