/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"encoding/json"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/v1/accesspolicy"
	"teleport-plugin-slack-access-request/internal/api/v1/accessrequest"
	"teleport-plugin-slack-access-request/internal/api/v1/accessreview"
	"teleport-plugin-slack-access-request/internal/slack/payload"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/container"

	slackapi "github.com/slack-go/slack"
)

type InteractionHandler struct {
	aPolicy  *accesspolicy.Handler
	aRequest *accessrequest.Handler
	aReview  *accessreview.Handler
	services *container.Services
}

func NewInteractionHandler(aPolicy *accesspolicy.Handler, aRequest *accessrequest.Handler, aReview *accessreview.Handler, s *container.Services) *InteractionHandler {
	return &InteractionHandler{
		aPolicy:  aPolicy,
		aRequest: aRequest,
		aReview:  aReview,
		services: s,
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
			i.aRequest.HandleModalSubmission(payloadStr, w)
		case "access_review_modal":
			i.aReview.HandleModalSubmission(payloadStr, w)
		case util.APolicyCallBackID:
			i.aPolicy.HandleModalSubmission(payloadStr, w)
		}
	default:
		http.Error(w, "unsupported interaction type", http.StatusBadRequest)
		return
	}
}

func (i *InteractionHandler) routeInteractionTypeBlockActions(callback payload.Callback, payloadStr string, w http.ResponseWriter) {
	switch callback.Actions[0].ActionID {
	case util.APolicyChanOptionBlockActionID:
		i.aPolicy.HandleChannelSelection(payloadStr, w)
	case util.APolicyRoleOptionBlockActionID:
		i.aPolicy.HandleRoleSelection(payloadStr, w)
	case util.APolicyUserOptionBlockActionID:
		i.aPolicy.HandleUserSelection(payloadStr, w)
	case util.APolicyStartDateBlockActionID:
		i.aPolicy.HandleStartDateSelection(payloadStr, w)
	case util.APolicyStartTimeBlockActionID:
		i.aPolicy.HandleStartTimeSelection(payloadStr, w)
	case util.APolicyEndDateBlockActionID:
		i.aPolicy.HandleEndDateSelection(payloadStr, w)
	case util.APolicyEndTimeBlockActionID:
		i.aPolicy.HandleEndTimeSelection(payloadStr, w)
	case util.APolicyAllowButtonBlockActionID:
		i.aPolicy.HandleEffectSelection(payloadStr, w)
	case util.APolicyDenyButtonBlockActionID:
		i.aPolicy.HandleEffectSelection(payloadStr, w)
	case "open_access_request_review_modal":
		i.HandleOpenAccessReviewModal(payloadStr, w)
	case util.ARequestRoleOptionBlockActionID:
		i.aRequest.HandleRoleSelection(payloadStr, w)
	case util.ARequestChannelOptionBlockActionID:
		i.aRequest.HandleChannelSelection(payloadStr, w)
	}
}
