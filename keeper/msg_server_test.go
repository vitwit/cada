package keeper_test

import (
	availkeeper "github.com/vitwit/avail-da-module/keeper"
	"github.com/vitwit/avail-da-module/types"
)

func (s *TestSuite) TestMsgServer_SubmitBlob() {

	testCases := []struct {
		name      string
		inputMsg  *types.MsgSubmitBlobRequest
		expectErr bool
	}{
		{
			"submit blob request",
			&types.MsgSubmitBlobRequest{
				ValidatorAddress: s.addrs[1].String(),
				BlocksRange: &types.Range{
					From: 1,
					To:   10,
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgserver.SubmitBlob(s.ctx, tc.inputMsg)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *TestSuite) TestMsgServer_UpdateBlobStatus() {
	err := s.keeper.SetProvenHeight(s.ctx, 10)
	s.Require().NoError(err)

	err = availkeeper.UpdateEndHeight(s.ctx, s.store, uint64(20))
	s.Require().NoError(err)

	availkeeper.UpdateBlobStatus(s.ctx, s.store, uint32(1))
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		inputMsg  *types.MsgUpdateBlobStatusRequest
		expectErr bool
	}{
		{
			"submit blob request",
			&types.MsgUpdateBlobStatusRequest{
				ValidatorAddress: s.addrs[1].String(),
				BlocksRange: &types.Range{
					From: 11,
					To:   20,
				},
				AvailHeight: 20,
				IsSuccess:   true,
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgserver.UpdateBlobStatus(s.ctx, tc.inputMsg)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
