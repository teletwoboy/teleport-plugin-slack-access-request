package accesspolicy

import (
	"context"
	"encoding/json"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
)

func (h *Handler) HandleDeletionOutbox(ctx context.Context, ob *model.Outbox) error {
	ctx, cancel := context.WithTimeout(ctx, constant.ProcessingTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.WorkerAccessPolicyDeletion)
	defer span.End()

	// 1. payload 역직렬화
	var payload model.AccessPolicyDeletion
	if err := json.Unmarshal([]byte(ob.Payload), &payload); err != nil {
		return err
	}
	inputChannelID := payload.InputChannelID
	ts := payload.MessageTimestamp

	// 핀 제거 처리하기
	if err := h.Services.Slack.RemovePinContext(ctx, inputChannelID, ts); err != nil {
		return err
	}

	// Done 처리하기
	return h.Services.Outbox.MarkDone(ctx, ob)
}
