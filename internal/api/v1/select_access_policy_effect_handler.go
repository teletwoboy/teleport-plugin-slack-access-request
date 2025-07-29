package v1

import (
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
)

func (i *InteractionHandler) HandleAccessPolicyEffectSelection(payloadStr string, w http.ResponseWriter) {
	payload, err := blockactions.ParseAccessPolicyEffectSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 3페이지 요약으로 넘어가기
	builder := modal.NewSummaryBuilder(payload)

	// 모달 푸시 하기
	if err := i.Services.Slack.PushModal(payload.TriggerID, builder); err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
