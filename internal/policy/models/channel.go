package models

import "time"

type Channel struct {
	PolicyChannelID int32
	AccessPolicy    *AccessPolicy
	ChannelID       string
	ChannelName     string
	UseYn           bool
	CreateCode      string
	CreateDate      time.Time
	UpdateCode      string
	UpdateDate      time.Time
	DeleteCode      string
	DeleteDate      time.Time
	Version         int64
}
