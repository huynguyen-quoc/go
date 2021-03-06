// Code generated by mockery v2.2.1. DO NOT EDIT.

package kafkawriter

import (
	context "context"

	kafka "github.com/huynguyen-quoc/go/streams/kafka"
	mock "github.com/stretchr/testify/mock"
)

// MockWriterInitialization is an autogenerated mock type for the WriterInitialization type
type MockWriterInitialization struct {
	mock.Mock
}

// NewReader provides a mock function with given fields: ctx, configurer
func (_m *MockWriterInitialization) NewReader(ctx context.Context, configurer kafka.Configurer) (Client, error) {
	ret := _m.Called(ctx, configurer)

	var r0 Client
	if rf, ok := ret.Get(0).(func(context.Context, kafka.Configurer) Client); ok {
		r0 = rf(ctx, configurer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, kafka.Configurer) error); ok {
		r1 = rf(ctx, configurer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
