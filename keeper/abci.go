package keeper

import (
	"bytes"
	"fmt"

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
	fmt.Println("coming hereee.........", req)

	currentBlockHeight := ctx.BlockHeight()
	if k.IsValidBlockToPostTODA(uint64(currentBlockHeight)) {
		return nil
	}

	fromHeight := k.GetProvenHeightFromStore(ctx) //Todo: change this get from ProvenHeight from store
	fmt.Println("from height..", fromHeight)
	endHeight := min(fromHeight+uint64(k.MaxBlocksForBlob), uint64(ctx.BlockHeight()))
	fmt.Println("end height..", endHeight)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := k.SetBlobStatusPending(sdkCtx, fromHeight, endHeight)
	if err != nil {
		fmt.Println("error while setting blob status...", err)
	}

	// only proposar should should run the this
	if bytes.Equal(req.ProposerAddress, k.proposerAddress) {
		// update blob status to success

		err = k.SetBlobStatusSuccess(sdkCtx, fromHeight, endHeight)
		if err != nil {
			return nil
		}

		// Todo: run the relayer routine
		// relayer doesn't have to make submitBlob Transaction, it should just start DA submission
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
