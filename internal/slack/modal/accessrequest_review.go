package modal

import (
	"fmt"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"time"

	"github.com/slack-go/slack"
)

type accessRequestReviewBuilder struct {
	accessRequest *teleportmodels.AccessRequest
	slackUser     *slackmodels.User
}

func NewAccessRequestReviewBuilder(a *teleportmodels.AccessRequest, s *slackmodels.User) Builder {
	return &accessRequestReviewBuilder{
		accessRequest: a,
		slackUser:     s,
	}
}

func (a *accessRequestReviewBuilder) Build() (*slack.ModalViewRequest, error) {
	section := a.BuildSectionBlock()
	radioBlock := a.BuildRadioBlock()
	reasonBlock := a.BuildReasonBlock()

	modal := slack.ModalViewRequest{
		Type:       slack.VTModal,
		Title:      slack.NewTextBlockObject("plain_text", "Review Access Request", false, false),
		Close:      slack.NewTextBlockObject("plain_text", "Close", false, false),
		Submit:     slack.NewTextBlockObject("plain_text", "Submit", false, false),
		CallbackID: "access_request_review_modal",
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				section,
				radioBlock,
				reasonBlock,
			},
		},
	}

	return &modal, nil
}

func (a *accessRequestReviewBuilder) BuildSectionBlock() *slack.SectionBlock {
	text := fmt.Sprintf(
		"👤 Requestor        : %s\n"+
			"🎯 Requested Role   : %s\n"+
			"📝 Request Reason   : %s\n"+
			"💬 Origin Channel   : #%s\n"+
			"📡 Review Channel   : #%s\n"+
			"⏳ Request Expiry   : %s\n"+
			"⏰ Role Expiry      : %s",
		a.slackUser.RealName,
		a.accessRequest.Role,
		a.accessRequest.Reason,
		a.accessRequest.InputChannelName,
		a.accessRequest.ReviewChannelName,
		a.accessRequest.Expires.Format(time.RFC3339),
		a.accessRequest.AccessDuration.Format(time.RFC3339),
	)

	section := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("```\n%s\n```", text), false, false),
		nil, nil,
	)
	return section
}

func (a *accessRequestReviewBuilder) BuildRadioBlock() *slack.InputBlock {
	radioOptions := []*slack.OptionBlockObject{
		slack.NewOptionBlockObject("allow", slack.NewTextBlockObject("plain_text", "✅ Allow", false, false), nil),
		slack.NewOptionBlockObject("deny", slack.NewTextBlockObject("plain_text", "⛔ Deny", false, false), nil),
	}
	radioElement := slack.NewRadioButtonsBlockElement("review_decision", radioOptions...)
	radioBlock := slack.NewInputBlock(
		"review_radio",
		slack.NewTextBlockObject("plain_text", "Choose Action", false, false),
		nil,
		radioElement,
	)
	return radioBlock
}

func (a *accessRequestReviewBuilder) BuildReasonBlock() *slack.InputBlock {
	reasonElement := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject("plain_text", "Write reason", false, false),
		"review_reason",
	)
	reasonBlock := slack.NewInputBlock(
		"reason_input",
		slack.NewTextBlockObject("plain_text", "Review Reason", false, false),
		nil,
		reasonElement,
	)
	return reasonBlock
}
