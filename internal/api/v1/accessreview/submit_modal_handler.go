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
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util/container"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

func (h *Handler) HandleModalSubmission(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 값 준비
	payload, err := viewsubmission.ParseAccessReviewModal(payloadStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2. 검증
	err = h.verifyAccessReviewModal(ctx, payload)
	if err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	// 3. 트랜잭션 시작하기
	tx, err := h.DB.Conn.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("failed to begin transaction", "err", err)
		res.ErrorMessageToSlack(h.Services.Slack, payload.ReviewerChannelID, err, w)
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
	qtx := h.DB.Queries.WithTx(tx)
	txRepos := container.NewRepositories(qtx)
	txServices := container.NewServices(h.Clients, txRepos)

	// 5. Review 수행하기
	err = performReview(ctx, txServices, payload)
	if err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	// 6.. 트랜잭션 종료하기
	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit transaction", "err", err)
		return
	}
	committed = true
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
	if err := slackVerifier.VerifyUserExistsInChannelByID(payload.ReviewerID, payload.ReviewerChannelID); err != nil {
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

func performReview(ctx context.Context, txServices *container.Services, payload *viewsubmission.AccessReviewModal) error {
	filterBuilder := accessrequest.NewFilterBuilder(payload.AccessRequestName)
	accessRequests, err := txServices.Teleport.FetchAccessRequests(ctx, filterBuilder)
	if err != nil {
		return fmt.Errorf("failed to fetch access requests: %w", err)
	}
	fetchedAccessRequest := accessRequests[0]

	// 1. database 에서 업데이트 대상 row 가져오기
	accessRequest, err := txServices.Teleport.GetAccessRequestByName(ctx, payload.AccessRequestName)
	if err != nil {
		return fmt.Errorf("failed to fetch access request: %w", err)
	}

	// 2. Access Request Row 업데이트하기
	accessRequest.Update(fetchedAccessRequest, payload.Decision)
	updatedAR, err := txServices.Teleport.UpdateAccessRequestStateByName(ctx, accessRequest)
	if err != nil {
		return fmt.Errorf("failed to update access request: %w", err)
	}

	// 3. Review Table에 저장하기
	//    1. Reviewer slackUser 정보 가저오기
	slackUser, err := txServices.Slack.GetUserByID(ctx, payload.ReviewerID)
	if err != nil {
		return fmt.Errorf("failed to get user by slack userID: %w", err)
	}

	//    2. Reviewer user 정보 가져오기
	user, err := txServices.User.GetUserBySlackUserID(ctx, slackUser.SlackUserID)
	if err != nil {
		return fmt.Errorf("failed to get user by slack userID: %w", err)
	}

	//    3. Access Review 저장하기
	accessReview := models.NewAccessReview(updatedAR.AccessRequestID, user.UserID, payload.Reason, payload.Decision)
	createdAccessReview, err := txServices.Teleport.CreateAccessReview(ctx, accessReview)
	if err != nil {
		return fmt.Errorf("failed to create access review: %w", err)
	}

	// 4. Teleport에 AccessRequest 업데이트 요청하기
	updateBuilder := accessrequest.NewUpdateBuilder(payload.AccessRequestName, updatedAR.State, payload.Reason)
	err = txServices.Teleport.SubmitAccessRequestState(ctx, updateBuilder)
	if err != nil {
		return fmt.Errorf("failed to submit access review state: %w", err)
	}

	// 5. 메시지에 띄울 requester 정보 가져오기
	requesterSlackUser, err := txServices.Slack.GetUserBySlackUserID(ctx, updatedAR.RequesterUserID)
	if err != nil {
		return fmt.Errorf("failed to get user by slack userID: %w", err)
	}

	// 6. Reviewer 에게 처리되었음을 알림
	builder := message.NewAccessReviewSubmissionBuilder(updatedAR, createdAccessReview, requesterSlackUser, slackUser)
	_, _, err = txServices.Slack.PostMessage(payload.ReviewerChannelID, builder)
	if err != nil {
		return fmt.Errorf("failed to post access review submission: %w", err)
	}

	// 7. Requestor 에게 처리되었음을 알림
	builder = message.NewAccessReviewToRequestorBuilder(accessRequest, accessReview, requesterSlackUser, slackUser)
	_, _, err = txServices.Slack.PostMessage(updatedAR.InputChannelID, builder)
	if err != nil {
		return fmt.Errorf("failed to post access review submission: %w", err)
	}
	return nil
}
