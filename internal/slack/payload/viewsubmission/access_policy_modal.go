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

package viewsubmission

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
)

type AccessPolicyModalPayload struct {
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
				ReasonBlock struct {
					AccessPolicyReasonInput struct {
						Type  string `json:"type"`
						Value string `json:"value,omitempty"`
					} `json:"access_policy_reason_input"`
				} `json:"reason_block"`

				TitleBlock struct {
					AccessPolicyTitleInput struct {
						Type  string `json:"type"`
						Value string `json:"value,omitempty"`
					} `json:"access_policy_title_input"`
				} `json:"title_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type AccessPolicyModal struct {
	RequesterChannelID   string
	RequesterChannelName string
	RequesterRealName    string
	RequesterID          string
	RequesterName        string
	SelectedChannelID    string
	SelectedChannelName  string
	SelectedRole         string
	SelectedRoleName     string
	SelectedUserID       string
	SelectedRealName     string
	SelectedStartDate    time.Time
	SelectedEndDate      time.Time
	SelectedEffect       string
	TriggerID            string
	ViewHash             string
	ViewID               string
	// new fields
	Title  string
	Reason string
}

func ParseAccessPolicyModal(payloadStr string) (*AccessPolicyModal, error) {
	var payload AccessPolicyModalPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %w", err)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata accesspolicy.SummaryPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid payload format: %w", err)
	}
	return &AccessPolicyModal{
		RequesterChannelID:   privateMetadata.ChannelID,
		RequesterChannelName: privateMetadata.ChannelName,
		RequesterRealName:    privateMetadata.RealName,
		RequesterID:          payload.User.ID,
		RequesterName:        payload.User.Name,
		SelectedChannelID:    privateMetadata.SelectedChannelID,
		SelectedChannelName:  privateMetadata.SelectedChannelName,
		SelectedRole:         privateMetadata.SelectedRole,
		SelectedRoleName:     privateMetadata.SelectedRoleName,
		SelectedUserID:       privateMetadata.SelectedUserID,
		SelectedRealName:     privateMetadata.SelectedRealName,
		SelectedStartDate:    privateMetadata.SelectedStartDate,
		SelectedEndDate:      privateMetadata.SelectedEndDate,
		SelectedEffect:       privateMetadata.SelectedEffect,
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
		Title:                payload.View.State.Values.TitleBlock.AccessPolicyTitleInput.Value,
		Reason:               payload.View.State.Values.ReasonBlock.AccessPolicyReasonInput.Value,
	}, nil
}
