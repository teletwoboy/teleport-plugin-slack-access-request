package accessrequest

import (
	"encoding/json"
	"fmt"
	"time"
)

type RequestTTLTimeSelectPayload struct {
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
				RequestTTLDateTimeBlock struct {
					AccessRequestRequestTTLTimeSelect struct {
						Type         string `json:"type"`
						SelectedTime string `json:"selected_time"`
					} `json:"access_request_request_ttl_time_select"`
				} `json:"request_ttl_date_time_block"`
			} `json:"values"`
		} `json:"state"`

		Hash string `json:"hash"`
	} `json:"view"`
}

type RequestTTLTimeSelectPrivateMetadataPayload struct {
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
	SelectedAccessDurationDate       string    `json:"selected_access_duration_date"`
	SelectedAccessDurationTime       string    `json:"selected_access_duration_time"`
	RequestTTL                       time.Time `json:"request_ttl"`
	SelectedRequestTTLOptionID       string    `json:"request_ttl_option_id"`
	SelectedRequestTTLOptionName     string    `json:"request_ttl_option_name"`
	SelectedRequestTTLDate           string    `json:"request_ttl_date"`
}

type RequestTTLTimeSelect struct {
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
	SelectedAccessDurationDate       string
	SelectedAccessDurationTime       string
	RequestTTL                       time.Time
	SelectedRequestTTLOptionID       string
	SelectedRequestTTLOptionName     string
	SelectedRequestTTLDate           string
	RequesterID                      string
	RequesterName                    string
	TriggerID                        string
	ViewHash                         string
	ViewID                           string
	// new fields
	RequestTTLTime string
}

func NewRequestTTLTimeWithFirstOpt(p *RequestTTLOptionSelect) *RequestTTLTimeSelect {
	return &RequestTTLTimeSelect{
		RequesterChannelID:               p.RequesterChannelID,
		RequesterChannelName:             p.RequesterChannelName,
		RequesterRealName:                p.RequesterRealName,
		RequireReason:                    p.RequireReason,
		SelectedRole:                     p.SelectedRole,
		SelectedChannelID:                p.SelectedChannelID,
		SelectedChannelName:              p.SelectedChannelName,
		SelectedStartDateOptionID:        p.SelectedStartDateOptionID,
		SelectedStartDateOptionName:      p.SelectedStartDateOptionName,
		TTL:                              p.TTL,
		SelectedStartDate:                p.SelectedStartDate,
		SelectedStartTime:                p.SelectedStartTime,
		SelectedAccessDurationOptionID:   p.SelectedAccessDurationOptionID,
		SelectedAccessDurationOptionName: p.SelectedAccessDurationOptionName,
		SelectedAccessDurationDate:       p.SelectedAccessDurationDate,
		SelectedAccessDurationTime:       p.SelectedAccessDurationTime,
		RequestTTL:                       p.RequestTTL,
		SelectedRequestTTLOptionID:       p.RequestTTLOptionID,
		SelectedRequestTTLOptionName:     p.RequestTTLOptionName,
		RequesterID:                      p.RequesterID,
		RequesterName:                    p.RequesterName,
		TriggerID:                        p.TriggerID,
		ViewHash:                         p.ViewHash,
		ViewID:                           p.ViewID,
	}
}

func ParseRequestTTLTimeSelect(payloadStr string) (*RequestTTLTimeSelect, error) {
	var payload RequestTTLTimeSelectPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload format: %s", payloadStr)
	}

	strPrivateMetadata := payload.View.PrivateMetadata
	var privateMetadata RequestTTLTimeSelectPrivateMetadataPayload
	if err := json.Unmarshal([]byte(strPrivateMetadata), &privateMetadata); err != nil {
		return nil, fmt.Errorf("invalid private metadata format: %s", strPrivateMetadata)
	}

	return &RequestTTLTimeSelect{
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
		SelectedAccessDurationDate:       privateMetadata.SelectedAccessDurationDate,
		SelectedAccessDurationTime:       privateMetadata.SelectedAccessDurationTime,
		RequestTTL:                       privateMetadata.RequestTTL,
		SelectedRequestTTLOptionID:       privateMetadata.SelectedRequestTTLOptionID,
		SelectedRequestTTLOptionName:     privateMetadata.SelectedRequestTTLOptionName,
		SelectedRequestTTLDate:           privateMetadata.SelectedRequestTTLDate,
		RequesterID:                      payload.User.ID,
		RequesterName:                    payload.User.Name,
		TriggerID:                        payload.TriggerID,
		ViewHash:                         payload.View.Hash,
		ViewID:                           payload.View.ID,
		RequestTTLTime:                   payload.View.State.Values.RequestTTLDateTimeBlock.AccessRequestRequestTTLTimeSelect.SelectedTime,
	}, nil
}
