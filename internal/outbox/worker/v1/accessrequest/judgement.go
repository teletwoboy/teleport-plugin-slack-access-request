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
	policymodels "teleport-plugin-slack-access-request/internal/policy/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/container"
	"time"

	"github.com/gravitational/teleport/api/types"
	"golang.org/x/sync/errgroup"
)

func (h *Handler) HandleJudgementOutbox(ctx context.Context, ob *model.Outbox) error {
	ctx, cancel := context.WithTimeout(ctx, constant.ProcessingTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.WorkerAccessRequestJudgement)
	defer span.End()

	// 1. payload 역직렬화
	var payload accessrequest.JudgementPayload
	if err := json.Unmarshal([]byte(ob.Payload), &payload); err != nil {
		return err
	}
	username := payload.Username

	// 2. reviewer 채널의 모든 Access Policy 가져오기
	policies, err := h.Services.Policy.GetAccessPoliciesByInputChannelID(ctx, payload.SelectedChannelID)
	if err != nil {
		return err
	}

	// 3. requester의 역할 가져오기
	teleportUser := teleportmodels.NewUserWithUsername(username)
	fetchedTeleportUser, err := h.Services.Teleport.FetchUserWithoutSecrets(ctx, teleportUser)
	if err != nil {
		return err
	}

	// 4. 하나씩 돌아보며 적용 가능한 policies 가져오기
	possiblePolicies, err := h.getAutoReviewablePolicies(ctx, payload, policies, fetchedTeleportUser)
	if err != nil {
		return err
	}
	if len(possiblePolicies) > 0 {
		return h.makeAutoReviewEvent(ctx, ob, policies, payload)
	}

	// 5. to requester 이벤트 생성하기
	newRequesterOB, err := accessrequest.NewOutboxWithToRequester(ob, payload)
	if err != nil {
		return err
	}

	// 6. to reviewer 이벤트 생성하기
	newReviewerOB, err := accessrequest.NewOutboxWithToReviewer(ob, payload)
	if err != nil {
		return err
	}

	// 7. 트랜잭션 시작하기
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

	g, gCtx := errgroup.WithContext(ctx)
	// 이벤트 저장하기
	g.Go(func() error {
		createdOB, err := txServices.Outbox.CreateOutbox(gCtx, newRequesterOB)
		if err != nil {
			return fmt.Errorf("failed to create to requester outbox: %w", err)
		}

		// Outbox Notification 생성
		obn, err := model.NewOutboxNotification(createdOB)
		if err != nil {
			return fmt.Errorf("failed to create outbox notification: %w", err)
		}
		if err := txServices.Outbox.Notify(ctx, obn); err != nil {
			return fmt.Errorf("failed to notify outbox: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		createdOB, err := txServices.Outbox.CreateOutbox(gCtx, newReviewerOB)
		if err != nil {
			return fmt.Errorf("failed to create to reviewer outbox: %w", err)
		}

		// Outbox Notification 생성
		obn, err := model.NewOutboxNotification(createdOB)
		if err != nil {
			return fmt.Errorf("failed to create outbox notification: %w", err)
		}
		if err := txServices.Outbox.Notify(ctx, obn); err != nil {
			return fmt.Errorf("failed to notify outbox: %w", err)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
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

func (h *Handler) getAutoReviewablePolicies(
	ctx context.Context,
	payload accessrequest.JudgementPayload,
	policies []*policymodels.AccessPolicy,
	teleportUser types.User,
) ([]*policymodels.AccessPolicy, error) {
	var accessPolicies []*policymodels.AccessPolicy
	now := time.Now().UTC()

	for _, policy := range policies {
		copiedPolicy := policy
		// 1. 시간이 벗어났는지 확인하기
		if now.After(copiedPolicy.EndDate) {
			// 1. Delete 처리 + Unpin 이벤트 저장하기
			if err := h.deletePolicy(ctx, copiedPolicy); err != nil {
				return nil, err
			}
			continue
		}

		// 2. 타겟 채널에 해당되는지 확인하기
		if copiedPolicy.TargetChannelID != util.APolicyAllOptionValue && copiedPolicy.TargetChannelID != payload.RequesterChannelID {
			continue
		}

		// 3. 타겟 역할에 해당되는지 확인하기
		isTargetRole := false
		for _, r := range teleportUser.GetRoles() {
			if copiedPolicy.TargetRole == util.APolicyAllOptionValue || copiedPolicy.TargetRole == r {
				isTargetRole = true
				break
			}
		}
		if !isTargetRole {
			continue
		}

		// 4. 타겟 유저에 해당되는지 확인하기
		if copiedPolicy.TargetSlackID != util.APolicyAllOptionValue && copiedPolicy.TargetSlackID != payload.RequesterID {
			continue
		}
		accessPolicies = append(accessPolicies, copiedPolicy)
	}
	return accessPolicies, nil
}

func (h *Handler) deletePolicy(ctx context.Context, policy *policymodels.AccessPolicy) error {
	// 1. 트랜잭션 시작하기
	tx, err := h.DB.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
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

	// 3. policy Delete 처리하기
	if err := txServices.Policy.DeleteAccessPolicyByAccessPolicyID(ctx, policy.AccessPolicyID); err != nil {
		return err
	}

	// 4. Unpin 이벤트 생성하기
	ob, err := model.NewOutboxWithAccessPolicyDeletion(policy)
	if err != nil {
		return err
	}
	createdOB, err := txServices.Outbox.CreateOutbox(ctx, ob)
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

	// 5. 트랜잭션 종료하기
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	committed = true
	return nil
}

func (h *Handler) getAutoReviewablePolicy(
	policies []*policymodels.AccessPolicy,
) *policymodels.AccessPolicy {
	// 1. Policy를 돌면서
	var allowPolicy *policymodels.AccessPolicy
	for _, policy := range policies {
		copiedPolicy := policy
		// 2. Policy가 한번이라도 Deny가 있다면
		if copiedPolicy.Effect == util.APolicyDenyButtonValue {
			return copiedPolicy
		}
		allowPolicy = policy
	}
	return allowPolicy
}

func (h *Handler) makeAutoReviewEvent(
	ctx context.Context,
	ob *model.Outbox,
	policies []*policymodels.AccessPolicy,
	payload accessrequest.JudgementPayload,
) error {
	possiblePolicy := h.getAutoReviewablePolicy(policies)

	// 2. Auto Review 이벤트 생성하기
	newOB, err := accessrequest.NewOutboxWithAutoReview(possiblePolicy, ob, payload)
	if err != nil {
		return err
	}

	// 2. 트랜잭션 시작하기
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

	// 3. 트랜잭션이 적용된 Repositories, Services 만들기
	qtx := h.DB.Queries.WithTx(tx)
	txRepos := container.NewRepositories(qtx)
	txServices := container.NewServices(h.Clients, txRepos)

	// 4. 이벤트 저장하기
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

	// 5. Done 처리하기
	if err := txServices.Outbox.MarkDone(ctx, ob); err != nil {
		return err
	}

	// 6. 트랜잭션 종료하기
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction : %w", err)
	}
	committed = true
	return nil
}
