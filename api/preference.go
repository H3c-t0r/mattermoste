// Copyright (c) 2015 Spinpunch, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	l4g "code.google.com/p/log4go"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mattermost/platform/model"
	"net/http"
)

func InitPreference(r *mux.Router) {
	l4g.Debug("Initializing preference api routes")

	sr := r.PathPrefix("/preferences").Subrouter()
	sr.Handle("/set", ApiAppHandler(setPreferences)).Methods("POST")
	sr.Handle("/{category:[A-Za-z0-9_]+}/{name:[A-Za-z0-9_]+}", ApiAppHandler(getPreferencesByName)).Methods("GET")
}

func setPreferences(c *Context, w http.ResponseWriter, r *http.Request) {
	var preferences []model.Preference

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&preferences); err != nil {
		c.Err = model.NewAppError("setPreferences", "Unable to decode preferences from request", err.Error())
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	// just attempt to save/update them one by one and abort if one fails
	// in the future, this could probably be done in a transaction, but that's unnecessary now
	for _, preference := range preferences {
		if c.Session.UserId != preference.UserId {
			c.Err = model.NewAppError("setPreferences", "Unable to set preferences for other user", "session.user_id="+c.Session.UserId+", preference.user_id="+preference.UserId)
			c.Err.StatusCode = http.StatusUnauthorized
			return
		}

		if result := <-Srv.Store.Preference().Save(&preference); result.Err != nil {
			if result = <-Srv.Store.Preference().Update(&preference); result.Err != nil {
				c.Err = result.Err
				return
			}
		}
	}

	w.Write([]byte("true"))
}

func getPreferencesByName(c *Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	category := params["category"]
	name := params["name"]

	if result := <-Srv.Store.Preference().GetByName(c.Session.UserId, category, name); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		w.Write([]byte(model.PreferenceListToJson(result.Data.([]*model.Preference))))
	}
}
