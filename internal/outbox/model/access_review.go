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

package model

import (
	"encoding/json"
	"fmt"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/constant"
	slackmodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/teleport/models"
)

type AccessReviewReviewerPayload struct {
	AccessRequest *teleportmodels.AccessRequest
	AccessReview  *teleportmodels.AccessReview
	Requester     *slackmodels.User
	Reviewer      *slackmodels.User
	MessageTs     string
}

func NewOutboxWithAccessReviewReviewer(
	aRequest *teleportmodels.AccessRequest,
	aReview *teleportmodels.AccessReview,
	requester *slackmodels.User,
	reviewer *slackmodels.User,
	messageTs string,
) (*Outbox, error) {
	payload := AccessReviewReviewerPayload{
		AccessRequest: aRequest,
		AccessReview:  aReview,
		Requester:     requester,
		Reviewer:      reviewer,
		MessageTs:     messageTs,
	}
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access review payload: %w", err)
	}

	outbox := &Outbox{
		EventType:     constant.AccessReviewReviewer,
		AggregateType: constant.AccessReview,
		AggregateID:   aReview.AccessReviewID,
		Payload:       string(marshaled),
		Status:        constant.Pending,
	}
	return outbox, nil
}

type AccessReviewRequesterPayload struct {
	AccessRequest *teleportmodels.AccessRequest
	AccessReview  *teleportmodels.AccessReview
	Requester     *slackmodels.User
	Reviewer      *slackmodels.User
}

func NewOutboxWithAccessReviewRequester(
	aRequest *teleportmodels.AccessRequest,
	aReview *teleportmodels.AccessReview,
	requester *slackmodels.User,
	reviewer *slackmodels.User,
) (*Outbox, error) {
	payload := AccessReviewRequesterPayload{
		AccessRequest: aRequest,
		AccessReview:  aReview,
		Requester:     requester,
		Reviewer:      reviewer,
	}
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access review payload: %w", err)
	}

	outbox := &Outbox{
		EventType:     constant.AccessReviewRequester,
		AggregateType: constant.AccessReview,
		AggregateID:   aReview.AccessReviewID,
		Payload:       string(marshaled),
		Status:        constant.Pending,
	}
	return outbox, nil
}
