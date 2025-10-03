package accesspolicy

import (
	"context"
	"encoding/json"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/util"

	"golang.org/x/sync/errgroup"
)

func (h *Handler) HandleCreationOutbox(ctx context.Context, ob *model.Outbox) error {
	ctx, cancel := context.WithTimeout(ctx, util.Timeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.WorkerAccessPolicyCreation)
	defer span.End()

	// 1. payload 역직렬화
	var payload model.AccessPolicyCreationPayload
	if err := json.Unmarshal([]byte(ob.Payload), &payload); err != nil {
		return err
	}
	accessPolicy := payload.AccessPolicy
	requesterRealName := payload.RequesterRealName

	// 2. 메시지 전송하기
	builder := message.NewAccessPolicySubmissionBuilder(accessPolicy, requesterRealName)
	channelID, timestamp, err := h.Services.Slack.PostMessageContext(ctx, accessPolicy.InputChannelID, builder)
	if err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		// 3. Pin 처리하기
		return h.Services.Slack.AddPinContext(gCtx, channelID, timestamp)
	})

	g.Go(func() error {
		// 4. timestamp 업데이트하기
		accessPolicy.UpdateTimestamp(timestamp)
		return h.Services.Policy.UpdateAccessPolicyMsgTs(gCtx, accessPolicy)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return h.Services.Outbox.MarkDone(ctx, ob)
}
