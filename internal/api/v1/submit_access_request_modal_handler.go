package v1

import (
	"context"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

func (i *InteractionHandler) HandleAccessRequestModalSubmission(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 값 준비
	payload, err := viewsubmission.ParseAccessRequestModal(payloadStr)
	if err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	// 1. 검증 -
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	slackVerifier := verifier.NewSlack(i.SlackSrv)
	if err := slackVerifier.VerifyExistsUserByID(ctx, payload.RequesterID, payload.RequesterName); err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelID, err, w)
		return
	}

	//     2. 요청을 보낸 채널이 리뷰 채널과 이름이 같은가?

	// 2. Teleport 서버로 Access Request 생성을 요청한다.
	//    1. Teleport Username 정보 가져오기
	slackUser, err := i.SlackSrv.GetUserByID(ctx, payload.RequesterID)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelID, err, w)
		return
	}

	//    2. 요청하기
	payload.Email = slackUser.Email
	builder := accessrequest.NewV3Builder(payload)
	summitedAccessRequest, err := i.TeleportSrv.SubmitAccessRequest(ctx, builder)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelID, err, w)
		return
	}

	// 3. payload, slack user, summitedAccessRequest 로 access_requests 테이블 row를 만든다.
	accessRequest := models.NewAccessRequestFromSubmission(summitedAccessRequest, payload, slackUser)
	createdAccessRequest, err := i.TeleportSrv.CreateAccessRequest(ctx, accessRequest)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelID, err, w)
		return
	}

	// 4. requesterChannel 로 access request 요청 처리되었음을 메시지로 보내기
	submissionBuilder := message.NewAccessRequestSubmissionBuilder(createdAccessRequest, slackUser)
	_, _, err = i.SlackSrv.PostMessage(payload.RequesterChannelID, submissionBuilder)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelID, err, w)
		return
	}

	// 5. reviewerChannel로 access request 리뷰 요청 및 리뷰 모달 열기 버튼 보내기
	toReviewersBuilder := message.NewAccessRequestToReviewersBuilder(createdAccessRequest, slackUser)
	_, _, err = i.SlackSrv.PostMessage(payload.ReviewersChannelID, toReviewersBuilder)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
