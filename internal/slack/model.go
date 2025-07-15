package slack

import "time"

type User struct {
	SlackUserID int32
	ID          string
	Name        string
	RealName    string
	Email       string
	Deleted     bool
	UseYn       bool
	CreateCode  string
	CreateDate  time.Time
	UpdateCode  string
	UpdateDate  time.Time
	DeleteCode  string
	DeleteDate  time.Time
	Version     int64
}

type TeamInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ReviewersChannel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsMember bool   `json:"is_member"`
}
