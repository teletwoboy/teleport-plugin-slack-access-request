package accesspolicy

import (
	"encoding/json"
	"fmt"
)

type RoleSelectPayload struct {
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
				RoleBlock struct {
					AccessPolicyRoleSelect struct {
						Type string `json:"type"`

						SelectedOption struct {
							Value string `json:"value"`

							Text struct {
								Type string `json:"type"`
								Text string `json:"text"`
							} `json:"text"`
						} `json:"selected_option"`
					} `json:"access_policy_role_select"`
				} `json:"role_block"`
			} `json:"Values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type RoleSelectPrivateMetadataPayload struct {
	ChannelID           string `json:"channel_id"`
	ChannelName         string `json:"channel_name"`
	RealName            string `json:"real_name"`
	TimeZone            string `json:"time_zone"`
	SelectedChannelID   string `json:"selected_channel_id"`
	SelectedChannelName string `json:"selected_channel_name"`
}

type RoleSelect struct {
	RequesterChannelID   string
	RequesterChannelName string
	RequesterRealName    string
	RequesterTimeZone    string
	RequesterID          string
	RequesterName        string
	SelectedChannelID    string
	SelectedChannelName  string
	TriggerID            string
	ViewHash             string
	ViewID               string
	// new fields
	Role     string
	RoleName string
}

func ParseRoleSelect(payloadStr string) (*RoleSelect, error) {
	var payload RoleSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %w", err)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata RoleSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %w", err)
	}

	return &RoleSelect{
		RequesterChannelID:   privateMetadata.ChannelID,
		RequesterChannelName: privateMetadata.ChannelName,
		RequesterRealName:    privateMetadata.RealName,
		RequesterTimeZone:    privateMetadata.TimeZone,
		RequesterID:          payload.User.ID,
		RequesterName:        payload.User.Name,
		SelectedChannelID:    privateMetadata.SelectedChannelID,
		SelectedChannelName:  privateMetadata.SelectedChannelName,
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
		Role:                 payload.View.State.Values.RoleBlock.AccessPolicyRoleSelect.SelectedOption.Value,
		RoleName:             payload.View.State.Values.RoleBlock.AccessPolicyRoleSelect.SelectedOption.Text.Text,
	}, nil
}
