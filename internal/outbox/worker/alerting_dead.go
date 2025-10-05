package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/util/container"
	"time"
)

func startAlertingDeadWorker(ctx context.Context, srv *container.Services) {
	ticker := time.NewTicker(constant.AlertingDeadInterval)
	defer ticker.Stop()

	// 작업 버퍼 큐
	queue := make(chan *model.Outbox, constant.AlteringDeadQueueSize)

	// 동시성 워커 풀
	var wg sync.WaitGroup
	wg.Add(constant.AlertingDeadMaxConcurrency)

	defer func() {
		close(queue)
		wg.Wait()
		slog.Info("altering dead worker stopped")
	}()

	for i := 0; i < constant.AlertingDeadMaxConcurrency; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case ob, ok := <-queue:
					if !ok {
						return
					}
					alertAttemptsExceeded(ctx, ob, srv)
				}
			}
		}()
	}
	slog.Info("starting alerting dead worker", "pool", constant.AlertingDeadMaxConcurrency)

	for {
		select {
		case <-ctx.Done():
			slog.Info("backup worker context canceled, shutting down")
			return
		case <-ticker.C:
			taskCtx, cancel := context.WithTimeout(ctx, constant.ClaimTimeout)
			obs, err := srv.Outbox.MarkDeadBatch(taskCtx, constant.AlertingDeadPullSize)
			cancel()
			if err != nil {
				slog.Error(err.Error())
			}
			if obs == nil {
				continue
			}
			for _, ob := range obs {
				copiedOB := *ob
				select {
				case <-ctx.Done():
					return
				case queue <- &copiedOB:
				}
			}
		}
	}
}

func alertAttemptsExceeded(ctx context.Context, ob *model.Outbox, srv *container.Services) {
	ctx, cancel := context.WithTimeout(ctx, constant.DeadTimeout)
	defer cancel()

	var err error
	if ob.Attempts > constant.MaxRetries {
		if ob.LastError == "" {
			err = fmt.Errorf("execution interrupted (unexpected server shutdown), outbox_id=%d", ob.OutboxID)
		} else {
			err = fmt.Errorf("outbox attempts exceeded max retries %d, outbox_id : %d", constant.MaxRetries, ob.OutboxID)
		}
		res.ErrorMessageToSlack(ctx, srv.Slack, config.Cfg.Slack.DefaultNotifChannelID, err)
		slog.Error(err.Error())
	}
}
