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

package accessreview

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
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/container"
	"teleport-plugin-slack-access-request/internal/util/verifier"

	"golang.org/x/sync/errgroup"
)

func (h *Handler) HandleModalSubmission(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.SlackTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.AReviewModalSubmission)
	defer span.End()

	// 1. 값 준비
	payload, err := viewsubmission.ParseAccessReviewModal(payloadStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2. 검증
	err = h.verifyAccessReviewModal(ctx, payload)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.ReviewerChannelID, err)
		return
	}

	// 5. Review 수행하기
	err = h.performTransaction(ctx, payload)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.ReviewerChannelID, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) verifyAccessReviewModal(ctx context.Context, payload *viewsubmission.AccessReviewModal) error {
	slackVerifier := verifier.NewSlack(h.Services.Slack)
	teleportVerifier := verifier.NewTeleport(h.Services.Teleport)
	// 1. 데이터베이스에 해당 유저가 존재하는가?
	if err := slackVerifier.VerifyUserExistsByID(ctx, payload.ReviewerID, payload.ReviewerName); err != nil {
		return fmt.Errorf("verify slack user exists by ID: %w", err)
	}

	// 2. 해당 유저가 ReviewersChannel 에 있는 사람이 맞는가?
	if err := slackVerifier.VerifyUserExistsInChannelByID(ctx, payload.ReviewerID, payload.ReviewerChannelID); err != nil {
		return fmt.Errorf("verify slack user exists in channel by ID: %w", err)
	}

	// 3. access request가 존재하며, 리뷰되지 않았는가? - teleport
	if err := teleportVerifier.VerifyAccessRequestFromCluster(ctx, payload.AccessRequestName); err != nil {
		return fmt.Errorf("verify access request from cluster: %w", err)
	}

	// 4. access request가 존재하며, 리뷰되지 않았는가? - database
	if err := teleportVerifier.VerifyAccessRequestFromDB(ctx, payload.AccessRequestName); err != nil {
		return fmt.Errorf("verify access request from db: %w", err)
	}
	return nil
}

func (h *Handler) performTransaction(ctx context.Context, payload *viewsubmission.AccessReviewModal) error {
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

	// 3. database 에서 업데이트 대상 row 가져오기
	accessRequest, err := txServices.Teleport.GetAccessRequestByName(ctx, payload.AccessRequestName)
	if err != nil {
		return fmt.Errorf("failed to fetch access request: %w", err)
	}

	// 4. Access Request Row 업데이트하기
	accessRequest.Update(payload.Decision)
	updatedAR, err := txServices.Teleport.UpdateAccessRequestStateByName(ctx, accessRequest)
	if err != nil {
		return fmt.Errorf("failed to update access request: %w", err)
	}

	// 5. Review Table에 저장하기
	//    1. Reviewer reviewerSlackUser 정보 가저오기
	reviewerSlackUser, err := txServices.Slack.GetUserByID(ctx, payload.ReviewerID)
	if err != nil {
		return fmt.Errorf("failed to get user by slack userID: %w", err)
	}

	//    2. Reviewer user 정보 가져오기
	user, err := txServices.User.GetUserBySlackUserID(ctx, reviewerSlackUser.SlackUserID)
	if err != nil {
		return fmt.Errorf("failed to get user by slack userID: %w", err)
	}

	// 6. 메시지에 띄울 requester 정보 가져오기
	requesterSlackUser, err := txServices.Slack.GetUserBySlackUserID(ctx, updatedAR.RequesterUserID)
	if err != nil {
		return fmt.Errorf("failed to get user by slack userID: %w", err)
	}

	//    3. Access Review 저장하기
	accessReview := models.NewAccessReview(updatedAR.AccessRequestID, user.UserID, payload.Reason, payload.Decision)
	createdAccessReview, err := txServices.Teleport.CreateAccessReview(ctx, accessReview)
	if err != nil {
		return fmt.Errorf("failed to create access review: %w", err)
	}

	g, gCtx := errgroup.WithContext(ctx)
	// 7. outbox 저장하기
	g.Go(func() error {
		// 1. Reviewer Channel 용 이벤트 객체 만들기
		ob, err := model.NewOutboxWithAccessReviewReviewer(updatedAR, createdAccessReview, requesterSlackUser, reviewerSlackUser, payload.MessageTs)
		if err != nil {
			return fmt.Errorf("failed to create outbox with access review in reviewer channel : %w", err)
		}

		if err := txServices.Outbox.CreateOutbox(gCtx, ob); err != nil {
			return fmt.Errorf("failed to create reviewer outbox: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		// 2. Requester Channel 용 이벤트 객체 만들기
		ob, err := model.NewOutboxWithAccessReviewRequester(updatedAR, createdAccessReview, requesterSlackUser, reviewerSlackUser)
		if err != nil {
			return fmt.Errorf("failed to create outbox with access review in requester channel : %w", err)
		}

		if err := txServices.Outbox.CreateOutbox(gCtx, ob); err != nil {
			return fmt.Errorf("failed to create requester outbox: %w", err)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to wait goroutines : %w", err)
	}

	// 8. 트랜잭션 종료하기
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction : %w", err)
	}
	committed = true
	return nil
}
