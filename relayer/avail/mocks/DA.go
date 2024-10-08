// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	avail "github.com/vitwit/avail-da-module/relayer/avail"
)

// DA is an autogenerated mock type for the DA type
type DA struct {
	mock.Mock
}

// GetBlock provides a mock function with given fields: availBlockHeight
func (_m *DA) GetBlock(availBlockHeight int) (avail.Block, error) {
	ret := _m.Called(availBlockHeight)

	if len(ret) == 0 {
		panic("no return value specified for GetBlock")
	}

	var r0 avail.Block
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (avail.Block, error)); ok {
		return rf(availBlockHeight)
	}
	if rf, ok := ret.Get(0).(func(int) avail.Block); ok {
		r0 = rf(availBlockHeight)
	} else {
		r0 = ret.Get(0).(avail.Block)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(availBlockHeight)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsDataAvailable provides a mock function with given fields: data, availBlockHeight
func (_m *DA) IsDataAvailable(data []byte, availBlockHeight int) (bool, error) {
	ret := _m.Called(data, availBlockHeight)

	if len(ret) == 0 {
		panic("no return value specified for IsDataAvailable")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, int) (bool, error)); ok {
		return rf(data, availBlockHeight)
	}
	if rf, ok := ret.Get(0).(func([]byte, int) bool); ok {
		r0 = rf(data, availBlockHeight)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func([]byte, int) error); ok {
		r1 = rf(data, availBlockHeight)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Submit provides a mock function with given fields: data
func (_m *DA) Submit(data []byte) (avail.BlockMetaData, error) {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for Submit")
	}

	var r0 avail.BlockMetaData
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (avail.BlockMetaData, error)); ok {
		return rf(data)
	}
	if rf, ok := ret.Get(0).(func([]byte) avail.BlockMetaData); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Get(0).(avail.BlockMetaData)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDA creates a new instance of DA. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDA(t interface {
	mock.TestingT
	Cleanup(func())
}) *DA {
	mock := &DA{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
