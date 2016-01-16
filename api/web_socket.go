// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mattermost/platform/i18n"
	"github.com/mattermost/platform/model"
	"net/http"
)

func InitWebSocket(r *mux.Router) {
	l4g.Debug(T("Initializing web socket api routes"))
	r.Handle("/websocket", ApiUserRequired(connect)).Methods("GET")
	hub.Start(T)
}

func connect(c *Context, w http.ResponseWriter, r *http.Request) {
	T := i18n.GetTranslations(w, r)
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l4g.Error(T("websocket connect err: %v"), err)
		c.Err = model.NewAppError("connect", T("Failed to upgrade websocket connection"), "")
		return
	}

	wc := NewWebConn(ws, c.Session.TeamId, c.Session.UserId, c.Session.Id)
	hub.Register(wc)
	go wc.writePump()
	wc.readPump()
}
