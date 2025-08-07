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
	policymodels "teleport-plugin-slack-access-request/internal/policy/models"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/container"
	"teleport-plugin-slack-access-request/internal/util/verifier"
	"time"
)

func (h *Handler) HandleModalSubmission(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

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
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 2. Teleport 서버로 Access Request 생성을 요청하기 위한 데이터 준비
	//    1. Slack, Teleport, User
	users, err := container.NewUsers(ctx, h.Services, payload.RequesterID)
	if err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 3. 트랜잭션 시작하기
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

	// 4. 트랜잭션이 적용된 Repositories, Services 만들기
	qtx := h.DB.Queries.WithTx(tx)
	txRepos := container.NewRepositories(qtx)
	txServices := container.NewServices(h.Clients, txRepos)

	// 5. Teleport 클러스터에 Access Request 생성 요청하기
	builder := accessrequest.NewV3Builder(payload, users.Teleport)
	summitedAccessRequest, err := txServices.Teleport.SubmitAccessRequest(ctx, builder)
	if err != nil {
		res.ErrorMessageToSlack(txServices.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 6. payload, slack user, summitedAccessRequest 로 access_requests 테이블 row를 만든다.
	accessRequest := teleportmodels.NewAccessRequest(summitedAccessRequest, payload, users.User.UserID)
	createdAccessRequest, err := txServices.Teleport.CreateAccessRequest(ctx, accessRequest)
	if err != nil {
		res.ErrorMessageToSlack(txServices.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 7. Auto Review 가능하면 수행하기
	//    1. 해당 채널의 가능한 AccessPolicies 가져오기
	accessPolicies, err := h.getAutoReviewableAccessPolicies(ctx, payload, users.Teleport)
	if err != nil {
		res.ErrorMessageToSlack(txServices.Slack, payload.RequesterChannelID, err, w)
		return
	}

	//    2. 실행하기
	performed, err := h.performAutoReview(ctx, txServices, accessPolicies, accessRequest, users)
	if err != nil {
		res.ErrorMessageToSlack(txServices.Slack, payload.RequesterChannelID, err, w)
		return
	}
	if performed {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(`{"response_action":"clear"}`))
		if err != nil {
			slog.Error("failed to write response", "err", err)
		}
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
	_, _, err = txServices.Slack.PostMessage(payload.SelectedChannelID, toReviewersBuilder)
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

func (h *Handler) performAutoReview(
	ctx context.Context,
	txServices *container.Services,
	accessPolicies []*policymodels.AccessPolicy,
	ar *teleportmodels.AccessRequest,
	users *container.Users,
) (bool, error) {
	// 1. Policy를 돌면서
	var allowPolicy *policymodels.AccessPolicy
	for _, policy := range accessPolicies {
		copiedPolicy := policy
		// 2. Policy가 한번이라도 Deny가 있다면
		if copiedPolicy.Effect == util.APolicyDenyButtonValue {
			if err := performReview(ctx, txServices, copiedPolicy, ar, users); err != nil {
				return false, fmt.Errorf("failed to perform review: %w", err)
			}
			return true, nil
		}
		allowPolicy = policy
	}
	if allowPolicy != nil {
		if err := performReview(ctx, txServices, allowPolicy, ar, users); err != nil {
			return false, fmt.Errorf("failed to perform review: %w", err)
		}
		return true, nil
	}
	return false, nil
}

func (h *Handler) getAutoReviewableAccessPolicies(
	ctx context.Context,
	payload *viewsubmission.AccessRequestModal,
	teleportUser *teleportmodels.User,
) ([]*policymodels.AccessPolicy, error) {
	var accessPolicies []*policymodels.AccessPolicy

	// 1. 리뷰어 채널에 있는 모든 Access Policy 가져오기
	policies, err := h.Services.Policy.GetAccessPoliciesByInputChannelID(ctx, payload.SelectedChannelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get access policies by channel id: %w", err)
	}

	fetchedTeleportUser, err := h.Services.Teleport.FetchUserWithoutSecrets(ctx, teleportUser)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user roles: %w", err)
	}

	now := time.Now().UTC()
	// 2. 하나씩 돌아보며 검사하기
	for _, accessPolicy := range policies {
		copiedPolicy := accessPolicy
		// 1. 날짜가 해당되는지

		//    2. 시간이 벗어났는지 비교하기
		if now.Before(copiedPolicy.StartDate) || now.After(copiedPolicy.EndDate) {
			// 1. Delete 처리 + Unpin 시키기
			_, err = h.Services.Policy.DeleteAccessPolicyByAccessPolicyID(ctx, copiedPolicy.AccessPolicyID)
			if err != nil {
				return nil, fmt.Errorf("failed to delete access policy: %w", err)
			}
			err = h.Services.Slack.RemovePin(copiedPolicy.InputChannelID, copiedPolicy.MessageTimestamp)
			if err != nil {
				return nil, fmt.Errorf("failed to remove message pin: %w", err)
			}
			continue
		}

		// 2. 타겟 채널에 해당되는지
		if copiedPolicy.TargetChannelID != "*" && copiedPolicy.TargetChannelID != payload.RequesterChannelID {
			continue
		}

		// 3. 타겟 역할에 해당되는지
		isTargetRole := false
		for _, r := range fetchedTeleportUser.GetRoles() {
			if copiedPolicy.TargetRole == "*" || copiedPolicy.TargetRole == r {
				isTargetRole = true
				break
			}
		}
		if !isTargetRole {
			continue
		}

		// 4. 타겟 유저에 해당되는지
		if copiedPolicy.TargetSlackID != "*" && copiedPolicy.TargetSlackID != payload.RequesterID {
			continue
		}
		accessPolicies = append(accessPolicies, copiedPolicy)
	}
	return accessPolicies, nil
}

func performReview(
	ctx context.Context,
	txServices *container.Services,
	policy *policymodels.AccessPolicy,
	ar *teleportmodels.AccessRequest,
	users *container.Users,
) error {
	// 1. access_requests 테이블 row 업데이트 하기
	ar.UpdateState(policy.Effect)
	updatedAR, err := txServices.Teleport.UpdateAccessRequestStateByName(ctx, ar)
	if err != nil {
		return fmt.Errorf("failed to update access request: %w", err)
	}

	// 2. Review Table에 저장하기
	accessReview := teleportmodels.NewAccessReview(ar.AccessRequestID, users.User.UserID, "Auto Review", policy.Effect)
	createdAccessReview, err := txServices.Teleport.CreateAccessReview(ctx, accessReview)
	if err != nil {
		return fmt.Errorf("failed to create access review: %w", err)
	}

	// 3. Reviewer 에게 처리되었음을 알림
	builder := message.NewAutoReviewToReviewersBuilder(updatedAR, createdAccessReview, users.Slack, policy)
	_, _, err = txServices.Slack.PostMessage(updatedAR.ReviewChannelID, builder)
	if err != nil {
		return fmt.Errorf("failed to post access review submission: %w", err)
	}

	// 4. Requestor 에게 처리되었음을 알림
	builder = message.NewAutoReviewToRequesterBuilder(updatedAR, createdAccessReview, users.Slack, policy)
	_, _, err = txServices.Slack.PostMessage(updatedAR.InputChannelID, builder)
	if err != nil {
		return fmt.Errorf("failed to post access review submission: %w", err)
	}

	// 5. Teleport에 AccessRequest 업데이트 요청하기
	updateBuilder := accessrequest.NewUpdateBuilder(updatedAR.Name, updatedAR.State, "Auto Review")
	if err := txServices.Teleport.SubmitAccessRequestState(ctx, updateBuilder); err != nil {
		return fmt.Errorf("failed to submit access review state: %w", err)
	}
	return nil
}
