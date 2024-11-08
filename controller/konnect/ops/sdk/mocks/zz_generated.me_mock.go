// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	operations "github.com/Kong/sdk-konnect-go/models/operations"
	mock "github.com/stretchr/testify/mock"
)

// MockMeSDK is an autogenerated mock type for the MeSDK type
type MockMeSDK struct {
	mock.Mock
}

type MockMeSDK_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMeSDK) EXPECT() *MockMeSDK_Expecter {
	return &MockMeSDK_Expecter{mock: &_m.Mock}
}

// GetOrganizationsMe provides a mock function with given fields: ctx, opts
func (_m *MockMeSDK) GetOrganizationsMe(ctx context.Context, opts ...operations.Option) (*operations.GetOrganizationsMeResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationsMe")
	}

	var r0 *operations.GetOrganizationsMeResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...operations.Option) (*operations.GetOrganizationsMeResponse, error)); ok {
		return rf(ctx, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...operations.Option) *operations.GetOrganizationsMeResponse); ok {
		r0 = rf(ctx, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*operations.GetOrganizationsMeResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...operations.Option) error); ok {
		r1 = rf(ctx, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMeSDK_GetOrganizationsMe_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationsMe'
type MockMeSDK_GetOrganizationsMe_Call struct {
	*mock.Call
}

// GetOrganizationsMe is a helper method to define mock.On call
//   - ctx context.Context
//   - opts ...operations.Option
func (_e *MockMeSDK_Expecter) GetOrganizationsMe(ctx interface{}, opts ...interface{}) *MockMeSDK_GetOrganizationsMe_Call {
	return &MockMeSDK_GetOrganizationsMe_Call{Call: _e.mock.On("GetOrganizationsMe",
		append([]interface{}{ctx}, opts...)...)}
}

func (_c *MockMeSDK_GetOrganizationsMe_Call) Run(run func(ctx context.Context, opts ...operations.Option)) *MockMeSDK_GetOrganizationsMe_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]operations.Option, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(operations.Option)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *MockMeSDK_GetOrganizationsMe_Call) Return(_a0 *operations.GetOrganizationsMeResponse, _a1 error) *MockMeSDK_GetOrganizationsMe_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMeSDK_GetOrganizationsMe_Call) RunAndReturn(run func(context.Context, ...operations.Option) (*operations.GetOrganizationsMeResponse, error)) *MockMeSDK_GetOrganizationsMe_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockMeSDK creates a new instance of MockMeSDK. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMeSDK(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMeSDK {
	mock := &MockMeSDK{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
