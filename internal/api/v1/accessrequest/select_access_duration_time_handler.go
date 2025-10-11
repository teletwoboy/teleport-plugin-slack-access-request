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
	"net/http"
	"time"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/api/res"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/metric/telemetry"
	accessrequestmodal "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/builder/modal/accessrequest"
	blockactions "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util/container"
)

func (h *Handler) HandleAccessDurationTimeSelection(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.SlackTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.ARequestAccessDurationTimeSelection)
	defer span.End()

	// 1. 값 준비
	payload, err := blockactions.ParseAccessDurationTimeSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	user, err := container.NewUsers(ctx, h.Services, payload.RequesterID)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}
	timezone := user.Slack.TimeZone
	role := payload.SelectedRole

	// 2. AccessDuration time.Time과 StartDate 조건부 time.Time 구하기
	aDDate := payload.SelectedAccessDurationDate
	aDTime := payload.AccessDurationTime
	aD, err := util.ParseDateTimeInLocation(aDDate, aDTime, timezone)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}

	var sD time.Time
	if payload.SelectedStartDateOptionID == util.ARequestStartDateSecondOption {
		sDate := payload.SelectedStartDate
		sTime := payload.SelectedStartTime
		sD, err = util.ParseDateTimeInLocation(sDate, sTime, timezone)
		if err != nil {
			res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
			return
		}
	}

	// 3. DryRun 요청하여 RequestTTL 구하기
	v3Builder := accessrequest.NewV3DryRunBuilder(role, sD, aD, time.Time{}, user.Teleport)
	submittedAccessRequest, err := h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}

	requestTTL, err := util.ParseTTLInLocation(submittedAccessRequest, timezone)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}

	// 4. 5단계 모달 빌드하기
	builder := accessrequestmodal.NewFifthStepBuilder(payload, requestTTL)

	// 5. 모달 업데이트하기
	if err := h.Services.Slack.UpdateModalContext(ctx, builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
