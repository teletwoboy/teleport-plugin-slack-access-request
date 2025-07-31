package blockactions

import (
	"encoding/json"
	"fmt"
)

type ChannelSelectPayload struct {
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
				ChannelBlock struct {
					AccessRequestChannelSelect struct {
						Type string `json:"type"`

						SelectedOption *struct {
							Value string `json:"value"`

							Text struct {
								Type string `json:"type"`
								Text string `json:"text"`
							} `json:"text"`
						} `json:"selected_option,omitempty"`
					} `json:"access_request_channel_select"`
				} `json:"channel_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type ChannelSelectPrivateMetadataPayload struct {
	ChannelID     string `json:"channel_id"`
	ChannelName   string `json:"channel_name"`
	RealName      string `json:"real_name"`
	RequireReason bool   `json:"require_reason"`
	SelectedRole  string `json:"selected_role"`
}

type ChannelSelect struct {
	RequesterChannelID   string
	RequesterChannelName string
	RequesterRealName    string
	RequireReason        bool
	RequesterID          string
	RequesterName        string
	SelectedRole         string
	TriggerID            string
	ViewHash             string
	ViewID               string
	// new fields
	ChannelID   string
	ChannelName string
}

func ParseChannelSelect(payloadStr string) (*ChannelSelect, error) {
	var payload ChannelSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", payloadStr)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata ChannelSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", strPrivateMetadata)
	}
	return &ChannelSelect{
		RequesterChannelID:   privateMetadata.ChannelID,
		RequesterChannelName: privateMetadata.ChannelName,
		RequesterRealName:    privateMetadata.RealName,
		RequireReason:        privateMetadata.RequireReason,
		RequesterID:          payload.User.ID,
		RequesterName:        payload.User.Name,
		SelectedRole:         privateMetadata.SelectedRole,
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
		ChannelID:            payload.View.State.Values.ChannelBlock.AccessRequestChannelSelect.SelectedOption.Value,
		ChannelName:          payload.View.State.Values.ChannelBlock.AccessRequestChannelSelect.SelectedOption.Text.Text,
	}, nil
}
