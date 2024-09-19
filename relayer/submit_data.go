package relayer

import (
	"encoding/base64"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/relayer/avail"
)

func (r *Relayer) SubmitDataToAvailClient(data []byte, blocks []int64, lightClientURL string) (avail.AvailBlockMetaData, error) {

	datab := base64.StdEncoding.EncodeToString(data)
	jsonData := []byte(fmt.Sprintf(`{"data":"%s"}`, datab))

	blockInfo, err := r.availDAClient.Submit(jsonData)

	if err != nil {
		r.logger.Error("Error while posting block(s) to Avail DA",
			"height_start", blocks[0],
			"height_end", blocks[len(blocks)-1],
			"appID", strconv.Itoa(r.AvailConfig.AppID),
		)

		return blockInfo, err
	}

	r.logger.Info("Successfully posted block(s) to Avail DA",
		"height_start", blocks[0],
		"height_end", blocks[len(blocks)-1],
		"appID", strconv.Itoa(r.AvailConfig.AppID),
		"block_hash", blockInfo.BlockHash,
		"block_number", blockInfo.BlockNumber,
		"hash", blockInfo.Hash,
	)

	return blockInfo, nil
}

// IsDataAvailable is to query the avail light client and check if the data is made available at the given height
func (r *Relayer) IsDataAvailable(ctx sdk.Context, from, to, availHeight uint64, lightClientURL string) (bool, error) {

	var blocks []int64
	for i := from; i <= to; i++ {
		blocks = append(blocks, int64(i))
	}

	cosmosBlocksData := r.GetBlocksDataFromLocal(ctx, blocks)

	return r.availDAClient.IsDataAvailable(cosmosBlocksData, int(availHeight))
}

// Define the struct that matches the JSON structure
type GetBlock struct {
	BlockNumber      int      `json:"block_number"`
	DataTransactions []string `json:"data_transactions"`
}
