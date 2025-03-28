// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	gen "github.com/jparkkennaby/docker-kafka-go/store/gen"

	mock "github.com/stretchr/testify/mock"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// AddEventRecord provides a mock function with given fields: ctx, request
func (_m *Store) AddEventRecord(ctx context.Context, request gen.AddEventRecordParams) (int64, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for AddEventRecord")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, gen.AddEventRecordParams) (int64, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, gen.AddEventRecordParams) int64); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, gen.AddEventRecordParams) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventRecord provides a mock function with given fields: ctx, id
func (_m *Store) GetEventRecord(ctx context.Context, id int64) (int64, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetEventRecord")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (int64, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) int64); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStore creates a new instance of Store. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *Store {
	mock := &Store{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
