// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package bleveengine

import (
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

func (b *BleveEngine) IndexPost(post *model.Post, teamId string) *model.AppError {
	b.Mutex.RLock()
	defer b.Mutex.RUnlock()

	blvPost := BLVPostFromPost(post, teamId)
	if err := b.PostIndex.Index(blvPost.Id, blvPost); err != nil {
		return model.NewAppError("Bleveengine.IndexPost", "bleveengine.index_post.error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (b *BleveEngine) SearchPosts(channels *model.ChannelList, searchParams []*model.SearchParams, page, perPage int) ([]string, model.PostSearchMatches, *model.AppError) {
	channelQueries := []query.Query{}
	for _, channel := range *channels {
		channelIdQ := bleve.NewTermQuery(channel.Id)
		channelIdQ.SetField("ChannelId")
		channelQueries = append(channelQueries, channelIdQ)
	}
	channelDisjunctionQ := bleve.NewDisjunctionQuery(channelQueries...)

	var termQueries []query.Query
	var notTermQueries []query.Query
	var filters []query.Query
	var notFilters []query.Query

	typeQ := bleve.NewTermQuery("")
	typeQ.SetField("Type")
	filters = append(filters, typeQ)

	for i, params := range searchParams {
		// Date, channels and FromUsers filters come in all
		// searchParams iteration, and as they are global to the
		// query, we only need to process them once
		if i == 0 {
			if len(params.InChannels) > 0 {
				inChannels := []query.Query{}
				for _, channelId := range params.InChannels {
					channelQ := bleve.NewTermQuery(channelId)
					channelQ.SetField("ChannelId")
					inChannels = append(inChannels, channelQ)
				}
				filters = append(filters, bleve.NewDisjunctionQuery(inChannels...))
			}

			if len(params.ExcludedChannels) > 0 {
				excludedChannels := []query.Query{}
				for _, channelId := range params.ExcludedChannels {
					channelQ := bleve.NewTermQuery(channelId)
					channelQ.SetField("ChannelId")
					excludedChannels = append(excludedChannels, channelQ)
				}
				notFilters = append(notFilters, bleve.NewDisjunctionQuery(excludedChannels...))
			}

			if len(params.FromUsers) > 0 {
				fromUsers := []query.Query{}
				for _, userId := range params.FromUsers {
					userQ := bleve.NewTermQuery(userId)
					userQ.SetField("UserId")
					fromUsers = append(fromUsers, userQ)
				}
				filters = append(filters, bleve.NewDisjunctionQuery(fromUsers...))
			}

			if len(params.ExcludedUsers) > 0 {
				excludedUsers := []query.Query{}
				for _, userId := range params.ExcludedUsers {
					userQ := bleve.NewTermQuery(userId)
					userQ.SetField("UserId")
					excludedUsers = append(excludedUsers, userQ)
				}
				notFilters = append(notFilters, bleve.NewDisjunctionQuery(excludedUsers...))
			}

			if params.OnDate != "" {
				before, after := params.GetOnDateMillis()
				beforeFloat64 := float64(before)
				afterFloat64 := float64(after)
				onDateQ := bleve.NewNumericRangeQuery(&beforeFloat64, &afterFloat64)
				onDateQ.SetField("CreateAt")
				filters = append(filters, onDateQ)
			} else {
				if params.AfterDate != "" || params.BeforeDate != "" {
					var min, max *float64
					if params.AfterDate != "" {
						minf := float64(params.GetAfterDateMillis())
						min = &minf
					}

					if params.BeforeDate != "" {
						maxf := float64(params.GetBeforeDateMillis())
						max = &maxf
					}

					dateQ := bleve.NewNumericRangeQuery(min, max)
					dateQ.SetField("CreateAt")
					filters = append(filters, dateQ)
				}

				if params.ExcludedAfterDate != "" {
					minf := float64(params.GetExcludedAfterDateMillis())
					dateQ := bleve.NewNumericRangeQuery(&minf, nil)
					dateQ.SetField("CreateAt")
					notFilters = append(notFilters, dateQ)
				}

				if params.ExcludedBeforeDate != "" {
					maxf := float64(params.GetExcludedBeforeDateMillis())
					dateQ := bleve.NewNumericRangeQuery(nil, &maxf)
					dateQ.SetField("CreateAt")
					notFilters = append(notFilters, dateQ)
				}

				if params.ExcludedDate != "" {
					before, after := params.GetExcludedDateMillis()
					beforef := float64(before)
					afterf := float64(after)
					onDateQ := bleve.NewNumericRangeQuery(&beforef, &afterf)
					onDateQ.SetField("CreateAt")
					notFilters = append(notFilters, onDateQ)
				}
			}
		}

		if params.IsHashtag {
			if params.Terms != "" {
				hashtagQ := bleve.NewMatchQuery(params.Terms)
				hashtagQ.SetField("Hashtags")
				termQueries = append(termQueries, hashtagQ)
			} else if params.ExcludedTerms != "" {
				hashtagQ := bleve.NewMatchQuery(params.ExcludedTerms)
				hashtagQ.SetField("Hashtags")
				notTermQueries = append(notTermQueries, hashtagQ)
			}
		} else {
			if len(params.Terms) > 0 {
				query := bleve.NewBooleanQuery()
				messageQ := bleve.NewMatchQuery(params.Terms)
				messageQ.SetField("Message")

				if searchParams[0].OrTerms {
					query.AddShould(messageQ)
				} else {
					query.AddMust(messageQ)
				}
				termQueries = append(termQueries, messageQ)
			}

			if len(params.ExcludedTerms) > 0 {
				messageQ := bleve.NewMatchQuery(params.ExcludedTerms)
				messageQ.SetField("Message")
				notTermQueries = append(notTermQueries, messageQ)
			}
		}
	}

	allTermsQ := bleve.NewBooleanQuery()
	allTermsQ.AddMustNot(notTermQueries...)
	if searchParams[0].OrTerms {
		allTermsQ.AddShould(termQueries...)
	} else {
		allTermsQ.AddMust(termQueries...)
	}

	query := bleve.NewBooleanQuery()
	query.AddMust(channelDisjunctionQ)

	if len(termQueries) > 0 || len(notTermQueries) > 0 {
		query.AddMust(allTermsQ)
	}

	if len(filters) > 0 {
		query.AddMust(bleve.NewConjunctionQuery(filters...))
	}
	if len(notFilters) > 0 {
		query.AddMustNot(notFilters...)
	}

	search := bleve.NewSearchRequest(query)
	results, err := b.PostIndex.Search(search)
	if err != nil {
		return nil, nil, model.NewAppError("Bleveengine.SearchPosts", "bleveengine.search_posts.error", nil, err.Error(), http.StatusInternalServerError)
	}

	postIds := []string{}
	matches := model.PostSearchMatches{}

	for _, r := range results.Hits {
		postIds = append(postIds, r.ID)
	}

	return postIds, matches, nil
}

func (b *BleveEngine) DeletePost(post *model.Post) *model.AppError {
	b.Mutex.RLock()
	defer b.Mutex.RUnlock()

	if err := b.PostIndex.Delete(post.Id); err != nil {
		return model.NewAppError("Bleveengine.DeletePost", "bleveengine.delete_post.error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (b *BleveEngine) IndexChannel(channel *model.Channel) *model.AppError {
	b.Mutex.RLock()
	defer b.Mutex.RUnlock()

	blvChannel := BLVChannelFromChannel(channel)
	if err := b.ChannelIndex.Index(blvChannel.Id, blvChannel); err != nil {
		return model.NewAppError("Bleveengine.IndexChannel", "bleveengine.index_channel.error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (b *BleveEngine) SearchChannels(teamId, term string) ([]string, *model.AppError) {
	teamIdQ := bleve.NewTermQuery(teamId)
	teamIdQ.SetField("TeamId")
	queries := []query.Query{teamIdQ}

	if term != "" {
		nameSuggestQ := bleve.NewPrefixQuery(strings.ToLower(term))
		nameSuggestQ.SetField("NameSuggest")
		queries = append(queries, nameSuggestQ)
	}

	query := bleve.NewSearchRequest(bleve.NewConjunctionQuery(queries...))
	results, err := b.ChannelIndex.Search(query)
	if err != nil {
		return nil, model.NewAppError("Bleveengine.SearchChannels", "bleveengine.search_channels.error", nil, err.Error(), http.StatusInternalServerError)
	}

	channelIds := []string{}
	for _, result := range results.Hits {
		channelIds = append(channelIds, result.ID)
	}

	return channelIds, nil
}

func (b *BleveEngine) DeleteChannel(channel *model.Channel) *model.AppError {
	b.Mutex.RLock()
	defer b.Mutex.RUnlock()

	if err := b.ChannelIndex.Delete(channel.Id); err != nil {
		return model.NewAppError("Bleveengine.DeleteChannel", "bleveengine.delete_channel.error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (b *BleveEngine) IndexUser(user *model.User, teamsIds, channelsIds []string) *model.AppError {
	b.Mutex.RLock()
	defer b.Mutex.RUnlock()

	blvUser := BLVUserFromUserAndTeams(user, teamsIds, channelsIds)
	if err := b.UserIndex.Index(blvUser.Id, blvUser); err != nil {
		return model.NewAppError("Bleveengine.IndexUser", "bleveengine.index_user.error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (b *BleveEngine) SearchUsersInChannel(teamId, channelId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) ([]string, []string, *model.AppError) {
	if restrictedToChannels != nil && len(restrictedToChannels) == 0 {
		return []string{}, []string{}, nil
	}

	// users in channel
	var queries []query.Query
	if term != "" {
		termQ := bleve.NewPrefixQuery(strings.ToLower(term))
		if options.AllowFullNames {
			termQ.SetField("SuggestionsWithFullname")
		} else {
			termQ.SetField("SuggestionsWithoutFullname")
		}
		queries = append(queries, termQ)
	}

	channelIdQ := bleve.NewTermQuery(channelId)
	channelIdQ.SetField("ChannelsIds")
	queries = append(queries, channelIdQ)

	query := bleve.NewConjunctionQuery(queries...)

	uchan, err := b.UserIndex.Search(bleve.NewSearchRequest(query))
	if err != nil {
		return nil, nil, model.NewAppError("Bleveengine.SearchUsersInChannel", "bleveengine.search_users_in_channel.uchan.error", nil, err.Error(), http.StatusInternalServerError)
	}

	// users not in channel
	boolQ := bleve.NewBooleanQuery()

	if term != "" {
		termQ := bleve.NewPrefixQuery(strings.ToLower(term))
		if options.AllowFullNames {
			termQ.SetField("SuggestionsWithFullname")
		} else {
			termQ.SetField("SuggestionsWithoutFullname")
		}
		boolQ.AddMust(termQ)
	}

	teamIdQ := bleve.NewTermQuery(teamId)
	teamIdQ.SetField("TeamsIds")
	boolQ.AddMust(teamIdQ)

	outsideChannelIdQ := bleve.NewTermQuery(channelId)
	outsideChannelIdQ.SetField("ChannelsIds")
	boolQ.AddMustNot(outsideChannelIdQ)

	if len(restrictedToChannels) > 0 {
		restrictedChannelsQ := bleve.NewDisjunctionQuery()
		for _, channelId := range restrictedToChannels {
			restrictedChannelQ := bleve.NewTermQuery(channelId)
			restrictedChannelsQ.AddQuery(restrictedChannelQ)
		}
		boolQ.AddMust(restrictedChannelsQ)
	}

	nuchan, err := b.UserIndex.Search(bleve.NewSearchRequest(boolQ))
	if err != nil {
		return nil, nil, model.NewAppError("Bleveengine.SearchUsersInChannel", "bleveengine.search_users_in_channel.nuchan.error", nil, err.Error(), http.StatusInternalServerError)
	}

	uchanIds := []string{}
	for _, result := range uchan.Hits {
		uchanIds = append(uchanIds, result.ID)
	}

	nuchanIds := []string{}
	for _, result := range nuchan.Hits {
		nuchanIds = append(nuchanIds, result.ID)
	}

	return uchanIds, nuchanIds, nil
}

func (b *BleveEngine) SearchUsersInTeam(teamId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) ([]string, *model.AppError) {
	if restrictedToChannels != nil && len(restrictedToChannels) == 0 {
		return []string{}, nil
	}

	var rootQ query.Query
	if term == "" && teamId == "" && restrictedToChannels == nil {
		rootQ = bleve.NewMatchAllQuery()
	} else {
		boolQ := bleve.NewBooleanQuery()

		if term != "" {
			termQ := bleve.NewPrefixQuery(strings.ToLower(term))
			if options.AllowFullNames {
				termQ.SetField("SuggestionsWithFullname")
			} else {
				termQ.SetField("SuggestionsWithoutFullname")
			}
			boolQ.AddMust(termQ)
		}

		if len(restrictedToChannels) > 0 {
			// restricted channels are already filtered by team, so we
			// can search only those matches
			restrictedChannelsQ := []query.Query{}
			for _, channelId := range restrictedToChannels {
				channelIdQ := bleve.NewTermQuery(channelId)
				channelIdQ.SetField("ChannelsIds")
				restrictedChannelsQ = append(restrictedChannelsQ, channelIdQ)
			}
			boolQ.AddMust(bleve.NewDisjunctionQuery(restrictedChannelsQ...))
		} else {
			// this means that we only need to restrict by team
			if teamId != "" {
				teamIdQ := bleve.NewTermQuery(teamId)
				teamIdQ.SetField("TeamsIds")
				boolQ.AddMust(teamIdQ)
			}
		}

		rootQ = boolQ
	}

	search := bleve.NewSearchRequest(rootQ)

	results, err := b.UserIndex.Search(search)
	if err != nil {
		return nil, model.NewAppError("Bleveengine.SearchUsersInTeam", "bleveengine.search_users_in_team.error", nil, err.Error(), http.StatusInternalServerError)
	}

	usersIds := []string{}
	for _, r := range results.Hits {
		usersIds = append(usersIds, r.ID)
	}

	return usersIds, nil
}

func (b *BleveEngine) DeleteUser(user *model.User) *model.AppError {
	b.Mutex.RLock()
	defer b.Mutex.RUnlock()

	if err := b.UserIndex.Delete(user.Id); err != nil {
		return model.NewAppError("Bleveengine.DeleteUser", "bleveengine.delete_user.error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}
