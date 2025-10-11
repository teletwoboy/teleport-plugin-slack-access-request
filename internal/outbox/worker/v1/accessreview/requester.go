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

package accessreview

import (
	"context"
	"encoding/json"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/metric/telemetry"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/constant"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/model"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/builder/message"
)

func (h *Handler) HandleRequesterOutbox(ctx context.Context, ob *model.Outbox) error {
	ctx, cancel := context.WithTimeout(ctx, constant.ProcessingTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.WorkerAccessReviewRequester)
	defer span.End()

	// 1. payload 역직렬화
	var payload model.AccessReviewRequesterPayload
	if err := json.Unmarshal([]byte(ob.Payload), &payload); err != nil {
		return err
	}
	aRequest := payload.AccessRequest
	aReview := payload.AccessReview
	reqSlackUser := payload.Requester
	revSlackUser := payload.Reviewer

	// 1. Requester 에게 처리되었음을 알림
	builder := message.NewAccessReviewToRequesterBuilder(aRequest, aReview, reqSlackUser, revSlackUser)
	_, _, err := h.Services.Slack.PostMessageContext(ctx, aRequest.InputChannelID, builder)
	if err != nil {
		return err
	}
	return h.Services.Outbox.MarkDone(ctx, ob)
}
