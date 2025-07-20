package message

import (
	"github.com/slack-go/slack"
)

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
