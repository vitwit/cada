package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	govtestutil "github.com/cosmos/cosmos-sdk/x/gov/testutil"
	availblob1 "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/keeper"
)

func setAvailKeeper(t *testing.T) (
	*keeper.Keeper,
	*govtestutil.MockAccountKeeper,
	*govtestutil.MockBankKeeper,
	*govtestutil.MockStakingKeeper,
	moduletestutil.TestEncodingConfig,
	sdk.Context,
) {

	key := storetypes.NewKVStoreKey(availblob1.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()
	banktypes.RegisterInterfaces(encCfg.InterfaceRegistry)
}
