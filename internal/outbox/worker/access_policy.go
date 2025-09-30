package worker

import (
	"context"
	"encoding/json"
	"golang.org/x/sync/errgroup"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/slack/builder/message"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/container"
)

func performAccessPolicy(ctx context.Context, ob *model.Outbox, srv *container.Services) error {
	ctx, cancel := context.WithTimeout(ctx, util.Timeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.WorkerAccessPolicy)
	defer span.End()

	// 1. payload 역직렬화
	var payload model.AccessPolicyPayload
	if err := json.Unmarshal([]byte(ob.Payload), &payload); err != nil {
		return err
	}
	accessPolicy := payload.AccessPolicy
	requesterRealName := payload.RequesterRealName

	// 2. 메시지 전송하기
	builder := message.NewAccessPolicySubmissionBuilder(accessPolicy, requesterRealName)
	channelID, timestamp, err := srv.Slack.PostMessageContext(ctx, accessPolicy.InputChannelID, builder)
	if err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		// 3. Pin 처리하기
		if err := srv.Slack.AddPinContext(gCtx, channelID, timestamp); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		// 4. timestamp 업데이트하기
		accessPolicy.UpdateTimestamp(timestamp)
		if err := srv.Policy.UpdateAccessPolicyMsgTs(gCtx, accessPolicy); err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	if err := srv.Outbox.MarkDone(ctx, ob); err != nil {
		return err
	}
	return nil
}
