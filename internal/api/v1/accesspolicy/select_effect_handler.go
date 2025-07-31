package accesspolicy

import (
	"fmt"
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

	// 2. 검증하기
	//    1. Start Date 가 End Date 보다 시간상 느린가?
	if payload.SelectedStartDate.After(payload.SelectedEndDate) {
		err := fmt.Errorf("start Date must be earlier than End Date. Please check your selection")
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
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
