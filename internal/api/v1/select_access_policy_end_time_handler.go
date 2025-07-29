package v1

import (
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
)

func (i *InteractionHandler) HandleAccessPolicyEndTimeSelection(payloadStr string, w http.ResponseWriter) {
	payload, err := blockactions.ParseAccessPolicyEndTimeSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 기존 모달에서 Access Policy Effect 추가하기
	builder := modal.NewSelectEffectBuilder(payload)

	// 모달 업데이트하기
	if err := i.Services.Slack.UpdateModal(builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
