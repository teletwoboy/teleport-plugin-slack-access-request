package accesspolicy

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/policy/models"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/util/container"
)

func (h *Handler) HandleModalSubmission(payloadStr string, w http.ResponseWriter) {
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
	slackUser, err := h.Services.Slack.GetUserByID(ctx, payload.RequesterID)
	if err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	//    2. User
	user, err := h.Services.User.GetUserBySlackUserID(ctx, slackUser.SlackUserID)
	if err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 3. 저장할 객체 만들기
	accessPolicy := models.NewAccessPolicy(payload, user)

	// 4. 트랜잭션 시작하기
	tx, err := h.DB.Conn.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("failed to begin transaction", "err", err)
		return
	}
	committed := false
	defer func(tx *sql.Tx) {
		if !committed {
			if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
				slog.Error("failed to rollback transaction", "err", err)
			}
		}
	}(tx)

	// 5. 트랜잭션이 적용된 Repositories, Services 만들기
	qtx := h.DB.Queries.WithTx(tx)
	txRepos := container.NewRepositories(qtx)
	txServices := container.NewServices(h.Clients, txRepos)

	// 6. 데이터 저장하기
	createdAccessPolicy, err := txServices.Policy.CreateAccessPolicy(ctx, accessPolicy)
	if err != nil {
		res.ErrorMessageToSlack(txServices.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 7. 메시지 전송하기
	builder := message.NewAccessPolicySubmissionBuilder(createdAccessPolicy, payload)
	channelID, timestamp, err := txServices.Slack.PostMessage(payload.RequesterChannelID, builder)
	if err != nil {
		res.ErrorMessageToSlack(txServices.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 8. Pin 처리하기
	err = txServices.Slack.AddPin(channelID, timestamp)
	if err != nil {
		res.ErrorMessageToSlack(txServices.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 9. timestamp 추가후 업데이트하기
	createdAccessPolicy.MessageTimestamp = timestamp
	err = h.Services.Policy.UpdateAccessPolicyMessageTimestamp(ctx, createdAccessPolicy)
	if err != nil {
		res.ErrorMessageToSlack(txServices.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 10. 트랜잭션 종료하기
	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit transaction", "err", err)
		return
	}
	committed = true

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"response_action":"clear"}`))
	if err != nil {
		slog.Error("failed to write response", "err", err)
	}
}
