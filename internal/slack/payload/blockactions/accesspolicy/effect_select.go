package accesspolicy

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/util"
	"time"
)

type EffectSelectPayload struct {
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

type EffectSelectPrivateMetadataPayload struct {
	ChannelID           string `json:"channel_id"`
	ChannelName         string `json:"channel_name"`
	RealName            string `json:"real_name"`
	TimeZone            string `json:"time_zone"`
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

type EffectSelect struct {
	RequesterChannelID   string
	RequesterChannelName string
	RequesterRealName    string
	RequesterTimeZone    string
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
	TriggerID            string
	ViewHash             string
	ViewID               string
	// new field
	Effect string
}

func ParseEffectSelect(payloadStr string) (*EffectSelect, error) {
	var payload EffectSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %w", err)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata EffectSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid payload format: %w", err)
	}

	startDate, err := parseStartDate(privateMetadata)
	if err != nil {
		return nil, fmt.Errorf("invalid payload format: %w", err)
	}

	endDate, err := parseEndDate(privateMetadata)
	if err != nil {
		return nil, fmt.Errorf("invalid payload format: %w", err)
	}
	return &EffectSelect{
		RequesterChannelID:   privateMetadata.ChannelID,
		RequesterChannelName: privateMetadata.ChannelName,
		RequesterRealName:    privateMetadata.RealName,
		RequesterTimeZone:    privateMetadata.TimeZone,
		RequesterID:          payload.User.ID,
		RequesterName:        payload.User.Name,
		SelectedChannelID:    privateMetadata.SelectedChannelID,
		SelectedChannelName:  privateMetadata.SelectedChannelName,
		SelectedRole:         privateMetadata.SelectedRole,
		SelectedRoleName:     privateMetadata.SelectedRoleName,
		SelectedUserID:       privateMetadata.SelectedUserID,
		SelectedRealName:     privateMetadata.SelectedRealName,
		SelectedStartDate:    startDate,
		SelectedEndDate:      endDate,
		TriggerID:            payload.TriggerID,
		ViewHash:             payload.View.Hash,
		ViewID:               payload.View.ID,
		Effect:               payload.Actions[0].Value,
	}, nil
}

func parseStartDate(pm EffectSelectPrivateMetadataPayload) (time.Time, error) {
	loc, err := time.LoadLocation(pm.TimeZone)
	if err != nil {
		return time.Time{}, err
	}

	startDateStr := pm.SelectedStartDate + " " + pm.SelectedStartTime
	startDate, err := time.ParseInLocation(util.MinuteTimeFormat, startDateStr, loc)
	if err != nil {
		return time.Time{}, err
	}
	return startDate.UTC(), nil
}

func parseEndDate(pm EffectSelectPrivateMetadataPayload) (time.Time, error) {
	loc, err := time.LoadLocation(pm.TimeZone)
	if err != nil {
		return time.Time{}, err
	}

	endDateStr := pm.SelectedEndDate + " " + pm.SelectedEndTime
	endDate, err := time.ParseInLocation(util.MinuteTimeFormat, endDateStr, loc)
	if err != nil {
		return time.Time{}, err
	}
	return endDate.UTC(), nil
}
