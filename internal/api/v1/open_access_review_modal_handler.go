package v1

import (
	"context"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

func (i *InteractionHandler) HandleOpenAccessReviewModal(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 값 준비
	payload, err := blockactions.ParseOpenAccessReviewModalPayload(payloadStr)
	if err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	// 2. 검증
	slackVerifier := verifier.NewSlack(i.Services.Slack)
	teleportVerifier := verifier.NewTeleport(i.Services.Teleport)
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	if err := slackVerifier.VerifyUserExistsByID(ctx, payload.ReviewerID, payload.ReviewerName); err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	//    2. 해당 유저가 ReviewersChannel 에 있는 사람이 맞는가?
	if err := slackVerifier.VerifyUserExistsInChannelByID(payload.ReviewerID, payload.ReviewerChannelID); err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	//    3. access request가 존재하며, 리뷰되지 않았는가? - teleport
	if err := teleportVerifier.VerifyAccessRequestFromCluster(ctx, payload.AccessRequestName); err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	//    4. access request가 존재하며, 리뷰되지 않았는가? - database
	if err := teleportVerifier.VerifyAccessRequestFromDB(ctx, payload.AccessRequestName); err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	// 3. 리뷰 모달 생성하기
	//    1. Access Request 정보 가져오기
	accessRequest, err := i.Services.Teleport.GetAccessRequestByName(ctx, payload.AccessRequestName)
	if err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	//    2. Slack User 정보 가져오기
	slackUser, err := i.Services.Slack.GetUserBySlackUserID(ctx, accessRequest.RequesterUserID)
	if err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	//    3. 모달 안에서도 사용자 요청 정보 볼 수 있게 설정
	accessRequestReviewBuilder := modal.NewAccessReviewBuilder(accessRequest, slackUser, payload.ReviewerChannelID)

	// 4. 모달 열기
	if err := i.Services.Slack.OpenModal(payload.TriggerID, accessRequestReviewBuilder); err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.ReviewerChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
