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

type StartDateSelectPayload struct {
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
				StartDateTimeBlock struct {
					AccessRequestStartDateSelect struct {
						Type         string `json:"type"`
						SelectedDate string `json:"selected_date"`
					} `json:"access_request_start_date_select"`
				} `json:"start_date_time_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type StartDateSelectPrivateMetadataPayload struct {
	ChannelID                   string    `json:"channel_id"`
	ChannelName                 string    `json:"channel_name"`
	RealName                    string    `json:"real_name"`
	RequireReason               bool      `json:"require_reason"`
	SelectedRole                string    `json:"selected_role"`
	SelectedChannelID           string    `json:"selected_channel_id"`
	SelectedChannelName         string    `json:"selected_channel_name"`
	SelectedStartDateOptionID   string    `json:"selected_date_option_id"`
	SelectedStartDateOptionName string    `json:"selected_date_option_name"`
	TTL                         time.Time `json:"ttl"`
}

type StartDateSelect struct {
	RequesterChannelID          string
	RequesterChannelName        string
	RequesterRealName           string
	RequireReason               bool
	SelectedRole                string
	SelectedChannelID           string
	SelectedChannelName         string
	SelectedStartDateOptionID   string
	SelectedStartDateOptionName string
	RequesterID                 string
	RequesterName               string
	TriggerID                   string
	ViewHash                    string
	ViewID                      string
	// new fields
	TTL       time.Time
	StartDate string
}

func ParseStartDateSelect(payloadStr string) (*StartDateSelect, error) {
	var payload StartDateSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", payloadStr)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata StartDateSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", strPrivateMetadata)
	}
	return &StartDateSelect{
		RequesterChannelID:          privateMetadata.ChannelID,
		RequesterChannelName:        privateMetadata.ChannelName,
		RequesterRealName:           privateMetadata.RealName,
		RequireReason:               privateMetadata.RequireReason,
		SelectedRole:                privateMetadata.SelectedRole,
		SelectedChannelID:           privateMetadata.SelectedChannelID,
		SelectedChannelName:         privateMetadata.SelectedChannelName,
		SelectedStartDateOptionID:   privateMetadata.SelectedStartDateOptionID,
		SelectedStartDateOptionName: privateMetadata.SelectedStartDateOptionName,
		RequesterID:                 payload.User.ID,
		RequesterName:               payload.User.Name,
		TriggerID:                   payload.TriggerID,
		ViewHash:                    payload.View.Hash,
		ViewID:                      payload.View.ID,
		TTL:                         privateMetadata.TTL,
		StartDate:                   payload.View.State.Values.StartDateTimeBlock.AccessRequestStartDateSelect.SelectedDate,
	}, nil
}
