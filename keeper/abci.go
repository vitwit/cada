package keeper

import (
	"fmt"
	"log"

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
	// fmt.Println("proposal err.........", err)
	// log.Fatal("errrorrrrr........", err)
	if err != nil {
		log.Fatal("error hereeeeeeeeeeeeeeeeee", err)
		return nil, err
	}

	// latestProvenHeight, err := h.keeper.GetProvenHeight(ctx)
	// fmt.Println("eerorrrrrrrrrrrr...........", err)
	// if err != nil {
	// 	return nil, err
	// }

	var latestProvenHeight int64 = 2
	// TODO : set latestproven height in store
	injectData := h.keeper.prepareInjectData(ctx, req.Time, latestProvenHeight)
	// fmt.Println("11111111111")
	injectDataBz := h.keeper.marshalMaxBytes(&injectData, req.MaxTxBytes, latestProvenHeight)
	// fmt.Println("2222222222222")
	resp.Txs = h.keeper.addAvailblobDataToTxs(injectDataBz, req.MaxTxBytes, resp.Txs)

	// fmt.Println("injected in prepareproposal.........", injectData)

	return resp, nil
}

func (h *ProofOfBlobProposalHandler) ProcessProposal(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
	injectedData := h.keeper.getInjectedData(req.Txs)
	if injectedData != nil {
		req.Txs = req.Txs[1:] // Pop the injected data for the default handler
		// if err := h.keeper.processCreateClient(ctx, injectedData.CreateClient); err != nil {
		// 	return nil, err
		// }
		// if err := h.keeper.processHeaders(ctx, injectedData.Headers); err != nil {
		// 	return nil, err
		// }
		// if err := h.keeper.processProofs(ctx, injectedData.Headers, injectedData.Proofs); err != nil {
		// 	return nil, err
		// }
		if err := h.keeper.processPendingBlocks(ctx, req.Time, &injectedData.PendingBlocks); err != nil {
			return nil, err
		}
	}
	return h.processProposalHandler(ctx, req)
}

func (k *Keeper) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) error {
	// injectedData := k.getInjectedData(req.Txs)
	// fmt.Println("above injected dataaa.........", req.Time, req.Height, ctx)
	injectedData := k.prepareInjectData(ctx, req.Time, req.Height)

	// fmt.Println("after injected dataaa.........", injectedData)

	// fmt.Println("111111111")
	injectDataBz := k.marshalMaxBytes(&injectedData, int64(req.Size()), req.Height)
	// fmt.Println("2222222222")
	_ = k.addAvailblobDataToTxs(injectDataBz, int64(req.Size()), req.Txs)
	// fmt.Println("resp hereee...........", resp)

	// fmt.Println("injected dataaa in preblockerrrr.........", injectedData)
	// fmt.Println("injected dataaa.........", injectedData, req.Txs)
	// if injectedData != nil {
	// if err := k.preblockerCreateClient(ctx, injectedData.CreateClient); err != nil {
	// 	return err
	// }
	// if err := k.preblockerHeaders(ctx, injectedData.Headers); err != nil {
	// 	return err
	// }
	// if err := k.preblockerProofs(ctx, injectedData.Proofs); err != nil {
	// 	return err
	// }

	fmt.Println("inside pre blocker")
	// fmt.Println("req", req)
	if err := k.preblockerPendingBlocks(ctx, req.Time, req.ProposerAddress, &injectedData.PendingBlocks); err != nil {
		return err
	}
	// }
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
