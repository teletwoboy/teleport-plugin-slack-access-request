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

	switch ob.EventType {
	case constant.AccessRequestSubmission:
		if err := h.ARequest.HandleSubmissionOutbox(ctx, ob); err != nil {
			slog.Error(err.Error())
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
				slog.Error(markErr.Error())
			}
		}
	case constant.AccessRequestJudgement:
		if err := h.ARequest.HandleJudgementOutbox(ctx, ob); err != nil {
			slog.Error(err.Error())
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
				slog.Error(markErr.Error())
			}
		}
	case constant.AccessRequestAutoReview:
		if err := h.ARequest.HandleAutoReviewOutbox(ctx, ob); err != nil {
			slog.Error(err.Error())
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
				slog.Error(markErr.Error())
			}
		}
	case constant.AccessRequestAutoReviewToRequester:
		if err := h.ARequest.HandleAutoReviewToRequesterOutbox(ctx, ob); err != nil {
			slog.Error(err.Error())
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
				slog.Error(markErr.Error())
			}
		}
	case constant.AccessRequestAutoReviewToReviewer:
		if err := h.ARequest.HandleAutoReviewToReviewerOutbox(ctx, ob); err != nil {
			slog.Error(err.Error())
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
				slog.Error(markErr.Error())
			}
		}
	case constant.AccessRequestToRequester:
		if err := h.ARequest.HandleToRequesterOutbox(ctx, ob); err != nil {
			slog.Error(err.Error())
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
				slog.Error(markErr.Error())
			}
		}
	case constant.AccessRequestToReviewer:
		if err := h.ARequest.HandleToReviewerOutbox(ctx, ob); err != nil {
			slog.Error(err.Error())
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
				slog.Error(markErr.Error())
			}
		}
	case constant.AccessReviewReviewer:
		if err := h.AReview.HandleReviewerOutbox(ctx, ob); err != nil {
			slog.Error(err.Error())
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
				slog.Error(markErr.Error())
			}
		}
	case constant.AccessReviewRequester:
		if err := h.AReview.HandleRequesterOutbox(ctx, ob); err != nil {
			slog.Error(err.Error())
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
				slog.Error(markErr.Error())
			}
		}
	case constant.AccessPolicyCreation:
		if err := h.APolicy.HandleCreationOutbox(ctx, ob); err != nil {
			slog.Error(err.Error())
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
			}
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
