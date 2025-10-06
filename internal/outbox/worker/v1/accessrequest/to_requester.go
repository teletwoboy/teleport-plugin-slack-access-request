package accessrequest

import (
	"context"
	"encoding/json"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/outbox/model/accessrequest"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
)

func (h *Handler) HandleToRequesterOutbox(ctx context.Context, ob *model.Outbox) error {
	ctx, cancel := context.WithTimeout(ctx, constant.ProcessingTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.WorkerAccessRequestToRequester)
	defer span.End()

	// payload 역직렬화
	var payload accessrequest.ToRequesterPayload
	if err := json.Unmarshal([]byte(ob.Payload), &payload); err != nil {
		return err
	}
	requesterChannelID := payload.RequesterChannelID
	slackUserID := payload.SlackUserID

	// Access Request 정보 가져오기
	aRequest, err := h.Services.Teleport.GetAccessRequestByAccessRequestID(ctx, ob.AggregateID)
	if err != nil {
		return err
	}

	// Slack User 정보 가져오기
	slackUser, err := h.Services.Slack.GetUserBySlackUserID(ctx, slackUserID)
	if err != nil {
		return err
	}

	// 메시지 전송하기
	builder := message.NewAccessRequestToRequesterBuilder(aRequest, slackUser)
	_, _, err = h.Services.Slack.PostMessageContext(ctx, requesterChannelID, builder)
	if err != nil {
		return err
	}

	// Done 처리하기
	return h.Services.Outbox.MarkDone(ctx, ob)
}
