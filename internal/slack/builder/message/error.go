package message

import "github.com/slack-go/slack"

type errorBuilder struct {
	Err error
}

func NewErrorBuilder(err error) Builder {
	return &errorBuilder{
		Err: err,
	}
}

func (e *errorBuilder) Build() slack.MsgOption {
	text := ":warning: Error Occurred \n```\n [Error] " + e.Err.Error() + "```"
	return slack.MsgOptionText(text, false)
}
