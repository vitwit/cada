package keeper_test

import (
	"path/filepath"
	"testing"

	addresstypes "cosmossdk.io/core/address"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"
	name "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/keeper"
	mocks "github.com/vitwit/avail-da-module/keeper/mocks"
	cada "github.com/vitwit/avail-da-module/module"
	"github.com/vitwit/avail-da-module/types"
)

type TestSuite struct {
	suite.Suite
	ctx             sdk.Context
	addrs           []sdk.AccAddress
	encodedAddrs    []string
	queryClient     types.QueryClient
	keeper          *keeper.Keeper
	msgserver       types.MsgServer
	encCfg          moduletestutil.TestEncodingConfig
	addressCodec    addresstypes.Codec
	baseApp         *baseapp.BaseApp
	mockKeeper      *mocks.MockKeeper
	appOptions      servertypes.AppOptions
	upgradeKeeper   upgradekeeper.Keeper
	SetupTestKeeper keeper.Keeper
	// cdc           codec.cdc wrong :(
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
func (s *TestSuite) SetupTest() {
	s.keeper, s.ctx = keeper.SetupTestKeeper()
	s.msgserver = keeper.NewMsgServerImpl(s.keeper)
	s.mockKeeper = new(mocks.MockKeeper)
	key := storetypes.NewKVStoreKey(name.ModuleName)
	// upgradekey := corestore.
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	s.encCfg = moduletestutil.MakeTestEncodingConfig(cada.AppModuleBasic{})
	s.addressCodec = address.NewBech32Codec("cosmos")
	s.baseApp = baseapp.NewBaseApp(
		"cada",
		log.NewNopLogger(),
		testCtx.DB,
		s.encCfg.TxConfig.TxDecoder(),
	)
	s.baseApp.SetCMS(testCtx.CMS)
	s.baseApp.SetInterfaceRegistry(s.encCfg.InterfaceRegistry)
	homeDir := filepath.Join(s.T().TempDir(), "x_upgrade_keeper_test")
	skipUpgradeHeights := make(map[int64]bool)
	upgradeKeeper := upgradekeeper.NewKeeper(skipUpgradeHeights, storeService, s.encCfg.Codec, homeDir, nil, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	s.upgradeKeeper = *upgradeKeeper
	s.keeper = keeper.NewKeeper(s.encCfg.Codec, s.appOptions, storeService, upgradeKeeper, key, 5, 1)
	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.encCfg.InterfaceRegistry)
	// types.RegisterQueryServer(queryHelper, s.keeper)
	queryClient := types.NewQueryClient(queryHelper)
	s.queryClient = queryClient
	s.msgserver = keeper.NewMsgServerImpl(s.keeper)
}
