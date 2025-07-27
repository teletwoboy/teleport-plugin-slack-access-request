package blockactions

import (
	"encoding/json"
	"fmt"
)

type AccessPolicyChannelSelectPayload struct {
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

	Email string
}

type AccessPolicyChannelSelectPrivateMetadataPayload struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}

type AccessPolicyChannelSelect struct {
	ChannelID            string
	ChannelName          string
	Email                string
	RequesterChannelID   string
	RequesterChannelName string
	RequesterID          string
	RequesterName        string
	TriggerID            string
	ViewHash             string
	ViewID               string
}

func ParseAccessPolicyChannelSelect(payloadStr string) (*AccessPolicyChannelSelect, error) {
	var payload AccessPolicyChannelSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", err)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata AccessPolicyChannelSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", err)
	}

	return &AccessPolicyChannelSelect{
		ChannelID:            payload.View.State.Values.ChannelBlock.AccessPolicyChannelSelect.SelectedOption.Value,
		ChannelName:          payload.View.State.Values.ChannelBlock.AccessPolicyChannelSelect.SelectedOption.Text.Text,
		RequesterChannelID:   privateMetadata.ChannelID,
		RequesterChannelName: privateMetadata.ChannelName,
		RequesterID:          payload.User.ID,
		RequesterName:        payload.User.Name,
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
	}, nil
}
