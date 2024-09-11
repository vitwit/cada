package keeper_test

import (
	// sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	// "github.com/vitwit/avail-da-module/types"
)

type MockKeeper struct {
	mock.Mock
}

// func (m *MockKeeper) GetProvenHeightFromStore(ctx sdk.Context) uint64 {
// 	args := m.Called(ctx)
// 	return args.Get(0).(uint64)
// }

// func (m *MockKeeper) GetEndHeightFromStore(ctx sdk.Context) uint64 {
// 	args := m.Called(ctx)
// 	return args.Get(0).(uint64)
// }

// func (m *MockKeeper) GetStatusFromStore(store sdk.KVStore) string {
// 	args := m.Called(store)
// 	return args.String(0)
// }

// func (m *MockKeeper) UpdateBlobStatus(ctx sdk.Context, req *types.MsgUpdateBlobStatusRequest) (*types.MsgUpdateBlobStatusResponse, error) {
// 	args := m.Called(ctx, req)
// 	return args.Get(0).(*types.MsgUpdateBlobStatusResponse), args.Error(1)
// }

// // Other methods of MockKeeper can be implemented as needed
