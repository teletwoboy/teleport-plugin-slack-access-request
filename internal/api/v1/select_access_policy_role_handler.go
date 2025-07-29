package v1

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
)

func (i *InteractionHandler) HandleAccessPolicyRoleSelection(payloadStr string, w http.ResponseWriter) {

	payload, err := blockactions.ParseAccessPolicyRoleSelect(payloadStr)
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
}
