package relayer

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/types"
)

// PostNextBlocks is called by the current proposing validator during PrepareProposal.
// If on the publish boundary, it will return the block heights that will be published
// It will not publish the block being proposed.

func (r *Relayer) NextBlocksToSumbit(ctx sdk.Context) (types.MsgSubmitBlobRequest, bool) {
	height := ctx.BlockHeight()
	// only publish new blocks on interval
	if height < 2 || (height-1)%int64(r.availPublishBlockInterval) != 0 {
		return types.MsgSubmitBlobRequest{}, false
	}

	return types.MsgSubmitBlobRequest{
		BlocksRange: &types.Range{
			From: uint64(height - int64(r.availPublishBlockInterval)),
			To:   uint64(height - 1),
		},
	}, true

}
func (r *Relayer) ProposePostNextBlocks(ctx sdk.Context, provenHeight int64) []int64 {
	height := ctx.BlockHeight()

	if height <= 1 {
		return nil
	}

	// only publish new blocks on interval
	if (height-1)%int64(r.availPublishBlockInterval) != 0 {
		return nil
	}

	var blocks []int64
	for block := height - int64(r.availPublishBlockInterval); block < height; block++ {
		// this could be false after a genesis restart
		if block > provenHeight {
			blocks = append(blocks, block)
		}
	}

	return blocks
}

// PostBlocks is call in the preblocker, the proposer will publish at this point with their block accepted
// func (r *Relayer) PostBlocks(ctx sdk.Context, blocks []int64) {
// 	go r.postBlocks(ctx, blocks)
// }

func (r *Relayer) PostBlocks(ctx sdk.Context, fromHeight uint64, endHeight uint64) {
	go r.postBlockss(ctx, fromHeight, endHeight)
}

// postBlocks will publish rollchain blocks to avail
// start height is inclusive, end height is exclusive
// func (r *Relayer) postBlocks(ctx sdk.Context, blocks []int64) {
// 	// process blocks instead of random data
// 	if len(blocks) == 0 {
// 		return
// 	}

// 	var bb []byte

// 	for _, height := range blocks {
// 		res, err := r.localProvider.GetBlockAtHeight(ctx, height)
// 		if err != nil {
// 			r.logger.Error("Error getting block", "height:", height, "error", err)
// 			return
// 		}

// 		blockProto, err := res.Block.ToProto()
// 		if err != nil {
// 			r.logger.Error("Error protoing block", "error", err)
// 			return
// 		}

// 		blockBz, err := blockProto.Marshal()
// 		if err != nil {
// 			r.logger.Error("Error marshaling block", "error", err)
// 			return
// 		}

// 		bb = append(bb, blockBz...)
// 	}

// 	err := r.SubmitDataToClient(r.rpcClient.config.Seed, r.rpcClient.config.AppID, bb, blocks, r.rpcClient.config.LightClientURL)
// 	if err != nil {
// 		r.logger.Error("Error while submitting block(s) to Avail DA",
// 			"height_start", blocks[0],
// 			"height_end", blocks[len(blocks)-1],
// 			"appID", string(r.rpcClient.config.AppID),
// 		)
// 	}
// }

func (r *Relayer) postBlockss(ctx sdk.Context, fromHeight uint64, endHeight uint64) {
	// process blocks instead of random data
	// if len(blocks) == 0 {
	// 	return
	// }

	var bb []byte

	for height := fromHeight; height <= endHeight; height++ {
		fmt.Printf("height: %v\n", height)
		res, err := r.localProvider.GetBlockAtHeight(ctx, int64(height))
		if err != nil {
			r.logger.Error("Error getting block", "height:", height, "error", err)
			return
		}

		blockProto, err := res.Block.ToProto()
		if err != nil {
			r.logger.Error("Error protoing block", "error", err)
			return
		}

		blockBz, err := blockProto.Marshal()
		if err != nil {
			r.logger.Error("Error marshaling block", "error", err)
			return
		}

		bb = append(bb, blockBz...)
	}

	err := r.SubmitDataToClient(r.rpcClient.config.Seed, r.rpcClient.config.AppID, bb, int64(fromHeight), int64(endHeight), r.rpcClient.config.LightClientURL)
	if err != nil {
		r.logger.Error("Error while submitting block(s) to Avail DA",
			"height_start", fromHeight,
			"height_end", endHeight,
			"appID", string(r.rpcClient.config.AppID),
		)
	}

}
