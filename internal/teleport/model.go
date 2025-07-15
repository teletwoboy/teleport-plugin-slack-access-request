package teleport

import "time"

type User struct {
	TeleportUserID int32
	Username       string
	UseYn          bool
	CreateCode     string
	CreateDate     time.Time
	UpdateCode     string
	UpdateDate     time.Time
	DeleteCode     string
	DeleteDate     time.Time
	Version        int64
}

type UserAccessInfo struct {
	Roles         []string
	RequireReason bool
}
