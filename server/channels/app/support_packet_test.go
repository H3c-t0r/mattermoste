// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/v8/channels/app/platform"
	fmocks "github.com/mattermost/mattermost/server/v8/platform/shared/filestore/mocks"
)

func TestCreatePluginsFile(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	// Happy path where we have a plugins file with no err
	fileData, err := th.App.createPluginsFile(th.Context)
	require.NotNil(t, fileData)
	assert.Equal(t, "plugins.json", fileData.Filename)
	assert.Positive(t, len(fileData.Body))
	assert.NoError(t, err)

	// Turn off plugins so we can get an error
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.PluginSettings.Enable = false
	})

	// Plugins off in settings so no fileData and we get a warning instead
	fileData, err = th.App.createPluginsFile(th.Context)
	assert.Nil(t, fileData)
	assert.ErrorContains(t, err, "failed to get plugin list for support package")
}

func TestGenerateSupportPacketYaml(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	licenseUsers := 100
	license := model.NewTestLicense()
	license.Features.Users = model.NewInt(licenseUsers)
	th.App.Srv().SetLicense(license)

	generateSupportPacket := func(t *testing.T) *model.SupportPacket {
		t.Helper()

		fileData, err := th.App.generateSupportPacketYaml(th.Context)
		require.NotNil(t, fileData)
		assert.Equal(t, "support_packet.yaml", fileData.Filename)
		assert.Positive(t, len(fileData.Body))
		assert.NoError(t, err)

		var packet model.SupportPacket
		require.NoError(t, yaml.Unmarshal(fileData.Body, &packet))
		require.NotNil(t, packet)
		return &packet
	}

	t.Run("Happy path", func(t *testing.T) {
		// Happy path where we have a support packet yaml file without any warnings
		packet := generateSupportPacket(t)

		/* Build information */
		assert.NotEmpty(t, packet.ServerOS)
		assert.NotEmpty(t, packet.ServerArchitecture)
		assert.Equal(t, model.CurrentVersion, packet.ServerVersion)
		// BuildHash is not present in tests

		/* DB */
		assert.NotEmpty(t, packet.DatabaseType)
		assert.NotEmpty(t, packet.DatabaseVersion)
		assert.NotEmpty(t, packet.DatabaseSchemaVersion)
		assert.Zero(t, packet.WebsocketConnections)
		assert.NotZero(t, packet.MasterDbConnections)
		assert.Zero(t, packet.ReplicaDbConnections)

		/* Cluster */
		assert.Empty(t, packet.ClusterID)

		/* File store */
		assert.Equal(t, "local", packet.FileDriver)
		assert.Equal(t, "OK", packet.FileStatus)

		/* LDAP */
		assert.Empty(t, packet.LdapVendorName)
		assert.Empty(t, packet.LdapVendorVersion)

		/* Elastic Search */
		assert.Empty(t, packet.ElasticServerVersion)
		assert.Empty(t, packet.ElasticServerPlugins)

		/* License */
		assert.Equal(t, "My awesome Company", packet.LicenseTo)
		assert.Equal(t, licenseUsers, packet.LicenseSupportedUsers)
		assert.Equal(t, false, packet.LicenseIsTrial)

		/* Server stats */
		assert.Equal(t, 3, packet.ActiveUsers) // from InitBasic()
		assert.Equal(t, 0, packet.DailyActiveUsers)
		assert.Equal(t, 0, packet.MonthlyActiveUsers)
		assert.Equal(t, 0, packet.InactiveUserCount)
		assert.Equal(t, 5, packet.TotalPosts)    // from InitBasic()
		assert.Equal(t, 3, packet.TotalChannels) // from InitBasic()
		assert.Equal(t, 1, packet.TotalTeams)    // from InitBasic()

		/* Jobs */
		assert.Empty(t, packet.DataRetentionJobs)
		assert.Empty(t, packet.MessageExportJobs)
		assert.Empty(t, packet.ElasticPostIndexingJobs)
		assert.Empty(t, packet.ElasticPostAggregationJobs)
		assert.Empty(t, packet.BlevePostIndexingJobs)
		assert.Empty(t, packet.LdapSyncJobs)
		assert.Empty(t, packet.MigrationJobs)
		assert.Empty(t, packet.ComplianceJobs)
	})

	t.Run("post count should be present if number of users extends AnalyticsSettings.MaxUsersForStatistics", func(t *testing.T) {
		th.App.UpdateConfig(func(cfg *model.Config) {
			cfg.AnalyticsSettings.MaxUsersForStatistics = model.NewInt(1)
		})

		for i := 0; i < 5; i++ {
			p := th.CreatePost(th.BasicChannel)
			require.NotNil(t, p)
		}

		// InitBasic() already creats 5 posts
		packet := generateSupportPacket(t)
		assert.Equal(t, 10, packet.TotalPosts)
	})

	t.Run("filestore fails", func(t *testing.T) {
		fb := &fmocks.FileBackend{}
		platform.SetFileStore(fb)(th.Server.Platform())
		fb.On("DriverName").Return("mock")
		fb.On("TestConnection").Return(errors.New("all broken"))

		packet := generateSupportPacket(t)

		assert.Equal(t, "mock", packet.FileDriver)
		assert.Equal(t, "FAIL: all broken", packet.FileStatus)
	})
}

func TestGenerateSupportPacket(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	d1 := []byte("hello\ngo\n")
	err := os.WriteFile("mattermost.log", d1, 0777)
	require.NoError(t, err)
	err = os.WriteFile("notifications.log", d1, 0777)
	require.NoError(t, err)

	fileDatas := th.App.GenerateSupportPacket(th.Context)
	var rFileNames []string
	testFiles := []string{
		"support_packet.yaml",
		"plugins.json",
		"sanitized_config.json",
		"mattermost.log",
		"notifications.log",
		"cpu.prof",
		"heap.prof",
		"goroutines",
	}
	for _, fileData := range fileDatas {
		require.NotNil(t, fileData)
		assert.Positive(t, len(fileData.Body))

		rFileNames = append(rFileNames, fileData.Filename)
	}
	assert.ElementsMatch(t, testFiles, rFileNames)

	// Remove these two files and ensure that warning.txt file is generated
	err = os.Remove("notifications.log")
	require.NoError(t, err)
	err = os.Remove("mattermost.log")
	require.NoError(t, err)
	fileDatas = th.App.GenerateSupportPacket(th.Context)
	testFiles = []string{
		"support_packet.yaml",
		"plugins.json",
		"sanitized_config.json",
		"cpu.prof",
		"heap.prof",
		"warning.txt",
		"goroutines",
	}
	rFileNames = nil
	for _, fileData := range fileDatas {
		require.NotNil(t, fileData)
		assert.Positive(t, len(fileData.Body))

		rFileNames = append(rFileNames, fileData.Filename)
	}
	assert.ElementsMatch(t, testFiles, rFileNames)
}

func TestGetNotificationsLog(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	// Disable notifications file to get an error
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.NotificationLogSettings.EnableFile = false
	})

	fileData, err := th.App.getNotificationsLog(th.Context)
	assert.Nil(t, fileData)
	assert.ErrorContains(t, err, "Unable to retrieve notifications.log because LogSettings: EnableFile is set to false")

	// Enable notifications file but delete any notifications file to get an error trying to read the file
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.NotificationLogSettings.EnableFile = true
	})

	// If any previous notifications.log file, lets delete it
	os.Remove("notifications.log")

	fileData, err = th.App.getNotificationsLog(th.Context)
	assert.Nil(t, fileData)
	assert.ErrorContains(t, err, "failed read notifcation log file at path")

	// Happy path where we have file and no error
	d1 := []byte("hello\ngo\n")
	err = os.WriteFile("notifications.log", d1, 0777)
	defer os.Remove("notifications.log")
	require.NoError(t, err)

	fileData, err = th.App.getNotificationsLog(th.Context)
	require.NotNil(t, fileData)
	assert.Equal(t, "notifications.log", fileData.Filename)
	assert.Positive(t, len(fileData.Body))
	assert.NoError(t, err)
}

func TestGetMattermostLog(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	// disable mattermost log file setting in config so we should get an warning
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.LogSettings.EnableFile = false
	})

	fileData, err := th.App.getMattermostLog(th.Context)
	assert.Nil(t, fileData)
	assert.ErrorContains(t, err, "Unable to retrieve mattermost.log because LogSettings: EnableFile is set to false")

	// We enable the setting but delete any mattermost log file
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.LogSettings.EnableFile = true
	})

	// If any previous mattermost.log file, lets delete it
	os.Remove("mattermost.log")

	fileData, err = th.App.getMattermostLog(th.Context)
	assert.Nil(t, fileData)
	assert.ErrorContains(t, err, "failed read mattermost log file at path mattermost.log")

	// Happy path where we get a log file and no warning
	d1 := []byte("hello\ngo\n")
	err = os.WriteFile("mattermost.log", d1, 0777)
	defer os.Remove("mattermost.log")
	require.NoError(t, err)

	fileData, err = th.App.getMattermostLog(th.Context)
	require.NotNil(t, fileData)
	assert.Equal(t, "mattermost.log", fileData.Filename)
	assert.Positive(t, len(fileData.Body))
	assert.NoError(t, err)
}

func TestCreateSanitizedConfigFile(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	// Happy path where we have a sanitized config file with no err
	fileData, err := th.App.createSanitizedConfigFile(th.Context)
	require.NotNil(t, fileData)
	assert.Equal(t, "sanitized_config.json", fileData.Filename)
	assert.Positive(t, len(fileData.Body))
	assert.NoError(t, err)
}
