package keeper

import (
	"bytes"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ProofOfBlobProposalHandler struct {
	keeper *Keeper

	prepareProposalHandler sdk.PrepareProposalHandler
	processProposalHandler sdk.ProcessProposalHandler
}

func NewProofOfBlobProposalHandler(
	k *Keeper,
	prepareProposalHandler sdk.PrepareProposalHandler,
	processProposalHandler sdk.ProcessProposalHandler,
) *ProofOfBlobProposalHandler {
	return &ProofOfBlobProposalHandler{
		keeper:                 k,
		prepareProposalHandler: prepareProposalHandler,
		processProposalHandler: processProposalHandler,
	}
}

func (h *ProofOfBlobProposalHandler) PrepareProposal(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
	h.keeper.proposerAddress = req.ProposerAddress

	resp, err := h.prepareProposalHandler(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *ProofOfBlobProposalHandler) ProcessProposal(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
	return h.processProposalHandler(ctx, req)
}

func (k *Keeper) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) error {

	currentBlockHeight := ctx.BlockHeight()
	if k.IsValidBlockToPostTODA(uint64(currentBlockHeight)) {
		return nil
	}

	fromHeight := k.GetProvenHeightFromStore(ctx) //Todo: change this get from ProvenHeight from store
	endHeight := min(fromHeight+uint64(k.MaxBlocksForBlob), uint64(ctx.BlockHeight()))

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.SetBlobStatusPending(sdkCtx, fromHeight, endHeight)

	// only proposar should should run the this
	if bytes.Equal(req.ProposerAddress, k.proposerAddress) {

		// Todo: run the relayer routine
		// relayer doesn't have to make submitBlob Transaction, it should just start DA submissio

	}

	return nil
}

func (k *Keeper) IsValidBlockToPostTODA(height uint64) bool {
	if uint64(height) <= uint64(1) {
		return false
	}

	if (height-1)%k.PublishToAvailBlockInterval != 0 {
		return false
	}

	return true
}
