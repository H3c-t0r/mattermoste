// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package store

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/go-gorp/gorp"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

type SqlChannelStore struct {
	*SqlStore
}

func NewSqlChannelStore(sqlStore *SqlStore) ChannelStore {
	s := &SqlChannelStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Channel{}, "Channels").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(26)
		table.ColMap("TeamId").SetMaxSize(26)
		table.ColMap("Type").SetMaxSize(1)
		table.ColMap("DisplayName").SetMaxSize(64)
		table.ColMap("Name").SetMaxSize(64)
		table.SetUniqueTogether("Name", "TeamId")
		table.ColMap("Header").SetMaxSize(1024)
		table.ColMap("Purpose").SetMaxSize(128)
		table.ColMap("CreatorId").SetMaxSize(26)

		tablem := db.AddTableWithName(model.ChannelMember{}, "ChannelMembers").SetKeys(false, "ChannelId", "UserId")
		tablem.ColMap("ChannelId").SetMaxSize(26)
		tablem.ColMap("UserId").SetMaxSize(26)
		tablem.ColMap("Roles").SetMaxSize(64)
		tablem.ColMap("NotifyProps").SetMaxSize(2000)
	}

	return s
}

func (s SqlChannelStore) UpgradeSchemaIfNeeded(T goi18n.TranslateFunc) {

	// BEGIN REMOVE AFTER 1.1.0
	if s.CreateColumnIfNotExists("ChannelMembers", "NotifyProps", "varchar(2000)", "varchar(2000)", "{}", T) {
		// populate NotifyProps from existing NotifyLevel field

		// set default values
		_, err := s.GetMaster().Exec(
			`UPDATE
				ChannelMembers
			SET
				NotifyProps = CONCAT('{"desktop":"', CONCAT(NotifyLevel, '","mark_unread":"` + model.CHANNEL_MARK_UNREAD_ALL + `"}'))`)
		if err != nil {
			l4g.Error(T("Unable to set default values for ChannelMembers.NotifyProps"))
			l4g.Error(err.Error())
		}

		// assume channels with all notifications enabled are just using the default settings
		_, err = s.GetMaster().Exec(
			`UPDATE
				ChannelMembers
			SET
				NotifyProps = '{"desktop":"` + model.CHANNEL_NOTIFY_DEFAULT + `","mark_unread":"` + model.CHANNEL_MARK_UNREAD_ALL + `"}'
			WHERE
				NotifyLevel = '` + model.CHANNEL_NOTIFY_ALL + `'`)
		if err != nil {
			l4g.Error(T("Unable to set values for ChannelMembers.NotifyProps when members previously had notifyLevel=all"))
			l4g.Error(err.Error())
		}

		// set quiet mode channels to have no notifications and only mark the channel unread on mentions
		_, err = s.GetMaster().Exec(
			`UPDATE
				ChannelMembers
			SET
				NotifyProps = '{"desktop":"` + model.CHANNEL_NOTIFY_NONE + `","mark_unread":"` + model.CHANNEL_MARK_UNREAD_MENTION + `"}'
			WHERE
				NotifyLevel = 'quiet'`)
		if err != nil {
			l4g.Error(T("Unable to set values for ChannelMembers.NotifyProps when members previously had notifyLevel=quiet"))
			l4g.Error(err.Error())
		}

		s.RemoveColumnIfExists("ChannelMembers", "NotifyLevel", T)
	}

	// BEGIN REMOVE AFTER 1.2.0
	s.RenameColumnIfExists("Channels", "Description", "Header", "varchar(1024)", T)
	s.CreateColumnIfNotExists("Channels", "Purpose", "varchar(1024)", "varchar(1024)", "", T)
	// END REMOVE AFTER 1.2.0
}

func (s SqlChannelStore) CreateIndexesIfNotExists(T goi18n.TranslateFunc) {
	s.CreateIndexIfNotExists("idx_channels_team_id", "Channels", "TeamId", T)
	s.CreateIndexIfNotExists("idx_channels_name", "Channels", "Name", T)

	s.CreateIndexIfNotExists("idx_channelmembers_channel_id", "ChannelMembers", "ChannelId", T)
	s.CreateIndexIfNotExists("idx_channelmembers_user_id", "ChannelMembers", "UserId", T)
}

func (s SqlChannelStore) Save(channel *model.Channel, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult
		if channel.Type == model.CHANNEL_DIRECT {
			result.Err = model.NewAppError("SqlChannelStore.Save", T("Use SaveDirectChannel to create a direct channel"), "")
		} else {
			if transaction, err := s.GetMaster().Begin(); err != nil {
				result.Err = model.NewAppError("SqlChannelStore.Save", T("Unable to open transaction"), err.Error())
			} else {
				result = s.saveChannelT(transaction, channel, T)
				if result.Err != nil {
					transaction.Rollback()
				} else {
					if err := transaction.Commit(); err != nil {
						result.Err = model.NewAppError("SqlChannelStore.Save", T("Unable to commit transaction"), err.Error())
					}
				}
			}
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) SaveDirectChannel(directchannel *model.Channel, member1 *model.ChannelMember, member2 *model.ChannelMember, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult

		if directchannel.Type != model.CHANNEL_DIRECT {
			result.Err = model.NewAppError("SqlChannelStore.SaveDirectChannel", T("Not a direct channel attempted to be created with SaveDirectChannel"), "")
		} else {
			if transaction, err := s.GetMaster().Begin(); err != nil {
				result.Err = model.NewAppError("SqlChannelStore.SaveDirectChannel", T("Unable to open transaction"), err.Error())
			} else {
				channelResult := s.saveChannelT(transaction, directchannel, T)

				if channelResult.Err != nil {
					transaction.Rollback()
					result.Err = channelResult.Err
				} else {
					newChannel := channelResult.Data.(*model.Channel)
					// Members need new channel ID
					member1.ChannelId = newChannel.Id
					member2.ChannelId = newChannel.Id

					member1Result := s.saveMemberT(transaction, member1, newChannel, T)
					member2Result := s.saveMemberT(transaction, member2, newChannel, T)

					if member1Result.Err != nil || member2Result.Err != nil {
						transaction.Rollback()
						details := ""
						if member1Result.Err != nil {
							details += "Member1Err: " + member1Result.Err.Message
						}
						if member2Result.Err != nil {
							details += "Member2Err: " + member2Result.Err.Message
						}
						result.Err = model.NewAppError("SqlChannelStore.SaveDirectChannel", T("Unable to add direct channel members"), details)
					} else {
						if err := transaction.Commit(); err != nil {
							result.Err = model.NewAppError("SqlChannelStore.SaveDirectChannel", T("Unable to commit transaction"), err.Error())
						} else {
							result = channelResult
						}
					}
				}
			}
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) saveChannelT(transaction *gorp.Transaction, channel *model.Channel, T goi18n.TranslateFunc) StoreResult {
	result := StoreResult{}

	if len(channel.Id) > 0 {
		result.Err = model.NewAppError("SqlChannelStore.Save", T("Must call update for exisiting channel"), "id="+channel.Id)
		return result
	}

	channel.PreSave()
	if result.Err = channel.IsValid(T); result.Err != nil {
		return result
	}

	if channel.Type != model.CHANNEL_DIRECT {
		if count, err := transaction.SelectInt("SELECT COUNT(0) FROM Channels WHERE TeamId = :TeamId AND DeleteAt = 0 AND (Type = 'O' OR Type = 'P')", map[string]interface{}{"TeamId": channel.TeamId}); err != nil {
			result.Err = model.NewAppError("SqlChannelStore.Save", T("Failed to get current channel count"), "teamId="+channel.TeamId+", "+err.Error())
			return result
		} else if count > 1000 {
			result.Err = model.NewAppError("SqlChannelStore.Save", T("You've reached the limit of the number of allowed channels."), "teamId="+channel.TeamId)
			return result
		}
	}

	if err := transaction.Insert(channel); err != nil {
		if IsUniqueConstraintError(err.Error(), "Name", "channels_name_teamid_key") {
			dupChannel := model.Channel{}
			s.GetMaster().SelectOne(&dupChannel, "SELECT * FROM Channels WHERE TeamId = :TeamId AND Name = :Name AND DeleteAt > 0", map[string]interface{}{"TeamId": channel.TeamId, "Name": channel.Name})
			if dupChannel.DeleteAt > 0 {
				result.Err = model.NewAppError("SqlChannelStore.Update", T("A channel with that URL was previously created"), "id="+channel.Id+", "+err.Error())
			} else {
				result.Err = model.NewAppError("SqlChannelStore.Update", T("A channel with that URL already exists"), "id="+channel.Id+", "+err.Error())
			}
		} else {
			result.Err = model.NewAppError("SqlChannelStore.Save", T("We couldn't save the channel"), "id="+channel.Id+", "+err.Error())
		}
	} else {
		result.Data = channel
	}

	return result
}

func (s SqlChannelStore) Update(channel *model.Channel, T goi18n.TranslateFunc) StoreChannel {

	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		channel.PreUpdate()

		if result.Err = channel.IsValid(T); result.Err != nil {
			storeChannel <- result
			close(storeChannel)
			return
		}

		if count, err := s.GetMaster().Update(channel); err != nil {
			if IsUniqueConstraintError(err.Error(), "Name", "channels_name_teamid_key") {
				dupChannel := model.Channel{}
				s.GetReplica().SelectOne(&dupChannel, "SELECT * FROM Channels WHERE TeamId = :TeamId AND Name= :Name AND DeleteAt > 0", map[string]interface{}{"TeamId": channel.TeamId, "Name": channel.Name})
				if dupChannel.DeleteAt > 0 {
					result.Err = model.NewAppError("SqlChannelStore.Update", T("A channel with that handle was previously created"), "id="+channel.Id+", "+err.Error())
				} else {
					result.Err = model.NewAppError("SqlChannelStore.Update", T("A channel with that handle already exists"), "id="+channel.Id+", "+err.Error())
				}
			} else {
				result.Err = model.NewAppError("SqlChannelStore.Update", T("We encounted an error updating the channel"), "id="+channel.Id+", "+err.Error())
			}
		} else if count != 1 {
			result.Err = model.NewAppError("SqlChannelStore.Update", T("We couldn't update the channel"), "id="+channel.Id)
		} else {
			result.Data = channel
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) extraUpdated(channel *model.Channel, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		channel.ExtraUpdated()

		_, err := s.GetMaster().Exec(
			`UPDATE
				Channels
			SET
				ExtraUpdateAt = :Time
			WHERE
				Id = :Id`,
			map[string]interface{}{"Id": channel.Id, "Time": channel.ExtraUpdateAt})

		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.extraUpdated", T("Problem updating members last updated time"), "id="+channel.Id+", "+err.Error())
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) Get(id string, T goi18n.TranslateFunc) StoreChannel {
	return s.get(id, false, T)
}

func (s SqlChannelStore) GetFromMaster(id string, T goi18n.TranslateFunc) StoreChannel {
	return s.get(id, true, T)
}

func (s SqlChannelStore) get(id string, master bool, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var db *gorp.DbMap
		if master {
			db = s.GetMaster()
		} else {
			db = s.GetReplica()
		}

		if obj, err := db.Get(model.Channel{}, id); err != nil {
			result.Err = model.NewAppError("SqlChannelStore.Get", T("We encountered an error finding the channel"), "id="+id+", "+err.Error())
		} else if obj == nil {
			result.Err = model.NewAppError("SqlChannelStore.Get", T("We couldn't find the existing channel"), "id="+id)
		} else {
			result.Data = obj.(*model.Channel)
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) Delete(channelId string, time int64, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		_, err := s.GetMaster().Exec("Update Channels SET DeleteAt = :Time, UpdateAt = :Time WHERE Id = :ChannelId", map[string]interface{}{"Time": time, "ChannelId": channelId})
		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.Delete", T("We couldn't delete the channel"), "id="+channelId+", err="+err.Error())
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) PermanentDeleteByTeam(teamId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		if _, err := s.GetMaster().Exec("DELETE FROM Channels WHERE TeamId = :TeamId", map[string]interface{}{"TeamId": teamId}); err != nil {
			result.Err = model.NewAppError("SqlChannelStore.PermanentDeleteByTeam", T("We couldn't delete the channels"), "teamId="+teamId+", "+err.Error())
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

type channelWithMember struct {
	model.Channel
	model.ChannelMember
}

func (s SqlChannelStore) GetChannels(teamId string, userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var data []channelWithMember
		_, err := s.GetReplica().Select(&data, "SELECT * FROM Channels, ChannelMembers WHERE Id = ChannelId AND TeamId = :TeamId AND UserId = :UserId AND DeleteAt = 0 ORDER BY DisplayName", map[string]interface{}{"TeamId": teamId, "UserId": userId})

		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.GetChannels", T("We couldn't get the channels"), "teamId="+teamId+", userId="+userId+", err="+err.Error())
		} else {
			channels := &model.ChannelList{make([]*model.Channel, len(data)), make(map[string]*model.ChannelMember)}
			for i := range data {
				v := data[i]
				channels.Channels[i] = &v.Channel
				channels.Members[v.Channel.Id] = &v.ChannelMember
			}

			if len(channels.Channels) == 0 {
				result.Err = model.NewAppError("SqlChannelStore.GetChannels", T("No channels were found"), "teamId="+teamId+", userId="+userId)
			} else {
				result.Data = channels
			}
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) GetMoreChannels(teamId string, userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var data []*model.Channel
		_, err := s.GetReplica().Select(&data,
			`SELECT 
			    *
			FROM
			    Channels
			WHERE
			    TeamId = :TeamId1
					AND Type IN ('O')
					AND DeleteAt = 0
			        AND Id NOT IN (SELECT 
			            Channels.Id
			        FROM
			            Channels,
			            ChannelMembers
			        WHERE
			            Id = ChannelId
			                AND TeamId = :TeamId2
			                AND UserId = :UserId
			                AND DeleteAt = 0)
			ORDER BY DisplayName`,
			map[string]interface{}{"TeamId1": teamId, "TeamId2": teamId, "UserId": userId})

		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.GetMoreChannels", T("We couldn't get the channels"), "teamId="+teamId+", userId="+userId+", err="+err.Error())
		} else {
			result.Data = &model.ChannelList{data, make(map[string]*model.ChannelMember)}
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

type channelIdWithCountAndUpdateAt struct {
	Id            string
	TotalMsgCount int64
	UpdateAt      int64
}

func (s SqlChannelStore) GetChannelCounts(teamId string, userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var data []channelIdWithCountAndUpdateAt
		_, err := s.GetReplica().Select(&data, "SELECT Id, TotalMsgCount, UpdateAt FROM Channels WHERE Id IN (SELECT ChannelId FROM ChannelMembers WHERE UserId = :UserId) AND TeamId = :TeamId AND DeleteAt = 0 ORDER BY DisplayName", map[string]interface{}{"TeamId": teamId, "UserId": userId})

		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.GetChannelCounts", T("We couldn't get the channel counts"), "teamId="+teamId+", userId="+userId+", err="+err.Error())
		} else {
			counts := &model.ChannelCounts{Counts: make(map[string]int64), UpdateTimes: make(map[string]int64)}
			for i := range data {
				v := data[i]
				counts.Counts[v.Id] = v.TotalMsgCount
				counts.UpdateTimes[v.Id] = v.UpdateAt
			}

			result.Data = counts
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

type totalChannelsByType struct {
	Total int64
}

func (s SqlChannelStore) GetTotalChannelsByType(teamType string, teamId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var data []totalChannelsByType
		_, err := s.GetReplica().Select(&data, "SELECT COUNT(DISTINCT Id) as Total FROM Channels WHERE TeamId = :TeamId AND Type = :TeamType AND DeleteAt=0",
			map[string]interface{}{"TeamId": teamId, "TeamType": teamType})

		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.GetChannelCounts", T("We couldn't get the channel counts"), "teamId="+teamId+", err="+err.Error())
		} else {
			counts := &model.ChannelCounts{Total: data[0].Total}
			result.Data = counts
		}

		storeChannel <- result
		close(storeChannel)

	}()

	return storeChannel
}

func (s SqlChannelStore) GetByName(teamId string, name string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		channel := model.Channel{}

		if err := s.GetReplica().SelectOne(&channel, "SELECT * FROM Channels WHERE TeamId = :TeamId AND Name= :Name AND DeleteAt = 0", map[string]interface{}{"TeamId": teamId, "Name": name}); err != nil {
			result.Err = model.NewAppError("SqlChannelStore.GetByName", T("We couldn't find the existing channel"), "teamId="+teamId+", "+"name="+name+", "+err.Error())
		} else {
			result.Data = &channel
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) SaveMember(member *model.ChannelMember, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult
		// Grab the channel we are saving this member to
		if cr := <-s.GetFromMaster(member.ChannelId, T); cr.Err != nil {
			result.Err = cr.Err
		} else {
			channel := cr.Data.(*model.Channel)

			if transaction, err := s.GetMaster().Begin(); err != nil {
				result.Err = model.NewAppError("SqlChannelStore.SaveMember", T("Unable to open transaction"), err.Error())
			} else {
				result = s.saveMemberT(transaction, member, channel, T)
				if result.Err != nil {
					transaction.Rollback()
				} else {
					if err := transaction.Commit(); err != nil {
						result.Err = model.NewAppError("SqlChannelStore.SaveMember", T("Unable to commit transaction"), err.Error())
					}
					// If sucessfull record members have changed in channel
					if mu := <-s.extraUpdated(channel, T); mu.Err != nil {
						result.Err = mu.Err
					}
				}
			}
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) saveMemberT(transaction *gorp.Transaction, member *model.ChannelMember, channel *model.Channel, T goi18n.TranslateFunc) StoreResult {
	result := StoreResult{}

	member.PreSave()
	if result.Err = member.IsValid(T); result.Err != nil {
		return result
	}

	if err := transaction.Insert(member); err != nil {
		if IsUniqueConstraintError(err.Error(), "ChannelId", "channelmembers_pkey") {
			result.Err = model.NewAppError("SqlChannelStore.SaveMember", T("A channel member with that id already exists"), "channel_id="+member.ChannelId+", user_id="+member.UserId+", "+err.Error())
		} else {
			result.Err = model.NewAppError("SqlChannelStore.SaveMember", T("We couldn't save the channel member"), "channel_id="+member.ChannelId+", user_id="+member.UserId+", "+err.Error())
		}
	} else {
		result.Data = member
	}

	return result
}

func (s SqlChannelStore) UpdateMember(member *model.ChannelMember, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		member.PreUpdate()

		if result.Err = member.IsValid(T); result.Err != nil {
			storeChannel <- result
			close(storeChannel)
			return
		}

		if _, err := s.GetMaster().Update(member); err != nil {
			result.Err = model.NewAppError("SqlChannelStore.UpdateMember", T("We encounted an error updating the channel member"),
				"channel_id="+member.ChannelId+", "+"user_id="+member.UserId+", "+err.Error())
		} else {
			result.Data = member
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) GetMembers(channelId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var members []model.ChannelMember
		_, err := s.GetReplica().Select(&members, "SELECT * FROM ChannelMembers WHERE ChannelId = :ChannelId", map[string]interface{}{"ChannelId": channelId})
		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.GetMembers", T("We couldn't get the channel members"), "channel_id="+channelId+err.Error())
		} else {
			result.Data = members
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) GetMember(channelId string, userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var member model.ChannelMember
		err := s.GetReplica().SelectOne(&member, "SELECT * FROM ChannelMembers WHERE ChannelId = :ChannelId AND UserId = :UserId", map[string]interface{}{"ChannelId": channelId, "UserId": userId})
		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.GetMember", T("We couldn't get the channel member"), "channel_id="+channelId+"user_id="+userId+","+err.Error())
		} else {
			result.Data = member
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) GetMemberCount(channelId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		count, err := s.GetReplica().SelectInt("SELECT count(*) FROM ChannelMembers WHERE ChannelId = :ChannelId", map[string]interface{}{"ChannelId": channelId})
		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.GetMemberCount", T("We couldn't get the channel member count"), "channel_id="+channelId+", "+err.Error())
		} else {
			result.Data = count
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) GetExtraMembers(channelId string, limit int, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var members []model.ExtraMember
		var err error

		if limit != -1 {
			_, err = s.GetReplica().Select(&members, "SELECT Id, Nickname, Email, ChannelMembers.Roles, Username FROM ChannelMembers, Users WHERE ChannelMembers.UserId = Users.Id AND ChannelId = :ChannelId LIMIT :Limit", map[string]interface{}{"ChannelId": channelId, "Limit": limit})
		} else {
			_, err = s.GetReplica().Select(&members, "SELECT Id, Nickname, Email, ChannelMembers.Roles, Username FROM ChannelMembers, Users WHERE ChannelMembers.UserId = Users.Id AND ChannelId = :ChannelId", map[string]interface{}{"ChannelId": channelId})
		}

		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.GetExtraMembers", T("We couldn't get the extra info for channel members"), "channel_id="+channelId+", "+err.Error())
		} else {
			for i := range members {
				members[i].Sanitize(utils.Cfg.GetSanitizeOptions())
			}
			result.Data = members
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) RemoveMember(channelId string, userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		// Grab the channel we are saving this member to
		if cr := <-s.Get(channelId, T); cr.Err != nil {
			result.Err = cr.Err
		} else {
			channel := cr.Data.(*model.Channel)

			_, err := s.GetMaster().Exec("DELETE FROM ChannelMembers WHERE ChannelId = :ChannelId AND UserId = :UserId", map[string]interface{}{"ChannelId": channelId, "UserId": userId})
			if err != nil {
				result.Err = model.NewAppError("SqlChannelStore.RemoveMember", T("We couldn't remove the channel member"), "channel_id="+channelId+", user_id="+userId+", "+err.Error())
			} else {
				// If sucessfull record members have changed in channel
				if mu := <-s.extraUpdated(channel, T); mu.Err != nil {
					result.Err = mu.Err
				}
			}
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) PermanentDeleteMembersByUser(userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		if _, err := s.GetMaster().Exec("DELETE FROM ChannelMembers WHERE UserId = :UserId", map[string]interface{}{"UserId": userId}); err != nil {
			result.Err = model.NewAppError("SqlChannelStore.RemoveMember", T("We couldn't remove the channel member"), "user_id="+userId+", "+err.Error())
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) CheckPermissionsTo(teamId string, channelId string, userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		count, err := s.GetReplica().SelectInt(
			`SELECT
			    COUNT(0)
			FROM
			    Channels,
			    ChannelMembers
			WHERE
			    Channels.Id = ChannelMembers.ChannelId
			        AND Channels.TeamId = :TeamId
			        AND Channels.DeleteAt = 0
			        AND ChannelMembers.ChannelId = :ChannelId
			        AND ChannelMembers.UserId = :UserId`,
			map[string]interface{}{"TeamId": teamId, "ChannelId": channelId, "UserId": userId})
		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.CheckPermissionsTo", T("We couldn't check the permissions"), "channel_id="+channelId+", user_id="+userId+", "+err.Error())
		} else {
			result.Data = count
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) CheckPermissionsToByName(teamId string, channelName string, userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		channelId, err := s.GetReplica().SelectStr(
			`SELECT
			    Channels.Id
			FROM
			    Channels,
			    ChannelMembers
			WHERE
			    Channels.Id = ChannelMembers.ChannelId
			        AND Channels.TeamId = :TeamId
			        AND Channels.Name = :Name
			        AND Channels.DeleteAt = 0
			        AND ChannelMembers.UserId = :UserId`,
			map[string]interface{}{"TeamId": teamId, "Name": channelName, "UserId": userId})
		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.CheckPermissionsToByName", T("We couldn't check the permissions"), "channel_id="+channelName+", user_id="+userId+", "+err.Error())
		} else {
			result.Data = channelId
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) CheckOpenChannelPermissions(teamId string, channelId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		count, err := s.GetReplica().SelectInt(
			`SELECT
			    COUNT(0)
			FROM
			    Channels
			WHERE
			    Channels.Id = :ChannelId
			        AND Channels.TeamId = :TeamId
			        AND Channels.Type = :ChannelType`,
			map[string]interface{}{"ChannelId": channelId, "TeamId": teamId, "ChannelType": model.CHANNEL_OPEN})
		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.CheckOpenChannelPermissions", T("We couldn't check the permissions"), "channel_id="+channelId+", "+err.Error())
		} else {
			result.Data = count
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) UpdateLastViewedAt(channelId string, userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var query string

		if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_POSTGRES {
			query = `UPDATE
				ChannelMembers
			SET
			    MentionCount = 0,
			    MsgCount = Channels.TotalMsgCount,
			    LastViewedAt = Channels.LastPostAt,
			    LastUpdateAt = Channels.LastPostAt
			FROM
				Channels
			WHERE
			    Channels.Id = ChannelMembers.ChannelId
			        AND UserId = :UserId
			        AND ChannelId = :ChannelId`
		} else if utils.Cfg.SqlSettings.DriverName == model.DATABASE_DRIVER_MYSQL {
			query = `UPDATE
				ChannelMembers, Channels
			SET
			    ChannelMembers.MentionCount = 0,
			    ChannelMembers.MsgCount = Channels.TotalMsgCount,
			    ChannelMembers.LastViewedAt = Channels.LastPostAt,
			    ChannelMembers.LastUpdateAt = Channels.LastPostAt
			WHERE
			    Channels.Id = ChannelMembers.ChannelId
			        AND UserId = :UserId
			        AND ChannelId = :ChannelId`
		}

		_, err := s.GetMaster().Exec(query, map[string]interface{}{"ChannelId": channelId, "UserId": userId})
		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.UpdateLastViewedAt", T("We couldn't update the last viewed at time"), "channel_id="+channelId+", user_id="+userId+", "+err.Error())
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) IncrementMentionCount(channelId string, userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		_, err := s.GetMaster().Exec(
			`UPDATE
				ChannelMembers
			SET
				MentionCount = MentionCount + 1
			WHERE
				UserId = :UserId
					AND ChannelId = :ChannelId`,
			map[string]interface{}{"ChannelId": channelId, "UserId": userId})
		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.IncrementMentionCount", T("We couldn't increment the mention count"), "channel_id="+channelId+", user_id="+userId+", "+err.Error())
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) GetForExport(teamId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var data []*model.Channel
		_, err := s.GetReplica().Select(&data, "SELECT * FROM Channels WHERE TeamId = :TeamId AND DeleteAt = 0 AND Type = 'O'", map[string]interface{}{"TeamId": teamId})

		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.GetAllChannels", T("We couldn't get all the channels"), "teamId="+teamId+", err="+err.Error())
		} else {
			result.Data = data
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlChannelStore) AnalyticsTypeCount(teamId string, channelType string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		v, err := s.GetReplica().SelectInt(
			`SELECT 
			    COUNT(Id) AS Value
			FROM
			    Channels
			WHERE
			    TeamId = :TeamId
			        AND Type = :ChannelType`,
			map[string]interface{}{"TeamId": teamId, "ChannelType": channelType})
		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.AnalyticsTypeCount", T("We couldn't get channel type counts"), err.Error())
		} else {
			result.Data = v
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}
