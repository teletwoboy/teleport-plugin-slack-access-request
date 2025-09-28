package model

import (
	"time"
)

type Outbox struct {
	OutboxID    int32
	EventType   string
	AggregateID int32
	Payload     string
	Status      string
	Attempts    int32
	ApiAttempts int32
	DBAttempts  int32
	NextTryAt   time.Time
	LastError   string
	UseYn       bool
	CreateCode  string
	CreateDate  time.Time
	UpdateCode  string
	UpdateDate  time.Time
	DeleteCode  string
	DeleteDate  time.Time
	Version     int64
}
