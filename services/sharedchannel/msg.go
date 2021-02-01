// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sharedchannel

import (
	"encoding/json"
	"time"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"
)

// syncMsg represents a change in content (post add/edit/delete, reaction add/remove, users).
// It is sent to remote clusters as the payload of a `RemoteClusterMsg`.
type syncMsg struct {
	ChannelId string            `json:"channel_id"`
	PostId    string            `json:"post_id"`
	Post      *model.Post       `json:"post"`
	Users     []*model.User     `json:"users"`
	Reactions []*model.Reaction `json:"reactions"`
}

func (sm syncMsg) ToJSON() ([]byte, error) {
	b, err := json.Marshal(sm)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (sm syncMsg) String() string {
	json, err := sm.ToJSON()
	if err != nil {
		return ""
	}
	return string(json)
}

// postsToMsg takes a slice of posts and converts to a `RemoteClusterMsg` which can be
// sent to a remote cluster.
func (scs *Service) postsToMsg(posts []*model.Post, rc *model.RemoteCluster, lastSyncAt int64) (model.RemoteClusterMsg, error) {
	syncMessages := make([]syncMsg, 0, len(posts))

	for _, p := range posts {
		if p.IsSystemMessage() { // don't sync system messages
			continue
		}

		// any reactions originating from the remote cluster are filtered out
		reactions, err := scs.server.GetStore().Reaction().GetForPostSince(p.Id, lastSyncAt, rc.RemoteId, true)
		if err != nil {
			return model.RemoteClusterMsg{}, err
		}

		postSync := p

		// TODO:  don't include the post if only the reactions changed. Unfortunately there is no way to reliably know the
		//        difference between an existing (synchronized) post with new reaction, and a brand new post (un-synchronized)
		//        with a reaction.
		//if p.EditAt < p.UpdateAt && p.CreateAt < p.UpdateAt && p.DeleteAt == 0 {
		//	postSync = nil
		//}

		// don't sync a post back to the remote it came from.
		if p.RemoteId != nil && *p.RemoteId == rc.RemoteId {
			postSync = nil
		}

		// Parse out all permalinks in the message.
		if postSync != nil {
			postSync.Message = scs.processPermalinkToRemote(postSync)
		}

		// any users originating from the remote cluster are filtered out
		users := scs.usersForPost(postSync, reactions, rc)

		// if everything was filtered out then don't send an empty message.
		if postSync == nil && len(reactions) == 0 && len(users) == 0 {
			continue
		}

		sm := syncMsg{
			ChannelId: p.ChannelId,
			PostId:    p.Id,
			Post:      postSync,
			Users:     users,
			Reactions: reactions,
		}
		syncMessages = append(syncMessages, sm)
	}

	if len(syncMessages) == 0 {
		// everything was filtered out and nothing to send.
		return model.RemoteClusterMsg{}, nil
	}

	json, err := json.Marshal(syncMessages)
	if err != nil {
		return model.RemoteClusterMsg{}, err
	}

	msg := model.NewRemoteClusterMsg(TopicSync, json)
	return msg, nil
}

// usersForPost provides a list of Users associated with the post that need to be synchronized.
func (scs *Service) usersForPost(post *model.Post, reactions []*model.Reaction, rc *model.RemoteCluster) []*model.User {
	userIds := make(map[string]struct{}) // avoid duplicates

	if post != nil {
		userIds[post.UserId] = struct{}{}
	}

	for _, r := range reactions {
		userIds[r.UserId] = struct{}{}
	}

	// TODO: extract @mentions to local users and sync those as well?

	users := make([]*model.User, 0)

	for id := range userIds {
		user, err := scs.server.GetStore().User().Get(id)
		if err == nil {
			if sync, err2 := scs.shouldUserSync(user, rc); err2 != nil {
				scs.server.GetLogger().Log(mlog.LvlSharedChannelServiceError, "Could not find user for post",
					mlog.String("user_id", id),
					mlog.Err(err2))
				continue
			} else if sync {
				users = append(users, sanitizeUserForSync(user))
			}
		} else {
			scs.server.GetLogger().Log(mlog.LvlSharedChannelServiceError, "Error checking if user should sync",
				mlog.String("user_id", id),
				mlog.Err(err))
		}
	}
	return users
}

func sanitizeUserForSync(user *model.User) *model.User {
	user.Password = model.NewId()
	user.AuthData = nil
	user.AuthService = ""
	user.Roles = "system_user"
	user.AllowMarketing = false
	user.Props = model.StringMap{}
	user.NotifyProps = model.StringMap{}
	user.LastPasswordUpdate = 0
	user.LastPictureUpdate = 0
	user.FailedAttempts = 0
	user.MfaActive = false
	user.MfaSecret = ""

	return user
}

// shouldUserSync determines if a user needs to be synchronized.
// User should be synchronized if it has no entry in the SharedChannelUsers table,
// or there is an entry but the LastSyncAt is less than user.UpdateAt
func (scs *Service) shouldUserSync(user *model.User, rc *model.RemoteCluster) (bool, error) {
	// don't sync users with the remote they originated from.
	if user.RemoteId != nil && *user.RemoteId == rc.RemoteId {
		return false, nil
	}

	scu, err := scs.server.GetStore().SharedChannel().GetUser(user.Id, rc.RemoteId)
	if err != nil {
		if _, ok := err.(errNotFound); !ok {
			return false, err
		}

		// user not in the SharedChannelUsers table, so we must add them.
		scu = &model.SharedChannelUser{
			UserId:          user.Id,
			RemoteClusterId: rc.RemoteId,
		}
		if _, err = scs.server.GetStore().SharedChannel().SaveUser(scu); err != nil {
			scs.server.GetLogger().Log(mlog.LvlSharedChannelServiceError, "Error adding user to shared channel users",
				mlog.String("remote_id", rc.RemoteId),
				mlog.String("user_id", user.Id))
		}
	} else if scu.LastSyncAt < user.UpdateAt {
		// Users.UpdateAt is updated quite frequently and for reasons we don't care about for sync.
		//
		// For now we'll mitigate that by not synchronizing users more than once per hour; alternatively we can start
		// storing the fields we care about in SharedChannelUsers and check if any changed, but in some use cases
		// the SharedChannelUsers table gets very large.
		//
		// Another option is to trigger sync on User update.
		//
		// Fields we care about for sync: username, email, nickname, firstname, lastname, position.
		if time.Since(model.GetTimeForMillis(scu.LastSyncAt)) < time.Hour {
			return false, nil
		}
	}
	return true, nil
}
