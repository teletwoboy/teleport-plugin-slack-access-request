package worker

import (
	"context"
	"go.opentelemetry.io/otel"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/util/container"
	"time"
)

var tracer = otel.Tracer(telemetry.WorkerService)

func StartWorker(ctx context.Context, srv *container.Services) {
	ticker := time.NewTicker(2 * time.Second)
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
	switch ob.EventType {
	case constant.AccessReviewReviewer:
		if err := performAccessReviewReviewerOutbox(ctx, ob, srv); err != nil {
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
				slog.Error(markErr.Error())
			}
		}
	case constant.AccessReviewRequester:
		if err := performAccessReviewRequesterOutbox(ctx, ob, srv); err != nil {
			if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
				slog.Error(markErr.Error())
			}
		}
	}
}
