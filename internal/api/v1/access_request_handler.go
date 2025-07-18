package v1

import (
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/slack/message"
	"teleport-plugin-slack-access-request/internal/slack/modal"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/user"
)

type AccessRequestHandler struct {
	SlackSrv    slack.Service
	TeleportSrv teleport.Service
	UserSrv     user.Service
}

func NewAccessRequestHandler(s slack.Service, t teleport.Service, u user.Service) *AccessRequestHandler {
	return &AccessRequestHandler{
		SlackSrv:    s,
		TeleportSrv: t,
		UserSrv:     u,
	}
}

func (a *AccessRequestHandler) HandleRequestModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. TriggerID와 SlackID를 요청 객체에서 파싱한다
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	triggerID := r.FormValue("trigger_id")
	userID := r.FormValue("user_id")
	userName := r.FormValue("user_name")
	channelID := r.FormValue("channel_id")
	channelName := r.FormValue("channel_name")
	command := r.FormValue("command")

	slog.Info("Received slash command",
		"command", command,
		"userMame", userName,
		"channelName", channelName,
	)

	// 2. SlackID로 GetUser을 통해 존재하는지 검증한다.
	exists, err := a.SlackSrv.ExistsUserByID(ctx, userID)
	if err != nil {
		slog.Error("failed to check existence of user")
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
		_, _, err := a.SlackSrv.PostMessage(channelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", channelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to post user not found message to slack", http.StatusInternalServerError)
		return
	}

	if !exists {
		slog.Error("User not found", "userID", userID, "err", err)
		userNotFoundMessageBuilder := message.NewUserNotFoundBuilder(userName)
		_, _, err := a.SlackSrv.PostMessage(channelID, userNotFoundMessageBuilder)
		if err != nil {
			slog.Error("failed to post user not found message to slack", "channelID", channelID, "err", err)
			http.Error(w, "failed to post user not found message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// 3. Access Request Modal을 만들기 위한 데이터를 수집한다.
	slackUser, err := a.SlackSrv.GetUserByID(ctx, userID)
	if err != nil {
		slog.Error("failed to get slack user by id",
			"ID", userID,
			"err", err,
		)
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
		_, _, err := a.SlackSrv.PostMessage(channelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", channelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	teleportUser, err := a.TeleportSrv.GetUserByUsername(ctx, slackUser.Email)
	if err != nil {
		slog.Error("failed to get teleport user by username",
			"username", slackUser.Email,
			"err", err,
		)
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
		_, _, err := a.SlackSrv.PostMessage(channelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", channelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	//    1. AccessInfo
	accessInfo, err := a.TeleportSrv.FetchUserAccessInfo(ctx, *teleportUser)
	if err != nil {
		slog.Error("failed to get teleport user access info",
			"username", teleportUser.Username,
			"err", err,
		)
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
		_, _, err := a.SlackSrv.PostMessage(channelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", channelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	//	  2. reviewersChannels
	channels, err := a.SlackSrv.FetchReviewersChannels()
	if err != nil {
		slog.Error("failed to get reviewers channels", "err", err)
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
		_, _, err := a.SlackSrv.PostMessage(channelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", channelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// 4.모달 builder를 만든다.
	builder := modal.NewAccessRequestBuilder(accessInfo, channels, channelID, channelName)
	// Service.OpenModal을 통해 모달을 연다
	err = a.SlackSrv.OpenModal(triggerID, builder)
	if err != nil {
		slog.Error("failed to open modal", "err", err)
		errorMessageBuilder := message.NewErrorMessageBuilder(err)
		_, _, err := a.SlackSrv.PostMessage(channelID, errorMessageBuilder)
		if err != nil {
			slog.Error("failed to post error message to slack", "channelID", channelID, "err", err)
			http.Error(w, "failed to post error message to slack", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
