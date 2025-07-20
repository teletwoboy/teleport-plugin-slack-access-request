package v1

import (
	"context"
	"encoding/json"
	slackapi "github.com/slack-go/slack"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/user"
)

type InteractionHandler struct {
	SlackSrv    slack.Service
	TeleportSrv teleport.Service
	UserSrv     user.Service
}

func NewInteractionHandler(s slack.Service, t teleport.Service, u user.Service) *InteractionHandler {
	return &InteractionHandler{
		SlackSrv:    s,
		TeleportSrv: t,
		UserSrv:     u,
	}
}

func (i *InteractionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	payloadStr := r.FormValue("payload")
	if payloadStr == "" {
		http.Error(w, "missing payload", http.StatusBadRequest)
		return
	}

	var callback payload.Callback
	if err := json.Unmarshal([]byte(payloadStr), &callback); err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	switch slackapi.InteractionType(callback.Type) {
	case slackapi.InteractionTypeBlockActions:
		// 메시지에서 버튼 클릭 처리
		i.HandleBlockActions(payloadStr, w)
	case slackapi.InteractionTypeViewSubmission:
		switch callback.View.CallbackID {
		case "access_request_modal":
			i.HandleViewSubmission(payloadStr, w)
		case "access_review_modal":
			i.HandleReviewModalSubmission(payloadStr, w)
		}
	default:
		http.Error(w, "unsupported interaction type", http.StatusBadRequest)
		return
	}
}

func (i *InteractionHandler) HandleViewSubmission(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 값 준비
	var callback viewsubmission.AccessRequestModalPayload
	if err := json.Unmarshal([]byte(payloadStr), &callback); err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	strPrivateMetadata := callback.View.PrivateMetadata
	var privateMetadata viewsubmission.AccessRequestModalPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	requesterID := callback.User.ID
	requesterName := callback.User.Name
	requesterChannelID := privateMetadata.ChannelID
	requesterChannelName := privateMetadata.ChannelName
	role := callback.View.State.Values.RoleBlock.RoleSelect.SelectedOption.Value
	reason := callback.View.State.Values.ReasonBlock.ReasonInput.Value
	reviewersChannelID := callback.View.State.Values.ChannelBlock.ChannelSelect.SelectedOption.Value
	reviewersChannelName := callback.View.State.Values.ChannelBlock.ChannelSelect.SelectedOption.Text.Text

	// 1. 검증 -
	//    1. 우리 DB에 있는 사용자가 맞는지
	exists, err := i.SlackSrv.ExistsUserByID(ctx, requesterID)
	if err != nil {
		slog.Error("failed to check existence of user")
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(requesterChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", requesterChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to post user not found message to slack", http.StatusInternalServerError)
		return
	}

	if !exists {
		slog.Error("User not found", "userID", requesterID, "err", err)
		userNotFoundMessageBuilder := message.NewUserNotFoundBuilder(requesterName)
		_, _, err := i.SlackSrv.PostMessage(privateMetadata.ChannelID, userNotFoundMessageBuilder)
		if err != nil {
			slog.Error("failed to post user not found message to slack", "channelID", requesterChannelID, "err", err)
			http.Error(w, "failed to post user not found message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	//     2. 같은 Role에 대한 요청이 존재한다면 5분 동안 요청 불가하기
	//     3. 요청을 보낸 채널과 리뷰어 채널이 동일한지

	// 2. Teleport 서버로 Access Request 생성 요청하기
	//    1. Teleport User 정보 가져오기
	slackUser, err := i.SlackSrv.GetUserByID(ctx, requesterID)
	if err != nil {
		slog.Error("failed to get slack user by id", "ID", requesterID, "err", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(requesterChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", requesterChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	callback.Email = slackUser.Email
	builder := accessrequest.NewV3Builder(&callback)
	summitedAccessRequest, err := i.TeleportSrv.SubmitAccessRequest(ctx, builder)
	if err != nil {
		slog.Error("failed to submit access request to teleport service", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(requesterChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", requesterChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// 3. Modal 값과 Access Request 생성 반환값으로 Access Request 테이블에 Row 만들기
	accessRequest := &models.AccessRequest{
		RequesterUserID:   slackUser.SlackUserID,
		Name:              summitedAccessRequest.GetName(),
		InputChannelID:    requesterChannelID,
		InputChannelName:  requesterChannelName,
		Role:              role,
		Reason:            reason,
		ReviewChannelID:   reviewersChannelID,
		ReviewChannelName: reviewersChannelName,
		State:             summitedAccessRequest.GetState().String(),
		Expires:           summitedAccessRequest.Expiry(),
		SessionTTL:        summitedAccessRequest.GetSessionTLL(),
		AccessDuration:    summitedAccessRequest.GetMaxDuration(),
		ExpiryDate:        summitedAccessRequest.GetAccessExpiry(),
	}

	createdAccessRequest, err := i.TeleportSrv.CreateAccessRequest(ctx, accessRequest)
	if err != nil {
		slog.Error("failed to submit access request to teleport service", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(requesterChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", requesterChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// 4. 생성이 잘되면, 요청자에게 Access Request 생성 처리되었음을 메시지로 보내기
	// 성공적으로 완료됨
	// 요청한 사람, 요청된 Role, 요청된 채널, 요청 사유
	submissionBuilder := message.NewAccessRequestSubmissionBuilder(createdAccessRequest, slackUser)
	_, _, err = i.SlackSrv.PostMessage(requesterChannelID, submissionBuilder)
	if err != nil {
		slog.Error("failed to post access request submission message to slack", "channelID", requesterChannelID, "err", err)
		http.Error(w, "failed to post access request submission message to slack", http.StatusInternalServerError)
		return
	}

	// 5. 리뷰어 채널로 Access Request 생성되었음과, 검토용 모달 열기 버튼 보내기
	// 요청한 사람, 요청된 Role, 요청된 채널, 요청 사유, 요청 만료 시각
	toReviewersBuilder := message.NewAccessRequestToReviewersBuilder(createdAccessRequest, slackUser)
	_, _, err = i.SlackSrv.PostMessage(reviewersChannelID, toReviewersBuilder)
	if err != nil {
		slog.Error("failed to post access request message for reviewers to slack", "channelID", reviewersChannelID, "err", err)
		http.Error(w, "failed to post access request message for reviewers to slack", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (i *InteractionHandler) HandleBlockActions(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 값 준비
	var callback blockactions.OpenAccessReviewModalPayload
	if err := json.Unmarshal([]byte(payloadStr), &callback); err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	reviewersChannelID := callback.Channel.ID
	reviewerID := callback.User.ID
	reviewerName := callback.User.Name
	accessRequestName := callback.Actions[0].Value
	triggerID := callback.TriggerID

	// 1. 검증
	//    1. 메시지 버튼을 누른 사람이 Reviewers 채널에 있는 사람이 맞는지 확인
	exists, err := i.SlackSrv.ExistsUserInChannelByID(reviewerID, reviewersChannelID)
	if err != nil {
		slog.Error("failed to check if user exists in slack channel", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if !exists {
		slog.Error("User not found", "userID", reviewerID, "err", err)
		userNotFoundMessageBuilder := message.NewUserNotFoundBuilder(reviewerName)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, userNotFoundMessageBuilder)
		if err != nil {
			slog.Error("failed to post user not found message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post user not found message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	//    2. 요청이 존재하는 지 확인
	exists, err = i.TeleportSrv.ExistsAccessRequestByName(ctx, accessRequestName)
	if err != nil {
		slog.Error("failed to check if access request exists in teleport service", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if !exists {
		slog.Error("AccessRequest not found", "userID", reviewerID, "err", err)
		accessRequestNotFoundBuilder := message.NewAccessRequestNotFoundBuilder(accessRequestName)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, accessRequestNotFoundBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	//    3. 요청이 이미 승인되었는지 확인
	accessRequest, err := i.TeleportSrv.GetAccessRequestByName(ctx, accessRequestName)
	if err != nil {
		slog.Error("failed to get access request state", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if accessRequest.State == "APPROVED" {
		alreadyApprovedBuilder := message.NewAccessRequestAlreadyApprovedBuilder(accessRequestName)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, alreadyApprovedBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	slackUser, err := i.SlackSrv.GetUserBySlackUserID(ctx, accessRequest.RequesterUserID)
	if err != nil {
		slog.Error("failed to get user from slack", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// 2. 리뷰 모달 생성하기
	//    다시 사용자가 요청한 내용 보여주기
	//    Allow / Deny 리스트
	//    Reason 칸
	accessRequestReviewBuilder := modal.NewAccessReviewBuilder(accessRequest, slackUser, reviewersChannelID)

	// 3. 모달 보내기
	err = i.SlackSrv.OpenModal(triggerID, accessRequestReviewBuilder)
	if err != nil {
		slog.Error("failed to open modal", "err", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (i *InteractionHandler) HandleReviewModalSubmission(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 값 준비
	var callback *viewsubmission.AccessReviewModalPayload
	err := json.Unmarshal([]byte(payloadStr), &callback)
	if err != nil {
		slog.Error("failed to unmarshal interaction callback", "error", err)
		http.Error(w, "failed to unmarshal interaction callback", http.StatusInternalServerError)
		return
	}

	strPrivateMetadata := callback.View.PrivateMetadata
	var privateMetadata viewsubmission.AccessReviewModalPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	reviewerID := callback.User.ID
	decision := callback.View.State.Values.ReviewRadio.ReviewDecision.SelectedOption.Value
	reason := callback.View.State.Values.ReasonInput.ReviewReason.Value
	reviewersChannelID := privateMetadata.ChannelID
	accessRequestName := privateMetadata.AccessRequestName
	slackUser, err := i.SlackSrv.GetUserByID(ctx, reviewerID)
	if err != nil {
		slog.Error("failed to get slack user by id", "ID", reviewerID, "err", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	reviewerName := slackUser.RealName

	// 1. 검증
	//    1. reviewers 채널에 있는 사람이 맞는지
	exists, err := i.SlackSrv.ExistsUserInChannelByID(reviewerID, reviewersChannelID)
	if err != nil {
		slog.Error("failed to check if user exists in slack channel", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if !exists {
		slog.Error("User not found", "userID", reviewerID, "err", err)
		userNotFoundMessageBuilder := message.NewUserNotFoundBuilder(reviewerName)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, userNotFoundMessageBuilder)
		if err != nil {
			slog.Error("failed to post user not found message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post user not found message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	//    2. Database에 있는 사람이 맞는지
	exists, err = i.SlackSrv.ExistsUserByID(ctx, reviewerID)
	if err != nil {
		slog.Error("failed to check existence of user")
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to post user not found message to slack", http.StatusInternalServerError)
		return
	}

	if !exists {
		slog.Error("User not found", "userID", reviewerID, "err", err)
		userNotFoundMessageBuilder := message.NewUserNotFoundBuilder(slackUser.RealName)
		_, _, err := i.SlackSrv.PostMessage(privateMetadata.ChannelID, userNotFoundMessageBuilder)
		if err != nil {
			slog.Error("failed to post user not found message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post user not found message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	//    3. 요청이 존재하는지
	exists, err = i.TeleportSrv.ExistsAccessRequestByName(ctx, accessRequestName)
	if err != nil {
		slog.Error("failed to check if access request exists in teleport service", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if !exists {
		slog.Error("AccessRequest not found", "userID", reviewerID, "err", err)
		accessRequestNotFoundBuilder := message.NewAccessRequestNotFoundBuilder(accessRequestName)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, accessRequestNotFoundBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	//    4. 요청이 이미 처리되었는지 Teleport 애서 확인하기

	//    5. 요청이 이미 처리되었는지 DB 에서 확인하기 - 로직 변경 필요
	accessRequest, err := i.TeleportSrv.GetAccessRequestByName(ctx, accessRequestName)
	if err != nil {
		slog.Error("failed to get access request state", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if accessRequest.State == "APPROVED" {
		alreadyApprovedBuilder := message.NewAccessRequestAlreadyApprovedBuilder(accessRequestName)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, alreadyApprovedBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// 2. Teleport에 AccessRequest 리뷰 요청하기
	updateBuilder := accessrequest.NewUpdateBuilder(accessRequestName, decision, reason)
	err = i.TeleportSrv.SubmitAccessRequestState(ctx, updateBuilder)
	if err != nil {
		slog.Error("failed to submit access review", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// 3. 잘되면 AccessRequest Row 업데이트 하기
	//    1. Fetch로 가져오기
	filterBuilder := accessrequest.NewFilterBuilder(accessRequestName)
	accessRequests, err := i.TeleportSrv.FetchAccessRequests(ctx, filterBuilder)
	if err != nil {
		slog.Error("failed to fetch access requests", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	fetchedAccessRequest := accessRequests[0]

	//    2. Access Request Row 업데이트하기
	accessRequest.State = fetchedAccessRequest.GetState().String()
	accessRequest.Expires = fetchedAccessRequest.Expiry()
	accessRequest.SessionTTL = fetchedAccessRequest.GetSessionTLL()
	accessRequest.AccessDuration = fetchedAccessRequest.GetMaxDuration()
	accessRequest.StartDate = *fetchedAccessRequest.GetAssumeStartTime()
	accessRequest.ExpiryDate = fetchedAccessRequest.GetAccessExpiry()

	updatedAccessRequest, err := i.TeleportSrv.UpdateAccessRequestStateByName(ctx, accessRequest)
	if err != nil {
		slog.Error("failed to update access request state", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// 4. Review Table에 저장하기
	user, err := i.UserSrv.GetUserBySlackUserID(ctx, slackUser.SlackUserID)
	if err != nil {
		slog.Error("failed to get user by slack user id", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	accessReview := &models.AccessReview{
		AccessRequestID: updatedAccessRequest.AccessRequestID,
		ReviewerUserID:  user.UserID,
		Reason:          reason,
		Decision:        decision,
	}

	createdAccessReview, err := i.TeleportSrv.CreateAccessReview(ctx, accessReview)
	if err != nil {
		slog.Error("failed to create access review", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// 5. Reviewer 에게 처리되었음을 알림
	builder := message.NewAccessReviewSubmissionBuilder(updatedAccessRequest, createdAccessReview, slackUser)
	_, _, err = i.SlackSrv.PostMessage(reviewersChannelID, builder)
	if err != nil {
		slog.Error("failed to post access review submission message to slack", "channelID", reviewersChannelID, "err", err)
		http.Error(w, "failed to post access review submission message slack", http.StatusInternalServerError)
		return
	}

	// 6. Requestor 에게 처리되었음을 알림
	requestorSlackUser, err := i.SlackSrv.GetUserBySlackUserID(ctx, updatedAccessRequest.RequesterUserID)
	if err != nil {
		slog.Error("failed to get user by slack user id", "error", err)
		errorMessageBuilder := message.NewErrorBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	builder = message.NewAccessReviewToRequestorBuilder(
		updatedAccessRequest.State,
		slackUser.RealName,
		updatedAccessRequest.ReviewChannelName,
		accessReview.Reason,
		requestorSlackUser.RealName,
		updatedAccessRequest.Role,
	)
	_, _, err = i.SlackSrv.PostMessage(updatedAccessRequest.InputChannelID, builder)
	if err != nil {
		slog.Error("failed to post access review message for requestor to slack", "channelID", updatedAccessRequest.InputChannelID, "err", err)
		http.Error(w, "failed to post access review message for requestor to slack", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
