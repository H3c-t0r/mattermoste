// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"net/http"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"
)

const (
	DAY_MILLISECONDS   = 24 * 60 * 60 * 1000
	MONTH_MILLISECONDS = 31 * DAY_MILLISECONDS
)

func (a *App) GetAnalytics(name string, teamId string) (model.AnalyticsRows, *model.AppError) {
	skipIntensiveQueries := false
	var systemUserCount int64
	systemUserCount, err := a.Srv().Store.User().Count(model.UserCountOptions{})
	if err != nil {
		return nil, err
	}

	if systemUserCount > int64(*a.Config().AnalyticsSettings.MaxUsersForStatistics) {
		mlog.Debug("More than limit users are on the system, intensive queries skipped", mlog.Int("limit", *a.Config().AnalyticsSettings.MaxUsersForStatistics))
		skipIntensiveQueries = true
	}

	if name == "standard" {
		var rows model.AnalyticsRows = make([]*model.AnalyticsRow, 11)
		rows[0] = &model.AnalyticsRow{Name: "channel_open_count", Value: 0}
		rows[1] = &model.AnalyticsRow{Name: "channel_private_count", Value: 0}
		rows[2] = &model.AnalyticsRow{Name: "post_count", Value: 0}
		rows[3] = &model.AnalyticsRow{Name: "unique_user_count", Value: 0}
		rows[4] = &model.AnalyticsRow{Name: "team_count", Value: 0}
		rows[5] = &model.AnalyticsRow{Name: "total_websocket_connections", Value: 0}
		rows[6] = &model.AnalyticsRow{Name: "total_master_db_connections", Value: 0}
		rows[7] = &model.AnalyticsRow{Name: "total_read_db_connections", Value: 0}
		rows[8] = &model.AnalyticsRow{Name: "daily_active_users", Value: 0}
		rows[9] = &model.AnalyticsRow{Name: "monthly_active_users", Value: 0}
		rows[10] = &model.AnalyticsRow{Name: "inactive_user_count", Value: 0}

		openChan := make(chan store.StoreResult, 1)
		privateChan := make(chan store.StoreResult, 1)
		go func() {
			count, err2 := a.Srv().Store.Channel().AnalyticsTypeCount(teamId, model.CHANNEL_OPEN)
			openChan <- store.StoreResult{Data: count, NErr: err2}
			close(openChan)
		}()
		go func() {
			count, err2 := a.Srv().Store.Channel().AnalyticsTypeCount(teamId, model.CHANNEL_PRIVATE)
			privateChan <- store.StoreResult{Data: count, NErr: err2}
			close(privateChan)
		}()

		var userChan chan store.StoreResult
		var userInactiveChan chan store.StoreResult
		if teamId == "" {
			userInactiveChan = make(chan store.StoreResult, 1)
			go func() {
				count, err2 := a.Srv().Store.User().AnalyticsGetInactiveUsersCount()
				userInactiveChan <- store.StoreResult{Data: count, Err: err2}
				close(userInactiveChan)
			}()
		} else {
			userChan = make(chan store.StoreResult, 1)
			go func() {
				count, err2 := a.Srv().Store.User().Count(model.UserCountOptions{TeamId: teamId})
				userChan <- store.StoreResult{Data: count, Err: err2}
				close(userChan)
			}()
		}

		var postChan chan store.StoreResult
		if !skipIntensiveQueries {
			postChan = make(chan store.StoreResult, 1)
			go func() {
				count, err2 := a.Srv().Store.Post().AnalyticsPostCount(teamId, false, false)
				postChan <- store.StoreResult{Data: count, Err: err2}
				close(postChan)
			}()
		}

		teamCountChan := make(chan store.StoreResult, 1)
		go func() {
			teamCount, err2 := a.Srv().Store.Team().AnalyticsTeamCount(false)
			teamCountChan <- store.StoreResult{Data: teamCount, NErr: err2}
			close(teamCountChan)
		}()

		dailyActiveChan := make(chan store.StoreResult, 1)
		go func() {
			dailyActive, err2 := a.Srv().Store.User().AnalyticsActiveCount(DAY_MILLISECONDS, model.UserCountOptions{IncludeBotAccounts: false, IncludeDeleted: false})
			dailyActiveChan <- store.StoreResult{Data: dailyActive, Err: err2}
			close(dailyActiveChan)
		}()

		monthlyActiveChan := make(chan store.StoreResult, 1)
		go func() {
			monthlyActive, err2 := a.Srv().Store.User().AnalyticsActiveCount(MONTH_MILLISECONDS, model.UserCountOptions{IncludeBotAccounts: false, IncludeDeleted: false})
			monthlyActiveChan <- store.StoreResult{Data: monthlyActive, Err: err2}
			close(monthlyActiveChan)
		}()

		r := <-openChan
		if r.NErr != nil {
			return nil, model.NewAppError("GetAnalytics", "app.channel.analytics_type_count.app_error", nil, r.NErr.Error(), http.StatusInternalServerError)
		}
		rows[0].Value = float64(r.Data.(int64))

		r = <-privateChan
		if r.NErr != nil {
			return nil, model.NewAppError("GetAnalytics", "app.channel.analytics_type_count.app_error", nil, r.NErr.Error(), http.StatusInternalServerError)
		}
		rows[1].Value = float64(r.Data.(int64))

		if postChan == nil {
			rows[2].Value = -1
		} else {
			r = <-postChan
			if r.Err != nil {
				return nil, r.Err
			}
			rows[2].Value = float64(r.Data.(int64))
		}

		if userChan == nil {
			rows[3].Value = float64(systemUserCount)
		} else {
			r = <-userChan
			if r.Err != nil {
				return nil, r.Err
			}
			rows[3].Value = float64(r.Data.(int64))
		}

		if userInactiveChan == nil {
			rows[10].Value = -1
		} else {
			r = <-userInactiveChan
			if r.Err != nil {
				return nil, r.Err
			}
			rows[10].Value = float64(r.Data.(int64))
		}

		r = <-teamCountChan
		if r.NErr != nil {
			return nil, model.NewAppError("GetAnalytics", "app.team.analytics_team_count.app_error", nil, r.NErr.Error(), http.StatusInternalServerError)
		}
		rows[4].Value = float64(r.Data.(int64))

		// If in HA mode then aggregate all the stats
		if a.Cluster() != nil && *a.Config().ClusterSettings.Enable {
			stats, err2 := a.Cluster().GetClusterStats()
			if err2 != nil {
				return nil, err2
			}

			totalSockets := a.TotalWebsocketConnections()
			totalMasterDb := a.Srv().Store.TotalMasterDbConnections()
			totalReadDb := a.Srv().Store.TotalReadDbConnections()

			for _, stat := range stats {
				totalSockets = totalSockets + stat.TotalWebsocketConnections
				totalMasterDb = totalMasterDb + stat.TotalMasterDbConnections
				totalReadDb = totalReadDb + stat.TotalReadDbConnections
			}

			rows[5].Value = float64(totalSockets)
			rows[6].Value = float64(totalMasterDb)
			rows[7].Value = float64(totalReadDb)

		} else {
			rows[5].Value = float64(a.TotalWebsocketConnections())
			rows[6].Value = float64(a.Srv().Store.TotalMasterDbConnections())
			rows[7].Value = float64(a.Srv().Store.TotalReadDbConnections())
		}

		r = <-dailyActiveChan
		if r.Err != nil {
			return nil, r.Err
		}
		rows[8].Value = float64(r.Data.(int64))

		r = <-monthlyActiveChan
		if r.Err != nil {
			return nil, r.Err
		}
		rows[9].Value = float64(r.Data.(int64))

		return rows, nil
	} else if name == "bot_post_counts_day" {
		if skipIntensiveQueries {
			rows := model.AnalyticsRows{&model.AnalyticsRow{Name: "", Value: -1}}
			return rows, nil
		}
		return a.Srv().Store.Post().AnalyticsPostCountsByDay(&model.AnalyticsPostCountsOptions{
			TeamId:        teamId,
			BotsOnly:      true,
			YesterdayOnly: false,
		})
	} else if name == "post_counts_day" {
		if skipIntensiveQueries {
			rows := model.AnalyticsRows{&model.AnalyticsRow{Name: "", Value: -1}}
			return rows, nil
		}
		return a.Srv().Store.Post().AnalyticsPostCountsByDay(&model.AnalyticsPostCountsOptions{
			TeamId:        teamId,
			BotsOnly:      false,
			YesterdayOnly: false,
		})
	} else if name == "user_counts_with_posts_day" {
		if skipIntensiveQueries {
			rows := model.AnalyticsRows{&model.AnalyticsRow{Name: "", Value: -1}}
			return rows, nil
		}

		return a.Srv().Store.Post().AnalyticsUserCountsWithPostsByDay(teamId)
	} else if name == "extra_counts" {
		var rows model.AnalyticsRows = make([]*model.AnalyticsRow, 6)
		rows[0] = &model.AnalyticsRow{Name: "file_post_count", Value: 0}
		rows[1] = &model.AnalyticsRow{Name: "hashtag_post_count", Value: 0}
		rows[2] = &model.AnalyticsRow{Name: "incoming_webhook_count", Value: 0}
		rows[3] = &model.AnalyticsRow{Name: "outgoing_webhook_count", Value: 0}
		rows[4] = &model.AnalyticsRow{Name: "command_count", Value: 0}
		rows[5] = &model.AnalyticsRow{Name: "session_count", Value: 0}

		iHookChan := make(chan store.StoreResult, 1)
		go func() {
			c, err2 := a.Srv().Store.Webhook().AnalyticsIncomingCount(teamId)
			iHookChan <- store.StoreResult{Data: c, NErr: err2}
			close(iHookChan)
		}()

		oHookChan := make(chan store.StoreResult, 1)
		go func() {
			c, err2 := a.Srv().Store.Webhook().AnalyticsOutgoingCount(teamId)
			oHookChan <- store.StoreResult{Data: c, NErr: err2}
			close(oHookChan)
		}()

		commandChan := make(chan store.StoreResult, 1)
		go func() {
			c, nErr := a.Srv().Store.Command().AnalyticsCommandCount(teamId)
			commandChan <- store.StoreResult{Data: c, NErr: nErr}
			close(commandChan)
		}()

		sessionChan := make(chan store.StoreResult, 1)
		go func() {
			count, err2 := a.Srv().Store.Session().AnalyticsSessionCount()
			sessionChan <- store.StoreResult{Data: count, NErr: err2}
			close(sessionChan)
		}()

		var fileChan chan store.StoreResult
		var hashtagChan chan store.StoreResult

		if !skipIntensiveQueries {
			fileChan = make(chan store.StoreResult, 1)
			go func() {
				count, err2 := a.Srv().Store.Post().AnalyticsPostCount(teamId, true, false)
				fileChan <- store.StoreResult{Data: count, Err: err2}
				close(fileChan)
			}()

			hashtagChan = make(chan store.StoreResult, 1)
			go func() {
				count, err2 := a.Srv().Store.Post().AnalyticsPostCount(teamId, false, true)
				hashtagChan <- store.StoreResult{Data: count, Err: err2}
				close(hashtagChan)
			}()
		}

		if fileChan == nil {
			rows[0].Value = -1
		} else {
			r := <-fileChan
			if r.Err != nil {
				return nil, r.Err
			}
			rows[0].Value = float64(r.Data.(int64))
		}

		if hashtagChan == nil {
			rows[1].Value = -1
		} else {
			r := <-hashtagChan
			if r.Err != nil {
				return nil, r.Err
			}
			rows[1].Value = float64(r.Data.(int64))
		}

		r := <-iHookChan
		if r.NErr != nil {
			return nil, model.NewAppError("GetAnalytics", "app.webhooks.analytics_incoming_count.app_error", nil, r.NErr.Error(), http.StatusInternalServerError)
		}
		rows[2].Value = float64(r.Data.(int64))

		r = <-oHookChan
		if r.NErr != nil {
			return nil, model.NewAppError("GetAnalytics", "app.webhooks.analytics_outgoing_count.app_error", nil, r.NErr.Error(), http.StatusInternalServerError)
		}
		rows[3].Value = float64(r.Data.(int64))

		r = <-commandChan
		if r.NErr != nil {
			return nil, model.NewAppError("GetAnalytics", "app.analytics.getanalytics.internal_error", nil, err.Error(), http.StatusInternalServerError)
		}
		rows[4].Value = float64(r.Data.(int64))

		r = <-sessionChan
		if r.NErr != nil {
			return nil, model.NewAppError("GetAnalytics", "app.session.analytics_session_count.app_error", nil, r.NErr.Error(), http.StatusInternalServerError)
		}
		rows[5].Value = float64(r.Data.(int64))

		return rows, nil
	}

	return nil, nil
}

func (a *App) GetRecentlyActiveUsersForTeam(teamId string) (map[string]*model.User, *model.AppError) {
	users, err := a.Srv().Store.User().GetRecentlyActiveUsersForTeam(teamId, 0, 100, nil)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*model.User)

	for _, user := range users {
		userMap[user.Id] = user
	}

	return userMap, nil
}

func (a *App) GetRecentlyActiveUsersForTeamPage(teamId string, page, perPage int, asAdmin bool, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError) {
	users, err := a.Srv().Store.User().GetRecentlyActiveUsersForTeam(teamId, page*perPage, perPage, viewRestrictions)
	if err != nil {
		return nil, err
	}

	return a.sanitizeProfiles(users, asAdmin), nil
}

func (a *App) GetNewUsersForTeamPage(teamId string, page, perPage int, asAdmin bool, viewRestrictions *model.ViewUsersRestrictions) ([]*model.User, *model.AppError) {
	users, err := a.Srv().Store.User().GetNewUsersForTeam(teamId, page*perPage, perPage, viewRestrictions)
	if err != nil {
		return nil, err
	}

	return a.sanitizeProfiles(users, asAdmin), nil
}
