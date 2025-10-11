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

package accesspolicy

import (
	"context"
	"net/http"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/api/res"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/metric/telemetry"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/builder/modal/accesspolicy"
	blockactions "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util"
)

func (h *Handler) HandleUserSelection(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.SlackTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.APolicyUserSelection)
	defer span.End()

	// 1. 값 준비
	payload, err := blockactions.ParseUserSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. 모달 생성하기
	builder := accesspolicy.NewFourthStepStartDateBuilder(payload)

	// 3. 모달 푸시하기
	if err := h.Services.Slack.PushModalContext(ctx, payload.TriggerID, builder); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
