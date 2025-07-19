package message

import (
	"fmt"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"time"

	"github.com/slack-go/slack"
)

/*
Builder 는 빌더 패턴을 따름
Message의 종류는 여러개이기 때문에 각각 Message를 보내기 위해
service.go 에서 모든 것을 구현하는 것은 매우 복잡한 일임

service 에선 하나의 Message 보내는 메서드( PostMessage() )만 만들고,
여러 가지 Message 종류를 받아서 Build 후 사용하기 위함
*/
type Builder interface {
	Build() slack.MsgOption
}

type userNotFoundBuilder struct {
	slackName string
}

func NewUserNotFoundBuilder(slackName string) Builder {
	return &userNotFoundBuilder{
		slackName: slackName,
	}
}

func (u *userNotFoundBuilder) Build() slack.MsgOption {
	text := ":warning: User Not Found \n```\n SlackName : " + u.slackName + "```"
	return slack.MsgOptionText(text, false)
}

// ---------------------------------------------------------------------------------------------

type errorBuilder struct {
	Err error
}

func NewErrorMessageBuilder(err error) Builder {
	return &errorBuilder{
		Err: err,
	}
}

func (e *errorBuilder) Build() slack.MsgOption {
	text := ":warning: Error Occurred \n```\n Error : " + e.Err.Error() + "```"
	return slack.MsgOptionText(text, false)
}

// -----------------------------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------------------------

type accessRequestNotFoundBuilder struct {
	name string
}

func NewAccessRequestNotFoundBuilder(name string) Builder {
	return &accessRequestNotFoundBuilder{
		name: name,
	}
}

func (a *accessRequestNotFoundBuilder) Build() slack.MsgOption {
	text := ":warning: Access Request Not Found \n```\n Name : " + a.name + "```"
	return slack.MsgOptionText(text, false)
}

// -----------------------------------------------------------------------------------------------

type accessRequestAlreadyApprovedBuilder struct {
	name string
}

func NewAccessRequestAlreadyApprovedBuilder(name string) Builder {
	return &accessRequestAlreadyApprovedBuilder{
		name: name,
	}
}

func (a *accessRequestAlreadyApprovedBuilder) Build() slack.MsgOption {
	text := ":warning: Access Request Already Approved \n```\n Name : " + a.name + "```"
	return slack.MsgOptionText(text, false)
}
