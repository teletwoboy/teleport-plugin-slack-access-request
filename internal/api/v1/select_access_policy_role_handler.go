package v1

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"log/slog"
	"net/http"
)

func (i *InteractionHandler) HandleAccessPolicyRoleSelection(payloadStr string, w http.ResponseWriter) {

	var payload slack.InteractionCallback
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pretty, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		slog.Error("Failed to format payload", "error", err)
	} else {
		fmt.Println("🔍 Formatted Slack Payload:\n", string(pretty))
	}
}
