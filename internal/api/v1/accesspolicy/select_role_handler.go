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
	"fmt"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal/accesspolicy"
	"teleport-plugin-slack-access-request/internal/slack/models"
	blockactions "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
)

const (
	AllChannelsAllRoles int = iota
	AllChannelsSpecificRole
	SpecificChannelAllRoles
	SpecificChannelSpecificRole
)

func (h *Handler) HandleRoleSelection(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	payload, err := blockactions.ParseRoleSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	var kind int
	switch {
	case payload.SelectedChannelID == "*" && payload.Role == "*":
		kind = AllChannelsAllRoles
	case payload.RequesterChannelID == "*" && payload.Role != "*":
		kind = AllChannelsSpecificRole
	case payload.RequesterChannelID != "*" && payload.Role == "*":
		kind = SpecificChannelAllRoles
	case payload.RequesterChannelID != "*" && payload.Role != "*":
		kind = SpecificChannelSpecificRole
	default:
		err := fmt.Errorf("invalid access policy channel, role kind")
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	var users []models.User
	switch kind {
	case AllChannelsSpecificRole:
		users, err = h.handleAllChannelsSpecificRole(ctx, payload)
		if err != nil {
			res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}
	case SpecificChannelSpecificRole:
		users, err = h.handleSpecificChannelSpecificRole(ctx, payload)
		if err != nil {
			res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}
	default:
		err := fmt.Errorf("invalid access policy channel, role kind")
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 모달 만들기
	builder := accesspolicy.NewThirdStepBuilder(payload, users)

	// 모달 업데이트 하기
	if err := h.Services.Slack.UpdateModal(builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(h.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleAllChannelsSpecificRole(ctx context.Context, payload *blockactions.RoleSelect) ([]models.User, error) {
	// 모든 채널의 특정 역할을 가진 유저 가져오기
	// 1. 모든 슬랙 유저 가져오기
	var users []models.User
	slackUsers, err := h.Services.Slack.FetchUsers()
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

func (h *Handler) handleSpecificChannelSpecificRole(ctx context.Context, payload *blockactions.RoleSelect) ([]models.User, error) {
	// 특정 채널의 특정 역할을 가진 유저 가져오기
	// 1. 특정 채널의 슬랙 유저 가져오기
	var users []models.User
	slackUserIDs, err := h.Services.Slack.FetchUsersInConversation(payload.SelectedChannelID)
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
