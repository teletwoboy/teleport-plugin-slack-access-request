package accessrequest

import (
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
)

func (h *Handler) HandleStartDateOptionSelection(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := tracer.Start(ctx, telemetry.ARequestStartDateOptionSelection)
	defer span.End()

	// 1. 값 준비
	payload, err := blockactions.ParseStartDateOptionSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. 요청자 User 정보 가져오기
	user, err := container.NewUsers(ctx, h.Services, payload.RequesterID)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	timezone := user.Slack.TimeZone
	role := payload.SelectedRole

	// 3. DryRun 으로 Access Request 요청 후 반환되는 값의 sessionTTL 값 가져오기
	v3Builder := accessrequest.NewV3DryRunBuilder(role, time.Time{}, time.Time{}, time.Time{}, user.Teleport)
	submittedAccessRequest, err := h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	ttl, err := util.ParseTTLInLocation(submittedAccessRequest, timezone)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	var builder modal.Builder
	// 4. Immediately 인지 Select DateTime 인지에 따라 나누어 빌드하기
	//    1. Immediately
	if payload.StartDateOptionID == util.ARequestStartDateFirstOption {
		// 1. StartTimeSelect 인스턴스 만들기
		startTimeSelect := blockactions.NewStartTimeSelectWithFirstOpt(payload, ttl)
		// 2. 4단계 모달 빌드하기
		builder = accessrequestmodal.NewFourthStepBuilder(startTimeSelect)
	} else { // 2. Select DateTime
		// 3. Date 모달 빌드하기
		builder = accessrequestmodal.NewThirdStepDateBuilder(payload, ttl)
	}

	// 5. 모달 업데이트하기
	if err := h.Services.Slack.UpdateModalContext(ctx, builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
