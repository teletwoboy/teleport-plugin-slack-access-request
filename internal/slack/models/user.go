package models

import "time"

type User struct {
	SlackUserID int32
	ID          string
	Name        string
	RealName    string
	Email       string
	UseYn       bool
	CreateCode  string
	CreateDate  time.Time
	UpdateCode  string
	UpdateDate  time.Time
	DeleteCode  string
	DeleteDate  time.Time
	Version     int64
}
