package v1

import (
	"context"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
	"teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

func (i *InteractionHandler) HandleAccessPolicyChannelSelection(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 값 준비
	payload, err := blockactions.ParseAccessPolicyChannelSelect(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. 검증
	slackVerifier := verifier.NewSlack(i.Services.Slack)
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	if err := slackVerifier.VerifyUserExistsByID(ctx, payload.RequesterID, payload.RequesterName); err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	//    2. 해당 유저가 Request Channel 에 있는 사람이 맞는가?
	if err := slackVerifier.VerifyUserInChannelExistsByID(payload.RequesterID, payload.RequesterChannelID); err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	// 3. 값 구분하기
	allOpt := false
	//    1. 모든 채널 "*" 인 경우
	if payload.ChannelID == "*" {
		allOpt = true
	}

	// 4. 업데이트 모달의 기본 데이터를 띄우기 위한 데이터 모으기
	//    1. 모든 Channel 가져오기
	allChannels, err := i.Services.Slack.FetchAllChannels()
	if err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}

	var builder modal.Builder

	// 모든 채널을 고른 경우
	if allOpt {
		// 1. 모든 Teleport User 가져오기
		teleportUsers, err := i.Services.Teleport.FetchUsersWithoutSecrets(ctx)
		if err != nil {
			res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}

		// 2. 모든 유저의 Roles 가져오기
		roles, err := i.Services.Teleport.FetchRoles(ctx, teleportUsers)
		if err != nil {
			res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}

		// 3. 업데이트 빌더 만들기
		builder = modal.NewSelectRoleBuilder(allChannels, payload, roles)
	} else {
		// 특정 채널을 고른 경우
		// 1. 해당 채널에 존재하는 모든 슬랙 유저를 가져오기
		channelUsers, err := i.Services.Slack.FetchUsersInConversation(payload.ChannelID)
		if err != nil {
			res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}

		// 2. 슬랙 유저 중 DB에 있는 유저들의 Teleport 정보 가져오기
		var teleportUsers []models.User
		for _, u := range channelUsers {
			copiedUser := u
			exists, err := i.Services.Slack.ExistsUserByID(ctx, copiedUser)
			if err != nil {
				res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
				return
			}
			if exists {
				slackUser, err := i.Services.Slack.GetUserByID(ctx, copiedUser)
				if err != nil {
					res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
					return
				}
				user, err := i.Services.User.GetUserBySlackUserID(ctx, slackUser.SlackUserID)
				if err != nil {
					res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
					return
				}
				teleportUser, err := i.Services.Teleport.GetUserByTeleportUserID(ctx, user.TeleportUser.TeleportUserID)
				if err != nil {
					res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
					return
				}
				teleportUsers = append(teleportUsers, *teleportUser)
			}
		}

		// 2. 각 유저의 모든 Role을 모아오기
		roles, err := i.Services.Teleport.FetchRoles(ctx, teleportUsers)
		if err != nil {
			res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
			return
		}

		// 3. 모달 업데이트 빌더 만들기
		builder = modal.NewSelectRoleBuilder(allChannels, payload, roles)
	}

	// 4. 모달 업데이트 하기
	if err := i.Services.Slack.UpdateModal(builder, "", payload.ViewHash, payload.ViewID); err != nil {
		res.ErrorMessageToSlack(i.Services.Slack, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
