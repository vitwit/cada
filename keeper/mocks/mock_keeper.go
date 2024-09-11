package mocks

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/vitwit/avail-da-module/types"
)

type MockKeeper struct {
	mock.Mock
}

func (m *MockKeeper) SubmitBlob(ctx sdk.Context, req *types.MsgSubmitBlobRequest) {
	m.Called(ctx, req)
}

type MockKVStore struct {
	mock.Mock
}

func (m *MockKVStore) Get(key []byte) []byte {
	args := m.Called(key)
	return args.Get(0).([]byte)
}

func (m *MockKVStore) Set(key, value []byte) {
	m.Called(key, value)
}
