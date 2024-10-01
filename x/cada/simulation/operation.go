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
	cadatypes "github.com/vitwit/avail-da-module/x/cada/types"
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
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, _ string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		ctx = ctx.WithBlockHeight(20)

		sender, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, sender.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(cadatypes.ModuleName, cadatypes.TypeMsgUpdateBlobStatus, "unable to generate fees"), nil, err
		}

		newStatus := true
		fromBlock := uint64(5)
		toBlock := uint64(20)
		availHeight := uint64(120)

		ran := cadatypes.Range{
			From: fromBlock,
			To:   toBlock,
		}

		msg := cadatypes.NewMsgUpdateBlobStatus(
			sender.Address.String(),
			ran,
			availHeight,
			newStatus,
		)

		store := ctx.KVStore(k.GetStoreKey())

		err = cadakeeper.UpdateEndHeight(ctx, store, uint64(20))
		if err != nil {
			return simtypes.NoOpMsg(cadatypes.ModuleName, cadatypes.TypeMsgUpdateBlobStatus, "unable to update end height"), nil, err
		}

		err = cadakeeper.UpdateProvenHeight(ctx, store, uint64(4))
		if err != nil {
			return simtypes.NoOpMsg(cadatypes.ModuleName, cadatypes.TypeMsgUpdateBlobStatus, "unable to update proven height"), nil, err
		}

		err = cadakeeper.UpdateBlobStatus(ctx, store, uint32(1))
		if err != nil {
			return simtypes.NoOpMsg(cadatypes.ModuleName, cadatypes.TypeMsgUpdateBlobStatus, "unable to update status to pending state"), nil, err
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			Context:       ctx,
			SimAccount:    sender,
			AccountKeeper: ak,
			ModuleName:    cadatypes.ModuleName,
		}

		return simulation.GenAndDeliverTx(txCtx, fees)
	}
}
