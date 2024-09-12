package keeper_test

// import (
// 	"encoding/json"
// 	"testing"

// 	"cosmossdk.io/log"
// 	abci "github.com/cometbft/cometbft/abci/types"
// 	"github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/vitwit/avail-da-module/keeper"
// )
// package keeper

// import (
// 	"encoding/json"
// 	"fmt"

// 	"cosmossdk.io/log"
// 	abci "github.com/cometbft/cometbft/abci/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// type VoteExtHandler struct {
// 	logger log.Logger

// 	Keeper *Keeper
// }

// // TODO: add required parameters like avail light client url, etc..
// func NewVoteExtHandler(
// 	logger log.Logger,
// 	keeper *Keeper,
// ) *VoteExtHandler {
// 	return &VoteExtHandler{
// 		logger: logger,
// 		Keeper: keeper,
// 	}
// }

// func Key(from, to uint64) string {
// 	return fmt.Sprintln(from, " ", to)
// }

// // TODO: change the Vote Extension to be actually usable
// type VoteExtension struct {
// 	Votes map[string]bool
// }

// func (h *VoteExtHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {

// 	return func(ctx sdk.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {

// 		// TODO: implement proper logic, this is for demo purpose only
// 		from := h.Keeper.GetStartHeightFromStore(ctx)
// 		end := h.Keeper.GetEndHeightFromStore(ctx)

// 		availHeight := h.Keeper.GetAvailHeightFromStore(ctx)

// 		pendingRangeKey := Key(from, end)

// 		blobStatus := h.Keeper.GetBlobStatus(ctx)
// 		currentHeight := ctx.BlockHeight()
// 		voteEndHeight := h.Keeper.GetVotingEndHeightFromStore(ctx)
// 		Votes := make(map[string]bool, 1)

// 		abciResponseVoteExt := &abci.ResponseExtendVote{}

// 		if currentHeight+1 != int64(voteEndHeight) || blobStatus != IN_VOTING_STATE {
// 			voteExt := VoteExtension{
// 				Votes: Votes,
// 			}

// 			//TODO: use better marshalling instead of json (eg: proto marshalling)
// 			votesBytes, err := json.Marshal(voteExt)
// 			if err != nil {
// 				return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
// 			}
// 			abciResponseVoteExt.VoteExtension = votesBytes
// 			return abciResponseVoteExt, nil
// 		}

// 		ok, err := h.Keeper.relayer.IsDataAvailable(ctx, from, end, availHeight, "http://localhost:8000")
// 		fmt.Println("checking light client...", ok, err)
// 		if ok {
// 			h.logger.Info("submitted data to Avail verified successfully at",
// 				"block_height", availHeight,
// 			)
// 		}

// 		// ok, checkLightClient()
// 		Votes[pendingRangeKey] = ok
// 		voteExt := VoteExtension{
// 			Votes: Votes,
// 		}

// 		//TODO: use proto marshalling instead
// 		votesBytes, err := json.Marshal(voteExt)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
// 		}

// 		return &abci.ResponseExtendVote{
// 			VoteExtension: votesBytes,
// 		}, nil
// 	}
// }

// func (h *VoteExtHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
// 	return func(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
// 		// TODO: write proper validation for the votes
// 		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
// 	}
// }

// // MockKeeper is a mock implementation of the Keeper interface
// type MockKeeper struct {
// 	mock.Mock
// }

// func (m *MockKeeper) GetStartHeightFromStore(ctx types.Context) uint64 {
// 	args := m.Called(ctx)
// 	return args.Get(0).(uint64)
// }

// func (m *MockKeeper) GetEndHeightFromStore(ctx types.Context) uint64 {
// 	args := m.Called(ctx)
// 	return args.Get(0).(uint64)
// }

// func (m *MockKeeper) GetAvailHeightFromStore(ctx types.Context) uint64 {
// 	args := m.Called(ctx)
// 	return args.Get(0).(uint64)
// }

// func (m *MockKeeper) GetBlobStatus(ctx types.Context) string {
// 	args := m.Called(ctx)
// 	return args.String(0)
// }

// func (m *MockKeeper) GetVotingEndHeightFromStore(ctx types.Context) uint64 {
// 	args := m.Called(ctx)
// 	return args.Get(0).(uint64)
// }

// // TestVoteExtHandler tests the VoteExtHandler methods
// func TestVoteExtHandler(t *testing.T) {
// 	logger := log.NewNopLogger()
// 	mockKeeper := new(MockKeeper)
// 	handler := NewVoteExtHandler(logger, mockKeeper)

// 	ctx := types.Context{} // Create a mock context

// 	// Test case 1: Voting not in progress
// 	mockKeeper.On("GetStartHeightFromStore", ctx).Return(uint64(1))
// 	mockKeeper.On("GetEndHeightFromStore", ctx).Return(uint64(2))
// 	mockKeeper.On("GetAvailHeightFromStore", ctx).Return(uint64(3))
// 	mockKeeper.On("GetBlobStatus", ctx).Return("NOT_IN_VOTING_STATE")
// 	mockKeeper.On("GetVotingEndHeightFromStore", ctx).Return(uint64(4))

// 	req := &abci.RequestExtendVote{}
// 	resp, err := handler.ExtendVoteHandler()(ctx, req)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)

// 	var voteExt VoteExtension
// 	err = json.Unmarshal(resp.VoteExtension, &voteExt)
// 	assert.NoError(t, err)
// 	assert.Empty(t, voteExt.Votes)

// 	// Test case 2: Voting in progress, data available
// 	mockKeeper.On("GetBlobStatus", ctx).Return("IN_VOTING_STATE")
// 	mockKeeper.On("GetVotingEndHeightFromStore", ctx).Return(uint64(5))
// 	mockKeeper.On("GetAvailHeightFromStore", ctx).Return(uint64(3))
// 	mockKeeper.On("relayer.IsDataAvailable", ctx, uint64(1), uint64(2), uint64(3), "http://localhost:8000").Return(true, nil)

// 	resp, err = handler.ExtendVoteHandler()(ctx, req)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)

// 	err = json.Unmarshal(resp.VoteExtension, &voteExt)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, voteExt.Votes)
// 	assert.True(t, voteExt.Votes[Key(1, 2)]) // Check if the vote is recorded

// 	// Test case 3: Error during marshalling
// 	mockKeeper.On("GetBlobStatus", ctx).Return("NOT_IN_VOTING_STATE")
// 	mockKeeper.On("GetVotingEndHeightFromStore", ctx).Return(uint64(4))
// 	handler.Keeper = nil // Set keeper to nil to trigger marshalling error
// 	resp, err = handler.ExtendVoteHandler()(ctx, req)

// 	assert.Error(t, err)
// 	assert.Nil(t, resp)
// }

// func Key(i1, i2 int) {
// 	panic("unimplemented")
// }

// // TestVerifyVoteExtensionHandler tests the VerifyVoteExtensionHandler method
// func TestVerifyVoteExtensionHandler(t *testing.T) {
// 	logger := log.NewNopLogger()
// 	mockKeeper := new(MockKeeper)
// 	handler := keeper.NewVoteExtHandler(logger, mockKeeper) // Use the exported NewVoteExtHandler function

// 	ctx := types.Context{}
// 	req := &abci.RequestVerifyVoteExtension{}

// 	resp, err := handler.VerifyVoteExtensionHandler()(ctx, req)

// 	assert.NoError(t, err)
// 	assert.Equal(t, abci.ResponseVerifyVoteExtension_ACCEPT, resp.Status)
// }
