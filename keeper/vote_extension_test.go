package keeper_test

// import (
// 	"encoding/json"
// 	"errors"
// 	"testing"

// 	"cosmossdk.io/log"
// 	"github.com/cometbft/cometbft/abci/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/vitwit/avail-da-module/keeper"
// )

// // MockKeeper embeds keeper.Keeper and overrides methods for mocking
// type MockKeeper struct {
// 	*keeper.Keeper // Embedding keeper.Keeper
// 	mock.Mock      // Embedding testify's Mock for stubbing methods
// }

// func (m *MockKeeper) GetStartHeightFromStore(ctx sdk.Context) uint64 {
// 	args := m.Called(ctx)
// 	return args.Get(0).(uint64)
// }

// func (m *MockKeeper) GetEndHeightFromStore(ctx sdk.Context) uint64 {
// 	args := m.Called(ctx)
// 	return args.Get(0).(uint64)
// }

// func (m *MockKeeper) GetAvailHeightFromStore(ctx sdk.Context) uint64 {
// 	args := m.Called(ctx)
// 	return args.Get(0).(uint64)
// }

// func (m *MockKeeper) GetBlobStatus(ctx sdk.Context) string {
// 	args := m.Called(ctx)
// 	return args.String(0)
// }

// func (m *MockKeeper) GetVotingEndHeightFromStore(ctx sdk.Context) uint64 {
// 	args := m.Called(ctx)
// 	return args.Get(0).(uint64)
// }

// func (m *MockKeeper) IsDataAvailable(ctx sdk.Context, from, to, height uint64, url string) (bool, error) {
// 	args := m.Called(ctx, from, to, height, url)
// 	return args.Bool(0), args.Error(1)
// }

// // MockLogger is a mock implementation of the Logger interface
// type MockLogger struct {
// 	mock.Mock
// }

// // Debug implements log.Logger.
// func (m *MockLogger) Debug(msg string, keyVals ...any) {
// 	panic("unimplemented")
// }

// // Error implements log.Logger.
// func (m *MockLogger) Error(msg string, keyVals ...any) {
// 	panic("unimplemented")
// }

// // Impl implements log.Logger.
// func (m *MockLogger) Impl() any {
// 	panic("unimplemented")
// }

// // Warn implements log.Logger.
// func (m *MockLogger) Warn(msg string, keyVals ...any) {
// 	panic("unimplemented")
// }

// // With implements log.Logger.
// func (m *MockLogger) With(keyVals ...any) log.Logger {
// 	panic("unimplemented")
// }

// func (m *MockLogger) Info(msg string, keyvals ...interface{}) {
// 	m.Called(msg, keyvals)
// }

// func TestVoteExtHandler(t *testing.T) {
// 	// Set up
// 	mockKeeper := &MockKeeper{}
// 	mockLogger := new(MockLogger)
// 	handler := keeper.NewVoteExtHandler(mockLogger, mockKeeper)

// 	ctx := sdk.Context{} // Initialize a suitable context for testing

// 	// Define test cases
// 	tests := []struct {
// 		name         string
// 		mockSetup    func()
// 		req          *types.RequestExtendVote
// 		expectedResp *types.ResponseExtendVote
// 		expectedErr  error
// 	}{
// 		{
// 			name: "Successful ExtendVoteHandler",
// 			mockSetup: func() {
// 				mockKeeper.On("GetStartHeightFromStore", ctx).Return(uint64(100))
// 				mockKeeper.On("GetEndHeightFromStore", ctx).Return(uint64(200))
// 				mockKeeper.On("GetAvailHeightFromStore", ctx).Return(uint64(150))
// 				mockKeeper.On("GetBlobStatus", ctx).Return("IN_VOTING_STATE")
// 				mockKeeper.On("GetVotingEndHeightFromStore", ctx).Return(uint64(100))
// 				mockKeeper.On("IsDataAvailable", ctx, uint64(100), uint64(200), uint64(150), "http://localhost:8000").Return(true, nil)
// 				mockLogger.On("Info", "submitted data to Avail verified successfully at", []interface{}{"block_height", uint64(150)})
// 			},
// 			req: &types.RequestExtendVote{},
// 			expectedResp: &types.ResponseExtendVote{
// 				VoteExtension: marshalVoteExtension(t, &keeper.VoteExtension{
// 					Votes: map[string]bool{"100 200": true},
// 				}),
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			name: "ExtendVoteHandler Error Marshal",
// 			mockSetup: func() {
// 				mockKeeper.On("GetStartHeightFromStore", ctx).Return(uint64(100))
// 				mockKeeper.On("GetEndHeightFromStore", ctx).Return(uint64(200))
// 				mockKeeper.On("GetAvailHeightFromStore", ctx).Return(uint64(150))
// 				mockKeeper.On("GetBlobStatus", ctx).Return("IN_VOTING_STATE")
// 				mockKeeper.On("GetVotingEndHeightFromStore", ctx).Return(uint64(100))
// 				mockKeeper.On("IsDataAvailable", ctx, uint64(100), uint64(200), uint64(150), "http://localhost:8000").Return(true, nil)
// 				mockLogger.On("Info", "submitted data to Avail verified successfully at", []interface{}{"block_height", uint64(150)})
// 			},
// 			req:          &types.RequestExtendVote{},
// 			expectedResp: nil,
// 			expectedErr:  errors.New("failed to marshal vote extension: json: unsupported type: map[interface {}]interface {}"),
// 		},
// 		{
// 			name:      "Successful VerifyVoteExtensionHandler",
// 			mockSetup: func() {},
// 			req:       &types.RequestVerifyVoteExtension{},
// 			expectedResp: &types.ResponseVerifyVoteExtension{
// 				Status: types.ResponseVerifyVoteExtension_ACCEPT,
// 			},
// 			expectedErr: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockSetup()
// 			resp, err := handler.ExtendVoteHandler()(ctx, tt.req)
// 			if tt.expectedErr != nil {
// 				assert.EqualError(t, err, tt.expectedErr.Error())
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.expectedResp, resp)
// 			}

// 			verResp, verErr := handler.VerifyVoteExtensionHandler()(ctx, tt.req)
// 			assert.NoError(t, verErr)
// 			assert.Equal(t, tt.expectedResp, verResp)
// 		})
// 	}
// }

// // Helper function to marshal VoteExtension for comparison
// func marshalVoteExtension(t *testing.T, voteExt *keeper.VoteExtension) []byte {
// 	data, err := json.Marshal(voteExt)
// 	if err != nil {
// 		t.Fatalf("failed to marshal vote extension: %v", err)
// 	}
// 	return data
// }
