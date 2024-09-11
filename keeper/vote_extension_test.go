package keeper_test

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vitwit/avail-da-module/keeper"
)

// MockKeeper is a mock implementation of the Keeper interface
type MockKeeper struct {
	mock.Mock
}

func (m *MockKeeper) GetStartHeightFromStore(ctx types.Context) uint64 {
	args := m.Called(ctx)
	return args.Get(0).(uint64)
}

func (m *MockKeeper) GetEndHeightFromStore(ctx types.Context) uint64 {
	args := m.Called(ctx)
	return args.Get(0).(uint64)
}

func (m *MockKeeper) GetAvailHeightFromStore(ctx types.Context) uint64 {
	args := m.Called(ctx)
	return args.Get(0).(uint64)
}

func (m *MockKeeper) GetBlobStatus(ctx types.Context) string {
	args := m.Called(ctx)
	return args.String(0)
}

func (m *MockKeeper) GetVotingEndHeightFromStore(ctx types.Context) uint64 {
	args := m.Called(ctx)
	return args.Get(0).(uint64)
}

// TestVoteExtHandler tests the VoteExtHandler methods
func TestVoteExtHandler(t *testing.T) {
	logger := log.NewNopLogger()
	mockKeeper := new(MockKeeper)
	handler := NewVoteExtHandler(logger, mockKeeper)

	ctx := types.Context{} // Create a mock context

	// Test case 1: Voting not in progress
	mockKeeper.On("GetStartHeightFromStore", ctx).Return(uint64(1))
	mockKeeper.On("GetEndHeightFromStore", ctx).Return(uint64(2))
	mockKeeper.On("GetAvailHeightFromStore", ctx).Return(uint64(3))
	mockKeeper.On("GetBlobStatus", ctx).Return("NOT_IN_VOTING_STATE")
	mockKeeper.On("GetVotingEndHeightFromStore", ctx).Return(uint64(4))

	req := &abci.RequestExtendVote{}
	resp, err := handler.ExtendVoteHandler()(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	var voteExt VoteExtension
	err = json.Unmarshal(resp.VoteExtension, &voteExt)
	assert.NoError(t, err)
	assert.Empty(t, voteExt.Votes)

	// Test case 2: Voting in progress, data available
	mockKeeper.On("GetBlobStatus", ctx).Return("IN_VOTING_STATE")
	mockKeeper.On("GetVotingEndHeightFromStore", ctx).Return(uint64(5))
	mockKeeper.On("GetAvailHeightFromStore", ctx).Return(uint64(3))
	mockKeeper.On("relayer.IsDataAvailable", ctx, uint64(1), uint64(2), uint64(3), "http://localhost:8000").Return(true, nil)

	resp, err = handler.ExtendVoteHandler()(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	err = json.Unmarshal(resp.VoteExtension, &voteExt)
	assert.NoError(t, err)
	assert.NotEmpty(t, voteExt.Votes)
	assert.True(t, voteExt.Votes[Key(1, 2)]) // Check if the vote is recorded

	// Test case 3: Error during marshalling
	mockKeeper.On("GetBlobStatus", ctx).Return("NOT_IN_VOTING_STATE")
	mockKeeper.On("GetVotingEndHeightFromStore", ctx).Return(uint64(4))
	handler.Keeper = nil // Set keeper to nil to trigger marshalling error
	resp, err = handler.ExtendVoteHandler()(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func Key(i1, i2 int) {
	panic("unimplemented")
}

// TestVerifyVoteExtensionHandler tests the VerifyVoteExtensionHandler method
func TestVerifyVoteExtensionHandler(t *testing.T) {
	logger := log.NewNopLogger()
	mockKeeper := new(MockKeeper)
	handler := keeper.NewVoteExtHandler(logger, mockKeeper) // Use the exported NewVoteExtHandler function

	ctx := types.Context{}
	req := &abci.RequestVerifyVoteExtension{}

	resp, err := handler.VerifyVoteExtensionHandler()(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, abci.ResponseVerifyVoteExtension_ACCEPT, resp.Status)
}
