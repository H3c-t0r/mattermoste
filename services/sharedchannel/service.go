// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannel

import (
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/services/filesstore"
	"github.com/mattermost/mattermost-server/v5/services/remotecluster"
	"github.com/mattermost/mattermost-server/v5/store"
)

const (
	TopicSync                    = "sharedchannel_sync"
	TopicChannelInvite           = "sharedchannel_invite"
	TopicUploadCreate            = "sharedchannel_upload"
	MaxRetries                   = 3
	MaxPostsPerSync              = 12 // a bit more than one typical screenfull of posts
	NotifyRemoteOfflineThreshold = time.Second * 10
	NotifyMinimumDelay           = time.Second * 2
)

// Mocks can be re-generated with `make sharedchannel-mocks`.

type ServerIface interface {
	Config() *model.Config
	IsLeader() bool
	AddClusterLeaderChangedListener(listener func()) string
	RemoveClusterLeaderChangedListener(id string)
	GetStore() store.Store
	GetLogger() mlog.LoggerIFace
	GetRemoteClusterService() *remotecluster.Service
}

type AppIface interface {
	SendEphemeralPost(userId string, post *model.Post) *model.Post
	CreateChannelWithUser(channel *model.Channel, userId string) (*model.Channel, *model.AppError)
	GetOrCreateDirectChannel(userId, otherUserId string, channelOptions ...model.ChannelOption) (*model.Channel, *model.AppError)
	AddUserToChannel(user *model.User, channel *model.Channel) (*model.ChannelMember, *model.AppError)
	AddUserToTeamByTeamId(teamId string, user *model.User) *model.AppError
	PermanentDeleteChannel(channel *model.Channel) *model.AppError
	CreatePost(post *model.Post, channel *model.Channel, triggerWebhooks bool, setOnline bool) (savedPost *model.Post, err *model.AppError)
	UpdatePost(post *model.Post, safeUpdate bool) (*model.Post, *model.AppError)
	SaveReactionForPost(reaction *model.Reaction) (*model.Reaction, *model.AppError)
	DeleteReactionForPost(reaction *model.Reaction) *model.AppError
	PatchChannelModerationsForChannel(channel *model.Channel, channelModerationsPatch []*model.ChannelModerationPatch) ([]*model.ChannelModeration, *model.AppError)
	CreateUploadSession(us *model.UploadSession) (*model.UploadSession, *model.AppError)
	FileReader(path string) (filesstore.ReadCloseSeeker, *model.AppError)
}

// errNotFound allows checking against Store.ErrNotFound errors without making Store a dependency.
type errNotFound interface {
	IsErrNotFound() bool
}

// errInvalidInput allows checking against Store.ErrInvalidInput errors without making Store a dependency.
type errInvalidInput interface {
	InvalidInputInfo() (entity string, field string, value interface{})
}

// Service provides shared channel synchronization.
type Service struct {
	server       ServerIface
	app          AppIface
	changeSignal chan struct{}

	// everything below guarded by `mux`
	mux                       sync.RWMutex
	active                    bool
	leaderListenerId          string
	connectionStateListenerId string
	done                      chan struct{}
	tasks                     map[string]syncTask
	syncTopicListenerId       string
	inviteTopicListenerId     string
	uploadTopicListenerId     string
	siteURL                   *url.URL
}

// NewSharedChannelService creates a RemoteClusterService instance.
func NewSharedChannelService(server ServerIface, app AppIface) (*Service, error) {
	service := &Service{
		server:       server,
		app:          app,
		changeSignal: make(chan struct{}, 1),
		tasks:        make(map[string]syncTask),
	}
	parsed, err := url.Parse(*server.Config().ServiceSettings.SiteURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse SiteURL: %w", err)
	}
	service.siteURL = parsed
	return service, nil
}

// Start is called by the server on server start-up.
func (scs *Service) Start() error {
	rcs := scs.server.GetRemoteClusterService()
	if rcs == nil {
		return errors.New("Shared Channel Service cannot activate: requires Remote Cluster Service")
	}

	scs.mux.Lock()
	scs.leaderListenerId = scs.server.AddClusterLeaderChangedListener(scs.onClusterLeaderChange)
	scs.syncTopicListenerId = rcs.AddTopicListener(TopicSync, scs.onReceiveSyncMessage)
	scs.inviteTopicListenerId = rcs.AddTopicListener(TopicChannelInvite, scs.onReceiveChannelInvite)
	scs.uploadTopicListenerId = rcs.AddTopicListener(TopicUploadCreate, scs.onReceiveUploadCreate)
	scs.connectionStateListenerId = rcs.AddConnectionStateListener(scs.onConnectionStateChange)
	scs.mux.Unlock()

	scs.onClusterLeaderChange()

	return nil
}

// Shutdown is called by the server on server shutdown.
func (scs *Service) Shutdown() error {
	rcs := scs.server.GetRemoteClusterService()
	if rcs == nil {
		return errors.New("Shared Channel Service cannot shutdown: requires Remote Cluster Service")
	}

	scs.mux.Lock()
	id := scs.leaderListenerId
	rcs.RemoveTopicListener(scs.syncTopicListenerId)
	scs.syncTopicListenerId = ""
	rcs.RemoveTopicListener(scs.inviteTopicListenerId)
	scs.inviteTopicListenerId = ""
	rcs.RemoveConnectionStateListener(scs.connectionStateListenerId)
	scs.connectionStateListenerId = ""
	scs.mux.Unlock()

	scs.server.RemoveClusterLeaderChangedListener(id)
	scs.pause()
	return nil
}

// Active determines whether the service is active on the node or not.
func (scs *Service) Active() bool {
	scs.mux.Lock()
	defer scs.mux.Unlock()

	return scs.active
}

func (scs *Service) sendEphemeralPost(channelId string, userId string, text string) {
	ephemeral := &model.Post{
		ChannelId: channelId,
		Message:   text,
		CreateAt:  model.GetMillis(),
	}
	scs.app.SendEphemeralPost(userId, ephemeral)
}

// onClusterLeaderChange is called whenever the cluster leader may have changed.
func (scs *Service) onClusterLeaderChange() {
	if scs.server.IsLeader() {
		scs.resume()
	} else {
		scs.pause()
	}
}

func (scs *Service) resume() {
	scs.mux.Lock()
	defer scs.mux.Unlock()

	if scs.active {
		return // already active
	}

	scs.active = true
	scs.done = make(chan struct{})

	go scs.syncLoop(scs.done)

	scs.server.GetLogger().Debug("Shared Channel Service active")
}

func (scs *Service) pause() {
	scs.mux.Lock()
	defer scs.mux.Unlock()

	if !scs.active {
		return // already inactive
	}

	scs.active = false
	close(scs.done)
	scs.done = nil

	scs.server.GetLogger().Debug("Shared Channel Service inactive")
}

// Makes the remote channel to be read-only(announcement mode, only admins can create posts and reactions).
func (scs *Service) makeChannelReadOnly(channel *model.Channel) *model.AppError {
	createPostPermission := model.ChannelModeratedPermissionsMap[model.PERMISSION_CREATE_POST.Id]
	createReactionPermission := model.ChannelModeratedPermissionsMap[model.PERMISSION_ADD_REACTION.Id]
	updateMap := model.ChannelModeratedRolesPatch{
		Guests:  model.NewBool(false),
		Members: model.NewBool(false),
	}

	readonlyChannelModerations := []*model.ChannelModerationPatch{
		{
			Name:  &createPostPermission,
			Roles: &updateMap,
		},
		{
			Name:  &createReactionPermission,
			Roles: &updateMap,
		},
	}

	_, err := scs.app.PatchChannelModerationsForChannel(channel, readonlyChannelModerations)
	return err
}

// onConnectionStateChange is called whenever the connection state of a remote cluster changes,
// for example when one comes back online.
func (scs *Service) onConnectionStateChange(rc *model.RemoteCluster, online bool) {
	if online {
		// when a previously offline remote comes back online force a sync.
		scs.ForceSyncForRemote(rc)
	}

	scs.server.GetLogger().Error("Remote cluster connection status changed",
		mlog.String("remoteId", rc.RemoteId),
		mlog.Bool("online", online),
	)
}
