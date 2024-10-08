// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"

	bytes "github.com/cometbft/cometbft/libs/bytes"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"

	mock "github.com/stretchr/testify/mock"
)

// RPCClient is an autogenerated mock type for the RPCClient type
type RPCClient struct {
	mock.Mock
}

// ABCIQuery provides a mock function with given fields: ctx, path, data
func (_m *RPCClient) ABCIQuery(ctx context.Context, path string, data bytes.HexBytes) (*coretypes.ResultABCIQuery, error) {
	ret := _m.Called(ctx, path, data)

	if len(ret) == 0 {
		panic("no return value specified for ABCIQuery")
	}

	var r0 *coretypes.ResultABCIQuery
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bytes.HexBytes) (*coretypes.ResultABCIQuery, error)); ok {
		return rf(ctx, path, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, bytes.HexBytes) *coretypes.ResultABCIQuery); ok {
		r0 = rf(ctx, path, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.ResultABCIQuery)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, bytes.HexBytes) error); ok {
		r1 = rf(ctx, path, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Block provides a mock function with given fields: ctx, height
func (_m *RPCClient) Block(ctx context.Context, height *int64) (*coretypes.ResultBlock, error) {
	ret := _m.Called(ctx, height)

	if len(ret) == 0 {
		panic("no return value specified for Block")
	}

	var r0 *coretypes.ResultBlock
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *int64) (*coretypes.ResultBlock, error)); ok {
		return rf(ctx, height)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *int64) *coretypes.ResultBlock); ok {
		r0 = rf(ctx, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.ResultBlock)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *int64) error); ok {
		r1 = rf(ctx, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Status provides a mock function with given fields: ctx
func (_m *RPCClient) Status(ctx context.Context) (*coretypes.ResultStatus, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Status")
	}

	var r0 *coretypes.ResultStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*coretypes.ResultStatus, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *coretypes.ResultStatus); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.ResultStatus)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRPCClient creates a new instance of RPCClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRPCClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *RPCClient {
	mock := &RPCClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
