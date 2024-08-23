package keeper_test

import (
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/test-go/testify/require"
	"github.com/vitwit/avail-da-module/keeper"
	"github.com/vitwit/avail-da-module/types"
)

func (s *KeeperTestSuite) TestPreBlocker() {

	require.NotNil(s.T(), s.app.AvailBlobKeeper)
	require.NotNil(s.T(), s.sdkCtx)

	req := abci.RequestFinalizeBlock{
		Time:            time.Now(),
		Height:          1,
		ProposerAddress: []byte{0x1, 0x2, 0x3, 0x4},
		Txs:             [][]byte{},
	}

	tests := []struct {
		name    string
		setup   func()
		req     abci.RequestFinalizeBlock
		wantErr bool
	}{
		{
			name:    "successful execution",
			setup:   func() {},
			req:     req,
			wantErr: false,
		},
		// {
		// 	name: "error during preblockerPendingBlocks",
		// 	setup: func() {
		// 		s.app.AvailBlobKeeper.PreblockerPendingBlocks = func(ctx sdk.Context, time time.Time, proposerAddr []byte, pendingBlocks *[]int64) error {
		// 			return errors.New("preblocker pending blocks error")
		// 		}
		// 	},
		// 	req:     req,
		// 	wantErr: true,
		// },
		// {
		// 	name: "handling larger tx sizes",
		// 	setup: func() {
		// 		req.Txs = append(req.Txs, []byte("larger transaction"))
		// 		req.Size = func() int { return len(req.Txs[0]) } // Simulate the size of txs.
		// 	},
		// 	req:     req,
		// 	wantErr: false,
		// },
	}

	// Run each test case
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := s.app.AvailBlobKeeper.PreBlocker(s.sdkCtx, &tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestPrepareProposal() {
	req := abci.RequestPrepareProposal{
		ProposerAddress: []byte{0x1, 0x2, 0x3, 0x4},
		MaxTxBytes:      1024,
	}

	tests := []struct {
		name    string
		setup   func()
		req     abci.RequestPrepareProposal
		wantErr bool
	}{
		{
			name: "mocked PrepareProposalHandler",
			setup: func() {
				mockPrepareProposalHandler := func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
					return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
				}
				suite.handler = keeper.NewProofOfBlobProposalHandler(
					suite.app.AvailBlobKeeper,
					mockPrepareProposalHandler,
					suite.handler.ProcessProposal,
				)
			},
			req:     req,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			tt.setup()
			resp, err := suite.handler.PrepareProposal(suite.sdkCtx, &tt.req)

			if tt.wantErr {
				suite.Require().Error(err)
				suite.Require().Nil(resp)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
			}
		})
	}
}

// func (suite *KeeperTestSuite) TestProcessProposal() {
// 	// Create a sample request for testing.
// 	req := abci.RequestProcessProposal{
// 		ProposerAddress: []byte{0x1, 0x2, 0x3, 0x4},
// 		Height:          1,
// 		Txs:             [][]byte{[]byte("injected_data"), []byte("tx1"), []byte("tx2")},
// 	}

// 	// Define the test cases
// 	tests := []struct {
// 		name    string
// 		setup   func()
// 		req     abci.RequestProcessProposal
// 		wantErr bool
// 		check   func(resp *abci.ResponseProcessProposal, err error)
// 	}{
// 		{
// 			name: "successful execution without injected data",
// 			setup: func() {
// 				suite.keeper = &MockKeeper{
// 					GetInjectedDataFn: func(txs [][]byte) *types.InjectedData {
// 						return nil
// 					},
// 				}
// 			},
// 			req:     req,
// 			wantErr: false,
// 			check: func(resp *abci.ResponseProcessProposal, err error) {
// 				suite.Require().NoError(err)
// 				suite.Require().NotNil(resp)
// 				suite.Require().Equal(abci.ResponseProcessProposal_ACCEPT, resp.Status)
// 				suite.Require().Equal([][]byte{[]byte("injected_data"), []byte("tx1"), []byte("tx2")}, req.Txs)
// 			},
// 		},
// 		{
// 			name: "successful execution with injected data",
// 			setup: func() {
// 				suite.keeper = &MockKeeper{
// 					GetInjectedDataFn: func(txs [][]byte) *types.InjectedData {
// 						return &types.InjectedData{PendingBlocks: types.PendingBlocks{
// 							BlockHeights: []int64{10, 20},
// 						}}
// 					},
// 					ProcessPendingBlocksFn: func(ctx sdk.Context, time time.Time, pendingBlocks *[]int64) error {
// 						return nil
// 					},
// 				}
// 			},
// 			req:     req,
// 			wantErr: false,
// 			check: func(resp *abci.ResponseProcessProposal, err error) {
// 				suite.Require().NoError(err)
// 				suite.Require().NotNil(resp)
// 				suite.Require().Equal(abci.ResponseProcessProposal_ACCEPT, resp.Status)
// 				suite.Require().Equal([][]byte{[]byte("tx1"), []byte("tx2")}, req.Txs)
// 			},
// 		},
// 		{
// 			name: "error during processPendingBlocks",
// 			setup: func() {
// 				suite.keeper = &MockKeeper{
// 					GetInjectedDataFn: func(txs [][]byte) *types.InjectedData {
// 						return &types.InjectedData{PendingBlocks: types.PendingBlocks{
// 							BlockHeights: []int64{10, 20},
// 						}}
// 					},
// 					ProcessPendingBlocksFn: func(ctx sdk.Context, time time.Time, pendingBlocks *[]int64) error {
// 						return errors.New("mock error")
// 					},
// 				}
// 			},
// 			req:     req,
// 			wantErr: true,
// 			check: func(resp *abci.ResponseProcessProposal, err error) {
// 				suite.Require().Error(err)
// 				suite.Require().Nil(resp)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		suite.T().Run(tt.name, func(t *testing.T) {
// 			tt.setup()

// 			resp, err := suite.handler.ProcessProposal(suite.sdkCtx, &tt.req)

// 			tt.check(resp, err)
// 		})
// 	}
// }

func (suite *KeeperTestSuite) TestProcessProposal() {
	req := abci.RequestProcessProposal{
		Txs: [][]byte{[]byte("mock_data"), []byte("tx1")},
	}

	// Mock behavior for GetInjectedData
	suite.keeper.On("GetInjectedData", req.Txs).Return(&types.InjectedData{
		PendingBlocks: types.PendingBlocks{
			BlockHeights: []int64{10, 20},
		},
	})

	// Mock behavior for ProcessPendingBlocks
	suite.keeper.On("ProcessPendingBlocks", suite.sdkCtx, mock.Anything, mock.Anything).Return(nil)

	// Call the method under test
	resp, err := suite.handler.ProcessProposal(suite.sdkCtx, &req)

	// Validate the response and behavior
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Equal(abci.ResponseProcessProposal_ACCEPT, resp.Status)
	//suite.Require().Equal([][]byte{[]byte("tx1")}, req.Txs)

	// Assert that the mocked methods were called as expected
	//suite.keeper.AssertExpectations(suite.T())

	suite.True(false)
}

type MockKeeper struct {
	mock.Mock
}

// Implement GetInjectedData method on MockKeeper
func (m *MockKeeper) GetInjectedData(txs [][]byte) *types.InjectedData {
	args := m.Called(txs)
	if data := args.Get(0); data != nil {
		return data.(*types.InjectedData)
	}
	return nil
}

// Implement ProcessPendingBlocks method on MockKeeper
func (m *MockKeeper) ProcessPendingBlocks(ctx sdk.Context, time time.Time, pendingBlocks *[]int64) error {
	args := m.Called(ctx, time, pendingBlocks)
	return args.Error(0)
}
