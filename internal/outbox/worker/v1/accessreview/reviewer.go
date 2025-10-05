package accessreview

import (
	"context"
	"encoding/json"
	"golang.org/x/sync/errgroup"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
)

func (h *Handler) HandleReviewerOutbox(ctx context.Context, ob *model.Outbox) error {
	ctx, cancel := context.WithTimeout(ctx, constant.ProcessingTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.WorkerAccessReviewReviewer)
	defer span.End()

	// 1. payload 역직렬화
	var payload model.AccessReviewReviewerPayload
	if err := json.Unmarshal([]byte(ob.Payload), &payload); err != nil {
		return err
	}
	aRequest := payload.AccessRequest
	aReview := payload.AccessReview
	reqSlackUser := payload.Requester
	revSlackUser := payload.Reviewer
	messageTs := payload.MessageTs

	// 2. Teleport에 AccessRequest 업데이트 요청하기
	updateBuilder := accessrequest.NewUpdateBuilder(aRequest.Name, aRequest.State, aReview.Reason)
	err := h.Services.Teleport.SubmitAccessRequestState(ctx, updateBuilder)
	if err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		// 3. 메시지에 띄울 permalink URL 정보 가져오기
		permalink, err := h.Services.Slack.GetPermalinkContext(gCtx, aRequest.ReviewChannelID, messageTs)
		if err != nil {
			return err
		}

		// 5. Reviewer 에게 처리되었음을 알림
		builder := message.NewAccessReviewToReviewersBuilder(aRequest, aReview, reqSlackUser, revSlackUser, permalink)
		_, _, err = h.Services.Slack.PostMessageContext(gCtx, aRequest.ReviewChannelID, builder)
		if err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		// 4.  검토 요청 메시지 내용 변경하기
		builder := message.NewToReviewersUpdateBuilder(aRequest, reqSlackUser, revSlackUser)
		_, _, _, err = h.Services.Slack.UpdateMessageContext(gCtx, aRequest.ReviewChannelID, messageTs, builder)
		if err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return h.Services.Outbox.MarkDone(ctx, ob)
}
