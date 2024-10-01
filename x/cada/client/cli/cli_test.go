package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	network "github.com/vitwit/avail-da-module/network"
	"github.com/vitwit/avail-da-module/x/cada/client/cli"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
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

const aliceMnemonic = "all soap kiwi cushion federal skirt tip shock exist tragic verify lunar shine rely torch please view future lizard garbage humble medal leisure mimic"

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	var err error

	appConfig := network.MinimumAppConfig()

	cfg, err := network.DefaultConfigWithAppConfig(appConfig)
	s.Require().NoError(err)

	s.cfg = cfg

	cfg.NumValidators = 1
	cfg.MinGasPrices = "0.000006stake"
	cfg.PublishBlobInterval = "5"
	cfg.LightClientURL = "http://127.0.0.1:8000"

	// Initialize the network
	s.network, err = network.New(s.T(), s.T().TempDir(), cfg)
	s.Require().NoError(err)

	kb := s.network.Validators[0].ClientCtx.Keyring
	path := sdk.GetConfig().GetFullBIP44Path()
	info, err := kb.NewAccount("alice", aliceMnemonic, "", path, hd.Secp256k1)
	s.Require().NoError(err)

	add, err := info.GetAddress()
	s.Require().NoError(err)
	s.addresses = append(s.addresses, add.String())

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestNewUpdateBlobStatusCmd() {
	val := s.network.Validators[0]

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"update blob status - success",
			[]string{
				"1",
				"10",
				"success",
				"120",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addresses[0]),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
			},
			false,
		},
		{
			"update blob status - failure",
			[]string{
				"1",
				"10",
				"failure",
				"120",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addresses[0]),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
			},
			false,
		},
		{
			"update blob status - invalid status",
			[]string{
				"1",
				"10",
				"invalid",
				"120",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addresses[0]),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.NewUpdateBlobStatusCmd()
			res, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)
			if tc.expectErr {
				if err != nil {
					s.Require().Error(err)
				}
			}

			s.Require().NoError(nil)
			s.Require().NotNil(res)
		})
	}
}
