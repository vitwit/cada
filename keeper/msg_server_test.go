package keeper_test

import (
	"github.com/vitwit/avail-da-module/types"
)

func (s *TestSuite) TestMsgUpdateBlobStatusRequest_test() {
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
		s.Run(tc.name, func() {
			_, err := s.msgserver.UpdateBlobStatus(s.ctx, tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
func (s *TestSuite) TestMsgSetAvailAddress() {
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
		s.Run(tc.name, func() {
			_, err := s.msgserver.SetAvailAddress(s.ctx, tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
func (s *TestSuite) TestMsgSubmitBlobRequest() {
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
		s.Run(tc.name, func() {
			_, err := s.msgserver.SubmitBlob(s.ctx, tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

// package keeper_test
// import (
//     "github.com/vitwit/avail-da-module/types"

// )
// func (s *TestSuite) TestMsgUpdateBlobStatusRequest_test() {
//     testCases := []struct {
//         name        string
//         blocksRange *types.Range
//         availHeight uint64
//         req         *types.MsgUpdateBlobStatusRequest
//         expectErr   bool
//     }{
//         {
//             name:        "Valid request",
//             blocksRange: &types.Range{From: 100, To: 200},
//             availHeight: 12345,
//             req: &types.MsgUpdateBlobStatusRequest{
//                 ValidatorAddress: "cosmos1h4hj28u89j8dj",
//                 BlocksRange:      &types.Range{From: 100, To: 200},
//                 AvailHeight:      12345,
//                 IsSuccess:        true,
//             },
//             expectErr: false,
//         },
//         {
//             name:        "Invalid request with range where From > To",
//             blocksRange: &types.Range{From: 300, To: 200},
//             availHeight: 10000,
//             req: &types.MsgUpdateBlobStatusRequest{
//                 ValidatorAddress: "cosmos1xyz123",
//                 BlocksRange:      &types.Range{From: 300, To: 200},
//                 AvailHeight:      10000,
//                 IsSuccess:        true,
//             },
//             expectErr: true,
//         },
//         {
//             name:        "Valid request with zero height",
//             blocksRange: &types.Range{From: 0, To: 50},
//             availHeight: 0,
//             req: &types.MsgUpdateBlobStatusRequest{
//                 ValidatorAddress: "cosmos1abcde456",
//                 BlocksRange:      &types.Range{From: 0, To: 50},
//                 AvailHeight:      0,
//                 IsSuccess:        false,
//             },
//             expectErr: false,
//         },
//     }
//     for _, tc := range testCases {
//         s.Run(tc.name, func() {
//             _, err := s.msgserver.UpdateBlobStatus(s.ctx, tc.req)
//             if tc.expectErr {
//                 s.Require().Error(err)
//             } else {
//                 s.Require().NoError(err)
//             }
//         })
//     }
// }
// func (s *TestSuite) TestMsgSetAvailAddress() {
//     testCases := []struct {
//         name          string
//         validatorAddr string
//         availAddr     string
//         req           *types.MsgSetAvailAddress
//         expectErr     bool
//     }{
//         {
//             name:          "Valid request",
//             validatorAddr: "cosmos1h4hj28u89j8dj",
//             availAddr:     "avail1avccvvxg5gt4mn6fw9dvzfrg2q7v",
//             req: &types.MsgSetAvailAddress{
//                 ValidatorAddress: "cosmos1h4hj28u89j8dj",
//                 AvailAddress:     "avail1avccvvxg5gt4mn6fw9dvzfrg2q7v",
//             },
//             expectErr: false,
//         },
//         {
//             name:          "Invalid request with empty ValidatorAddress",
//             validatorAddr: "",
//             availAddr:     "avail1avccvvxg5gt4mn6fw9dvzfrg2q7v",
//             req: &types.MsgSetAvailAddress{
//                 ValidatorAddress: "",
//                 AvailAddress:     "avail1avccvvxg5gt4mn6fw9dvzfrg2q7v",
//             },
//             expectErr: true,
//         },
//         {
//             name:          "Invalid request with empty AvailAddress",
//             validatorAddr: "cosmos1xyz123",
//             availAddr:     "",
//             req: &types.MsgSetAvailAddress{
//                 ValidatorAddress: "cosmos1xyz123",
//                 AvailAddress:     "",
//             },
//             expectErr: true,
//         },
//         {
//             name:          "Invalid request with both addresses empty",
//             validatorAddr: "",
//             availAddr:     "",
//             req: &types.MsgSetAvailAddress{
//                 ValidatorAddress: "",
//                 AvailAddress:     "",
//             },
//             expectErr: true,
//         },
//     }
//     for _, tc := range testCases {
//         s.Run(tc.name, func() {
//             _, err := s.msgserver.SetAvailAddress(s.ctx, tc.req)
//             if tc.expectErr {
//                 s.Require().Error(err)
//             } else {
//                 s.Require().NoError(err)
//             }
//         })
//     }
// }
// func (s *TestSuite) TestMsgSubmitBlobRequest() {
//     testCases := []struct {
//         name          string
//         validatorAddr string
//         blocksRange   *types.Range
//         req           *types.MsgSubmitBlobRequest
//         expectErr     bool
//     }{
//         {
//             name:          "Valid request",
//             validatorAddr: "cosmos1h4hj28u89j8dj",
//             blocksRange:   &types.Range{From: 100, To: 200},
//             req: &types.MsgSubmitBlobRequest{
//                 ValidatorAddress: "cosmos1h4hj28u89j8dj",
//                 BlocksRange:      &types.Range{From: 100, To: 200},
//             },
//             expectErr: false,
//         },
//         {
//             name:          "Invalid request with empty ValidatorAddress",
//             validatorAddr: "",
//             blocksRange:   &types.Range{From: 100, To: 200},
//             req: &types.MsgSubmitBlobRequest{
//                 ValidatorAddress: "",
//                 BlocksRange:      &types.Range{From: 100, To: 200},
//             },
//             expectErr: true,
//         },
//         {
//             name:          "Invalid request with range where From > To",
//             validatorAddr: "cosmos1xyz123",
//             blocksRange:   &types.Range{From: 300, To: 200},
//             req: &types.MsgSubmitBlobRequest{
//                 ValidatorAddress: "cosmos1xyz123",
//                 BlocksRange:      &types.Range{From: 300, To: 200},
//             },
//             expectErr: true,
//         },
//         {
//             name:          "Valid request with empty range",
//             validatorAddr: "cosmos1abcde456",
//             blocksRange:   nil,
//             req: &types.MsgSubmitBlobRequest{
//                 ValidatorAddress: "cosmos1abcde456",
//                 BlocksRange:      nil,
//             },
//             expectErr: false,
//         },
//         {
//             name:          "Invalid request with both ValidatorAddress and Range empty",
//             validatorAddr: "",
//             blocksRange:   nil,
//             req: &types.MsgSubmitBlobRequest{
//                 ValidatorAddress: "",
//                 BlocksRange:      nil,
//             },
//             expectErr: true,
//         },
//     }
//     for _, tc := range testCases {
//         s.Run(tc.name, func() {
//             _, err := s.msgserver.SubmitBlob(s.ctx, tc.req)
//             if tc.expectErr {
//                 s.Require().Error(err)
//             } else {
//                 s.Require().NoError(err)
//             }
//         })
//     }
// }

// package keeper_test

// import (
// 	"errors"
// 	"fmt"
// 	"testing"

// 	"cosmossdk.io/store"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/require"

// 	"github.com/vitwit/avail-da-module/keeper" // Adjust this import path to where your Keeper is defined
// 	"github.com/vitwit/avail-da-module/types"  // Adjust this import path to where your MsgUpdateBlobStatusRequest and MsgUpdateBlobStatusResponse are defined
// )

// // Mock constants for status
// const (
// 	PENDING_STATE   = "PENDING"
// 	IN_VOTING_STATE = "VOTING"
// 	FAILURE_STATE   = "FAILURE"
// )

// // MockContext creates a mock sdk.Context for testing.
// func MockContext() sdk.Context {
// 	memDB := store.NewCommitMultiStore()
// 	memDB.MountStoreWithDB(types.NewKVStoreKey("testStore"), sdk.StoreTypeIAVL, nil)
// 	_ = memDB.LoadLatestVersion()

// 	return sdk.NewContext(memDB, sdk.Header{Height: 1}, false, nil)
// }

// // Mock store key for testing
// var storeKey = sdk.NewKVStoreKey("testStore") // Replace "testStore" with the actual store key name

// // Helper functions for setting up and getting mocked data
// func mockGetProvenHeightFromStore(ctx sdk.Context) uint64 {
// 	return 100
// }

// func mockGetEndHeightFromStore(ctx sdk.Context) uint64 {
// 	return 200
// }

// func mockGetStatusFromStore(store sdk.KVStore) string {
// 	return PENDING_STATE
// }

// // Mock functions for testing Keeper
// func mockUpdateAvailHeight(ctx sdk.Context, store sdk.KVStore, availHeight uint64) {
// 	// Mock implementation if needed
// }

// func mockUpdateVotingEndHeight(ctx sdk.Context, store sdk.KVStore, endHeight uint64) {
// 	// Mock implementation if needed
// }

// func mockUpdateBlobStatus(ctx sdk.Context, store sdk.KVStore, status string) {
// 	// Mock implementation if needed
// }

// // TestUpdateBlobStatus contains the test cases for the UpdateBlobStatus function.
// func TestUpdateBlobStatus(t *testing.T) {
// 	ctx := MockContext()

// 	// Initialize a mock Keeper with the necessary attributes
// 	k := &keeper.Keeper{
// 		StoreKey:       storeKey,
// 		VotingInterval: 5, // Example interval; set according to your requirements
// 	}

// 	// Define test scenarios
// 	tests := []struct {
// 		name         string
// 		setupMocks   func()
// 		req          *types.MsgUpdateBlobStatusRequest
// 		expectedResp *types.MsgUpdateBlobStatusResponse
// 		expectedErr  error
// 	}{
// 		{
// 			name: "Invalid block range request",
// 			setupMocks: func() {
// 				keeper.GetProvenHeightFromStore = mockGetProvenHeightFromStore
// 				keeper.GetEndHeightFromStore = mockGetEndHeightFromStore
// 				keeper.GetStatusFromStore = mockGetStatusFromStore
// 			},
// 			req: &types.MsgUpdateBlobStatusRequest{
// 				BlocksRange: types.BlocksRange{From: 101, To: 199}, // Invalid range
// 				IsSuccess:   true,
// 				AvailHeight: 300,
// 			},
// 			expectedResp: nil,
// 			expectedErr:  fmt.Errorf("invalid blocks range request: expected range [101 -> 200], got [101 -> 199]"),
// 		},
// 		{
// 			name: "Non-pending state",
// 			setupMocks: func() {
// 				keeper.GetProvenHeightFromStore = mockGetProvenHeightFromStore
// 				keeper.GetEndHeightFromStore = mockGetEndHeightFromStore
// 				keeper.GetStatusFromStore = func(store sdk.KVStore) string { return IN_VOTING_STATE } // Status is not pending
// 			},
// 			req: &types.MsgUpdateBlobStatusRequest{
// 				BlocksRange: types.BlocksRange{From: 101, To: 200},
// 				IsSuccess:   true,
// 				AvailHeight: 300,
// 			},
// 			expectedResp: nil,
// 			expectedErr:  errors.New("can update the status if it is not pending"),
// 		},
// 		{
// 			name: "Successful status update",
// 			setupMocks: func() {
// 				keeper.GetProvenHeightFromStore = mockGetProvenHeightFromStore
// 				keeper.GetEndHeightFromStore = mockGetEndHeightFromStore
// 				keeper.GetStatusFromStore = mockGetStatusFromStore // Status is pending
// 			},
// 			req: &types.MsgUpdateBlobStatusRequest{
// 				BlocksRange: types.BlocksRange{From: 101, To: 200},
// 				IsSuccess:   true,
// 				AvailHeight: 300,
// 			},
// 			expectedResp: &types.MsgUpdateBlobStatusResponse{},
// 			expectedErr:  nil,
// 		},
// 		{
// 			name: "Failure status update",
// 			setupMocks: func() {
// 				keeper.GetProvenHeightFromStore = mockGetProvenHeightFromStore
// 				keeper.GetEndHeightFromStore = mockGetEndHeightFromStore
// 				keeper.GetStatusFromStore = mockGetStatusFromStore // Status is pending
// 			},
// 			req: &types.MsgUpdateBlobStatusRequest{
// 				BlocksRange: types.BlocksRange{From: 101, To: 200},
// 				IsSuccess:   false,
// 				AvailHeight: 300,
// 			},
// 			expectedResp: &types.MsgUpdateBlobStatusResponse{},
// 			expectedErr:  nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setupMocks()

// 			resp, err := k.UpdateBlobStatus(ctx, tt.req)

// 			if tt.expectedErr != nil {
// 				require.Error(t, err)
// 				require.EqualError(t, err, tt.expectedErr.Error())
// 			} else {
// 				require.NoError(t, err)
// 				require.Equal(t, tt.expectedResp, resp)
// 			}
// 		})
// 	}
// }
