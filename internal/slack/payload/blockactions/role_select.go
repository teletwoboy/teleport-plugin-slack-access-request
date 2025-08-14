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

package blockactions

import (
	"encoding/json"
	"fmt"
)

type RoleSelectPayload struct {
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
				RoleBlock struct {
					AccessRequestRoleSelect struct {
						Type string `json:"type"`

						SelectedOption *struct {
							Value string `json:"value"`

							Text struct {
								Type string `json:"type"`
								Text string `json:"text"`
							} `json:"text"`
						} `json:"selected_option,omitempty"`
					} `json:"access_request_role_select"`
				} `json:"role_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type RoleSelectPrivateMetadataPayload struct {
	ChannelID     string `json:"channel_id"`
	ChannelName   string `json:"channel_name"`
	RealName      string `json:"real_name"`
	RequireReason bool   `json:"require_reason"`
}

type RoleSelect struct {
	RequesterChannelID   string
	RequesterChannelName string
	RequesterRealName    string
	RequireReason        bool
	RequesterID          string
	RequesterName        string
	TriggerID            string
	ViewHash             string
	ViewID               string
	// new fields
	Role string
}

func ParseRoleSelect(payloadStr string) (*RoleSelect, error) {
	var payload RoleSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", payloadStr)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata RoleSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", strPrivateMetadata)
	}
	return &RoleSelect{
		RequesterChannelID:   privateMetadata.ChannelID,
		RequesterChannelName: privateMetadata.ChannelName,
		RequesterRealName:    privateMetadata.RealName,
		RequireReason:        privateMetadata.RequireReason,
		RequesterID:          payload.User.ID,
		RequesterName:        payload.User.Name,
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
		Role:                 payload.View.State.Values.RoleBlock.AccessRequestRoleSelect.SelectedOption.Value,
	}, nil
}
