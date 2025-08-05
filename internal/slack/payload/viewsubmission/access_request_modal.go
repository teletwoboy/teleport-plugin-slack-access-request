package viewsubmission

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions"
)

type AccessRequestModalPayload struct {
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
					AccessRequestReasonInput struct {
						Type  string `json:"type"`
						Value string `json:"value,omitempty"`
					} `json:"access_request_reason_input"`
				} `json:"reason_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`

	Email string
}

type AccessRequestModal struct {
	RequesterChannelID   string
	RequesterChannelName string
	RequesterRealName    string
	RequireReason        bool
	RequesterID          string
	RequesterName        string
	SelectedRole         string
	SelectedChannelID    string
	SelectedChannelName  string
	TriggerID            string
	ViewHash             string
	ViewID               string
	// new fields
	Reason string
}

func ParseAccessRequestModal(payloadStr string) (*AccessRequestModal, error) {
	var payload AccessRequestModalPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", payloadStr)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata blockactions.SummaryPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", strPrivateMetadata)
	}
	return &AccessRequestModal{
		RequesterChannelID:   privateMetadata.ChannelID,
		RequesterChannelName: privateMetadata.ChannelName,
		RequesterRealName:    privateMetadata.RealName,
		RequireReason:        privateMetadata.RequireReason,
		RequesterID:          payload.User.ID,
		RequesterName:        payload.User.Name,
		SelectedRole:         privateMetadata.SelectedRole,
		SelectedChannelID:    privateMetadata.SelectedChannelID,
		SelectedChannelName:  privateMetadata.SelectedChannelName,
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
		Reason:               payload.View.State.Values.ReasonBlock.AccessRequestReasonInput.Value,
	}, nil
}
