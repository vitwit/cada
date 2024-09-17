package keeper_test

import (
	availkeeper "github.com/vitwit/avail-da-module/keeper"
	"github.com/vitwit/avail-da-module/types"
)

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

// func (s *TestSuite) TestUpdateBlobStatus() {

// 	testCases := []struct {
// 		name         string
// 		inputMsg     *types.MsgUpdateBlobStatusRequest
// 		provenHeight uint64
// 		endHeight    uint64
// 		status       uint32
// 		expectErr    bool
// 	}{
// 		{
// 			"update blob status",
// 			&types.MsgUpdateBlobStatusRequest{
// 				ValidatorAddress: s.addrs[1].String(),
// 				BlocksRange: &types.Range{
// 					From: 11,
// 					To:   20,
// 				},
// 				AvailHeight: 20,
// 				IsSuccess:   true,
// 			},
// 			10,
// 			20,
// 			1,
// 			false,
// 		},
// 		{
// 			"submit invalid block range request",
// 			&types.MsgUpdateBlobStatusRequest{
// 				ValidatorAddress: s.addrs[1].String(),
// 				BlocksRange: &types.Range{
// 					From: 11,
// 					To:   20,
// 				},
// 				AvailHeight: 20,
// 				IsSuccess:   true,
// 			},
// 			10,
// 			15,
// 			1,
// 			true,
// 		},
// 		{
// 			"status is not in pending state",
// 			&types.MsgUpdateBlobStatusRequest{
// 				ValidatorAddress: s.addrs[1].String(),
// 				BlocksRange: &types.Range{
// 					From: 11,
// 					To:   20,
// 				},
// 				AvailHeight: 20,
// 				IsSuccess:   true,
// 			},
// 			10,
// 			20,
// 			3,
// 			true,
// 		},
// 		{
// 			"update blob status request is not success",
// 			&types.MsgUpdateBlobStatusRequest{
// 				ValidatorAddress: s.addrs[1].String(),
// 				BlocksRange: &types.Range{
// 					From: 11,
// 					To:   20,
// 				},
// 				AvailHeight: 20,
// 				IsSuccess:   false,
// 			},
// 			10,
// 			20,
// 			1,
// 			false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(tc.name, func() {
// 			err := s.keeper.SetProvenHeight(s.ctx, tc.provenHeight)
// 			s.Require().NoError(err)

// 			err = availkeeper.UpdateEndHeight(s.ctx, s.store, tc.endHeight)
// 			s.Require().NoError(err)

// 			availkeeper.UpdateBlobStatus(s.ctx, s.store, tc.status)
// 			s.Require().NoError(err)

// 			err = store.UpdateBlobStatus(s.ctx, s.store, tc.status)
// 			if tc.expectErr {
// 				s.Require().Error(err)
// 			} else {
// 				s.Require().NoError(err)
// 			}
// 		})
// 	}
// }
