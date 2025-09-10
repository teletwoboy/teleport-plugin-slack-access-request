package accessrequest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	accessrequestmodal "teleport-plugin-slack-access-request/internal/slack/builder/modal/accessrequest"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"
)

func (h *Handler) HandleRequestTTLOptionSelection(payloadStr string, w http.ResponseWriter) {
	// 1. 값 준비
	payload, err := blockactions.ParseRequestTTLOptionSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	pretty, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		slog.Error("Failed to format payload", "error", err)
	} else {
		fmt.Println("🔍 Formatted Slack Payload:\n", string(pretty))
	}

	var builder modal.Builder
	switch payload.RequestTTLOptionID {
	case util.ARequestRequestTTLFirstOption:
	case util.ARequestRequestTTLSecondOption:
		// 1. Date 모달 빌드하기
		builder = accessrequestmodal.NewFifthStepDateBuilder(payload)
	}

	if err := h.Services.Slack.UpdateModal(builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
