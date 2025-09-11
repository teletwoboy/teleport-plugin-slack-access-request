package accessrequest

import "time"

type SummaryPrivateMetadataPayload struct {
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
	SelectedRequestTTLTime           string    `json:"request_ttl_time"`
}
