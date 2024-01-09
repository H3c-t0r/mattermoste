// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/v8/channels/audit"
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

const (
	whereOutgoingOAuthConnection = "outgoingOAuthConnections"
)

func (api *API) InitOutgoingOAuthConnection() {
	api.BaseRoutes.OutgoingOAuthConnections.Handle("", api.APISessionRequired(listOutgoingOAuthConnections)).Methods("GET")
	api.BaseRoutes.OutgoingOAuthConnections.Handle("", api.APISessionRequired(createOutgoingOAuthConnection)).Methods("POST")
	api.BaseRoutes.OutgoingOAuthConnection.Handle("", api.APISessionRequired(getOutgoingOAuthConnection)).Methods("GET")
	api.BaseRoutes.OutgoingOAuthConnection.Handle("", api.APISessionRequired(updateOutgoingOAuthConnection)).Methods("PUT")
	api.BaseRoutes.OutgoingOAuthConnection.Handle("", api.APISessionRequired(deleteOutgoingOAuthConnection)).Methods("DELETE")
}

// checkOutgoingOAuthConnectionReadPermissions checks if the user has the permissions to read outgoing oauth connections.
// The user needs to have the permissions to manage outgoing webhooks and slash commands in order to read outgoing oauth connections
// so that they can use them in their outgoing webhooks and slash commands.
// This is made in this way so only users with the management permission can setup the outgoing oauth connections and then
// other users can use them in their outgoing webhooks and slash commands if they have permissions to manage those.
func checkOutgoingOAuthConnectionReadPermissions(c *Context) bool {
	if c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageOutgoingWebhooks) && c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSlashCommands) {
		return true
	}

	c.SetPermissionError(model.PermissionManageOutgoingWebhooks, model.PermissionManageSlashCommands)
	return false
}

// checkOutgoingOAuthConnectionWritePermissions checks if the user has the permissions to write outgoing oauth connections.
// This is a more granular permissions intended for system admins to manage (setup) outgoing oauth connections.
func checkOutgoingOAuthConnectionWritePermissions(c *Context) bool {
	if c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.OutgoingOAuthConnectionManagementPermission) {
		return true
	}

	c.SetPermissionError(model.OutgoingOAuthConnectionManagementPermission)
	return false
}

func ensureOutgoingOAuthConnectionInterface(c *Context, where string) (einterfaces.OutgoingOAuthConnectionInterface, bool) {
	if !c.App.Config().FeatureFlags.OutgoingOAuthConnections {
		c.Err = model.NewAppError(where, "api.context.outgoing_oauth_connection.not_available.feature_flag", nil, "", http.StatusNotImplemented)
		return nil, false
	}

	if c.App.OutgoingOAuthConnections() == nil || c.App.License() == nil || c.App.License().SkuShortName != model.LicenseShortSkuEnterprise {
		c.Err = model.NewAppError(where, "api.license.upgrade_needed.app_error", nil, "", http.StatusNotImplemented)
		return nil, false
	}
	return c.App.OutgoingOAuthConnections(), true
}

type listOutgoingOAuthConnectionsQuery struct {
	FromID string
	Limit  int
}

// SetDefaults sets the default values for the query.
func (q *listOutgoingOAuthConnectionsQuery) SetDefaults() {
	// Set default values
	if q.Limit == 0 {
		q.Limit = 10
	}
}

// IsValid validates the query.
func (q *listOutgoingOAuthConnectionsQuery) IsValid() error {
	if q.Limit < 1 || q.Limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}
	return nil
}

// ToFilter converts the query to a filter that can be used to query the database.
func (q *listOutgoingOAuthConnectionsQuery) ToFilter() model.OutgoingOAuthConnectionGetConnectionsFilter {
	return model.OutgoingOAuthConnectionGetConnectionsFilter{
		OffsetId: q.FromID,
		Limit:    q.Limit,
	}
}

func NewListOutgoingOAuthConnectionsQueryFromURLQuery(values url.Values) (*listOutgoingOAuthConnectionsQuery, error) {
	query := &listOutgoingOAuthConnectionsQuery{}
	query.SetDefaults()

	fromID := values.Get("from_id")
	if fromID != "" {
		query.FromID = fromID
	}

	limit := values.Get("limit")
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}
		query.Limit = limitInt
	}

	return query, nil
}

func listOutgoingOAuthConnections(c *Context, w http.ResponseWriter, r *http.Request) {
	if !checkOutgoingOAuthConnectionReadPermissions(c) {
		return
	}

	service, ok := ensureOutgoingOAuthConnectionInterface(c, whereOutgoingOAuthConnection)
	if !ok {
		return
	}

	query, err := NewListOutgoingOAuthConnectionsQueryFromURLQuery(r.URL.Query())
	if err != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.list_connections.input_error", nil, err.Error(), http.StatusBadRequest)
		return
	}

	if errValid := query.IsValid(); errValid != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.list_connections.input_error", nil, errValid.Error(), http.StatusBadRequest)
		return
	}

	connections, errList := service.GetConnections(c.AppContext, query.ToFilter())
	if errList != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.list_connections.app_error", nil, errList.Error(), http.StatusInternalServerError)
		return
	}

	service.SanitizeConnections(connections)

	if errJSON := json.NewEncoder(w).Encode(connections); errJSON != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.list_connections.app_error", nil, errJSON.Error(), http.StatusInternalServerError)
		return
	}
}

func getOutgoingOAuthConnection(c *Context, w http.ResponseWriter, r *http.Request) {
	if !checkOutgoingOAuthConnectionReadPermissions(c) {
		return
	}

	service, ok := ensureOutgoingOAuthConnectionInterface(c, whereOutgoingOAuthConnection)
	if !ok {
		return
	}

	c.RequireOutgoingOAuthConnectionId()

	connection, err := service.GetConnection(c.AppContext, c.Params.OutgoingOAuthConnectionID)
	if err != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.list_connections.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	service.SanitizeConnection(connection)

	if err := json.NewEncoder(w).Encode(connection); err != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.list_connections.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createOutgoingOAuthConnection(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("createOutgoingOauthConnection", audit.Fail)
	defer c.LogAuditRec(auditRec)
	c.LogAudit("attempt")

	if !checkOutgoingOAuthConnectionWritePermissions(c) {
		return
	}

	service, ok := ensureOutgoingOAuthConnectionInterface(c, whereOutgoingOAuthConnection)
	if !ok {
		return
	}

	var inputConnection model.OutgoingOAuthConnection
	if err := json.NewDecoder(r.Body).Decode(&inputConnection); err != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.create_connection.input_error", nil, err.Error(), http.StatusBadRequest)
		return
	}

	audit.AddEventParameterAuditable(auditRec, "outgoing_oauth_connection", &inputConnection)

	inputConnection.CreatorId = c.AppContext.Session().UserId

	connection, err := service.SaveConnection(c.AppContext, &inputConnection)
	if err != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.create_connection.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	auditRec.Success()
	auditRec.AddEventResultState(connection)
	auditRec.AddEventObjectType("outgoing_oauth_connection")
	c.LogAudit("client_id=" + connection.ClientId)

	service.SanitizeConnection(connection)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(connection); err != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.create_connection.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateOutgoingOAuthConnection(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("updateOutgoingOAuthConnection", audit.Fail)
	defer c.LogAuditRec(auditRec)
	audit.AddEventParameter(auditRec, "outgoing_oauth_connection_id", c.Params.OutgoingOAuthConnectionID)
	c.LogAudit("attempt")

	if !checkOutgoingOAuthConnectionWritePermissions(c) {
		return
	}

	service, ok := ensureOutgoingOAuthConnectionInterface(c, whereOutgoingOAuthConnection)
	if !ok {
		return
	}

	c.RequireOutgoingOAuthConnectionId()
	if c.Err != nil {
		return
	}

	var inputConnection model.OutgoingOAuthConnection
	if err := json.NewDecoder(r.Body).Decode(&inputConnection); err != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.update_connection.input_error", nil, err.Error(), http.StatusBadRequest)
		return
	}

	if inputConnection.Id != c.Params.OutgoingOAuthConnectionID {
		c.SetInvalidParam("id")
		return
	}

	currentConnection, err := service.GetConnection(c.AppContext, c.Params.OutgoingOAuthConnectionID)
	if err != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.update_connection.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
	auditRec.AddEventPriorState(currentConnection)

	connection, err := service.UpdateConnection(c.AppContext, &inputConnection)
	if err != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.update_connection.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	auditRec.AddEventObjectType("outgoing_oauth_connection")
	auditRec.AddEventResultState(connection)
	auditRec.Success()
	auditLogExtraInfo := "success"
	// Audit log changes to clientID/Client Secret
	if connection.ClientId != currentConnection.ClientId {
		auditLogExtraInfo += " new_client_id=" + connection.ClientId
	}
	if connection.ClientSecret != currentConnection.ClientSecret {
		auditLogExtraInfo += " new_client_secret"
	}
	c.LogAudit(auditLogExtraInfo)
	service.SanitizeConnection(connection)

	if err := json.NewEncoder(w).Encode(connection); err != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.update_connection.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteOutgoingOAuthConnection(c *Context, w http.ResponseWriter, r *http.Request) {
	auditRec := c.MakeAuditRecord("deleteOutgoingOAuthConnection", audit.Fail)
	defer c.LogAuditRec(auditRec)
	audit.AddEventParameter(auditRec, "outgoing_oauth_connection_id", c.Params.OutgoingOAuthConnectionID)
	c.LogAudit("attempt")

	if !checkOutgoingOAuthConnectionWritePermissions(c) {
		return
	}

	service, ok := ensureOutgoingOAuthConnectionInterface(c, whereOutgoingOAuthConnection)
	if !ok {
		return
	}

	c.RequireOutgoingOAuthConnectionId()
	if c.Err != nil {
		return
	}

	connection, err := service.GetConnection(c.AppContext, c.Params.OutgoingOAuthConnectionID)
	if err != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.delete_connection.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
	auditRec.AddEventPriorState(connection)

	if err := service.DeleteConnection(c.AppContext, c.Params.OutgoingOAuthConnectionID); err != nil {
		c.Err = model.NewAppError(whereOutgoingOAuthConnection, "api.context.outgoing_oauth_connection.delete_connection.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	auditRec.AddEventObjectType("outgoing_oauth_connection")
	auditRec.Success()

	ReturnStatusOK(w)
}
