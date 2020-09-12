// Code generated by mockery v2.2.1. DO NOT EDIT.

package core

import (
	common "github.com/huynguyen-quoc/go/streams/schema/common"
	mock "github.com/stretchr/testify/mock"
)

// MockMessage is an autogenerated mock type for the Message type
type MockMessage struct {
	mock.Mock
}

// Descriptor provides a mock function with given fields:
func (_m *MockMessage) Descriptor() ([]byte, []int) {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 []int
	if rf, ok := ret.Get(1).(func() []int); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]int)
		}
	}

	return r0, r1
}

// GetStreamInfo provides a mock function with given fields:
func (_m *MockMessage) GetStreamInfo() common.StreamInfoEntity {
	ret := _m.Called()

	var r0 common.StreamInfoEntity
	if rf, ok := ret.Get(0).(func() common.StreamInfoEntity); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(common.StreamInfoEntity)
	}

	return r0
}

// ProtoMessage provides a mock function with given fields:
func (_m *MockMessage) ProtoMessage() {
	_m.Called()
}

// Reset provides a mock function with given fields:
func (_m *MockMessage) Reset() {
	_m.Called()
}

// String provides a mock function with given fields:
func (_m *MockMessage) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}