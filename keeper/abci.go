package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/types"
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

	// var latestProvenHeight int64 = 1
	// // TODO : set latestproven height in store
	// injectData := h.keeper.prepareInjectData(ctx, req.Time, latestProvenHeight)
	// injectDataBz := h.keeper.marshalMaxBytes(&injectData, req.MaxTxBytes, latestProvenHeight)
	// resp.Txs = h.keeper.addAvailblobDataToTxs(injectDataBz, req.MaxTxBytes, resp.Txs)

	return resp, nil
}

func (h *ProofOfBlobProposalHandler) ProcessProposal(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
	// // fmt.Println("length of transactions: ", len(req.Txs), ctx.BlockHeight())
	// injectedData := h.keeper.getInjectedData(req.Txs)
	// if injectedData != nil {
	// 	req.Txs = req.Txs[1:] // Pop the injected data for the default handler

	// 	if err := h.keeper.processPendingBlocks(ctx, req.Time, &injectedData.PendingBlocks); err != nil {
	// 		return nil, err
	// 	}
	// }
	return h.processProposalHandler(ctx, req)
}

func (k *Keeper) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) error {
	// injectedData := k.prepareInjectData(ctx, req.Time, req.Height)

	// injectDataBz := k.marshalMaxBytes(&injectedData, int64(req.Size()), req.Height)
	// _ = k.addAvailblobDataToTxs(injectDataBz, int64(req.Size()), req.Txs)

	// if err := k.preblockerPendingBlocks(ctx, req.Time, req.ProposerAddress, &injectedData.PendingBlocks); err != nil {
	// 	return err
	// }
	// // }
	return nil
}

func (k *Keeper) getInjectedData(txs [][]byte) *types.InjectedData {
	if len(txs) != 0 {
		var injectedData types.InjectedData
		err := k.cdc.Unmarshal(txs[0], &injectedData)
		if err == nil {
			return &injectedData
		}
	}
	return nil
}
