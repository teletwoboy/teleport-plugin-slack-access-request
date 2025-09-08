package accessrequest

import (
	"encoding/json"
	"fmt"
)

type StartDateOptionSelectPayload struct {
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
				StartDateOptionBlock struct {
					AccessRequestStartDateOptionSelect struct {
						Type string `json:"type"`

						SelectedOption *struct {
							Value string `json:"value"`

							Text struct {
								Type string `json:"type"`
								Text string `json:"text"`
							} `json:"text"`
						} `json:"selected_option,omitempty"`
					} `json:"access_request_start_date_option_select"`
				} `json:"start_date_option_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type StartDateOptionSelectPrivateMetadataPayload struct {
	ChannelID           string `json:"channel_id"`
	ChannelName         string `json:"channel_name"`
	RealName            string `json:"real_name"`
	RequireReason       bool   `json:"require_reason"`
	SelectedRole        string `json:"selected_role"`
	SelectedChannelID   string `json:"selected_channel_id"`
	SelectedChannelName string `json:"selected_channel_name"`
}

type StartDateOptionSelect struct {
	RequesterChannelID   string
	RequesterChannelName string
	RequesterID          string
	RequesterName        string
	RequesterRealName    string
	RequireReason        bool
	SelectedRole         string
	SelectedChannelID    string
	SelectedChannelName  string
	TriggerID            string
	ViewHash             string
	ViewID               string
	// new fields
	StartDateOptionID   string
	StartDateOptionName string
}

func ParseStartOptionSelect(payloadStr string) (*StartDateOptionSelect, error) {
	var payload StartDateOptionSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", payloadStr)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata StartDateOptionSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", strPrivateMetadata)
	}
	return &StartDateOptionSelect{
		RequesterChannelID:   privateMetadata.ChannelID,
		RequesterChannelName: privateMetadata.ChannelName,
		RequesterID:          payload.User.ID,
		RequesterName:        payload.User.Name,
		RequesterRealName:    privateMetadata.RealName,
		RequireReason:        privateMetadata.RequireReason,
		SelectedRole:         privateMetadata.SelectedRole,
		SelectedChannelID:    privateMetadata.SelectedChannelID,
		SelectedChannelName:  privateMetadata.SelectedChannelName,
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
		StartDateOptionID:    payload.View.State.Values.StartDateOptionBlock.AccessRequestStartDateOptionSelect.SelectedOption.Value,
		StartDateOptionName:  payload.View.State.Values.StartDateOptionBlock.AccessRequestStartDateOptionSelect.SelectedOption.Text.Text,
	}, nil
}
