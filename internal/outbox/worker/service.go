package worker

import (
	"context"
	"fmt"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	v1 "teleport-plugin-slack-access-request/internal/outbox/worker/v1"
	"teleport-plugin-slack-access-request/internal/util/container"
	"time"
)

func StartWorker(ctx context.Context, db *database.DB, clients *container.Clients, srv *container.Services) {
	ticker := time.NewTicker(constant.PollInterval)
	defer ticker.Stop()

	h := v1.NewHandler(db, clients, srv)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Outbox Worker context canceled, shutting down")
			return
		case <-ticker.C:
			ob, err := srv.Outbox.ClaimNextOutbox(ctx)
			if err != nil {
				slog.Error(err.Error())
			}
			if ob != nil {
				handle(ctx, ob, h, srv)
			}
		}
	}
}

func handle(ctx context.Context, ob *model.Outbox, h *v1.Handler, srv *container.Services) {
	// 기본 알림 채널로 보내는게 나을지, 요청/검토 채널에 보내는게 나을지?
	if err := validateOutboxAttempts(ob); err != nil {
		res.ErrorMessageToSlack(ctx, srv.Slack, config.Cfg.Slack.DefaultNotifChannelID, err)
		slog.Error(err.Error())
		if err := srv.Outbox.MarkDead(ctx, ob); err != nil {
			slog.Error(err.Error())
		}
	}

	switch ob.AggregateType {
	case constant.AccessRequest:
		handleAccessRequest(ctx, ob, h, srv)
	case constant.AccessReview:
		handleAccessReview(ctx, ob, h, srv)
	case constant.AccessPolicy:
		handleAccessPolicy(ctx, ob, h, srv)
	}
}

func handleAccessRequest(ctx context.Context, ob *model.Outbox, h *v1.Handler, srv *container.Services) {
	var err error
	switch ob.EventType {
	case constant.AccessRequestSubmission:
		err = h.ARequest.HandleSubmissionOutbox(ctx, ob)
	case constant.AccessRequestJudgement:
		err = h.ARequest.HandleJudgementOutbox(ctx, ob)
	case constant.AccessRequestAutoReview:
		err = h.ARequest.HandleAutoReviewOutbox(ctx, ob)
	case constant.AccessRequestAutoReviewToRequester:
		err = h.ARequest.HandleAutoReviewToRequesterOutbox(ctx, ob)
	case constant.AccessRequestAutoReviewToReviewer:
		err = h.ARequest.HandleAutoReviewToReviewerOutbox(ctx, ob)
	case constant.AccessRequestToRequester:
		err = h.ARequest.HandleToRequesterOutbox(ctx, ob)
	case constant.AccessRequestToReviewer:
		err = h.ARequest.HandleToReviewerOutbox(ctx, ob)
	}
	if err != nil {
		slog.Error(err.Error())
		if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
			slog.Error(markErr.Error())
		}
	}
}

func handleAccessReview(ctx context.Context, ob *model.Outbox, h *v1.Handler, srv *container.Services) {
	var err error
	switch ob.EventType {
	case constant.AccessReviewReviewer:
		err = h.AReview.HandleReviewerOutbox(ctx, ob)
	case constant.AccessReviewRequester:
		err = h.AReview.HandleRequesterOutbox(ctx, ob)
	}
	if err != nil {
		slog.Error(err.Error())
		if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
			slog.Error(markErr.Error())
		}
	}
}

func handleAccessPolicy(ctx context.Context, ob *model.Outbox, h *v1.Handler, srv *container.Services) {
	var err error
	switch ob.EventType {
	case constant.AccessPolicyCreation:
		err = h.APolicy.HandleCreationOutbox(ctx, ob)
	case constant.AccessPolicyDeletion:
		err = h.APolicy.HandleDeletionOutbox(ctx, ob)
	}
	if err != nil {
		slog.Error(err.Error())
		if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
			slog.Error(markErr.Error())
		}
	}
}

func validateOutboxAttempts(ob *model.Outbox) error {
	if ob.Attempts > constant.MaxRetries {
		if ob.LastError == "" {
			return fmt.Errorf("execution interrupted (unexpected server shutdown), outbox_id=%d", ob.OutboxID)
		}
		return fmt.Errorf("outbox attempts exceeded max retries %d, outbox_id : %d", constant.MaxRetries, ob.OutboxID)
	}
	return nil
}
