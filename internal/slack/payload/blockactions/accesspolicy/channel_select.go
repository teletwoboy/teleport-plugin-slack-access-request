package accesspolicy

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
					AccessPolicyChannelSelect struct {
						Type string `json:"type"`

						SelectedOption struct {
							Value string `json:"value"`

							Text struct {
								Type string `json:"type"`
								Text string `json:"text"`
							} `json:"text"`
						} `json:"selected_option"`
					} `json:"access_policy_channel_select"`
				} `json:"channel_block"`
			} `json:"Values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type ChannelSelectPrivateMetadataPayload struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	RealName    string `json:"real_name"`
	TimeZone    string `json:"time_zone"`
}

type ChannelSelect struct {
	RequesterChannelID   string
	RequesterChannelName string
	RequesterRealName    string
	RequesterTimeZone    string
	RequesterID          string
	RequesterName        string
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
		return nil, fmt.Errorf("invalid payload format: %w", err)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata ChannelSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %w", err)
	}

	return &ChannelSelect{
		RequesterChannelID:   privateMetadata.ChannelID,
		RequesterChannelName: privateMetadata.ChannelName,
		RequesterRealName:    privateMetadata.RealName,
		RequesterTimeZone:    privateMetadata.TimeZone,
		RequesterID:          payload.User.ID,
		RequesterName:        payload.User.Name,
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
		ChannelID:            payload.View.State.Values.ChannelBlock.AccessPolicyChannelSelect.SelectedOption.Value,
		ChannelName:          payload.View.State.Values.ChannelBlock.AccessPolicyChannelSelect.SelectedOption.Text.Text,
	}, nil
}
