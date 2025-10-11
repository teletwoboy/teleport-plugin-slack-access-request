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

package models

import (
	"time"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/user/models"
)

type AccessPolicy struct {
	AccessPolicyID    int32
	UserID            int32
	InputChannelID    string
	InputChannelName  string
	Title             string
	Reason            string
	StartDate         time.Time
	EndDate           time.Time
	Effect            string
	TargetChannelID   string
	TargetChannelName string
	TargetRole        string
	TargetRoleName    string
	TargetSlackID     string
	TargetRealName    string
	MessageTimestamp  string
	UseYn             bool
	CreateCode        string
	CreateDate        time.Time
	UpdateCode        string
	UpdateDate        time.Time
	DeleteCode        string
	DeleteDate        time.Time
	Version           int64
}

func NewAccessPolicy(payload *viewsubmission.AccessPolicyModal, user *models.User) *AccessPolicy {
	return &AccessPolicy{
		UserID:            user.UserID,
		InputChannelID:    payload.RequesterChannelID,
		InputChannelName:  payload.RequesterChannelName,
		Title:             payload.Title,
		Reason:            payload.Reason,
		StartDate:         payload.SelectedStartDate,
		EndDate:           payload.SelectedEndDate,
		Effect:            payload.SelectedEffect,
		TargetChannelID:   payload.SelectedChannelID,
		TargetChannelName: payload.SelectedChannelName,
		TargetRole:        payload.SelectedRole,
		TargetRoleName:    payload.SelectedRoleName,
		TargetSlackID:     payload.SelectedUserID,
		TargetRealName:    payload.SelectedRealName,
	}
}

func (a *AccessPolicy) UpdateTimestamp(ts string) {
	a.MessageTimestamp = ts
}
