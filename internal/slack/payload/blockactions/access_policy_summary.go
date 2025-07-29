package blockactions

type AccessPolicySummaryPrivateMetadataPayload struct {
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
	SelectedEffect      string `json:"selected_effect"`
}
