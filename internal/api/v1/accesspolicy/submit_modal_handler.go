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

package accesspolicy

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/policy/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/container"
)

func (h *Handler) HandleModalSubmission(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.SlackTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.APolicyModalSubmission)
	defer span.End()

	// 1. 갑 준비
	payload, err := viewsubmission.ParseAccessPolicyModal(payloadStr)
	if err != nil {
		slog.Error("failed to parse access modal", "err", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := h.performTransaction(ctx, payload); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"response_action":"clear"}`))
	if err != nil {
		slog.Error("failed to write response", "err", err)
	}
}

func (h *Handler) performTransaction(ctx context.Context, payload *viewsubmission.AccessPolicyModal) error {
	// 1. 트랜잭션 시작하기
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

	// 2. 트랜잭션이 적용된 Repositories, Services 만들기
	qtx := h.DB.Queries.WithTx(tx)
	txRepos := container.NewRepositories(qtx)
	txServices := container.NewServices(h.Clients, txRepos)

	// 3. Access Policy 객체를 만들기 위한 데이터 가져오기
	//    1. Slack User
	slackUser, err := h.Services.Slack.GetUserByID(ctx, payload.RequesterID)
	if err != nil {
		return fmt.Errorf("failed to get user by id : %w", err)
	}

	//    2. User
	user, err := h.Services.User.GetUserBySlackUserID(ctx, slackUser.SlackUserID)
	if err != nil {
		return fmt.Errorf("failed to get user by slack id : %w", err)
	}

	// 4. Access Policy 객체 만들기
	accessPolicy := models.NewAccessPolicy(payload, user)

	// 5. 데이터 저장하기
	createdAccessPolicy, err := txServices.Policy.CreateAccessPolicy(ctx, accessPolicy)
	if err != nil {
		return fmt.Errorf("failed to create access policy : %w", err)
	}

	// 6. Policy 이벤트 객체 만들기
	ob, err := model.NewOutboxWithAccessPolicyCreation(createdAccessPolicy, payload.RequesterRealName)

	if err := txServices.Outbox.CreateOutbox(ctx, ob); err != nil {
		return fmt.Errorf("failed to create requester outbox: %w", err)
	}

	// 7. 트랜잭션 종료하기
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction : %w", err)
	}
	committed = true
	return nil
}
