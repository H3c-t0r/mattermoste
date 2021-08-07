// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/mattermost/mattermost-server/v6/einterfaces/mocks"
	"github.com/mattermost/mattermost-server/v6/utils"
	mocks2 "github.com/mattermost/mattermost-server/v6/utils/mocks"
	"github.com/mattermost/mattermost-server/v6/utils/testutils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/v6/model"
)

func TestGetOldClientLicense(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()
	client := th.Client

	license, _, err = client.GetOldClientLicense("")
	require.NoError(t, err)

	require.NotEqual(t, license["IsLicensed"], "", "license not returned correctly")

	client.Logout()

	_, _, err = client.GetOldClientLicense("")
	require.NoError(t, err)

	_, err := client.DoApiGet("/license/client", "")
	require.NotNil(t, err, "get /license/client did not return an error")
	require.Equal(t, err.StatusCode, http.StatusNotImplemented,
		"expected 501 Not Implemented")

	_, err = client.DoApiGet("/license/client?format=junk", "")
	require.NotNil(t, err, "get /license/client?format=junk did not return an error")
	require.Equal(t, err.StatusCode, http.StatusBadRequest,
		"expected 400 Bad Request")

	license, _, err = th.SystemAdminClient.GetOldClientLicense("")
	require.NoError(t, err)

	require.NotEmpty(t, license["IsLicensed"], "license not returned correctly")
}

func TestUploadLicenseFile(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()
	client := th.Client
	LocalClient := th.LocalClient

	t.Run("as system user", func(t *testing.T) {
		ok, resp, err := client.UploadLicenseFile([]byte{})
		require.Error(t, err)
		CheckForbiddenStatus(t, resp)
		require.False(t, ok)
	})

	th.TestForSystemAdminAndLocal(t, func(t *testing.T, c *model.Client4) {
		ok, resp, _ := c.UploadLicenseFile([]byte{})
		CheckBadRequestStatus(t, resp)
		require.False(t, ok)
	}, "as system admin user")

	t.Run("as restricted system admin user", func(t *testing.T) {
		th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ExperimentalSettings.RestrictSystemAdmin = true })

		ok, resp, err := th.SystemAdminClient.UploadLicenseFile([]byte{})
		require.Error(t, err)
		CheckForbiddenStatus(t, resp)
		require.False(t, ok)
	})

	t.Run("restricted admin setting not honoured through local client", func(t *testing.T) {
		th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ExperimentalSettings.RestrictSystemAdmin = true })
		ok, resp, _ := LocalClient.UploadLicenseFile([]byte{})
		CheckBadRequestStatus(t, resp)
		require.False(t, ok)
	})

	t.Run("server has already gone through trial", func(t *testing.T) {
		th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ExperimentalSettings.RestrictSystemAdmin = false })
		mockLicenseValidator := mocks2.LicenseValidatorIface{}
		defer testutils.ResetLicenseValidator()

		//startTimestamp, err := time.Parse("2 Jan 2006 3:04 pm", "1 Jan 2021 12:00 am")
		//require.Nil(t, err)

		userCount := 100
		mills := model.GetMillis()

		license := model.License{
			Id: "AAAAAAAAAAAAAAAAAAAAAAAAAA",
			Features: &model.Features{
				Users: &userCount,
			},
			Customer: &model.Customer{
				Name: "Test",
			},
			StartsAt:  mills + 100,
			ExpiresAt: mills + 100 + (30*(time.Hour*24) + (time.Hour * 8)).Milliseconds(),
		}

		mockLicenseValidator.On("LicenseFromBytes", mock.Anything).Return(&license, nil).Once()
		licenseBytes, _ := json.Marshal(license)
		licenseStr := string(licenseBytes)

		mockLicenseValidator.On("ValidateLicense", mock.Anything).Return(true, licenseStr)
		utils.LicenseValidator = &mockLicenseValidator

		licenseManagerMock := &mocks.LicenseInterface{}
		licenseManagerMock.On("CanStartTrial").Return(false, nil).Once()
		th.App.Srv().LicenseManager = licenseManagerMock

		ok, resp, err := th.SystemAdminClient.UploadLicenseFile([]byte("sadasdasdasdasdasdsa"))
		require.False(t, ok)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		CheckErrorID(t, err, "api.license.request-trial.can-start-trial.not-allowed")
	})

	t.Run("allow uploading sanctioned trials even if server already gone through trial", func(t *testing.T) {
		mockLicenseValidator := mocks2.LicenseValidatorIface{}
		defer testutils.ResetLicenseValidator()

		userCount := 100
		mills := model.GetMillis()

		license := model.License{
			Id: "PPPPPPPPPPPPPPPPPPPPPPPPPP",
			Features: &model.Features{
				Users: &userCount,
			},
			Customer: &model.Customer{
				Name: "Test",
			},
			IsTrial:   true,
			StartsAt:  mills + 100,
			ExpiresAt: mills + 100 + (29*(time.Hour*24) + (time.Hour * 8)).Milliseconds(),
		}

		mockLicenseValidator.On("LicenseFromBytes", mock.Anything).Return(&license, nil).Once()

		licenseBytes, _ := json.Marshal(license)
		licenseStr := string(licenseBytes)

		mockLicenseValidator.On("ValidateLicense", mock.Anything).Return(true, licenseStr)

		utils.LicenseValidator = &mockLicenseValidator

		licenseManagerMock := &mocks.LicenseInterface{}
		licenseManagerMock.On("CanStartTrial").Return(false, nil).Once()
		th.App.Srv().LicenseManager = licenseManagerMock

		ok, resp, err := th.SystemAdminClient.UploadLicenseFile([]byte("sadasdasdasdasdasdsa"))
		require.False(t, ok)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NoError(t, err)
	})
}

func TestRemoveLicenseFile(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()
	client := th.Client
	LocalClient := th.LocalClient

	t.Run("as system user", func(t *testing.T) {
		ok, resp, err := client.RemoveLicenseFile()
		require.Error(t, err)
		CheckForbiddenStatus(t, resp)
		require.False(t, ok)
	})

	th.TestForSystemAdminAndLocal(t, func(t *testing.T, c *model.Client4) {
		ok, _, err = c.RemoveLicenseFile()
		require.NoError(t, err)
		require.True(t, ok)
	}, "as system admin user")

	t.Run("as restricted system admin user", func(t *testing.T) {
		th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ExperimentalSettings.RestrictSystemAdmin = true })

		ok, resp, err := th.SystemAdminClient.RemoveLicenseFile()
		require.Error(t, err)
		CheckForbiddenStatus(t, resp)
		require.False(t, ok)
	})

	t.Run("restricted admin setting not honoured through local client", func(t *testing.T) {
		th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ExperimentalSettings.RestrictSystemAdmin = true })

		ok, _, err = LocalClient.RemoveLicenseFile()
		require.NoError(t, err)
		require.True(t, ok)
	})
}

func TestRequestTrialLicense(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	licenseManagerMock := &mocks.LicenseInterface{}
	licenseManagerMock.On("CanStartTrial").Return(true, nil)
	th.App.Srv().LicenseManager = licenseManagerMock

	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.SiteURL = "http://localhost:8065/" })

	t.Run("permission denied", func(t *testing.T) {
		ok, resp, err := th.Client.RequestTrialLicense(1000)
		require.Error(t, err)
		CheckForbiddenStatus(t, resp)
		require.False(t, ok)
	})

	t.Run("blank site url", func(t *testing.T) {
		th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.SiteURL = "" })
		defer th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.SiteURL = "http://localhost:8065/" })
		ok, resp, err := th.SystemAdminClient.RequestTrialLicense(1000)
		CheckBadRequestStatus(t, resp)
		CheckErrorID(t, err, "api.license.request_trial_license.no-site-url.app_error")
		require.False(t, ok)
	})

	t.Run("trial license user count less than current users", func(t *testing.T) {
		ok, resp, err := th.SystemAdminClient.RequestTrialLicense(1)
		CheckBadRequestStatus(t, resp)
		CheckErrorID(t, err, "api.license.add_license.unique_users.app_error")
		require.False(t, ok)
	})

	th.App.Srv().LicenseManager = nil
	t.Run("trial license should fail if LicenseManager is nil", func(t *testing.T) {
		ok, resp, err := th.SystemAdminClient.RequestTrialLicense(1)
		CheckForbiddenStatus(t, resp)
		require.False(t, ok)
		CheckErrorID(t, err, "api.license.upgrade_needed.app_error")
	})
}
