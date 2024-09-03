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
	fmt.Println("coming hereee.........", ctx.BlockHeight())

	currentBlockHeight := ctx.BlockHeight()
	if !k.IsValidBlockToPostTODA(uint64(currentBlockHeight)) {
		return nil
	}

	// fmt.Printf("Ctx.........%+v\n", ctx)

	fmt.Println("block heighttt.........", ctx.BlockHeight(), ctx.ExecMode(), ctx.IsCheckTx(), ctx.IsReCheckTx())

	provenHeight := k.GetProvenHeightFromStore(ctx)
	fromHeight := provenHeight + 1
	fmt.Println("from height..", fromHeight)
	endHeight := min(fromHeight+uint64(k.MaxBlocksForBlob), uint64(ctx.BlockHeight())) //exclusive i.e [fromHeight, endHeight)
	fmt.Println("end height..", endHeight-1)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := k.SetBlobStatusPending(sdkCtx, fromHeight, endHeight-1)
	if err != nil {
		fmt.Println("error while setting blob status...", err)
		return nil
	}

	var blocksToSumit []int64

	for i := fromHeight; i < endHeight; i++ {
		blocksToSumit = append(blocksToSumit, int64(i))
	}

	fmt.Println("blocks to submittttttttt.........", blocksToSumit)

	// only proposar should should run the this
	if bytes.Equal(req.ProposerAddress, k.proposerAddress) {
		// update blob status to success
		// err = k.SetBlobStatusSuccess(sdkCtx, fromHeight, endHeight)
		// if err != nil {
		// 	return nil
		// }

		// Todo: run the relayer routine
		// relayer doesn't have to make submitBlob Transaction, it should just start DA submission
		k.relayer.PostBlocks(ctx, blocksToSumit, k.cdc, req.ProposerAddress)

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
