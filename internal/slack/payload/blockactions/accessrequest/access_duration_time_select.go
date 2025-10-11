/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package accessrequest

import (
	"encoding/json"
	"fmt"
	"time"
)

type AccessDurationTimeSelectPayload struct {
	Type      string `json:"type"`
	TriggerID string `json:"trigger_id"`

	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`

	View struct {
		ID              string `json:"id"`
		PrivateMetadata string `json:"private_metadata"`
		CallbackID      string `json:"callback_id"`

		State struct {
			Values struct {
				AccessDurationDateTimeBlock struct {
					AccessRequestAccessDurationTimeSelect struct {
						Type         string `json:"type"`
						SelectedTime string `json:"selected_time"`
					} `json:"access_request_access_duration_time_select"`
				} `json:"access_duration_date_time_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type AccessDurationTimeSelectPrivateMetadataPayload struct {
	ChannelID                        string    `json:"channel_id"`
	ChannelName                      string    `json:"channel_name"`
	RealName                         string    `json:"real_name"`
	RequireReason                    bool      `json:"require_reason"`
	SelectedRole                     string    `json:"selected_role"`
	SelectedChannelID                string    `json:"selected_channel_id"`
	SelectedChannelName              string    `json:"selected_channel_name"`
	SelectedStartDateOptionID        string    `json:"selected_date_option_id"`
	SelectedStartDateOptionName      string    `json:"selected_date_option_name"`
	TTL                              time.Time `json:"ttl"`
	SelectedStartDate                string    `json:"selected_start_date"`
	SelectedStartTime                string    `json:"selected_start_time"`
	SelectedAccessDurationOptionID   string    `json:"selected_access_duration_option_id"`
	SelectedAccessDurationOptionName string    `json:"selected_access_duration_option_name"`
	SelectedAccessDurationDate       string    `json:"selected_access_duration_date"`
}

type AccessDurationTimeSelect struct {
	RequesterChannelID               string
	RequesterChannelName             string
	RequesterRealName                string
	RequireReason                    bool
	SelectedRole                     string
	SelectedChannelID                string
	SelectedChannelName              string
	SelectedStartDateOptionID        string
	SelectedStartDateOptionName      string
	TTL                              time.Time
	SelectedStartDate                string
	SelectedStartTime                string
	SelectedAccessDurationOptionID   string
	SelectedAccessDurationOptionName string
	SelectedAccessDurationDate       string
	RequesterID                      string
	RequesterName                    string
	TriggerID                        string
	ViewHash                         string
	ViewID                           string
	// new fields
	AccessDurationTime string
}

func NewAccessDurationTimeSelectWithFirstOpt(payload *AccessDurationOptionSelect) *AccessDurationTimeSelect {
	return &AccessDurationTimeSelect{
		RequesterChannelID:               payload.RequesterChannelID,
		RequesterChannelName:             payload.RequesterChannelName,
		RequesterRealName:                payload.RequesterRealName,
		RequireReason:                    payload.RequireReason,
		SelectedRole:                     payload.SelectedRole,
		SelectedChannelID:                payload.SelectedChannelID,
		SelectedChannelName:              payload.SelectedChannelName,
		SelectedStartDateOptionID:        payload.SelectedStartDateOptionID,
		SelectedStartDateOptionName:      payload.SelectedStartDateOptionName,
		TTL:                              payload.TTL,
		SelectedStartDate:                payload.SelectedStartDate,
		SelectedStartTime:                payload.SelectedStartTime,
		SelectedAccessDurationOptionID:   payload.AccessDurationOptionID,
		SelectedAccessDurationOptionName: payload.AccessDurationOptionName,
		RequesterID:                      payload.RequesterID,
		RequesterName:                    payload.RequesterName,
		TriggerID:                        payload.TriggerID,
		ViewHash:                         payload.ViewHash,
		ViewID:                           payload.ViewID,
	}
}

func ParseAccessDurationTimeSelect(payloadStr string) (*AccessDurationTimeSelect, error) {
	var payload AccessDurationTimeSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", payloadStr)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata AccessDurationTimeSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", strPrivateMetadata)
	}

	return &AccessDurationTimeSelect{
		RequesterChannelID:               privateMetadata.ChannelID,
		RequesterChannelName:             privateMetadata.ChannelName,
		RequesterRealName:                privateMetadata.RealName,
		RequireReason:                    privateMetadata.RequireReason,
		SelectedRole:                     privateMetadata.SelectedRole,
		SelectedChannelID:                privateMetadata.SelectedChannelID,
		SelectedChannelName:              privateMetadata.SelectedChannelName,
		SelectedStartDateOptionID:        privateMetadata.SelectedStartDateOptionID,
		SelectedStartDateOptionName:      privateMetadata.SelectedStartDateOptionName,
		TTL:                              privateMetadata.TTL,
		SelectedStartDate:                privateMetadata.SelectedStartDate,
		SelectedStartTime:                privateMetadata.SelectedStartTime,
		SelectedAccessDurationOptionID:   privateMetadata.SelectedAccessDurationOptionID,
		SelectedAccessDurationOptionName: privateMetadata.SelectedAccessDurationOptionName,
		SelectedAccessDurationDate:       privateMetadata.SelectedAccessDurationDate,
		RequesterID:                      payload.User.ID,
		RequesterName:                    payload.User.Name,
		TriggerID:                        payload.TriggerID,
		ViewHash:                         payload.View.Hash,
		ViewID:                           payload.View.ID,
		AccessDurationTime:               payload.View.State.Values.AccessDurationDateTimeBlock.AccessRequestAccessDurationTimeSelect.SelectedTime,
	}, nil
}
