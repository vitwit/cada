package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/types"
)

func (k Keeper) processPendingBlocks(ctx sdk.Context, currentBlockTime time.Time, pendingBlocks *types.PendingBlocks) error {
	if pendingBlocks != nil {
		height := ctx.BlockHeight()
		numBlocks := len(pendingBlocks.BlockHeights)
		if numBlocks > 2 && numBlocks > k.publishToAvailBlockInterval {
			return fmt.Errorf("process pending blocks, included pending blocks (%d) exceeds limit (%d)", numBlocks, k.publishToAvailBlockInterval)
		}
		for _, pendingBlock := range pendingBlocks.BlockHeights {
			if pendingBlock <= 0 {
				return fmt.Errorf("process pending blocks, invalid block: %d", pendingBlock)
			}
			if pendingBlock >= height {
				return fmt.Errorf("process pending blocks, start (%d) cannot be >= this block height (%d)", pendingBlock, height)
			}
			// Check if already pending, if so, is it expired?
			if k.IsBlockPending(ctx, pendingBlock) && !k.IsBlockExpired(ctx, currentBlockTime, pendingBlock) {
				return fmt.Errorf("process pending blocks, block height (%d) is pending, but not expired", pendingBlock)
			}
			// Check if we have a proof for this block
			// if k.relayer.HasCachedProof(pendingBlock) {
			// 	return fmt.Errorf("process pending blocks, cached proof exists for block %d", pendingBlock)
			// }
		}
		// Ensure publish boundries includes new blocks, once they are on-chain, they will be tracked appropriately
		provenHeight, err := k.GetProvenHeight(ctx)
		if err != nil {
			return fmt.Errorf("process pending blocks, getting proven height, %v", err)
		}
		newBlocks := k.relayer.ProposePostNextBlocks(ctx, provenHeight)
		for i, newBlock := range newBlocks {
			if newBlock != pendingBlocks.BlockHeights[i] {
				return fmt.Errorf("process pending blocks, block (%d) must be included", newBlock)
			}
		}
	}

	return nil
}
