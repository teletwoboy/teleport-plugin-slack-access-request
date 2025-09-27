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
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal/accesspolicy"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
	"teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

func (h *Handler) HandleChannelSelection(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), util.SlackTimeout)
	defer cancel()

	ctx, span := tracer.Start(ctx, telemetry.APolicyChannelSelection)
	defer span.End()

	// 1. 값 준비
	payload, err := blockactions.ParseChannelSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. 검증
	if err := h.verifyChannelSelection(ctx, payload); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}

	// 3. roles 섹션을 위한 데이터 모으기
	roles, err := h.getRolesForRoleSection(ctx, payload)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}

	// 4. 업데이트 빌더 만들기
	builder := accesspolicy.NewSecondStepBuilder(payload, roles)

	// 5. 모달 업데이트 하기
	if err := h.Services.Slack.UpdateModalContext(ctx, builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) verifyChannelSelection(ctx context.Context, payload *blockactions.ChannelSelect) error {
	slackVerifier := verifier.NewSlack(h.Services.Slack)
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	if err := slackVerifier.VerifyUserExistsByID(ctx, payload.RequesterID, payload.RequesterName); err != nil {
		return err
	}

	//    2. 해당 유저가 Request Channel 에 있는 사람이 맞는가?
	return slackVerifier.VerifyUserExistsInChannelByID(ctx, payload.RequesterID, payload.RequesterChannelID)
}

func (h *Handler) getRolesForRoleSection(ctx context.Context, payload *blockactions.ChannelSelect) (map[string]struct{}, error) {
	var roles map[string]struct{}
	if payload.ChannelID == util.APolicyAllOptionValue {
		roles, err := h.handleAllChannels(ctx)
		if err != nil {
			return nil, err
		}
		return roles, nil
	}
	roles, err := h.handleSpecificChannels(ctx, payload)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (h *Handler) handleAllChannels(ctx context.Context) (map[string]struct{}, error) {
	// 1. 모든 Teleport User 가져오기
	teleportUsers, err := h.Services.Teleport.FetchUsersWithoutSecrets(ctx)
	if err != nil {
		return nil, err
	}

	// 2. 모든 유저의 Roles 가져오기
	roles, err := h.Services.Teleport.FetchAllUsersRole(ctx, teleportUsers)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (h *Handler) handleSpecificChannels(ctx context.Context, payload *blockactions.ChannelSelect) (map[string]struct{}, error) {
	// 1. 해당 채널에 존재하는 모든 슬랙 유저를 가져오기
	channelUsers, err := h.Services.Slack.FetchUsersInConversationContext(ctx, payload.ChannelID)
	if err != nil {
		return nil, err
	}

	// 2. 슬랙 유저 중 DB에 있는 유저들의 Teleport 정보 가져오기
	var teleportUsers []models.User
	for _, u := range channelUsers {
		copiedUser := u
		exists, err := h.Services.Slack.ExistsUserByID(ctx, copiedUser)
		if err != nil {
			return nil, err
		}
		if exists {
			slackUser, err := h.Services.Slack.GetUserByID(ctx, copiedUser)
			if err != nil {
				return nil, err
			}
			user, err := h.Services.User.GetUserBySlackUserID(ctx, slackUser.SlackUserID)
			if err != nil {
				return nil, err
			}
			teleportUser, err := h.Services.Teleport.GetUserByTeleportUserID(ctx, user.TeleportUser.TeleportUserID)
			if err != nil {
				return nil, err
			}
			teleportUsers = append(teleportUsers, *teleportUser)
		}
	}

	// 2. 각 유저의 모든 Role을 모아오기
	roles, err := h.Services.Teleport.FetchAllUsersRole(ctx, teleportUsers)
	if err != nil {
		return nil, err
	}
	return roles, nil
}
