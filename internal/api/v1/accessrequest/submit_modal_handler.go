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
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/model/accessrequest"
	"teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/container"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

func (h *Handler) HandleModalSubmission(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.SlackTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.ARequestModalSubmission)
	defer span.End()

	// 1. 값 준비
	payload, err := viewsubmission.ParseAccessRequestModal(payloadStr)
	if err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	// 2. 검증
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	slackVerifier := verifier.NewSlack(h.Services.Slack)
	if err := slackVerifier.VerifyUserExistsByID(ctx, payload.RequesterID, payload.RequesterName); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}

	// 3. Request 수행하기
	if err := h.performTransaction(ctx, payload); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"response_action":"clear"}`))
	if err != nil {
		slog.Error("failed to write response", "err", err)
	}
}

func (h *Handler) performTransaction(ctx context.Context, payload *viewsubmission.AccessRequestModal) error {
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

	// 3. Slack, Teleport, User 정보 가져오기
	users, err := container.NewUsers(ctx, h.Services, payload.RequesterID)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	// 4. Access Request 저장하기 - 최소한의 정보만
	aRequest := teleportmodels.NewAccessRequestChg(payload, users.User.UserID)
	createdARequest, err := txServices.Teleport.CreateAccessRequest(ctx, aRequest)
	if err != nil {
		return fmt.Errorf("failed to create access request: %w", err)
	}

	// 5. UTC 변환
	err = ParseTime(payload, users.Slack)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	// 6. Teleport Access Request 생성 이벤트 발행
	ob, err := accessrequest.NewOutboxWithSubmission(payload, users.Slack, users.Teleport, users.User, createdARequest.AccessRequestID)
	if err != nil {
		return fmt.Errorf("failed to create outbox with access request creation : %w", err)
	}
	if err := txServices.Outbox.CreateOutbox(ctx, ob); err != nil {
		return fmt.Errorf("failed to create outbox: %w", err)
	}

	// 7. 트랜잭션 종료하기
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	committed = true
	return nil
}

func ParseTime(p *viewsubmission.AccessRequestModal, s *models.User) error {
	if p.SelectedStartDateOptionID == util.ARequestStartDateSecondOption {
		sDDate := p.SelectedStartDate
		sDTime := p.SelectedStartTime
		sD, err := util.ParseDateTimeInLocation(sDDate, sDTime, s.TimeZone)
		if err != nil {
			return fmt.Errorf("failed to parse start date: %w", err)
		}
		p.SelectedStartDateTime = sD
	}

	if p.SelectedAccessDurationOptionID == util.ARequestAccessDurationSecondOption {
		aDDate := p.SelectedAccessDurationDate
		aDTime := p.SelectedAccessDurationTime
		aD, err := util.ParseDateTimeInLocation(aDDate, aDTime, s.TimeZone)
		if err != nil {
			return fmt.Errorf("failed to parse acccess duration: %w", err)
		}
		p.SelectedAccessDurationDateTime = aD
	}

	if p.SelectedRequestTTLOptionID == util.ARequestRequestTTLSecondOption {
		rTDate := p.SelectedRequestTTLDate
		rTTime := p.SelectedRequestTTLTime
		rT, err := util.ParseDateTimeInLocation(rTDate, rTTime, s.TimeZone)
		if err != nil {
			return fmt.Errorf("failed to parse request ttl: %w", err)
		}
		p.SelectedRequestTTLDateTime = rT
	}
	return nil
}
