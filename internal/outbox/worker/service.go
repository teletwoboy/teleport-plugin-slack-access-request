package worker

import (
	"context"
	"log/slog"
	"sync"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	v1 "teleport-plugin-slack-access-request/internal/outbox/worker/v1"
	"teleport-plugin-slack-access-request/internal/util/container"
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
