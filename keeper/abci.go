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
	store := ctx.KVStore(k.storeKey)

	currentBlockHeight := ctx.BlockHeight()
	fmt.Printf("currentBlockHeight: %v\n", currentBlockHeight)
	if k.IsValidBlockToPostTODA(uint64(currentBlockHeight)) {
		return nil
	}

	fromHeight := k.GetProvenHeightFromStore(ctx)
	endHeight := min(fromHeight+uint64(k.MaxBlocksForBlob), uint64(ctx.BlockHeight()))

	fmt.Printf("fromHeight: %v\n", fromHeight)
	fmt.Printf("endHeight: %v\n", endHeight)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := k.SetBlobStatusPending(sdkCtx, fromHeight, endHeight)
	if err != nil {
		return err
	}

	fmt.Printf("req.ProposerAddress: %v\n", req.ProposerAddress)
	fmt.Printf("k.proposerAddress: %v\n", k.proposerAddress)
	// only proposar should should run the this
	if bytes.Equal(req.ProposerAddress, k.proposerAddress) {

		k.relayer.PostBlocks(sdkCtx, fromHeight, endHeight)

		// Todo: run the relayer routine
		// relayer doesn't have to make submitBlob Transaction, it should just start DA submissio

	}

	err = UpdateProvenHeight(sdkCtx, store, endHeight)
	if err != nil {
		return err
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
