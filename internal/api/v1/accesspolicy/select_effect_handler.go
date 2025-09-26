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
	"fmt"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal/accesspolicy"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
)

func (h *Handler) HandleEffectSelection(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := tracer.Start(ctx, telemetry.APolicyEffectSelection)
	defer span.End()

	// 1. 값 준비
	payload, err := blockactions.ParseEffectSelect(payloadStr)
	if err != nil {
		slog.Error("failed to parse payload", "err", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. 검증하기
	//    1. Start Date 가 End Date 보다 시간상 느린가?
	if payload.SelectedStartDate.After(payload.SelectedEndDate) {
		err := fmt.Errorf("start Date must be earlier than End Date. Please check your selection")
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 2. Summary 모달 생성하기
	builder := accesspolicy.NewSummaryBuilder(payload)

	// 3. 모달 푸시하기
	if err := h.Services.Slack.PushModalContext(ctx, payload.TriggerID, builder); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
