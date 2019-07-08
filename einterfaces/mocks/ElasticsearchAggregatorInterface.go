// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make einterfaces-mocks`.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/mattermost/mattermost-server/model"

// ElasticsearchAggregatorInterface is an autogenerated mock type for the ElasticsearchAggregatorInterface type
type ElasticsearchAggregatorInterface struct {
	mock.Mock
}

// MakeScheduler provides a mock function with given fields:
func (_m *ElasticsearchAggregatorInterface) MakeScheduler() model.Scheduler {
	ret := _m.Called()

	var r0 model.Scheduler
	if rf, ok := ret.Get(0).(func() model.Scheduler); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Scheduler)
		}
	}

	return r0
}

// MakeWorker provides a mock function with given fields:
func (_m *ElasticsearchAggregatorInterface) MakeWorker() model.Worker {
	ret := _m.Called()

	var r0 model.Worker
	if rf, ok := ret.Get(0).(func() model.Worker); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Worker)
		}
	}

	return r0
}
