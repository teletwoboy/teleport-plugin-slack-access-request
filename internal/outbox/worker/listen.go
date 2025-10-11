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
	"encoding/json"
	"errors"
	"log/slog"
	"sync"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/database"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/constant"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/model"
	v1 "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/worker/v1"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util/container"

	"github.com/jackc/pgx/v5"
)

func startListenWorker(ctx context.Context, h *v1.Handler, srv *container.Services) {
	dsn := database.MakeDsn()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		slog.Error("failed to connect to database for listening", "err", err.Error())
		return
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}(conn, ctx)

	query := "LISTEN " + constant.OutboxChannel
	// Listen을 위한 단일 커넥션 생성
	_, err = conn.Exec(ctx, query)
	if err != nil {
		slog.Error("failed to listen", "err", err.Error(), "channel", constant.OutboxChannel)
		return
	}
	slog.Info("starting postgresql LISTEN on channel " + constant.OutboxChannel)

	// 작업 버퍼 큐
	queue := make(chan int32, constant.ListenQueueSize)

	// 동시성 워커 풀
	var wg sync.WaitGroup
	wg.Add(constant.ListenMaxConcurrency)

	defer func() {
		close(queue)
		wg.Wait()
		slog.Info("listen worker stopped")
	}()

	for i := 0; i < constant.ListenMaxConcurrency; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case payload, ok := <-queue:
					if !ok {
						return
					}
					taskCtx, cancel := context.WithTimeout(ctx, constant.ClaimTimeout)
					ob, err := srv.Outbox.ClaimOutboxByOutboxID(taskCtx, payload)
					cancel()
					if err != nil {
						slog.Error(err.Error())
						continue
					}
					if ob != nil {
						handle(ctx, ob, h, srv)
					}
				}
			}
		}()
	}
	slog.Info("starting listen worker", "pool", constant.ListenMaxConcurrency)

	for {
		n, err := conn.WaitForNotification(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				slog.Info("listening stopped (context done)")
				break
			}
			slog.Error("failed to wait for notification", "err", err)
			continue
		}

		var payload model.OutboxNotificationPayload
		if err := json.Unmarshal([]byte(n.Payload), &payload); err != nil {
			slog.Error(err.Error())
			continue
		}

		// 큐가 꽉 차면 대기
		select {
		case <-ctx.Done():
			slog.Info("listen worker context canceled, shutting down")
			return
		case queue <- payload.OutboxID:
		}
	}
}
