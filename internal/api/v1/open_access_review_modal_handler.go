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
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

func (i *InteractionHandler) HandleOpenAccessReviewModal(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 값 준비
	payload, err := blockactions.ParseOpenAccessReviewModalPayload(payloadStr)
	if err != nil {
		http.Error(w, "invalid payload format", http.StatusBadRequest)
		return
	}

	// 2. 검증
	slackVerifier := verifier.NewSlack(i.services.Slack)
	teleportVerifier := verifier.NewTeleport(i.services.Teleport)
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	if err := slackVerifier.VerifyUserExistsByID(ctx, payload.ReviewerID, payload.ReviewerName); err != nil {
		res.ErrorMessageToSlack(i.services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	//    2. 해당 유저가 ReviewersChannel 에 있는 사람이 맞는가?
	if err := slackVerifier.VerifyUserExistsInChannelByID(payload.ReviewerID, payload.ReviewerChannelID); err != nil {
		res.ErrorMessageToSlack(i.services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	//    3. access request가 존재하며, 리뷰되지 않았는가? - teleport
	if err := teleportVerifier.VerifyAccessRequestFromCluster(ctx, payload.AccessRequestName); err != nil {
		res.ErrorMessageToSlack(i.services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	//    4. access request가 존재하며, 리뷰되지 않았는가? - database
	if err := teleportVerifier.VerifyAccessRequestFromDB(ctx, payload.AccessRequestName); err != nil {
		res.ErrorMessageToSlack(i.services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	// 3. 리뷰 모달 생성하기
	//    1. Access Request 정보 가져오기
	accessRequest, err := i.services.Teleport.GetAccessRequestByName(ctx, payload.AccessRequestName)
	if err != nil {
		res.ErrorMessageToSlack(i.services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	//    2. Slack User 정보 가져오기
	slackUser, err := i.services.Slack.GetUserBySlackUserID(ctx, accessRequest.RequesterUserID)
	if err != nil {
		res.ErrorMessageToSlack(i.services.Slack, payload.ReviewerChannelID, err, w)
		return
	}

	//    3. 모달 안에서도 사용자 요청 정보 볼 수 있게 설정
	accessRequestReviewBuilder := modal.NewAccessReviewBuilder(accessRequest, slackUser, payload.ReviewerChannelID)

	// 4. 모달 열기
	if err := i.services.Slack.OpenModal(payload.TriggerID, accessRequestReviewBuilder); err != nil {
		res.ErrorMessageToSlack(i.services.Slack, payload.ReviewerChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
