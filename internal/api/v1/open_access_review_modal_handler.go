package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

func (i *InteractionHandler) HandleOpenAccessReviewModal(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 값 준비
	var callback blockactions.OpenAccessReviewModalPayload
	if err := json.Unmarshal([]byte(payloadStr), &callback); err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	payload, err := blockactions.ParseOpenAccessReviewModalPayload(payloadStr)
	if err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	reviewersChannelID := callback.Channel.ID
	triggerID := callback.TriggerID

	// 1. 검증
	slackVerifier := verifier.NewSlack(i.SlackSrv)
	teleportVerifier := verifier.NewTeleport(i.TeleportSrv)
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	if err := slackVerifier.VerifyExistsUserByID(ctx, payload.ReviewerID, payload.ReviewerName); err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	//    2. 해당 유저가 ReviewersChannel 에 있는 사람이 맞는가?
	if err := slackVerifier.VerifyExistsUserInChannelByID(payload.ReviewerID, payload.ReviewerChannelID); err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	//    2. 요청이 존재하는가?
	if err := teleportVerifier.VerifyExistsAccessRequestsByName(ctx, payload.AccessRequestName); err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	//    3. 요청이 이미 리뷰 되었는가? - teleport

	//    4. 요청이 이미 리뷰 되었는가? - database
	if err := teleportVerifier.VerifyReviewedAccessRequestByName(ctx, payload.AccessRequestName); err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	// 2. 리뷰 모달 생성하기
	//    1. Access Request 정보 가져오기
	accessRequest, err := i.TeleportSrv.GetAccessRequestByName(ctx, payload.AccessRequestName)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	//    2. Slack User 정보 가져오기
	slackUser, err := i.SlackSrv.GetUserBySlackUserID(ctx, accessRequest.RequesterUserID)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	//    3. 모달 안에서도 사용자 요청 정보 볼 수 있게 설정
	accessRequestReviewBuilder := modal.NewAccessReviewBuilder(accessRequest, slackUser, reviewersChannelID)

	// 3. 모달 열기
	err = i.SlackSrv.OpenModal(triggerID, accessRequestReviewBuilder)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
