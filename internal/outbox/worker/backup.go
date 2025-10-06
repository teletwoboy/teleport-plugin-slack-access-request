package worker

import (
	"context"
	"log/slog"
	"sync"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	v1 "teleport-plugin-slack-access-request/internal/outbox/worker/v1"
	"teleport-plugin-slack-access-request/internal/util/container"
	"time"
)

func startBackupWorker(ctx context.Context, h *v1.Handler, srv *container.Services) {
	ticker := time.NewTicker(constant.BackupInterval)
	defer ticker.Stop()

	// 작업 버퍼 큐
	queue := make(chan *model.Outbox, constant.BackupQueueSize)

	// 동시성 워커 풀
	var wg sync.WaitGroup
	wg.Add(constant.BackupMaxConcurrency)

	defer func() {
		close(queue)
		wg.Wait()
		slog.Info("backup worker stopped")
	}()

	for i := 0; i < constant.BackupMaxConcurrency; i++ {
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
					handle(ctx, ob, h, srv)
				}
			}
		}()
	}
	slog.Info("starting backup worker", "pool", constant.BackupMaxConcurrency)

	for {
		select {
		case <-ctx.Done():
			slog.Info("backup worker context canceled, shutting down")
			return
		case <-ticker.C:
			taskCtx, cancel := context.WithTimeout(ctx, constant.ClaimTimeout)
			obs, err := srv.Outbox.ClaimOutboxes(taskCtx, constant.BackupPullSize)
			cancel()
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			if len(obs) == 0 {
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
