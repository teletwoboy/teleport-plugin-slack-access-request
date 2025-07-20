package models

import "time"

type AccessRequest struct {
	AccessRequestID   int32
	RequesterUserID   int32
	Name              string
	InputChannelID    string
	InputChannelName  string
	Role              string
	Reason            string
	ReviewChannelID   string
	ReviewChannelName string
	Status            string
	Expires           time.Time
	SessionTTL        time.Time
	AccessDuration    time.Time
	StartDate         time.Time
	ExpiryDate        time.Time
	UseYn             bool
	CreateCode        string
	CreateDate        time.Time
	UpdateCode        string
	UpdateDate        time.Time
	DeleteCode        string
	DeleteDate        time.Time
	Version           int64
}
