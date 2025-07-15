package user

import (
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
)

type User struct {
	TeleportUser *teleport.User `json:"teleport_user"`
	SlackUser    *slack.User    `json:"slack_user"`
}
