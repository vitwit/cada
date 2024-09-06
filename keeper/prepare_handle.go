package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/vitwit/avail-da-module/types"
)

const DelayAfterUpgrade = int64(10)

// TODO: not using irrelevant
func (k *Keeper) prepareInjectData(ctx sdk.Context, currentBlockTime time.Time, latestProvenHeight int64) types.InjectedData {
	// return types.InjectedData{
	// 	PendingBlocks: k.preparePostBlocks(ctx, currentBlockTime),
	// }
	return types.InjectedData{}
}

// TODO: not using irrelevant
func (k *Keeper) addAvailblobDataToTxs(injectDataBz []byte, maxTxBytes int64, txs [][]byte) [][]byte {
	if injectDataBz != nil && len(injectDataBz) > 0 {
		var finalTxs [][]byte
		totalTxBytes := int64(len(injectDataBz))
		finalTxs = append(finalTxs, injectDataBz)
		for _, tx := range txs {
			totalTxBytes += int64(len(tx))
			// Append as many transactions as will fit
			if totalTxBytes <= maxTxBytes {
				finalTxs = append(finalTxs, tx)
			} else {
				break
			}
		}
		return finalTxs
	}
	return txs
}

// TODO: not using irrelevant
func (k *Keeper) preparePostBlocks(ctx sdk.Context, currentBlockTime time.Time) types.PendingBlocks {
	// provenHeight, err := k.GetProvenHeight(ctx)
	// if err != nil {
	// 	return types.PendingBlocks{}
	// }
	// newBlocks := k.relayer.ProposePostNextBlocks(ctx, provenHeight)

	// If there are no new blocks to propose, check for expired blocks
	// Additionally, if the block interval is 1, we need to also be able to re-publish an expired block
	// if len(newBlocks) < 2 && k.shouldGetExpiredBlock(ctx) {
	// 	expiredBlocks := k.GetExpiredBlocks(ctx, currentBlockTime)
	// 	for _, expiredBlock := range expiredBlocks {
	// 		// Check if we have a proof for this block already
	// 		if k.relayer.HasCachedProof(expiredBlock) {
	// 			continue
	// 		}
	// 		// Add it to the list respecting the publish limit
	// 		if len(newBlocks) < 2 || len(newBlocks) < k.publishToAvailBlockInterval {
	// 			newBlocks = append(newBlocks, expiredBlock)
	// 		}
	// 	}
	// }

	// return types.PendingBlocks{
	// 	BlockHeights: newBlocks,
	// }
	return types.PendingBlocks{}
}

// shouldGetExpiredBlocks checks if this chain has recently upgraded.    // TODO: not using irrelevant
// If so, it will delay publishing expired blocks so that the relayer has time to populate block proof cache first
func (k *Keeper) shouldGetExpiredBlock(ctx sdk.Context) bool {
	_, lastUpgradeHeight, _ := k.upgradeKeeper.GetLastCompletedUpgrade(ctx)
	return ctx.BlockHeight() >= lastUpgradeHeight+DelayAfterUpgrade
}

// TODO: not using irrelevant
func (k *Keeper) marshalMaxBytes(injectData *types.InjectedData, maxBytes int64, latestProvenHeight int64) []byte {
	if len(injectData.PendingBlocks.BlockHeights) == 0 {
		return nil
	}

	injectDataBz, err := k.cdc.Marshal(injectData)
	if err != nil {
		return nil
	}

	proofLimit := k.injectedProofsLimit
	for int64(len(injectDataBz)) > maxBytes {
		proofLimit = proofLimit - 1

		injectDataBz, err = k.cdc.Marshal(injectData)
		if err != nil {
			return nil
		}
	}

	return injectDataBz
}
