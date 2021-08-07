// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/v6/model"
)

func TestGetRole(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	role := &model.Role{
		Name:          model.NewId(),
		DisplayName:   model.NewId(),
		Description:   model.NewId(),
		Permissions:   []string{"manage_system", "create_public_channel"},
		SchemeManaged: true,
	}

	role, err := th.App.Srv().Store.Role().Save(role)
	require.NoError(t, err)
	defer th.App.Srv().Store.Job().Delete(role.Id)

	th.TestForAllClients(t, func(t *testing.T, client *model.Client4) {
		received, _, err := client.GetRole(role.Id)
		require.NoError(t, err)

		assert.Equal(t, received.Id, role.Id)
		assert.Equal(t, received.Name, role.Name)
		assert.Equal(t, received.DisplayName, role.DisplayName)
		assert.Equal(t, received.Description, role.Description)
		assert.EqualValues(t, received.Permissions, role.Permissions)
		assert.Equal(t, received.SchemeManaged, role.SchemeManaged)
	})

	th.TestForSystemAdminAndLocal(t, func(t *testing.T, client *model.Client4) {
		_, resp, _ := client.GetRole("1234")
		CheckBadRequestStatus(t, resp)

		_, resp, _ = client.GetRole(model.NewId())
		CheckNotFoundStatus(t, resp)
	})
}

func TestGetRoleByName(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	role := &model.Role{
		Name:          model.NewId(),
		DisplayName:   model.NewId(),
		Description:   model.NewId(),
		Permissions:   []string{"manage_system", "create_public_channel"},
		SchemeManaged: true,
	}

	role, err := th.App.Srv().Store.Role().Save(role)
	assert.NoError(t, err)
	defer th.App.Srv().Store.Job().Delete(role.Id)

	th.TestForAllClients(t, func(t *testing.T, client *model.Client4) {
		received, _, err = client.GetRoleByName(role.Name)
		require.NoError(t, err)

		assert.Equal(t, received.Id, role.Id)
		assert.Equal(t, received.Name, role.Name)
		assert.Equal(t, received.DisplayName, role.DisplayName)
		assert.Equal(t, received.Description, role.Description)
		assert.EqualValues(t, received.Permissions, role.Permissions)
		assert.Equal(t, received.SchemeManaged, role.SchemeManaged)
	})

	th.TestForSystemAdminAndLocal(t, func(t *testing.T, client *model.Client4) {
		_, resp, _ := client.GetRoleByName(strings.Repeat("abcdefghij", 10))
		CheckBadRequestStatus(t, resp)

		_, resp, _ = client.GetRoleByName(model.NewId())
		CheckNotFoundStatus(t, resp)
	})
}

func TestGetRolesByNames(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	role1 := &model.Role{
		Name:          model.NewId(),
		DisplayName:   model.NewId(),
		Description:   model.NewId(),
		Permissions:   []string{"manage_system", "create_public_channel"},
		SchemeManaged: true,
	}
	role2 := &model.Role{
		Name:          model.NewId(),
		DisplayName:   model.NewId(),
		Description:   model.NewId(),
		Permissions:   []string{"manage_system", "delete_private_channel"},
		SchemeManaged: true,
	}
	role3 := &model.Role{
		Name:          model.NewId(),
		DisplayName:   model.NewId(),
		Description:   model.NewId(),
		Permissions:   []string{"manage_system", "manage_public_channel_properties"},
		SchemeManaged: true,
	}

	role1, err := th.App.Srv().Store.Role().Save(role1)
	assert.NoError(t, err)
	defer th.App.Srv().Store.Job().Delete(role1.Id)

	role2, err = th.App.Srv().Store.Role().Save(role2)
	assert.NoError(t, err)
	defer th.App.Srv().Store.Job().Delete(role2.Id)

	role3, err = th.App.Srv().Store.Role().Save(role3)
	assert.NoError(t, err)
	defer th.App.Srv().Store.Job().Delete(role3.Id)

	th.TestForAllClients(t, func(t *testing.T, client *model.Client4) {
		// Check all three roles can be found.
		received, _, err = client.GetRolesByNames([]string{role1.Name, role2.Name, role3.Name})
		require.NoError(t, err)

		assert.Contains(t, received, role1)
		assert.Contains(t, received, role2)
		assert.Contains(t, received, role3)

		// Check a list of non-existent roles.
		_, _, err = client.GetRolesByNames([]string{model.NewId(), model.NewId()})
		require.NoError(t, err)
	})

	th.TestForSystemAdminAndLocal(t, func(t *testing.T, client *model.Client4) {
		// Empty list should error.
		_, resp, _ := client.GetRolesByNames([]string{})
		CheckBadRequestStatus(t, resp)
	})

	th.TestForAllClients(t, func(t *testing.T, client *model.Client4) {
		// Invalid role name should error.
		_, resp, _ := client.GetRolesByNames([]string{model.NewId(), model.NewId(), "!!!!!!"})
		CheckBadRequestStatus(t, resp)

		// Empty/whitespace rolenames should be ignored.
		_, _, err = client.GetRolesByNames([]string{model.NewId(), model.NewId(), "", "    "})
		require.NoError(t, err)
	})

}

func TestPatchRole(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	role := &model.Role{
		Name:          model.NewId(),
		DisplayName:   model.NewId(),
		Description:   model.NewId(),
		Permissions:   []string{"manage_system", "create_public_channel", "manage_slash_commands"},
		SchemeManaged: true,
	}

	role, err := th.App.Srv().Store.Role().Save(role)
	assert.NoError(t, err)
	defer th.App.Srv().Store.Job().Delete(role.Id)

	patch := &model.RolePatch{
		Permissions: &[]string{"manage_system", "create_public_channel", "manage_incoming_webhooks", "manage_outgoing_webhooks"},
	}

	th.TestForSystemAdminAndLocal(t, func(t *testing.T, client *model.Client4) {

		// Cannot edit a system admin
		adminRole, err := th.App.Srv().Store.Role().GetByName(context.Background(), "system_admin")
		assert.NoError(t, err)
		defer th.App.Srv().Store.Job().Delete(adminRole.Id)

		_, resp, _ := client.PatchRole(adminRole.Id, patch)
		CheckNotImplementedStatus(t, resp)

		// Cannot give other roles read / write to system roles or manage roles because only system admin can do these actions
		systemManager, err := th.App.Srv().Store.Role().GetByName(context.Background(), "system_manager")
		assert.NoError(t, err)
		defer th.App.Srv().Store.Job().Delete(systemManager.Id)

		patchWriteSystemRoles := &model.RolePatch{
			Permissions: &[]string{model.PermissionSysconsoleWriteUserManagementSystemRoles.Id},
		}

		_, resp, _ = client.PatchRole(systemManager.Id, patchWriteSystemRoles)
		CheckNotImplementedStatus(t, resp)

		patchReadSystemRoles := &model.RolePatch{
			Permissions: &[]string{model.PermissionSysconsoleReadUserManagementSystemRoles.Id},
		}

		_, resp, _ = client.PatchRole(systemManager.Id, patchReadSystemRoles)
		CheckNotImplementedStatus(t, resp)

		patchManageRoles := &model.RolePatch{
			Permissions: &[]string{model.PermissionManageRoles.Id},
		}

		_, resp, _ = client.PatchRole(systemManager.Id, patchManageRoles)
		CheckNotImplementedStatus(t, resp)
	})

	th.TestForSystemAdminAndLocal(t, func(t *testing.T, client *model.Client4) {
		received, _, err = client.PatchRole(role.Id, patch)
		require.NoError(t, err)

		assert.Equal(t, received.Id, role.Id)
		assert.Equal(t, received.Name, role.Name)
		assert.Equal(t, received.DisplayName, role.DisplayName)
		assert.Equal(t, received.Description, role.Description)
		perms := []string{"manage_system", "create_public_channel", "manage_incoming_webhooks", "manage_outgoing_webhooks"}
		sort.Strings(perms)
		assert.EqualValues(t, received.Permissions, perms)
		assert.Equal(t, received.SchemeManaged, role.SchemeManaged)

		// Check a no-op patch succeeds.
		_, _, err = client.PatchRole(role.Id, patch)
		require.NoError(t, err)

		_, resp, _ = client.PatchRole("junk", patch)
		CheckBadRequestStatus(t, resp)
	})

	_, resp, _ := th.Client.PatchRole(model.NewId(), patch)
	CheckNotFoundStatus(t, resp)

	_, resp, _ = th.Client.PatchRole(role.Id, patch)
	CheckForbiddenStatus(t, resp)

	patch = &model.RolePatch{
		Permissions: &[]string{"manage_system", "manage_incoming_webhooks", "manage_outgoing_webhooks"},
	}

	th.TestForSystemAdminAndLocal(t, func(t *testing.T, client *model.Client4) {
		received, _, err = client.PatchRole(role.Id, patch)
		require.NoError(t, err)

		assert.Equal(t, received.Id, role.Id)
		assert.Equal(t, received.Name, role.Name)
		assert.Equal(t, received.DisplayName, role.DisplayName)
		assert.Equal(t, received.Description, role.Description)
		perms := []string{"manage_system", "manage_incoming_webhooks", "manage_outgoing_webhooks"}
		sort.Strings(perms)
		assert.EqualValues(t, received.Permissions, perms)
		assert.Equal(t, received.SchemeManaged, role.SchemeManaged)

		t.Run("Check guest permissions editing without E20 license", func(t *testing.T) {
			license := model.NewTestLicense()
			license.Features.GuestAccountsPermissions = model.NewBool(false)
			th.App.Srv().SetLicense(license)

			guestRole, err := th.App.Srv().Store.Role().GetByName(context.Background(), "system_guest")
			require.NoError(t, err)
			received, resp, _ = client.PatchRole(guestRole.Id, patch)
			CheckNotImplementedStatus(t, resp)
		})

		t.Run("Check guest permissions editing with E20 license", func(t *testing.T) {
			license := model.NewTestLicense()
			license.Features.GuestAccountsPermissions = model.NewBool(true)
			th.App.Srv().SetLicense(license)
			guestRole, err := th.App.Srv().Store.Role().GetByName(context.Background(), "system_guest")
			require.NoError(t, err)
			_, _, err = client.PatchRole(guestRole.Id, patch)
			require.NoError(t, err)
		})
	})
}
