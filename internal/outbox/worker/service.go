package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	v1 "teleport-plugin-slack-access-request/internal/outbox/worker/v1"
	"teleport-plugin-slack-access-request/internal/util/container"
	"time"

	"github.com/jackc/pgx/v5"
)

func StartWorker(ctx context.Context, db *database.DB, clients *container.Clients, srv *container.Services) {
	h := v1.NewHandler(db, clients, srv)

	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		startListenWorker(ctx, h, srv)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		startBackupWorker(ctx, h, srv)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		startAlertingDeadWorker(ctx, srv)
	}()
}

func startListenWorker(ctx context.Context, h *v1.Handler, srv *container.Services) {
	dsn := database.MakeDsn()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		slog.Error("failed to connect to database for listening: %w", err)
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
		slog.Error("failed to listen: %w, channel: %s", err, constant.OutboxChannel)
		return
	}
	slog.Info("starting postgresql LISTEN on channel " + constant.OutboxChannel)

	// 동시성 제한용 세마포어 & 종료 대기용 WaitGroup
	sem := make(chan struct{}, constant.MaxConcurrent)
	var wg sync.WaitGroup

	// 종료 시 현재 처리 중인 작업 마무리 대기
	defer func() {
		wg.Wait()
	}()
	slog.Info("starting listen worker")

	for {
		// err := conn.PgConn().WaitForNotification(ctx) 차이점 분석하기
		n, err := conn.WaitForNotification(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				slog.Info("listening stopped (context done)")
				return
			}
			slog.Error("failed to wait for notification: %w", err)
			return
		}

		select {
		case sem <- struct{}{}:
		default:
			// 슬롯이 가득 차면 드롭(이후 변경 예정)
			slog.Warn("drop notification due to max concurrency")
			continue
		}

		wg.Add(1)
		go func() {
			defer func() {
				<-sem
				wg.Done()
			}()

			var payload model.OutboxNotificationPayload
			if err := json.Unmarshal([]byte(n.Payload), &payload); err != nil {
				slog.Error(err.Error())
				return
			}

			ob, err := srv.Outbox.ClaimOutboxByOutboxID(ctx, payload.OutboxID)
			if err != nil {
				slog.Error(err.Error())
			}
			if ob != nil {
				handle(ctx, ob, h, srv)
			}
		}()
	}
}

func startBackupWorker(ctx context.Context, h *v1.Handler, srv *container.Services) {
	ticker := time.NewTicker(constant.BackupInterval)
	defer ticker.Stop()

	// 동시성 제한용 세마포어 & 종료 대기용 WaitGroup
	sem := make(chan struct{}, constant.BackupMaxConcurrency)
	var wg sync.WaitGroup

	// 종료 시 현재 처리 중인 작업 마무리 대기
	defer func() {
		wg.Wait()
	}()
	slog.Info("starting backup worker")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Outbox Worker context canceled, shutting down")
			return
		case <-ticker.C:
			obs, err := srv.Outbox.ClaimOutboxes(ctx, constant.BackupPullSize)
			if err != nil {
				slog.Error(err.Error())
			}
			if obs == nil {
				continue
			}

			for _, ob := range obs {
				copiedOB := *ob

				select {
				case sem <- struct{}{}:
				default:
					// 슬롯이 가득 차면 드롭(이후 변경 예정)
					slog.Warn("drop notification due to max concurrency")
					continue
				}

				wg.Add(1)
				go func() {
					defer func() {
						<-sem
						wg.Done()
					}()
					handle(ctx, &copiedOB, h, srv)
				}()
			}
		}
	}
}

func startAlertingDeadWorker(ctx context.Context, srv *container.Services) {
	ticker := time.NewTicker(constant.AlertingDeadInterval)
	defer ticker.Stop()

	// 동시성 제한용 세마포어 & 종료 대기용 WaitGroup
	sem := make(chan struct{}, constant.AlertingDeadMaxConcurrency)
	var wg sync.WaitGroup

	// 종료 시 현재 처리 중인 작업 마무리 대기
	defer func() {
		wg.Wait()
	}()
	slog.Info("starting alerting dead worker")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Cleaning Worker context canceled, shutting down")
			return
		case <-ticker.C:
			obs, err := srv.Outbox.MarkDeadBatch(ctx, constant.AlertingDeadPullSize)
			if err != nil {
				slog.Error(err.Error())
			}
			if obs == nil {
				continue
			}

			for _, ob := range obs {
				copiedOB := *ob

				select {
				case sem <- struct{}{}:
				default:
					// 슬롯이 가득 차면 드롭(이후 변경 예정)
					slog.Warn("drop notification due to max concurrency")
					continue
				}

				wg.Add(1)
				go func() {
					defer func() {
						<-sem
						wg.Done()
					}()
					alertAttemptsExceeded(ctx, &copiedOB, srv)
				}()
			}
		}
	}
}

func handle(ctx context.Context, ob *model.Outbox, h *v1.Handler, srv *container.Services) {
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

func alertAttemptsExceeded(ctx context.Context, ob *model.Outbox, srv *container.Services) {
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
