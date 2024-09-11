package relayer_test

import (
	"testing"

	addresstypes "cosmossdk.io/core/address"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/suite"
	cada "github.com/vitwit/avail-da-module"
	module "github.com/vitwit/avail-da-module/module"
	relayer "github.com/vitwit/avail-da-module/relayer"
)

type RelayerTestSuite struct {
	suite.Suite

	ctx          sdk.Context
	httpHandler  relayer.HTTPClientHandler
	addrs        []sdk.AccAddress
	encCfg       moduletestutil.TestEncodingConfig
	addressCodec addresstypes.Codec
	baseApp      *baseapp.BaseApp
	relayer      *relayer.Relayer
}

func TestRelayerTestSuite(t *testing.T) {
	suite.Run(t, new(RelayerTestSuite))
}

func (s *RelayerTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(cada.ModuleName)
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

	s.httpHandler = *relayer.NewHTTPClientHandler()

	s.relayer = &relayer.Relayer{}

}
