// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	gorestapi "github.com/snowzach/gorestapi/gorestapi"
	mock "github.com/stretchr/testify/mock"

	queryp "github.com/snowzach/queryp"
)

// GRStore is an autogenerated mock type for the GRStore type
type GRStore struct {
	mock.Mock
}

// ThingDeleteByID provides a mock function with given fields: ctx, id
func (_m *GRStore) ThingDeleteByID(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ThingGetByID provides a mock function with given fields: ctx, id
func (_m *GRStore) ThingGetByID(ctx context.Context, id string) (*gorestapi.Thing, error) {
	ret := _m.Called(ctx, id)

	var r0 *gorestapi.Thing
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*gorestapi.Thing, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *gorestapi.Thing); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorestapi.Thing)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ThingSave provides a mock function with given fields: ctx, thing
func (_m *GRStore) ThingSave(ctx context.Context, thing *gorestapi.Thing) error {
	ret := _m.Called(ctx, thing)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorestapi.Thing) error); ok {
		r0 = rf(ctx, thing)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ThingsFind provides a mock function with given fields: ctx, qp
func (_m *GRStore) ThingsFind(ctx context.Context, qp *queryp.QueryParameters) ([]*gorestapi.Thing, *int64, error) {
	ret := _m.Called(ctx, qp)

	var r0 []*gorestapi.Thing
	var r1 *int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *queryp.QueryParameters) ([]*gorestapi.Thing, *int64, error)); ok {
		return rf(ctx, qp)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *queryp.QueryParameters) []*gorestapi.Thing); ok {
		r0 = rf(ctx, qp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*gorestapi.Thing)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *queryp.QueryParameters) *int64); ok {
		r1 = rf(ctx, qp)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*int64)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *queryp.QueryParameters) error); ok {
		r2 = rf(ctx, qp)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// WidgetDeleteByID provides a mock function with given fields: ctx, id
func (_m *GRStore) WidgetDeleteByID(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WidgetGetByID provides a mock function with given fields: ctx, id
func (_m *GRStore) WidgetGetByID(ctx context.Context, id string) (*gorestapi.Widget, error) {
	ret := _m.Called(ctx, id)

	var r0 *gorestapi.Widget
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*gorestapi.Widget, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *gorestapi.Widget); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorestapi.Widget)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WidgetSave provides a mock function with given fields: ctx, thing
func (_m *GRStore) WidgetSave(ctx context.Context, thing *gorestapi.Widget) error {
	ret := _m.Called(ctx, thing)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorestapi.Widget) error); ok {
		r0 = rf(ctx, thing)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WidgetsFind provides a mock function with given fields: ctx, qp
func (_m *GRStore) WidgetsFind(ctx context.Context, qp *queryp.QueryParameters) ([]*gorestapi.Widget, *int64, error) {
	ret := _m.Called(ctx, qp)

	var r0 []*gorestapi.Widget
	var r1 *int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *queryp.QueryParameters) ([]*gorestapi.Widget, *int64, error)); ok {
		return rf(ctx, qp)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *queryp.QueryParameters) []*gorestapi.Widget); ok {
		r0 = rf(ctx, qp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*gorestapi.Widget)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *queryp.QueryParameters) *int64); ok {
		r1 = rf(ctx, qp)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*int64)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *queryp.QueryParameters) error); ok {
		r2 = rf(ctx, qp)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewGRStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewGRStore creates a new instance of GRStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGRStore(t mockConstructorTestingTNewGRStore) *GRStore {
	mock := &GRStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
