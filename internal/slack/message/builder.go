package message

import (
	"fmt"
	"strconv"
	"teleport-plugin-slack-access-request/internal/teleport/models"
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

type UserNotFoundBuilder struct {
	slackName string
}

func NewUserNotFoundBuilder(slackName string) *UserNotFoundBuilder {
	return &UserNotFoundBuilder{
		slackName: slackName,
	}
}

func (u *UserNotFoundBuilder) Build() slack.MsgOption {
	text := ":warning: User Not Found \n```\n SlackName : " + u.slackName + "```"
	return slack.MsgOptionText(text, false)
}

// ---------------------------------------------------------------------------------------------

type ErrorBuilder struct {
	Err error
}

func NewErrorMessageBuilder(err error) *ErrorBuilder {
	return &ErrorBuilder{
		Err: err,
	}
}

func (e *ErrorBuilder) Build() slack.MsgOption {
	text := ":warning: Error Occurred \n```\n Error : " + e.Err.Error() + "```"
	return slack.MsgOptionText(text, false)
}

// -----------------------------------------------------------------------------------------------

type AccessRequestSubmissionBuilder struct {
	Username      string
	AccessRequest *models.AccessRequest
}

func NewAccessRequestSubmissionBuilder(username string, a *models.AccessRequest) *AccessRequestSubmissionBuilder {
	return &AccessRequestSubmissionBuilder{
		Username:      username,
		AccessRequest: a,
	}
}

func (a *AccessRequestSubmissionBuilder) Build() slack.MsgOption {
	text := "*🔐 Successfully submitted Access Request*\n"
	text += "\n```\n"
	text += fmt.Sprintf("👤 Request User    : %s\n", a.Username)
	text += fmt.Sprintf("🎯 Request Role    : %s\n", a.AccessRequest.Role)
	text += fmt.Sprintf("📡 Request Channel : #%s\n", a.AccessRequest.ReviewChannelName)
	text += fmt.Sprintf("📝 Request Reason  : %s\n", a.AccessRequest.Reason)
	text += "```\n"
	return slack.MsgOptionText(text, false)
}

// -----------------------------------------------------------------------------------------------

type AccessRequestToReviewersBuilder struct {
	Username      string
	AccessRequest *models.AccessRequest
}

func NewAccessRequestToReviewersBuilder(username string, a *models.AccessRequest) *AccessRequestToReviewersBuilder {
	return &AccessRequestToReviewersBuilder{
		Username:      username,
		AccessRequest: a,
	}
}

func (a *AccessRequestToReviewersBuilder) Build() slack.MsgOption {
	text := "*🔐 Someone submitted Access Request*\n"
	text += "\n```\n"
	text += fmt.Sprintf("👤 Requestor        : %s\n", a.Username)
	text += fmt.Sprintf("🎯 Requested Role   : %s\n", a.AccessRequest.Role)
	text += fmt.Sprintf("📝 Request Reason   : %s\n", a.AccessRequest.Reason)
	text += fmt.Sprintf("💬 Origin Channel   : #%s\n", a.AccessRequest.InputChannelName)
	text += fmt.Sprintf("📡 Review Channel   : #%s\n", a.AccessRequest.ReviewChannelName)
	text += fmt.Sprintf("⏳ Request Expiry   : %s\n", a.AccessRequest.Expires.Format(time.RFC3339))
	text += fmt.Sprintf("⏰ Role Expiry      : %s\n", a.AccessRequest.AccessDuration.Format(time.RFC3339))
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
				"open_modal",
				strconv.Itoa(int(a.AccessRequest.AccessRequestID)),
				slack.NewTextBlockObject("plain_text", "Review Request", false, false),
			).WithStyle("primary"),
		),
	}

	return slack.MsgOptionBlocks(blocks...)
}
