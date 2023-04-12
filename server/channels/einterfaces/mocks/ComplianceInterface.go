// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make einterfaces-mocks`.

package mocks

import (
	model "github.com/mattermost/mattermost-server/v6/model"
	mock "github.com/stretchr/testify/mock"
)

// ComplianceInterface is an autogenerated mock type for the ComplianceInterface type
type ComplianceInterface struct {
	mock.Mock
}

// RunComplianceJob provides a mock function with given fields: job
func (_m *ComplianceInterface) RunComplianceJob(job *model.Compliance) *model.AppError {
	ret := _m.Called(job)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(*model.Compliance) *model.AppError); ok {
		r0 = rf(job)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// StartComplianceDailyJob provides a mock function with given fields:
func (_m *ComplianceInterface) StartComplianceDailyJob() {
	_m.Called()
}

type mockConstructorTestingTNewComplianceInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewComplianceInterface creates a new instance of ComplianceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewComplianceInterface(t mockConstructorTestingTNewComplianceInterface) *ComplianceInterface {
	mock := &ComplianceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
