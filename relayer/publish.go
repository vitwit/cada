package relayer

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/types"

	"github.com/cosmos/cosmos-sdk/codec"
	dacli "github.com/vitwit/avail-da-module/chainclient"
)

// PostNextBlocks is called by the current proposing validator during PrepareProposal.
// If on the publish boundary, it will return the block heights that will be published
// It will not publish the block being proposed.

func (r *Relayer) NextBlocksToSumbit(ctx sdk.Context) (types.MsgSubmitBlobRequest, bool) {
	height := ctx.BlockHeight() // TODO: not using deprecate
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
func (r *Relayer) ProposePostNextBlocks(ctx sdk.Context, provenHeight int64) []int64 { // TODO: deprecated code
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

	fmt.Println("is it coming here where we post to DA")

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
		fmt.Println("proposer addressss........", sdk.AccAddress.String(proposer),
			uint64(blocks[0]), uint64(blocks[len(blocks)-1]), uint64(blockInfo.BlockNumber))

		msg := types.MsgUpdateBlobStatusRequest{ValidatorAddress: sdk.AccAddress.String(proposer),
			BlocksRange: &types.Range{
				From: uint64(blocks[0]),
				To:   uint64(blocks[len(blocks)-1]),
			},
			AvailHeight: uint64(blockInfo.BlockNumber),
			IsSuccess:   true}

		fmt.Println("submit blocks msg.......", msg)

		// TODO : execute tx about successfull submission
		err = dacli.ExecuteTX(ctx, msg, cdc)
		if err != nil {
			fmt.Println("error while submitting tx...", err)
		}
	}
}

// var availdHomePath = xfilepath.JoinFromHome(xfilepath.Path("availsdk"))

// availdHomePath := filepath.Join(os.Getenv("HOME"), "availsdk")

// func GetBinPath() string {
// 	homeDir, err := os.UserHomeDir()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	availdHomePath := filepath.Join(homeDir, ".availsdk")
// 	fmt.Println("availdHonmePath.......", availdHomePath)
// 	return availdHomePath
// }

// func ExecuteTX(ctx sdk.Context, msg types.MsgUpdateBlobStatusRequest, cdc codec.BinaryCodec) error {
// 	// Define keyring and RPC client configuration

// 	// homePath := "/home/vitwit/.availsdk"
// 	homePath := GetBinPath()
// 	fmt.Println("get key namee.........", os.Getenv("KEY"))
// 	keyName := os.Getenv("KEY")
// 	rpcAddress := "http://localhost:26657"

// 	// Create a keyring
// 	kr, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, homePath, os.Stdin, cdc.(codec.Codec))
// 	if err != nil {
// 		return fmt.Errorf("error creating keyring: %w", err)
// 	}

// 	// List all keys in the keyring
// 	// keys, err := kr.List()
// 	// if err != nil {
// 	// 	fmt.Println("error listing keys:", err)
// 	// }

// 	info, err := kr.Key(keyName)
// 	valAddr, err := info.GetAddress()
// 	fmt.Println("after address................", valAddr)

// 	// valAddr, err := sdk.AccAddressFromBech32(addr.String())
// 	// fmt.Println("val addr, err..", valAddr, err, addr)

// 	// fmt.Println("keysss........", keys)

// 	// // Print out the keys
// 	// for _, keyInfo := range keys {
// 	// 	addr, err := keyInfo.GetAddress()
// 	// 	fmt.Println("err..", err)
// 	// 	fmt.Printf("Name: %s, Address: %s\n", keyInfo.Name, addr)
// 	// }

// 	// Create an RPC client
// 	rpcClient, err := cometrpc.NewWithTimeout(rpcAddress, "/websocket", 3)
// 	if err != nil {
// 		return fmt.Errorf("error creating RPC client: %w", err)
// 	}

// 	// Create a new client context
// 	clientCtx := NewClientCtx(kr, rpcClient, ctx.ChainID(), cdc, homePath, valAddr)

// 	// Retrieve the validator address (replace with actual logic to get the address)
// 	// valAddr, err = sdk.AccAddressFromBech32("cosmos1fhqer4tc50nut2evvnj6yegcah2yfu3s844n9a")
// 	// if err != nil {
// 	// 	return fmt.Errorf("error parsing validator address: %w", err)
// 	// }

// 	// Set the client context's from fields
// 	clientCtx.FromName = keyName
// 	clientCtx.FromAddress = valAddr

// 	// Fetch account number and sequence from the blockchain
// 	accountRetriever := authtypes.AccountRetriever{}
// 	account, err := accountRetriever.GetAccount(clientCtx, valAddr)
// 	if err != nil {
// 		return fmt.Errorf("error retrieving account: %w", err)
// 	}

// 	fmt.Println("account details......", account.GetAccountNumber(), account.GetSequence())

// 	// Set the correct account number and sequence
// 	factory := NewFactory(clientCtx).
// 		WithAccountNumber(account.GetAccountNumber()).
// 		WithSequence(account.GetSequence())

// 	// Create a transaction factory and set the validator address in the message
// 	// factory := NewFactory(clientCtx)
// 	msg.ValidatorAddress = valAddr.String()
// 	// time.Sleep(10 * time.Second)

// 	// Generate and broadcast the transaction
// 	if err := clitx.GenerateOrBroadcastTxWithFactory(clientCtx, factory, &msg); err != nil {
// 		return fmt.Errorf("error broadcasting transaction: %w", err)
// 	}

// 	return nil
// }
