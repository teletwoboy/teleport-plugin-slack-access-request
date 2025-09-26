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

package v1

import (
	"context"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal/accesspolicy"
	"teleport-plugin-slack-access-request/internal/slack/payload/slashcommands"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/container"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

type OpenAccessPolicyModalHandler struct {
	Services *container.Services
}

func NewOpenAccessPolicyModalHandler(s *container.Services) *OpenAccessPolicyModalHandler {
	return &OpenAccessPolicyModalHandler{
		Services: s,
	}
}

func (o *OpenAccessPolicyModalHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.SlackTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.OpenModalAccessPolicy)
	defer span.End()

	// 1. 값 준비
	payload, err := slashcommands.ParseAccessPolicy(r, w)
	if err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	// 2. Slack 에서 요청자 정보 가져오기
	fetchedSlackUser, err := o.Services.Slack.FetchUserInfoContext(ctx, payload.UserID)
	if err != nil {
		res.ErrorMessageToSlack(ctx, o.Services.Slack, payload.ChannelID, err, w)
		return
	}

	// 2. 검증
	slackVerifier := verifier.NewSlack(o.Services.Slack)
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	if err := slackVerifier.VerifyUserExistsByID(ctx, payload.UserID, fetchedSlackUser.RealName); err != nil {
		res.ErrorMessageToSlack(ctx, o.Services.Slack, payload.ChannelID, err, w)
		return
	}

	//    2. 요청이 Reviewers 채널에서 온것인가?
	if err := slackVerifier.VerifyChanIsReviewersChan(payload.ChannelName); err != nil {
		res.ErrorMessageToSlack(ctx, o.Services.Slack, payload.ChannelID, err, w)
		return
	}

	//    3. 요청 유저가 해당 채널에 존재하는가?
	if err := slackVerifier.VerifyUserExistsInChannelByID(ctx, payload.UserID, payload.ChannelID); err != nil {
		res.ErrorMessageToSlack(ctx, o.Services.Slack, payload.ChannelID, err, w)
		return
	}

	// 3. Access Policy Modal을 만들기 위한 데이터를 수집한다.
	//    1. 모든 채널 목록
	allChannels, err := o.Services.Slack.FetchAllChannelsContext(ctx)
	if err != nil {
		res.ErrorMessageToSlack(ctx, o.Services.Slack, payload.ChannelID, err, w)
		return
	}

	//    2. 요청자 슬랙 정보
	slackUser, err := o.Services.Slack.GetUserByID(ctx, payload.UserID)
	if err != nil {
		res.ErrorMessageToSlack(ctx, o.Services.Slack, payload.ChannelID, err, w)
		return
	}

	// 4. 모달 builder를 만든다.
	builder := accesspolicy.NewFirstStepBuilder(allChannels, payload, slackUser)

	// 5. 모달을 보낸다.
	if err := o.Services.Slack.OpenModalContext(ctx, payload.TriggerID, builder); err != nil {
		res.ErrorMessageToSlack(ctx, o.Services.Slack, payload.ChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
