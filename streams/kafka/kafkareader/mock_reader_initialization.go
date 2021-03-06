// Code generated by mockery v2.2.1. DO NOT EDIT.

package kafkareader

import (
	context "context"

	kafka "github.com/huynguyen-quoc/go/streams/kafka"
	mock "github.com/stretchr/testify/mock"
)

// MockReaderInitialization is an autogenerated mock type for the ReaderInitialization type
type MockReaderInitialization struct {
	mock.Mock
}

// NewReader provides a mock function with given fields: ctx, entity, configurer
func (_m *MockReaderInitialization) NewReader(ctx context.Context, entity kafka.Entity, configurer kafka.Configurer) (Client, error) {
	ret := _m.Called(ctx, entity, configurer)

	var r0 Client
	if rf, ok := ret.Get(0).(func(context.Context, kafka.Entity, kafka.Configurer) Client); ok {
		r0 = rf(ctx, entity, configurer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, kafka.Entity, kafka.Configurer) error); ok {
		r1 = rf(ctx, entity, configurer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
