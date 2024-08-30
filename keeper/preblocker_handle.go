package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/types"
)

func (k *Keeper) preblockerPendingBlocks(ctx sdk.Context, blockTime time.Time, proposerAddr []byte, pendingBlocks *types.PendingBlocks) error {
	if pendingBlocks != nil {
		// if reflect.DeepEqual(k.proposerAddress, proposerAddr) {
		// 	k.relayer.PostBlocks(ctx, pendingBlocks.BlockHeights)
		// }

		for _, pendingBlock := range pendingBlocks.BlockHeights {
			if err := k.AddUpdatePendingBlock(ctx, pendingBlock, blockTime); err != nil {
				return fmt.Errorf("preblocker pending blocks, %v", err)
			}
		}
	}

	return nil
}

func (k *Keeper) notifyProvenHeight(ctx sdk.Context, previousProvenHeight int64) {

	// TODO
	// provenHeight, err := k.GetProvenHeight(ctx)
	// if err != nil {
	// 	fmt.Println("unable to get proven height", err)
	// 	return
	// }

	// //k.relayer.NotifyProvenHeight(provenHeight)
	// go k.relayer.PruneBlockStore(previousProvenHeight)
}
