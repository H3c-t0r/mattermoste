// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"github.com/mattermost/mattermost-server/v6/model/sort"
	"net/http"
	"regexp"
)

// const (
// 	MaxAddMembersBatch    = 256
// 	MaximumBulkImportSize = 10 * 1024 * 1024
// 	groupIDsParamPattern  = "[^a-zA-Z0-9,]*"
// )

// var groupIDsQueryParamRegex *regexp.Regexp

func init() {
	groupIDsQueryParamRegex = regexp.MustCompile(groupIDsParamPattern)
}

func (api *API) InitHashtag() {

	api.BaseRoutes.HashTag.Handle("", api.APISessionRequired(api.suggestHashTag)).Methods("GET")
	api.BaseRoutes.HashTags.Handle("", api.APISessionRequired(api.getHashTags)).Methods("GET")
}

func (api *API) getHashTags(c *Context, w http.ResponseWriter, r *http.Request) {
	querySort := r.URL.Query().Get("sort")

	var sortToUse sort.Sort
	if querySort == "messages[asc]" {
		sortToUse = sort.Asc
	} else if querySort == "messages[desc]" {
		sortToUse = sort.Desc
	} else {
		sortToUse = ""
	}

	if sortToUse != "" {
		hashtags, _ := c.App.Srv().GetStore().Hashtag().GetMostCommon(sortToUse)
		response, _ := json.Marshal(hashtags)
		w.Write(response)
		return
	}

	hashtags, _ := c.App.Srv().GetStore().Hashtag().GetAll()
	response, _ := json.Marshal(hashtags)
	w.Write(response)
}

func (api *API) suggestHashTag(c *Context, w http.ResponseWriter, r *http.Request) {
	hashtags, _ := c.App.Srv().GetStore().Hashtag().SearchForUser(c.Params.HashTagQuery, c.AppContext.Session().UserId)

	response, _ := json.Marshal(hashtags)
	w.Write(response)
}
