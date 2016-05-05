// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"encoding/json"
	"io"
)

type InitialLoad struct {
	User              *User                  `json:"user"`
	TeamMembers       []*TeamMember          `json:"team_members"`
	Teams             []*Team                `json:"teams"`
	DirectProfiles    map[string]*User       `json:"direct_profiles"`
	Preferences       Preferences            `json:"preferences"`
	DirectTeamMembers map[string]*TeamMember `json:"direct_team_members"`
	ClientCfg         map[string]string      `json:"client_cfg"`
	LicenseCfg        map[string]string      `json:"license_cfg"`
	NoAccounts        bool                   `json:"no_accounts"`
}

func (me *InitialLoad) ToJson() string {
	b, err := json.Marshal(me)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func InitialLoadFromJson(data io.Reader) *InitialLoad {
	decoder := json.NewDecoder(data)
	var o InitialLoad
	err := decoder.Decode(&o)
	if err == nil {
		return &o
	} else {
		return nil
	}
}
