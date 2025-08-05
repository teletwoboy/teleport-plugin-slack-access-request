package message

import "github.com/slack-go/slack"

type Builder interface {
	Build() slack.MsgOption
}
