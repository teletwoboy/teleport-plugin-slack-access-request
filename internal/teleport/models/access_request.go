package models

import (
	"teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"time"

	"github.com/gravitational/teleport/api/types"
)

type AccessRequest struct {
	AccessRequestID   int32
	RequesterUserID   int32
	Name              string
	InputChannelID    string
	InputChannelName  string
	Role              string
	Reason            string
	ReviewChannelID   string
	ReviewChannelName string
	State             string
	Expires           time.Time
	SessionTTL        time.Time
	AccessDuration    time.Time
	StartDate         time.Time
	ExpiryDate        time.Time
	UseYn             bool
	CreateCode        string
	CreateDate        time.Time
	UpdateCode        string
	UpdateDate        time.Time
	DeleteCode        string
	DeleteDate        time.Time
	Version           int64
}

func NewAccessRequest(ar types.AccessRequest, payload *viewsubmission.AccessRequestModal, slackUser *models.User) *AccessRequest {
	return &AccessRequest{
		RequesterUserID:   slackUser.SlackUserID,
		Name:              ar.GetName(),
		InputChannelID:    payload.RequesterChannelID,
		InputChannelName:  payload.RequesterChannelName,
		Role:              payload.SelectedRole,
		Reason:            payload.Reason,
		ReviewChannelID:   payload.SelectedChannelID,
		ReviewChannelName: payload.SelectedChannelName,
		State:             ar.GetState().String(),
		Expires:           ar.Expiry(),
		SessionTTL:        ar.GetSessionTLL(),
		AccessDuration:    ar.GetMaxDuration(),
		ExpiryDate:        ar.GetAccessExpiry(),
	}
}

func (ar *AccessRequest) Update(a types.AccessRequest) {
	ar.AccessDuration = a.GetMaxDuration()
	ar.Expires = a.Expiry()
	ar.ExpiryDate = a.GetAccessExpiry()
	ar.SessionTTL = a.GetSessionTLL()
	ar.StartDate = *a.GetAssumeStartTime()
	ar.State = a.GetState().String()
}
