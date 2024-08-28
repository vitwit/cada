package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/vitwit/avail-da-module/types"
)

func (s *KeeperTestSuite) TestSubmitBlob() {
	addrs := simtestutil.AddTestAddrsIncremental(s.app.BankKeeper, s.app.StakingKeeper, s.sdkCtx, 2, math.NewInt(30000000))
	testCases := []struct {
		name      string
		inputMsg  *types.MsgSubmitBlobRequest
		expectErr bool
	}{
		{
			"send blob request",
			&types.MsgSubmitBlobRequest{
				ValidatorAddress: addrs[1].String(),
				BlocksRange: &types.Range{
					From: uint64(1),
					To:   uint64(5),
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := s.msgServer.SubmitBlob(s.sdkCtx, tc.inputMsg)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

			}
		})
	}
}
