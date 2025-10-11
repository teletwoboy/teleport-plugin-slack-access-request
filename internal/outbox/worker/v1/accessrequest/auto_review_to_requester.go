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
	"context"
	"encoding/json"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/metric/telemetry"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/constant"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/model"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/model/accessrequest"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/builder/message"
)

func (h *Handler) HandleAutoReviewToRequesterOutbox(ctx context.Context, ob *model.Outbox) error {
	ctx, cancel := context.WithTimeout(ctx, constant.ProcessingTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.WorkerAccessRequestAutoReviewToRequester)
	defer span.End()

	var payload accessrequest.AutoReviewToRequesterPayload
	if err := json.Unmarshal([]byte(ob.Payload), &payload); err != nil {
		return err
	}
	aPolicyID := payload.AccessPolicyID
	aRequestID := payload.AccessRequestID
	aReviewID := payload.AccessReviewID
	slackUserID := payload.SlackUserID

	aPolicy, err := h.Services.Policy.GetAccessPoliciesByAccessPolicyID(ctx, aPolicyID)
	if err != nil {
		return err
	}

	aRequest, err := h.Services.Teleport.GetAccessRequestByAccessRequestID(ctx, aRequestID)
	if err != nil {
		return err
	}

	aReview, err := h.Services.Teleport.GetAccessReviewByAccessReviewID(ctx, aReviewID)
	if err != nil {
		return err
	}

	slackUser, err := h.Services.Slack.GetUserBySlackUserID(ctx, slackUserID)
	if err != nil {
		return err
	}

	// 메시지 전송하기
	builder := message.NewAutoReviewToRequesterBuilder(aRequest, aReview, slackUser, aPolicy)
	_, _, err = h.Services.Slack.PostMessageContext(ctx, aRequest.InputChannelID, builder)
	if err != nil {
		return err
	}

	// Done 처리하기
	return h.Services.Outbox.MarkDone(ctx, ob)
}
