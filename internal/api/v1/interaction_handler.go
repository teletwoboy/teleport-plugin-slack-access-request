package v1

import (
	"encoding/json"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/v1/accesspolicy"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload"
	"teleport-plugin-slack-access-request/internal/util/container"

	slackapi "github.com/slack-go/slack"
)

type InteractionHandler struct {
	DB       *database.DB
	Clients  *container.Clients
	Repos    *container.Repositories
	Services *container.Services
}

func NewInteractionHandler(db *database.DB, c *container.Clients, r *container.Repositories, s *container.Services) *InteractionHandler {
	return &InteractionHandler{
		DB:       db,
		Clients:  c,
		Repos:    r,
		Services: s,
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

	apHandler := accesspolicy.NewHandler(i.Services)
	switch slackapi.InteractionType(callback.Type) {
	case slackapi.InteractionTypeBlockActions:
		switch callback.Actions[0].ActionID {
		case modal.APChannelOptionBlockActionID:
			apHandler.HandleChannelSelection(payloadStr, w)
		case modal.APRoleOptionBlockActionID:
			apHandler.HandleRoleSelection(payloadStr, w)
		case modal.APUserOptionBlockActionID:
			apHandler.HandleUserSelection(payloadStr, w)
		case modal.APTimeZoneOptionBlockActionID:
			apHandler.HandleTimeZoneSelection(payloadStr, w)
		case modal.APStartDateBlockActionID:
			apHandler.HandleStartDateSelection(payloadStr, w)
		case modal.APStartTimeBlockActionID:
			apHandler.HandleStartTimeSelection(payloadStr, w)
		case modal.APEndDateBlockActionID:
			apHandler.HandleEndDateSelection(payloadStr, w)
		case modal.APEndTimeBlockActionID:
			apHandler.HandleEndTimeSelection(payloadStr, w)
		case modal.APAllowButtonBlockActionID:
			apHandler.HandleEffectSelection(payloadStr, w)
		case modal.APDenyButtonBlockActionID:
			apHandler.HandleEffectSelection(payloadStr, w)
		case "open_access_request_review_modal":
			i.HandleOpenAccessReviewModal(payloadStr, w)
		case "role_select":
			i.HandleAccessRoleModalSelection(payloadStr, w)
		}
	case slackapi.InteractionTypeViewSubmission:
		switch callback.View.CallbackID {
		case "access_request_modal":
			i.HandleAccessRequestModalSubmission(payloadStr, w)
		case "access_review_modal":
			i.HandleAccessReviewModalSubmission(payloadStr, w)
		case modal.APCallBackID:
			i.SubmitAccessPolicyModalHandler(payloadStr, w)
		}
	default:
		http.Error(w, "unsupported interaction type", http.StatusBadRequest)
		return
	}
}
