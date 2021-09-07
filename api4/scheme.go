// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-server/v6/audit"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
)

func (api *API) InitScheme() {
	api.BaseRoutes.Schemes.Handle("", api.APISessionRequired(getSchemes)).Methods("GET")
	api.BaseRoutes.Schemes.Handle("", api.APISessionRequired(createScheme)).Methods("POST")
	api.BaseRoutes.Schemes.Handle("/{scheme_id:[A-Za-z0-9]+}", api.APISessionRequired(deleteScheme)).Methods("DELETE")
	api.BaseRoutes.Schemes.Handle("/{scheme_id:[A-Za-z0-9]+}", api.APISessionRequiredTrustRequester(getScheme)).Methods("GET")
	api.BaseRoutes.Schemes.Handle("/{scheme_id:[A-Za-z0-9]+}/patch", api.APISessionRequired(patchScheme)).Methods("PUT")
	api.BaseRoutes.Schemes.Handle("/{scheme_id:[A-Za-z0-9]+}/teams", api.APISessionRequiredTrustRequester(getTeamsForScheme)).Methods("GET")
	api.BaseRoutes.Schemes.Handle("/{scheme_id:[A-Za-z0-9]+}/channels", api.APISessionRequiredTrustRequester(getChannelsForScheme)).Methods("GET")
}

func createScheme(c *Context, w http.ResponseWriter, r *http.Request) {
	var scheme model.Scheme
	if jsonErr := json.NewDecoder(r.Body).Decode(&scheme); jsonErr != nil {
		c.SetInvalidParam("scheme")
		return
	}

	auditRec := c.MakeAuditRecord("createScheme", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("scheme", scheme)

	if c.App.Srv().License() == nil || !*c.App.Srv().License().Features.CustomPermissionsSchemes {
		c.Err = model.NewAppError("Api4.CreateScheme", "api.scheme.create_scheme.license.error", nil, "", http.StatusNotImplemented)
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleWriteUserManagementPermissions) {
		c.SetPermissionError(model.PermissionSysconsoleWriteUserManagementPermissions)
		return
	}

	returnedScheme, err := c.App.CreateScheme(&scheme)
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	auditRec.AddMeta("scheme", returnedScheme) // overwrite meta

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(returnedScheme); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getScheme(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireSchemeId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadUserManagementPermissions) {
		c.SetPermissionError(model.PermissionSysconsoleReadUserManagementPermissions)
		return
	}

	scheme, err := c.App.GetScheme(c.Params.SchemeId)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(scheme); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getSchemes(c *Context, w http.ResponseWriter, r *http.Request) {
	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadUserManagementPermissions) {
		c.SetPermissionError(model.PermissionSysconsoleReadUserManagementPermissions)
		return
	}

	scope := c.Params.Scope
	if scope != "" && scope != model.SchemeScopeTeam && scope != model.SchemeScopeChannel {
		c.SetInvalidParam("scope")
		return
	}

	schemes, err := c.App.GetSchemesPage(c.Params.Scope, c.Params.Page, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	}

	js, jsonErr := json.Marshal(schemes)
	if jsonErr != nil {
		c.Err = model.NewAppError("getSchemes", "api.marshal_error", nil, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func getTeamsForScheme(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireSchemeId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadUserManagementTeams) {
		c.SetPermissionError(model.PermissionSysconsoleReadUserManagementTeams)
		return
	}

	scheme, err := c.App.GetScheme(c.Params.SchemeId)
	if err != nil {
		c.Err = err
		return
	}

	if scheme.Scope != model.SchemeScopeTeam {
		c.Err = model.NewAppError("Api4.GetTeamsForScheme", "api.scheme.get_teams_for_scheme.scope.error", nil, "", http.StatusBadRequest)
		return
	}

	teams, err := c.App.GetTeamsForSchemePage(scheme, c.Params.Page, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	}

	js, jsonErr := json.Marshal(teams)
	if jsonErr != nil {
		c.Err = model.NewAppError("getTeamsForScheme", "api.marshal_error", nil, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func getChannelsForScheme(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireSchemeId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadUserManagementChannels) {
		c.SetPermissionError(model.PermissionSysconsoleReadUserManagementChannels)
		return
	}

	scheme, err := c.App.GetScheme(c.Params.SchemeId)
	if err != nil {
		c.Err = err
		return
	}

	if scheme.Scope != model.SchemeScopeChannel {
		c.Err = model.NewAppError("Api4.GetChannelsForScheme", "api.scheme.get_channels_for_scheme.scope.error", nil, "", http.StatusBadRequest)
		return
	}

	channels, err := c.App.GetChannelsForSchemePage(scheme, c.Params.Page, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(channels); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func patchScheme(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireSchemeId()
	if c.Err != nil {
		return
	}

	var patch model.SchemePatch
	if jsonErr := json.NewDecoder(r.Body).Decode(&patch); jsonErr != nil {
		c.SetInvalidParam("scheme")
		return
	}

	auditRec := c.MakeAuditRecord("patchScheme", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if c.App.Srv().License() == nil || !*c.App.Srv().License().Features.CustomPermissionsSchemes {
		c.Err = model.NewAppError("Api4.PatchScheme", "api.scheme.patch_scheme.license.error", nil, "", http.StatusNotImplemented)
		return
	}

	scheme, err := c.App.GetScheme(c.Params.SchemeId)
	if err != nil {
		c.Err = err
		return
	}
	auditRec.AddMeta("scheme", scheme)

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleWriteUserManagementPermissions) {
		c.SetPermissionError(model.PermissionSysconsoleWriteUserManagementPermissions)
		return
	}

	scheme, err = c.App.PatchScheme(scheme, &patch)
	if err != nil {
		c.Err = err
		return
	}
	auditRec.AddMeta("patch", scheme)

	auditRec.Success()
	c.LogAudit("")

	if err := json.NewEncoder(w).Encode(scheme); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func deleteScheme(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireSchemeId()
	if c.Err != nil {
		return
	}

	auditRec := c.MakeAuditRecord("deleteScheme", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if c.App.Srv().License() == nil || !*c.App.Srv().License().Features.CustomPermissionsSchemes {
		c.Err = model.NewAppError("Api4.DeleteScheme", "api.scheme.delete_scheme.license.error", nil, "", http.StatusNotImplemented)
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleWriteUserManagementPermissions) {
		c.SetPermissionError(model.PermissionSysconsoleWriteUserManagementPermissions)
		return
	}

	scheme, err := c.App.DeleteScheme(c.Params.SchemeId)
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	auditRec.AddMeta("scheme", scheme)

	ReturnStatusOK(w)
}
