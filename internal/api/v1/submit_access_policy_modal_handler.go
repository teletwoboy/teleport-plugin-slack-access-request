package v1

import (
	"context"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/policy/models"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
)

func (i *InteractionHandler) SubmitAccessPolicyModalHandler(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 갑 준비
	payload, err := viewsubmission.ParseAccessPolicyModal(payloadStr)
	if err != nil {
		slog.Error("failed to parse access modal", "err", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. Access Policy 객체를 만들기 위한 데이터 가져오기
	//    1. Slack User
	slackUser, err := i.Services.Slack.GetUserByID(ctx, payload.RequesterID)
	if err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	//    2. User
	user, err := i.Services.User.GetUserBySlackUserID(ctx, slackUser.SlackUserID)
	if err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 3. 저장할 객체 만들기
	accessPolicy := models.NewAccessPolicy(payload, user)

	// 4. 데이터 저장하기
	createdAccessPolicy, err := i.Services.Policy.CreateAccessPolicy(ctx, accessPolicy)
	if err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 5. 메시지 전송하기
	builder := message.NewAccessPolicySubmissionBuilder(createdAccessPolicy, payload)
	channelID, timestamp, err := i.Services.Slack.PostMessage(payload.RequesterChannelID, builder)
	if err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 6. Pin 처리하기
	err = i.Services.Slack.AddPin(channelID, timestamp)
	if err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"response_action":"clear"}`))
	if err != nil {
		slog.Error("failed to write response", "err", err)
	}
}
