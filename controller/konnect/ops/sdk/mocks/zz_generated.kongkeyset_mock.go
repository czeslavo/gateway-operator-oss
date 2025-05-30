// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	components "github.com/Kong/sdk-konnect-go/models/components"

	mock "github.com/stretchr/testify/mock"

	operations "github.com/Kong/sdk-konnect-go/models/operations"
)

// MockKeySetsSDK is an autogenerated mock type for the KeySetsSDK type
type MockKeySetsSDK struct {
	mock.Mock
}

type MockKeySetsSDK_Expecter struct {
	mock *mock.Mock
}

func (_m *MockKeySetsSDK) EXPECT() *MockKeySetsSDK_Expecter {
	return &MockKeySetsSDK_Expecter{mock: &_m.Mock}
}

// CreateKeySet provides a mock function with given fields: ctx, controlPlaneID, keySet, opts
func (_m *MockKeySetsSDK) CreateKeySet(ctx context.Context, controlPlaneID string, keySet components.KeySetInput, opts ...operations.Option) (*operations.CreateKeySetResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, controlPlaneID, keySet)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CreateKeySet")
	}

	var r0 *operations.CreateKeySetResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, components.KeySetInput, ...operations.Option) (*operations.CreateKeySetResponse, error)); ok {
		return rf(ctx, controlPlaneID, keySet, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, components.KeySetInput, ...operations.Option) *operations.CreateKeySetResponse); ok {
		r0 = rf(ctx, controlPlaneID, keySet, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*operations.CreateKeySetResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, components.KeySetInput, ...operations.Option) error); ok {
		r1 = rf(ctx, controlPlaneID, keySet, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockKeySetsSDK_CreateKeySet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateKeySet'
type MockKeySetsSDK_CreateKeySet_Call struct {
	*mock.Call
}

// CreateKeySet is a helper method to define mock.On call
//   - ctx context.Context
//   - controlPlaneID string
//   - keySet components.KeySetInput
//   - opts ...operations.Option
func (_e *MockKeySetsSDK_Expecter) CreateKeySet(ctx interface{}, controlPlaneID interface{}, keySet interface{}, opts ...interface{}) *MockKeySetsSDK_CreateKeySet_Call {
	return &MockKeySetsSDK_CreateKeySet_Call{Call: _e.mock.On("CreateKeySet",
		append([]interface{}{ctx, controlPlaneID, keySet}, opts...)...)}
}

func (_c *MockKeySetsSDK_CreateKeySet_Call) Run(run func(ctx context.Context, controlPlaneID string, keySet components.KeySetInput, opts ...operations.Option)) *MockKeySetsSDK_CreateKeySet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]operations.Option, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(operations.Option)
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(components.KeySetInput), variadicArgs...)
	})
	return _c
}

func (_c *MockKeySetsSDK_CreateKeySet_Call) Return(_a0 *operations.CreateKeySetResponse, _a1 error) *MockKeySetsSDK_CreateKeySet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockKeySetsSDK_CreateKeySet_Call) RunAndReturn(run func(context.Context, string, components.KeySetInput, ...operations.Option) (*operations.CreateKeySetResponse, error)) *MockKeySetsSDK_CreateKeySet_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteKeySet provides a mock function with given fields: ctx, controlPlaneID, keySetID, opts
func (_m *MockKeySetsSDK) DeleteKeySet(ctx context.Context, controlPlaneID string, keySetID string, opts ...operations.Option) (*operations.DeleteKeySetResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, controlPlaneID, keySetID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteKeySet")
	}

	var r0 *operations.DeleteKeySetResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, ...operations.Option) (*operations.DeleteKeySetResponse, error)); ok {
		return rf(ctx, controlPlaneID, keySetID, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, ...operations.Option) *operations.DeleteKeySetResponse); ok {
		r0 = rf(ctx, controlPlaneID, keySetID, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*operations.DeleteKeySetResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, ...operations.Option) error); ok {
		r1 = rf(ctx, controlPlaneID, keySetID, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockKeySetsSDK_DeleteKeySet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteKeySet'
type MockKeySetsSDK_DeleteKeySet_Call struct {
	*mock.Call
}

// DeleteKeySet is a helper method to define mock.On call
//   - ctx context.Context
//   - controlPlaneID string
//   - keySetID string
//   - opts ...operations.Option
func (_e *MockKeySetsSDK_Expecter) DeleteKeySet(ctx interface{}, controlPlaneID interface{}, keySetID interface{}, opts ...interface{}) *MockKeySetsSDK_DeleteKeySet_Call {
	return &MockKeySetsSDK_DeleteKeySet_Call{Call: _e.mock.On("DeleteKeySet",
		append([]interface{}{ctx, controlPlaneID, keySetID}, opts...)...)}
}

func (_c *MockKeySetsSDK_DeleteKeySet_Call) Run(run func(ctx context.Context, controlPlaneID string, keySetID string, opts ...operations.Option)) *MockKeySetsSDK_DeleteKeySet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]operations.Option, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(operations.Option)
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockKeySetsSDK_DeleteKeySet_Call) Return(_a0 *operations.DeleteKeySetResponse, _a1 error) *MockKeySetsSDK_DeleteKeySet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockKeySetsSDK_DeleteKeySet_Call) RunAndReturn(run func(context.Context, string, string, ...operations.Option) (*operations.DeleteKeySetResponse, error)) *MockKeySetsSDK_DeleteKeySet_Call {
	_c.Call.Return(run)
	return _c
}

// ListKeySet provides a mock function with given fields: ctx, request, opts
func (_m *MockKeySetsSDK) ListKeySet(ctx context.Context, request operations.ListKeySetRequest, opts ...operations.Option) (*operations.ListKeySetResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, request)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListKeySet")
	}

	var r0 *operations.ListKeySetResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, operations.ListKeySetRequest, ...operations.Option) (*operations.ListKeySetResponse, error)); ok {
		return rf(ctx, request, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, operations.ListKeySetRequest, ...operations.Option) *operations.ListKeySetResponse); ok {
		r0 = rf(ctx, request, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*operations.ListKeySetResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, operations.ListKeySetRequest, ...operations.Option) error); ok {
		r1 = rf(ctx, request, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockKeySetsSDK_ListKeySet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListKeySet'
type MockKeySetsSDK_ListKeySet_Call struct {
	*mock.Call
}

// ListKeySet is a helper method to define mock.On call
//   - ctx context.Context
//   - request operations.ListKeySetRequest
//   - opts ...operations.Option
func (_e *MockKeySetsSDK_Expecter) ListKeySet(ctx interface{}, request interface{}, opts ...interface{}) *MockKeySetsSDK_ListKeySet_Call {
	return &MockKeySetsSDK_ListKeySet_Call{Call: _e.mock.On("ListKeySet",
		append([]interface{}{ctx, request}, opts...)...)}
}

func (_c *MockKeySetsSDK_ListKeySet_Call) Run(run func(ctx context.Context, request operations.ListKeySetRequest, opts ...operations.Option)) *MockKeySetsSDK_ListKeySet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]operations.Option, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(operations.Option)
			}
		}
		run(args[0].(context.Context), args[1].(operations.ListKeySetRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockKeySetsSDK_ListKeySet_Call) Return(_a0 *operations.ListKeySetResponse, _a1 error) *MockKeySetsSDK_ListKeySet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockKeySetsSDK_ListKeySet_Call) RunAndReturn(run func(context.Context, operations.ListKeySetRequest, ...operations.Option) (*operations.ListKeySetResponse, error)) *MockKeySetsSDK_ListKeySet_Call {
	_c.Call.Return(run)
	return _c
}

// UpsertKeySet provides a mock function with given fields: ctx, request, opts
func (_m *MockKeySetsSDK) UpsertKeySet(ctx context.Context, request operations.UpsertKeySetRequest, opts ...operations.Option) (*operations.UpsertKeySetResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, request)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for UpsertKeySet")
	}

	var r0 *operations.UpsertKeySetResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, operations.UpsertKeySetRequest, ...operations.Option) (*operations.UpsertKeySetResponse, error)); ok {
		return rf(ctx, request, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, operations.UpsertKeySetRequest, ...operations.Option) *operations.UpsertKeySetResponse); ok {
		r0 = rf(ctx, request, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*operations.UpsertKeySetResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, operations.UpsertKeySetRequest, ...operations.Option) error); ok {
		r1 = rf(ctx, request, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockKeySetsSDK_UpsertKeySet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpsertKeySet'
type MockKeySetsSDK_UpsertKeySet_Call struct {
	*mock.Call
}

// UpsertKeySet is a helper method to define mock.On call
//   - ctx context.Context
//   - request operations.UpsertKeySetRequest
//   - opts ...operations.Option
func (_e *MockKeySetsSDK_Expecter) UpsertKeySet(ctx interface{}, request interface{}, opts ...interface{}) *MockKeySetsSDK_UpsertKeySet_Call {
	return &MockKeySetsSDK_UpsertKeySet_Call{Call: _e.mock.On("UpsertKeySet",
		append([]interface{}{ctx, request}, opts...)...)}
}

func (_c *MockKeySetsSDK_UpsertKeySet_Call) Run(run func(ctx context.Context, request operations.UpsertKeySetRequest, opts ...operations.Option)) *MockKeySetsSDK_UpsertKeySet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]operations.Option, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(operations.Option)
			}
		}
		run(args[0].(context.Context), args[1].(operations.UpsertKeySetRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockKeySetsSDK_UpsertKeySet_Call) Return(_a0 *operations.UpsertKeySetResponse, _a1 error) *MockKeySetsSDK_UpsertKeySet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockKeySetsSDK_UpsertKeySet_Call) RunAndReturn(run func(context.Context, operations.UpsertKeySetRequest, ...operations.Option) (*operations.UpsertKeySetResponse, error)) *MockKeySetsSDK_UpsertKeySet_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockKeySetsSDK creates a new instance of MockKeySetsSDK. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockKeySetsSDK(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockKeySetsSDK {
	mock := &MockKeySetsSDK{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
