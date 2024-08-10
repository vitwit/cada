package relayer

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CelestiaPublishKeyName  = "blob"
	CelestiaFeegrantKeyName = "feegrant"
	celestiaBech32Prefix    = "celestia"
	celestiaBlobPostMemo    = "Posted by tiablob https://rollchains.com"
)

// PostNextBlocks is called by the current proposing validator during PrepareProposal.
// If on the publish boundary, it will return the block heights that will be published
// It will not publish the block being proposed.
func (r *Relayer) ProposePostNextBlocks(ctx sdk.Context, provenHeight int64) []int64 {
	height := ctx.BlockHeight()
	fmt.Println("publish block intervallllllllllll............", r.availPublishBlockInterval)

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
func (r *Relayer) PostBlocks(ctx sdk.Context, blocks []int64) {
	go r.postBlocks(ctx, blocks)
}

// postBlocks will publish rollchain blocks to avail
// start height is inclusive, end height is exclusive
func (r *Relayer) postBlocks(ctx sdk.Context, blocks []int64) {
	// process blocks instead of random data
	if len(blocks) == 0 {
		// fmt.Println("Empty blocks dataa...", len(blocks))
		return
	}

	var bb []byte

	for _, height := range blocks {
		res, err := r.localProvider.GetBlockAtHeight(ctx, height)
		if err != nil {
			r.logger.Error("Error getting block", "height:", height, "error", err)
			return
		}

		// fmt.Println("blocks res...........", res)

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

	// fmt.Println("block dataaa.......", ctx.)
	// hash, err := r.rpcClient.SubmitData(r.rpcClient.config.AppRpcURL, r.rpcClient.config.Seed, r.rpcClient.config.AppID, bb)
	// if err == nil {
	// 	r.logger.Info("Posted block(s) to Avail DA",
	// 		"height_start", blocks[0],
	// 		"height_end", blocks[len(blocks)-1],
	// 		"appID", string(r.rpcClient.config.AppID),
	// 		"hash", hash,
	// 		// "tx_hash", hex.EncodeToString(res.Hash),
	// 		// "url", fmt.Sprintf("https://mocha.celenium.io/tx/%s", hex.EncodeToString(res.Hash)),
	// 	)
	// 	// return
	// }

	err := r.SubmitDataToClient(r.rpcClient.config.AppRpcURL, r.rpcClient.config.Seed, r.rpcClient.config.AppID, bb, blocks)
	if err != nil {
		r.logger.Error("Error while submitting block(s) to Avail DA",
			"height_start", blocks[0],
			"height_end", blocks[len(blocks)-1],
			"appID", string(r.rpcClient.config.AppID),
		)
	}
}
