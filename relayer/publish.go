package relayer

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dacli "github.com/vitwit/avail-da-module/chainclient"
	"github.com/vitwit/avail-da-module/types"
)

func (r *Relayer) ProposePostNextBlocks(ctx sdk.Context, provenHeight int64) []int64 {
	height := ctx.BlockHeight()

	if height <= 1 {
		return nil
	}

	// only publish new blocks on interval
	if (height-1)%int64(r.AvailConfig.PublishBlobInterval) != 0 {
		return nil
	}

	var blocks []int64
	for block := height - int64(r.AvailConfig.PublishBlobInterval); block < height; block++ {
		// this could be false after a genesis restart
		if block > provenHeight {
			blocks = append(blocks, block)
		}
	}

	return blocks
}

// PostBlocks is called in the PreBlocker. The proposer will publish the blocks at this point
// once their block has been accepted. The method launches the posting process in a separate
// goroutine to handle the submission of blocks asynchronously.
func (r *Relayer) PostBlocks(ctx sdk.Context, blocks []int64, cdc codec.BinaryCodec, proposer []byte) {
	go r.postBlocks(ctx, blocks, cdc, proposer)
}

// GetBlocksDataFromLocal retrieves block data from the local provider for the specified block heights.
func (r *Relayer) GetBlocksDataFromLocal(ctx sdk.Context, blocks []int64) []byte {
	if len(blocks) == 0 {
		return []byte{}
	}

	var bb []byte

	for _, height := range blocks {
		res, err := r.LocalProvider.GetBlockAtHeight(ctx, height)
		if err != nil {
			r.Logger.Error("Error getting block", "height:", height, "error", err)
			return []byte{}
		}

		blockProto, err := res.Block.ToProto()
		if err != nil {
			r.Logger.Error("Error protoing block", "error", err)
			return []byte{}
		}

		blockBz, err := blockProto.Marshal()
		if err != nil {
			r.Logger.Error("Error marshaling block", "error", err)
			return []byte{}
		}

		bb = append(bb, blockBz...)
	}

	return bb
}

// postBlocks will publish rollchain blocks to avail
// start height is inclusive, end height is exclusive
func (r *Relayer) postBlocks(ctx sdk.Context, blocks []int64, cdc codec.BinaryCodec, proposer []byte) {
	// process blocks instead of random data
	if len(blocks) == 0 {
		return
	}

	bb := r.GetBlocksDataFromLocal(ctx, blocks)

	blockInfo, err := r.SubmitDataToAvailClient(bb, blocks)
	if err != nil {
		r.Logger.Error("Error while submitting block(s) to Avail DA",
			"height_start", blocks[0],
			"height_end", blocks[len(blocks)-1],
			"appID", strconv.Itoa(r.AvailConfig.AppID), err,
		)

		// execute tx about failure submission
		err = dacli.ExecuteTX(ctx, types.MsgUpdateBlobStatusRequest{
			ValidatorAddress: sdk.AccAddress.String(proposer),
			BlocksRange: &types.Range{
				From: uint64(blocks[0]),
				To:   uint64(blocks[len(blocks)-1]),
			},
			// AvailHeight: uint64(blockInfo.BlockNumber),
			IsSuccess: false,
		}, cdc, r.AvailConfig, r.NodeDir)
		if err != nil {
			fmt.Println("error while submitting tx...", err)
		}

		return
	}

	if blockInfo.BlockNumber != 0 {
		msg := types.MsgUpdateBlobStatusRequest{
			ValidatorAddress: sdk.AccAddress.String(proposer),
			BlocksRange: &types.Range{
				From: uint64(blocks[0]),
				To:   uint64(blocks[len(blocks)-1]),
			},
			AvailHeight: uint64(blockInfo.BlockNumber),
			IsSuccess:   true,
		}

		// execute tx about successful submission
		err = dacli.ExecuteTX(ctx, msg, cdc, r.AvailConfig, r.NodeDir)
		if err != nil {
			fmt.Println("error while submitting tx...", err)
		}
	}
}
