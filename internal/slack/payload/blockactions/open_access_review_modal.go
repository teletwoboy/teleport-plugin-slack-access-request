package blockactions

import (
	"encoding/json"
	"fmt"
)

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

type OpenAccessReviewModal struct {
	AccessRequestName string
	ReviewerChannelID string
	ReviewerID        string
	ReviewerName      string
	TriggerID         string
}

func ParseOpenAccessReviewModalPayload(payloadStr string) (*OpenAccessReviewModal, error) {
	var payload OpenAccessReviewModalPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return nil, fmt.Errorf("failed to parse open access review modal payload: %w", err)
	}
	return &OpenAccessReviewModal{
		AccessRequestName: payload.Actions[0].Value,
		ReviewerChannelID: payload.Channel.ID,
		ReviewerID:        payload.User.ID,
		ReviewerName:      payload.User.Name,
		TriggerID:         payload.TriggerID,
	}, nil
}
