package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/stretchr/testify/suite"
	network "github.com/vitwit/avail-da-module/network"

	app "github.com/vitwit/avail-da-module/simapp"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	cli "github.com/vitwit/avail-da-module/client/cli"
)

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

type IntegrationTestSuite struct {
	suite.Suite

	cfg       network.Config
	network   *network.Network
	addresses []string
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	var err error

	cfg := network.DefaultConfig(app.NewTestNetworkFixture())
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network, err = network.New(s.T(), s.T().TempDir(), cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(10)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestNewSubmitBlobCmd() {
	val := s.network.Validators[0]

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"submit blocks",
			[]string{
				"1",
				"10",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.NewSubmitBlobCmd()
			res, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)
			if tc.expectErr {
				if err != nil {
					s.Require().Error(err)
				}
			}

			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}

}
