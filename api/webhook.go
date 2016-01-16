// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/gorilla/mux"
	"github.com/mattermost/platform/i18n"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
	"net/http"
)

func InitWebhook(r *mux.Router) {
	l4g.Debug(T("Initializing webhook api routes"))

	sr := r.PathPrefix("/hooks").Subrouter()
	sr.Handle("/incoming/create", ApiUserRequired(createIncomingHook)).Methods("POST")
	sr.Handle("/incoming/delete", ApiUserRequired(deleteIncomingHook)).Methods("POST")
	sr.Handle("/incoming/list", ApiUserRequired(getIncomingHooks)).Methods("GET")

	sr.Handle("/outgoing/create", ApiUserRequired(createOutgoingHook)).Methods("POST")
	sr.Handle("/outgoing/regen_token", ApiUserRequired(regenOutgoingHookToken)).Methods("POST")
	sr.Handle("/outgoing/delete", ApiUserRequired(deleteOutgoingHook)).Methods("POST")
	sr.Handle("/outgoing/list", ApiUserRequired(getOutgoingHooks)).Methods("GET")
}

func createIncomingHook(c *Context, w http.ResponseWriter, r *http.Request) {
	T := i18n.GetTranslations(w, r)
	if !utils.Cfg.ServiceSettings.EnableIncomingWebhooks {
		c.Err = model.NewAppError("createIncomingHook", T("Incoming webhooks have been disabled by the system admin."), "")
		c.Err.StatusCode = http.StatusNotImplemented
		return
	}

	c.LogAudit(T("attempt"), T)

	hook := model.IncomingWebhookFromJson(r.Body)

	if hook == nil {
		c.SetInvalidParam("createIncomingHook", "webhook")
		return
	}

	cchan := Srv.Store.Channel().Get(hook.ChannelId, T)
	pchan := Srv.Store.Channel().CheckPermissionsTo(c.Session.TeamId, hook.ChannelId, c.Session.UserId, T)

	hook.UserId = c.Session.UserId
	hook.TeamId = c.Session.TeamId

	var channel *model.Channel
	if result := <-cchan; result.Err != nil {
		c.Err = result.Err
		return
	} else {
		channel = result.Data.(*model.Channel)
	}

	if !c.HasPermissionsToChannel(pchan, "createIncomingHook", T) {
		if channel.Type != model.CHANNEL_OPEN || channel.TeamId != c.Session.TeamId {
			c.LogAudit(T("fail - bad channel permissions"), T)
			return
		}
	}

	if result := <-Srv.Store.Webhook().SaveIncoming(hook, T); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		c.LogAudit(T("success"), T)
		rhook := result.Data.(*model.IncomingWebhook)
		w.Write([]byte(rhook.ToJson()))
	}
}

func deleteIncomingHook(c *Context, w http.ResponseWriter, r *http.Request) {
	T := i18n.GetTranslations(w, r)
	if !utils.Cfg.ServiceSettings.EnableIncomingWebhooks {
		c.Err = model.NewAppError("deleteIncomingHook", T("Incoming webhooks have been disabled by the system admin."), "")
		c.Err.StatusCode = http.StatusNotImplemented
		return
	}

	c.LogAudit(T("attempt"), T)

	props := model.MapFromJson(r.Body)

	id := props["id"]
	if len(id) == 0 {
		c.SetInvalidParam("deleteIncomingHook", "id")
		return
	}

	if result := <-Srv.Store.Webhook().GetIncoming(id, T); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		if c.Session.UserId != result.Data.(*model.IncomingWebhook).UserId && !c.IsTeamAdmin() {
			c.LogAudit(T("fail - inappropriate conditions"), T)
			c.Err = model.NewAppError("deleteIncomingHook", T("Inappropriate permissions to delete incoming webhook"), "user_id="+c.Session.UserId)
			return
		}
	}

	if err := (<-Srv.Store.Webhook().DeleteIncoming(id, model.GetMillis(), T)).Err; err != nil {
		c.Err = err
		return
	}

	c.LogAudit(T("success"), T)
	w.Write([]byte(model.MapToJson(props)))
}

func getIncomingHooks(c *Context, w http.ResponseWriter, r *http.Request) {
	T := i18n.GetTranslations(w, r)
	if !utils.Cfg.ServiceSettings.EnableIncomingWebhooks {
		c.Err = model.NewAppError("getIncomingHooks", T("Incoming webhooks have been disabled by the system admin."), "")
		c.Err.StatusCode = http.StatusNotImplemented
		return
	}

	if result := <-Srv.Store.Webhook().GetIncomingByUser(c.Session.UserId, T); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		hooks := result.Data.([]*model.IncomingWebhook)
		w.Write([]byte(model.IncomingWebhookListToJson(hooks)))
	}
}

func createOutgoingHook(c *Context, w http.ResponseWriter, r *http.Request) {
	T := i18n.GetTranslations(w, r)
	if !utils.Cfg.ServiceSettings.EnableOutgoingWebhooks {
		c.Err = model.NewAppError("createOutgoingHook", T("Outgoing webhooks have been disabled by the system admin."), "")
		c.Err.StatusCode = http.StatusNotImplemented
		return
	}

	c.LogAudit(T("attempt"), T)

	hook := model.OutgoingWebhookFromJson(r.Body)

	if hook == nil {
		c.SetInvalidParam("createOutgoingHook", "webhook")
		return
	}

	hook.CreatorId = c.Session.UserId
	hook.TeamId = c.Session.TeamId

	if len(hook.ChannelId) != 0 {
		cchan := Srv.Store.Channel().Get(hook.ChannelId, T)
		pchan := Srv.Store.Channel().CheckPermissionsTo(c.Session.TeamId, hook.ChannelId, c.Session.UserId, T)

		var channel *model.Channel
		if result := <-cchan; result.Err != nil {
			c.Err = result.Err
			return
		} else {
			channel = result.Data.(*model.Channel)
		}

		if channel.Type != model.CHANNEL_OPEN {
			c.LogAudit(T("fail - not open channel"), T)
		}

		if !c.HasPermissionsToChannel(pchan, "createOutgoingHook", T) {
			if channel.Type != model.CHANNEL_OPEN || channel.TeamId != c.Session.TeamId {
				c.LogAudit(T("fail - bad channel permissions"), T)
				return
			}
		}
	} else if len(hook.TriggerWords) == 0 {
		c.Err = model.NewAppError("createOutgoingHook", T("Either trigger_words or channel_id must be set"), "")
		return
	}

	if result := <-Srv.Store.Webhook().SaveOutgoing(hook, T); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		c.LogAudit(T("success"), T)
		rhook := result.Data.(*model.OutgoingWebhook)
		w.Write([]byte(rhook.ToJson()))
	}
}

func getOutgoingHooks(c *Context, w http.ResponseWriter, r *http.Request) {
	T := i18n.GetTranslations(w, r)
	if !utils.Cfg.ServiceSettings.EnableOutgoingWebhooks {
		c.Err = model.NewAppError("getOutgoingHooks", T("Outgoing webhooks have been disabled by the system admin."), "")
		c.Err.StatusCode = http.StatusNotImplemented
		return
	}

	if result := <-Srv.Store.Webhook().GetOutgoingByCreator(c.Session.UserId, T); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		hooks := result.Data.([]*model.OutgoingWebhook)
		w.Write([]byte(model.OutgoingWebhookListToJson(hooks)))
	}
}

func deleteOutgoingHook(c *Context, w http.ResponseWriter, r *http.Request) {
	T := i18n.GetTranslations(w, r)
	if !utils.Cfg.ServiceSettings.EnableIncomingWebhooks {
		c.Err = model.NewAppError("deleteOutgoingHook", T("Outgoing webhooks have been disabled by the system admin."), "")
		c.Err.StatusCode = http.StatusNotImplemented
		return
	}

	c.LogAudit(T("attempt"), T)

	props := model.MapFromJson(r.Body)

	id := props["id"]
	if len(id) == 0 {
		c.SetInvalidParam("deleteIncomingHook", "id")
		return
	}

	if result := <-Srv.Store.Webhook().GetOutgoing(id, T); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		if c.Session.UserId != result.Data.(*model.OutgoingWebhook).CreatorId && !c.IsTeamAdmin() {
			c.LogAudit(T("fail - inappropriate permissions"), T)
			c.Err = model.NewAppError("deleteOutgoingHook", T("Inappropriate permissions to delete outcoming webhook"), "user_id="+c.Session.UserId)
			return
		}
	}

	if err := (<-Srv.Store.Webhook().DeleteOutgoing(id, model.GetMillis(), T)).Err; err != nil {
		c.Err = err
		return
	}

	c.LogAudit(T("success"), T)
	w.Write([]byte(model.MapToJson(props)))
}

func regenOutgoingHookToken(c *Context, w http.ResponseWriter, r *http.Request) {
	T := i18n.GetTranslations(w, r)
	if !utils.Cfg.ServiceSettings.EnableIncomingWebhooks {
		c.Err = model.NewAppError("regenOutgoingHookToken", T("Outgoing webhooks have been disabled by the system admin."), "")
		c.Err.StatusCode = http.StatusNotImplemented
		return
	}

	c.LogAudit(T("attempt"), T)

	props := model.MapFromJson(r.Body)

	id := props["id"]
	if len(id) == 0 {
		c.SetInvalidParam("regenOutgoingHookToken", "id")
		return
	}

	var hook *model.OutgoingWebhook
	if result := <-Srv.Store.Webhook().GetOutgoing(id, T); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		hook = result.Data.(*model.OutgoingWebhook)

		if c.Session.UserId != hook.CreatorId && !c.IsTeamAdmin() {
			c.LogAudit(T("fail - inappropriate permissions"), T)
			c.Err = model.NewAppError("regenOutgoingHookToken", T("Inappropriate permissions to regenerate outcoming webhook token"), "user_id="+c.Session.UserId)
			return
		}
	}

	hook.Token = model.NewId()

	if result := <-Srv.Store.Webhook().UpdateOutgoing(hook, T); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		w.Write([]byte(result.Data.(*model.OutgoingWebhook).ToJson()))
	}
}
