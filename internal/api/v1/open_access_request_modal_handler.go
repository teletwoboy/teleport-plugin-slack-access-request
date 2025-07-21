package v1

import (
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/slashcommands"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/user"
	"teleport-plugin-slack-access-request/internal/util/verifier"
)

type OpenAccessRequestModalHandler struct {
	SlackSrv    slack.Service
	TeleportSrv teleport.Service
	UserSrv     user.Service
}

func NewOpenAccessRequestModalHandler(s slack.Service, t teleport.Service, u user.Service) *OpenAccessRequestModalHandler {
	return &OpenAccessRequestModalHandler{
		SlackSrv:    s,
		TeleportSrv: t,
		UserSrv:     u,
	}
}

func (a *OpenAccessRequestModalHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. 페이로드 추출
	payload, err := slashcommands.ParseAccessRequest(r, w)
	if err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	// 2. 검증
	//    1. 데이터베이스에 해당 유저가 존재하는가?
	slackVerifier := verifier.NewSlack(a.SlackSrv)
	if err := slackVerifier.VerifyExistsUserByID(ctx, payload.UserID, payload.UserName); err != nil {
		res.ErrorMessageToSlack(a.SlackSrv, payload.ChannelID, err, w)
		return
	}

	// 3. Access Request Modal을 만들기 위한 데이터를 수집한다.
	//    1. slack user
	slackUser, err := a.SlackSrv.GetUserByID(ctx, payload.UserID)
	if err != nil {
		res.ErrorMessageToSlack(a.SlackSrv, payload.ChannelID, err, w)
		return
	}

	//    2. teleport user
	teleportUser, err := a.TeleportSrv.GetUserByUsername(ctx, slackUser.Email)
	if err != nil {
		res.ErrorMessageToSlack(a.SlackSrv, payload.ChannelID, err, w)
		return
	}

	//    1. AccessInfo
	accessInfo, err := a.TeleportSrv.FetchUserAccessInfo(ctx, teleportUser)
	if err != nil {
		res.ErrorMessageToSlack(a.SlackSrv, payload.ChannelID, err, w)
		return
	}

	//	  2. reviewersChannels
	channels, err := a.SlackSrv.FetchReviewersChannels()
	if err != nil {
		res.ErrorMessageToSlack(a.SlackSrv, payload.ChannelID, err, w)
		return
	}

	// 4. 모달 builder를 만든다.
	builder := modal.NewAccessRequestBuilder(accessInfo, channels, payload.ChannelID, payload.ChannelName)

	// 5. 모달을 보낸다
	if err := a.SlackSrv.OpenModal(payload.TriggerID, builder); err != nil {
		res.ErrorMessageToSlack(a.SlackSrv, payload.ChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
