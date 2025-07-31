package models

import (
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/user/models"
	"time"
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
