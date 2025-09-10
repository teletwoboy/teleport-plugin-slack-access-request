package accessrequest

import (
	"context"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	accessrequestmodal "teleport-plugin-slack-access-request/internal/slack/builder/modal/accessrequest"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/container"
	"time"
)

func (h *Handler) HandleStartTimeSelection(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 값 준비
	payload, err := blockactions.ParseStartTimeSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. DryRun 으로 데이터 정합성 판단
	user, err := container.NewUsers(ctx, h.Services, payload.RequesterID)
	if err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	timezone := user.Slack.TimeZone
	role := payload.SelectedRole
	sDate := payload.SelectedStartDate
	sTime := payload.StartTime
	sD, err := util.ParseDateTimeInLocation(sDate, sTime, timezone)

	v3Builder := accessrequest.NewV3DryRunBuilder(role, sD, time.Time{}, user.Teleport)
	_, err = h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
	if err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 3. 4단계 모달 빌드하기
	builder := accessrequestmodal.NewFourthStepBuilder(payload)

	// 4. 모달 업데이트하기
	if err := h.Services.Slack.UpdateModal(builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
