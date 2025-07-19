package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/slack/message"
	"teleport-plugin-slack-access-request/internal/slack/modal"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/teleport/accessrequest"
	"teleport-plugin-slack-access-request/internal/teleport/models"

	slackapi "github.com/slack-go/slack"
)

type InteractionHandler struct {
	SlackSrv    slack.Service
	TeleportSrv teleport.Service
}

func NewInteractionHandler(s slack.Service, t teleport.Service) *InteractionHandler {
	return &InteractionHandler{
		SlackSrv:    s,
		TeleportSrv: t,
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

	var callbackType modal.TypePayload
	if err := json.Unmarshal([]byte(payloadStr), &callbackType); err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	switch slackapi.InteractionType(callbackType.Type) {
	case slackapi.InteractionTypeBlockActions:
		// 메시지에서 버튼 클릭 처리
		i.HandleBlockActions(payloadStr, w)
	case slackapi.InteractionTypeViewSubmission:
		// 모달에서 submit 처리
		// callback_id 로 분기 처리 가능
		i.HandleViewSubmission(payloadStr, w)
	default:
		http.Error(w, "unsupported interaction type", http.StatusBadRequest)
		return
	}
}

func (i *InteractionHandler) HandleViewSubmission(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 값 준비
	var callback modal.AccessRequestViewSubmissionPayload
	if err := json.Unmarshal([]byte(payloadStr), &callback); err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	strPrivateMetadata := callback.View.PrivateMetadata
	var privateMetadata modal.PrivateMetadataPayload
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
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
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
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
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
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
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
		Status:            summitedAccessRequest.GetState().String(),
		Expires:           summitedAccessRequest.Expiry(),
		SessionTTL:        summitedAccessRequest.GetSessionTLL(),
		AccessDuration:    summitedAccessRequest.GetMaxDuration(),
		ExpiryDate:        summitedAccessRequest.GetAccessExpiry(),
	}

	createdAccessRequest, err := i.TeleportSrv.CreateAccessRequest(ctx, accessRequest)
	if err != nil {
		slog.Error("failed to submit access request to teleport service", "error", err)
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
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
	submissionBuilder := message.NewAccessRequestSubmissionBuilder(requesterName, createdAccessRequest)
	_, _, err = i.SlackSrv.PostMessage(requesterChannelID, submissionBuilder)
	if err != nil {
		slog.Error("failed to post access request submission message to slack", "channelID", requesterChannelID, "err", err)
		http.Error(w, "failed to post access request submission message to slack", http.StatusInternalServerError)
		return
	}

	// 5. 리뷰어 채널로 Access Request 생성되었음과, 검토용 모달 열기 버튼 보내기
	// 요청한 사람, 요청된 Role, 요청된 채널, 요청 사유, 요청 만료 시각
	toReviewersBuilder := message.NewAccessRequestToReviewersBuilder(requesterName, createdAccessRequest)
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
	var callback modal.AccessRequestBlockActionsPayload
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
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
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
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
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
		slog.Error("failed to get access request status", "error", err)
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if accessRequest.Status == "APPROVED" {
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
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
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
	accessRequestReviewBuilder := modal.NewAccessRequestReviewBuilder(accessRequest, slackUser)

	// 3. 모달 보내기
	err = i.SlackSrv.OpenModal(triggerID, accessRequestReviewBuilder)
	if err != nil {
		slog.Error("failed to open modal", "err", err)
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
		_, _, err := i.SlackSrv.PostMessage(reviewersChannelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", reviewersChannelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
