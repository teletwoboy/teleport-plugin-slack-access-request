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

package accessrequest

import (
	"context"
	"fmt"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal/accessrequest"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

func (h *Handler) HandleRoleSelection(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := tracer.Start(ctx, telemetry.ARequestRoleSelection)
	defer span.End()

	// 1. 값 준비
	payload, err := blockactions.ParseRoleSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. 검증
	if err := h.verifyRoleSelection(ctx, payload); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 3. Access Request Modal을 만들기 위한 데이터를 수집한다.
	//	  1. reviewersChannels
	channels, err := h.Services.Slack.FetchReviewersChannelByRole(ctx, payload.Role)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 4. 모달 Builder를 만든다.
	builder := accessrequest.NewSecondStepBuilder(channels, payload)

	// 5. 모달을 보낸다.
	if err := h.Services.Slack.UpdateModalContext(ctx, builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) verifyRoleSelection(ctx context.Context, payload *blockactions.RoleSelect) error {
	slackVerifier := verifier.NewSlack(h.Services.Slack)
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	if err := slackVerifier.VerifyUserExistsByID(ctx, payload.RequesterID, payload.RequesterName); err != nil {
		return fmt.Errorf("failed to verify slack user exists by ID: %w", err)
	}

	//    2. 해당 유저가 Request Channel 에 있는 사람이 맞는가?
	if err := slackVerifier.VerifyUserExistsInChannelByID(ctx, payload.RequesterID, payload.RequesterChannelID); err != nil {
		return fmt.Errorf("failed to verify slack user exists in channel by ID: %w", err)
	}
	return nil
}
