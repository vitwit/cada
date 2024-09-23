package keeper_test

import (
	availkeeper "github.com/vitwit/avail-da-module/keeper"
	"github.com/vitwit/avail-da-module/types"
)

func (s *TestSuite) TestMsgServer_UpdateBlobStatus() {

	err := availkeeper.UpdateEndHeight(s.ctx, s.store, uint64(20))
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		inputMsg  *types.MsgUpdateBlobStatusRequest
		status    uint32
		expectErr bool
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
			1,
			false,
		},
		{
			"status is not in pending",
			&types.MsgUpdateBlobStatusRequest{
				ValidatorAddress: s.addrs[1].String(),
				BlocksRange: &types.Range{
					From: 11,
					To:   20,
				},
				AvailHeight: 20,
				IsSuccess:   true,
			},
			3,
			true,
		},
		{
			"data posting is not succeeded",
			&types.MsgUpdateBlobStatusRequest{
				ValidatorAddress: s.addrs[1].String(),
				BlocksRange: &types.Range{
					From: 11,
					To:   20,
				},
				AvailHeight: 20,
				IsSuccess:   false,
			},
			1,
			false,
		},
		{
			"invalid blocks range",
			&types.MsgUpdateBlobStatusRequest{
				ValidatorAddress: s.addrs[1].String(),
				BlocksRange: &types.Range{
					From: 12,
					To:   20,
				},
				AvailHeight: 20,
				IsSuccess:   true,
			},
			1,
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			availkeeper.UpdateBlobStatus(s.ctx, s.store, tc.status)
			s.Require().NoError(err)

			_, err := s.msgserver.UpdateBlobStatus(s.ctx, tc.inputMsg)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
