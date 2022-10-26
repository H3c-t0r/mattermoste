// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package suite

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/product"
)

// Ensure preferences service wrapper implements `product.PreferencesService`
var _ product.PreferencesService = (*SuiteService)(nil)

func (s *SuiteService) UpdatePreferencesForUser(userID string, preferences model.Preferences) *model.AppError {
	return s.UpdatePreferences(userID, preferences)
}

func (s *SuiteService) DeletePreferencesForUser(userID string, preferences model.Preferences) *model.AppError {
	return s.DeletePreferences(userID, preferences)
}

func (s *SuiteService) GetPreferencesForUser(userID string) (model.Preferences, *model.AppError) {
	preferences, err := s.platform.Store.Preference().GetAll(userID)
	if err != nil {
		return nil, model.NewAppError("GetPreferencesForUser", "app.preference.get_all.app_error", nil, "", http.StatusBadRequest).Wrap(err)
	}
	return preferences, nil
}

func (s *SuiteService) GetPreferenceByCategoryForUser(userID string, category string) (model.Preferences, *model.AppError) {
	preferences, err := s.platform.Store.Preference().GetCategory(userID, category)
	if err != nil {
		return nil, model.NewAppError("GetPreferenceByCategoryForUser", "app.preference.get_category.app_error", nil, "", http.StatusBadRequest).Wrap(err)
	}
	if len(preferences) == 0 {
		err := model.NewAppError("GetPreferenceByCategoryForUser", "api.preference.preferences_category.get.app_error", nil, "", http.StatusNotFound)
		return nil, err
	}
	return preferences, nil
}

func (s *SuiteService) GetPreferenceByCategoryAndNameForUser(userID string, category string, preferenceName string) (*model.Preference, *model.AppError) {
	res, err := s.platform.Store.Preference().Get(userID, category, preferenceName)
	if err != nil {
		return nil, model.NewAppError("GetPreferenceByCategoryAndNameForUser", "app.preference.get.app_error", nil, "", http.StatusBadRequest).Wrap(err)
	}
	return res, nil
}

func (s *SuiteService) UpdatePreferences(userID string, preferences model.Preferences) *model.AppError {
	for _, preference := range preferences {
		if userID != preference.UserId {
			return model.NewAppError("savePreferences", "api.preference.update_preferences.set.app_error", nil,
				"userId="+userID+", preference.UserId="+preference.UserId, http.StatusForbidden)
		}
	}

	if err := s.platform.Store.Preference().Save(preferences); err != nil {
		var appErr *model.AppError
		switch {
		case errors.As(err, &appErr):
			return appErr
		default:
			return model.NewAppError("UpdatePreferences", "app.preference.save.updating.app_error", nil, "", http.StatusBadRequest).Wrap(err)
		}
	}

	if err := s.platform.Store.Channel().UpdateSidebarChannelsByPreferences(preferences); err != nil {
		return model.NewAppError("UpdatePreferences", "api.preference.update_preferences.update_sidebar.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	message := model.NewWebSocketEvent(model.WebsocketEventSidebarCategoryUpdated, "", "", userID, nil, "")
	// TODO this needs to be updated to include information on which categories changed
	s.platform.Publish(message)

	message = model.NewWebSocketEvent(model.WebsocketEventPreferencesChanged, "", "", userID, nil, "")
	prefsJSON, jsonErr := json.Marshal(preferences)
	if jsonErr != nil {
		return model.NewAppError("UpdatePreferences", "api.marshal_error", nil, "", http.StatusInternalServerError).Wrap(jsonErr)
	}
	message.Add("preferences", string(prefsJSON))
	s.platform.Publish(message)

	return nil
}

func (s *SuiteService) DeletePreferences(userID string, preferences model.Preferences) *model.AppError {
	for _, preference := range preferences {
		if userID != preference.UserId {
			err := model.NewAppError("DeletePreferences", "api.preference.delete_preferences.delete.app_error", nil,
				"userId="+userID+", preference.UserId="+preference.UserId, http.StatusForbidden)
			return err
		}
	}

	for _, preference := range preferences {
		if err := s.platform.Store.Preference().Delete(userID, preference.Category, preference.Name); err != nil {
			return model.NewAppError("DeletePreferences", "app.preference.delete.app_error", nil, "", http.StatusBadRequest).Wrap(err)
		}
	}

	if err := s.platform.Store.Channel().DeleteSidebarChannelsByPreferences(preferences); err != nil {
		return model.NewAppError("DeletePreferences", "api.preference.delete_preferences.update_sidebar.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	message := model.NewWebSocketEvent(model.WebsocketEventSidebarCategoryUpdated, "", "", userID, nil, "")
	// TODO this needs to be updated to include information on which categories changed
	s.platform.Publish(message)

	message = model.NewWebSocketEvent(model.WebsocketEventPreferencesDeleted, "", "", userID, nil, "")
	prefsJSON, jsonErr := json.Marshal(preferences)
	if jsonErr != nil {
		return model.NewAppError("DeletePreferences", "api.marshal_error", nil, "", http.StatusInternalServerError).Wrap(jsonErr)
	}
	message.Add("preferences", string(prefsJSON))
	s.platform.Publish(message)

	return nil
}
