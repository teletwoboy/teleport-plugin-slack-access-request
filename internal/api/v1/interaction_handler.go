package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/slack/payload"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/user"

	slackapi "github.com/slack-go/slack"
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
		i.HandleOpenAccessReviewModal(payloadStr, w)
	case slackapi.InteractionTypeViewSubmission:
		switch callback.View.CallbackID {
		case "access_request_modal":
			i.HandleAccessRequestModalSubmission(payloadStr, w)
		case "access_review_modal":
			i.HandleReviewModalSubmission(payloadStr, w)
		}
	default:
		http.Error(w, "unsupported interaction type", http.StatusBadRequest)
		return
	}
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
