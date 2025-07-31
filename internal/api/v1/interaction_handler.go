package v1

import (
	"encoding/json"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/v1/accesspolicy"
	"teleport-plugin-slack-access-request/internal/api/v1/accessrequest"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/slack/payload"
	"teleport-plugin-slack-access-request/internal/util"
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

	switch slackapi.InteractionType(callback.Type) {
	case slackapi.InteractionTypeBlockActions:
		i.routeInteractionTypeBlockActions(callback, payloadStr, w)
	case slackapi.InteractionTypeViewSubmission:
		switch callback.View.CallbackID {
		case util.ARequestCallBackID:
			i.HandleAccessRequestModalSubmission(payloadStr, w)
		case "access_review_modal":
			i.HandleAccessReviewModalSubmission(payloadStr, w)
		case util.APolicyCallBackID:
			i.SubmitAccessPolicyModalHandler(payloadStr, w)
		}
	default:
		http.Error(w, "unsupported interaction type", http.StatusBadRequest)
		return
	}
}

func (i *InteractionHandler) routeInteractionTypeBlockActions(callback payload.Callback, payloadStr string, w http.ResponseWriter) {
	apHandler := accesspolicy.NewHandler(i.Services)
	arHandler := accessrequest.NewHandler(i.Services)
	switch callback.Actions[0].ActionID {
	case util.APolicyChanOptionBlockActionID:
		apHandler.HandleChannelSelection(payloadStr, w)
	case util.APolicyRoleOptionBlockActionID:
		apHandler.HandleRoleSelection(payloadStr, w)
	case util.APolicyUserOptionBlockActionID:
		apHandler.HandleUserSelection(payloadStr, w)
	case util.APolicyStartDateBlockActionID:
		apHandler.HandleStartDateSelection(payloadStr, w)
	case util.APolicyStartTimeBlockActionID:
		apHandler.HandleStartTimeSelection(payloadStr, w)
	case util.APolicyEndDateBlockActionID:
		apHandler.HandleEndDateSelection(payloadStr, w)
	case util.APolicyEndTimeBlockActionID:
		apHandler.HandleEndTimeSelection(payloadStr, w)
	case util.APolicyAllowButtonBlockActionID:
		apHandler.HandleEffectSelection(payloadStr, w)
	case util.APolicyDenyButtonBlockActionID:
		apHandler.HandleEffectSelection(payloadStr, w)
	case "open_access_request_review_modal":
		i.HandleOpenAccessReviewModal(payloadStr, w)
	case util.ARequestRoleOptionBlockActionID:
		arHandler.HandleRoleSelection(payloadStr, w)
	case util.ARequestChannelOptionBlockActionID:
		arHandler.HandleChannelSelection(payloadStr, w)
	}
}
