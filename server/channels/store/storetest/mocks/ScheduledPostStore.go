// Code generated by mockery v2.42.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/mattermost/mattermost/server/public/model"
	mock "github.com/stretchr/testify/mock"
)

// ScheduledPostStore is an autogenerated mock type for the ScheduledPostStore type
type ScheduledPostStore struct {
	mock.Mock
}

// CreateScheduledPost provides a mock function with given fields: scheduledPost
func (_m *ScheduledPostStore) CreateScheduledPost(scheduledPost *model.ScheduledPost) (*model.ScheduledPost, error) {
	ret := _m.Called(scheduledPost)

	if len(ret) == 0 {
		panic("no return value specified for CreateScheduledPost")
	}

	var r0 *model.ScheduledPost
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.ScheduledPost) (*model.ScheduledPost, error)); ok {
		return rf(scheduledPost)
	}
	if rf, ok := ret.Get(0).(func(*model.ScheduledPost) *model.ScheduledPost); ok {
		r0 = rf(scheduledPost)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ScheduledPost)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.ScheduledPost) error); ok {
		r1 = rf(scheduledPost)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetScheduledPosts provides a mock function with given fields: beforeTime, lastScheduledPostId, perPage
func (_m *ScheduledPostStore) GetScheduledPosts(beforeTime int64, lastScheduledPostId string, perPage uint64) ([]*model.ScheduledPost, error) {
	ret := _m.Called(beforeTime, lastScheduledPostId, perPage)

	if len(ret) == 0 {
		panic("no return value specified for GetScheduledPosts")
	}

	var r0 []*model.ScheduledPost
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, string, uint64) ([]*model.ScheduledPost, error)); ok {
		return rf(beforeTime, lastScheduledPostId, perPage)
	}
	if rf, ok := ret.Get(0).(func(int64, string, uint64) []*model.ScheduledPost); ok {
		r0 = rf(beforeTime, lastScheduledPostId, perPage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ScheduledPost)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, string, uint64) error); ok {
		r1 = rf(beforeTime, lastScheduledPostId, perPage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetScheduledPostsForUser provides a mock function with given fields: userId, teamId
func (_m *ScheduledPostStore) GetScheduledPostsForUser(userId string, teamId string) ([]*model.ScheduledPost, error) {
	ret := _m.Called(userId, teamId)

	if len(ret) == 0 {
		panic("no return value specified for GetScheduledPostsForUser")
	}

	var r0 []*model.ScheduledPost
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) ([]*model.ScheduledPost, error)); ok {
		return rf(userId, teamId)
	}
	if rf, ok := ret.Get(0).(func(string, string) []*model.ScheduledPost); ok {
		r0 = rf(userId, teamId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ScheduledPost)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(userId, teamId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewScheduledPostStore creates a new instance of ScheduledPostStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewScheduledPostStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *ScheduledPostStore {
	mock := &ScheduledPostStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
