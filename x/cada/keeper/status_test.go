package keeper_test

import (
	"github.com/vitwit/avail-da-module/x/cada/types"
)

func (s *TestSuite) TestSetBlobStatusPending() {
	testCases := []struct {
		name         string
		startHeight  uint64
		endHeight    uint64
		expectOutput bool
	}{
		{
			"set blob status as pending",
			10,
			20,
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			res := s.keeper.SetBlobStatusPending(s.ctx, tc.startHeight, tc.endHeight)
			status, err := s.queryClient.SubmittedBlobStatus(s.ctx, &types.QuerySubmittedBlobStatusRequest{})
			s.Require().NoError(err)
			if tc.expectOutput {
				s.Require().Equal(status.Status, "PENDING_STATE")
				s.Require().True(res)
			}
		})
	}
}

func (s *TestSuite) TestSetBlobStatus() {
	testCases := []struct {
		name        string
		startHeight uint64
		endHeight   uint64
		status      uint32
		expectErr   bool
	}{
		{
			"set blob status as pending",
			10,
			20,
			0,
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			sErr := s.keeper.SetBlobStatus(s.ctx, tc.status)
			status, err := s.queryClient.SubmittedBlobStatus(s.ctx, &types.QuerySubmittedBlobStatusRequest{})
			s.Require().NoError(err)
			if tc.expectErr {
				s.Require().NoError(sErr)
				s.Require().Equal(status.Status, "READY_STATE")
			}
		})
	}
}
