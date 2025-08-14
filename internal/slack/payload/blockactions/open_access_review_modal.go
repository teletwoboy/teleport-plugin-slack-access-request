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

type OpenAccessReviewModalPayload struct {
	Type      string `json:"type"`
	TriggerID string `json:"trigger_id"`

	Channel struct {
		ID string `json:"id"`
	} `json:"channel"`

	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`

	Actions []struct {
		ActionID string `json:"action_id"`
		BlockID  string `json:"block_id"`
		Type     string `json:"type"`

		Text struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"text"`

		Value string `json:"value"`
	} `json:"actions"`
}

type OpenAccessReviewModal struct {
	AccessRequestName string
	ReviewerChannelID string
	ReviewerID        string
	ReviewerName      string
	TriggerID         string
}

func ParseOpenAccessReviewModalPayload(payloadStr string) (*OpenAccessReviewModal, error) {
	var payload OpenAccessReviewModalPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("failed to parse open access review modal payload: %w", err)
	}
	return &OpenAccessReviewModal{
		AccessRequestName: payload.Actions[0].Value,
		ReviewerChannelID: payload.Channel.ID,
		ReviewerID:        payload.User.ID,
		ReviewerName:      payload.User.Name,
		TriggerID:         payload.TriggerID,
	}, nil
}
