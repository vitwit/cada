package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtestutil "github.com/cosmos/cosmos-sdk/x/gov/testutil"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/test-go/testify/require"
	"github.com/test-go/testify/suite"
	availblob1 "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/keeper"
	AvailBlobKeeper "github.com/vitwit/avail-da-module/keeper"
	availdapp "github.com/vitwit/avail-da-module/simapp/app"
	"github.com/vitwit/avail-da-module/types"
)

type KeeperTestSuite struct {
	suite.Suite

	cdc           codec.Codec
	ctx           sdk.Context
	govKeeper     *keeper.Keeper
	acctKeeper    *govtestutil.MockAccountKeeper
	bankKeeper    *govtestutil.MockBankKeeper
	stakingKeeper *govtestutil.MockStakingKeeper
	distKeeper    *govtestutil.MockDistributionKeeper
	addrs         []sdk.AccAddress
	msgServer     types.MsgServer
	legacyMsgSrvr v1beta1.MsgServer
}

func (suite *KeeperTestSuite) SetupSuite() {
	suite.reset()
}

func (suite *KeeperTestSuite) reset() {
	govKeeper, acctKeeper, bankKeeper, stakingKeeper, distKeeper, encCfg, ctx := setupGovKeeper(suite.T())

	// Populate the gov account with some coins, as the TestProposal we have
	// is a MsgSend from the gov account.
	coins := sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100000)))
	err := bankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.NoError(err)
	err = bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, types.ModuleName, coins)
	suite.NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	v1.RegisterQueryServer(queryHelper, keeper.NewQueryServer(govKeeper))
	legacyQueryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	v1beta1.RegisterQueryServer(legacyQueryHelper, keeper.NewLegacyQueryServer(govKeeper))
	queryClient := v1.NewQueryClient(queryHelper)
	legacyQueryClient := v1beta1.NewQueryClient(legacyQueryHelper)

	suite.ctx = ctx
	suite.govKeeper = govKeeper
	suite.acctKeeper = acctKeeper
	suite.bankKeeper = bankKeeper
	suite.stakingKeeper = stakingKeeper
	suite.distKeeper = distKeeper
	suite.cdc = encCfg.Codec
	suite.queryClient = queryClient
	suite.legacyQueryClient = legacyQueryClient
	suite.msgSrvr = keeper.NewMsgServerImpl(suite.govKeeper)

	suite.legacyMsgSrvr = keeper.NewLegacyMsgServerImpl(govAcct.String(), suite.msgSrvr)
	suite.addrs = simtestutil.AddTestAddrsIncremental(bankKeeper, stakingKeeper, ctx, 3, sdkmath.NewInt(30000000))

	suite.acctKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("cosmos")).AnyTimes()
}

func (suite *KeeperTestSuite) SetupTest() {
	// Initialize the application
	suite.app = availdapp.Setup(suite.T(), false)

	key := suite.app.GetKey(availblob1.StoreKey)
	testCtx := testutil.DefaultContextWithDB(suite.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	suite.sdkCtx = ctx

	suite.msgServer = AvailBlobKeeper.NewMsgServerImpl(suite.keeper)

	// Ensure that the keeper is properly initialized
	require.NotNil(suite.T(), suite.app.AvailBlobKeeper, "AvailBlobKeeper should not be nil")

}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
