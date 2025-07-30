package accesspolicy

import (
	"encoding/json"
	"fmt"
)

type TimeZoneSelectPayload struct {
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
				TimeZoneBlock struct {
					AccessPolicyTimeZoneSelect struct {
						Type string `json:"type"`

						SelectedOption struct {
							Value string `json:"value"`

							Text struct {
								Type string `json:"type"`
								Text string `json:"text"`
							} `json:"text"`
						} `json:"selected_option"`
					} `json:"access_policy_time_zone_select"`
				} `json:"time_zone_block"`
			} `json:"Values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type TimeZoneSelectPrivateMetadataPayload struct {
	ChannelID           string `json:"channel_id"`
	ChannelName         string `json:"channel_name"`
	RealName            string `json:"real_name"`
	SelectedChannelID   string `json:"selected_channel_id"`
	SelectedChannelName string `json:"selected_channel_name"`
	SelectedRole        string `json:"selected_role"`
	SelectedRoleName    string `json:"selected_role_name"`
	SelectedUserID      string `json:"selected_user_id"`
	SelectedRealName    string `json:"selected_real_name"`
}

type TimeZoneSelect struct {
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
	TriggerID            string
	ViewHash             string
	ViewID               string
	// new fields
	TimeZone string
}

func ParseTimeZoneSelect(payloadStr string) (*TimeZoneSelect, error) {
	var payload TimeZoneSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %w", err)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata TimeZoneSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid payload format: %w", err)
	}
	return &TimeZoneSelect{
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
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
		TimeZone:             payload.View.State.Values.TimeZoneBlock.AccessPolicyTimeZoneSelect.SelectedOption.Value,
	}, nil
}
