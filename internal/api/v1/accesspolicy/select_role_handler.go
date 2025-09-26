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
	"teleport-plugin-slack-access-request/internal/slack/models"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
	"teleport-plugin-slack-access-request/internal/util"
)

func (h *Handler) HandleRoleSelection(payloadStr string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := tracer.Start(ctx, telemetry.APolicyRoleSelection)
	defer span.End()

	// 1. 값 준비하기
	payload, err := blockactions.ParseRoleSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. user 정보 가져오기
	users, err := h.getUsersForUserSection(ctx, payload)
	if err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 3. 모달 만들기
	builder := accesspolicy.NewThirdStepBuilder(payload, users)

	// 4. 모달 업데이트 하기
	if err := h.Services.Slack.UpdateModalContext(ctx, builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(ctx, h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getUsersForUserSection(ctx context.Context, payload *blockactions.RoleSelect) ([]models.User, error) {
	if payload.SelectedChannelID == util.APolicyAllOptionValue {
		if payload.Role == util.APolicyAllOptionValue {
			return h.handleAllChannelsAllRole(ctx)
		}
		return h.handleAllChannelsSpecificRole(ctx, payload)
	}
	if payload.Role == util.APolicyAllOptionValue {
		return h.handleSpecificChannelAllRole(ctx, payload)
	}
	return h.handleSpecificChannelSpecificRole(ctx, payload)
}

func (h *Handler) handleAllChannelsAllRole(ctx context.Context) ([]models.User, error) {
	// 모든 채널의 특정 역할을 가진 유저 가져오기
	// 1. 모든 슬랙 유저 가져오기
	var users []models.User
	slackUsers, err := h.Services.Slack.FetchUsersContext(ctx)
	if err != nil {
		return nil, err
	}

	// 2. DB에 존재하는 유저만 거르기
	var existsUsers []models.User
	for _, u := range slackUsers {
		copiedUser := u
		exists, err := h.Services.Slack.ExistsUserByID(ctx, copiedUser.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			existsUsers = append(existsUsers, copiedUser)
		}
	}

	// 3. 모든 유저 추가하기
	for _, u := range existsUsers {
		copiedUser := u
		// 1. slack user
		slackUser, err := h.Services.Slack.GetUserByID(ctx, copiedUser.ID)
		if err != nil {
			return nil, err
		}

		// 2. slack user 추가하기
		users = append(users, *slackUser)
	}
	return users, nil
}

func (h *Handler) handleAllChannelsSpecificRole(ctx context.Context, payload *blockactions.RoleSelect) ([]models.User, error) {
	// 모든 채널의 특정 역할을 가진 유저 가져오기
	// 1. 모든 슬랙 유저 가져오기
	var users []models.User
	slackUsers, err := h.Services.Slack.FetchUsersContext(ctx)
	if err != nil {
		return nil, err
	}

	// 2. DB에 존재하는 유저만 거르기
	var existsUsers []models.User
	for _, u := range slackUsers {
		copiedUser := u
		exists, err := h.Services.Slack.ExistsUserByID(ctx, copiedUser.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			existsUsers = append(existsUsers, copiedUser)
		}
	}

	// 3. 각 유저의 Teleport 정보 가져온 후 역할 정보 가져와서 일치하면 추가하기
	for _, u := range existsUsers {
		copiedUser := u
		// 1. slack user
		slackUser, err := h.Services.Slack.GetUserByID(ctx, copiedUser.ID)
		if err != nil {
			return nil, err
		}

		// 2. user
		user, err := h.Services.User.GetUserBySlackUserID(ctx, slackUser.SlackUserID)
		if err != nil {
			return nil, err
		}

		// 3. teleport user
		teleportUser, err := h.Services.Teleport.GetUserByTeleportUserID(ctx, user.TeleportUser.TeleportUserID)
		if err != nil {
			return nil, err
		}

		// 4. teleport fetch user
		fetchedTeleportUser, err := h.Services.Teleport.FetchUserWithoutSecrets(ctx, teleportUser)
		if err != nil {
			return nil, err
		}

		for _, r := range fetchedTeleportUser.GetRoles() {
			copiedRole := r
			if payload.Role == copiedRole {
				users = append(users, *slackUser)
			}
		}
	}
	return users, nil
}

func (h *Handler) handleSpecificChannelAllRole(ctx context.Context, payload *blockactions.RoleSelect) ([]models.User, error) {
	// 특정 채널의 특정 역할을 가진 유저 가져오기
	// 1. 특정 채널의 슬랙 유저 가져오기
	var users []models.User
	slackUserIDs, err := h.Services.Slack.FetchUsersInConversationContext(ctx, payload.SelectedChannelID)
	if err != nil {
		return nil, err
	}

	// 2. DB에 존재하는 유저만 거르기
	var existsUserIDs []string
	for _, u := range slackUserIDs {
		copiedUserID := u
		exists, err := h.Services.Slack.ExistsUserByID(ctx, copiedUserID)
		if err != nil {
			return nil, err
		}
		if exists {
			existsUserIDs = append(existsUserIDs, copiedUserID)
		}
	}

	// 3. 모든 유저 추가하기
	for _, u := range existsUserIDs {
		copiedUserID := u
		// 1. slack user
		slackUser, err := h.Services.Slack.GetUserByID(ctx, copiedUserID)
		if err != nil {
			return nil, err
		}

		// 2. slack user 추가하기
		users = append(users, *slackUser)
	}
	return users, nil
}

func (h *Handler) handleSpecificChannelSpecificRole(ctx context.Context, payload *blockactions.RoleSelect) ([]models.User, error) {
	// 특정 채널의 특정 역할을 가진 유저 가져오기
	// 1. 특정 채널의 슬랙 유저 가져오기
	var users []models.User
	slackUserIDs, err := h.Services.Slack.FetchUsersInConversationContext(ctx, payload.SelectedChannelID)
	if err != nil {
		return nil, err
	}

	// 2. DB에 존재하는 유저만 거르기
	var existsUserIDs []string
	for _, u := range slackUserIDs {
		copiedUserID := u
		exists, err := h.Services.Slack.ExistsUserByID(ctx, copiedUserID)
		if err != nil {
			return nil, err
		}
		if exists {
			existsUserIDs = append(existsUserIDs, copiedUserID)
		}
	}

	// 3. 각 유저의 Teleport 정보 가져온 후 역할 정보 가져와서 일치하면 추가하기
	for _, u := range existsUserIDs {
		copiedUserID := u
		// 1. slack user
		slackUser, err := h.Services.Slack.GetUserByID(ctx, copiedUserID)
		if err != nil {
			return nil, err
		}

		// 2. user
		user, err := h.Services.User.GetUserBySlackUserID(ctx, slackUser.SlackUserID)
		if err != nil {
			return nil, err
		}

		// 3. teleport user
		teleportUser, err := h.Services.Teleport.GetUserByTeleportUserID(ctx, user.TeleportUser.TeleportUserID)
		if err != nil {
			return nil, err
		}

		// 4. teleport fetch user
		fetchedTeleportUser, err := h.Services.Teleport.FetchUserWithoutSecrets(ctx, teleportUser)
		if err != nil {
			return nil, err
		}

		for _, r := range fetchedTeleportUser.GetRoles() {
			copiedRole := r
			if payload.Role == copiedRole {
				users = append(users, *slackUser)
			}
		}
	}
	return users, nil
}
