// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mattermost/mattermost-server/v5/audit"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/services/remotecluster"
)

func (api *API) InitRemoteCluster() {
	api.BaseRoutes.RemoteCluster.Handle("/ping", api.ApiHandler(remoteClusterPing)).Methods("POST")
	api.BaseRoutes.RemoteCluster.Handle("/msg", api.ApiHandler(remoteClusterAcceptMessage)).Methods("POST")
	api.BaseRoutes.RemoteCluster.Handle("/confirm_invite", api.ApiHandler(remoteClusterConfirmInvite)).Methods("POST")
}

func remoteClusterPing(c *Context, w http.ResponseWriter, r *http.Request) {
	// make sure remote cluster service is running.
	if _, appErr := c.App.GetRemoteClusterService(); appErr != nil {
		c.Err = appErr
		return
	}

	frame, appErr := model.RemoteClusterFrameFromJSON(r.Body)
	if appErr != nil {
		c.Err = appErr
		return
	}

	if appErr = frame.IsValid(); appErr != nil {
		c.Err = appErr
		return
	}

	auditRec := c.MakeAuditRecord("remoteClusterPing", audit.Fail)
	defer c.LogAuditRec(auditRec)

	rc, err := c.App.GetRemoteCluster(frame.RemoteId)
	if err != nil {
		c.SetInvalidRemoteClusterIdError(frame.RemoteId)
		return
	}
	auditRec.AddMeta("remoteCluster", rc)

	if rc.Token != frame.Token {
		c.SetInvalidRemoteClusterTokenError()
		return
	}

	ping, err := model.RemoteClusterPingFromRawJSON(frame.Msg.Payload)
	if err != nil {
		c.SetInvalidParam("msg.payload")
		return
	}
	ping.RecvAt = model.GetMillis()

	if metrics := c.App.Metrics(); metrics != nil {
		metrics.IncrementRemoteClusterMsgReceivedCounter(rc.RemoteId)
	}

	auditRec.AddMeta("SentAt", ping.SentAt)
	auditRec.AddMeta("RecvAt", ping.RecvAt)

	if err := c.App.SetRemoteClusterLastPingAt(rc.RemoteId); err != nil {
		auditRec.AddMeta("err", err)
		c.Err = err
		return
	}

	auditRec.Success()

	resp, _ := json.Marshal(ping)
	w.Write(resp)
}

func remoteClusterAcceptMessage(c *Context, w http.ResponseWriter, r *http.Request) {
	// make sure remote cluster service is running.
	service, appErr := c.App.GetRemoteClusterService()
	if appErr != nil {
		c.Err = appErr
		return
	}

	frame, appErr := model.RemoteClusterFrameFromJSON(r.Body)
	if appErr != nil {
		c.Err = appErr
		return
	}

	if appErr = frame.IsValid(); appErr != nil {
		c.Err = appErr
		return
	}

	auditRec := c.MakeAuditRecord("remoteClusterAcceptMessage", audit.Fail)
	defer c.LogAuditRec(auditRec)

	rc, err := c.App.GetRemoteCluster(frame.RemoteId)
	if err != nil {
		c.SetInvalidRemoteClusterIdError(frame.RemoteId)
		return
	}
	auditRec.AddMeta("remoteCluster", rc)

	if rc.Token != frame.Token {
		c.SetInvalidRemoteClusterTokenError()
		return
	}

	// pass message to Remote Cluster Service and write response
	resp := service.ReceiveIncomingMsg(rc, frame.Msg)

	b, errMarshall := json.Marshal(resp)
	if errMarshall != nil {
		c.Err = model.NewAppError("remoteClusterAcceptMessage", "api.marshal_error", nil, errMarshall.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func remoteClusterConfirmInvite(c *Context, w http.ResponseWriter, r *http.Request) {
	// make sure remote cluster service is running.
	if _, appErr := c.App.GetRemoteClusterService(); appErr != nil {
		c.Err = appErr
		return
	}

	frame, appErr := model.RemoteClusterFrameFromJSON(r.Body)
	if appErr != nil {
		c.Err = appErr
		return
	}

	if appErr = frame.IsValid(); appErr != nil {
		c.Err = appErr
		return
	}

	auditRec := c.MakeAuditRecord("remoteClusterAcceptInvite", audit.Fail)
	defer c.LogAuditRec(auditRec)

	rc, err := c.App.GetRemoteCluster(frame.RemoteId)
	if err != nil {
		c.SetInvalidRemoteClusterIdError(frame.RemoteId)
		return
	}
	auditRec.AddMeta("remoteCluster", rc)

	if time.Since(model.GetTimeForMillis(rc.CreateAt)) > remotecluster.InviteExpiresAfter {
		c.Err = model.NewAppError("remoteClusterAcceptMessage", "api.context.invitation_expired.error", nil, "", http.StatusBadRequest)
		return
	}

	if rc.Token != frame.Token {
		c.SetInvalidRemoteClusterTokenError()
		return
	}

	confirm, appErr := model.RemoteClusterInviteFromRawJSON(frame.Msg.Payload)
	if appErr != nil {
		c.Err = appErr
		return
	}

	rc.RemoteTeamId = confirm.RemoteTeamId
	rc.SiteURL = confirm.SiteURL
	rc.RemoteToken = confirm.Token

	if _, err := c.App.UpdateRemoteCluster(rc); err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	ReturnStatusOK(w)
}
