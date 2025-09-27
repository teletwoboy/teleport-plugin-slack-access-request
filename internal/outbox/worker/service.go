package worker

import (
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"
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
	case model.AccessReview:
		if err := performAccessReviewOutbox(ctx, ob, srv); err != nil {
			slog.Error(err.Error())
		}
	}
}

func performAccessReviewOutbox(ctx context.Context, ob *model.Outbox, srv *container.Services) error {
	ctx, cancel := context.WithTimeout(ctx, util.Timeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.WorkerAccessReview)
	defer span.End()

	// 1. payload 역직렬화
	var payload model.AccessReviewPayload
	if err := json.Unmarshal([]byte(ob.Payload), &payload); err != nil {
		if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
			slog.Error(markErr.Error())
		}
		return err
	}
	aRequest := payload.AccessRequest
	aReview := payload.AccessReview
	reqSlackUser := payload.Requester
	revSlackUser := payload.Reviewer
	messageTs := payload.MessageTs

	// 2. Teleport에 AccessRequest 업데이트 요청하기
	updateBuilder := accessrequest.NewUpdateBuilder(aRequest.Name, aRequest.State, aReview.Reason)
	err := srv.Teleport.SubmitAccessRequestState(ctx, updateBuilder)
	if err != nil {
		if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
			slog.Error(markErr.Error())
		}
		return err
	}

	// 3. 메시지에 띄울 permalink URL 정보 가져오기
	permalink, err := srv.Slack.GetPermalinkContext(ctx, aRequest.ReviewChannelID, messageTs)
	if err != nil {
		if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
			slog.Error(markErr.Error())
		}
		return err
	}

	// 4.  검토 요청 메시지 내용 변경하기
	builder := message.NewToReviewersUpdateBuilder(aRequest, reqSlackUser, revSlackUser)
	_, _, _, err = srv.Slack.UpdateMessageContext(ctx, aRequest.ReviewChannelID, messageTs, builder)
	if err != nil {
		if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
			slog.Error(markErr.Error())
		}
		return err
	}

	// 5. Reviewer 에게 처리되었음을 알림
	builder = message.NewAccessReviewSubmissionBuilder(aRequest, aReview, reqSlackUser, revSlackUser, permalink)
	_, _, err = srv.Slack.PostMessageContext(ctx, aRequest.ReviewChannelID, builder)
	if err != nil {
		if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
			slog.Error(markErr.Error())
		}
		return err
	}

	// 6. Requestor 에게 처리되었음을 알림
	builder = message.NewAccessReviewToRequestorBuilder(aRequest, aReview, reqSlackUser, revSlackUser)
	_, _, err = srv.Slack.PostMessageContext(ctx, aRequest.InputChannelID, builder)
	if err != nil {
		if markErr := srv.Outbox.MarkFailed(ctx, ob, err); markErr != nil {
			slog.Error(markErr.Error())
		}
		return err
	}
	return nil
}
