package database

import "time"

const (
	CreateCode = "teleport-plugin-slack-access-request"
	UpdateCode = "teleport-plugin-slack-access-request"
)

type BaseEntity struct {
	UseYn      bool
	CreateCode string
	CreateDate time.Time
	UpdateCode string
	UpdateDate time.Time
	DeleteCode string
	DeleteDate time.Time
	Version    int64
}

func MarkCreate() *BaseEntity {
	return &BaseEntity{
		UseYn:      true,
		CreateCode: CreateCode,
		CreateDate: time.Now(),
		Version:    0,
	}
}

func MarkUpdate() *BaseEntity {
	return &BaseEntity{
		UpdateCode: UpdateCode,
		UpdateDate: time.Now(),
	}
}
