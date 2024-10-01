package keeper_test

import (
	"testing"

	addresstypes "cosmossdk.io/core/address"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/suite"
	relayer "github.com/vitwit/avail-da-module/relayer"
	"github.com/vitwit/avail-da-module/x/cada/keeper"
	module "github.com/vitwit/avail-da-module/x/cada/module"
	types "github.com/vitwit/avail-da-module/x/cada/types"
)

type TestSuite struct {
	suite.Suite

	ctx                        sdk.Context
	addrs                      []sdk.AccAddress
	queryClient                types.QueryClient
	keeper                     keeper.Keeper
	msgserver                  types.MsgServer
	encCfg                     moduletestutil.TestEncodingConfig
	addressCodec               addresstypes.Codec
	baseApp                    *baseapp.BaseApp
	upgradeKeeper              upgradekeeper.Keeper
	store                      storetypes.KVStore
	queryServer                types.QueryServer
	voteExtensionHandler       keeper.VoteExtHandler
	logger                     log.Logger
	proofofBlobProposerHandler keeper.ProofOfBlobProposalHandler
	appOptions                 servertypes.AppOptions
	relayer                    relayer.Relayer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(types.ModuleName)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	s.encCfg = moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})
	s.addressCodec = address.NewBech32Codec("cosmos")

	s.baseApp = baseapp.NewBaseApp(
		"cada",
		log.NewNopLogger(),
		testCtx.DB,
		s.encCfg.TxConfig.TxDecoder(),
	)

	s.baseApp.SetCMS(testCtx.CMS)
	s.baseApp.SetInterfaceRegistry(s.encCfg.InterfaceRegistry)
	s.addrs = simtestutil.CreateIncrementalAccounts(7)

	s.keeper = *keeper.NewKeeper(s.encCfg.Codec, storeService, &s.upgradeKeeper, key, s.appOptions, s.logger, &s.relayer)

	s.store = s.ctx.KVStore(key)

	s.queryServer = keeper.NewQueryServerImpl(&s.keeper)
	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, s.queryServer)

	queryClient := types.NewQueryClient(queryHelper)
	s.queryClient = queryClient

	s.msgserver = keeper.NewMsgServerImpl(&s.keeper)

	s.voteExtensionHandler = *keeper.NewVoteExtHandler(s.logger, &s.keeper)

	s.relayer.AvailConfig.PublishBlobInterval = 5

	prepareProposalHandler := func(_ sdk.Context, _ *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		return &abci.ResponsePrepareProposal{}, nil
	}

	processProposalHandler := func(_ sdk.Context, _ *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		return &abci.ResponseProcessProposal{}, nil
	}

	s.proofofBlobProposerHandler = *keeper.NewProofOfBlobProposalHandler(&s.keeper,
		prepareProposalHandler, processProposalHandler, s.voteExtensionHandler)
}
