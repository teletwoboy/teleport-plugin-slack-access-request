package res

import (
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
)

func ErrorMessageToSlack(s slack.Service, channelID string, err error, w http.ResponseWriter) {
	slog.Error("operation failed", "err", err)
	msg := message.NewErrorBuilder(err)

	_, _, postErr := s.PostMessage(channelID, msg)
	if postErr != nil {
		slog.Error("failed to post msg to slack", "err", postErr)
		http.Error(w, "failed to post msg to slack", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
