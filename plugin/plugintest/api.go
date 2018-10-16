// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make plugin-mocks`.

package plugintest

import mock "github.com/stretchr/testify/mock"
import model "github.com/mattermost/mattermost-server/model"

// API is an autogenerated mock type for the API type
type API struct {
	mock.Mock
}

// AddChannelMember provides a mock function with given fields: channelId, userId
func (_m *API) AddChannelMember(channelId string, userId string) (*model.ChannelMember, *model.AppError) {
	ret := _m.Called(channelId, userId)

	var r0 *model.ChannelMember
	if rf, ok := ret.Get(0).(func(string, string) *model.ChannelMember); ok {
		r0 = rf(channelId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ChannelMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(channelId, userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// AddReaction provides a mock function with given fields: reaction
func (_m *API) AddReaction(reaction *model.Reaction) (*model.Reaction, *model.AppError) {
	ret := _m.Called(reaction)

	var r0 *model.Reaction
	if rf, ok := ret.Get(0).(func(*model.Reaction) *model.Reaction); ok {
		r0 = rf(reaction)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Reaction)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Reaction) *model.AppError); ok {
		r1 = rf(reaction)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// CopyFileInfos provides a mock function with given fields: userId, fileIds
func (_m *API) CopyFileInfos(userId string, fileIds []string) ([]string, *model.AppError) {
	ret := _m.Called(userId, fileIds)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string, []string) []string); ok {
		r0 = rf(userId, fileIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, []string) *model.AppError); ok {
		r1 = rf(userId, fileIds)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// CreateChannel provides a mock function with given fields: channel
func (_m *API) CreateChannel(channel *model.Channel) (*model.Channel, *model.AppError) {
	ret := _m.Called(channel)

	var r0 *model.Channel
	if rf, ok := ret.Get(0).(func(*model.Channel) *model.Channel); ok {
		r0 = rf(channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Channel)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Channel) *model.AppError); ok {
		r1 = rf(channel)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// CreatePost provides a mock function with given fields: post
func (_m *API) CreatePost(post *model.Post) (*model.Post, *model.AppError) {
	ret := _m.Called(post)

	var r0 *model.Post
	if rf, ok := ret.Get(0).(func(*model.Post) *model.Post); ok {
		r0 = rf(post)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Post)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Post) *model.AppError); ok {
		r1 = rf(post)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// CreateTeam provides a mock function with given fields: team
func (_m *API) CreateTeam(team *model.Team) (*model.Team, *model.AppError) {
	ret := _m.Called(team)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(*model.Team) *model.Team); ok {
		r0 = rf(team)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Team) *model.AppError); ok {
		r1 = rf(team)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// CreateTeamMember provides a mock function with given fields: teamId, userId
func (_m *API) CreateTeamMember(teamId string, userId string) (*model.TeamMember, *model.AppError) {
	ret := _m.Called(teamId, userId)

	var r0 *model.TeamMember
	if rf, ok := ret.Get(0).(func(string, string) *model.TeamMember); ok {
		r0 = rf(teamId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TeamMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(teamId, userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// CreateTeamMembers provides a mock function with given fields: teamId, userIds, requestorId
func (_m *API) CreateTeamMembers(teamId string, userIds []string, requestorId string) ([]*model.TeamMember, *model.AppError) {
	ret := _m.Called(teamId, userIds, requestorId)

	var r0 []*model.TeamMember
	if rf, ok := ret.Get(0).(func(string, []string, string) []*model.TeamMember); ok {
		r0 = rf(teamId, userIds, requestorId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, []string, string) *model.AppError); ok {
		r1 = rf(teamId, userIds, requestorId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: user
func (_m *API) CreateUser(user *model.User) (*model.User, *model.AppError) {
	ret := _m.Called(user)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(*model.User) *model.User); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.User) *model.AppError); ok {
		r1 = rf(user)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// DeleteChannel provides a mock function with given fields: channelId
func (_m *API) DeleteChannel(channelId string) *model.AppError {
	ret := _m.Called(channelId)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string) *model.AppError); ok {
		r0 = rf(channelId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// DeleteChannelMember provides a mock function with given fields: channelId, userId
func (_m *API) DeleteChannelMember(channelId string, userId string) *model.AppError {
	ret := _m.Called(channelId, userId)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string, string) *model.AppError); ok {
		r0 = rf(channelId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// DeletePost provides a mock function with given fields: postId
func (_m *API) DeletePost(postId string) *model.AppError {
	ret := _m.Called(postId)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string) *model.AppError); ok {
		r0 = rf(postId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// DeleteTeam provides a mock function with given fields: teamId
func (_m *API) DeleteTeam(teamId string) *model.AppError {
	ret := _m.Called(teamId)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string) *model.AppError); ok {
		r0 = rf(teamId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// DeleteTeamMember provides a mock function with given fields: teamId, userId, requestorId
func (_m *API) DeleteTeamMember(teamId string, userId string, requestorId string) *model.AppError {
	ret := _m.Called(teamId, userId, requestorId)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string, string, string) *model.AppError); ok {
		r0 = rf(teamId, userId, requestorId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// DeleteUser provides a mock function with given fields: userId
func (_m *API) DeleteUser(userId string) *model.AppError {
	ret := _m.Called(userId)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string) *model.AppError); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// GetChannel provides a mock function with given fields: channelId
func (_m *API) GetChannel(channelId string) (*model.Channel, *model.AppError) {
	ret := _m.Called(channelId)

	var r0 *model.Channel
	if rf, ok := ret.Get(0).(func(string) *model.Channel); ok {
		r0 = rf(channelId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Channel)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(channelId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetChannelByName provides a mock function with given fields: teamId, name, includeDeleted
func (_m *API) GetChannelByName(teamId string, name string, includeDeleted bool) (*model.Channel, *model.AppError) {
	ret := _m.Called(teamId, name, includeDeleted)

	var r0 *model.Channel
	if rf, ok := ret.Get(0).(func(string, string, bool) *model.Channel); ok {
		r0 = rf(teamId, name, includeDeleted)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Channel)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string, bool) *model.AppError); ok {
		r1 = rf(teamId, name, includeDeleted)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetChannelByNameForTeamName provides a mock function with given fields: teamName, channelName, includeDeleted
func (_m *API) GetChannelByNameForTeamName(teamName string, channelName string, includeDeleted bool) (*model.Channel, *model.AppError) {
	ret := _m.Called(teamName, channelName, includeDeleted)

	var r0 *model.Channel
	if rf, ok := ret.Get(0).(func(string, string, bool) *model.Channel); ok {
		r0 = rf(teamName, channelName, includeDeleted)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Channel)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string, bool) *model.AppError); ok {
		r1 = rf(teamName, channelName, includeDeleted)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetChannelMember provides a mock function with given fields: channelId, userId
func (_m *API) GetChannelMember(channelId string, userId string) (*model.ChannelMember, *model.AppError) {
	ret := _m.Called(channelId, userId)

	var r0 *model.ChannelMember
	if rf, ok := ret.Get(0).(func(string, string) *model.ChannelMember); ok {
		r0 = rf(channelId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ChannelMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(channelId, userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetChannelMembers provides a mock function with given fields: channelId, page, perPage
func (_m *API) GetChannelMembers(channelId string, page int, perPage int) (*model.ChannelMembers, *model.AppError) {
	ret := _m.Called(channelId, page, perPage)

	var r0 *model.ChannelMembers
	if rf, ok := ret.Get(0).(func(string, int, int) *model.ChannelMembers); ok {
		r0 = rf(channelId, page, perPage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ChannelMembers)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, int, int) *model.AppError); ok {
		r1 = rf(channelId, page, perPage)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetConfig provides a mock function with given fields:
func (_m *API) GetConfig() *model.Config {
	ret := _m.Called()

	var r0 *model.Config
	if rf, ok := ret.Get(0).(func() *model.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Config)
		}
	}

	return r0
}

// GetDirectChannel provides a mock function with given fields: userId1, userId2
func (_m *API) GetDirectChannel(userId1 string, userId2 string) (*model.Channel, *model.AppError) {
	ret := _m.Called(userId1, userId2)

	var r0 *model.Channel
	if rf, ok := ret.Get(0).(func(string, string) *model.Channel); ok {
		r0 = rf(userId1, userId2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Channel)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(userId1, userId2)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetFileInfo provides a mock function with given fields: fileId
func (_m *API) GetFileInfo(fileId string) (*model.FileInfo, *model.AppError) {
	ret := _m.Called(fileId)

	var r0 *model.FileInfo
	if rf, ok := ret.Get(0).(func(string) *model.FileInfo); ok {
		r0 = rf(fileId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FileInfo)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(fileId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetGroupChannel provides a mock function with given fields: userIds
func (_m *API) GetGroupChannel(userIds []string) (*model.Channel, *model.AppError) {
	ret := _m.Called(userIds)

	var r0 *model.Channel
	if rf, ok := ret.Get(0).(func([]string) *model.Channel); ok {
		r0 = rf(userIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Channel)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func([]string) *model.AppError); ok {
		r1 = rf(userIds)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetLDAPUserAttributes provides a mock function with given fields: userId, attributes
func (_m *API) GetLDAPUserAttributes(userId string, attributes []string) (map[string]string, *model.AppError) {
	ret := _m.Called(userId, attributes)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(string, []string) map[string]string); ok {
		r0 = rf(userId, attributes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, []string) *model.AppError); ok {
		r1 = rf(userId, attributes)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetPost provides a mock function with given fields: postId
func (_m *API) GetPost(postId string) (*model.Post, *model.AppError) {
	ret := _m.Called(postId)

	var r0 *model.Post
	if rf, ok := ret.Get(0).(func(string) *model.Post); ok {
		r0 = rf(postId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Post)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(postId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetPostsForChannel provides a mock function with given fields: channelId, page, perPage
func (_m *API) GetPostsForChannel(channelId string, page int, perPage int) (*model.PostList, *model.AppError) {
	ret := _m.Called(channelId, page, perPage)

	var r0 *model.PostList
	if rf, ok := ret.Get(0).(func(string, int, int) *model.PostList); ok {
		r0 = rf(channelId, page, perPage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PostList)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, int, int) *model.AppError); ok {
		r1 = rf(channelId, page, perPage)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetPublicChannelsForTeam provides a mock function with given fields: teamId, offset, limit
func (_m *API) GetPublicChannelsForTeam(teamId string, offset int, limit int) (*model.ChannelList, *model.AppError) {
	ret := _m.Called(teamId, offset, limit)

	var r0 *model.ChannelList
	if rf, ok := ret.Get(0).(func(string, int, int) *model.ChannelList); ok {
		r0 = rf(teamId, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ChannelList)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, int, int) *model.AppError); ok {
		r1 = rf(teamId, offset, limit)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetReactions provides a mock function with given fields: postId
func (_m *API) GetReactions(postId string) ([]*model.Reaction, *model.AppError) {
	ret := _m.Called(postId)

	var r0 []*model.Reaction
	if rf, ok := ret.Get(0).(func(string) []*model.Reaction); ok {
		r0 = rf(postId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Reaction)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(postId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetServerVersion provides a mock function with given fields:
func (_m *API) GetServerVersion() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetSession provides a mock function with given fields: sessionId
func (_m *API) GetSession(sessionId string) (*model.Session, *model.AppError) {
	ret := _m.Called(sessionId)

	var r0 *model.Session
	if rf, ok := ret.Get(0).(func(string) *model.Session); ok {
		r0 = rf(sessionId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Session)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(sessionId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetTeam provides a mock function with given fields: teamId
func (_m *API) GetTeam(teamId string) (*model.Team, *model.AppError) {
	ret := _m.Called(teamId)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(string) *model.Team); ok {
		r0 = rf(teamId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(teamId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetTeamByName provides a mock function with given fields: name
func (_m *API) GetTeamByName(name string) (*model.Team, *model.AppError) {
	ret := _m.Called(name)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(string) *model.Team); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(name)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetTeamMember provides a mock function with given fields: teamId, userId
func (_m *API) GetTeamMember(teamId string, userId string) (*model.TeamMember, *model.AppError) {
	ret := _m.Called(teamId, userId)

	var r0 *model.TeamMember
	if rf, ok := ret.Get(0).(func(string, string) *model.TeamMember); ok {
		r0 = rf(teamId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TeamMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(teamId, userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetTeamMembers provides a mock function with given fields: teamId, offset, limit
func (_m *API) GetTeamMembers(teamId string, offset int, limit int) ([]*model.TeamMember, *model.AppError) {
	ret := _m.Called(teamId, offset, limit)

	var r0 []*model.TeamMember
	if rf, ok := ret.Get(0).(func(string, int, int) []*model.TeamMember); ok {
		r0 = rf(teamId, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, int, int) *model.AppError); ok {
		r1 = rf(teamId, offset, limit)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetTeams provides a mock function with given fields:
func (_m *API) GetTeams() ([]*model.Team, *model.AppError) {
	ret := _m.Called()

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func() []*model.Team); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func() *model.AppError); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetTeamsUnreadForUser provides a mock function with given fields: userId
func (_m *API) GetTeamsUnreadForUser(userId string) ([]*model.TeamUnread, *model.AppError) {
	ret := _m.Called(userId)

	var r0 []*model.TeamUnread
	if rf, ok := ret.Get(0).(func(string) []*model.TeamUnread); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamUnread)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: userId
func (_m *API) GetUser(userId string) (*model.User, *model.AppError) {
	ret := _m.Called(userId)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string) *model.User); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetUserByEmail provides a mock function with given fields: email
func (_m *API) GetUserByEmail(email string) (*model.User, *model.AppError) {
	ret := _m.Called(email)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string) *model.User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(email)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetUserByUsername provides a mock function with given fields: name
func (_m *API) GetUserByUsername(name string) (*model.User, *model.AppError) {
	ret := _m.Called(name)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string) *model.User); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(name)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetUserStatus provides a mock function with given fields: userId
func (_m *API) GetUserStatus(userId string) (*model.Status, *model.AppError) {
	ret := _m.Called(userId)

	var r0 *model.Status
	if rf, ok := ret.Get(0).(func(string) *model.Status); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Status)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetUserStatusesByIds provides a mock function with given fields: userIds
func (_m *API) GetUserStatusesByIds(userIds []string) ([]*model.Status, *model.AppError) {
	ret := _m.Called(userIds)

	var r0 []*model.Status
	if rf, ok := ret.Get(0).(func([]string) []*model.Status); ok {
		r0 = rf(userIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Status)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func([]string) *model.AppError); ok {
		r1 = rf(userIds)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// HasPermissionTo provides a mock function with given fields: userId, permission
func (_m *API) HasPermissionTo(userId string, permission *model.Permission) bool {
	ret := _m.Called(userId, permission)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, *model.Permission) bool); ok {
		r0 = rf(userId, permission)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// HasPermissionToChannel provides a mock function with given fields: userId, channelId, permission
func (_m *API) HasPermissionToChannel(userId string, channelId string, permission *model.Permission) bool {
	ret := _m.Called(userId, channelId, permission)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string, *model.Permission) bool); ok {
		r0 = rf(userId, channelId, permission)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// HasPermissionToTeam provides a mock function with given fields: userId, teamId, permission
func (_m *API) HasPermissionToTeam(userId string, teamId string, permission *model.Permission) bool {
	ret := _m.Called(userId, teamId, permission)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string, *model.Permission) bool); ok {
		r0 = rf(userId, teamId, permission)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// KVDelete provides a mock function with given fields: key
func (_m *API) KVDelete(key string) *model.AppError {
	ret := _m.Called(key)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string) *model.AppError); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// KVDeleteAll provides a mock function with given fields:
func (_m *API) KVDeleteAll() *model.AppError {
	ret := _m.Called()

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func() *model.AppError); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// KVGet provides a mock function with given fields: key
func (_m *API) KVGet(key string) ([]byte, *model.AppError) {
	ret := _m.Called(key)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(key)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// KVList provides a mock function with given fields: page, perPage
func (_m *API) KVList(page int, perPage int) ([]string, *model.AppError) {
	ret := _m.Called(page, perPage)

	var r0 []string
	if rf, ok := ret.Get(0).(func(int, int) []string); ok {
		r0 = rf(page, perPage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(int, int) *model.AppError); ok {
		r1 = rf(page, perPage)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// KVSet provides a mock function with given fields: key, value
func (_m *API) KVSet(key string, value []byte) *model.AppError {
	ret := _m.Called(key, value)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string, []byte) *model.AppError); ok {
		r0 = rf(key, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// KVSetWithExpiry provides a mock function with given fields: key, value, expireInSeconds
func (_m *API) KVSetWithExpiry(key string, value []byte, expireInSeconds int64) *model.AppError {
	ret := _m.Called(key, value, expireInSeconds)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string, []byte, int64) *model.AppError); ok {
		r0 = rf(key, value, expireInSeconds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// LoadPluginConfiguration provides a mock function with given fields: dest
func (_m *API) LoadPluginConfiguration(dest interface{}) error {
	ret := _m.Called(dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LogDebug provides a mock function with given fields: msg, keyValuePairs
func (_m *API) LogDebug(msg string, keyValuePairs ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, keyValuePairs...)
	_m.Called(_ca...)
}

// LogError provides a mock function with given fields: msg, keyValuePairs
func (_m *API) LogError(msg string, keyValuePairs ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, keyValuePairs...)
	_m.Called(_ca...)
}

// LogInfo provides a mock function with given fields: msg, keyValuePairs
func (_m *API) LogInfo(msg string, keyValuePairs ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, keyValuePairs...)
	_m.Called(_ca...)
}

// LogWarn provides a mock function with given fields: msg, keyValuePairs
func (_m *API) LogWarn(msg string, keyValuePairs ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, keyValuePairs...)
	_m.Called(_ca...)
}

// PublishWebSocketEvent provides a mock function with given fields: event, payload, broadcast
func (_m *API) PublishWebSocketEvent(event string, payload map[string]interface{}, broadcast *model.WebsocketBroadcast) {
	_m.Called(event, payload, broadcast)
}

// ReadFile provides a mock function with given fields: path
func (_m *API) ReadFile(path string) ([]byte, *model.AppError) {
	ret := _m.Called(path)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(path)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// RegisterCommand provides a mock function with given fields: command
func (_m *API) RegisterCommand(command *model.Command) error {
	ret := _m.Called(command)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Command) error); ok {
		r0 = rf(command)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveReaction provides a mock function with given fields: reaction
func (_m *API) RemoveReaction(reaction *model.Reaction) *model.AppError {
	ret := _m.Called(reaction)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(*model.Reaction) *model.AppError); ok {
		r0 = rf(reaction)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// SaveConfig provides a mock function with given fields: config
func (_m *API) SaveConfig(config *model.Config) *model.AppError {
	ret := _m.Called(config)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(*model.Config) *model.AppError); ok {
		r0 = rf(config)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// SendEphemeralPost provides a mock function with given fields: userId, post
func (_m *API) SendEphemeralPost(userId string, post *model.Post) *model.Post {
	ret := _m.Called(userId, post)

	var r0 *model.Post
	if rf, ok := ret.Get(0).(func(string, *model.Post) *model.Post); ok {
		r0 = rf(userId, post)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Post)
		}
	}

	return r0
}

// UnregisterCommand provides a mock function with given fields: teamId, trigger
func (_m *API) UnregisterCommand(teamId string, trigger string) error {
	ret := _m.Called(teamId, trigger)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(teamId, trigger)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateChannel provides a mock function with given fields: channel
func (_m *API) UpdateChannel(channel *model.Channel) (*model.Channel, *model.AppError) {
	ret := _m.Called(channel)

	var r0 *model.Channel
	if rf, ok := ret.Get(0).(func(*model.Channel) *model.Channel); ok {
		r0 = rf(channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Channel)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Channel) *model.AppError); ok {
		r1 = rf(channel)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// UpdateChannelMemberNotifications provides a mock function with given fields: channelId, userId, notifications
func (_m *API) UpdateChannelMemberNotifications(channelId string, userId string, notifications map[string]string) (*model.ChannelMember, *model.AppError) {
	ret := _m.Called(channelId, userId, notifications)

	var r0 *model.ChannelMember
	if rf, ok := ret.Get(0).(func(string, string, map[string]string) *model.ChannelMember); ok {
		r0 = rf(channelId, userId, notifications)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ChannelMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string, map[string]string) *model.AppError); ok {
		r1 = rf(channelId, userId, notifications)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// UpdateChannelMemberRoles provides a mock function with given fields: channelId, userId, newRoles
func (_m *API) UpdateChannelMemberRoles(channelId string, userId string, newRoles string) (*model.ChannelMember, *model.AppError) {
	ret := _m.Called(channelId, userId, newRoles)

	var r0 *model.ChannelMember
	if rf, ok := ret.Get(0).(func(string, string, string) *model.ChannelMember); ok {
		r0 = rf(channelId, userId, newRoles)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ChannelMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string, string) *model.AppError); ok {
		r1 = rf(channelId, userId, newRoles)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// UpdatePost provides a mock function with given fields: post
func (_m *API) UpdatePost(post *model.Post) (*model.Post, *model.AppError) {
	ret := _m.Called(post)

	var r0 *model.Post
	if rf, ok := ret.Get(0).(func(*model.Post) *model.Post); ok {
		r0 = rf(post)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Post)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Post) *model.AppError); ok {
		r1 = rf(post)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// UpdateTeam provides a mock function with given fields: team
func (_m *API) UpdateTeam(team *model.Team) (*model.Team, *model.AppError) {
	ret := _m.Called(team)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(*model.Team) *model.Team); ok {
		r0 = rf(team)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.Team) *model.AppError); ok {
		r1 = rf(team)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// UpdateTeamMemberRoles provides a mock function with given fields: teamId, userId, newRoles
func (_m *API) UpdateTeamMemberRoles(teamId string, userId string, newRoles string) (*model.TeamMember, *model.AppError) {
	ret := _m.Called(teamId, userId, newRoles)

	var r0 *model.TeamMember
	if rf, ok := ret.Get(0).(func(string, string, string) *model.TeamMember); ok {
		r0 = rf(teamId, userId, newRoles)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TeamMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string, string) *model.AppError); ok {
		r1 = rf(teamId, userId, newRoles)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: user
func (_m *API) UpdateUser(user *model.User) (*model.User, *model.AppError) {
	ret := _m.Called(user)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(*model.User) *model.User); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.User) *model.AppError); ok {
		r1 = rf(user)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// UpdateUserStatus provides a mock function with given fields: userId, status
func (_m *API) UpdateUserStatus(userId string, status string) (*model.Status, *model.AppError) {
	ret := _m.Called(userId, status)

	var r0 *model.Status
	if rf, ok := ret.Get(0).(func(string, string) *model.Status); ok {
		r0 = rf(userId, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Status)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(userId, status)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}
