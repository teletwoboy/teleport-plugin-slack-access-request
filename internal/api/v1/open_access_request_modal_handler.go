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
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal/accessrequest"
	"teleport-plugin-slack-access-request/internal/slack/payload/slashcommands"
	"teleport-plugin-slack-access-request/internal/util/container"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

type OpenAccessRequestModalHandler struct {
	Services *container.Services
}

func NewOpenAccessRequestModalHandler(s *container.Services) *OpenAccessRequestModalHandler {
	return &OpenAccessRequestModalHandler{
		Services: s,
	}
}

func (o *OpenAccessRequestModalHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. 페이로드 추출
	payload, err := slashcommands.ParseAccessRole(r, w)
	if err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	// 2. 검증
	slackVerifier := verifier.NewSlack(o.Services.Slack)
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	if err := slackVerifier.VerifyUserExistsByID(ctx, payload.UserID, payload.UserName); err != nil {
		res.ErrorMessageToSlack(o.Services.Slack, payload.ChannelID, err, w)
		return
	}

	// 3. Access Role Select Modal을 만들기 위한 데이터를 수집한다.
	//    1. Slack, Teleport, User
	users, err := container.NewUsers(ctx, o.Services, payload.UserID)
	if err != nil {
		res.ErrorMessageToSlack(o.Services.Slack, payload.ChannelID, err, w)
		return
	}

	//    2. AccessInfo
	accessInfo, err := o.Services.Teleport.FetchUserAccessInfo(ctx, users.Teleport)
	if err != nil {
		res.ErrorMessageToSlack(o.Services.Slack, payload.ChannelID, err, w)
		return
	}

	// 4. 모달 builder를 만든다.
	builder := accessrequest.NewFirstStepBuilder(accessInfo, payload, users.Slack)

	// 5. 모달을 보낸다
	if err := o.Services.Slack.OpenModal(payload.TriggerID, builder); err != nil {
		res.ErrorMessageToSlack(o.Services.Slack, payload.ChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
