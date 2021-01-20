// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package storetest

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"
)

func TestSharedChannelStore(t *testing.T, ss store.Store, s SqlStore) {
	t.Run("SaveSharedChannel", func(t *testing.T) { testSaveSharedChannel(t, ss) })
	t.Run("GetSharedChannel", func(t *testing.T) { testGetSharedChannel(t, ss) })
	t.Run("GetSharedChannels", func(t *testing.T) { testGetSharedChannels(t, ss) })
	t.Run("UpdateSharedChannel", func(t *testing.T) { testUpdateSharedChannel(t, ss) })
	t.Run("DeleteSharedChannel", func(t *testing.T) { testDeleteSharedChannel(t, ss) })

	t.Run("SaveSharedChannelRemote", func(t *testing.T) { testSaveSharedChannelRemote(t, ss) })
	t.Run("UpdateSharedChannelRemote", func(t *testing.T) { testUpdateSharedChannelRemote(t, ss) })
	t.Run("GetSharedChannelRemote", func(t *testing.T) { testGetSharedChannelRemote(t, ss) })
	t.Run("GetSharedChannelRemoteByIds", func(t *testing.T) { testGetSharedChannelRemoteByIds(t, ss) })
	t.Run("GetSharedChannelRemotes", func(t *testing.T) { testGetSharedChannelRemotes(t, ss) })
	t.Run("HasRemote", func(t *testing.T) { testHasRemote(t, ss) })
	t.Run("UpdateSharedChannelRemoteLastSyncAt", func(t *testing.T) { testUpdateSharedChannelRemoteLastSyncAt(t, ss) })
	t.Run("DeleteSharedChannelRemote", func(t *testing.T) { testDeleteSharedChannelRemote(t, ss) })
}

func testSaveSharedChannel(t *testing.T, ss store.Store) {
	t.Run("Save shared channel (home)", func(t *testing.T) {
		channel, err := createTestChannel(ss, "test_save")
		require.Nil(t, err)

		sc := &model.SharedChannel{
			ChannelId: channel.Id,
			TeamId:    channel.TeamId,
			CreatorId: model.NewId(),
			ShareName: "testshare",
			Home:      true,
		}

		scSaved, err := ss.SharedChannel().Save(sc)
		require.NoError(t, err, "couldn't save shared channel")

		require.Equal(t, sc.ChannelId, scSaved.ChannelId)
		require.Equal(t, sc.TeamId, scSaved.TeamId)
		require.Equal(t, sc.CreatorId, scSaved.CreatorId)

		// ensure channel's Shared flag is set
		channelMod, err := ss.Channel().Get(channel.Id, false)
		require.NoError(t, err)
		require.True(t, channelMod.IsShared())
	})

	t.Run("Save shared channel (remote)", func(t *testing.T) {
		channel, err := createTestChannel(ss, "test_save2")
		require.Nil(t, err)

		sc := &model.SharedChannel{
			ChannelId:       channel.Id,
			TeamId:          channel.TeamId,
			CreatorId:       model.NewId(),
			ShareName:       "testshare",
			RemoteClusterId: model.NewId(),
		}

		scSaved, err := ss.SharedChannel().Save(sc)
		require.Nil(t, err, "couldn't save shared channel", err)

		require.Equal(t, sc.ChannelId, scSaved.ChannelId)
		require.Equal(t, sc.TeamId, scSaved.TeamId)
		require.Equal(t, sc.CreatorId, scSaved.CreatorId)

		// ensure channel's Shared flag is set
		channelMod, err := ss.Channel().Get(channel.Id, false)
		require.NoError(t, err)
		require.True(t, channelMod.IsShared())
	})

	t.Run("Save invalid shared channel", func(t *testing.T) {
		sc := &model.SharedChannel{
			ChannelId: "",
			TeamId:    model.NewId(),
			CreatorId: model.NewId(),
			ShareName: "testshare",
			Home:      true,
		}

		_, err := ss.SharedChannel().Save(sc)
		require.NotNil(t, err, "should error saving invalid shared channel", err)
	})

	t.Run("Save with invalid channel id", func(t *testing.T) {
		sc := &model.SharedChannel{
			ChannelId:       model.NewId(),
			TeamId:          model.NewId(),
			CreatorId:       model.NewId(),
			ShareName:       "testshare",
			RemoteClusterId: model.NewId(),
		}

		_, err := ss.SharedChannel().Save(sc)
		require.Error(t, err, "expected error for invalid channel id")
	})
}

func testGetSharedChannel(t *testing.T, ss store.Store) {
	channel, err := createTestChannel(ss, "test_get")
	require.Nil(t, err)

	sc := &model.SharedChannel{
		ChannelId: channel.Id,
		TeamId:    channel.TeamId,
		CreatorId: model.NewId(),
		ShareName: "testshare",
		Home:      true,
	}

	scSaved, err := ss.SharedChannel().Save(sc)
	require.Nil(t, err, "couldn't save shared channel", err)

	t.Run("Get existing shared channel", func(t *testing.T) {
		sc, err := ss.SharedChannel().Get(scSaved.ChannelId)
		require.Nil(t, err, "couldn't get shared channel", err)

		require.Equal(t, sc.ChannelId, scSaved.ChannelId)
		require.Equal(t, sc.TeamId, scSaved.TeamId)
		require.Equal(t, sc.CreatorId, scSaved.CreatorId)
	})

	t.Run("Get non-existent shared channel", func(t *testing.T) {
		sc, err := ss.SharedChannel().Get(model.NewId())
		require.NotNil(t, err)
		require.Nil(t, sc)
	})
}

func testGetSharedChannels(t *testing.T, ss store.Store) {
	clearSharedChannels(ss)

	creator := model.NewId()
	team1 := model.NewId()
	team2 := model.NewId()
	rid := model.NewId()

	data := []model.SharedChannel{
		{CreatorId: creator, TeamId: team1, ShareName: "test1", Home: true},
		{CreatorId: creator, TeamId: team1, ShareName: "test2", Home: false, RemoteClusterId: rid},
		{CreatorId: creator, TeamId: team1, ShareName: "test3", Home: false, RemoteClusterId: rid},
		{CreatorId: creator, TeamId: team1, ShareName: "test4", Home: true},
		{CreatorId: creator, TeamId: team2, ShareName: "test5", Home: true},
		{CreatorId: creator, TeamId: team2, ShareName: "test6", Home: false, RemoteClusterId: rid},
		{CreatorId: creator, TeamId: team2, ShareName: "test7", Home: false, RemoteClusterId: rid},
		{CreatorId: creator, TeamId: team2, ShareName: "test8", Home: true},
		{CreatorId: creator, TeamId: team2, ShareName: "test9", Home: true},
	}

	for i, sc := range data {
		channel, err := createTestChannel(ss, "test_get2_"+strconv.Itoa(i))
		require.Nil(t, err)

		sc.ChannelId = channel.Id

		_, err = ss.SharedChannel().Save(&sc)
		require.Nil(t, err, "error saving shared channel")
	}

	t.Run("Get shared channels home only", func(t *testing.T) {
		opts := store.SharedChannelFilterOpts{
			ExcludeRemote: true,
			CreatorId:     creator,
		}

		count, err := ss.SharedChannel().GetAllCount(opts)
		require.Nil(t, err, "error getting shared channels count")

		home, err := ss.SharedChannel().GetAll(0, 100, opts)
		require.Nil(t, err, "error getting shared channels")

		require.Equal(t, int(count), len(home))
		require.Len(t, home, 5, "should be 5 home channels")
		for _, sc := range home {
			require.True(t, sc.Home, "should be home channel")
		}
	})

	t.Run("Get shared channels remote only", func(t *testing.T) {
		opts := store.SharedChannelFilterOpts{
			ExcludeHome: true,
		}

		count, err := ss.SharedChannel().GetAllCount(opts)
		require.Nil(t, err, "error getting shared channels count")

		remotes, err := ss.SharedChannel().GetAll(0, 100, opts)
		require.Nil(t, err, "error getting shared channels")

		require.Equal(t, int(count), len(remotes))
		require.Len(t, remotes, 4, "should be 4 remote channels")
		for _, sc := range remotes {
			require.False(t, sc.Home, "should be remote channel")
		}
	})

	t.Run("Get shared channels bad opts", func(t *testing.T) {
		opts := store.SharedChannelFilterOpts{
			ExcludeHome:   true,
			ExcludeRemote: true,
		}
		_, err := ss.SharedChannel().GetAll(0, 100, opts)
		require.NotNil(t, err, "error expected")
	})

	t.Run("Get shared channels by team", func(t *testing.T) {
		opts := store.SharedChannelFilterOpts{
			TeamId: team1,
		}

		count, err := ss.SharedChannel().GetAllCount(opts)
		require.Nil(t, err, "error getting shared channels count")

		remotes, err := ss.SharedChannel().GetAll(0, 100, opts)
		require.Nil(t, err, "error getting shared channels")

		require.Equal(t, int(count), len(remotes))
		require.Len(t, remotes, 4, "should be 4 matching channels")
		for _, sc := range remotes {
			require.Equal(t, team1, sc.TeamId)
		}
	})

	t.Run("Get shared channels invalid pagnation", func(t *testing.T) {
		opts := store.SharedChannelFilterOpts{
			TeamId: team1,
		}

		_, err := ss.SharedChannel().GetAll(-1, 100, opts)
		require.NotNil(t, err)

		_, err = ss.SharedChannel().GetAll(0, -100, opts)
		require.NotNil(t, err)
	})
}

func testUpdateSharedChannel(t *testing.T, ss store.Store) {
	channel, err := createTestChannel(ss, "test_update")
	require.Nil(t, err)

	sc := &model.SharedChannel{
		ChannelId: channel.Id,
		TeamId:    channel.TeamId,
		CreatorId: model.NewId(),
		ShareName: "testshare",
		Home:      true,
	}

	scSaved, err := ss.SharedChannel().Save(sc)
	require.Nil(t, err, "couldn't save shared channel", err)

	t.Run("Update existing shared channel", func(t *testing.T) {
		id := model.NewId()
		scMod := scSaved // copy struct (contains basic types only)
		scMod.ShareName = "newname"
		scMod.ShareDisplayName = "For testing"
		scMod.ShareHeader = "This is a header."
		scMod.RemoteClusterId = id

		scUpdated, err := ss.SharedChannel().Update(scMod)
		require.Nil(t, err, "couldn't update shared channel", err)

		require.Equal(t, "newname", scUpdated.ShareName)
		require.Equal(t, "For testing", scUpdated.ShareDisplayName)
		require.Equal(t, "This is a header.", scUpdated.ShareHeader)
		require.Equal(t, id, scUpdated.RemoteClusterId)
	})

	t.Run("Update non-existent shared channel", func(t *testing.T) {
		sc := &model.SharedChannel{
			ChannelId: model.NewId(),
			TeamId:    model.NewId(),
			CreatorId: model.NewId(),
			ShareName: "missingshare",
		}
		_, err := ss.SharedChannel().Update(sc)
		require.NotNil(t, err, "should error when updating non-existent shared channel", err)
	})
}

func testDeleteSharedChannel(t *testing.T, ss store.Store) {
	channel, err := createTestChannel(ss, "test_delete")
	require.Nil(t, err)

	sc := &model.SharedChannel{
		ChannelId:       channel.Id,
		TeamId:          channel.TeamId,
		CreatorId:       model.NewId(),
		ShareName:       "testshare",
		RemoteClusterId: model.NewId(),
	}

	_, err = ss.SharedChannel().Save(sc)
	require.Nil(t, err, "couldn't save shared channel", err)

	// add some remotes
	for i := 0; i < 10; i++ {
		remote := &model.SharedChannelRemote{
			ChannelId:       channel.Id,
			Description:     "remote_" + strconv.Itoa(i),
			CreatorId:       model.NewId(),
			RemoteClusterId: model.NewId(),
		}
		_, err := ss.SharedChannel().SaveRemote(remote)
		require.Nil(t, err, "couldn't add remote", err)
	}

	t.Run("Delete existing shared channel", func(t *testing.T) {
		deleted, err := ss.SharedChannel().Delete(channel.Id)
		require.Nil(t, err, "delete existing shared channel should not error", err)
		require.True(t, deleted, "expected true from delete shared channel")

		sc, err := ss.SharedChannel().Get(channel.Id)
		require.NotNil(t, err)
		require.Nil(t, sc)

		// make sure the remotes were deleted.
		remotes, err := ss.SharedChannel().GetRemotes(channel.Id)
		require.Nil(t, err)
		require.Len(t, remotes, 0, "expected empty remotes list")

		// ensure channel's Shared flag is unset
		channelMod, err := ss.Channel().Get(channel.Id, false)
		require.NoError(t, err)
		require.False(t, channelMod.IsShared())
	})

	t.Run("Delete non-existent shared channel", func(t *testing.T) {
		deleted, err := ss.SharedChannel().Delete(model.NewId())
		require.Nil(t, err, "delete non-existent shared channel should not error", err)
		require.False(t, deleted, "expected false from delete shared channel")
	})
}

func testSaveSharedChannelRemote(t *testing.T, ss store.Store) {
	t.Run("Save shared channel remote", func(t *testing.T) {
		channel, err := createTestChannel(ss, "test_save_remote")
		require.Nil(t, err)

		remote := &model.SharedChannelRemote{
			ChannelId:       channel.Id,
			Description:     "test_remote",
			CreatorId:       model.NewId(),
			RemoteClusterId: model.NewId(),
		}

		remoteSaved, err := ss.SharedChannel().SaveRemote(remote)
		require.Nil(t, err, "couldn't save shared channel remote", err)

		require.Equal(t, remote.ChannelId, remoteSaved.ChannelId)
		require.Equal(t, remote.CreatorId, remoteSaved.CreatorId)
	})

	t.Run("Save invalid shared channel remote", func(t *testing.T) {
		remote := &model.SharedChannelRemote{
			ChannelId:       "",
			Description:     "test_remote",
			CreatorId:       model.NewId(),
			RemoteClusterId: model.NewId(),
		}

		_, err := ss.SharedChannel().SaveRemote(remote)
		require.NotNil(t, err, "should error saving invalid remote", err)
	})

	t.Run("Save shared channel remote with invalid channel id", func(t *testing.T) {
		remote := &model.SharedChannelRemote{
			ChannelId:       model.NewId(),
			Description:     "test_remote",
			CreatorId:       model.NewId(),
			RemoteClusterId: model.NewId(),
		}

		_, err := ss.SharedChannel().SaveRemote(remote)
		require.Error(t, err, "expected error for invalid channel id")
	})
}

func testUpdateSharedChannelRemote(t *testing.T, ss store.Store) {
	t.Run("Update shared channel remote", func(t *testing.T) {
		channel, err := createTestChannel(ss, "test_update_remote")
		require.Nil(t, err)

		remote := &model.SharedChannelRemote{
			ChannelId:       channel.Id,
			Description:     "test_remote_update",
			CreatorId:       model.NewId(),
			RemoteClusterId: model.NewId(),
		}

		remoteSaved, err := ss.SharedChannel().SaveRemote(remote)
		require.Nil(t, err, "couldn't save shared channel remote", err)

		remoteSaved.IsInviteAccepted = true
		remoteSaved.IsInviteConfirmed = true
		remoteSaved.Description = "new_desc"

		remoteUpdated, err := ss.SharedChannel().UpdateRemote(remoteSaved)
		require.Nil(t, err, "couldn't update shared channel remote", err)

		require.Equal(t, true, remoteUpdated.IsInviteAccepted)
		require.Equal(t, true, remoteUpdated.IsInviteConfirmed)
		require.Equal(t, "new_desc", remoteUpdated.Description)
	})

	t.Run("Update invalid shared channel remote", func(t *testing.T) {
		remote := &model.SharedChannelRemote{
			ChannelId:       "",
			Description:     "test_remote",
			CreatorId:       model.NewId(),
			RemoteClusterId: model.NewId(),
		}

		_, err := ss.SharedChannel().UpdateRemote(remote)
		require.NotNil(t, err, "should error updating invalid remote", err)
	})

	t.Run("Update shared channel remote with invalid channel id", func(t *testing.T) {
		remote := &model.SharedChannelRemote{
			ChannelId:       model.NewId(),
			Description:     "test_remote",
			CreatorId:       model.NewId(),
			RemoteClusterId: model.NewId(),
		}

		_, err := ss.SharedChannel().UpdateRemote(remote)
		require.Error(t, err, "expected error for invalid channel id")
	})
}

func testGetSharedChannelRemote(t *testing.T, ss store.Store) {
	channel, err := createTestChannel(ss, "test_remote_get")
	require.Nil(t, err)

	remote := &model.SharedChannelRemote{
		ChannelId:       channel.Id,
		Description:     "test_remote",
		CreatorId:       model.NewId(),
		RemoteClusterId: model.NewId(),
	}

	remoteSaved, err := ss.SharedChannel().SaveRemote(remote)
	require.Nil(t, err, "couldn't save remote", err)

	t.Run("Get existing shared channel remote", func(t *testing.T) {
		r, err := ss.SharedChannel().GetRemote(remoteSaved.Id)
		require.Nil(t, err, "could not get shared channel remote", err)

		require.Equal(t, remoteSaved.Id, r.Id)
		require.Equal(t, remoteSaved.ChannelId, r.ChannelId)
		require.Equal(t, remoteSaved.Description, r.Description)
		require.Equal(t, remoteSaved.CreatorId, r.CreatorId)
		require.Equal(t, remoteSaved.RemoteClusterId, r.RemoteClusterId)
	})

	t.Run("Get non-existent shared channel remote", func(t *testing.T) {
		r, err := ss.SharedChannel().GetRemote(model.NewId())
		require.NotNil(t, err)
		require.Nil(t, r)
	})
}

func testGetSharedChannelRemoteByIds(t *testing.T, ss store.Store) {
	channel, err := createTestChannel(ss, "test_remote_get_by_ids")
	require.Nil(t, err)

	remote := &model.SharedChannelRemote{
		ChannelId:       channel.Id,
		Description:     "test_remote_by_ids",
		CreatorId:       model.NewId(),
		RemoteClusterId: model.NewId(),
	}

	remoteSaved, err := ss.SharedChannel().SaveRemote(remote)
	require.Nil(t, err, "could not save remote", err)

	t.Run("Get existing shared channel remote by ids", func(t *testing.T) {
		r, err := ss.SharedChannel().GetRemoteByIds(remoteSaved.ChannelId, remoteSaved.RemoteClusterId)
		require.Nil(t, err, "couldn't get shared channel remote by ids", err)

		require.Equal(t, remoteSaved.Id, r.Id)
		require.Equal(t, remoteSaved.ChannelId, r.ChannelId)
		require.Equal(t, remoteSaved.Description, r.Description)
		require.Equal(t, remoteSaved.CreatorId, r.CreatorId)
		require.Equal(t, remoteSaved.RemoteClusterId, r.RemoteClusterId)
	})

	t.Run("Get non-existent shared channel remote by ids", func(t *testing.T) {
		r, err := ss.SharedChannel().GetRemoteByIds(model.NewId(), model.NewId())
		require.NotNil(t, err)
		require.Nil(t, r)
	})
}

func testGetSharedChannelRemotes(t *testing.T, ss store.Store) {
	channel, err := createTestChannel(ss, "test_remotes_get2")
	require.Nil(t, err)

	creator := model.NewId()

	data := []model.SharedChannelRemote{
		{ChannelId: channel.Id, CreatorId: creator, Description: "r1", RemoteClusterId: model.NewId()},
		{ChannelId: channel.Id, CreatorId: creator, Description: "r2", RemoteClusterId: model.NewId()},
		{ChannelId: channel.Id, CreatorId: creator, Description: "r3", RemoteClusterId: model.NewId()},
		{ChannelId: "", CreatorId: creator, Description: "r4", RemoteClusterId: model.NewId()},
		{ChannelId: "", CreatorId: creator, Description: "r5", RemoteClusterId: model.NewId()},
		{ChannelId: "", CreatorId: creator, Description: "r6", RemoteClusterId: model.NewId()},
	}

	for i, r := range data {
		if r.ChannelId == "" {
			c, err := createTestChannel(ss, "test_remotes_get2_"+strconv.Itoa(i))
			require.Nil(t, err)
			r.ChannelId = c.Id
		}
		_, err := ss.SharedChannel().SaveRemote(&r)
		require.Nil(t, err, "error saving shared channel remote")
	}

	t.Run("Get shared channel remotes by channel_id", func(t *testing.T) {
		remotes, err := ss.SharedChannel().GetRemotes(channel.Id)
		require.Nil(t, err, "should not error", err)
		require.Len(t, remotes, 3)
		for _, r := range remotes {
			require.Contains(t, []string{"r1", "r2", "r3"}, r.Description)
		}
	})

	t.Run("Get shared channel remotes by invalid channel_id", func(t *testing.T) {
		remotes, err := ss.SharedChannel().GetRemotes(model.NewId())
		require.Nil(t, err, "should not error", err)
		require.Len(t, remotes, 0)
	})
}

func testHasRemote(t *testing.T, ss store.Store) {
	channel, err := createTestChannel(ss, "test_remotes_get2")
	require.Nil(t, err)

	creator := model.NewId()
	data := []model.SharedChannelRemote{
		{ChannelId: channel.Id, CreatorId: creator, Description: "r1", RemoteClusterId: model.NewId()},
		{ChannelId: channel.Id, CreatorId: creator, Description: "r2", RemoteClusterId: model.NewId()},
	}

	for _, r := range data {
		_, err := ss.SharedChannel().SaveRemote(&r)
		require.Nil(t, err, "error saving shared channel remote")
	}

	t.Run("has channel", func(t *testing.T) {
		has, err := ss.SharedChannel().HasRemote(channel.Id)
		require.NoError(t, err)
		assert.True(t, has)
	})

	t.Run("does not have channel", func(t *testing.T) {
		has, err := ss.SharedChannel().HasRemote(model.NewId())
		require.NoError(t, err)
		assert.False(t, has)
	})
}

func testUpdateSharedChannelRemoteLastSyncAt(t *testing.T, ss store.Store) {
	channel, err := createTestChannel(ss, "test_remote_update_last_sync_at")
	require.NoError(t, err)

	remote := &model.SharedChannelRemote{
		ChannelId:       channel.Id,
		Description:     "test_remote",
		CreatorId:       model.NewId(),
		RemoteClusterId: model.NewId(),
	}

	remoteSaved, err := ss.SharedChannel().SaveRemote(remote)
	require.NoError(t, err, "couldn't save remote", err)

	future := model.GetMillis() + 3600000 // 1 hour in the future

	t.Run("Update LastSyncAt", func(t *testing.T) {
		err := ss.SharedChannel().UpdateRemoteLastSyncAt(remoteSaved.Id, future)
		require.Nil(t, err, "update LasTSyncAt should not error", err)

		r, err := ss.SharedChannel().GetRemote(remoteSaved.Id)
		require.NoError(t, err)
		require.Equal(t, future, r.LastSyncAt)
	})

	t.Run("Update LastSyncAt for non-existent shared channel remote", func(t *testing.T) {
		err := ss.SharedChannel().UpdateRemoteLastSyncAt(model.NewId(), future)
		require.Error(t, err, "update non-existent remote should error", err)
	})
}

func testDeleteSharedChannelRemote(t *testing.T, ss store.Store) {
	channel, err := createTestChannel(ss, "test_remote_delete")
	require.NoError(t, err)

	remote := &model.SharedChannelRemote{
		ChannelId:       channel.Id,
		Description:     "test_remote",
		CreatorId:       model.NewId(),
		RemoteClusterId: model.NewId(),
	}

	remoteSaved, err := ss.SharedChannel().SaveRemote(remote)
	require.Nil(t, err, "couldn't save remote", err)

	t.Run("Delete existing shared channel remote", func(t *testing.T) {
		deleted, err := ss.SharedChannel().DeleteRemote(remoteSaved.Id)
		require.Nil(t, err, "delete existing remote should not error", err)
		require.True(t, deleted, "expected true from delete remote")

		r, err := ss.SharedChannel().GetRemote(remoteSaved.Id)
		require.NotNil(t, err)
		require.Nil(t, r)
	})

	t.Run("Delete non-existent shared channel remote", func(t *testing.T) {
		deleted, err := ss.SharedChannel().DeleteRemote(model.NewId())
		require.Nil(t, err, "delete non-existent remote should not error", err)
		require.False(t, deleted, "expected false from delete remote")
	})
}

func createTestChannel(ss store.Store, name string) (*model.Channel, error) {
	channel := &model.Channel{
		TeamId:      model.NewId(),
		Type:        model.CHANNEL_OPEN,
		Name:        name,
		DisplayName: name + " display name",
		Header:      name + " header",
		Purpose:     name + "purpose",
		CreatorId:   model.NewId(),
	}
	return ss.Channel().Save(channel, 10000)
}

func clearSharedChannels(ss store.Store) error {
	opts := store.SharedChannelFilterOpts{}
	all, err := ss.SharedChannel().GetAll(0, 1000, opts)
	if err != nil {
		return err
	}

	for _, sc := range all {
		if _, err := ss.SharedChannel().Delete(sc.ChannelId); err != nil {
			return err
		}
	}
	return nil
}
