package v1

import (
	"context"
	"net/http"
	"teleport-plugin-slack-access-request/internal/api/res"
	"teleport-plugin-slack-access-request/internal/slack/builder/modal"
	"teleport-plugin-slack-access-request/internal/slack/payload/slashcommands"
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
	ctx := context.Background()

	// 1. 값 준비
	payload, err := slashcommands.ParseAccessPolicy(r, w)
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

	//    2. 요청이 Reviewers 채널에서 온것인가?
	if err := slackVerifier.VerifyChanIsReviewersChan(payload.ChannelName); err != nil {
		res.ErrorMessageToSlack(o.Services.Slack, payload.ChannelID, err, w)
		return
	}

	//    3. 요청 유저가 해당 채널에 존재하는가?
	if err := slackVerifier.VerifyUserInChannelExistsByID(payload.UserID, payload.ChannelID); err != nil {
		res.ErrorMessageToSlack(o.Services.Slack, payload.ChannelID, err, w)
		return
	}

	// 3. Access Policy Modal을 만들기 위한 데이터를 수집한다.
	//    1. 모든 채널 목록
	allChannels, err := o.Services.Slack.FetchAllChannels()
	if err != nil {
		res.ErrorMessageToSlack(o.Services.Slack, payload.ChannelID, err, w)
		return
	}

	// 4. 모달 builder를 만든다.
	builder := modal.NewAccessPolicyBuilder(allChannels, payload)

	// 5. 모달을 보낸다.
	if err := o.Services.Slack.OpenModal(payload.TriggerID, builder); err != nil {
		res.ErrorMessageToSlack(o.Services.Slack, payload.ChannelID, err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
