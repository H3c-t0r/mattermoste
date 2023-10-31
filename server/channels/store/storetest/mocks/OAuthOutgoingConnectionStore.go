// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/mattermost/mattermost/server/public/model"
	request "github.com/mattermost/mattermost/server/public/shared/request"
	mock "github.com/stretchr/testify/mock"
)

// OAuthOutgoingConnectionStore is an autogenerated mock type for the OAuthOutgoingConnectionStore type
type OAuthOutgoingConnectionStore struct {
	mock.Mock
}

// DeleteConnection provides a mock function with given fields: c, id
func (_m *OAuthOutgoingConnectionStore) DeleteConnection(c request.CTX, id string) error {
	ret := _m.Called(c, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(request.CTX, string) error); ok {
		r0 = rf(c, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetConnection provides a mock function with given fields: c, id
func (_m *OAuthOutgoingConnectionStore) GetConnection(c request.CTX, id string) (*model.OAuthOutgoingConnection, error) {
	ret := _m.Called(c, id)

	var r0 *model.OAuthOutgoingConnection
	var r1 error
	if rf, ok := ret.Get(0).(func(request.CTX, string) (*model.OAuthOutgoingConnection, error)); ok {
		return rf(c, id)
	}
	if rf, ok := ret.Get(0).(func(request.CTX, string) *model.OAuthOutgoingConnection); ok {
		r0 = rf(c, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OAuthOutgoingConnection)
		}
	}

	if rf, ok := ret.Get(1).(func(request.CTX, string) error); ok {
		r1 = rf(c, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConnections provides a mock function with given fields: c, offset, limit
func (_m *OAuthOutgoingConnectionStore) GetConnections(c request.CTX, offset int, limit int) ([]*model.OAuthOutgoingConnection, error) {
	ret := _m.Called(c, offset, limit)

	var r0 []*model.OAuthOutgoingConnection
	var r1 error
	if rf, ok := ret.Get(0).(func(request.CTX, int, int) ([]*model.OAuthOutgoingConnection, error)); ok {
		return rf(c, offset, limit)
	}
	if rf, ok := ret.Get(0).(func(request.CTX, int, int) []*model.OAuthOutgoingConnection); ok {
		r0 = rf(c, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.OAuthOutgoingConnection)
		}
	}

	if rf, ok := ret.Get(1).(func(request.CTX, int, int) error); ok {
		r1 = rf(c, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveConnection provides a mock function with given fields: c, conn
func (_m *OAuthOutgoingConnectionStore) SaveConnection(c request.CTX, conn *model.OAuthOutgoingConnection) (*model.OAuthOutgoingConnection, error) {
	ret := _m.Called(c, conn)

	var r0 *model.OAuthOutgoingConnection
	var r1 error
	if rf, ok := ret.Get(0).(func(request.CTX, *model.OAuthOutgoingConnection) (*model.OAuthOutgoingConnection, error)); ok {
		return rf(c, conn)
	}
	if rf, ok := ret.Get(0).(func(request.CTX, *model.OAuthOutgoingConnection) *model.OAuthOutgoingConnection); ok {
		r0 = rf(c, conn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OAuthOutgoingConnection)
		}
	}

	if rf, ok := ret.Get(1).(func(request.CTX, *model.OAuthOutgoingConnection) error); ok {
		r1 = rf(c, conn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateConnection provides a mock function with given fields: c, conn
func (_m *OAuthOutgoingConnectionStore) UpdateConnection(c request.CTX, conn *model.OAuthOutgoingConnection) (*model.OAuthOutgoingConnection, error) {
	ret := _m.Called(c, conn)

	var r0 *model.OAuthOutgoingConnection
	var r1 error
	if rf, ok := ret.Get(0).(func(request.CTX, *model.OAuthOutgoingConnection) (*model.OAuthOutgoingConnection, error)); ok {
		return rf(c, conn)
	}
	if rf, ok := ret.Get(0).(func(request.CTX, *model.OAuthOutgoingConnection) *model.OAuthOutgoingConnection); ok {
		r0 = rf(c, conn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OAuthOutgoingConnection)
		}
	}

	if rf, ok := ret.Get(1).(func(request.CTX, *model.OAuthOutgoingConnection) error); ok {
		r1 = rf(c, conn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewOAuthOutgoingConnectionStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewOAuthOutgoingConnectionStore creates a new instance of OAuthOutgoingConnectionStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOAuthOutgoingConnectionStore(t mockConstructorTestingTNewOAuthOutgoingConnectionStore) *OAuthOutgoingConnectionStore {
	mock := &OAuthOutgoingConnectionStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
