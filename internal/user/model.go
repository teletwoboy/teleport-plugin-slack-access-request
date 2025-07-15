package user

import (
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
	"time"
)

type User struct {
	UserID       int32
	TeleportUser *teleport.User
	SlackUser    *slack.User
	UseYn        bool
	CreateCode   string
	CreateDate   time.Time
	UpdateCode   string
	UpdateDate   time.Time
	DeleteCode   string
	DeleteDate   time.Time
	Version      int64
}
