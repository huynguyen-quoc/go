// Code generated by mockery v2.2.1. DO NOT EDIT.

package core

import mock "github.com/stretchr/testify/mock"

// MockStreamProducer is an autogenerated mock type for the StreamProducer type
type MockStreamProducer struct {
	mock.Mock
}

// Done provides a mock function with given fields:
func (_m *MockStreamProducer) Done() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// SaveBytesToStream provides a mock function with given fields: partitionBytes, dataBytes
func (_m *MockStreamProducer) SaveBytesToStream(partitionBytes []byte, dataBytes []byte) error {
	ret := _m.Called(partitionBytes, dataBytes)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, []byte) error); ok {
		r0 = rf(partitionBytes, dataBytes)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveToStream provides a mock function with given fields: dto
func (_m *MockStreamProducer) SaveToStream(dto WriterDTO) error {
	ret := _m.Called(dto)

	var r0 error
	if rf, ok := ret.Get(0).(func(WriterDTO) error); ok {
		r0 = rf(dto)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Shutdown provides a mock function with given fields:
func (_m *MockStreamProducer) Shutdown() {
	_m.Called()
}
