package slack

import "github.com/slack-go/slack"

/*
MessageBuilder 는 빌더 패턴을 따름
Message의 종류는 여러개이기 때문에 각각 Message를 보내기 위해
service.go 에서 모든 것을 구현하는 것은 매우 복잡한 일임

service 에선 하나의 Message 보내는 메서드( PostMessage() )만 만들고,
여러 가지 Message 종류를 받아서 Build 후 사용하기 위함
*/
type MessageBuilder interface {
	Build() slack.MsgOption
}

type UserNotFoundMessageBuilder struct {
	slackName string
}

func NewUserNotFoundMessageBuilder(slackName string) *UserNotFoundMessageBuilder {
	return &UserNotFoundMessageBuilder{
		slackName: slackName,
	}
}

func (u *UserNotFoundMessageBuilder) Build() slack.MsgOption {
	text := ":warning: User Not Found \n```\n SlackName : " + u.slackName + "```"
	return slack.MsgOptionText(text, false)
}

type ErrorMessageBuilder struct {
	Err error
}

func NewErrorMessageBuilder(err error) *ErrorMessageBuilder {
	return &ErrorMessageBuilder{
		Err: err,
	}
}

func (e *ErrorMessageBuilder) Build() slack.MsgOption {
	text := ":warning: Error Occurred \n```\n Error : " + e.Err.Error() + "```"
	return slack.MsgOptionText(text, false)
}
