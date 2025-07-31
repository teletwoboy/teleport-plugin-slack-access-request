package v1

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util/container"
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

	// 2. 검증
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	slackVerifier := verifier.NewSlack(i.Services.Slack)
	if err := slackVerifier.VerifyUserExistsByID(ctx, payload.RequesterID, payload.RequesterName); err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 2. Teleport 서버로 Access Request 생성을 요청하기 위한 데이터 준비
	//    1. Slack, Teleport, User
	users, err := container.NewUsers(ctx, i.Services, payload.RequesterID)
	if err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 3. 트랜잭션 시작하기
	tx, err := i.DB.Conn.BeginTx(ctx, nil)
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

	// 4. 트랜잭션이 적용된 Repositories, Services 만들기
	qtx := i.DB.Queries.WithTx(tx)
	txRepos := container.NewRepositories(qtx)
	txServices := container.NewServices(i.Clients, txRepos)

	// 5. Teleport 클러스터에 Access Request 생성 요청하기
	builder := accessrequest.NewV3Builder(payload, users.Teleport)
	summitedAccessRequest, err := txServices.Teleport.SubmitAccessRequest(ctx, builder)
	if err != nil {
		res.ErrorMessageToSlack(txServices.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 7. payload, slack user, summitedAccessRequest 로 access_requests 테이블 row를 만든다.
	accessRequest := teleportmodels.NewAccessRequest(summitedAccessRequest, payload, users.Slack)
	createdAccessRequest, err := txServices.Teleport.CreateAccessRequest(ctx, accessRequest)
	if err != nil {
		res.ErrorMessageToSlack(txServices.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 8. requesterChannel 로 access request 요청 처리되었음을 메시지로 보내기
	submissionBuilder := message.NewAccessRequestSubmissionBuilder(createdAccessRequest, users.Slack)
	_, _, err = txServices.Slack.PostMessage(payload.RequesterChannelID, submissionBuilder)
	if err != nil {
		res.ErrorMessageToSlack(txServices.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 9. reviewerChannel로 access request 리뷰 요청 및 리뷰 모달 열기 버튼 보내기
	toReviewersBuilder := message.NewAccessRequestToReviewersBuilder(createdAccessRequest, users.Slack)
	_, _, err = txServices.Slack.PostMessage(payload.RequesterChannelID, toReviewersBuilder)
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
