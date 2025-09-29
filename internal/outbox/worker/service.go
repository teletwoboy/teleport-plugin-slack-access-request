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
	ticker := time.NewTicker(constant.OutboxPollInterval)
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
	}
}

func validateOutboxAttempts(ob *model.Outbox) error {
	if ob.Attempts > constant.MaxAttempts {
		return fmt.Errorf("outbox attempts exceeds max limit %d, outbox_id : %d", constant.MaxAttempts, ob.OutboxID)
	}
	return nil
}
