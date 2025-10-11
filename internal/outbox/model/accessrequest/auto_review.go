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

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/constant"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/model"
	policymodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/policy/models"
)

type AutoReviewPayload struct {
	AccessPolicyID     int32
	RequesterID        string
	RequesterChannelID string
	SelectedChannelID  string
	SlackUserID        int32
	UserID             int32
	Username           string
}

func NewOutboxWithAutoReview(
	policy *policymodels.AccessPolicy,
	ob *model.Outbox,
	payload JudgementPayload,
) (*model.Outbox, error) {
	p := AutoReviewPayload{
		AccessPolicyID:     policy.AccessPolicyID,
		RequesterID:        payload.RequesterID,
		RequesterChannelID: payload.RequesterChannelID,
		SelectedChannelID:  payload.SelectedChannelID,
		SlackUserID:        payload.SlackUserID,
		UserID:             payload.UserID,
		Username:           payload.Username,
	}
	marshaled, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access request auto review payload: %w", err)
	}

	outbox := &model.Outbox{
		EventType:     constant.AccessRequestAutoReview,
		AggregateType: constant.AccessRequest,
		AggregateID:   ob.AggregateID,
		Payload:       string(marshaled),
		Status:        constant.Pending,
	}
	return outbox, nil
}
