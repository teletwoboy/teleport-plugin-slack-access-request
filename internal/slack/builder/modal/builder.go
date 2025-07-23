package modal

import "github.com/slack-go/slack"

type Builder interface {
	Build() (*slack.ModalViewRequest, error)
}
