/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/constant"
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
