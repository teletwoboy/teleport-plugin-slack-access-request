package accessrequest

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
)

type ToReviewerPayload struct {
	ReviewerChannelID string
	SlackUserID       int32
}

func NewOutboxWithToReviewer(
	ob *model.Outbox,
	payload JudgementPayload,
) (*model.Outbox, error) {
	p := ToReviewerPayload{
		ReviewerChannelID: payload.SelectedChannelID,
		SlackUserID:       payload.SlackUserID,
	}
	marshaled, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access request to reviewer payload: %w", err)
	}

	outbox := &model.Outbox{
		EventType:   constant.AccessRequestToReviewer,
		AggregateID: ob.AggregateID,
		Payload:     string(marshaled),
		Status:      constant.Pending,
	}
	return outbox, nil
}
