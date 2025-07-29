package models

import "time"

type User struct {
	PolicyUserID int32
	AccessPolicy *AccessPolicy
	UseYn        bool
	CreateCode   string
	CreateDate   time.Time
	UpdateCode   string
	UpdateDate   time.Time
	DeleteCode   string
	DeleteDate   time.Time
	Version      int64
}
