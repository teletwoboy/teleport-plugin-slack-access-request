package v1

import (
	"encoding/json"
	"net/http"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/slack/payload"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/user"

	slackapi "github.com/slack-go/slack"
)

type InteractionHandler struct {
	SlackSrv    slack.Service
	TeleportSrv teleport.Service
	UserSrv     user.Service
}

func NewInteractionHandler(s slack.Service, t teleport.Service, u user.Service) *InteractionHandler {
	return &InteractionHandler{
		SlackSrv:    s,
		TeleportSrv: t,
		UserSrv:     u,
	}
}

func (i *InteractionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	payloadStr := r.FormValue("payload")
	if payloadStr == "" {
		http.Error(w, "missing payload", http.StatusBadRequest)
		return
	}

	var callback payload.Callback
	if err := json.Unmarshal([]byte(payloadStr), &callback); err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	switch slackapi.InteractionType(callback.Type) {
	case slackapi.InteractionTypeBlockActions:
		i.HandleOpenAccessReviewModal(payloadStr, w)
	case slackapi.InteractionTypeViewSubmission:
		switch callback.View.CallbackID {
		case "access_request_modal":
			i.HandleAccessRequestModalSubmission(payloadStr, w)
		case "access_review_modal":
			i.HandleAccessReviewModalSubmission(payloadStr, w)
		}
	default:
		http.Error(w, "unsupported interaction type", http.StatusBadRequest)
		return
	}
}
