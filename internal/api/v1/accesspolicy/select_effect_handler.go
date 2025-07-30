package accesspolicy

import (
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal/accesspolicy"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
)

func (h *Handler) HandleEffectSelection(payloadStr string, w http.ResponseWriter) {
	// 1. 값 준비
	payload, err := blockactions.ParseEffectSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. Summary 모달 생성하기
	builder := accesspolicy.NewSummaryBuilder(payload)

	// 3. 모달 푸시하기
	if err := h.Services.Slack.PushModal(payload.TriggerID, builder); err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
