package message

import (
	"fmt"
	"github.com/slack-go/slack"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"time"
)

type accessRequestSubmissionBuilder struct {
	accessRequest *teleportmodels.AccessRequest
	slackUser     *slackmodels.User
}

func NewAccessRequestSubmissionBuilder(a *teleportmodels.AccessRequest, s *slackmodels.User) Builder {
	return &accessRequestSubmissionBuilder{
		accessRequest: a,
		slackUser:     s,
	}
}

func (a *accessRequestSubmissionBuilder) Build() slack.MsgOption {
	text := "*🔐 Successfully submitted Access Request*\n"
	text += "\n```\n"
	text += fmt.Sprintf("👤 Request User    : %s\n", a.slackUser.RealName)
	text += fmt.Sprintf("🎯 Request Role    : %s\n", a.accessRequest.Role)
	text += fmt.Sprintf("📡 Request Channel : #%s\n", a.accessRequest.ReviewChannelName)
	text += fmt.Sprintf("📝 Request Reason  : %s\n", a.accessRequest.Reason)
	text += "```\n"
	return slack.MsgOptionText(text, false)
}

// -- To reviewers
type accessRequestToReviewersBuilder struct {
	accessRequest *teleportmodels.AccessRequest
	slackUser     *slackmodels.User
}

func NewAccessRequestToReviewersBuilder(a *teleportmodels.AccessRequest, s *slackmodels.User) Builder {
	return &accessRequestToReviewersBuilder{
		accessRequest: a,
		slackUser:     s,
	}
}

func (a *accessRequestToReviewersBuilder) Build() slack.MsgOption {
	text := "*🔐 Someone submitted Access Request*\n"
	text += "\n```\n"
	text += fmt.Sprintf("👤 Requestor        : %s\n", a.slackUser.RealName)
	text += fmt.Sprintf("🎯 Requested Role   : %s\n", a.accessRequest.Role)
	text += fmt.Sprintf("📝 Request Reason   : %s\n", a.accessRequest.Reason)
	text += fmt.Sprintf("💬 Origin Channel   : #%s\n", a.accessRequest.InputChannelName)
	text += fmt.Sprintf("📡 Review Channel   : #%s\n", a.accessRequest.ReviewChannelName)
	text += fmt.Sprintf("⏳ Request Expiry   : %s\n", a.accessRequest.Expires.Format(time.RFC3339))
	text += fmt.Sprintf("⏰ Role Expiry      : %s\n", a.accessRequest.AccessDuration.Format(time.RFC3339))
	text += "```"
	text += "\n👉 Click the button below to review this request."

	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", text, false, false),
			nil,
			nil,
		),
		slack.NewActionBlock(
			"access_request_actions",
			slack.NewButtonBlockElement(
				"open_access_request_review_modal",
				a.accessRequest.Name,
				slack.NewTextBlockObject("plain_text", "Review Request", false, false),
			).WithStyle("primary"),
		),
	}

	return slack.MsgOptionBlocks(blocks...)
}
