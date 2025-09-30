package worker

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/util/container"
	"time"
)

var tracer = otel.Tracer(telemetry.WorkerService)

func StartWorker(ctx context.Context, srv *container.Services) {
	ticker := time.NewTicker(constant.PollInterval)
	defer ticker.Stop()

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
				handle(ctx, ob, srv)
			}
		}
	}
}

func handle(ctx context.Context, ob *model.Outbox, srv *container.Services) {
	// 기본 알림 채널로 보내는게 나을지, 요청/검토 채널에 보내는게 나을지?
	if err := validateOutboxAttempts(ob); err != nil {
		res.ErrorMessageToSlack(ctx, srv.Slack, config.Cfg.Slack.DefaultNotifChannelID, err)
		slog.Error(err.Error())
		if err := srv.Outbox.MarkDead(ctx, ob); err != nil {
			slog.Error(err.Error())
		}
	}

	switch ob.EventType {
	case constant.AccessReviewReviewer:
		if err := performAccessReviewReviewerOutbox(ctx, ob, srv); err != nil {
			slog.Error(err.Error())
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
				slog.Error(markErr.Error())
			}
		}
	case constant.AccessReviewRequester:
		if err := performAccessReviewRequesterOutbox(ctx, ob, srv); err != nil {
			slog.Error(err.Error())
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
				slog.Error(markErr.Error())
			}
		}
	case constant.AccessPolicy:
		if err := performAccessPolicy(ctx, ob, srv); err != nil {
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
