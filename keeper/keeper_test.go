package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/test-go/testify/suite"
	availdapp "github.com/vitwit/avail-da-module/simapp/app"
)

type KeeperTestSuite struct {
	suite.Suite

	app    *availdapp.ChainApp
	sdkCtx sdk.Context
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = availdapp.Setup(s.T(), false)
	s.app.BaseApp.NewContext(false)

}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
