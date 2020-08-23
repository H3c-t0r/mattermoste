// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"fmt"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store/storetest/mocks"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNoticeValidation(t *testing.T) {
	th := SetupWithStoreMock(t)
	mockStore := th.App.Srv().Store.(*mocks.Store)
	mockRoleStore := mocks.RoleStore{}
	mockSystemStore := mocks.SystemStore{}
	mockUserStore := mocks.UserStore{}
	mockPostStore := mocks.PostStore{}
	mockStore.On("Role").Return(&mockRoleStore)
	mockStore.On("System").Return(&mockSystemStore)
	mockStore.On("User").Return(&mockUserStore)
	mockStore.On("Post").Return(&mockPostStore)
	mockSystemStore.On("SaveOrUpdate", &model.System{Name: "ActiveLicenseId", Value: ""}).Return(nil)
	mockUserStore.On("Count", model.UserCountOptions{IncludeBotAccounts: false, IncludeDeleted: true, ExcludeRegularUsers: false, TeamId: "", ChannelId: "", ViewRestrictions: (*model.ViewUsersRestrictions)(nil), Roles: []string(nil), ChannelRoles: []string(nil), TeamRoles: []string(nil)}).Return(int64(1), nil)
	defer th.TearDown()

	type args struct {
		client               model.NoticeClientType
		clientVersion        string
		locale               string
		sku                  string
		postCount, userCount int64
		cloud                bool
		teamAdmin            bool
		systemAdmin          bool
		notice               *model.ProductNotice
	}
	messages := map[string]model.NoticeMessageInternal{
		"enUS": {
			Description: "descr",
			Title:       "title",
		},
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantOk  bool
	}{
		{
			name: "general notice",
			args: args{
				client:        "mobile",
				clientVersion: "1.2.3",

				notice: &model.ProductNotice{
					Conditions:        model.Conditions{},
					ID:                "123",
					LocalizedMessages: messages,
				},
			},
			wantErr: false,
			wantOk:  true,
		},
		{
			name: "mobile notice",
			args: args{
				client:        "desktop",
				clientVersion: "1.2.3",

				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						ClientType: model.NewNoticeClientType(model.NoticeClientType_Mobile),
					},
				},
			},
			wantErr: false,
			wantOk:  false,
		},
		{
			name: "notice with config check",
			args: args{
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						ServerConfig: map[string]interface{}{"ServiceSettings.LetsEncryptCertificateCacheFile": "./config/letsencrypt.cache"},
					},
				},
			},
			wantErr: false,
			wantOk:  true,
		},
		{
			name: "notice with failing config check",
			args: args{
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						ServerConfig: map[string]interface{}{"ServiceSettings.ZZ": "test"},
					},
				},
			},
			wantErr: false,
			wantOk:  false,
		},
		{
			name: "notice with server version check",
			args: args{
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						ServerVersion: []string{"> 4.0.0 < 99.0.0"},
					},
				},
			},
			wantErr: false,
			wantOk:  true,
		},
		{
			name: "notice with server version check that doesn't match",
			args: args{
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						ServerVersion: []string{"> 99.0.0"},
					},
				},
			},
			wantErr: false,
			wantOk:  false,
		},
		{
			name: "notice with server version check that is invalid",
			args: args{
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						ServerVersion: []string{"99.0.0 + 1.0.0"},
					},
				},
			},
			wantErr: true,
			wantOk:  false,
		},
		{
			name: "notice with user count",
			args: args{
				userCount: 300,
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						NumberOfUsers: model.NewInt64(400),
					},
				},
			},
			wantErr: false,
			wantOk:  false,
		},
		{
			name: "notice with good user count and bad post count",
			args: args{
				userCount: 500,
				postCount: 2000,
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						NumberOfUsers: model.NewInt64(400),
						NumberOfPosts: model.NewInt64(3000),
					},
				},
			},
			wantErr: false,
			wantOk:  false,
		},
		{
			name: "notice with date check",
			args: args{
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						DisplayDate: model.NewString("> 2000-03-01T00:00:00Z <= 2999-04-01T00:00:00Z"),
					},
				},
			},
			wantErr: false,
			wantOk:  true,
		},

		{
			name: "notice with date check that doesn't match",
			args: args{
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						DisplayDate: model.NewString("> 2999-03-01T00:00:00Z <= 3000-04-01T00:00:00Z"),
					},
				},
			},
			wantErr: false,
			wantOk:  false,
		},
		{
			name: "notice with bad date check",
			args: args{
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						DisplayDate: model.NewString("> 2000 -03-01T00:00:00Z <= 2999-04-01T00:00:00Z"),
					},
				},
			},
			wantErr: true,
			wantOk:  false,
		},
		{
			name: "notice with audience check",
			args: args{
				systemAdmin: true,
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						Audience: model.NewNoticeAudience(model.NoticeAudience_Sysadmin),
					},
				},
			},
			wantErr: false,
			wantOk:  true,
		},
		{
			name: "notice with correct sku",
			args: args{
				sku: "e20",
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						Sku: model.NewNoticeSKU(model.NoticeSKU_E20),
					},
				},
			},
			wantErr: false,
			wantOk:  true,
		},
		{
			name: "notice with incorrect sku",
			args: args{
				sku: "e20",
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						Sku: model.NewNoticeSKU(model.NoticeSKU_E10),
					},
				},
			},
			wantErr: false,
			wantOk:  false,
		},
		{
			name: "notice with sku check for all",
			args: args{
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						Sku: model.NewNoticeSKU(model.NoticeSKU_All),
					},
				},
			},
			wantErr: false,
			wantOk:  true,
		},
		{
			name: "notice with instance check cloud",
			args: args{
				cloud: true,
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						InstanceType: model.NewNoticeInstanceType(model.NoticeInstanceType_Cloud),
					},
				},
			},
			wantErr: false,
			wantOk:  true,
		},
		{
			name: "notice with instance check both",
			args: args{
				notice: &model.ProductNotice{
					Conditions: model.Conditions{
						InstanceType: model.NewNoticeInstanceType(model.NoticeInstanceType_Both),
					},
				},
			},
			wantErr: false,
			wantOk:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientVersion := tt.args.clientVersion
			if clientVersion == "" {
				clientVersion = "1.2.3"
			}
			if ok, err := noticeMatchesConditions(th.App, tt.args.client, clientVersion, tt.args.locale, tt.args.postCount, tt.args.userCount, tt.args.systemAdmin, tt.args.teamAdmin, tt.args.cloud, tt.args.sku, tt.args.notice); (err != nil) != tt.wantErr {
				t.Errorf("noticeMatchesConditions() error = %v, wantErr %v", err, tt.wantErr)
			} else if ok != tt.wantOk {
				t.Errorf("noticeMatchesConditions() result = %v, wantOk %v", ok, tt.wantOk)
			}
		})
	}
}

func TestNoticeFetch(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	notices := model.ProductNotices{model.ProductNotice{
		Conditions: model.Conditions{},
		ID:         "123",
		LocalizedMessages: map[string]model.NoticeMessageInternal{
			"enUS": {
				Description: "description",
				Title:       "title",
			},
		},
		Repeatable: nil,
	}}
	noticesBytes, appErr := notices.Marshal()
	require.NoError(t, appErr)

	notices2 := model.ProductNotices{model.ProductNotice{
		Conditions: model.Conditions{
			NumberOfPosts: model.NewInt64(99999),
		},
		ID: "333",
		LocalizedMessages: map[string]model.NoticeMessageInternal{
			"enUS": {
				Description: "description",
				Title:       "title",
			},
		},
		Repeatable: nil,
	}}
	noticesBytes2, appErr := notices2.Marshal()
	require.NoError(t, appErr)
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "notices.json") {
			w.Write(noticesBytes)
		} else {
			w.Write(noticesBytes2)
		}
	}))
	defer server1.Close()

	NOTICES_JSON_URL = fmt.Sprintf("http://%s/notices.json", server1.Listener.Addr().String())

	// fetch fake notices
	appErr = th.App.UpdateProductNotices()
	require.Nil(t, appErr)

	now := time.Now().UTC().Unix()
	// get them for specified user
	messages, appErr := th.App.GetProductNotices(0, th.BasicUser.Id, th.BasicTeam.Id, model.NoticeClientType_All, "1.2.3", "enUS")
	require.Nil(t, appErr)
	require.Len(t, messages, 1)

	// mark notices as viewed
	appErr = th.App.UpdateViewedProductNotices(th.BasicUser.Id, []string{messages[0].ID})
	require.Nil(t, appErr)

	// get them again, see that none are returned
	messages, appErr = th.App.GetProductNotices(now, th.BasicUser.Id, th.BasicTeam.Id, model.NoticeClientType_All, "1.2.3", "enUS")
	require.Nil(t, appErr)
	require.Len(t, messages, 0)

	// validate views table
	views, err := th.App.Srv().Store.ProductNotices().GetViews(th.BasicUser.Id)
	require.Nil(t, err)
	require.Len(t, views, 1)

	// fetch another set
	NOTICES_JSON_URL = fmt.Sprintf("http://%s/notices2.json", server1.Listener.Addr().String())

	// fetch fake notices
	appErr = th.App.UpdateProductNotices()
	require.Nil(t, appErr)

	// get them again, since conditions don't match we should be zero
	messages, appErr = th.App.GetProductNotices(now, th.BasicUser.Id, th.BasicTeam.Id, model.NoticeClientType_All, "1.2.3", "enUS")
	require.Nil(t, appErr)
	require.Len(t, messages, 0)

	// even though UpdateViewedProductNotices was called previously, the table should be empty, since there's cleanup done during UpdateProductNotices
	views, err = th.App.Srv().Store.ProductNotices().GetViews(th.BasicUser.Id)
	require.Nil(t, err)
	require.Len(t, views, 0)
}
