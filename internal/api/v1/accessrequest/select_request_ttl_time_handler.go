package accessrequest

import (
	"context"
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

func (h *Handler) HandleRequestTTLTimeSelection(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 값 준비
	payload, err := blockactions.ParseRequestTTLTimeSelect(payloadStr)
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
	case util.ARequestStartDateFirstOption:
		switch payload.SelectedAccessDurationOptionID {
		case util.ARequestAccessDurationFirstOption:
			// 1. RequestTTLDate/Time 값을 time.Time 으로 만들기
			rTDate := payload.SelectedRequestTTLDate
			rTTime := payload.RequestTTLTime
			rT, err := util.ParseDateTimeInLocation(rTDate, rTTime, timezone)
			if err != nil {
				res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
				return
			}

			// 2. DryRun 으로 Access Request 요청하여 값 검증하기
			v3Builder := accessrequest.NewV3DryRunBuilder(role, time.Time{}, time.Time{}, rT, user.Teleport)
			_, err = h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
			if err != nil {
				res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
				return
			}

			// 3. 6단계 모달 빌드하기
			builder = accessrequestmodal.NewSixthStepBuilder(payload)
		case util.ARequestAccessDurationSecondOption:
			// 1. AccessDurationDate/Time, RequestTTLDate/Time 값을 time.Time 으로 만들기
			aDDate := payload.SelectedAccessDurationDate
			aDTime := payload.SelectedAccessDurationTime
			aD, err := util.ParseDateTimeInLocation(aDDate, aDTime, timezone)
			if err != nil {
				res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
				return
			}
			rTDate := payload.SelectedRequestTTLDate
			rTTime := payload.RequestTTLTime
			rT, err := util.ParseDateTimeInLocation(rTDate, rTTime, timezone)
			if err != nil {
				res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
				return
			}

			// 2. DryRun 으로 Access Request 요청하여 값 검증하기
			v3Builder := accessrequest.NewV3DryRunBuilder(role, time.Time{}, aD, rT, user.Teleport)
			_, err = h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
			if err != nil {
				res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
				return
			}

			// 3. 6단계 모달 빌드하기
			builder = accessrequestmodal.NewSixthStepBuilder(payload)
		}
	case util.ARequestStartDateSecondOption:
		switch payload.SelectedAccessDurationOptionID {
		case util.ARequestAccessDurationFirstOption:
			// 1. StartDate/Time, RequestTTLDate/Time 값을 time.Time 으로 만들기
			sDate := payload.SelectedStartDate
			sTime := payload.SelectedStartTime
			sD, err := util.ParseDateTimeInLocation(sDate, sTime, timezone)
			if err != nil {
				res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
				return
			}
			rTDate := payload.SelectedRequestTTLDate
			rTTime := payload.RequestTTLTime
			rT, err := util.ParseDateTimeInLocation(rTDate, rTTime, timezone)
			if err != nil {
				res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
				return
			}

			// 2. DryRun 으로 Access Request 요청하여 값 검증하기
			v3Builder := accessrequest.NewV3DryRunBuilder(role, sD, time.Time{}, rT, user.Teleport)
			_, err = h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
			if err != nil {
				res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
				return
			}

			// 3. 6단계 모달 빌드하기
			builder = accessrequestmodal.NewSixthStepBuilder(payload)
		case util.ARequestAccessDurationSecondOption:
			// 1. StartDate/Time, AccessDurationDate/Time, RequestTTLDate/Time 값을 time.Time 으로 만들기
			sDate := payload.SelectedStartDate
			sTime := payload.SelectedStartTime
			sD, err := util.ParseDateTimeInLocation(sDate, sTime, timezone)
			if err != nil {
				res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
				return
			}
			aDDate := payload.SelectedAccessDurationDate
			aDTime := payload.SelectedAccessDurationTime
			aD, err := util.ParseDateTimeInLocation(aDDate, aDTime, timezone)
			if err != nil {
				res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
				return
			}
			rTDate := payload.SelectedRequestTTLDate
			rTTime := payload.RequestTTLTime
			rT, err := util.ParseDateTimeInLocation(rTDate, rTTime, timezone)
			if err != nil {
				res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
				return
			}

			// 2. DryRun 으로 Access Request 요청하여 값 검증하기
			v3Builder := accessrequest.NewV3DryRunBuilder(role, sD, aD, rT, user.Teleport)
			_, err = h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
			if err != nil {
				res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
				return
			}

			// 3. 6단계 모달 빌드하기
			builder = accessrequestmodal.NewSixthStepBuilder(payload)
		}
	}

	// 3. 모달 업데이트하기
	if err := h.Services.Slack.UpdateModal(builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
