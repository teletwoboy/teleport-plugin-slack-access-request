package accessrequest

import (
	"encoding/json"
	"fmt"
	"time"
)

type AccessDurationDateSelectPayload struct {
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
				AccessDurationDateTimeBlock struct {
					AccessRequestAccessDurationDateSelect struct {
						Type         string `json:"type"`
						SelectedDate string `json:"selected_date"`
					} `json:"access_request_access_duration_date_select"`
				} `json:"access_duration_date_time_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type AccessDurationDateSelectPrivateMetadataPayload struct {
	ChannelID                        string    `json:"channel_id"`
	ChannelName                      string    `json:"channel_name"`
	RealName                         string    `json:"real_name"`
	RequireReason                    bool      `json:"require_reason"`
	SelectedRole                     string    `json:"selected_role"`
	SelectedChannelID                string    `json:"selected_channel_id"`
	SelectedChannelName              string    `json:"selected_channel_name"`
	SelectedStartDateOptionID        string    `json:"selected_date_option_id"`
	SelectedStartDateOptionName      string    `json:"selected_date_option_name"`
	TTL                              time.Time `json:"ttl"`
	SelectedStartDate                string    `json:"selected_start_date"`
	SelectedStartTime                string    `json:"selected_start_time"`
	SelectedAccessDurationOptionID   string    `json:"selected_access_duration_option_id"`
	SelectedAccessDurationOptionName string    `json:"selected_access_duration_option_name"`
}

type AccessDurationDateSelect struct {
	RequesterChannelID               string
	RequesterChannelName             string
	RequesterRealName                string
	RequireReason                    bool
	SelectedRole                     string
	SelectedChannelID                string
	SelectedChannelName              string
	SelectedStartDateOptionID        string
	SelectedStartDateOptionName      string
	TTL                              time.Time
	SelectedStartDate                string
	SelectedStartTime                string
	SelectedAccessDurationOptionID   string
	SelectedAccessDurationOptionName string
	RequesterID                      string
	RequesterName                    string
	TriggerID                        string
	ViewHash                         string
	ViewID                           string
	// new fields
	AccessDurationDate string
}

func ParseAccessDurationDateSelect(payloadStr string) (*AccessDurationDateSelect, error) {
	var payload AccessDurationDateSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", payloadStr)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata AccessDurationDateSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", strPrivateMetadata)
	}
	return &AccessDurationDateSelect{
		RequesterChannelID:               privateMetadata.ChannelID,
		RequesterChannelName:             privateMetadata.ChannelName,
		RequesterRealName:                privateMetadata.RealName,
		RequireReason:                    privateMetadata.RequireReason,
		SelectedRole:                     privateMetadata.SelectedRole,
		SelectedChannelID:                privateMetadata.SelectedChannelID,
		SelectedChannelName:              privateMetadata.SelectedChannelName,
		SelectedStartDateOptionID:        privateMetadata.SelectedStartDateOptionID,
		SelectedStartDateOptionName:      privateMetadata.SelectedStartDateOptionName,
		TTL:                              privateMetadata.TTL,
		SelectedStartDate:                privateMetadata.SelectedStartDate,
		SelectedStartTime:                privateMetadata.SelectedStartTime,
		SelectedAccessDurationOptionID:   privateMetadata.SelectedAccessDurationOptionID,
		SelectedAccessDurationOptionName: privateMetadata.SelectedAccessDurationOptionName,
		RequesterID:                      payload.User.ID,
		RequesterName:                    payload.User.Name,
		TriggerID:                        payload.TriggerID,
		ViewHash:                         payload.View.Hash,
		ViewID:                           payload.View.ID,
		AccessDurationDate:               payload.View.State.Values.AccessDurationDateTimeBlock.AccessRequestAccessDurationDateSelect.SelectedDate,
	}, nil
}
