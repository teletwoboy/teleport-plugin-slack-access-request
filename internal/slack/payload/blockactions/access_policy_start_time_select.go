package blockactions

import (
	"encoding/json"
	"fmt"
)

type AccessPolicyStartTimeSelectPayload struct {
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
				StartDateTimeBlock struct {
					AccessPolicyStartTimeSelect struct {
						Type         string `json:"type"`
						SelectedTime string `json:"selected_time"`
					} `json:"access_policy_start_time_select"`
				} `json:"start_date_time_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type AccessPolicyStartTimeSelectPrivateMetadataPayload struct {
	ChannelID           string `json:"channel_id"`
	ChannelName         string `json:"channel_name"`
	RealName            string `json:"real_name"`
	SelectedChannelID   string `json:"selected_channel_id"`
	SelectedChannelName string `json:"selected_channel_name"`
	SelectedRole        string `json:"selected_role"`
	SelectedRoleName    string `json:"selected_role_name"`
	SelectedUserID      string `json:"selected_user_id"`
	SelectedRealName    string `json:"selected_real_name"`
	SelectedStartDate   string `json:"selected_start_date"`
}

type AccessPolicyStartTimeSelect struct {
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
	SelectedStartDate    string
	TriggerID            string
	ViewHash             string
	ViewID               string
	// new field
	StartTime string
}

func ParseAccessPolicyStartTimeSelect(payloadStr string) (*AccessPolicyStartTimeSelect, error) {
	var payload AccessPolicyStartTimeSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", err)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata AccessPolicyStartTimeSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", err)
	}
	return &AccessPolicyStartTimeSelect{
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
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
		StartTime:            payload.View.State.Values.StartDateTimeBlock.AccessPolicyStartTimeSelect.SelectedTime,
	}, nil
}
