package database

import "time"

const (
	CreateCode = "teleport-plugin-slack-access-request"
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

func PrePersist() *BaseEntity {
	return &BaseEntity{
		UseYn:      true,
		CreateCode: CreateCode,
		CreateDate: time.Now(),
		Version:    0,
	}
}
