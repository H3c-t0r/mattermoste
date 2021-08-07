// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/v6/model"
)

func TestCreateOAuthApp(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()
	client := th.Client
	adminClient := th.SystemAdminClient

	defaultRolePermissions := th.SaveDefaultRolePermissions()
	enableOAuthServiceProvider := th.App.Config().ServiceSettings.EnableOAuthServiceProvider
	defer func() {
		th.RestoreDefaultRolePermissions(defaultRolePermissions)
		th.App.UpdateConfig(func(cfg *model.Config) { cfg.ServiceSettings.EnableOAuthServiceProvider = enableOAuthServiceProvider })
	}()

	// Grant permission to regular users.
	th.AddPermissionToRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = true })

	oapp := &model.OAuthApp{Name: GenerateTestAppName(), Homepage: "https://nowhere.com", Description: "test", CallbackUrls: []string{"https://nowhere.com"}, IsTrusted: true}

	rapp, resp, err := adminClient.CreateOAuthApp(oapp)
	require.NoError(t, err)
	CheckCreatedStatus(t, resp)
	assert.Equal(t, oapp.Name, rapp.Name, "names did not match")
	assert.Equal(t, oapp.IsTrusted, rapp.IsTrusted, "trusted did no match")

	// Revoke permission from regular users.
	th.RemovePermissionFromRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)

	_, resp, _ = client.CreateOAuthApp(oapp)
	CheckForbiddenStatus(t, resp)
	// Grant permission to regular users.
	th.AddPermissionToRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)

	rapp, resp, err = client.CreateOAuthApp(oapp)
	require.NoError(t, err)
	CheckCreatedStatus(t, resp)

	assert.False(t, rapp.IsTrusted, "trusted should be false - created by non admin")

	oapp.Name = ""
	_, resp, _ = adminClient.CreateOAuthApp(oapp)
	CheckBadRequestStatus(t, resp)

	r, err := client.DoApiPost("/oauth/apps", "garbage")
	require.NotNil(t, err, "expected error from garbage post")
	assert.Equal(t, http.StatusBadRequest, r.StatusCode)

	client.Logout()
	_, resp, _ = client.CreateOAuthApp(oapp)
	CheckUnauthorizedStatus(t, resp)

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = false })
	oapp.Name = GenerateTestAppName()
	_, resp, _ = adminClient.CreateOAuthApp(oapp)
	CheckNotImplementedStatus(t, resp)
}

func TestUpdateOAuthApp(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	client := th.Client
	adminClient := th.SystemAdminClient

	defaultRolePermissions := th.SaveDefaultRolePermissions()
	enableOAuthServiceProvider := th.App.Config().ServiceSettings.EnableOAuthServiceProvider
	defer func() {
		th.RestoreDefaultRolePermissions(defaultRolePermissions)
		th.App.UpdateConfig(func(cfg *model.Config) { cfg.ServiceSettings.EnableOAuthServiceProvider = enableOAuthServiceProvider })
	}()

	// Grant permission to regular users.
	th.AddPermissionToRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)
	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = true })

	oapp := &model.OAuthApp{
		Name:         "oapp",
		IsTrusted:    false,
		IconURL:      "https://nowhere.com/img",
		Homepage:     "https://nowhere.com",
		Description:  "test",
		CallbackUrls: []string{"https://callback.com"},
	}

	oapp, _, _ = adminClient.CreateOAuthApp(oapp)

	oapp.Name = "oapp_update"
	oapp.IsTrusted = true
	oapp.IconURL = "https://nowhere.com/img_update"
	oapp.Homepage = "https://nowhere_update.com"
	oapp.Description = "test_update"
	oapp.CallbackUrls = []string{"https://callback_update.com", "https://another_callback.com"}

	updatedApp, _, err = adminClient.UpdateOAuthApp(oapp)
	require.NoError(t, err)
	assert.Equal(t, oapp.Id, updatedApp.Id, "Id should have not updated")
	assert.Equal(t, oapp.CreatorId, updatedApp.CreatorId, "CreatorId should have not updated")
	assert.Equal(t, oapp.CreateAt, updatedApp.CreateAt, "CreateAt should have not updated")
	assert.NotEqual(t, oapp.UpdateAt, updatedApp.UpdateAt, "UpdateAt should have updated")
	assert.Equal(t, oapp.ClientSecret, updatedApp.ClientSecret, "ClientSecret should have not updated")
	assert.Equal(t, oapp.Name, updatedApp.Name, "Name should have updated")
	assert.Equal(t, oapp.Description, updatedApp.Description, "Description should have updated")
	assert.Equal(t, oapp.IconURL, updatedApp.IconURL, "IconURL should have updated")

	if len(updatedApp.CallbackUrls) == len(oapp.CallbackUrls) {
		for i, callbackUrl := range updatedApp.CallbackUrls {
			assert.Equal(t, oapp.CallbackUrls[i], callbackUrl, "Description should have updated")
		}
	}
	assert.Equal(t, oapp.Homepage, updatedApp.Homepage, "Homepage should have updated")
	assert.Equal(t, oapp.IsTrusted, updatedApp.IsTrusted, "IsTrusted should have updated")

	th.LoginBasic2()
	updatedApp.CreatorId = th.BasicUser2.Id
	_, resp, _ = client.UpdateOAuthApp(oapp)
	CheckForbiddenStatus(t, resp)

	th.LoginBasic()

	// Revoke permission from regular users.
	th.RemovePermissionFromRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)

	_, resp, _ = client.UpdateOAuthApp(oapp)
	CheckForbiddenStatus(t, resp)

	oapp.Id = "zhk9d1ggatrqz236c7h87im7bc"
	_, resp, _ = adminClient.UpdateOAuthApp(oapp)
	CheckNotFoundStatus(t, resp)

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = false })

	_, resp, _ = adminClient.UpdateOAuthApp(oapp)
	CheckNotImplementedStatus(t, resp)

	client.Logout()
	_, resp, _ = client.UpdateOAuthApp(oapp)
	CheckUnauthorizedStatus(t, resp)

	oapp.Id = "junk"
	_, resp, _ = adminClient.UpdateOAuthApp(oapp)
	CheckBadRequestStatus(t, resp)

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = true })
	th.AddPermissionToRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)
	th.LoginBasic()

	userOapp := &model.OAuthApp{
		Name:         "useroapp",
		IsTrusted:    false,
		IconURL:      "https://nowhere.com/img",
		Homepage:     "https://nowhere.com",
		Description:  "test",
		CallbackUrls: []string{"https://callback.com"},
	}

	userOapp, _, err = client.CreateOAuthApp(userOapp)
	require.NoError(t, err)

	userOapp.IsTrusted = true
	userOapp, _, err = client.UpdateOAuthApp(userOapp)
	require.NoError(t, err)
	assert.False(t, userOapp.IsTrusted)

	userOapp.IsTrusted = true
	userOapp, _, err = adminClient.UpdateOAuthApp(userOapp)
	require.NoError(t, err)
	assert.True(t, userOapp.IsTrusted)

	userOapp.IsTrusted = false
	userOapp, _, err = client.UpdateOAuthApp(userOapp)
	require.NoError(t, err)
	assert.True(t, userOapp.IsTrusted)
}

func TestGetOAuthApps(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()
	client := th.Client
	adminClient := th.SystemAdminClient

	defaultRolePermissions := th.SaveDefaultRolePermissions()
	enableOAuthServiceProvider := th.App.Config().ServiceSettings.EnableOAuthServiceProvider
	defer func() {
		th.RestoreDefaultRolePermissions(defaultRolePermissions)
		th.App.UpdateConfig(func(cfg *model.Config) { cfg.ServiceSettings.EnableOAuthServiceProvider = enableOAuthServiceProvider })
	}()

	// Grant permission to regular users.
	th.AddPermissionToRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)
	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = true })

	oapp := &model.OAuthApp{Name: GenerateTestAppName(), Homepage: "https://nowhere.com", Description: "test", CallbackUrls: []string{"https://nowhere.com"}}

	rapp, _, err = adminClient.CreateOAuthApp(oapp)
	require.NoError(t, err)

	oapp.Name = GenerateTestAppName()
	rapp2, _, err = client.CreateOAuthApp(oapp)
	require.NoError(t, err)

	apps, _, err = adminClient.GetOAuthApps(0, 1000)
	require.NoError(t, err)

	found1 := false
	found2 := false
	for _, a := range apps {
		if a.Id == rapp.Id {
			found1 = true
		}
		if a.Id == rapp2.Id {
			found2 = true
		}
	}
	assert.Truef(t, found1, "missing oauth app %v", rapp.Id)
	assert.Truef(t, found2, "missing oauth app %v", rapp2.Id)

	apps, _, err = adminClient.GetOAuthApps(1, 1)
	require.NoError(t, err)
	require.Equal(t, 1, len(apps), "paging failed")

	apps, _, err = client.GetOAuthApps(0, 1000)
	require.NoError(t, err)
	require.True(t, len(apps) == 1 || apps[0].Id == rapp2.Id, "wrong apps returned")

	// Revoke permission from regular users.
	th.RemovePermissionFromRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)

	_, resp, _ = client.GetOAuthApps(0, 1000)
	CheckForbiddenStatus(t, resp)

	client.Logout()

	_, resp, _ = client.GetOAuthApps(0, 1000)
	CheckUnauthorizedStatus(t, resp)

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = false })
	_, resp, _ = adminClient.GetOAuthApps(0, 1000)
	CheckNotImplementedStatus(t, resp)
}

func TestGetOAuthApp(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()
	client := th.Client
	adminClient := th.SystemAdminClient

	defaultRolePermissions := th.SaveDefaultRolePermissions()
	enableOAuthServiceProvider := th.App.Config().ServiceSettings.EnableOAuthServiceProvider
	defer func() {
		th.RestoreDefaultRolePermissions(defaultRolePermissions)
		th.App.UpdateConfig(func(cfg *model.Config) { cfg.ServiceSettings.EnableOAuthServiceProvider = enableOAuthServiceProvider })
	}()

	// Grant permission to regular users.
	th.AddPermissionToRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)
	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = true })

	oapp := &model.OAuthApp{Name: GenerateTestAppName(), Homepage: "https://nowhere.com", Description: "test", CallbackUrls: []string{"https://nowhere.com"}}

	rapp, _, err = adminClient.CreateOAuthApp(oapp)
	require.NoError(t, err)

	oapp.Name = GenerateTestAppName()
	rapp2, _, err = client.CreateOAuthApp(oapp)
	require.NoError(t, err)

	rrapp, _, err = adminClient.GetOAuthApp(rapp.Id)
	require.NoError(t, err)
	assert.Equal(t, rapp.Id, rrapp.Id, "wrong app")
	assert.NotEqual(t, "", rrapp.ClientSecret, "should not be sanitized")

	rrapp2, _, err = adminClient.GetOAuthApp(rapp2.Id)
	require.NoError(t, err)
	assert.Equal(t, rapp2.Id, rrapp2.Id, "wrong app")
	assert.NotEqual(t, "", rrapp2.ClientSecret, "should not be sanitized")

	_, _, err = client.GetOAuthApp(rapp2.Id)
	require.NoError(t, err)

	_, resp, _ = client.GetOAuthApp(rapp.Id)
	CheckForbiddenStatus(t, resp)

	// Revoke permission from regular users.
	th.RemovePermissionFromRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)

	_, resp, _ = client.GetOAuthApp(rapp2.Id)
	CheckForbiddenStatus(t, resp)

	client.Logout()

	_, resp, _ = client.GetOAuthApp(rapp2.Id)
	CheckUnauthorizedStatus(t, resp)

	_, resp, _ = adminClient.GetOAuthApp("junk")
	CheckBadRequestStatus(t, resp)

	_, resp, _ = adminClient.GetOAuthApp(model.NewId())
	CheckNotFoundStatus(t, resp)

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = false })
	_, resp, _ = adminClient.GetOAuthApp(rapp.Id)
	CheckNotImplementedStatus(t, resp)
}

func TestGetOAuthAppInfo(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()
	client := th.Client
	adminClient := th.SystemAdminClient

	defaultRolePermissions := th.SaveDefaultRolePermissions()
	enableOAuthServiceProvider := th.App.Config().ServiceSettings.EnableOAuthServiceProvider
	defer func() {
		th.RestoreDefaultRolePermissions(defaultRolePermissions)
		th.App.UpdateConfig(func(cfg *model.Config) { cfg.ServiceSettings.EnableOAuthServiceProvider = enableOAuthServiceProvider })
	}()

	// Grant permission to regular users.
	th.AddPermissionToRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)
	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = true })

	oapp := &model.OAuthApp{Name: GenerateTestAppName(), Homepage: "https://nowhere.com", Description: "test", CallbackUrls: []string{"https://nowhere.com"}}

	rapp, _, err = adminClient.CreateOAuthApp(oapp)
	require.NoError(t, err)

	oapp.Name = GenerateTestAppName()
	rapp2, _, err = client.CreateOAuthApp(oapp)
	require.NoError(t, err)

	rrapp, _, err = adminClient.GetOAuthAppInfo(rapp.Id)
	require.NoError(t, err)
	assert.Equal(t, rapp.Id, rrapp.Id, "wrong app")
	assert.Equal(t, "", rrapp.ClientSecret, "should be sanitized")

	rrapp2, _, err = adminClient.GetOAuthAppInfo(rapp2.Id)
	require.NoError(t, err)
	assert.Equal(t, rapp2.Id, rrapp2.Id, "wrong app")
	assert.Equal(t, "", rrapp2.ClientSecret, "should be sanitized")

	_, _, err = client.GetOAuthAppInfo(rapp2.Id)
	require.NoError(t, err)

	_, _, err = client.GetOAuthAppInfo(rapp.Id)
	require.NoError(t, err)

	// Revoke permission from regular users.
	th.RemovePermissionFromRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)

	_, _, err = client.GetOAuthAppInfo(rapp2.Id)
	require.NoError(t, err)

	client.Logout()

	_, resp, _ = client.GetOAuthAppInfo(rapp2.Id)
	CheckUnauthorizedStatus(t, resp)

	_, resp, _ = adminClient.GetOAuthAppInfo("junk")
	CheckBadRequestStatus(t, resp)

	_, resp, _ = adminClient.GetOAuthAppInfo(model.NewId())
	CheckNotFoundStatus(t, resp)

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = false })
	_, resp, _ = adminClient.GetOAuthAppInfo(rapp.Id)
	CheckNotImplementedStatus(t, resp)
}

func TestDeleteOAuthApp(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()
	client := th.Client
	adminClient := th.SystemAdminClient

	defaultRolePermissions := th.SaveDefaultRolePermissions()
	enableOAuthServiceProvider := th.App.Config().ServiceSettings.EnableOAuthServiceProvider
	defer func() {
		th.RestoreDefaultRolePermissions(defaultRolePermissions)
		th.App.UpdateConfig(func(cfg *model.Config) { cfg.ServiceSettings.EnableOAuthServiceProvider = enableOAuthServiceProvider })
	}()

	// Grant permission to regular users.
	th.AddPermissionToRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)
	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = true })

	oapp := &model.OAuthApp{Name: GenerateTestAppName(), Homepage: "https://nowhere.com", Description: "test", CallbackUrls: []string{"https://nowhere.com"}}

	rapp, _, err = adminClient.CreateOAuthApp(oapp)
	require.NoError(t, err)

	oapp.Name = GenerateTestAppName()
	rapp2, _, err = client.CreateOAuthApp(oapp)
	require.NoError(t, err)

	pass, _, err = adminClient.DeleteOAuthApp(rapp.Id)
	require.NoError(t, err)
	assert.True(t, pass, "should have passed")

	_, _, err = adminClient.DeleteOAuthApp(rapp2.Id)
	require.NoError(t, err)

	rapp, _, err = adminClient.CreateOAuthApp(oapp)
	require.NoError(t, err)

	oapp.Name = GenerateTestAppName()
	rapp2, _, err = client.CreateOAuthApp(oapp)
	require.NoError(t, err)

	_, resp, _ = client.DeleteOAuthApp(rapp.Id)
	CheckForbiddenStatus(t, resp)

	_, _, err = client.DeleteOAuthApp(rapp2.Id)
	require.NoError(t, err)

	// Revoke permission from regular users.
	th.RemovePermissionFromRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)

	_, resp, _ = client.DeleteOAuthApp(rapp.Id)
	CheckForbiddenStatus(t, resp)

	client.Logout()
	_, resp, _ = client.DeleteOAuthApp(rapp.Id)
	CheckUnauthorizedStatus(t, resp)

	_, resp, _ = adminClient.DeleteOAuthApp("junk")
	CheckBadRequestStatus(t, resp)

	_, resp, _ = adminClient.DeleteOAuthApp(model.NewId())
	CheckNotFoundStatus(t, resp)

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = false })
	_, resp, _ = adminClient.DeleteOAuthApp(rapp.Id)
	CheckNotImplementedStatus(t, resp)
}

func TestRegenerateOAuthAppSecret(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()
	client := th.Client
	adminClient := th.SystemAdminClient

	defaultRolePermissions := th.SaveDefaultRolePermissions()
	enableOAuthServiceProvider := th.App.Config().ServiceSettings.EnableOAuthServiceProvider
	defer func() {
		th.RestoreDefaultRolePermissions(defaultRolePermissions)
		th.App.UpdateConfig(func(cfg *model.Config) { cfg.ServiceSettings.EnableOAuthServiceProvider = enableOAuthServiceProvider })
	}()

	// Grant permission to regular users.
	th.AddPermissionToRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)
	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = true })

	oapp := &model.OAuthApp{Name: GenerateTestAppName(), Homepage: "https://nowhere.com", Description: "test", CallbackUrls: []string{"https://nowhere.com"}}

	rapp, _, err = adminClient.CreateOAuthApp(oapp)
	require.NoError(t, err)

	oapp.Name = GenerateTestAppName()
	rapp2, _, err = client.CreateOAuthApp(oapp)
	require.NoError(t, err)

	rrapp, _, err = adminClient.RegenerateOAuthAppSecret(rapp.Id)
	require.NoError(t, err)
	assert.Equal(t, rrapp.Id, rapp.Id, "wrong app")
	assert.NotEqual(t, rapp.ClientSecret, rrapp.ClientSecret, "secret didn't change")

	_, _, err = adminClient.RegenerateOAuthAppSecret(rapp2.Id)
	require.NoError(t, err)

	rapp, _, err = adminClient.CreateOAuthApp(oapp)
	require.NoError(t, err)

	oapp.Name = GenerateTestAppName()
	rapp2, _, err = client.CreateOAuthApp(oapp)
	require.NoError(t, err)

	_, resp, _ = client.RegenerateOAuthAppSecret(rapp.Id)
	CheckForbiddenStatus(t, resp)

	_, _, err = client.RegenerateOAuthAppSecret(rapp2.Id)
	require.NoError(t, err)

	// Revoke permission from regular users.
	th.RemovePermissionFromRole(model.PermissionManageOAuth.Id, model.SystemUserRoleId)

	_, resp, _ = client.RegenerateOAuthAppSecret(rapp.Id)
	CheckForbiddenStatus(t, resp)

	client.Logout()
	_, resp, _ = client.RegenerateOAuthAppSecret(rapp.Id)
	CheckUnauthorizedStatus(t, resp)

	_, resp, _ = adminClient.RegenerateOAuthAppSecret("junk")
	CheckBadRequestStatus(t, resp)

	_, resp, _ = adminClient.RegenerateOAuthAppSecret(model.NewId())
	CheckNotFoundStatus(t, resp)

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = false })
	_, resp, _ = adminClient.RegenerateOAuthAppSecret(rapp.Id)
	CheckNotImplementedStatus(t, resp)
}

func TestGetAuthorizedOAuthAppsForUser(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	client := th.Client
	adminClient := th.SystemAdminClient

	enableOAuth := th.App.Config().ServiceSettings.EnableOAuthServiceProvider
	defer func() {
		th.App.UpdateConfig(func(cfg *model.Config) { cfg.ServiceSettings.EnableOAuthServiceProvider = enableOAuth })
	}()
	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.EnableOAuthServiceProvider = true })

	oapp := &model.OAuthApp{Name: GenerateTestAppName(), Homepage: "https://nowhere.com", Description: "test", CallbackUrls: []string{"https://nowhere.com"}}

	rapp, _, err = adminClient.CreateOAuthApp(oapp)
	require.NoError(t, err)

	authRequest := &model.AuthorizeRequest{
		ResponseType: model.AuthCodeResponseType,
		ClientId:     rapp.Id,
		RedirectUri:  rapp.CallbackUrls[0],
		Scope:        "",
		State:        "123",
	}

	_, _, err = client.AuthorizeOAuthApp(authRequest)
	require.NoError(t, err)

	apps, _, err = client.GetAuthorizedOAuthAppsForUser(th.BasicUser.Id, 0, 1000)
	require.NoError(t, err)

	found := false
	for _, a := range apps {
		if a.Id == rapp.Id {
			found = true
		}
		assert.Equal(t, "", a.ClientSecret, "not sanitized")
	}
	require.True(t, found, "missing app")

	_, resp, _ = client.GetAuthorizedOAuthAppsForUser(th.BasicUser2.Id, 0, 1000)
	CheckForbiddenStatus(t, resp)

	_, resp, _ = client.GetAuthorizedOAuthAppsForUser("junk", 0, 1000)
	CheckBadRequestStatus(t, resp)

	client.Logout()
	_, resp, _ = client.GetAuthorizedOAuthAppsForUser(th.BasicUser.Id, 0, 1000)
	CheckUnauthorizedStatus(t, resp)

	_, _, err = adminClient.GetAuthorizedOAuthAppsForUser(th.BasicUser.Id, 0, 1000)
	require.NoError(t, err)
}

func closeBody(r *http.Response) {
	if r != nil && r.Body != nil {
		ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
}

func TestNilAuthorizeOAuthApp(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	client := th.Client

	_, _, err := client.AuthorizeOAuthApp(nil)
	require.Error(t, err)
	CheckErrorID(t, err, "api.context.invalid_body_param.app_error")
}
