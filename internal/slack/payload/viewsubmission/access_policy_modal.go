package viewsubmission

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"
	"time"
)

const (
	layout = "2006-01-02 15:04"
)

type AccessPolicyModalPayload struct {
	Type      string `json:"type"`
	TriggerID string `json:"trigger_id"`

	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`

	View struct {
		ID              string `json:"id"`
		PrivateMetadata string `json:"private_metadata"`
		CallbackID      string `json:"callback_id"`

		State struct {
			Values struct {
				ReasonBlock struct {
					AccessPolicyReasonInput struct {
						Type  string `json:"type"`
						Value string `json:"value,omitempty"`
					} `json:"access_policy_reason_input"`
				} `json:"reason_block"`

				TitleBlock struct {
					AccessPolicyTitleInput struct {
						Type  string `json:"type"`
						Value string `json:"value,omitempty"`
					} `json:"access_policy_title_input"`
				} `json:"title_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type AccessPolicyModal struct {
	RequesterChannelID   string
	RequesterChannelName string
	RequesterRealName    string
	RequesterID          string
	RequesterName        string
	SelectedChannelID    string
	SelectedChannelName  string
	SelectedRole         string
	SelectedRoleName     string
	SelectedUserID       string
	SelectedRealName     string
	SelectedStartDate    time.Time
	SelectedEndDate      time.Time
	SelectedEffect       string
	TriggerID            string
	ViewHash             string
	ViewID               string
	// new fields
	Title  string
	Reason string
}

func ParseAccessPolicyModal(payloadStr string) (*AccessPolicyModal, error) {
	var payload AccessPolicyModalPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %w", err)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata accesspolicy.SummaryPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid payload format: %w", err)
	}
	return &AccessPolicyModal{
		RequesterChannelID:   privateMetadata.ChannelID,
		RequesterChannelName: privateMetadata.ChannelName,
		RequesterRealName:    privateMetadata.RealName,
		RequesterID:          payload.User.ID,
		RequesterName:        payload.User.Name,
		SelectedChannelID:    privateMetadata.SelectedChannelID,
		SelectedChannelName:  privateMetadata.SelectedChannelName,
		SelectedRole:         privateMetadata.SelectedRole,
		SelectedRoleName:     privateMetadata.SelectedRoleName,
		SelectedUserID:       privateMetadata.SelectedUserID,
		SelectedRealName:     privateMetadata.SelectedRealName,
		SelectedStartDate:    privateMetadata.SelectedStartDate,
		SelectedEndDate:      privateMetadata.SelectedEndDate,
		SelectedEffect:       privateMetadata.SelectedEffect,
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
		Title:                payload.View.State.Values.TitleBlock.AccessPolicyTitleInput.Value,
		Reason:               payload.View.State.Values.ReasonBlock.AccessPolicyReasonInput.Value,
	}, nil
}
