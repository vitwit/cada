// keeper_test.go
package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitwit/avail-da-module/types"
)

func TestKeeper(t *testing.T) {
	suite := new(TestSuite)
	suite.SetupTest()

	t.Run("TestMsgUpdateBlobStatusRequest", suite.TestMsgUpdateBlobStatusRequest)
	t.Run("TestMsgSetAvailAddress", suite.TestMsgSetAvailAddress)
	t.Run("TestMsgSubmitBlobRequest", suite.TestMsgSubmitBlobRequest)
}

func (s *TestSuite) TestMsgUpdateBlobStatusRequest(t *testing.T) {
	testCases := []struct {
		name        string
		blocksRange *types.Range
		availHeight uint64
		req         *types.MsgUpdateBlobStatusRequest
		expectErr   bool
	}{
		{
			name:        "Valid request",
			blocksRange: &types.Range{From: 100, To: 200},
			availHeight: 12345,
			req: &types.MsgUpdateBlobStatusRequest{
				ValidatorAddress: "cosmos1h4hj28u89j8dj",
				BlocksRange:      &types.Range{From: 100, To: 200},
				AvailHeight:      12345,
				IsSuccess:        true,
			},
			expectErr: false,
		},
		{
			name:        "Invalid request with range where From > To",
			blocksRange: &types.Range{From: 300, To: 200},
			availHeight: 10000,
			req: &types.MsgUpdateBlobStatusRequest{
				ValidatorAddress: "cosmos1xyz123",
				BlocksRange:      &types.Range{From: 300, To: 200},
				AvailHeight:      10000,
				IsSuccess:        true,
			},
			expectErr: true,
		},
		{
			name:        "Valid request with zero height",
			blocksRange: &types.Range{From: 0, To: 50},
			availHeight: 0,
			req: &types.MsgUpdateBlobStatusRequest{
				ValidatorAddress: "cosmos1abcde456",
				BlocksRange:      &types.Range{From: 0, To: 50},
				AvailHeight:      0,
				IsSuccess:        false,
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.msgserver.UpdateBlobStatus(s.ctx, tc.req)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func (s *TestSuite) TestMsgSetAvailAddress(t *testing.T) {
	testCases := []struct {
		name          string
		validatorAddr string
		availAddr     string
		req           *types.MsgSetAvailAddress
		expectErr     bool
	}{
		{
			name:          "Valid request",
			validatorAddr: "cosmos1h4hj28u89j8dj",
			availAddr:     "avail1avccvvxg5gt4mn6fw9dvzfrg2q7v",
			req: &types.MsgSetAvailAddress{
				ValidatorAddress: "cosmos1h4hj28u89j8dj",
				AvailAddress:     "avail1avccvvxg5gt4mn6fw9dvzfrg2q7v",
			},
			expectErr: false,
		},
		{
			name:          "Invalid request with empty ValidatorAddress",
			validatorAddr: "",
			availAddr:     "avail1avccvvxg5gt4mn6fw9dvzfrg2q7v",
			req: &types.MsgSetAvailAddress{
				ValidatorAddress: "",
				AvailAddress:     "avail1avccvvxg5gt4mn6fw9dvzfrg2q7v",
			},
			expectErr: true,
		},
		{
			name:          "Invalid request with empty AvailAddress",
			validatorAddr: "cosmos1xyz123",
			availAddr:     "",
			req: &types.MsgSetAvailAddress{
				ValidatorAddress: "cosmos1xyz123",
				AvailAddress:     "",
			},
			expectErr: true,
		},
		{
			name:          "Invalid request with both addresses empty",
			validatorAddr: "",
			availAddr:     "",
			req: &types.MsgSetAvailAddress{
				ValidatorAddress: "",
				AvailAddress:     "",
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.msgserver.SetAvailAddress(s.ctx, tc.req)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func (s *TestSuite) TestMsgSubmitBlobRequest(t *testing.T) {
	testCases := []struct {
		name          string
		validatorAddr string
		blocksRange   *types.Range
		req           *types.MsgSubmitBlobRequest
		expectErr     bool
	}{
		{
			name:          "Valid request",
			validatorAddr: "cosmos1h4hj28u89j8dj",
			blocksRange:   &types.Range{From: 100, To: 200},
			req: &types.MsgSubmitBlobRequest{
				ValidatorAddress: "cosmos1h4hj28u89j8dj",
				BlocksRange:      &types.Range{From: 100, To: 200},
			},
			expectErr: false,
		},
		{
			name:          "Invalid request with empty ValidatorAddress",
			validatorAddr: "",
			blocksRange:   &types.Range{From: 100, To: 200},
			req: &types.MsgSubmitBlobRequest{
				ValidatorAddress: "",
				BlocksRange:      &types.Range{From: 100, To: 200},
			},
			expectErr: true,
		},
		{
			name:          "Invalid request with range where From > To",
			validatorAddr: "cosmos1xyz123",
			blocksRange:   &types.Range{From: 300, To: 200},
			req: &types.MsgSubmitBlobRequest{
				ValidatorAddress: "cosmos1xyz123",
				BlocksRange:      &types.Range{From: 300, To: 200},
			},
			expectErr: true,
		},
		{
			name:          "Valid request with empty range",
			validatorAddr: "cosmos1abcde456",
			blocksRange:   nil,
			req: &types.MsgSubmitBlobRequest{
				ValidatorAddress: "cosmos1abcde456",
				BlocksRange:      nil,
			},
			expectErr: false,
		},
		{
			name:          "Invalid request with both ValidatorAddress and Range empty",
			validatorAddr: "",
			blocksRange:   nil,
			req: &types.MsgSubmitBlobRequest{
				ValidatorAddress: "",
				BlocksRange:      nil,
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.msgserver.SubmitBlob(s.ctx, tc.req)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
