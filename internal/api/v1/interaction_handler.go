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

	var callback modal.AccessRequestViewSubmissionPayload
	if err := json.Unmarshal([]byte(payloadStr), &callback); err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	switch slackapi.InteractionType(callback.Type) {
	case slackapi.InteractionTypeViewSubmission:
		// 버튼 클릭 등 처리
		// callback_id 로 분기 처리 가능
		i.HandleViewSubmission(&callback, w)
	default:
		http.Error(w, "unsupported interaction type", http.StatusBadRequest)
		return
	}
}

func (i *InteractionHandler) HandleViewSubmission(callback *modal.AccessRequestViewSubmissionPayload, w http.ResponseWriter) {
	ctx := context.Background()
	// 1. 들어온 값 검증 - 이런 방식 말고, 체인 형식으로 특정 요청에 대한 검증 핸들러 등록이 가능함
	//    1. Slack 에서 온 요청이 맞는지 <- 얘만 가능함.
	//    2. 우리 DB에 있는 사용자가 맞는지 <- 슬래시 커맨드와 인터랙션의 페이로드 방식이 너무 달라서 각자 해주는게 적절
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
	builder := accessrequest.NewV3Builder(callback)
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
