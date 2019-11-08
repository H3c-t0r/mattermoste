// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package plugin

import (
	"bytes"
	"github.com/blang/semver"
	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Helpers interface {
	// EnsureBot either returns an existing bot user matching the given bot, or creates a bot user from the given bot.
	// Returns the id of the resulting bot.
	//
	// Minimum server version: 5.10
	EnsureBot(bot *model.Bot) (string, error)

	// KVSetJSON stores a key-value pair, unique per plugin, marshalling the given value as a JSON string.
	//
	// Deprecated: Use p.API.KVSetWithOptions instead.
	//
	// Minimum server version: 5.2
	KVSetJSON(key string, value interface{}) error

	// KVCompareAndSetJSON updates a key-value pair, unique per plugin, but only if the current value matches the given oldValue after marshalling as a JSON string.
	// Inserts a new key if oldValue == nil.
	// Returns (false, err) if DB error occurred
	// Returns (false, nil) if current value != oldValue or key already exists when inserting
	// Returns (true, nil) if current value == oldValue or new key is inserted
	//
	// Deprecated: Use p.API.KVSetWithOptions instead.
	//
	// Minimum server version: 5.12
	KVCompareAndSetJSON(key string, oldValue interface{}, newValue interface{}) (bool, error)

	// KVCompareAndDeleteJSON deletes a key-value pair, unique per plugin, but only if the current value matches the given oldValue after marshalling as a JSON string.
	// Returns (false, err) if DB error occurred
	// Returns (false, nil) if current value != oldValue or the key was already deleted
	// Returns (true, nil) if current value == oldValue
	//
	// Minimum server version: 5.16
	KVCompareAndDeleteJSON(key string, oldValue interface{}) (bool, error)

	// KVGetJSON retrieves a value based on the key, unique per plugin, unmarshalling the previously set JSON string into the given value. Returns true if the key exists.
	//
	// Minimum server version: 5.2
	KVGetJSON(key string, value interface{}) (bool, error)

	// KVSetWithExpiryJSON stores a key-value pair with an expiry time, unique per plugin, marshalling the given value as a JSON string.
	//
	// Deprecated: Use p.API.KVSetWithOptions instead.
	//
	// Minimum server version: 5.6
	KVSetWithExpiryJSON(key string, value interface{}, expireInSeconds int64) error

	// Minimum server version: 5.18
	InstallPluginFromUrl(url string, replace bool) (*model.Manifest, *model.AppError)
}
type HelpersImpl struct {
	API API
}

func (p *HelpersImpl) ensureServerVersion(required string) error {
	serverVersion := p.API.GetServerVersion()
	currentVersion := semver.MustParse(serverVersion)
	requiredVersion := semver.MustParse(required)

	if currentVersion.LT(requiredVersion) {
		return errors.Errorf("incompatible server version for plugin, minimum required version: %s, current version: %s", required, serverVersion)
	}
	return nil
}

func (p *HelpersImpl) InstallPluginFromUrl(downloadUrl string, replace bool) (*model.Manifest, *model.AppError) {
	parsedUrl, _ := url.Parse(downloadUrl)
	if !*p.API.GetConfig().PluginSettings.AllowInsecureDownloadUrl && parsedUrl.Scheme != "https" {
		return nil, model.NewAppError("InstallPluginFromUrl", "api.plugin.install.insecure_url.app_error", nil, "", http.StatusBadRequest)
	}

	client := &http.Client{Timeout: 60 * time.Minute}
	response, err := client.Get(downloadUrl)
	if err != nil {
		return nil, model.NewAppError("InstallPluginFromUrlHelper", "api.plugin.install.download_failed.app_error", nil, err.Error(), http.StatusBadRequest)
	}
	defer response.Body.Close()

	fileBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, model.NewAppError("InstallPluginFromUrlHelper", "api.plugin.install.reading_stream_failed.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	manifest, appError := p.API.InstallPlugin(bytes.NewReader(fileBytes), true)
	if appError != nil {
		return nil, appError
	}

	return manifest, nil
}
