package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	cadakeeper "github.com/vitwit/avail-da-module/x/cada/keeper"
	availtypes "github.com/vitwit/avail-da-module/x/cada/types"
)

const (
	OpWeightMsgUpdateBlobStatusRequest = "op_weight_msg_update_blob_status"

	DefaultWeightMsgUpdateStatusRequest = 100
)

func WeightedOperations(
	appParams simtypes.AppParams, _ codec.JSONCodec, _ client.TxConfig,
	ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k cadakeeper.Keeper,
) simulation.WeightedOperations {
	var weightMsgUpdateBlobStatusRequest int
	appParams.GetOrGenerate(OpWeightMsgUpdateBlobStatusRequest, &weightMsgUpdateBlobStatusRequest, nil, func(_ *rand.Rand) {
		weightMsgUpdateBlobStatusRequest = DefaultWeightMsgUpdateStatusRequest
	})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgUpdateBlobStatusRequest,
			SimulateMsgUpdateBlobStatus(ak, bk, k),
		),
	}
}

func SimulateMsgUpdateBlobStatus(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k cadakeeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		ctx = ctx.WithBlockHeight(20)
		// Randomly select a sender account
		sender, _ := simtypes.RandomAcc(r, accs)

		// Ensure the sender has sufficient balance
		account := ak.GetAccount(ctx, sender.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		// Generate random fees for the transaction
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(availtypes.ModuleName, availtypes.TypeMsgUpdateBlobStatus, "unable to generate fees"), nil, err
		}

		// Prepare a random blob status update
		newStatus := true      // You can randomize this value as needed
		fromBlock := uint64(5) // Example block range start
		toBlock := uint64(20)  // Example block range end
		availHeight := uint64(120)

		ran := availtypes.Range{
			From: fromBlock,
			To:   toBlock,
		}

		msg := availtypes.NewMsgUpdateBlobStatus(
			sender.Address.String(),
			ran,
			availHeight,
			newStatus,
		)

		store := ctx.KVStore(k.GetStoreKey())
		cadakeeper.UpdateEndHeight(ctx, store, uint64(20))

		cadakeeper.UpdateProvenHeight(ctx, store, uint64(4))

		cadakeeper.UpdateBlobStatus(ctx, store, uint32(1))

		// Set up the transaction context
		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			Context:       ctx,
			SimAccount:    sender,
			AccountKeeper: ak,
			ModuleName:    availtypes.ModuleName,
		}

		// Generate and deliver the transaction
		return simulation.GenAndDeliverTx(txCtx, fees)
	}
}
