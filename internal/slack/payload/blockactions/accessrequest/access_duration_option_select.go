package accessrequest

import (
	"encoding/json"
	"fmt"
	"time"
)

type AccessDurationOptionSelectPayload struct {
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
				AccessDurationOptionBlock struct {
					AccessRequestAccessDurationOptionSelect struct {
						Type string `json:"type"`

						SelectedOption *struct {
							Value string `json:"value"`

							Text struct {
								Type string `json:"type"`
								Text string `json:"text"`
							} `json:"text"`
						} `json:"selected_option,omitempty"`
					} `json:"access_request_access_duration_option_select"`
				} `json:"access_duration_option_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type AccessDurationOptionSelectPrivateMetadataPayload struct {
	ChannelID                   string    `json:"channel_id"`
	ChannelName                 string    `json:"channel_name"`
	RealName                    string    `json:"real_name"`
	RequireReason               bool      `json:"require_reason"`
	SelectedRole                string    `json:"selected_role"`
	SelectedChannelID           string    `json:"selected_channel_id"`
	SelectedChannelName         string    `json:"selected_channel_name"`
	SelectedStartDateOptionID   string    `json:"selected_date_option_id"`
	SelectedStartDateOptionName string    `json:"selected_date_option_name"`
	TTL                         time.Time `json:"ttl"`
	SelectedStartDate           string    `json:"selected_start_date"`
	SelectedStartTime           string    `json:"selected_start_time"`
}

type AccessDurationOptionSelect struct {
	RequesterChannelID          string
	RequesterChannelName        string
	RequesterRealName           string
	RequireReason               bool
	SelectedRole                string
	SelectedChannelID           string
	SelectedChannelName         string
	SelectedStartDateOptionID   string
	SelectedStartDateOptionName string
	TTL                         time.Time
	SelectedStartDate           string
	SelectedStartTime           string
	RequesterID                 string
	RequesterName               string
	TriggerID                   string
	ViewHash                    string
	ViewID                      string
	// new fields
	AccessDurationOptionID   string
	AccessDurationOptionName string
}

func ParseAccessDurationOptionSelect(payloadStr string) (*AccessDurationOptionSelect, error) {
	var payload AccessDurationOptionSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", payloadStr)
	}
	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata AccessDurationOptionSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", strPrivateMetadata)
	}
	return &AccessDurationOptionSelect{
		RequesterChannelID:          privateMetadata.ChannelID,
		RequesterChannelName:        privateMetadata.ChannelName,
		RequesterRealName:           privateMetadata.RealName,
		RequireReason:               privateMetadata.RequireReason,
		SelectedRole:                privateMetadata.SelectedRole,
		SelectedChannelID:           privateMetadata.SelectedChannelID,
		SelectedChannelName:         privateMetadata.SelectedChannelName,
		SelectedStartDateOptionID:   privateMetadata.SelectedStartDateOptionID,
		SelectedStartDateOptionName: privateMetadata.SelectedStartDateOptionName,
		TTL:                         privateMetadata.TTL,
		SelectedStartDate:           privateMetadata.SelectedStartDate,
		SelectedStartTime:           privateMetadata.SelectedStartTime,
		RequesterID:                 payload.User.ID,
		RequesterName:               payload.User.Name,
		TriggerID:                   payload.TriggerID,
		ViewHash:                    payload.View.Hash,
		ViewID:                      payload.View.ID,
		AccessDurationOptionID:      payload.View.State.Values.AccessDurationOptionBlock.AccessRequestAccessDurationOptionSelect.SelectedOption.Value,
		AccessDurationOptionName:    payload.View.State.Values.AccessDurationOptionBlock.AccessRequestAccessDurationOptionSelect.SelectedOption.Text.Text,
	}, nil
}
