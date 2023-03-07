// Code generated by mockery v2.21.1. DO NOT EDIT.

package mocks

import (
	core "sls-go/src/items/core"

	mock "github.com/stretchr/testify/mock"
)

// BatchPersister is an autogenerated mock type for the BatchPersister type
type BatchPersister struct {
	mock.Mock
}

// PersistBatch provides a mock function with given fields: _a0
func (_m *BatchPersister) PersistBatch(_a0 []*core.Item) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*core.Item) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewBatchPersister interface {
	mock.TestingT
	Cleanup(func())
}

// NewBatchPersister creates a new instance of BatchPersister. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBatchPersister(t mockConstructorTestingTNewBatchPersister) *BatchPersister {
	mock := &BatchPersister{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
