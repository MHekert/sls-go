// Code generated by mockery v2.21.1. DO NOT EDIT.

package mocks

import (
	context "context"
	core "sls-go/src/items/core"

	mock "github.com/stretchr/testify/mock"
)

// GetImportItemsChannel is an autogenerated mock type for the GetImportItemsChannel type
type GetImportItemsChannel struct {
	mock.Mock
}

// GetImportItemsChannel provides a mock function with given fields: ctx, key, importChannel
func (_m *GetImportItemsChannel) GetImportItemsChannel(ctx context.Context, key string, importChannel chan<- core.Item) error {
	ret := _m.Called(ctx, key, importChannel)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, chan<- core.Item) error); ok {
		r0 = rf(ctx, key, importChannel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewGetImportItemsChannel interface {
	mock.TestingT
	Cleanup(func())
}

// NewGetImportItemsChannel creates a new instance of GetImportItemsChannel. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGetImportItemsChannel(t mockConstructorTestingTNewGetImportItemsChannel) *GetImportItemsChannel {
	mock := &GetImportItemsChannel{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
