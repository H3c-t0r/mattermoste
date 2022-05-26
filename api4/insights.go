// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-server/v6/model"
)

func (api *API) InitInsights() {
	// Reactions
	api.BaseRoutes.InsightsForTeam.Handle("/reactions", api.APISessionRequired(getTopReactionsForTeamSince)).Methods("GET")
	api.BaseRoutes.InsightsForUser.Handle("/reactions", api.APISessionRequired(getTopReactionsForUserSince)).Methods("GET")

	// Channels
	api.BaseRoutes.InsightsForTeam.Handle("/channels", api.APISessionRequired(getTopChannelsForTeamSince)).Methods("GET")
	api.BaseRoutes.InsightsForUser.Handle("/channels", api.APISessionRequired(getTopChannelsForUserSince)).Methods("GET")

	// Threads
	api.BaseRoutes.InsightsForTeam.Handle("/threads", api.APISessionRequired(requireLicense(getTopThreadsForTeamSince))).Methods("GET")
	api.BaseRoutes.InsightsForUser.Handle("/threads", api.APISessionRequired(requireLicense(getTopThreadsForUserSince))).Methods("GET")
}

// Top Reactions

func getTopReactionsForTeamSince(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	team, err := c.App.GetTeam(c.Params.TeamId)
	if err != nil {
		c.Err = err
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
		c.SetPermissionError(model.PermissionViewTeam)
		return
	}

	startTime, err := model.GetStartUnixMilliForTimeRange(c.Params.TimeRange)
	if err != nil {
		c.Err = err
		return
	}

	topReactionList, err := c.App.GetTopReactionsForTeamSince(c.Params.TeamId, c.AppContext.Session().UserId, &model.InsightsOpts{
		StartUnixMilli: startTime,
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if err != nil {
		c.Err = err
		return
	}

	js, jsonErr := json.Marshal(topReactionList)
	if jsonErr != nil {
		c.Err = model.NewAppError("getTopReactionsForTeamSince", "api.marshal_error", nil, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func getTopReactionsForUserSince(c *Context, w http.ResponseWriter, r *http.Request) {
	c.Params.TeamId = r.URL.Query().Get("team_id")

	// TeamId is an optional parameter
	if c.Params.TeamId != "" {
		if !model.IsValidId(c.Params.TeamId) {
			c.SetInvalidURLParam("team_id")
			return
		}

		team, teamErr := c.App.GetTeam(c.Params.TeamId)
		if teamErr != nil {
			c.Err = teamErr
			return
		}

		if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
			c.SetPermissionError(model.PermissionViewTeam)
			return
		}
	}

	startTime, err := model.GetStartUnixMilliForTimeRange(c.Params.TimeRange)
	if err != nil {
		c.Err = err
		return
	}

	topReactionList, err := c.App.GetTopReactionsForUserSince(c.AppContext.Session().UserId, c.Params.TeamId, &model.InsightsOpts{
		StartUnixMilli: startTime,
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if err != nil {
		c.Err = err
		return
	}

	js, jsonErr := json.Marshal(topReactionList)
	if jsonErr != nil {
		c.Err = model.NewAppError("getTopReactionsForUserSince", "api.marshal_error", nil, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

// Top Channels

func getTopChannelsForTeamSince(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	team, err := c.App.GetTeam(c.Params.TeamId)
	if err != nil {
		c.Err = err
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
		c.SetPermissionError(model.PermissionViewTeam)
		return
	}

	startTime, err := model.GetStartUnixMilliForTimeRange(c.Params.TimeRange)
	if err != nil {
		c.Err = err
		return
	}

	topChannels, err := c.App.GetTopChannelsForTeamSince(c.Params.TeamId, c.AppContext.Session().UserId, &model.InsightsOpts{
		StartUnixMilli: startTime,
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if err != nil {
		c.Err = err
		return
	}

	js, jsonErr := json.Marshal(topChannels)
	if jsonErr != nil {
		c.Err = model.NewAppError("getTopChannelsForTeamSince", "api.marshal_error", nil, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func getTopChannelsForUserSince(c *Context, w http.ResponseWriter, r *http.Request) {
	c.Params.TeamId = r.URL.Query().Get("team_id")

	// TeamId is an optional parameter
	if c.Params.TeamId != "" {
		if !model.IsValidId(c.Params.TeamId) {
			c.SetInvalidURLParam("team_id")
			return
		}

		team, teamErr := c.App.GetTeam(c.Params.TeamId)
		if teamErr != nil {
			c.Err = teamErr
			return
		}

		if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
			c.SetPermissionError(model.PermissionViewTeam)
			return
		}
	}

	startTime, err := model.GetStartUnixMilliForTimeRange(c.Params.TimeRange)
	if err != nil {
		c.Err = err
		return
	}

	topChannels, err := c.App.GetTopChannelsForUserSince(c.AppContext.Session().UserId, c.Params.TeamId, &model.InsightsOpts{
		StartUnixMilli: startTime,
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})

	if err != nil {
		c.Err = err
		return
	}

	js, jsonErr := json.Marshal(topChannels)
	if jsonErr != nil {
		c.Err = model.NewAppError("getTopChannelsForUserSince", "api.marshal_error", nil, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

// Top Threads
func getTopThreadsForTeamSince(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	team, err := c.App.GetTeam(c.Params.TeamId)
	if err != nil {
		c.Err = err
		return
	}

	// license check
	lic := c.App.Srv().License()
	if lic.SkuShortName != model.LicenseShortSkuProfessional && lic.SkuShortName != model.LicenseShortSkuEnterprise {
		c.Err = model.NewAppError("", "api.insights.license_error", nil, "", http.StatusNotImplemented)
		return
	}

	// restrict guests and users with no access to team
	user, err := c.App.GetUser(c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) || user.IsGuest() {
		c.SetPermissionError(model.PermissionViewTeam)
		return
	}

	startTime, err := model.GetStartUnixMilliForTimeRange(c.Params.TimeRange)
	if err != nil {
		c.Err = err
		return
	}

	topThreads, err := c.App.GetTopThreadsForTeamSince(c.Params.TeamId, c.AppContext.Session().UserId, &model.InsightsOpts{
		StartUnixMilli: startTime,
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if err != nil {
		c.Err = err
		return
	}

	js, jsonErr := json.Marshal(topThreads)
	if jsonErr != nil {
		c.Err = model.NewAppError("getTopThreadsForTeamSince", "api.marshal_error", nil, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func getTopThreadsForUserSince(c *Context, w http.ResponseWriter, r *http.Request) {
	c.Params.TeamId = r.URL.Query().Get("team_id")

	// TeamId is an optional parameter
	if c.Params.TeamId != "" {
		if !model.IsValidId(c.Params.TeamId) {
			c.SetInvalidURLParam("team_id")
			return
		}

		team, teamErr := c.App.GetTeam(c.Params.TeamId)
		if teamErr != nil {
			c.Err = teamErr
			return
		}

		// license check
		lic := c.App.Srv().License()
		if lic.SkuShortName != model.LicenseShortSkuProfessional && lic.SkuShortName != model.LicenseShortSkuEnterprise {
			c.Err = model.NewAppError("", "api.insights.license_error", nil, "", http.StatusNotImplemented)
			return
		}

		// restrict guests and users with no access to team
		user, err := c.App.GetUser(c.AppContext.Session().UserId)
		if err != nil {
			c.Err = err
			return
		}

		if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) || user.IsGuest() {
			c.SetPermissionError(model.PermissionViewTeam)
			return
		}
	}

	startTime, err := model.GetStartUnixMilliForTimeRange(c.Params.TimeRange)
	if err != nil {
		c.Err = err
		return
	}

	topThreads, err := c.App.GetTopThreadsForUserSince(c.Params.TeamId, c.AppContext.Session().UserId, &model.InsightsOpts{
		StartUnixMilli: startTime,
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if err != nil {
		c.Err = err
		return
	}

	js, jsonErr := json.Marshal(topThreads)
	if jsonErr != nil {
		c.Err = model.NewAppError("getTopThreadsForUserSince", "api.marshal_error", nil, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}
