// Code generated by mockery v2.2.1. DO NOT EDIT.

package core

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// MockConsumerMessage is an autogenerated mock type for the ConsumerMessage type
type MockConsumerMessage struct {
	mock.Mock
}

// Data provides a mock function with given fields:
func (_m *MockConsumerMessage) Data() []byte {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// Offset provides a mock function with given fields:
func (_m *MockConsumerMessage) Offset() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// Partition provides a mock function with given fields:
func (_m *MockConsumerMessage) Partition() int32 {
	ret := _m.Called()

	var r0 int32
	if rf, ok := ret.Get(0).(func() int32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int32)
	}

	return r0
}

// Timestamp provides a mock function with given fields:
func (_m *MockConsumerMessage) Timestamp() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}
