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

package accesspolicy

import (
	"context"
	"encoding/json"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/metric/telemetry"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/constant"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/model"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/builder/message"

	"golang.org/x/sync/errgroup"
)

func (h *Handler) HandleCreationOutbox(ctx context.Context, ob *model.Outbox) error {
	ctx, cancel := context.WithTimeout(ctx, constant.ProcessingTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.WorkerAccessPolicyCreation)
	defer span.End()

	// 1. payload 역직렬화
	var payload model.AccessPolicyCreationPayload
	if err := json.Unmarshal([]byte(ob.Payload), &payload); err != nil {
		return err
	}
	accessPolicy := payload.AccessPolicy
	requesterRealName := payload.RequesterRealName

	// 2. 메시지 전송하기
	builder := message.NewAccessPolicyToReviewersBuilder(accessPolicy, requesterRealName)
	channelID, timestamp, err := h.Services.Slack.PostMessageContext(ctx, accessPolicy.InputChannelID, builder)
	if err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		// 3. Pin 처리하기
		return h.Services.Slack.AddPinContext(gCtx, channelID, timestamp)
	})

	g.Go(func() error {
		// 4. timestamp 업데이트하기
		accessPolicy.UpdateTimestamp(timestamp)
		return h.Services.Policy.UpdateAccessPolicyMsgTs(gCtx, accessPolicy)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return h.Services.Outbox.MarkDone(ctx, ob)
}
