package message

import "github.com/slack-go/slack"

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
