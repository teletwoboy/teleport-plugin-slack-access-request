package v1

import (
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/slashcommands"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

type OpenAccessRoleModalHandler struct {
	SlackSrv    slack.Service
	TeleportSrv teleport.Service
}

func NewOpenAccessRoleModalHandler(s slack.Service, t teleport.Service) *OpenAccessRoleModalHandler {
	return &OpenAccessRoleModalHandler{
		SlackSrv:    s,
		TeleportSrv: t,
	}
}

func (o *OpenAccessRoleModalHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. 페이로드 추출
	payload, err := slashcommands.ParseAccessRole(r, w)
	if err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	// 2. 검증
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	slackVerifier := verifier.NewSlack(o.SlackSrv)
	if err := slackVerifier.VerifyUserExistsByID(ctx, payload.UserID, payload.UserName); err != nil {
		res.ErrorMessageToSlack(o.SlackSrv, payload.ChannelID, err, w)
		return
	}

	// 3. Access Role Select Modal을 만들기 위한 데이터를 수집한다.
	//    1. Slack User
	slackUser, err := o.SlackSrv.GetUserByID(ctx, payload.UserID)
	if err != nil {
		res.ErrorMessageToSlack(o.SlackSrv, payload.ChannelID, err, w)
		return
	}

	//    2. teleport user
	teleportUser, err := o.TeleportSrv.GetUserByUsername(ctx, slackUser.Email)
	if err != nil {
		res.ErrorMessageToSlack(o.SlackSrv, payload.ChannelID, err, w)
		return
	}

	//    3. AccessInfo
	accessInfo, err := o.TeleportSrv.FetchUserAccessInfo(ctx, teleportUser)
	if err != nil {
		res.ErrorMessageToSlack(o.SlackSrv, payload.ChannelID, err, w)
		return
	}

	// 4. 모달 builder를 만든다.
	builder := modal.NewAccessRoleBuilder(accessInfo, payload, slackUser)

	// 5. 모달을 보낸다
	if err := o.SlackSrv.OpenModal(payload.TriggerID, builder); err != nil {
		res.ErrorMessageToSlack(o.SlackSrv, payload.ChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
