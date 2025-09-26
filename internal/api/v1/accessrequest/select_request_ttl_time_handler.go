package accessrequest

import (
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	accessrequestmodal "teleport-plugin-slack-access-request/internal/slack/builder/modal/accessrequest"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/container"
	"time"
)

func (h *Handler) HandleRequestTTLTimeSelection(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	timezone := user.Slack.TimeZone
	role := payload.SelectedRole

	rTDate := payload.SelectedRequestTTLDate
	rTTime := payload.RequestTTLTime
	rT, err := util.ParseDateTimeInLocation(rTDate, rTTime, timezone)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	var sD time.Time
	if payload.SelectedStartDateOptionID == util.ARequestStartDateSecondOption {
		sDate := payload.SelectedStartDate
		sTime := payload.SelectedStartTime
		sD, err = util.ParseDateTimeInLocation(sDate, sTime, timezone)
		if err != nil {
			res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}
	}

	var aD time.Time
	if payload.SelectedAccessDurationOptionID == util.ARequestAccessDurationSecondOption {
		aDDate := payload.SelectedAccessDurationDate
		aDTime := payload.SelectedAccessDurationTime
		aD, err = util.ParseDateTimeInLocation(aDDate, aDTime, timezone)
		if err != nil {
			res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}
	}

	v3Builder := accessrequest.NewV3DryRunBuilder(role, sD, aD, rT, user.Teleport)
	_, err = h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	builder := accessrequestmodal.NewSixthStepBuilder(payload)

	// 3. 모달 업데이트하기
	if err := h.Services.Slack.UpdateModalContext(ctx, builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
