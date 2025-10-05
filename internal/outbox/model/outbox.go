package model

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"time"
)

type Outbox struct {
	OutboxID      int32
	EventType     string
	AggregateType string
	AggregateID   int32
	Payload       string
	Status        string
	Attempts      int32
	NextTryAt     time.Time
	LastError     string
	UseYn         bool
	CreateCode    string
	CreateDate    time.Time
	UpdateCode    string
	UpdateDate    time.Time
	DeleteCode    string
	DeleteDate    time.Time
	Version       int64
}

type OutboxNotification struct {
	Channel string
	Payload string
}

type OutboxNotificationPayload struct {
	OutboxID int32
}

func NewOutboxNotification(ob *Outbox) (*OutboxNotification, error) {
	p := OutboxNotificationPayload{
		OutboxID: ob.OutboxID,
	}
	marshaled, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal outbox notification payload: %w", err)
	}
	return &OutboxNotification{
		Channel: constant.OutboxChannel,
		Payload: string(marshaled),
	}, nil
}
