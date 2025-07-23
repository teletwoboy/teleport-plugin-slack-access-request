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

func (i *InteractionHandler) HandleAccessReviewModalSubmission(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 값 준비
	payload, err := viewsubmission.ParseAccessReviewModal(payloadStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2. 검증
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

	//    3. access request가 존재하며, 리뷰되지 않았는가? - teleport
	if err := teleportVerifier.VerifyAccessRequestFromCluster(ctx, payload.AccessRequestName); err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	//    4. access request가 존재하며, 리뷰되지 않았는가? - database
	if err := teleportVerifier.VerifyAccessRequestFromDB(ctx, payload.AccessRequestName); err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	// 3. Teleport에 AccessRequest 업데이트 요청하기
	updateBuilder := accessrequest.NewUpdateBuilder(payload.AccessRequestName, payload.Decision, payload.Reason)
	err = i.TeleportSrv.SubmitAccessRequestState(ctx, updateBuilder)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	// 4. access_requests 테이블 row 업데이트 하기
	//    1. Fetch로 teleport 에서 업데이트된 Access Request 정보 가져오기
	filterBuilder := accessrequest.NewFilterBuilder(payload.AccessRequestName)
	accessRequests, err := i.TeleportSrv.FetchAccessRequests(ctx, filterBuilder)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
	}
	fetchedAccessRequest := accessRequests[0]

	//    2. database 에서 업데이트 대상 row 가져오기
	accessRequest, err := i.TeleportSrv.GetAccessRequestByName(ctx, payload.AccessRequestName)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	//    2. Access Request Row 업데이트하기
	accessRequest.Update(fetchedAccessRequest)
	updatedAccessRequest, err := i.TeleportSrv.UpdateAccessRequestStateByName(ctx, accessRequest)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	// 5. Review Table에 저장하기
	//    1. slackUser 정보 가저오기
	slackUser, err := i.SlackSrv.GetUserBySlackUserID(ctx, accessRequest.RequesterUserID)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	//    2. user 정보 가져오기
	user, err := i.UserSrv.GetUserBySlackUserID(ctx, slackUser.SlackUserID)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	accessReview := models.NewAccessReviewFromSubmission(accessRequest.AccessRequestID, user.UserID, payload)
	createdAccessReview, err := i.TeleportSrv.CreateAccessReview(ctx, accessReview)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
	}

	// 6. 메시지에 띄울 requester 정보 가져오기
	requesterSlackUser, err := i.SlackSrv.GetUserBySlackUserID(ctx, updatedAccessRequest.RequesterUserID)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	// 7. Reviewer 에게 처리되었음을 알림
	builder := message.NewAccessReviewSubmissionBuilder(updatedAccessRequest, createdAccessReview, requesterSlackUser, slackUser)
	_, _, err = i.SlackSrv.PostMessage(payload.ReviewerChannelID, builder)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}

	// 8. Requestor 에게 처리되었음을 알림
	builder = message.NewAccessReviewToRequestorBuilder(accessRequest, accessReview, requesterSlackUser, slackUser)
	_, _, err = i.SlackSrv.PostMessage(updatedAccessRequest.InputChannelID, builder)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.ReviewerChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
