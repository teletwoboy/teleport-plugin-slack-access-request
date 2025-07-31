package blockactions

type SummaryPrivateMetadataPayload struct {
	ChannelID           string `json:"channel_id"`
	ChannelName         string `json:"channel_name"`
	RealName            string `json:"real_name"`
	RequireReason       bool   `json:"require_reason"`
	SelectedRole        string `json:"selected_role"`
	SelectedChannelID   string `json:"selected_channel_id"`
	SelectedChannelName string `json:"selected_channel_name"`
}
