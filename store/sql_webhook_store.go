// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package store

import (
	"github.com/mattermost/platform/model"
	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

type SqlWebhookStore struct {
	*SqlStore
}

func NewSqlWebhookStore(sqlStore *SqlStore) WebhookStore {
	s := &SqlWebhookStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.IncomingWebhook{}, "IncomingWebhooks").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(26)
		table.ColMap("UserId").SetMaxSize(26)
		table.ColMap("ChannelId").SetMaxSize(26)
		table.ColMap("TeamId").SetMaxSize(26)

		tableo := db.AddTableWithName(model.OutgoingWebhook{}, "OutgoingWebhooks").SetKeys(false, "Id")
		tableo.ColMap("Id").SetMaxSize(26)
		tableo.ColMap("Token").SetMaxSize(26)
		tableo.ColMap("CreatorId").SetMaxSize(26)
		tableo.ColMap("ChannelId").SetMaxSize(26)
		tableo.ColMap("TeamId").SetMaxSize(26)
		tableo.ColMap("TriggerWords").SetMaxSize(1024)
		tableo.ColMap("CallbackURLs").SetMaxSize(1024)
	}

	return s
}

func (s SqlWebhookStore) UpgradeSchemaIfNeeded() {
}

func (s SqlWebhookStore) CreateIndexesIfNotExists(T goi18n.TranslateFunc) {
	s.CreateIndexIfNotExists("idx_incoming_webhook_user_id", "IncomingWebhooks", "UserId", T)
	s.CreateIndexIfNotExists("idx_incoming_webhook_team_id", "IncomingWebhooks", "TeamId", T)
	s.CreateIndexIfNotExists("idx_outgoing_webhook_team_id", "OutgoingWebhooks", "TeamId", T)
}

func (s SqlWebhookStore) SaveIncoming(webhook *model.IncomingWebhook, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		if len(webhook.Id) > 0 {
			result.Err = model.NewAppError("SqlWebhookStore.SaveIncoming",
				T("You cannot overwrite an existing IncomingWebhook"), "id="+webhook.Id)
			storeChannel <- result
			close(storeChannel)
			return
		}

		webhook.PreSave()
		if result.Err = webhook.IsValid(T); result.Err != nil {
			storeChannel <- result
			close(storeChannel)
			return
		}

		if err := s.GetMaster().Insert(webhook); err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.SaveIncoming", T("We couldn't save the IncomingWebhook"), "id="+webhook.Id+", "+err.Error())
		} else {
			result.Data = webhook
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlWebhookStore) GetIncoming(id string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var webhook model.IncomingWebhook

		if err := s.GetReplica().SelectOne(&webhook, "SELECT * FROM IncomingWebhooks WHERE Id = :Id AND DeleteAt = 0", map[string]interface{}{"Id": id}); err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.GetIncoming", T("We couldn't get the webhook"), "id="+id+", err="+err.Error())
		}

		result.Data = &webhook

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlWebhookStore) DeleteIncoming(webhookId string, time int64, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		_, err := s.GetMaster().Exec("Update IncomingWebhooks SET DeleteAt = :DeleteAt, UpdateAt = :UpdateAt WHERE Id = :Id", map[string]interface{}{"DeleteAt": time, "UpdateAt": time, "Id": webhookId})
		if err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.DeleteIncoming", T("We couldn't delete the webhook"), "id="+webhookId+", err="+err.Error())
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlWebhookStore) PermanentDeleteIncomingByUser(userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		_, err := s.GetMaster().Exec("DELETE FROM IncomingWebhooks WHERE UserId = :UserId", map[string]interface{}{"UserId": userId})
		if err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.DeleteIncomingByUser", T("We couldn't delete the webhook"), "id="+userId+", err="+err.Error())
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlWebhookStore) GetIncomingByUser(userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var webhooks []*model.IncomingWebhook

		if _, err := s.GetReplica().Select(&webhooks, "SELECT * FROM IncomingWebhooks WHERE UserId = :UserId AND DeleteAt = 0", map[string]interface{}{"UserId": userId}); err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.GetIncomingByUser", T("We couldn't get the webhook"), "userId="+userId+", err="+err.Error())
		}

		result.Data = webhooks

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlWebhookStore) GetIncomingByChannel(channelId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var webhooks []*model.IncomingWebhook

		if _, err := s.GetReplica().Select(&webhooks, "SELECT * FROM IncomingWebhooks WHERE ChannelId = :ChannelId AND DeleteAt = 0", map[string]interface{}{"ChannelId": channelId}); err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.GetIncomingByChannel", T("We couldn't get the webhooks"), "channelId="+channelId+", err="+err.Error())
		}

		result.Data = webhooks

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlWebhookStore) SaveOutgoing(webhook *model.OutgoingWebhook, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		if len(webhook.Id) > 0 {
			result.Err = model.NewAppError("SqlWebhookStore.SaveOutgoing",
				T("You cannot overwrite an existing OutgoingWebhook"), "id="+webhook.Id)
			storeChannel <- result
			close(storeChannel)
			return
		}

		webhook.PreSave()
		if result.Err = webhook.IsValid(T); result.Err != nil {
			storeChannel <- result
			close(storeChannel)
			return
		}

		if err := s.GetMaster().Insert(webhook); err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.SaveOutgoing", T("We couldn't save the OutgoingWebhook"), "id="+webhook.Id+", "+err.Error())
		} else {
			result.Data = webhook
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlWebhookStore) GetOutgoing(id string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var webhook model.OutgoingWebhook

		if err := s.GetReplica().SelectOne(&webhook, "SELECT * FROM OutgoingWebhooks WHERE Id = :Id AND DeleteAt = 0", map[string]interface{}{"Id": id}); err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.GetOutgoing", T("We couldn't get the webhook"), "id="+id+", err="+err.Error())
		}

		result.Data = &webhook

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlWebhookStore) GetOutgoingByCreator(userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var webhooks []*model.OutgoingWebhook

		if _, err := s.GetReplica().Select(&webhooks, "SELECT * FROM OutgoingWebhooks WHERE CreatorId = :UserId AND DeleteAt = 0", map[string]interface{}{"UserId": userId}); err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.GetOutgoingByCreator", T("We couldn't get the webhooks"), "userId="+userId+", err="+err.Error())
		}

		result.Data = webhooks

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlWebhookStore) GetOutgoingByChannel(channelId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var webhooks []*model.OutgoingWebhook

		if _, err := s.GetReplica().Select(&webhooks, "SELECT * FROM OutgoingWebhooks WHERE ChannelId = :ChannelId AND DeleteAt = 0", map[string]interface{}{"ChannelId": channelId}); err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.GetOutgoingByChannel", T("We couldn't get the webhooks"), "channelId="+channelId+", err="+err.Error())
		}

		result.Data = webhooks

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlWebhookStore) GetOutgoingByTeam(teamId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		var webhooks []*model.OutgoingWebhook

		if _, err := s.GetReplica().Select(&webhooks, "SELECT * FROM OutgoingWebhooks WHERE TeamId = :TeamId AND DeleteAt = 0", map[string]interface{}{"TeamId": teamId}); err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.GetOutgoingByTeam", T("We couldn't get the webhooks"), "teamId="+teamId+", err="+err.Error())
		}

		result.Data = webhooks

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlWebhookStore) DeleteOutgoing(webhookId string, time int64, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		_, err := s.GetMaster().Exec("Update OutgoingWebhooks SET DeleteAt = :DeleteAt, UpdateAt = :UpdateAt WHERE Id = :Id", map[string]interface{}{"DeleteAt": time, "UpdateAt": time, "Id": webhookId})
		if err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.DeleteOutgoing", T("We couldn't delete the webhook"), "id="+webhookId+", err="+err.Error())
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlWebhookStore) PermanentDeleteOutgoingByUser(userId string, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		_, err := s.GetMaster().Exec("DELETE FROM OutgoingWebhooks WHERE CreatorId = :UserId", map[string]interface{}{"UserId": userId})
		if err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.DeleteOutgoingByUser", T("We couldn't delete the webhook"), "id="+userId+", err="+err.Error())
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s SqlWebhookStore) UpdateOutgoing(hook *model.OutgoingWebhook, T goi18n.TranslateFunc) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := StoreResult{}

		hook.UpdateAt = model.GetMillis()

		if _, err := s.GetMaster().Update(hook); err != nil {
			result.Err = model.NewAppError("SqlWebhookStore.UpdateOutgoing", T("We couldn't update the webhook"), "id="+hook.Id+", "+err.Error())
		} else {
			result.Data = hook
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}
