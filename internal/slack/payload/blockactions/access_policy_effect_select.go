package blockactions

import (
	"encoding/json"
	"fmt"
)

type AccessPolicyEffectSelectPayload struct {
	Type      string `json:"type"`
	TriggerID string `json:"trigger_id"`

	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`

	Actions []struct {
		ActionID string `json:"action_id"`
		BlockID  string `json:"block_id"`
		Type     string `json:"type"`

		Text struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"text"`

		Value string `json:"value"`
	} `json:"actions"`

	View struct {
		ID              string `json:"id"`
		PrivateMetadata string `json:"private_metadata"`
		CallbackID      string `json:"callback_id"`
		Hash            string `json:"hash"`
	} `json:"view"`
}

type AccessPolicyEffectSelectPrivateMetadataPayload struct {
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
	SelectedStartTime   string `json:"selected_start_time"`
	SelectedEndDate     string `json:"selected_end_date"`
	SelectedEndTime     string `json:"selected_end_time"`
	SelectedButtonValue string `json:"selected_button_value"`
}

type AccessPolicyEffectSelect struct {
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
	SelectedStartTime    string
	SelectedEndDate      string
	SelectedEndTime      string
	TriggerID            string
	ViewHash             string
	ViewID               string
	// new field
	Effect string
}

func ParseAccessPolicyEffectSelect(payloadStr string) (*AccessPolicyEffectSelect, error) {
	var payload AccessPolicyEffectSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", err)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata AccessPolicyEffectSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", err)
	}
	return &AccessPolicyEffectSelect{
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
		SelectedStartTime:    privateMetadata.SelectedStartTime,
		SelectedEndDate:      privateMetadata.SelectedEndDate,
		SelectedEndTime:      privateMetadata.SelectedEndTime,
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
		Effect:               payload.Actions[0].Value,
	}, nil
}
