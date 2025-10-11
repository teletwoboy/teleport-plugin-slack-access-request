/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package accessrequest

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/metric/telemetry"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/constant"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/model"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/model/accessrequest"
	teleport "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util/container"
)

func (h *Handler) HandleSubmissionOutbox(ctx context.Context, ob *model.Outbox) error {
	ctx, cancel := context.WithTimeout(ctx, constant.ProcessingTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.WorkerAccessRequestSubmission)
	defer span.End()

	// 1. payload 역직렬화
	var payload accessrequest.SubmissionPayload
	if err := json.Unmarshal([]byte(ob.Payload), &payload); err != nil {
		return err
	}
	p := payload.Payload
	slackUserID := payload.SlackUserID
	username := payload.Username
	userID := payload.UserID

	// 3. Teleport에 Access Request 생성 요청하기
	builder := teleport.NewV3Builder(p, username)
	submittedAccessRequest, err := h.Services.Teleport.SubmitAccessRequest(ctx, builder)
	if err != nil {
		return err
	}

	// 4. Submit된 Access Request로 DB 업데이트하기
	//    1. 해당 row 가져오기
	aRequest, err := h.Services.Teleport.GetAccessRequestByAccessRequestID(ctx, ob.AggregateID)
	if err != nil {
		return err
	}

	// 5. 트랜잭션 시작하기
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

	// 6. 트랜잭션이 적용된 Repositories, Services 만들기
	qtx := h.DB.Queries.WithTx(tx)
	txRepos := container.NewRepositories(qtx)
	txServices := container.NewServices(h.Clients, txRepos)

	//    2. 업데이트하기
	aRequest.UpdateAfterSubmission(submittedAccessRequest, p)
	if err := txServices.Teleport.UpdateAccessRequestByAccessRequestID(ctx, aRequest); err != nil {
		return err
	}

	// 7. 자동 검토 여부 이벤트 생성하기
	newOB, err := accessrequest.NewOutboxWithJudgement(ob, p, slackUserID, userID, username)
	if err != nil {
		return err
	}
	createdOB, err := txServices.Outbox.CreateOutbox(ctx, newOB)
	if err != nil {
		return err
	}

	// Outbox Notification 생성
	obn, err := model.NewOutboxNotification(createdOB)
	if err != nil {
		return fmt.Errorf("failed to create outbox notification: %w", err)
	}
	if err := txServices.Outbox.Notify(ctx, obn); err != nil {
		return fmt.Errorf("failed to notify outbox: %w", err)
	}

	// 8. Done 처리하기
	if err := txServices.Outbox.MarkDone(ctx, ob); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction : %w", err)
	}
	committed = true
	return nil
}
