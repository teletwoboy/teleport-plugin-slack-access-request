package accessrequest

import (
	"context"
	"fmt"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	accessrequestmodal "teleport-plugin-slack-access-request/internal/slack/builder/modal/accessrequest"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/container"
	"time"
)

func (h *Handler) HandleAccessDurationTimeSelection(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 값 준비
	payload, err := blockactions.ParseAccessDurationTimeSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	user, err := container.NewUsers(ctx, h.Services, payload.RequesterID)
	if err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	timezone := user.Slack.TimeZone
	role := payload.SelectedRole

	var builder modal.Builder
	switch payload.SelectedStartDateOptionID {
	case util.ARequestStartDateFirstOption: // 1. StartDate: Immediately
		// 1. AccessDurationDate/Time 값을 time.Time 으로 만들기
		aDDate := payload.SelectedAccessDurationDate
		aDTime := payload.AccessDurationTime
		fmt.Println(aDDate, aDTime)
		aD, err := util.ParseDateTimeInLocation(aDDate, aDTime, timezone)
		if err != nil {
			res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}

		// 2. DryRun 으로 Access Request 요청 후 반환되는 값의 RequestTTL 값 가져오기
		v3Builder := accessrequest.NewV3DryRunBuilder(role, time.Time{}, aD, user.Teleport)
		submittedAccessRequest, err := h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
		if err != nil {
			res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}

		requestTTL, err := util.ParseTTLInLocation(submittedAccessRequest, timezone)
		if err != nil {
			res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}

		// 3. 5단계 모달 빌드하기
		builder = accessrequestmodal.NewFifthStepBuilder(payload, requestTTL)
	case util.ARequestStartDateSecondOption: // 2. StartDate: Select DateTime
		// 1. StartDate/Time, AccessDurationDate/Time 값을 time.Time 으로 만들기
		sDate := payload.SelectedStartDate
		sTime := payload.SelectedStartTime
		sD, err := util.ParseDateTimeInLocation(sDate, sTime, timezone)
		if err != nil {
			res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}
		aDDate := payload.SelectedAccessDurationDate
		aDTime := payload.AccessDurationTime
		aD, err := util.ParseDateTimeInLocation(aDDate, aDTime, timezone)
		if err != nil {
			res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}

		// 2. DryRun 으로 Access Request 요청 후 반환되는 값의 RequestTTL 값 가져오기
		v3Builder := accessrequest.NewV3DryRunBuilder(role, sD, aD, user.Teleport)
		submittedAccessRequest, err := h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
		if err != nil {
			res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}

		requestTTL, err := util.ParseTTLInLocation(submittedAccessRequest, timezone)
		if err != nil {
			res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}

		// 3. 5단계 모달 빌드하기
		builder = accessrequestmodal.NewFifthStepBuilder(payload, requestTTL)
	}

	// 5. 모달 업데이트하기
	if err := h.Services.Slack.UpdateModal(builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
