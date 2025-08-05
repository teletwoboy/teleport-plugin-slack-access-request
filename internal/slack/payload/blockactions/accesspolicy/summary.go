package accesspolicy

import "time"

type SummaryPrivateMetadataPayload struct {
	ChannelID           string    `json:"channel_id"`
	ChannelName         string    `json:"channel_name"`
	RealName            string    `json:"real_name"`
	TimeZone            string    `json:"time_zone"`
	SelectedChannelID   string    `json:"selected_channel_id"`
	SelectedChannelName string    `json:"selected_channel_name"`
	SelectedRole        string    `json:"selected_role"`
	SelectedRoleName    string    `json:"selected_role_name"`
	SelectedUserID      string    `json:"selected_user_id"`
	SelectedRealName    string    `json:"selected_real_name"`
	SelectedStartDate   time.Time `json:"selected_start_date"`
	SelectedEndDate     time.Time `json:"selected_end_date"`
	SelectedEffect      string    `json:"selected_effect"`
}
