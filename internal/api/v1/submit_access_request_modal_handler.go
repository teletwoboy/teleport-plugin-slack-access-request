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

	// 6. Access Policy 체크 및 수행하기
	//performed, err := i.performAutoReview(ctx, txServices, payload, summitedAccessRequest, users.Teleport)
	//if err != nil {
	//	res.ErrorMessageToSlack(txServices.Slack, payload.RequesterChannelID, err, w)
	//	return
	//}
	//if performed {
	//	return
	//}

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

//func (i *InteractionHandler) performAutoReview(ctx context.Context, txServices *container.Services, payload *viewsubmission.AccessRequestModal, accessRequest types.AccessRequest, teleportUser *teleportmodels.User) (bool, error) {
//	accessPolicies, err := getAutoReviewableAccessPolicies(ctx, txServices, payload, teleportUser)
//	if err != nil {
//		return false, err
//	}
//
//	if len(accessPolicies) == 0 {
//		return false, nil
//	}
//
//	for _, policy := range accessPolicies {
//		copiedPolicy := policy
//		if copiedPolicy.Effect == util.APDenyButtonValue {
//		}
//	}
//	return true, nil
//}
//
//func getAutoReviewableAccessPolicies(ctx context.Context, txServices *container.Services, payload *viewsubmission.AccessRequestModal, teleportUser *teleportmodels.User) ([]*policymodels.AccessPolicy, error) {
//	var accessPolicies []*policymodels.AccessPolicy
//
//	// 1. 리뷰어 채널에 있는 모든 Access Policy 가져오기
//	policies, err := txServices.Policy.GetAccessPoliciesByInputChannelID(ctx, payload.SelectedChannelID)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get access policies by channel id: %w", err)
//	}
//
//	fetchedTeleportUser, err := txServices.Teleport.FetchUserWithoutSecrets(ctx, teleportUser)
//	if err != nil {
//		return nil, fmt.Errorf("failed to fetch user roles: %w", err)
//	}
//
//	// 2. 하나씩 돌아보며 검사하기
//	for _, accessPolicy := range policies {
//		copiedPolicy := accessPolicy
//		// 1. 날짜가 해당되는지
//		//    1. 현재 시각을 DB의 포맷에 맞게 포맷하기
//		loc, err := time.LoadLocation(copiedPolicy.TimeZone)
//		if err != nil {
//			return nil, fmt.Errorf("failed to load timezone: %w", err)
//		}
//		now := time.Now().In(loc)
//
//		//    2. 시간이 벗어났는지 비교하기
//		if now.Before(copiedPolicy.StartDate) || now.After(copiedPolicy.EndDate) {
//			// 1. Delete 처리 + Unpin 시키기
//			_, err = txServices.Policy.DeleteAccessPolicyByAccessPolicyID(ctx, copiedPolicy.AccessPolicyID)
//			if err != nil {
//				return nil, fmt.Errorf("failed to delete access policy: %w", err)
//			}
//			err = txServices.Slack.RemovePin(copiedPolicy.InputChannelID, copiedPolicy.MessageTimestamp)
//			if err != nil {
//				return nil, fmt.Errorf("failed to remove message pin: %w", err)
//			}
//			continue
//		}
//
//		// 2. 타겟 채널에 해당되는지
//		if copiedPolicy.TargetChannelID != "*" && copiedPolicy.TargetChannelID != payload.RequesterChannelID {
//			continue
//		}
//
//		// 3. 타겟 역할에 해당되는지
//		isTargetRole := false
//		for _, r := range fetchedTeleportUser.GetRoles() {
//			if copiedPolicy.TargetRole == "*" || copiedPolicy.TargetRole == r {
//				isTargetRole = true
//				break
//			}
//		}
//		if !isTargetRole {
//			continue
//		}
//
//		// 4. 타겟 유저에 해당되는지
//		if copiedPolicy.TargetSlackID != "*" && copiedPolicy.TargetSlackID != payload.RequesterID {
//			continue
//		}
//		accessPolicies = append(accessPolicies, copiedPolicy)
//	}
//	return accessPolicies, nil
//}
