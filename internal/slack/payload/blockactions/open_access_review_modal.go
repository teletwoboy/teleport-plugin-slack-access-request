package blockactions

type OpenAccessReviewModalPayload struct {
	Type      string `json:"type"`
	TriggerID string `json:"trigger_id"`

	Channel struct {
		ID string `json:"id"`
	} `json:"channel"`

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
}
