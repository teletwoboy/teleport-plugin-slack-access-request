package models

import (
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"time"
)

type User struct {
	UserID       int32
	TeleportUser *teleportmodels.User
	SlackUser    *slackmodels.User
	UseYn        bool
	CreateCode   string
	CreateDate   time.Time
	UpdateCode   string
	UpdateDate   time.Time
	DeleteCode   string
	DeleteDate   time.Time
	Version      int64
}
