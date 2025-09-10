package accessrequest

import (
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal/accessrequest"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
)

func (h *Handler) HandleStartDateSelection(payloadStr string, w http.ResponseWriter) {
	// 1. 값 준비
	payload, err := blockactions.ParseStartDateSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. 기존 모달에 Time 추가하기
	builder := accessrequest.NewThirdStepTimeBuilder(payload)

	// 3. 모달 업데이트하기
	if err := h.Services.Slack.UpdateModal(builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
