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

package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/api/res"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/config"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/constant"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/model"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util/container"
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
