package message

import (
	"fmt"

	"github.com/slack-go/slack"
)

type successCreateUser struct {
	realName string
	username string
}

func NewSuccessCreateUser(realName, username string) Builder {
	return &successCreateUser{realName: realName, username: username}
}

func (u *successCreateUser) Build() slack.MsgOption {
	text := "*🤗 Successfully Added User*\n"
	text += "\n```\n"
	text += fmt.Sprintf("Slack Name        : %s\n", u.realName)
	text += fmt.Sprintf("Teleport Username : %s\n", u.username)
	text += "```\n"
	return slack.MsgOptionText(text, false)
}

type successDeleteUser struct {
	realName string
	username string
}

func NewSuccessDeleteUser(realName, username string) Builder {
	return &successDeleteUser{realName: realName, username: username}
}

func (u *successDeleteUser) Build() slack.MsgOption {
	text := "*🤗 Successfully Deleted User*\n"
	text += "\n```\n"
	text += fmt.Sprintf("Deleted Slack Name        : %s\n", u.realName)
	text += fmt.Sprintf("Deleted Teleport Username : %s\n", u.username)
	text += "```\n"
	return slack.MsgOptionText(text, false)
}
