package message

import "github.com/slack-go/slack"

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
