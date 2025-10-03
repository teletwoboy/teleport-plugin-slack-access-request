package accessrequest

import (
	"context"
	"encoding/json"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/outbox/model/accessrequest"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/util"
)

func (h *Handler) HandleAutoReviewToRequesterOutbox(ctx context.Context, ob *model.Outbox) error {
	ctx, cancel := context.WithTimeout(ctx, util.Timeout)
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
	if err := h.Services.Outbox.MarkDone(ctx, ob); err != nil {
		return err
	}
	return nil
}
