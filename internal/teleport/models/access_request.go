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
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/util"
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

func NewAccessRequest(ar types.AccessRequest, payload *viewsubmission.AccessRequestModal, userID int32) *AccessRequest {
	return &AccessRequest{
		RequesterUserID:   userID,
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

func (ar *AccessRequest) UpdateState(effect string) {
	switch effect {
	case util.APolicyAllowButtonValue:
		ar.State = types.RequestState_APPROVED.String()
	case util.APolicyDenyButtonValue:
		ar.State = types.RequestState_DENIED.String()
	}
}
