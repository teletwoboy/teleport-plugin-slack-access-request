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
)

type ToReviewerPayload struct {
	ReviewerChannelID string
	SlackUserID       int32
}

func NewOutboxWithToReviewer(
	ob *model.Outbox,
	payload JudgementPayload,
) (*model.Outbox, error) {
	p := ToReviewerPayload{
		ReviewerChannelID: payload.SelectedChannelID,
		SlackUserID:       payload.SlackUserID,
	}
	marshaled, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access request to reviewer payload: %w", err)
	}

	outbox := &model.Outbox{
		EventType:     constant.AccessRequestToReviewer,
		AggregateType: constant.AccessRequest,
		AggregateID:   ob.AggregateID,
		Payload:       string(marshaled),
		Status:        constant.Pending,
	}
	return outbox, nil
}
