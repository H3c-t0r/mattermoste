// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	l4g "github.com/alecthomas/log4go"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
)

const (
	TRIGGERWORDS_FULL       = 0
	TRIGGERWORDS_STARTSWITH = 1
)

func handleWebhookEvents(post *model.Post, team *model.Team, channel *model.Channel, user *model.User) {
	if !utils.Cfg.ServiceSettings.EnableOutgoingWebhooks {
		return
	}

	if channel.Type != model.CHANNEL_OPEN {
		return
	}

	hchan := Srv.Store.Webhook().GetOutgoingByTeam(team.Id)
	result := <-hchan
	if result.Err != nil {
		l4g.Error(utils.T("api.post.handle_webhook_events_and_forget.getting.error"), result.Err)
		return
	}

	hooks := result.Data.([]*model.OutgoingWebhook)
	if len(hooks) == 0 {
		return
	}

	splitWords := strings.Fields(post.Message)
	if len(splitWords) == 0 {
		return
	}
	firstWord := splitWords[0]

	relevantHooks := []*model.OutgoingWebhook{}
	for _, hook := range hooks {
		if hook.ChannelId == post.ChannelId || len(hook.ChannelId) == 0 {
			if hook.ChannelId == post.ChannelId && len(hook.TriggerWords) == 0 {
				relevantHooks = append(relevantHooks, hook)
			} else if hook.TriggerWhen == TRIGGERWORDS_FULL && hook.HasTriggerWord(firstWord) {
				relevantHooks = append(relevantHooks, hook)
			} else if hook.TriggerWhen == TRIGGERWORDS_STARTSWITH && hook.TriggerWordStartsWith(firstWord) {
				relevantHooks = append(relevantHooks, hook)
			}
		}
	}

	for _, hook := range relevantHooks {
		go func(hook *model.OutgoingWebhook) {
			payload := &model.OutgoingWebhookPayload{
				Token:       hook.Token,
				TeamId:      hook.TeamId,
				TeamDomain:  team.Name,
				ChannelId:   post.ChannelId,
				ChannelName: channel.Name,
				Timestamp:   post.CreateAt,
				UserId:      post.UserId,
				UserName:    user.Username,
				PostId:      post.Id,
				Text:        post.Message,
				TriggerWord: firstWord,
			}
			var body io.Reader
			var contentType string
			if hook.ContentType == "application/json" {
				body = strings.NewReader(payload.ToJSON())
				contentType = "application/json"
			} else {
				body = strings.NewReader(payload.ToFormValues())
				contentType = "application/x-www-form-urlencoded"
			}
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: *utils.Cfg.ServiceSettings.EnableInsecureOutgoingConnections},
			}
			client := &http.Client{Transport: tr}

			for _, url := range hook.CallbackURLs {
				go func(url string) {
					req, _ := http.NewRequest("POST", url, body)
					req.Header.Set("Content-Type", contentType)
					req.Header.Set("Accept", "application/json")
					if resp, err := client.Do(req); err != nil {
						l4g.Error(utils.T("api.post.handle_webhook_events_and_forget.event_post.error"), err.Error())
					} else {
						defer func() {
							ioutil.ReadAll(resp.Body)
							resp.Body.Close()
						}()
						respProps := model.MapFromJson(resp.Body)

						if text, ok := respProps["text"]; ok {
							if _, err := CreateWebhookPost(hook.CreatorId, hook.TeamId, post.ChannelId, text, respProps["username"], respProps["icon_url"], post.Props, post.Type); err != nil {
								l4g.Error(utils.T("api.post.handle_webhook_events_and_forget.create_post.error"), err)
							}
						}
					}
				}(url)
			}

		}(hook)
	}
}

func CreateWebhookPost(userId, teamId, channelId, text, overrideUsername, overrideIconUrl string, props model.StringInterface, postType string) (*model.Post, *model.AppError) {
	post := &model.Post{UserId: userId, ChannelId: channelId, Message: text, Type: postType}
	post.AddProp("from_webhook", "true")

	if utils.Cfg.ServiceSettings.EnablePostUsernameOverride {
		if len(overrideUsername) != 0 {
			post.AddProp("override_username", overrideUsername)
		} else {
			post.AddProp("override_username", model.DEFAULT_WEBHOOK_USERNAME)
		}
	}

	if utils.Cfg.ServiceSettings.EnablePostIconOverride {
		if len(overrideIconUrl) != 0 {
			post.AddProp("override_icon_url", overrideIconUrl)
		}
	}

	post.Message = parseSlackLinksToMarkdown(post.Message)

	if len(props) > 0 {
		for key, val := range props {
			if key == "attachments" {
				parseSlackAttachment(post, val)
			} else if key != "override_icon_url" && key != "override_username" && key != "from_webhook" {
				post.AddProp(key, val)
			}
		}
	}

	if _, err := CreatePost(post, teamId, false); err != nil {
		return nil, model.NewLocAppError("CreateWebhookPost", "api.post.create_webhook_post.creating.app_error", nil, "err="+err.Message)
	}

	return post, nil
}
