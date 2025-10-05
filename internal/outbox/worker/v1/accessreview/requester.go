package accessreview

import (
	"context"
	"encoding/json"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
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
