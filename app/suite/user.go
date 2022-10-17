// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package suite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v6/app/email"
	"github.com/mattermost/mattermost-server/v6/app/imaging"
	"github.com/mattermost/mattermost-server/v6/app/request"
	"github.com/mattermost/mattermost-server/v6/app/users"
	"github.com/mattermost/mattermost-server/v6/einterfaces"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/mattermost/mattermost-server/v6/shared/i18n"
	"github.com/mattermost/mattermost-server/v6/shared/mfa"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"github.com/mattermost/mattermost-server/v6/store"
)

const (
	TokenTypePasswordRecovery  = "password_recovery"
	TokenTypeVerifyEmail       = "verify_email"
	TokenTypeTeamInvitation    = "team_invitation"
	TokenTypeGuestInvitation   = "guest_invitation"
	TokenTypeCWSAccess         = "cws_access_token"
	PasswordRecoverExpiryTime  = 1000 * 60 * 60 * 24 // 24 hours
	InvitationExpiryTime       = 1000 * 60 * 60 * 48 // 48 hours
	ImageProfilePixelDimension = 128
)

func (ss *SuiteService) CreateUserWithToken(c request.CTX, user *model.User, token *model.Token) (*model.User, *model.AppError) {
	if err := ss.IsUserSignUpAllowed(); err != nil {
		return nil, err
	}

	if token.Type != TokenTypeTeamInvitation && token.Type != TokenTypeGuestInvitation {
		return nil, model.NewAppError("CreateUserWithToken", "api.user.create_user.signup_link_invalid.app_error", nil, "", http.StatusBadRequest)
	}

	if model.GetMillis()-token.CreateAt >= InvitationExpiryTime {
		ss.DeleteToken(token)
		return nil, model.NewAppError("CreateUserWithToken", "api.user.create_user.signup_link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	tokenData := model.MapFromJSON(strings.NewReader(token.Extra))

	team, nErr := ss.platform.Store.Team().Get(tokenData["teamId"])
	if nErr != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(nErr, &nfErr):
			return nil, model.NewAppError("CreateUserWithToken", "app.team.get.find.app_error", nil, "", http.StatusNotFound).Wrap(nErr)
		default:
			return nil, model.NewAppError("CreateUserWithToken", "app.team.get.finding.app_error", nil, "", http.StatusInternalServerError).Wrap(nErr)
		}
	}

	emailFromToken := tokenData["email"]
	if emailFromToken != user.Email {
		return nil, model.NewAppError("CreateUserWithToken", "api.user.create_user.bad_token_email_data.app_error", nil, "", http.StatusBadRequest)
	}

	user.Email = tokenData["email"]
	user.EmailVerified = true

	var ruser *model.User
	var err *model.AppError
	if token.Type == TokenTypeTeamInvitation {
		ruser, err = ss.CreateUser(c, user)
	} else {
		ruser, err = ss.CreateGuest(c, user)
	}
	if err != nil {
		return nil, err
	}

	if _, err := ss.JoinUserToTeam(c, team, ruser, ""); err != nil {
		return nil, err
	}

	err = ss.channels.AddToDefaultChannelsWithToken(c, team.Id, ruser.Id, token)
	if err != nil {
		c.Logger().Warn("Failed to add channel member", mlog.Err(err))
	}

	if err := ss.DeleteToken(token); err != nil {
		c.Logger().Warn("Error while deleting token", mlog.Err(err))
	}

	return ruser, nil
}

func (a *SuiteService) CreateUserWithInviteId(c request.CTX, user *model.User, inviteId, redirect string) (*model.User, *model.AppError) {
	if err := a.IsUserSignUpAllowed(); err != nil {
		return nil, err
	}

	team, nErr := a.platform.Store.Team().GetByInviteId(inviteId)
	if nErr != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(nErr, &nfErr):
			return nil, model.NewAppError("CreateUserWithInviteId", "app.team.get_by_invite_id.finding.app_error", nil, "", http.StatusNotFound).Wrap(nErr)
		default:
			return nil, model.NewAppError("CreateUserWithInviteId", "app.team.get_by_invite_id.finding.app_error", nil, "", http.StatusInternalServerError).Wrap(nErr)
		}
	}

	if team.IsGroupConstrained() {
		return nil, model.NewAppError("CreateUserWithInviteId", "app.team.invite_id.group_constrained.error", nil, "", http.StatusForbidden)
	}

	if !users.CheckUserDomain(user, team.AllowedDomains) {
		return nil, model.NewAppError("CreateUserWithInviteId", "api.team.invite_members.invalid_email.app_error", map[string]any{"Addresses": team.AllowedDomains}, "", http.StatusForbidden)
	}

	user.EmailVerified = false

	ruser, err := a.CreateUser(c, user)
	if err != nil {
		return nil, err
	}

	if _, err := a.JoinUserToTeam(c, team, ruser, ""); err != nil {
		return nil, err
	}

	a.channels.AddDirectChannels(c, team.Id, ruser.Id)

	if err := a.email.SendWelcomeEmail(ruser.Id, ruser.Email, ruser.EmailVerified, ruser.DisableWelcomeEmail, ruser.Locale, a.GetSiteURL(), redirect); err != nil {
		c.Logger().Warn("Failed to send welcome email on create user with inviteId", mlog.Err(err))
	}

	return ruser, nil
}

func (a *SuiteService) CreateUserAsAdmin(c request.CTX, user *model.User, redirect string) (*model.User, *model.AppError) {
	ruser, err := a.CreateUser(c, user)
	if err != nil {
		return nil, err
	}

	if err := a.email.SendWelcomeEmail(ruser.Id, ruser.Email, ruser.EmailVerified, ruser.DisableWelcomeEmail, ruser.Locale, a.GetSiteURL(), redirect); err != nil {
		c.Logger().Warn("Failed to send welcome email to the new user, created by system admin", mlog.Err(err))
	}

	return ruser, nil
}

func (a *SuiteService) CreateUserFromSignup(c request.CTX, user *model.User, redirect string) (*model.User, *model.AppError) {
	if err := a.IsUserSignUpAllowed(); err != nil {
		return nil, err
	}

	if !a.IsFirstUserAccount() && !*a.platform.Config().TeamSettings.EnableOpenServer {
		err := model.NewAppError("CreateUserFromSignup", "api.user.create_user.no_open_server", nil, "email="+user.Email, http.StatusForbidden)
		return nil, err
	}

	user.EmailVerified = false

	ruser, err := a.CreateUser(c, user)
	if err != nil {
		return nil, err
	}

	if err := a.email.SendWelcomeEmail(ruser.Id, ruser.Email, ruser.EmailVerified, ruser.DisableWelcomeEmail, ruser.Locale, a.GetSiteURL(), redirect); err != nil {
		c.Logger().Warn("Failed to send welcome email on create user from signup", mlog.Err(err))
	}

	return ruser, nil
}

func (a *SuiteService) IsUserSignUpAllowed() *model.AppError {
	if !*a.platform.Config().EmailSettings.EnableSignUpWithEmail || !*a.platform.Config().TeamSettings.EnableUserCreation {
		err := model.NewAppError("IsUserSignUpAllowed", "api.user.create_user.signup_email_disabled.app_error", nil, "", http.StatusNotImplemented)
		return err
	}
	return nil
}

func (a *SuiteService) IsFirstUserAccount() bool {
	return a.platform.IsFirstUserAccount()
}

// CreateUser creates a user and sets several fields of the returned User struct to
// their zero values.
func (a *SuiteService) CreateUser(c request.CTX, user *model.User) (*model.User, *model.AppError) {
	return a.createUserOrGuest(c, user, false)
}

// CreateGuest creates a guest and sets several fields of the returned User struct to
// their zero values.
func (a *SuiteService) CreateGuest(c request.CTX, user *model.User) (*model.User, *model.AppError) {
	return a.createUserOrGuest(c, user, true)
}

func (a *SuiteService) createUserOrGuest(c request.CTX, user *model.User, guest bool) (*model.User, *model.AppError) {
	if err := a.isUniqueToGroupNames(user.Username); err != nil {
		err.Where = "createUserOrGuest"
		return nil, err
	}

	ruser, nErr := a.createUserWithOptions(user, UserCreateOptions{Guest: guest})
	if nErr != nil {
		var appErr *model.AppError
		var invErr *store.ErrInvalidInput
		var nfErr *users.ErrInvalidPassword
		switch {
		case errors.As(nErr, &appErr):
			return nil, appErr
		case errors.Is(nErr, users.AcceptedDomainError):
			return nil, model.NewAppError("createUserOrGuest", "api.user.create_user.accepted_domain.app_error", nil, "", http.StatusBadRequest).Wrap(nErr)
		case errors.As(nErr, &nfErr):
			return nil, model.NewAppError("createUserOrGuest", "api.user.check_user_password.invalid.app_error", nil, "", http.StatusBadRequest).Wrap(nErr)
		case errors.Is(nErr, users.UserStoreIsEmptyError):
			return nil, model.NewAppError("createUserOrGuest", "app.user.store_is_empty.app_error", nil, "", http.StatusInternalServerError).Wrap(nErr)
		case errors.As(nErr, &invErr):
			switch invErr.Field {
			case "email":
				return nil, model.NewAppError("createUserOrGuest", "app.user.save.email_exists.app_error", nil, "", http.StatusBadRequest).Wrap(nErr)
			case "username":
				return nil, model.NewAppError("createUserOrGuest", "app.user.save.username_exists.app_error", nil, "", http.StatusBadRequest).Wrap(nErr)
			default:
				return nil, model.NewAppError("createUserOrGuest", "app.user.save.existing.app_error", nil, "", http.StatusBadRequest).Wrap(nErr)
			}
		default:
			return nil, model.NewAppError("createUserOrGuest", "app.user.save.app_error", nil, "", http.StatusInternalServerError).Wrap(nErr)
		}
	}

	if user.EmailVerified {
		a.InvalidateCacheForUser(ruser.Id)

		nUser, err := a.getUser(ruser.Id)
		if err != nil {
			var nfErr *store.ErrNotFound
			switch {
			case errors.As(err, &nfErr):
				return nil, model.NewAppError("createUserOrGuest", MissingAccountError, nil, "", http.StatusNotFound).Wrap(err)
			default:
				return nil, model.NewAppError("createUserOrGuest", "app.user.get.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
			}
		}

		a.sendUpdatedUserEvent(*nUser)
	}

	recommendedNextStepsPref := model.Preference{UserId: ruser.Id, Category: model.PreferenceRecommendedNextSteps, Name: "hide", Value: "false"}
	tutorialStepPref := model.Preference{UserId: ruser.Id, Category: model.PreferenceCategoryTutorialSteps, Name: ruser.Id, Value: "0"}

	preferences := model.Preferences{recommendedNextStepsPref, tutorialStepPref}

	if a.platform.Config().FeatureFlags.InsightsEnabled {
		// We don't want to show the insights intro modal for new users
		preferences = append(preferences, model.Preference{UserId: ruser.Id, Category: model.PreferenceCategoryInsights, Name: model.PreferenceNameInsights, Value: "{\"insights_modal_viewed\":true}"})
	} else {
		preferences = append(preferences, model.Preference{UserId: ruser.Id, Category: model.PreferenceCategoryInsights, Name: model.PreferenceNameInsights, Value: "{\"insights_modal_viewed\":false}"})
	}

	if err := a.platform.Store.Preference().Save(preferences); err != nil {
		c.Logger().Warn("Encountered error saving user preferences", mlog.Err(err))
	}

	go a.UpdateViewedProductNoticesForNewUser(ruser.Id)

	// This message goes to everyone, so the teamID, channelID and userID are irrelevant
	message := model.NewWebSocketEvent(model.WebsocketEventNewUser, "", "", "", nil, "")
	message.Add("user_id", ruser.Id)
	a.platform.Publish(message)

	if pluginsEnvironment := a.GetPluginsEnvironment(); pluginsEnvironment != nil {
		a.platform.Go(func() {
			pluginContext := pluginContext(c)
			pluginsEnvironment.RunMultiPluginHook(func(hooks plugin.Hooks) bool {
				hooks.UserHasBeenCreated(pluginContext, ruser)
				return true
			}, plugin.UserHasBeenCreatedID)
		})
	}

	return ruser, nil
}

func (a *SuiteService) CreateOAuthUser(c *request.Context, service string, userData io.Reader, teamID string, tokenUser *model.User) (*model.User, *model.AppError) {
	if !*a.platform.Config().TeamSettings.EnableUserCreation {
		return nil, model.NewAppError("CreateOAuthUser", "api.user.create_user.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	provider, e := a.getSSOProvider(service)
	if e != nil {
		return nil, e
	}
	user, err1 := provider.GetUserFromJSON(userData, tokenUser)
	if err1 != nil {
		return nil, model.NewAppError("CreateOAuthUser", "api.user.create_oauth_user.create.app_error", map[string]any{"Service": service}, "", http.StatusInternalServerError).Wrap(err1)
	}
	if user.AuthService == "" {
		user.AuthService = service
	}

	found := true
	count := 0
	for found {
		if found = a.IsUsernameTaken(user.Username); found {
			user.Username = user.Username + strconv.Itoa(count)
			count++
		}
	}

	userByAuth, _ := a.getUserByAuth(user.AuthData, service)
	if userByAuth != nil {
		return userByAuth, nil
	}

	userByEmail, _ := a.getUserByEmail(user.Email)
	if userByEmail != nil {
		if userByEmail.AuthService == "" {
			return nil, model.NewAppError("CreateOAuthUser", "api.user.create_oauth_user.already_attached.app_error", map[string]any{"Service": service, "Auth": model.UserAuthServiceEmail}, "email="+user.Email, http.StatusBadRequest)
		}
		if provider.IsSameUser(userByEmail, user) {
			if _, err := a.platform.Store.User().UpdateAuthData(userByEmail.Id, user.AuthService, user.AuthData, "", false); err != nil {
				// if the user is not updated, write a warning to the log, but don't prevent user login
				c.Logger().Warn("Error attempting to update user AuthData", mlog.Err(err))
			}
			return userByEmail, nil
		}
		return nil, model.NewAppError("CreateOAuthUser", "api.user.create_oauth_user.already_attached.app_error", map[string]any{"Service": service, "Auth": userByEmail.AuthService}, "email="+user.Email+" authData="+*user.AuthData, http.StatusBadRequest)
	}

	user.EmailVerified = true

	ruser, err := a.CreateUser(c, user)
	if err != nil {
		return nil, err
	}

	if teamID != "" {
		err = a.AddUserToTeamByTeamId(c, teamID, user)
		if err != nil {
			return nil, err
		}

		err = a.channels.AddDirectChannels(c, teamID, user.Id)
		if err != nil {
			c.Logger().Warn("Failed to add direct channels", mlog.Err(err))
		}
	}

	return ruser, nil
}

func (a *SuiteService) GetUser(userID string) (*model.User, *model.AppError) {
	user, err := a.getUser(userID)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetUser", MissingAccountError, nil, "", http.StatusNotFound).Wrap(err)
		default:
			return nil, model.NewAppError("GetUser", "app.user.get_by_username.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	return user, nil
}

func (a *SuiteService) GetUsers(userIDs []string) ([]*model.User, *model.AppError) {
	users, err := a.getUsers(userIDs)
	if err != nil {
		return nil, model.NewAppError("GetUsers", "app.user.get.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

func (a *SuiteService) GetUserByUsername(username string) (*model.User, *model.AppError) {
	result, err := a.getUserByUsername(username)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetUserByUsername", "app.user.get_by_username.app_error", nil, "", http.StatusNotFound).Wrap(err)
		default:
			return nil, model.NewAppError("GetUserByUsername", "app.user.get_by_username.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}
	return result, nil
}

func (a *SuiteService) GetUserByEmail(email string) (*model.User, *model.AppError) {
	user, err := a.getUserByEmail(email)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetUserByEmail", MissingAccountError, nil, "", http.StatusNotFound).Wrap(err)
		default:
			return nil, model.NewAppError("GetUserByEmail", MissingAccountError, nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}
	return user, nil
}

func (a *SuiteService) GetUserByAuth(authData *string, authService string) (*model.User, *model.AppError) {
	user, err := a.getUserByAuth(authData, authService)
	if err != nil {
		var invErr *store.ErrInvalidInput
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &invErr):
			return nil, model.NewAppError("GetUserByAuth", MissingAuthAccountError, nil, "", http.StatusBadRequest).Wrap(err)
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetUserByAuth", MissingAuthAccountError, nil, "", http.StatusInternalServerError).Wrap(err)
		default:
			return nil, model.NewAppError("GetUserByAuth", "app.user.get_by_auth.other.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	return user, nil
}

func (a *SuiteService) GetUsersFromProfiles(options *model.UserGetOptions) ([]*model.User, *model.AppError) {
	users, err := a.getUsersFromProfiles(options)
	if err != nil {
		return nil, model.NewAppError("GetUsers", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

func (a *SuiteService) GetUsersPage(options *model.UserGetOptions, asAdmin bool) ([]*model.User, *model.AppError) {
	users, err := a.getUsersPage(options, asAdmin)
	if err != nil {
		return nil, model.NewAppError("GetUsersPage", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

func (a *SuiteService) GetUsersEtag(restrictionsHash string) string {
	return a.getUsersEtag(restrictionsHash)
}

func (a *SuiteService) GetUsersInTeam(options *model.UserGetOptions) ([]*model.User, *model.AppError) {
	users, err := a.getUsersInTeam(options)
	if err != nil {
		return nil, model.NewAppError("GetUsersInTeam", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

func (a *SuiteService) GetUsersNotInTeam(teamID string, groupConstrained bool, offset int, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError) {
	users, err := a.getUsersNotInTeam(teamID, groupConstrained, offset, limit, viewRestrictions)
	if err != nil {
		return nil, model.NewAppError("GetUsersNotInTeam", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

func (a *SuiteService) GetUsersInTeamPage(options *model.UserGetOptions, asAdmin bool) ([]*model.User, *model.AppError) {
	users, err := a.getUsersInTeamPage(options, asAdmin)
	if err != nil {
		return nil, model.NewAppError("GetUsersInTeamPage", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return a.sanitizeProfiles(users, asAdmin), nil
}

func (a *SuiteService) GetUsersNotInTeamPage(teamID string, groupConstrained bool, page int, perPage int, asAdmin bool, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError) {
	users, err := a.getUsersNotInTeamPage(teamID, groupConstrained, page*perPage, perPage, asAdmin, viewRestrictions)
	if err != nil {
		return nil, model.NewAppError("GetUsersNotInTeamPage", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return a.sanitizeProfiles(users, asAdmin), nil
}

func (a *SuiteService) GetUsersInTeamEtag(teamID string, restrictionsHash string) string {
	return a.getUsersInTeamEtag(teamID, restrictionsHash)
}

func (a *SuiteService) GetUsersNotInTeamEtag(teamID string, restrictionsHash string) string {
	return a.getUsersNotInTeamEtag(teamID, restrictionsHash)
}

func (a *SuiteService) GetUsersInChannel(options *model.UserGetOptions) ([]*model.User, *model.AppError) {
	users, err := a.platform.Store.User().GetProfilesInChannel(options)
	if err != nil {
		return nil, model.NewAppError("GetUsersInChannel", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

func (a *SuiteService) GetUsersInChannelByStatus(options *model.UserGetOptions) ([]*model.User, *model.AppError) {
	users, err := a.platform.Store.User().GetProfilesInChannelByStatus(options)
	if err != nil {
		return nil, model.NewAppError("GetUsersInChannelByStatus", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

func (a *SuiteService) GetUsersInChannelByAdmin(options *model.UserGetOptions) ([]*model.User, *model.AppError) {
	users, err := a.platform.Store.User().GetProfilesInChannelByAdmin(options)
	if err != nil {
		return nil, model.NewAppError("GetUsersInChannelByAdmin", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

func (a *SuiteService) GetUsersInChannelMap(options *model.UserGetOptions, asAdmin bool) (map[string]*model.User, *model.AppError) {
	users, err := a.GetUsersInChannel(options)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*model.User, len(users))

	for _, user := range users {
		a.SanitizeProfile(user, asAdmin)
		userMap[user.Id] = user
	}

	return userMap, nil
}

func (a *SuiteService) GetUsersInChannelPage(options *model.UserGetOptions, asAdmin bool) ([]*model.User, *model.AppError) {
	users, err := a.GetUsersInChannel(options)
	if err != nil {
		return nil, err
	}
	return a.sanitizeProfiles(users, asAdmin), nil
}

func (a *SuiteService) GetUsersInChannelPageByStatus(options *model.UserGetOptions, asAdmin bool) ([]*model.User, *model.AppError) {
	users, err := a.GetUsersInChannelByStatus(options)
	if err != nil {
		return nil, err
	}
	return a.sanitizeProfiles(users, asAdmin), nil
}

func (a *SuiteService) GetUsersInChannelPageByAdmin(options *model.UserGetOptions, asAdmin bool) ([]*model.User, *model.AppError) {
	users, err := a.GetUsersInChannelByAdmin(options)
	if err != nil {
		return nil, err
	}
	return a.sanitizeProfiles(users, asAdmin), nil
}

func (a *SuiteService) GetUsersNotInChannel(teamID string, channelID string, groupConstrained bool, offset int, limit int, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError) {
	users, err := a.platform.Store.User().GetProfilesNotInChannel(teamID, channelID, groupConstrained, offset, limit, viewRestrictions)
	if err != nil {
		return nil, model.NewAppError("GetUsersNotInChannel", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

func (a *SuiteService) GetUsersNotInChannelMap(teamID string, channelID string, groupConstrained bool, offset int, limit int, asAdmin bool, viewRestrictions *model.ViewUsersRestrictions) (map[string]*model.User, *model.AppError) {
	users, err := a.GetUsersNotInChannel(teamID, channelID, groupConstrained, offset, limit, viewRestrictions)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*model.User, len(users))

	for _, user := range users {
		a.SanitizeProfile(user, asAdmin)
		userMap[user.Id] = user
	}

	return userMap, nil
}

func (a *SuiteService) GetUsersNotInChannelPage(teamID string, channelID string, groupConstrained bool, page int, perPage int, asAdmin bool, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError) {
	users, err := a.GetUsersNotInChannel(teamID, channelID, groupConstrained, page*perPage, perPage, viewRestrictions)
	if err != nil {
		return nil, err
	}

	return a.sanitizeProfiles(users, asAdmin), nil
}

func (a *SuiteService) GetUsersWithoutTeamPage(options *model.UserGetOptions, asAdmin bool) ([]*model.User, *model.AppError) {
	users, err := a.getUsersWithoutTeamPage(options, asAdmin)
	if err != nil {
		return nil, model.NewAppError("GetUsersWithoutTeamPage", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return a.sanitizeProfiles(users, asAdmin), nil
}

func (a *SuiteService) GetUsersWithoutTeam(options *model.UserGetOptions) ([]*model.User, *model.AppError) {
	users, err := a.getUsersWithoutTeam(options)
	if err != nil {
		return nil, model.NewAppError("GetUsersWithoutTeam", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

// GetTeamGroupUsers returns the users who are associated to the team via GroupTeams and GroupMembers.
func (a *SuiteService) GetTeamGroupUsers(teamID string) ([]*model.User, *model.AppError) {
	users, err := a.platform.Store.User().GetTeamGroupUsers(teamID)
	if err != nil {
		return nil, model.NewAppError("GetTeamGroupUsers", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

// GetChannelGroupUsers returns the users who are associated to the channel via GroupChannels and GroupMembers.
func (a *SuiteService) GetChannelGroupUsers(channelID string) ([]*model.User, *model.AppError) {
	users, err := a.platform.Store.User().GetChannelGroupUsers(channelID)
	if err != nil {
		return nil, model.NewAppError("GetChannelGroupUsers", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

func (a *SuiteService) GetUsersByIds(userIDs []string, options *store.UserGetByIdsOpts) ([]*model.User, *model.AppError) {
	users, err := a.getUsersByIds(userIDs, options)
	if err != nil {
		return nil, model.NewAppError("GetUsersByIds", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

func (a *SuiteService) GetUsersByGroupChannelIds(c *request.Context, channelIDs []string, asAdmin bool) (map[string][]*model.User, *model.AppError) {
	usersByChannelId, err := a.platform.Store.User().GetProfileByGroupChannelIdsForUser(c.Session().UserId, channelIDs)
	if err != nil {
		return nil, model.NewAppError("GetUsersByGroupChannelIds", "app.user.get_profile_by_group_channel_ids_for_user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	for channelID, userList := range usersByChannelId {
		usersByChannelId[channelID] = a.sanitizeProfiles(userList, asAdmin)
	}

	return usersByChannelId, nil
}

func (a *SuiteService) GetUsersByUsernames(usernames []string, asAdmin bool, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError) {
	users, err := a.getUsersByUsernames(usernames, &model.UserGetOptions{ViewRestrictions: viewRestrictions})
	if err != nil {
		return nil, model.NewAppError("GetUsersByUsernames", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	return a.sanitizeProfiles(users, asAdmin), nil
}

func (a *SuiteService) sanitizeProfiles(users []*model.User, asAdmin bool) []*model.User {
	for _, u := range users {
		a.SanitizeProfile(u, asAdmin)
	}

	return users
}

func (a *SuiteService) GenerateMfaSecret(userID string) (*model.MfaSecret, *model.AppError) {
	user, appErr := a.GetUser(userID)
	if appErr != nil {
		return nil, appErr
	}

	if !*a.platform.Config().ServiceSettings.EnableMultifactorAuthentication {
		return nil, model.NewAppError("GenerateMfaSecret", "mfa.mfa_disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	mfaSecret, err := a.generateMfaSecret(user)
	if err != nil {
		return nil, model.NewAppError("GenerateMfaSecret", "mfa.generate_qr_code.create_code.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return mfaSecret, nil
}

func (a *SuiteService) ActivateMfa(userID, token string) *model.AppError {
	user, appErr := a.GetUser(userID)
	if appErr != nil {
		return appErr
	}

	if user.AuthService != "" && user.AuthService != model.UserAuthServiceLdap {
		return model.NewAppError("ActivateMfa", "api.user.activate_mfa.email_and_ldap_only.app_error", nil, "", http.StatusBadRequest)
	}

	if !*a.platform.Config().ServiceSettings.EnableMultifactorAuthentication {
		return model.NewAppError("ActivateMfa", "mfa.mfa_disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	if err := a.activateMfa(user, token); err != nil {
		switch {
		case errors.Is(err, mfa.InvalidToken):
			return model.NewAppError("ActivateMfa", "mfa.activate.bad_token.app_error", nil, "", http.StatusUnauthorized)
		default:
			return model.NewAppError("ActivateMfa", "mfa.activate.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	// Make sure old MFA status is not cached locally or in cluster nodes.
	a.InvalidateCacheForUser(userID)

	return nil
}

func (a *SuiteService) DeactivateMfa(userID string) *model.AppError {
	user, appErr := a.GetUser(userID)
	if appErr != nil {
		return appErr
	}

	if err := a.deactivateMfa(user); err != nil {
		return model.NewAppError("DeactivateMfa", "mfa.deactivate.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	// Make sure old MFA status is not cached locally or in cluster nodes.
	a.InvalidateCacheForUser(userID)

	return nil
}

func (a *SuiteService) SetDefaultProfileImage(c request.CTX, user *model.User) *model.AppError {
	img, err := a.GetDefaultProfileImage(user)
	if err != nil {
		switch {
		case errors.Is(err, users.DefaultFontError):
			return model.NewAppError("SetDefaultProfileImage", "api.user.create_profile_image.default_font.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		case errors.Is(err, users.UserInitialsError):
			return model.NewAppError("SetDefaultProfileImage", "api.user.create_profile_image.initial.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		default:
			return model.NewAppError("SetDefaultProfileImage", "api.user.create_profile_image.encode.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	path := getProfileImagePath(user.Id)
	if _, err := a.platform.FileBackend().WriteFile(bytes.NewReader(img), path); err != nil {
		return model.NewAppError("WriteFile", "api.file.write_file.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.User().ResetLastPictureUpdate(user.Id); err != nil {
		c.Logger().Warn("Failed to reset last picture update", mlog.Err(err))
	}

	a.InvalidateCacheForUser(user.Id)

	updatedUser, appErr := a.GetUser(user.Id)
	if appErr != nil {
		c.Logger().Warn("Error in getting users profile forcing logout", mlog.String("user_id", user.Id), mlog.Err(appErr))
		return nil
	}

	options := a.platform.Config().GetSanitizeOptions()
	updatedUser.SanitizeProfile(options)

	message := model.NewWebSocketEvent(model.WebsocketEventUserUpdated, "", "", "", nil, "")
	message.Add("user", updatedUser)
	a.platform.Publish(message)

	return nil
}

func (a *SuiteService) SetProfileImage(c request.CTX, userID string, imageData *multipart.FileHeader) *model.AppError {
	file, err := imageData.Open()
	if err != nil {
		return model.NewAppError("SetProfileImage", "api.user.upload_profile_user.open.app_error", nil, "", http.StatusBadRequest).Wrap(err)
	}
	defer file.Close()
	return a.SetProfileImageFromMultiPartFile(c, userID, file)
}

func (a *SuiteService) SetProfileImageFromMultiPartFile(c request.CTX, userID string, file multipart.File) *model.AppError {
	if limitErr := checkImageLimits(file, *a.platform.Config().FileSettings.MaxImageResolution); limitErr != nil {
		return model.NewAppError("SetProfileImage", "api.user.upload_profile_user.check_image_limits.app_error", nil, "", http.StatusBadRequest)
	}

	return a.SetProfileImageFromFile(c, userID, file)
}

func (a *SuiteService) AdjustImage(file io.Reader) (*bytes.Buffer, *model.AppError) {
	// Decode image into Image object
	img, _, err := a.imgDecoder.Decode(file)
	if err != nil {
		return nil, model.NewAppError("SetProfileImage", "api.user.upload_profile_user.decode.app_error", nil, "", http.StatusBadRequest).Wrap(err)
	}

	orientation, _ := imaging.GetImageOrientation(file)
	img = imaging.MakeImageUpright(img, orientation)

	// Scale profile image
	profileWidthAndHeight := 128
	img = imaging.FillCenter(img, profileWidthAndHeight, profileWidthAndHeight)

	buf := new(bytes.Buffer)
	err = a.imgEncoder.EncodePNG(buf, img)
	if err != nil {
		return nil, model.NewAppError("SetProfileImage", "api.user.upload_profile_user.encode.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	return buf, nil
}

func (a *SuiteService) SetProfileImageFromFile(c request.CTX, userID string, file io.Reader) *model.AppError {
	buf, err := a.AdjustImage(file)
	if err != nil {
		return err
	}

	path := getProfileImagePath(userID)
	if storedData, err := a.ReadFile(path); err == nil && bytes.Equal(storedData, buf.Bytes()) {
		return nil
	}

	if _, err := a.platform.FileBackend().WriteFile(buf, path); err != nil {
		return model.NewAppError("SetProfileImage", "api.user.upload_profile_user.upload_profile.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.User().UpdateLastPictureUpdate(userID); err != nil {
		c.Logger().Warn("Error with updating last picture update", mlog.Err(err))
	}
	a.invalidateUserCacheAndPublish(userID)
	a.onUserProfileChange(userID)

	return nil
}

func (a *SuiteService) UpdatePasswordAsUser(c request.CTX, userID, currentPassword, newPassword string) *model.AppError {
	user, err := a.GetUser(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return model.NewAppError("updatePassword", "api.user.update_password.valid_account.app_error", nil, "", http.StatusBadRequest)
	}

	if user.AuthData != nil && *user.AuthData != "" {
		return model.NewAppError("updatePassword", "api.user.update_password.oauth.app_error", nil, "auth_service="+user.AuthService, http.StatusBadRequest)
	}

	if err := a.DoubleCheckPassword(user, currentPassword); err != nil {
		if err.Id == "api.user.check_user_password.invalid.app_error" {
			err = model.NewAppError("updatePassword", "api.user.update_password.incorrect.app_error", nil, "", http.StatusBadRequest)
		}
		return err
	}

	T := i18n.GetUserTranslations(user.Locale)

	return a.UpdatePasswordSendEmail(c, user, newPassword, T("api.user.update_password.menu"))
}

func (a *SuiteService) userDeactivated(c request.CTX, userID string) *model.AppError {
	a.SetStatusOffline(userID, false)

	user, err := a.GetUser(userID)
	if err != nil {
		return err
	}

	// when disable a user, userDeactivated is called for the user and the
	// bots the user owns. Only notify once, when the user is the owner, not the
	// owners bots
	if !user.IsBot {
		a.notifySysadminsBotOwnerDeactivated(c, userID)
	}

	if *a.platform.Config().ServiceSettings.DisableBotsWhenOwnerIsDeactivated {
		a.disableUserBots(c, userID)
	}

	return nil
}

func (a *SuiteService) invalidateUserChannelMembersCaches(c request.CTX, userID string) *model.AppError {
	teamsForUser, err := a.GetTeamsForUser(userID)
	if err != nil {
		return err
	}

	for _, team := range teamsForUser {
		channelsForUser, err := a.channels.GetChannelsForTeamForUser(c, team.Id, userID, &model.ChannelSearchOpts{
			IncludeDeleted: false,
			LastDeleteAt:   0,
		})
		if err != nil {
			return err
		}

		for _, channel := range channelsForUser {
			a.platform.InvalidateCacheForChannelMembers(channel.Id)
		}
	}

	return nil
}

func (a *SuiteService) UpdateActive(c request.CTX, user *model.User, active bool) (*model.User, *model.AppError) {
	user.UpdateAt = model.GetMillis()
	if active {
		user.DeleteAt = 0
	} else {
		user.DeleteAt = user.UpdateAt
	}

	userUpdate, err := a.updateUser(user, true)
	if err != nil {
		var appErr *model.AppError
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(err, &appErr):
			return nil, appErr
		case errors.As(err, &invErr):
			return nil, model.NewAppError("UpdateActive", "app.user.update.find.app_error", nil, "", http.StatusBadRequest).Wrap(err)
		default:
			return nil, model.NewAppError("UpdateActive", "app.user.update.finding.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}
	ruser := userUpdate.New

	if !active {
		if err := a.RevokeAllSessions(ruser.Id); err != nil {
			return nil, err
		}
		if err := a.userDeactivated(c, ruser.Id); err != nil {
			return nil, err
		}
	}

	a.invalidateUserChannelMembersCaches(c, user.Id)
	a.InvalidateCacheForUser(user.Id)

	a.sendUpdatedUserEvent(*ruser)

	return ruser, nil
}

func (a *SuiteService) DeactivateGuests(c *request.Context) *model.AppError {
	userIDs, err := a.deactivateAllGuests()
	if err != nil {
		return model.NewAppError("DeactivateGuests", "app.user.update_active_for_multiple_users.updating.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	for _, userID := range userIDs {
		if err := a.platform.RevokeAllSessions(userID); err != nil {
			return model.NewAppError("DeactivateGuests", "app.user.update_active_for_multiple_users.updating.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	for _, userID := range userIDs {
		if err := a.userDeactivated(c, userID); err != nil {
			return err
		}
	}

	a.platform.Store.Channel().ClearCaches()
	a.platform.Store.User().ClearCaches()

	message := model.NewWebSocketEvent(model.WebsocketEventGuestsDeactivated, "", "", "", nil, "")
	a.platform.Publish(message)

	return nil
}

func (a *SuiteService) GetSanitizeOptions(asAdmin bool) map[string]bool {
	return a.getUserSanitizeOptions(asAdmin)
}

func (a *SuiteService) SanitizeProfile(user *model.User, asAdmin bool) {
	options := a.getUserSanitizeOptions(asAdmin)

	user.SanitizeProfile(options)
}

func (a *SuiteService) UpdateUserAsUser(c request.CTX, user *model.User, asAdmin bool) (*model.User, *model.AppError) {
	updatedUser, err := a.UpdateUser(c, user, true)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// CheckProviderAttributes returns the empty string if the patch can be applied without
// overriding attributes set by the user's login provider; otherwise, the name of the offending
// field is returned.
func (a *SuiteService) CheckProviderAttributes(user *model.User, patch *model.UserPatch) string {
	tryingToChange := func(userValue *string, patchValue *string) bool {
		return patchValue != nil && *patchValue != *userValue
	}

	// If any login provider is used, then the username may not be changed
	if user.AuthService != "" && tryingToChange(&user.Username, patch.Username) {
		return "username"
	}

	LdapSettings := &a.platform.Config().LdapSettings
	SamlSettings := &a.platform.Config().SamlSettings

	conflictField := ""
	if a.ldap != nil &&
		(user.IsLDAPUser() || (user.IsSAMLUser() && *SamlSettings.EnableSyncWithLdap)) {
		conflictField = a.ldap.CheckProviderAttributes(LdapSettings, user, patch)
	} else if a.saml != nil && user.IsSAMLUser() {
		conflictField = a.saml.CheckProviderAttributes(SamlSettings, user, patch)
	} else if user.IsOAuthUser() {
		if tryingToChange(&user.FirstName, patch.FirstName) || tryingToChange(&user.LastName, patch.LastName) {
			conflictField = "full name"
		}
	}

	return conflictField
}

func (a *SuiteService) PatchUser(c request.CTX, userID string, patch *model.UserPatch, asAdmin bool) (*model.User, *model.AppError) {
	user, err := a.GetUser(userID)
	if err != nil {
		return nil, err
	}

	user.Patch(patch)

	updatedUser, err := a.UpdateUser(c, user, true)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (a *SuiteService) UpdateUserAuth(userID string, userAuth *model.UserAuth) (*model.UserAuth, *model.AppError) {
	userAuth.Password = ""
	if _, err := a.platform.Store.User().UpdateAuthData(userID, userAuth.AuthService, userAuth.AuthData, "", false); err != nil {
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(err, &invErr):
			return nil, model.NewAppError("UpdateUserAuth", "app.user.update_auth_data.email_exists.app_error", nil, "", http.StatusBadRequest).Wrap(err)
		default:
			return nil, model.NewAppError("UpdateUserAuth", "app.user.update_auth_data.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	return userAuth, nil
}

func (a *SuiteService) sendUpdatedUserEvent(user model.User) {
	// exclude event creator user from admin, member user broadcast
	omitUsers := make(map[string]bool, 1)
	omitUsers[user.Id] = true

	// declare admin and unsanitized copy of user
	adminCopyOfUser := user.DeepCopy()
	unsanitizedCopyOfUser := user.DeepCopy()

	a.SanitizeProfile(adminCopyOfUser, true)
	adminMessage := model.NewWebSocketEvent(model.WebsocketEventUserUpdated, "", "", "", omitUsers, "")
	adminMessage.Add("user", adminCopyOfUser)
	adminMessage.GetBroadcast().ContainsSensitiveData = true
	a.platform.Publish(adminMessage)

	a.SanitizeProfile(&user, false)
	message := model.NewWebSocketEvent(model.WebsocketEventUserUpdated, "", "", "", omitUsers, "")
	message.Add("user", &user)
	message.GetBroadcast().ContainsSanitizedData = true
	a.platform.Publish(message)

	// send unsanitized user to event creator
	sourceUserMessage := model.NewWebSocketEvent(model.WebsocketEventUserUpdated, "", "", unsanitizedCopyOfUser.Id, nil, "")
	sourceUserMessage.Add("user", unsanitizedCopyOfUser)
	a.platform.Publish(sourceUserMessage)
}

func (a *SuiteService) isUniqueToGroupNames(val string) *model.AppError {
	if val == "" {
		return nil
	}
	var notFoundErr *store.ErrNotFound
	group, err := a.platform.Store.Group().GetByName(val, model.GroupSearchOpts{})
	if err != nil && !errors.As(err, &notFoundErr) {
		return model.NewAppError("isUniqueToGroupNames", model.NoTranslation, nil, "", http.StatusInternalServerError).Wrap(err)
	}
	if group != nil {
		return model.NewAppError("isUniqueToGroupNames", model.NoTranslation, nil, fmt.Sprintf("group name %s exists", val), http.StatusBadRequest)
	}
	return nil
}

func (a *SuiteService) UpdateUser(c request.CTX, user *model.User, sendNotifications bool) (*model.User, *model.AppError) {
	prev, err := a.getUser(user.Id)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("UpdateUser", MissingAccountError, nil, "", http.StatusNotFound).Wrap(err)
		default:
			return nil, model.NewAppError("UpdateUser", "app.user.get.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	if prev.CreateAt != user.CreateAt {
		user.CreateAt = prev.CreateAt
	}

	if user.Username != prev.Username {
		if err := a.isUniqueToGroupNames(user.Username); err != nil {
			err.Where = "UpdateUser"
			return nil, err
		}
	}

	var newEmail string
	if user.Email != prev.Email {
		if !users.CheckUserDomain(user, *a.platform.Config().TeamSettings.RestrictCreationToDomains) {
			if !prev.IsGuest() && !prev.IsLDAPUser() && !prev.IsSAMLUser() {
				return nil, model.NewAppError("UpdateUser", "api.user.update_user.accepted_domain.app_error", nil, "", http.StatusBadRequest)
			}
		}

		if !users.CheckUserDomain(user, *a.platform.Config().GuestAccountsSettings.RestrictCreationToDomains) {
			if prev.IsGuest() && !prev.IsLDAPUser() && !prev.IsSAMLUser() {
				return nil, model.NewAppError("UpdateUser", "api.user.update_user.accepted_guest_domain.app_error", nil, "", http.StatusBadRequest)
			}
		}

		if *a.platform.Config().EmailSettings.RequireEmailVerification {
			newEmail = user.Email
			// Don't set new eMail on user account if email verification is required, this will be done as a post-verification action
			// to avoid users being able to set non-controlled eMails as their account email
			if _, appErr := a.GetUserByEmail(newEmail); appErr == nil {
				return nil, model.NewAppError("UpdateUser", "app.user.save.email_exists.app_error", nil, "user_id="+user.Id, http.StatusBadRequest)
			}

			//  When a bot is created, prev.Email will be an autogenerated faked email,
			//  which will not match a CLI email input during bot to user conversions.
			//  To update a bot users email, do not set the email to the faked email
			//  stored in prev.Email.  Allow using the email defined in the CLI
			if !user.IsBot {
				user.Email = prev.Email
			}
		}
	}

	userUpdate, err := a.updateUser(user, false)
	if err != nil {
		var appErr *model.AppError
		var invErr *store.ErrInvalidInput
		var conErr *store.ErrConflict
		switch {
		case errors.As(err, &appErr):
			return nil, appErr
		case errors.As(err, &invErr):
			return nil, model.NewAppError("UpdateUser", "app.user.update.find.app_error", nil, "", http.StatusBadRequest).Wrap(err)
		case errors.As(err, &conErr):
			if conErr.Resource == "Username" {
				return nil, model.NewAppError("UpdateUser", "app.user.save.username_exists.app_error", nil, "", http.StatusBadRequest).Wrap(err)
			}
			return nil, model.NewAppError("UpdateUser", "app.user.save.email_exists.app_error", nil, "", http.StatusBadRequest).Wrap(err)
		default:
			return nil, model.NewAppError("UpdateUser", "app.user.update.finding.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	if sendNotifications {
		if userUpdate.New.Email != userUpdate.Old.Email || newEmail != "" {
			if *a.platform.Config().EmailSettings.RequireEmailVerification {
				a.platform.Go(func() {
					if err := a.SendEmailVerification(userUpdate.New, newEmail, ""); err != nil {
						c.Logger().Error("Failed to send email verification", mlog.Err(err))
					}
				})
			} else {
				a.platform.Go(func() {
					if err := a.email.SendEmailChangeEmail(userUpdate.Old.Email, userUpdate.New.Email, userUpdate.New.Locale, a.GetSiteURL()); err != nil {
						c.Logger().Error("Failed to send email change email", mlog.Err(err))
					}
				})
			}
		}

		if userUpdate.New.Username != userUpdate.Old.Username {
			a.platform.Go(func() {
				if err := a.email.SendChangeUsernameEmail(userUpdate.New.Username, userUpdate.New.Email, userUpdate.New.Locale, a.GetSiteURL()); err != nil {
					c.Logger().Error("Failed to send change username email", mlog.Err(err))
				}
			})
		}
		a.sendUpdatedUserEvent(*userUpdate.New)
	}

	a.InvalidateCacheForUser(user.Id)
	a.onUserProfileChange(user.Id)

	return userUpdate.New, nil
}

func (a *SuiteService) UpdateUserActive(c request.CTX, userID string, active bool) *model.AppError {
	user, err := a.GetUser(userID)

	if err != nil {
		return err
	}
	if _, err = a.UpdateActive(c, user, active); err != nil {
		return err
	}

	return nil
}

func (a *SuiteService) UpdateUserNotifyProps(userID string, props map[string]string) *model.AppError {
	err := a.updateUserNotifyProps(userID, props)
	if err != nil {
		var appErr *model.AppError
		switch {
		case errors.As(err, &appErr):
			return appErr
		default:
			return model.NewAppError("UpdateUser", "app.user.update.finding.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	a.InvalidateCacheForUser(userID)
	a.onUserProfileChange(userID)

	return nil
}

func (a *SuiteService) UpdateMfa(c request.CTX, activate bool, userID, token string) *model.AppError {
	if activate {
		if err := a.ActivateMfa(userID, token); err != nil {
			return err
		}
	} else {
		if err := a.DeactivateMfa(userID); err != nil {
			return err
		}
	}

	a.platform.Go(func() {
		user, err := a.GetUser(userID)
		if err != nil {
			c.Logger().Error("Failed to get user", mlog.Err(err))
			return
		}

		if err := a.email.SendMfaChangeEmail(user.Email, activate, user.Locale, a.GetSiteURL()); err != nil {
			c.Logger().Error("Failed to send mfa change email", mlog.Err(err))
		}
	})

	return nil
}

func (a *SuiteService) UpdatePasswordByUserIdSendEmail(c request.CTX, userID, newPassword, method string) *model.AppError {
	user, err := a.GetUser(userID)
	if err != nil {
		return err
	}

	return a.UpdatePasswordSendEmail(c, user, newPassword, method)
}

func (a *SuiteService) UpdatePassword(user *model.User, newPassword string) *model.AppError {
	if err := a.IsPasswordValid(newPassword); err != nil {
		return err
	}

	hashedPassword := model.HashPassword(newPassword)

	if err := a.platform.Store.User().UpdatePassword(user.Id, hashedPassword); err != nil {
		return model.NewAppError("UpdatePassword", "api.user.update_password.failed.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	a.InvalidateCacheForUser(user.Id)

	return nil
}

func (a *SuiteService) UpdatePasswordSendEmail(c request.CTX, user *model.User, newPassword, method string) *model.AppError {
	if err := a.UpdatePassword(user, newPassword); err != nil {
		return err
	}

	a.platform.Go(func() {
		if err := a.email.SendPasswordChangeEmail(user.Email, method, user.Locale, a.GetSiteURL()); err != nil {
			c.Logger().Error("Failed to send password change email", mlog.Err(err))
		}
	})

	return nil
}

func (a *SuiteService) UpdateHashedPasswordByUserId(userID, newHashedPassword string) *model.AppError {
	user, err := a.GetUser(userID)
	if err != nil {
		return err
	}

	return a.UpdateHashedPassword(user, newHashedPassword)
}

func (a *SuiteService) UpdateHashedPassword(user *model.User, newHashedPassword string) *model.AppError {
	if err := a.platform.Store.User().UpdatePassword(user.Id, newHashedPassword); err != nil {
		return model.NewAppError("UpdatePassword", "api.user.update_password.failed.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	a.InvalidateCacheForUser(user.Id)

	return nil
}

func (a *SuiteService) ResetPasswordFromToken(c request.CTX, userSuppliedTokenString, newPassword string) *model.AppError {
	return a.resetPasswordFromToken(c, userSuppliedTokenString, newPassword, model.GetMillis())
}

func (a *SuiteService) resetPasswordFromToken(c request.CTX, userSuppliedTokenString, newPassword string, nowMilli int64) *model.AppError {
	token, err := a.GetPasswordRecoveryToken(userSuppliedTokenString)
	if err != nil {
		return err
	}
	if nowMilli-token.CreateAt >= PasswordRecoverExpiryTime {
		return model.NewAppError("resetPassword", "api.user.reset_password.link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	tokenData := struct {
		UserId string
		Email  string
	}{}

	err2 := json.Unmarshal([]byte(token.Extra), &tokenData)
	if err2 != nil {
		return model.NewAppError("resetPassword", "api.user.reset_password.token_parse.error", nil, "", http.StatusInternalServerError)
	}

	user, err := a.GetUser(tokenData.UserId)
	if err != nil {
		return err
	}

	if user.Email != tokenData.Email {
		return model.NewAppError("resetPassword", "api.user.reset_password.link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	if user.IsSSOUser() {
		return model.NewAppError("ResetPasswordFromCode", "api.user.reset_password.sso.app_error", nil, "userId="+user.Id, http.StatusBadRequest)
	}

	T := i18n.GetUserTranslations(user.Locale)

	if err := a.UpdatePasswordSendEmail(c, user, newPassword, T("api.user.reset_password.method")); err != nil {
		return err
	}

	if err := a.DeleteToken(token); err != nil {
		c.Logger().Warn("Failed to delete token", mlog.Err(err))
	}

	return nil
}

func (a *SuiteService) SendPasswordReset(email string, siteURL string) (bool, *model.AppError) {
	user, err := a.GetUserByEmail(email)
	if err != nil {
		return false, nil
	}

	if user.AuthData != nil && *user.AuthData != "" {
		return false, model.NewAppError("SendPasswordReset", "api.user.send_password_reset.sso.app_error", nil, "userId="+user.Id, http.StatusBadRequest)
	}

	token, err := a.CreatePasswordRecoveryToken(user.Id, user.Email)
	if err != nil {
		return false, err
	}

	result, eErr := a.email.SendPasswordResetEmail(user.Email, token, user.Locale, siteURL)
	if eErr != nil {
		return result, model.NewAppError("SendPasswordReset", "api.user.send_password_reset.send.app_error", nil, "err="+eErr.Error(), http.StatusInternalServerError)
	}

	return result, nil
}

func (a *SuiteService) CreatePasswordRecoveryToken(userID, email string) (*model.Token, *model.AppError) {
	tokenExtra := struct {
		UserId string
		Email  string
	}{
		userID,
		email,
	}
	jsonData, err := json.Marshal(tokenExtra)
	if err != nil {
		return nil, model.NewAppError("CreatePasswordRecoveryToken", "api.user.create_password_token.error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	token := model.NewToken(TokenTypePasswordRecovery, string(jsonData))

	if err := a.platform.Store.Token().Save(token); err != nil {
		var appErr *model.AppError
		switch {
		case errors.As(err, &appErr):
			return nil, appErr
		default:
			return nil, model.NewAppError("CreatePasswordRecoveryToken", "app.recover.save.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	return token, nil
}

func (a *SuiteService) GetPasswordRecoveryToken(token string) (*model.Token, *model.AppError) {
	rtoken, err := a.platform.Store.Token().GetByToken(token)
	if err != nil {
		return nil, model.NewAppError("GetPasswordRecoveryToken", "api.user.reset_password.invalid_link.app_error", nil, "", http.StatusBadRequest).Wrap(err)
	}
	if rtoken.Type != TokenTypePasswordRecovery {
		return nil, model.NewAppError("GetPasswordRecoveryToken", "api.user.reset_password.broken_token.app_error", nil, "", http.StatusBadRequest)
	}
	return rtoken, nil
}

func (a *SuiteService) GetTokenById(token string) (*model.Token, *model.AppError) {
	rtoken, err := a.platform.Store.Token().GetByToken(token)

	if err != nil {
		var status int

		switch err.(type) {
		case *store.ErrNotFound:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}

		return nil, model.NewAppError("GetTokenById", "api.user.create_user.signup_link_invalid.app_error", nil, "", status).Wrap(err)
	}

	return rtoken, nil
}

func (a *SuiteService) DeleteToken(token *model.Token) *model.AppError {
	err := a.platform.Store.Token().Delete(token.Token)
	if err != nil {
		return model.NewAppError("DeleteToken", "app.recover.delete.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	return nil
}

func (a *SuiteService) UpdateUserRoles(c request.CTX, userID string, newRoles string, sendWebSocketEvent bool) (*model.User, *model.AppError) {
	user, err := a.GetUser(userID)
	if err != nil {
		err.StatusCode = http.StatusBadRequest
		return nil, err
	}

	return a.UpdateUserRolesWithUser(c, user, newRoles, sendWebSocketEvent)
}

func (a *SuiteService) UpdateUserRolesWithUser(c request.CTX, user *model.User, newRoles string, sendWebSocketEvent bool) (*model.User, *model.AppError) {

	if err := a.CheckRolesExist(strings.Fields(newRoles)); err != nil {
		return nil, err
	}

	user.Roles = newRoles
	uchan := make(chan store.StoreResult, 1)
	go func() {
		userUpdate, err := a.platform.Store.User().Update(user, true)
		uchan <- store.StoreResult{Data: userUpdate, NErr: err}
		close(uchan)
	}()

	schan := make(chan store.StoreResult, 1)
	go func() {
		id, err := a.platform.Store.Session().UpdateRoles(user.Id, newRoles)
		schan <- store.StoreResult{Data: id, NErr: err}
		close(schan)
	}()

	result := <-uchan
	if result.NErr != nil {
		var appErr *model.AppError
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(result.NErr, &appErr):
			return nil, appErr
		case errors.As(result.NErr, &invErr):
			return nil, model.NewAppError("UpdateUserRoles", "app.user.update.find.app_error", nil, "", http.StatusBadRequest).Wrap(result.NErr)
		default:
			return nil, model.NewAppError("UpdateUserRoles", "app.user.update.finding.app_error", nil, "", http.StatusInternalServerError).Wrap(result.NErr)
		}
	}
	ruser := result.Data.(*model.UserUpdate).New

	if result := <-schan; result.NErr != nil {
		// soft error since the user roles were still updated
		c.Logger().Warn("Failed during updating user roles", mlog.Err(result.NErr))
	}

	a.InvalidateCacheForUser(user.Id)
	a.ClearSessionCacheForUser(user.Id)

	if sendWebSocketEvent {
		message := model.NewWebSocketEvent(model.WebsocketEventUserRoleUpdated, "", "", user.Id, nil, "")
		message.Add("user_id", user.Id)
		message.Add("roles", newRoles)
		a.platform.Publish(message)
	}

	return ruser, nil
}

func (a *SuiteService) PermanentDeleteUser(c *request.Context, user *model.User) *model.AppError {
	c.Logger().Warn("Attempting to permanently delete account", mlog.String("user_id", user.Id), mlog.String("user_email", user.Email))
	if user.IsInRole(model.SystemAdminRoleId) {
		c.Logger().Warn("You are deleting a user that is a system administrator.  You may need to set another account as the system administrator using the command line tools.", mlog.String("user_email", user.Email))
	}

	if _, err := a.UpdateActive(c, user, false); err != nil {
		return err
	}

	if err := a.platform.Store.Session().PermanentDeleteSessionsByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.session.permanent_delete_sessions_by_user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.UserAccessToken().DeleteAllForUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.user_access_token.delete.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.OAuth().PermanentDeleteAuthDataByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.oauth.permanent_delete_auth_data_by_user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.Webhook().PermanentDeleteIncomingByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.webhooks.permanent_delete_incoming_by_user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.Webhook().PermanentDeleteOutgoingByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.webhooks.permanent_delete_outgoing_by_user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.Command().PermanentDeleteByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.user.permanentdeleteuser.internal_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.Preference().PermanentDeleteByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.preference.permanent_delete_by_user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.Channel().PermanentDeleteMembersByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.channel.permanent_delete_members_by_user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.Group().PermanentDeleteMembersByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.group.permanent_delete_members_by_user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.Post().PermanentDeleteByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.post.permanent_delete_by_user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.Bot().PermanentDelete(user.Id); err != nil {
		var invErr *store.ErrInvalidInput
		switch {
		case errors.As(err, &invErr):
			return model.NewAppError("PermanentDeleteUser", "app.bot.permenent_delete.bad_id", map[string]any{"user_id": invErr.Value}, "", http.StatusBadRequest).Wrap(err)
		default: // last fallback in case it doesn't map to an existing app error.
			return model.NewAppError("PermanentDeleteUser", "app.bot.permanent_delete.internal_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	infos, err := a.platform.Store.FileInfo().GetForUser(user.Id)
	if err != nil {
		c.Logger().Warn("Error getting file list for user from FileInfoStore", mlog.Err(err))
	}

	for _, info := range infos {
		res, err := a.platform.FileBackend().FileExists(info.Path)
		if err != nil {
			c.Logger().Warn(
				"Error checking existence of file",
				mlog.String("path", info.Path),
				mlog.Err(err),
			)
			continue
		}

		if !res {
			c.Logger().Warn("File not found", mlog.String("path", info.Path))
			continue
		}

		err = a.platform.FileBackend().RemoveFile(info.Path)

		if err != nil {
			c.Logger().Warn(
				"Unable to remove file",
				mlog.String("path", info.Path),
				mlog.Err(err),
			)
		}
	}

	// delete directory containing user's profile image
	profileImageDirectory := getProfileImageDirectory(user.Id)
	profileImagePath := getProfileImagePath(user.Id)
	resProfileImageExists, errProfileImageExists := a.platform.FileBackend().FileExists(profileImagePath)

	fileHandlingErrorsFound := false

	if errProfileImageExists != nil {
		fileHandlingErrorsFound = true
		mlog.Warn(
			"Error checking existence of profile image.",
			mlog.String("path", profileImagePath),
			mlog.Err(errProfileImageExists),
		)
	}

	if resProfileImageExists {
		errRemoveDirectory := a.platform.FileBackend().RemoveDirectory(profileImageDirectory)

		if errRemoveDirectory != nil {
			fileHandlingErrorsFound = true
			mlog.Warn(
				"Unable to remove profile image directory",
				mlog.String("path", profileImageDirectory),
				mlog.Err(errRemoveDirectory),
			)
		}
	}

	if _, err := a.platform.Store.FileInfo().PermanentDeleteByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.file_info.permanent_delete_by_user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.User().PermanentDelete(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.user.permanent_delete.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.Audit().PermanentDeleteByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.audit.permanent_delete_by_user.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	if err := a.platform.Store.Team().RemoveAllMembersByUser(user.Id); err != nil {
		return model.NewAppError("PermanentDeleteUser", "app.team.remove_member.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	a.InvalidateCacheForUser(user.Id)

	if fileHandlingErrorsFound {
		return model.NewAppError("PermanentDeleteUser", "app.file_info.permanent_delete_by_user.app_error", nil, "Couldn't delete profile image of the user.", http.StatusAccepted)
	}

	c.Logger().Warn("Permanently deleted account", mlog.String("user_email", user.Email), mlog.String("user_id", user.Id))

	return nil
}

func (a *SuiteService) PermanentDeleteAllUsers(c *request.Context) *model.AppError {
	users, err := a.platform.Store.User().GetAll()
	if err != nil {
		return model.NewAppError("PermanentDeleteAllUsers", "app.user.get.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	for _, user := range users {
		a.PermanentDeleteUser(c, user)
	}

	return nil
}

func (a *SuiteService) SendEmailVerification(user *model.User, newEmail, redirect string) *model.AppError {
	token, err := a.email.CreateVerifyEmailToken(user.Id, newEmail)
	if err != nil {
		switch {
		case errors.Is(err, email.CreateEmailTokenError):
			return model.NewAppError("CreateVerifyEmailToken", "api.user.create_email_token.error", nil, "", http.StatusInternalServerError)
		default:
			return model.NewAppError("CreateVerifyEmailToken", "app.recover.save.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
	}

	if _, err := a.GetStatus(user.Id); err != nil {
		if err.StatusCode != http.StatusNotFound {
			return err
		}
		eErr := a.email.SendVerifyEmail(newEmail, user.Locale, a.GetSiteURL(), token.Token, redirect)
		if eErr != nil {
			return model.NewAppError("SendVerifyEmail", "api.user.send_verify_email_and_forget.failed.error", nil, "", http.StatusInternalServerError).Wrap(eErr)
		}

		return nil
	}

	if err := a.email.SendEmailChangeVerifyEmail(newEmail, user.Locale, a.GetSiteURL(), token.Token); err != nil {
		return model.NewAppError("sendEmailChangeVerifyEmail", "api.user.send_email_change_verify_email_and_forget.error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return nil
}

func (a *SuiteService) VerifyEmailFromToken(c request.CTX, userSuppliedTokenString string) *model.AppError {
	token, err := a.GetVerifyEmailToken(userSuppliedTokenString)
	if err != nil {
		return err
	}
	if model.GetMillis()-token.CreateAt >= PasswordRecoverExpiryTime {
		return model.NewAppError("VerifyEmailFromToken", "api.user.verify_email.link_expired.app_error", nil, "", http.StatusBadRequest)
	}

	tokenData := struct {
		UserId string
		Email  string
	}{}

	err2 := json.Unmarshal([]byte(token.Extra), &tokenData)
	if err2 != nil {
		return model.NewAppError("VerifyEmailFromToken", "api.user.verify_email.token_parse.error", nil, "", http.StatusInternalServerError)
	}

	user, err := a.GetUser(tokenData.UserId)
	if err != nil {
		return err
	}

	tokenData.Email = strings.ToLower(tokenData.Email)
	if err := a.VerifyUserEmail(tokenData.UserId, tokenData.Email); err != nil {
		return err
	}

	if user.Email != tokenData.Email {
		a.platform.Go(func() {
			if err := a.email.SendEmailChangeEmail(user.Email, tokenData.Email, user.Locale, a.GetSiteURL()); err != nil {
				mlog.Error("Failed to send email change email", mlog.Err(err))
			}
		})
	}

	if err := a.DeleteToken(token); err != nil {
		c.Logger().Warn("Failed to delete token", mlog.Err(err))
	}

	return nil
}

func (a *SuiteService) GetVerifyEmailToken(token string) (*model.Token, *model.AppError) {
	rtoken, err := a.platform.Store.Token().GetByToken(token)
	if err != nil {
		return nil, model.NewAppError("GetVerifyEmailToken", "api.user.verify_email.bad_link.app_error", nil, "", http.StatusBadRequest).Wrap(err)
	}
	if rtoken.Type != TokenTypeVerifyEmail {
		return nil, model.NewAppError("GetVerifyEmailToken", "api.user.verify_email.broken_token.app_error", nil, "", http.StatusBadRequest)
	}
	return rtoken, nil
}

// GetTotalUsersStats is used for the DM list total
func (a *SuiteService) GetTotalUsersStats(viewRestrictions *model.ViewUsersRestrictions) (*model.UsersStats, *model.AppError) {
	count, err := a.platform.Store.User().Count(model.UserCountOptions{
		IncludeBotAccounts: true,
		ViewRestrictions:   viewRestrictions,
	})
	if err != nil {
		return nil, model.NewAppError("GetTotalUsersStats", "app.user.get_total_users_count.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	stats := &model.UsersStats{
		TotalUsersCount: count,
	}
	return stats, nil
}

// GetFilteredUsersStats is used to get a count of users based on the set of filters supported by UserCountOptions.
func (a *SuiteService) GetFilteredUsersStats(options *model.UserCountOptions) (*model.UsersStats, *model.AppError) {
	count, err := a.platform.Store.User().Count(*options)
	if err != nil {
		return nil, model.NewAppError("GetFilteredUsersStats", "app.user.get_total_users_count.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	stats := &model.UsersStats{
		TotalUsersCount: count,
	}
	return stats, nil
}

func (a *SuiteService) VerifyUserEmail(userID, email string) *model.AppError {
	if _, err := a.platform.Store.User().VerifyEmail(userID, email); err != nil {
		return model.NewAppError("VerifyUserEmail", "app.user.verify_email.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	a.InvalidateCacheForUser(userID)

	user, err := a.GetUser(userID)

	if err != nil {
		return err
	}

	a.sendUpdatedUserEvent(*user)

	return nil
}

func (a *SuiteService) SearchUsers(props *model.UserSearch, options *model.UserSearchOptions) ([]*model.User, *model.AppError) {
	if props.WithoutTeam {
		return a.SearchUsersWithoutTeam(props.Term, options)
	}
	if props.InChannelId != "" {
		return a.SearchUsersInChannel(props.InChannelId, props.Term, options)
	}
	if props.NotInChannelId != "" {
		return a.SearchUsersNotInChannel(props.TeamId, props.NotInChannelId, props.Term, options)
	}
	if props.NotInTeamId != "" {
		return a.SearchUsersNotInTeam(props.NotInTeamId, props.Term, options)
	}
	if props.InGroupId != "" {
		return a.SearchUsersInGroup(props.InGroupId, props.Term, options)
	}
	if props.NotInGroupId != "" {
		return a.SearchUsersNotInGroup(props.NotInGroupId, props.Term, options)
	}
	return a.SearchUsersInTeam(props.TeamId, props.Term, options)
}

func (a *SuiteService) SearchUsersInChannel(channelID string, term string, options *model.UserSearchOptions) ([]*model.User, *model.AppError) {
	term = strings.TrimSpace(term)
	users, err := a.platform.Store.User().SearchInChannel(channelID, term, options)
	if err != nil {
		return nil, model.NewAppError("SearchUsersInChannel", "app.user.search.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	for _, user := range users {
		a.SanitizeProfile(user, options.IsAdmin)
	}

	return users, nil
}

func (a *SuiteService) SearchUsersNotInChannel(teamID string, channelID string, term string, options *model.UserSearchOptions) ([]*model.User, *model.AppError) {
	term = strings.TrimSpace(term)
	users, err := a.platform.Store.User().SearchNotInChannel(teamID, channelID, term, options)
	if err != nil {
		return nil, model.NewAppError("SearchUsersNotInChannel", "app.user.search.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	for _, user := range users {
		a.SanitizeProfile(user, options.IsAdmin)
	}

	return users, nil
}

func (a *SuiteService) SearchUsersInTeam(teamID, term string, options *model.UserSearchOptions) ([]*model.User, *model.AppError) {
	term = strings.TrimSpace(term)

	users, err := a.platform.Store.User().Search(teamID, term, options)
	if err != nil {
		return nil, model.NewAppError("SearchUsersInTeam", "app.user.search.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	for _, user := range users {
		a.SanitizeProfile(user, options.IsAdmin)
	}

	return users, nil
}

func (a *SuiteService) SearchUsersNotInTeam(notInTeamId string, term string, options *model.UserSearchOptions) ([]*model.User, *model.AppError) {
	term = strings.TrimSpace(term)
	users, err := a.platform.Store.User().SearchNotInTeam(notInTeamId, term, options)
	if err != nil {
		return nil, model.NewAppError("SearchUsersNotInTeam", "app.user.search.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	for _, user := range users {
		a.SanitizeProfile(user, options.IsAdmin)
	}

	return users, nil
}

func (a *SuiteService) SearchUsersWithoutTeam(term string, options *model.UserSearchOptions) ([]*model.User, *model.AppError) {
	term = strings.TrimSpace(term)
	users, err := a.platform.Store.User().SearchWithoutTeam(term, options)
	if err != nil {
		return nil, model.NewAppError("SearchUsersWithoutTeam", "app.user.search.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	for _, user := range users {
		a.SanitizeProfile(user, options.IsAdmin)
	}

	return users, nil
}

func (a *SuiteService) SearchUsersInGroup(groupID string, term string, options *model.UserSearchOptions) ([]*model.User, *model.AppError) {
	term = strings.TrimSpace(term)
	users, err := a.platform.Store.User().SearchInGroup(groupID, term, options)
	if err != nil {
		return nil, model.NewAppError("SearchUsersInGroup", "app.user.search.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	for _, user := range users {
		a.SanitizeProfile(user, options.IsAdmin)
	}

	return users, nil
}

func (a *SuiteService) SearchUsersNotInGroup(groupID string, term string, options *model.UserSearchOptions) ([]*model.User, *model.AppError) {
	term = strings.TrimSpace(term)
	users, err := a.platform.Store.User().SearchNotInGroup(groupID, term, options)
	if err != nil {
		return nil, model.NewAppError("SearchUsersNotInGroup", "app.user.search.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	for _, user := range users {
		a.SanitizeProfile(user, options.IsAdmin)
	}

	return users, nil
}

func (a *SuiteService) AutocompleteUsersInChannel(teamID string, channelID string, term string, options *model.UserSearchOptions) (*model.UserAutocompleteInChannel, *model.AppError) {
	term = strings.TrimSpace(term)

	autocomplete, err := a.platform.Store.User().AutocompleteUsersInChannel(teamID, channelID, term, options)
	if err != nil {
		return nil, model.NewAppError("AutocompleteUsersInChannel", "app.user.search.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	for _, user := range autocomplete.InChannel {
		a.SanitizeProfile(user, options.IsAdmin)
	}

	for _, user := range autocomplete.OutOfChannel {
		a.SanitizeProfile(user, options.IsAdmin)
	}

	return autocomplete, nil
}

func (a *SuiteService) AutocompleteUsersInTeam(teamID string, term string, options *model.UserSearchOptions) (*model.UserAutocompleteInTeam, *model.AppError) {
	term = strings.TrimSpace(term)

	users, err := a.platform.Store.User().Search(teamID, term, options)
	if err != nil {
		return nil, model.NewAppError("AutocompleteUsersInTeam", "app.user.search.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	for _, user := range users {
		a.SanitizeProfile(user, options.IsAdmin)
	}

	autocomplete := &model.UserAutocompleteInTeam{}
	autocomplete.InTeam = users
	return autocomplete, nil
}

func (a *SuiteService) UpdateOAuthUserAttrs(userData io.Reader, user *model.User, provider einterfaces.OAuthProvider, service string, tokenUser *model.User) *model.AppError {
	oauthUser, err1 := provider.GetUserFromJSON(userData, tokenUser)
	if err1 != nil {
		return model.NewAppError("UpdateOAuthUserAttrs", "api.user.update_oauth_user_attrs.get_user.app_error", map[string]any{"Service": service}, "", http.StatusBadRequest).Wrap(err1)
	}

	userAttrsChanged := false

	if oauthUser.Username != user.Username {
		if existingUser, _ := a.GetUserByUsername(oauthUser.Username); existingUser == nil {
			user.Username = oauthUser.Username
			userAttrsChanged = true
		}
	}

	if oauthUser.GetFullName() != user.GetFullName() {
		user.FirstName = oauthUser.FirstName
		user.LastName = oauthUser.LastName
		userAttrsChanged = true
	}

	if oauthUser.Email != user.Email {
		if existingUser, _ := a.GetUserByEmail(oauthUser.Email); existingUser == nil {
			user.Email = oauthUser.Email
			userAttrsChanged = true
		}
	}

	if user.DeleteAt > 0 {
		// Make sure they are not disabled
		user.DeleteAt = 0
		userAttrsChanged = true
	}

	if userAttrsChanged {
		users, err := a.platform.Store.User().Update(user, true)
		if err != nil {
			var appErr *model.AppError
			var invErr *store.ErrInvalidInput
			switch {
			case errors.As(err, &appErr):
				return appErr
			case errors.As(err, &invErr):
				return model.NewAppError("UpdateOAuthUserAttrs", "app.user.update.find.app_error", nil, "", http.StatusBadRequest).Wrap(err)
			default:
				return model.NewAppError("UpdateOAuthUserAttrs", "app.user.update.finding.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
			}
		}

		user = users.New
		a.InvalidateCacheForUser(user.Id)
	}

	return nil
}

func (a *SuiteService) RestrictUsersGetByPermissions(userID string, options *model.UserGetOptions) (*model.UserGetOptions, *model.AppError) {
	restrictions, err := a.GetViewUsersRestrictions(userID)
	if err != nil {
		return nil, err
	}

	options.ViewRestrictions = restrictions
	return options, nil
}

// FilterNonGroupTeamMembers returns the subset of the given user IDs of the users who are not members of groups
// associated to the team excluding bots.
func (a *SuiteService) FilterNonGroupTeamMembers(userIDs []string, team *model.Team) ([]string, error) {
	teamGroupUsers, err := a.GetTeamGroupUsers(team.Id)
	if err != nil {
		return nil, err
	}
	return a.filterNonGroupUsers(userIDs, teamGroupUsers)
}

// FilterNonGroupChannelMembers returns the subset of the given user IDs of the users who are not members of groups
// associated to the channel excluding bots
func (a *SuiteService) FilterNonGroupChannelMembers(userIDs []string, channel *model.Channel) ([]string, error) {
	channelGroupUsers, err := a.GetChannelGroupUsers(channel.Id)
	if err != nil {
		return nil, err
	}
	return a.filterNonGroupUsers(userIDs, channelGroupUsers)
}

// filterNonGroupUsers is a helper function that takes a list of user ids and a list of users
// and returns the list of normal users present in userIDs but not in groupUsers.
func (a *SuiteService) filterNonGroupUsers(userIDs []string, groupUsers []*model.User) ([]string, error) {
	nonMemberIds := []string{}
	users, err := a.platform.Store.User().GetProfileByIds(context.Background(), userIDs, nil, false)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		userIsMember := user.IsBot

		for _, pu := range groupUsers {
			if pu.Id == user.Id {
				userIsMember = true
				break
			}
		}
		if !userIsMember {
			nonMemberIds = append(nonMemberIds, user.Id)
		}
	}

	return nonMemberIds, nil
}

func (a *SuiteService) RestrictUsersSearchByPermissions(userID string, options *model.UserSearchOptions) (*model.UserSearchOptions, *model.AppError) {
	restrictions, err := a.GetViewUsersRestrictions(userID)
	if err != nil {
		return nil, err
	}

	options.ViewRestrictions = restrictions
	return options, nil
}

func (a *SuiteService) UserCanSeeOtherUser(userID string, otherUserId string) (bool, *model.AppError) {
	if userID == otherUserId {
		return true, nil
	}

	restrictions, err := a.GetViewUsersRestrictions(userID)
	if err != nil {
		return false, err
	}

	if restrictions == nil {
		return true, nil
	}

	if len(restrictions.Teams) > 0 {
		result, err := a.platform.Store.Team().UserBelongsToTeams(otherUserId, restrictions.Teams)
		if err != nil {
			return false, model.NewAppError("UserCanSeeOtherUser", "app.team.user_belongs_to_teams.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
		}
		if result {
			return true, nil
		}
	}

	if len(restrictions.Channels) > 0 {
		result, err := a.userBelongsToChannels(otherUserId, restrictions.Channels)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}

	return false, nil
}

func (a *SuiteService) userBelongsToChannels(userID string, channelIDs []string) (bool, *model.AppError) {
	belongs, err := a.platform.Store.Channel().UserBelongsToChannels(userID, channelIDs)
	if err != nil {
		return false, model.NewAppError("userBelongsToChannels", "app.channel.user_belongs_to_channels.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return belongs, nil
}

func (a *SuiteService) GetViewUsersRestrictions(userID string) (*model.ViewUsersRestrictions, *model.AppError) {
	if a.HasPermissionTo(userID, model.PermissionViewMembers) {
		return nil, nil
	}

	teamIDs, nErr := a.platform.Store.Team().GetUserTeamIds(userID, true)
	if nErr != nil {
		return nil, model.NewAppError("GetViewUsersRestrictions", "app.team.get_user_team_ids.app_error", nil, "", http.StatusInternalServerError).Wrap(nErr)
	}

	teamIDsWithPermission := []string{}
	for _, teamID := range teamIDs {
		if a.HasPermissionToTeam(userID, teamID, model.PermissionViewMembers) {
			teamIDsWithPermission = append(teamIDsWithPermission, teamID)
		}
	}

	userChannelMembers, err := a.platform.Store.Channel().GetAllChannelMembersForUser(userID, true, true)
	if err != nil {
		return nil, model.NewAppError("GetViewUsersRestrictions", "app.channel.get_channels.get.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	channelIDs := []string{}
	for channelID := range userChannelMembers {
		channelIDs = append(channelIDs, channelID)
	}

	return &model.ViewUsersRestrictions{Teams: teamIDsWithPermission, Channels: channelIDs}, nil
}

// PromoteGuestToUser Convert user's roles and all his membership's roles from
// guest roles to regular user roles.
func (a *SuiteService) PromoteGuestToUser(c *request.Context, user *model.User, requestorId string) *model.AppError {
	nErr := a.promoteGuestToUser(user)
	a.InvalidateCacheForUser(user.Id)
	if nErr != nil {
		return model.NewAppError("PromoteGuestToUser", "app.user.promote_guest.user_update.app_error", nil, "", http.StatusInternalServerError).Wrap(nErr)
	}
	userTeams, nErr := a.platform.Store.Team().GetTeamsByUserId(user.Id)
	if nErr != nil {
		return model.NewAppError("PromoteGuestToUser", "app.team.get_all.app_error", nil, "", http.StatusInternalServerError).Wrap(nErr)
	}

	for _, team := range userTeams {
		// Soft error if there is an issue joining the default channels
		if err := a.channels.JoinDefaultChannels(c, team.Id, user, false, requestorId); err != nil {
			c.Logger().Warn("Failed to join default channels", mlog.String("user_id", user.Id), mlog.String("team_id", team.Id), mlog.String("requestor_id", requestorId), mlog.Err(err))
		}
	}

	promotedUser, err := a.GetUser(user.Id)
	if err != nil {
		c.Logger().Warn("Failed to get user on promote guest to user", mlog.Err(err))
	} else {
		a.sendUpdatedUserEvent(*promotedUser)
		if uErr := a.platform.UpdateSessionsIsGuest(promotedUser.Id, promotedUser.IsGuest()); uErr != nil {
			c.Logger().Warn("Unable to update user sessions", mlog.String("user_id", promotedUser.Id), mlog.Err(uErr))
		}
	}

	teamMembers, err := a.GetTeamMembersForUser(user.Id, "", true)
	if err != nil {
		c.Logger().Warn("Failed to get team members for user on promote guest to user", mlog.Err(err))
	}

	for _, member := range teamMembers {
		a.sendUpdatedMemberRoleEvent(user.Id, member)

		channelMembers, appErr := a.channels.GetChannelMembersForUser(c, member.TeamId, user.Id)
		if appErr != nil {
			c.Logger().Warn("Failed to get channel members for user on promote guest to user", mlog.Err(appErr))
		}

		for _, member := range channelMembers {
			a.platform.InvalidateCacheForChannelMembers(member.ChannelId)

			evt := model.NewWebSocketEvent(model.WebsocketEventChannelMemberUpdated, "", "", user.Id, nil, "")
			memberJSON, jsonErr := json.Marshal(member)
			if jsonErr != nil {
				return model.NewAppError("PromoteGuestToUser", "api.marshal_error", nil, "", http.StatusInternalServerError).Wrap(jsonErr)
			}
			evt.Add("channelMember", string(memberJSON))
			a.platform.Publish(evt)
		}
	}

	a.ClearSessionCacheForUser(user.Id)
	return nil
}

// DemoteUserToGuest Convert user's roles and all his membership's roles from
// regular user roles to guest roles.
func (a *SuiteService) DemoteUserToGuest(c request.CTX, user *model.User) *model.AppError {
	demotedUser, nErr := a.demoteUserToGuest(user)
	a.InvalidateCacheForUser(user.Id)
	if nErr != nil {
		return model.NewAppError("DemoteUserToGuest", "app.user.demote_user_to_guest.user_update.app_error", nil, "", http.StatusInternalServerError).Wrap(nErr)
	}

	a.sendUpdatedUserEvent(*demotedUser)
	if uErr := a.platform.UpdateSessionsIsGuest(demotedUser.Id, demotedUser.IsGuest()); uErr != nil {
		c.Logger().Warn("Unable to update user sessions", mlog.String("user_id", demotedUser.Id), mlog.Err(uErr))
	}

	teamMembers, err := a.GetTeamMembersForUser(user.Id, "", true)
	if err != nil {
		c.Logger().Warn("Failed to get team members for users on demote user to guest", mlog.Err(err))
	}

	for _, member := range teamMembers {
		a.sendUpdatedMemberRoleEvent(user.Id, member)

		channelMembers, appErr := a.channels.GetChannelMembersForUser(c, member.TeamId, user.Id)
		if appErr != nil {
			c.Logger().Warn("Failed to get channel members for users on demote user to guest", mlog.Err(appErr))
			continue
		}

		for _, member := range channelMembers {
			a.platform.InvalidateCacheForChannelMembers(member.ChannelId)

			evt := model.NewWebSocketEvent(model.WebsocketEventChannelMemberUpdated, "", "", user.Id, nil, "")
			memberJSON, jsonErr := json.Marshal(member)
			if jsonErr != nil {
				return model.NewAppError("DemoteUserToGuest", "api.marshal_error", nil, "", http.StatusInternalServerError).Wrap(jsonErr)
			}
			evt.Add("channelMember", string(memberJSON))
			a.platform.Publish(evt)
		}
	}

	a.ClearSessionCacheForUser(user.Id)
	return nil
}

func (a *SuiteService) PublishUserTyping(userID, channelID, parentId string) *model.AppError {
	omitUsers := make(map[string]bool, 1)
	omitUsers[userID] = true

	event := model.NewWebSocketEvent(model.WebsocketEventTyping, "", channelID, "", omitUsers, "")
	event.Add("parent_id", parentId)
	event.Add("user_id", userID)
	a.platform.Publish(event)

	return nil
}

// invalidateUserCacheAndPublish Invalidates cache for a user and publishes user updated event
func (a *SuiteService) invalidateUserCacheAndPublish(userID string) {
	a.InvalidateCacheForUser(userID)

	user, userErr := a.GetUser(userID)
	if userErr != nil {
		mlog.Error("Error in getting users profile", mlog.String("user_id", userID), mlog.Err(userErr))
		return
	}

	options := a.platform.Config().GetSanitizeOptions()
	user.SanitizeProfile(options)

	message := model.NewWebSocketEvent(model.WebsocketEventUserUpdated, "", "", "", nil, "")
	message.Add("user", user)
	a.platform.Publish(message)
}

// GetKnownUsers returns the list of user ids of users with any direct
// relationship with a user. That means any user sharing any channel, including
// direct and group channels.
func (a *SuiteService) GetKnownUsers(userID string) ([]string, *model.AppError) {
	users, err := a.platform.Store.User().GetKnownUsers(userID)
	if err != nil {
		return nil, model.NewAppError("GetKnownUsers", "app.user.get_known_users.get_users.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

// ConvertBotToUser converts a bot to user.
func (a *SuiteService) ConvertBotToUser(c request.CTX, bot *model.Bot, userPatch *model.UserPatch, sysadmin bool) (*model.User, *model.AppError) {
	user, nErr := a.platform.Store.User().Get(c.Context(), bot.UserId)
	if nErr != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(nErr, &nfErr):
			return nil, model.NewAppError("ConvertBotToUser", MissingAccountError, nil, "", http.StatusNotFound).Wrap(nErr)
		default:
			return nil, model.NewAppError("ConvertBotToUser", "app.user.get.app_error", nil, "", http.StatusInternalServerError).Wrap(nErr)
		}
	}

	if sysadmin && !user.IsInRole(model.SystemAdminRoleId) {
		_, appErr := a.UpdateUserRoles(c,
			user.Id,
			fmt.Sprintf("%s %s", user.Roles, model.SystemAdminRoleId),
			false)
		if appErr != nil {
			return nil, appErr
		}
	}

	user.Patch(userPatch)

	user, err := a.UpdateUser(c, user, false)
	if err != nil {
		return nil, err
	}

	err = a.UpdatePassword(user, *userPatch.Password)
	if err != nil {
		return nil, err
	}

	appErr := a.platform.Store.Bot().PermanentDelete(bot.UserId)
	if appErr != nil {
		return nil, model.NewAppError("ConvertBotToUser", "app.user.convert_bot_to_user.app_error", nil, "", http.StatusInternalServerError).Wrap(appErr)
	}

	return user, nil
}

func (a *SuiteService) GetUsersWithInvalidEmails(page int, perPage int) ([]*model.User, *model.AppError) {
	users, err := a.platform.Store.User().GetUsersWithInvalidEmails(page, perPage, *a.platform.Config().TeamSettings.RestrictCreationToDomains)
	if err != nil {
		return nil, model.NewAppError("GetUsersPage", "app.user.get_profiles.app_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return users, nil
}

func getProfileImagePath(userID string) string {
	return filepath.Join("users", userID, "profile.png")
}

func getProfileImageDirectory(userID string) string {
	return filepath.Join("users", userID)
}

// onUserProfileChange is called when a user's profile has changed
// (username, email, profile image, ...)
func (a *SuiteService) onUserProfileChange(userID string) {
	syncService := a.sharedChannelSyncService
	if syncService == nil || !syncService.Active() {
		return
	}
	syncService.NotifyUserProfileChanged(userID)
}
