package accessrequest

import (
	"context"
	"encoding/json"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/outbox/model/accessrequest"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	teleport "teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
)

func (h *Handler) HandleAutoReviewToReviewerOutbox(ctx context.Context, ob *model.Outbox) error {
	ctx, cancel := context.WithTimeout(ctx, constant.ProcessingTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.WorkerAccessRequestAutoReviewToReviewer)
	defer span.End()

	var payload accessrequest.AutoReviewToReviewerPayload
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

	// Teleport에 state 업데이트 요청하기
	updateBuilder := teleport.NewUpdateBuilder(aRequest.Name, aRequest.State, "Auto Review")
	if err := h.Services.Teleport.SubmitAccessRequestState(ctx, updateBuilder); err != nil {
		return err
	}

	// 메시지 전송하기
	builder := message.NewAutoReviewToReviewersBuilder(aRequest, aReview, slackUser, aPolicy)
	_, _, err = h.Services.Slack.PostMessageContext(ctx, aRequest.ReviewChannelID, builder)
	if err != nil {
		return err
	}

	// Done 처리하기
	return h.Services.Outbox.MarkDone(ctx, ob)
}
