package v1

import (
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
)

func (i *InteractionHandler) HandleAccessPolicyUserSelection(payloadStr string, w http.ResponseWriter) {
	// 1. 값 준비
	payload, err := blockactions.ParseAccessPolicyUserSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. 모달 생성하기
	builder := modal.NewSelectStartDateBuilder(payload)

	// 3. 모달 푸시하기
	if err := i.Services.Slack.PushModal(payload.TriggerID, builder); err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
