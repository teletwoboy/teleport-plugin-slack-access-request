package accessrequest

import (
	"context"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	accessrequestmodal "teleport-plugin-slack-access-request/internal/slack/builder/modal/accessrequest"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/container"
	"time"

	"github.com/gravitational/teleport/api/types"
)

func (h *Handler) HandleAccessDurationOptionSelection(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.SlackTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.ARequestAccessDurationOptionSelection)
	defer span.End()

	// 1. 값 준비
	payload, err := blockactions.ParseAccessDurationOptionSelect(payloadStr)
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

	var builder modal.Builder
	switch payload.AccessDurationOptionID {
	case util.ARequestAccessDurationFirstOption: // AccessDuration: Default
		var submittedAccessRequest types.AccessRequest
		switch payload.SelectedStartDateOptionID {
		case util.ARequestStartDateFirstOption: // StartDate: Immediately
			// 1. DryRun 으로 Access Request 요청 후 반환되는 값의 RequestTTL 값 가져오기
			v3Builder := accessrequest.NewV3DryRunBuilder(role, time.Time{}, time.Time{}, time.Time{}, user.Teleport)
			submittedAccessRequest, err = h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
			if err != nil {
				res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
				return
			}

			requestTTL, err := util.ParseTTLInLocation(submittedAccessRequest, timezone)
			if err != nil {
				res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
				return
			}

			// 2. AccessDurationTimeSelect 인스턴스 만들기
			accessDurationTimeSelect := blockactions.NewAccessDurationTimeSelectWithFirstOpt(payload)

			// 3. 5단계 모달 빌드하기
			builder = accessrequestmodal.NewFifthStepBuilder(accessDurationTimeSelect, requestTTL)
		case util.ARequestStartDateSecondOption: // StartDate: Select DateTime
			// 1. StartDate/Time 값을 time.Time 으로 만들기
			sDate := payload.SelectedStartDate
			sTime := payload.SelectedStartTime
			sD, err := util.ParseDateTimeInLocation(sDate, sTime, timezone)
			if err != nil {
				res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
				return
			}

			// 2. DryRun 으로 Access Request 요청 후 반환되는 값의 RequestTTL 값 가져오기
			v3Builder := accessrequest.NewV3DryRunBuilder(role, sD, time.Time{}, time.Time{}, user.Teleport)
			submittedAccessRequest, err = h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
			if err != nil {
				res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
				return
			}

			requestTTL, err := util.ParseTTLInLocation(submittedAccessRequest, timezone)
			if err != nil {
				res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
				return
			}

			// 3. AccessDurationTimeSelect 인스턴스 만들기
			accessDurationTimeSelect := blockactions.NewAccessDurationTimeSelectWithFirstOpt(payload)

			// 4. 5단계 모달 빌드하기
			builder = accessrequestmodal.NewFifthStepBuilder(accessDurationTimeSelect, requestTTL)
		}
	case util.ARequestAccessDurationSecondOption: // AccessDuration: Select DateTime
		// 1. Date 모달 빌드하기
		builder = accessrequestmodal.NewFourthStepDateBuilder(payload)
	}

	// 3. 모달 업데이트하기
	if err := h.Services.Slack.UpdateModalContext(ctx, builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
