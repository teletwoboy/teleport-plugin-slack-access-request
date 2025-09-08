package accessrequest

import (
	"context"
	"github.com/gravitational/teleport/api/types"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	modal "teleport-plugin-slack-access-request/internal/slack/builder/modal"
	accessrequestmodal "teleport-plugin-slack-access-request/internal/slack/builder/modal/accessrequest"
	"teleport-plugin-slack-access-request/internal/slack/models"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/container"
	"time"
)

const (
	maxTTLDuration = time.Hour * 30
)

func (h *Handler) HandleStartDateOptionSelection(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 값 준비
	payload, err := blockactions.ParseStartOptionSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	var builder modal.Builder
	// 2. Immediately 인지 Select DateTime 인지에 따라 나누기
	//    1. Immediately
	if payload.StartDateOptionID == util.ARequestStartDateFirstOption {
		// 1. 사용자가 요청한 Role의 max_duration 정보 가져오기
		//roleMaxDuration, err := h.Services.Teleport.FetchRoleMaxDuration(ctx, payload.SelectedRole)
		//if err != nil {
		//	res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		//	return
		//}
		//
		//// 2. 정보에 따라 분류하기
		////    1. 값이 없다면 14일로
		//var accessDuration time.Duration
		//if roleMaxDuration == time.Duration(0) {
		//	accessDuration = 14 * 24 * time.Hour
		//} else { // 2. 비어 있지 않으면 해당 값으로
		//	accessDuration = roleMaxDuration
		//}
		//// 3. 해당 정보를 바탕으로 4단계 모달 빌드하기

	} else { // 2. Select DateTime
		// 1. DryRun 으로 Access Request 요청 후 반환되는 값의 sessionTTL 및 requestTTL 값 가져오기
		// 1. 요청자 Teleport User 정보 가져오기
		user, err := container.NewUsers(ctx, h.Services, payload.RequesterID)
		if err != nil {
			res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}

		// 2. DryRun 으로 Access Request 요청 후 반환되는 값의 sessionTTL 및 requestTTL 값 가져오기
		v3Builder := accessrequest.NewV3DryRunBuilder(payload, user.Teleport)
		submittedAccessRequest, err := h.Services.Teleport.SubmitAccessRequest(ctx, v3Builder)
		if err != nil {
			res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}

		maxTTL := calculateMaxTTL(submittedAccessRequest, user.Slack)

		// 3. 해당 정보를 바탕으로 Date 모달 빌드하기
		builder = accessrequestmodal.NewThirdStepDateBuilder(payload, maxTTL)
	}
	// 3. 모달 푸시하기
	if err := h.Services.Slack.UpdateModal(builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func calculateMaxTTL(a types.AccessRequest, s *models.User) time.Time {
	timezone := s.TimeZone
	utcNow := time.Now().UTC()
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}

	localTime := utcNow.In(loc)
	duration := a.GetSessionTLL().Sub(utcNow)
	ttl := localTime.Add(duration)
	TTLDuration := ttl.Sub(localTime)

	if TTLDuration > maxTTLDuration {
		return ttl.Add(maxTTLDuration)
	}
	return ttl
}
