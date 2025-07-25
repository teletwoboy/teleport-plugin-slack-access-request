package v1

import (
	"context"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

func (i *InteractionHandler) HandleAccessRoleModalSelection(payloadStr string, w http.ResponseWriter) {
	ctx := context.Background()

	// 1. 값 준비
	payload, err := blockactions.ParseAccessRoleModal(payloadStr)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. 검증
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	slackVerifier := verifier.NewSlack(i.SlackSrv)
	if err := slackVerifier.VerifyUserExistsByID(ctx, payload.RequesterID, payload.RequesterName); err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelName, err, w)
		return
	}

	//    2. 해당 유저가 Request Channel 에 있는 사람이 맞는가?
	if err := slackVerifier.VerifyUserInChannelExistsByID(payload.RequesterID, payload.RequesterChannelID); err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelName, err, w)
		return
	}

	// 3. Access Request Modal을 만들기 위한 데이터를 수집한다.
	//    1. Slack User
	slackUser, err := i.SlackSrv.GetUserByID(ctx, payload.RequesterID)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelID, err, w)
		return
	}

	//    2. User
	user, err := i.UserSrv.GetUserBySlackUserID(ctx, slackUser.SlackUserID)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelID, err, w)
		return
	}

	//    3. Teleport User
	teleportUser, err := i.TeleportSrv.GetUserByTeleportUserID(ctx, user.TeleportUser.TeleportUserID)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelID, err, w)
		return
	}

	//    3. AccessInfo
	accessInfo, err := i.TeleportSrv.FetchUserAccessInfo(ctx, teleportUser)
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelID, err, w)
		return
	}

	//	  4. reviewersChannels
	channels, err := i.SlackSrv.FetchReviewersChannels()
	if err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelID, err, w)
		return
	}

	// 4. 모달 Builder를 만든다.
	builder := modal.NewAccessRequestBuilder(accessInfo, channels, payload, slackUser)

	// 5. 모달을 보낸다.
	if err := i.SlackSrv.PushModal(payload.TriggerID, builder); err != nil {
		res.ErrorMessageToSlack(i.SlackSrv, payload.RequesterChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
