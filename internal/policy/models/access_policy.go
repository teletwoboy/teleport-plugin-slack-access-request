package models

import "time"

type AccessPolicy struct {
	AccessPolicyID int32
	UserID         int32
	ChannelID      string
	ChannelName    string
	Name           string
	Reason         string
	StartTime      time.Time
	EndTime        time.Time
	AutoReview     bool
	Effect         string
	UseYn          bool
	CreateCode     string
	CreateDate     time.Time
	UpdateCode     string
	UpdateDate     time.Time
	DeleteCode     string
	DeleteDate     time.Time
	Version        int64
}
