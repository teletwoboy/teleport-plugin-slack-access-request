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
	teleportmodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/teleport/models"
)

type AutoReviewToRequesterPayload struct {
	AccessPolicyID  int32
	AccessRequestID int32
	AccessReviewID  int32
	SlackUserID     int32
}

func NewOutboxWithAutoReviewToRequester(
	accessPolicy *policymodels.AccessPolicy,
	accessRequest *teleportmodels.AccessRequest,
	accessReview *teleportmodels.AccessReview,
	slackUserID int32,
) (*model.Outbox, error) {
	p := AutoReviewToRequesterPayload{
		AccessPolicyID:  accessPolicy.AccessPolicyID,
		AccessRequestID: accessRequest.AccessRequestID,
		AccessReviewID:  accessReview.AccessReviewID,
		SlackUserID:     slackUserID,
	}
	marshaled, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access request creation payload: %w", err)
	}

	outbox := model.Outbox{
		EventType:     constant.AccessRequestAutoReviewToRequester,
		AggregateType: constant.AccessRequest,
		AggregateID:   accessRequest.AccessRequestID,
		Payload:       string(marshaled),
		Status:        constant.Pending,
	}
	return &outbox, nil
}
