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
	slackmodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/models"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	teleportmodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/teleport/models"
	usermodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/user/models"
)

type SubmissionPayload struct {
	Payload     *viewsubmission.AccessRequestModal
	SlackUserID int32
	UserID      int32
	Username    string
}

func NewOutboxWithSubmission(
	p *viewsubmission.AccessRequestModal,
	slackUser *slackmodels.User,
	teleportUser *teleportmodels.User,
	user *usermodels.User,
	aRequestID int32,
) (*model.Outbox, error) {
	payload := SubmissionPayload{
		Payload:     p,
		SlackUserID: slackUser.SlackUserID,
		UserID:      user.UserID,
		Username:    teleportUser.Username,
	}
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access request creation payload: %w", err)
	}

	outbox := &model.Outbox{
		EventType:     constant.AccessRequestSubmission,
		AggregateType: constant.AccessRequest,
		AggregateID:   aRequestID,
		Payload:       string(marshaled),
		Status:        constant.Pending,
	}
	return outbox, nil
}
