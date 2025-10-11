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

func (h *Handler) HandleRequestTTLTimeSelection(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.SlackTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.ARequestRequestTTLTimeSelection)
	defer span.End()

	// 1. 값 준비
	payload, err := blockactions.ParseRequestTTLTimeSelect(payloadStr)
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

	rTDate := payload.SelectedRequestTTLDate
	rTTime := payload.RequestTTLTime
	rT, err := util.ParseDateTimeInLocation(rTDate, rTTime, timezone)
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

	var aD time.Time
	if payload.SelectedAccessDurationOptionID == util.ARequestAccessDurationSecondOption {
		aDDate := payload.SelectedAccessDurationDate
		aDTime := payload.SelectedAccessDurationTime
		aD, err = util.ParseDateTimeInLocation(aDDate, aDTime, timezone)
		if err != nil {
			res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
			return
		}
	}

	v3Builder := accessrequest.NewV3DryRunBuilder(role, sD, aD, rT, user.Teleport)
	_, err = h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}
	builder := accessrequestmodal.NewSixthStepBuilder(payload)

	// 3. 모달 업데이트하기
	if err := h.Services.Slack.UpdateModalContext(ctx, builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
