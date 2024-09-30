package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/vitwit/avail-da-module/x/cada/keeper"
	availtypes "github.com/vitwit/avail-da-module/x/cada/types"
)

const (
	OpWeightMsgUpdateBlobStatusRequest = "op_weight_msg_update_blob_status"

	DefaultWeightMsgUpdateStatusRequest = 100
)

func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, txConfig client.TxConfig,
	ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k keeper.Keeper,
) simulation.WeightedOperations {
	var weightMsgUpdateBlobStatusRequest int
	appParams.GetOrGenerate(OpWeightMsgUpdateBlobStatusRequest, &weightMsgUpdateBlobStatusRequest, nil, func(_ *rand.Rand) {
		weightMsgUpdateBlobStatusRequest = DefaultWeightMsgUpdateStatusRequest
	})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgUpdateBlobStatusRequest,
			SimulateUpdateBlobStatus(txConfig, cdc, ak, bk, k),
		),
	}
}

func SimulateUpdateBlobStatus(txConfig client.TxConfig, _ codec.JSONCodec, ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper, _ keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simaAccount, _ := simtypes.RandomAcc(r, accs)

		if simaAccount.Address.Empty() {
			fmt.Printf("simaAccount address is empty\n")
			return simtypes.NoOpMsg(availtypes.ModuleName, availtypes.TypeMsgUpdateBlobStatus, "simaAccount address is empty"), nil, nil
		}

		// Ensure TxGen is properly initialized
		if txConfig == nil {
			return simtypes.NoOpMsg(availtypes.ModuleName, availtypes.TypeMsgUpdateBlobStatus, "TxGen is nil"), nil, nil
		}

		// account := ak.GetAccount(ctx, simaAccount.Address)
		// if account == nil {
		// 	fmt.Printf("account is nil for address: %s\n", simaAccount.Address.String())
		// 	return simtypes.NoOpMsg(availtypes.ModuleName, availtypes.TypeMsgUpdateBlobStatus, "account is nil"), nil, nil
		// }

		// fmt.Printf("account: %v\n", account)

		// pubKey := account.GetPubKey()
		// if pubKey == nil {
		// 	fmt.Printf("account has no pubkey: %s\n", account.GetAddress().String())
		// 	return simtypes.NoOpMsg(availtypes.ModuleName, availtypes.TypeMsgUpdateBlobStatus, "account has no pubkey"), nil, nil
		// }

		//fmt.Printf("account.GetPubKey(): %v\n", pubKey)

		fromBlock := r.Uint64()%100 + 1
		toBlock := fromBlock + r.Uint64()%10

		isSuccess := r.Intn(2) == 1

		availHeight := r.Uint64()%100 + 1

		blockRange := availtypes.Range{
			From: fromBlock,
			To:   toBlock,
		}

		fmt.Printf("fromBlock: %v\n", fromBlock)
		fmt.Printf("toBlock: %v\n", toBlock)
		fmt.Printf("availHeight: %v\n", availHeight)

		// Fetch spendable coins to simulate transaction fees (even if just dummy fees)
		spendable := bk.SpendableCoins(ctx, simaAccount.Address)
		fmt.Printf("spendable:::::::::::::::::::::::: %v\n", spendable)
		if spendable.Empty() {
			return simtypes.NoOpMsg(availtypes.ModuleName, availtypes.TypeMsgUpdateBlobStatus, "account has no spendable coins"), nil, nil
		}

		fmt.Printf("txConfig:::::::::::::::::::::::::::::: %v\n", txConfig)

		fmt.Printf("simaAccount.Address.String(): %v\n", simaAccount.Address.String())

		msg := availtypes.NewMsgUpdateBlobStatus(simaAccount.Address.String(), blockRange, availHeight, isSuccess)

		fmt.Printf("r: %v\n", r)
		fmt.Printf("app: %v\n", app)
		fmt.Printf("msg: %v\n", msg)
		fmt.Printf("ctx: %v\n", ctx)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           txConfig,
			Msg:             msg,
			Context:         ctx,
			SimAccount:      simaAccount,
			ModuleName:      availtypes.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		opMsg, futureOps, err := simulation.GenAndDeliverTxWithRandFees(txCtx)
		if err != nil {
			fmt.Printf("errrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr: %v\n", err)
			return simtypes.NoOpMsg(availtypes.ModuleName, availtypes.TypeMsgUpdateBlobStatus, "error generating or delivering tx"), nil, err
		}

		return opMsg, futureOps, nil

		// Generate and deliver the transaction
		//return simulation.GenAndDeliverTxWithRandFees(txCtx)

		//return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
