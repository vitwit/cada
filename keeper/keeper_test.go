package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/test-go/testify/require"
	"github.com/test-go/testify/suite"
	availblob1 "github.com/vitwit/avail-da-module"
	availabci "github.com/vitwit/avail-da-module/keeper"
	availdapp "github.com/vitwit/avail-da-module/simapp/app"
)

type KeeperTestSuite struct {
	suite.Suite

	app     *availdapp.ChainApp
	sdkCtx  sdk.Context
	handler *availabci.ProofOfBlobProposalHandler
	keeper  *MockKeeper
}

func (suite *KeeperTestSuite) SetupTest() {
	// Initialize the application
	suite.app = availdapp.Setup(suite.T(), false)

	key := suite.app.GetKey(availblob1.StoreKey)
	testCtx := testutil.DefaultContextWithDB(suite.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	suite.sdkCtx = ctx
	suite.keeper = &MockKeeper{}

	// Mock handlers
	mockPrepareHandler := func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		return &abci.ResponsePrepareProposal{}, nil
	}
	mockProcessHandler := func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}

	suite.handler = availabci.NewProofOfBlobProposalHandler(
		suite.app.AvailBlobKeeper,
		mockPrepareHandler,
		mockProcessHandler,
	)

	// Ensure that the keeper is properly initialized
	require.NotNil(suite.T(), suite.app.AvailBlobKeeper, "AvailBlobKeeper should not be nil")

}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
