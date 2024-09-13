package keeper_test

import (
	"testing"

	cadaApp "simapp/app"

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
	cada "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/keeper"
	availkeeper "github.com/vitwit/avail-da-module/keeper"
	mocks "github.com/vitwit/avail-da-module/keeper/mocks"
	module "github.com/vitwit/avail-da-module/module"
	"github.com/vitwit/avail-da-module/types"
)

type TestSuite struct {
	suite.Suite

	ctx                        sdk.Context
	addrs                      []sdk.AccAddress
	encodedAddrs               []string
	queryClient                types.QueryClient
	keeper                     keeper.Keeper
	msgserver                  types.MsgServer
	encCfg                     moduletestutil.TestEncodingConfig
	addressCodec               addresstypes.Codec
	baseApp                    *baseapp.BaseApp
	mockKeeper                 *mocks.MockKeeper
	appOpts                    servertypes.AppOptions
	app                        *cadaApp.ChainApp
	upgradeKeeper              upgradekeeper.Keeper
	store                      storetypes.KVStore
	queryServer                types.QueryServer
	voteExtensionHandler       keeper.VoteExtHandler
	logger                     log.Logger
	proofofBlobProposerHandler keeper.ProofOfBlobProposalHandler
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(cada.ModuleName)
	s.mockKeeper = new(mocks.MockKeeper)
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

	s.keeper = *keeper.NewKeeper(s.encCfg.Codec, s.appOpts, storeService, &s.upgradeKeeper, key, 10, 1)

	s.store = s.ctx.KVStore(key)

	s.queryServer = keeper.NewQueryServerImpl(&s.keeper)
	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, s.queryServer)

	queryClient := types.NewQueryClient(queryHelper)
	s.queryClient = queryClient

	s.msgserver = keeper.NewMsgServerImpl(&s.keeper)

	s.voteExtensionHandler = *keeper.NewVoteExtHandler(s.logger, &s.keeper)

	prepareProposalHandler := func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		return &abci.ResponsePrepareProposal{}, nil
	}

	processProposalHandler := func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		return &abci.ResponseProcessProposal{}, nil
	}

	s.proofofBlobProposerHandler = *keeper.NewProofOfBlobProposalHandler(&s.keeper,
		prepareProposalHandler, processProposalHandler, s.voteExtensionHandler)

}

func (s *TestSuite) TestUpdateBlobStatus() {

	testCases := []struct {
		name         string
		inputMsg     *types.MsgUpdateBlobStatusRequest
		provenHeight uint64
		endHeight    uint64
		status       uint32
		expectErr    bool
	}{
		{
			"update blob status",
			&types.MsgUpdateBlobStatusRequest{
				ValidatorAddress: s.addrs[1].String(),
				BlocksRange: &types.Range{
					From: 11,
					To:   20,
				},
				AvailHeight: 20,
				IsSuccess:   true,
			},
			10,
			20,
			1,
			false,
		},
		{
			"submit invalid block range request",
			&types.MsgUpdateBlobStatusRequest{
				ValidatorAddress: s.addrs[1].String(),
				BlocksRange: &types.Range{
					From: 11,
					To:   20,
				},
				AvailHeight: 20,
				IsSuccess:   true,
			},
			10,
			15,
			1,
			true,
		},
		{
			"status is not in pending state",
			&types.MsgUpdateBlobStatusRequest{
				ValidatorAddress: s.addrs[1].String(),
				BlocksRange: &types.Range{
					From: 11,
					To:   20,
				},
				AvailHeight: 20,
				IsSuccess:   true,
			},
			10,
			20,
			3,
			true,
		},
		{
			"update blob status request is not success",
			&types.MsgUpdateBlobStatusRequest{
				ValidatorAddress: s.addrs[1].String(),
				BlocksRange: &types.Range{
					From: 11,
					To:   20,
				},
				AvailHeight: 20,
				IsSuccess:   false,
			},
			10,
			20,
			1,
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.keeper.SetProvenHeight(s.ctx, tc.provenHeight)
			s.Require().NoError(err)

			err = availkeeper.UpdateEndHeight(s.ctx, s.store, tc.endHeight)
			s.Require().NoError(err)

			availkeeper.UpdateBlobStatus(s.ctx, s.store, tc.status)
			s.Require().NoError(err)

			res, err := s.keeper.UpdateBlobStatus(s.ctx, tc.inputMsg)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(res)
			}
		})
	}
}
