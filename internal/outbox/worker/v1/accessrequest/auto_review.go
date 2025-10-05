package accessrequest

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/outbox/model/accessrequest"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util/container"
)

func (h *Handler) HandleAutoReviewOutbox(ctx context.Context, ob *model.Outbox) error {
	ctx, cancel := context.WithTimeout(ctx, constant.ProcessingTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.WorkerAccessRequestAutoReview)
	defer span.End()

	var payload accessrequest.AutoReviewPayload
	if err := json.Unmarshal([]byte(ob.Payload), &payload); err != nil {
		return err
	}
	accessPolicyID := payload.AccessPolicyID
	accessRequestID := ob.AggregateID
	userID := payload.UserID
	slackUserID := payload.SlackUserID

	// Access Request 정보 가져오기
	aRequest, err := h.Services.Teleport.GetAccessRequestByAccessRequestID(ctx, accessRequestID)
	if err != nil {
		return err
	}

	// Access Policy 정보 가져오기
	aPolicy, err := h.Services.Policy.GetAccessPoliciesByAccessPolicyID(ctx, accessPolicyID)
	if err != nil {
		return err
	}

	// 트랜잭션 시작하기
	tx, err := h.DB.Conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction : %w", err)
	}
	committed := false
	defer func(tx *sql.Tx) {
		if !committed {
			if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
				slog.Error("failed to rollback transaction", "err", err)
			}
		}
	}(tx)

	// 트랜잭션이 적용된 Repositories, Services 만들기
	qtx := h.DB.Queries.WithTx(tx)
	txRepos := container.NewRepositories(qtx)
	txServices := container.NewServices(h.Clients, txRepos)

	// access request 테이블 row 업데이트하기
	aRequest.Update(aPolicy.Effect)
	_, err = txServices.Teleport.UpdateAccessRequestStateByName(ctx, aRequest)
	if err != nil {
		return err
	}

	// review row 저장하기
	aReview := teleportmodels.NewAccessReview(accessRequestID, userID, "Auto Review", aPolicy.Effect)
	createdAReview, err := txServices.Teleport.CreateAccessReview(ctx, aReview)
	if err != nil {
		return err
	}

	// 이벤트 저장하기
	// 1. To Requester
	reqOB, err := accessrequest.NewOutboxWithAutoReviewToRequester(aPolicy, aRequest, createdAReview, slackUserID)
	if err != nil {
		return err
	}
	createdReqOB, err := txServices.Outbox.CreateOutbox(ctx, reqOB)
	if err != nil {
		return err
	}

	// Outbox Notification 생성
	reqObn, err := model.NewOutboxNotification(createdReqOB)
	if err != nil {
		return fmt.Errorf("failed to create outbox notification: %w", err)
	}
	if err := txServices.Outbox.Notify(ctx, reqObn); err != nil {
		return fmt.Errorf("failed to notify outbox: %w", err)
	}

	// 2. To Reviewer
	revOB, err := accessrequest.NewOutboxWithAutoReviewToReviewer(aPolicy, aRequest, createdAReview, slackUserID)
	if err != nil {
		return err
	}
	createdRevOB, err := txServices.Outbox.CreateOutbox(ctx, revOB)
	if err != nil {
		return err
	}

	// Outbox Notification 생성
	revObn, err := model.NewOutboxNotification(createdRevOB)
	if err != nil {
		return fmt.Errorf("failed to create outbox notification: %w", err)
	}
	if err := txServices.Outbox.Notify(ctx, revObn); err != nil {
		return fmt.Errorf("failed to notify outbox: %w", err)
	}

	// Done 처리하기
	if err := txServices.Outbox.MarkDone(ctx, ob); err != nil {
		return err
	}

	// 트랜잭션 종료하기
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction : %w", err)
	}
	committed = true
	return nil
}
