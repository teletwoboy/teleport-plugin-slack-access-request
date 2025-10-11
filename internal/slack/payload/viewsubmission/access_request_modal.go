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

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
)

type AccessRequestModalPayload struct {
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
					AccessRequestReasonInput struct {
						Type  string `json:"type"`
						Value string `json:"value,omitempty"`
					} `json:"access_request_reason_input"`
				} `json:"reason_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`

	Email string
}

type AccessRequestModal struct {
	RequesterChannelID               string
	RequesterChannelName             string
	RequesterRealName                string
	RequireReason                    bool
	RequesterID                      string
	RequesterName                    string
	SelectedRole                     string
	SelectedChannelID                string
	SelectedChannelName              string
	SelectedStartDateOptionID        string
	SelectedStartDateOptionName      string
	TTL                              string
	SelectedStartDate                string
	SelectedStartTime                string
	SelectedAccessDurationOptionID   string
	SelectedAccessDurationOptionName string
	SelectedAccessDurationDate       string
	SelectedAccessDurationTime       string
	RequestTTL                       string
	SelectedRequestTTLOptionID       string
	SelectedRequestTTLOptionName     string
	SelectedRequestTTLDate           string
	SelectedRequestTTLTime           string
	TriggerID                        string
	ViewHash                         string
	ViewID                           string
	// new fields
	SelectedStartDateTime          time.Time
	SelectedAccessDurationDateTime time.Time
	SelectedRequestTTLDateTime     time.Time
	Reason                         string
}

func ParseAccessRequestModal(payloadStr string) (*AccessRequestModal, error) {
	var payload AccessRequestModalPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", payloadStr)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata accessrequest.SummaryPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", strPrivateMetadata)
	}

	return &AccessRequestModal{
		RequesterChannelID:               privateMetadata.ChannelID,
		RequesterChannelName:             privateMetadata.ChannelName,
		RequesterRealName:                privateMetadata.RealName,
		RequireReason:                    privateMetadata.RequireReason,
		RequesterID:                      payload.User.ID,
		RequesterName:                    payload.User.Name,
		SelectedRole:                     privateMetadata.SelectedRole,
		SelectedChannelID:                privateMetadata.SelectedChannelID,
		SelectedChannelName:              privateMetadata.SelectedChannelName,
		SelectedStartDateOptionID:        privateMetadata.SelectedStartDateOptionID,
		SelectedStartDateOptionName:      privateMetadata.SelectedStartDateOptionName,
		SelectedStartDate:                privateMetadata.SelectedStartDate,
		SelectedStartTime:                privateMetadata.SelectedStartTime,
		SelectedAccessDurationOptionID:   privateMetadata.SelectedAccessDurationOptionID,
		SelectedAccessDurationOptionName: privateMetadata.SelectedAccessDurationOptionName,
		SelectedAccessDurationDate:       privateMetadata.SelectedAccessDurationDate,
		SelectedAccessDurationTime:       privateMetadata.SelectedAccessDurationTime,
		SelectedRequestTTLOptionID:       privateMetadata.SelectedRequestTTLOptionID,
		SelectedRequestTTLOptionName:     privateMetadata.SelectedRequestTTLOptionName,
		SelectedRequestTTLDate:           privateMetadata.SelectedRequestTTLDate,
		SelectedRequestTTLTime:           privateMetadata.SelectedRequestTTLTime,
		TriggerID:                        payload.TriggerID,
		ViewHash:                         payload.View.Hash,
		ViewID:                           payload.View.ID,
		SelectedStartDateTime:            time.Time{},
		SelectedAccessDurationDateTime:   time.Time{},
		SelectedRequestTTLDateTime:       time.Time{},
		Reason:                           payload.View.State.Values.ReasonBlock.AccessRequestReasonInput.Value,
	}, nil
}
