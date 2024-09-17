package relayer

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/types"

	"github.com/cosmos/cosmos-sdk/codec"
	dacli "github.com/vitwit/avail-da-module/chainclient"
)

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
func (r *Relayer) PostBlocks(ctx sdk.Context, blocks []int64, cdc codec.BinaryCodec, proposer []byte) {
	go r.postBlocks(ctx, blocks, cdc, proposer)
}

func (r *Relayer) GetBlocksDataFromLocal(ctx sdk.Context, blocks []int64) []byte {
	if len(blocks) == 0 {
		return []byte{}
	}

	var bb []byte

	for _, height := range blocks {
		res, err := r.localProvider.GetBlockAtHeight(ctx, height)
		if err != nil {
			r.logger.Error("Error getting block", "height:", height, "error", err)
			return []byte{}
		}

		blockProto, err := res.Block.ToProto()
		if err != nil {
			r.logger.Error("Error protoing block", "error", err)
			return []byte{}
		}

		blockBz, err := blockProto.Marshal()
		if err != nil {
			r.logger.Error("Error marshaling block", "error", err)
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

	blockInfo, err := r.SubmitDataToAvailClient(r.rpcClient.config.Seed, r.rpcClient.config.AppID, bb, blocks, r.rpcClient.config.LightClientURL)

	if err != nil {
		r.logger.Error("Error while submitting block(s) to Avail DA",
			"height_start", blocks[0],
			"height_end", blocks[len(blocks)-1],
			"appID", string(r.rpcClient.config.AppID),
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
		}, cdc)
		if err != nil {
			fmt.Println("error while submitting tx...", err)
		}

		return
	}

	if blockInfo.BlockNumber != 0 {
		msg := types.MsgUpdateBlobStatusRequest{ValidatorAddress: sdk.AccAddress.String(proposer),
			BlocksRange: &types.Range{
				From: uint64(blocks[0]),
				To:   uint64(blocks[len(blocks)-1]),
			},
			AvailHeight: uint64(blockInfo.BlockNumber),
			IsSuccess:   true}

		// TODO : execute tx about successfull submission
		err = dacli.ExecuteTX(ctx, msg, cdc)
		if err != nil {
			fmt.Println("error while submitting tx...", err)
		}
	}
}
