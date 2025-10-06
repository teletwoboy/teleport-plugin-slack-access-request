package accessrequest

import (
	"context"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	accessrequestmodal "teleport-plugin-slack-access-request/internal/slack/builder/modal/accessrequest"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"
)

func (h *Handler) HandleRequestTTLOptionSelection(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.SlackTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.ARequestRequestTTLOptionSelection)
	defer span.End()

	// 1. 값 준비
	payload, err := blockactions.ParseRequestTTLOptionSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	var builder modal.Builder
	switch payload.RequestTTLOptionID {
	case util.ARequestRequestTTLFirstOption:
		// 1. RequestTTLTimeSelect 인스턴스 만들기
		requestTTLTime := blockactions.NewRequestTTLTimeWithFirstOpt(payload)
		// 2. 6단계 모달 빌드하기
		builder = accessrequestmodal.NewSixthStepBuilder(requestTTLTime)
	case util.ARequestRequestTTLSecondOption:
		// 1. Date 모달 빌드하기
		builder = accessrequestmodal.NewFifthStepDateBuilder(payload)
	}

	if err := h.Services.Slack.UpdateModalContext(ctx, builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
