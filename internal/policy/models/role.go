package models

import "time"

type Role struct {
	PolicyRoleID int32
	AccessPolicy *AccessPolicy
	Role         string
	UseYn        bool
	CreateCode   string
	CreateDate   time.Time
	UpdateCode   string
	UpdateDate   time.Time
	DeleteCode   string
	DeleteDate   time.Time
	Version      int64
}
