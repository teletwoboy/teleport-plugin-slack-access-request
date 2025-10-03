package accessrequest

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
)

type ToRequesterPayload struct {
	RequesterChannelID string
	SlackUserID        int32
}

func NewOutboxWithToRequester(
	ob *model.Outbox,
	payload JudgementPayload,
) (*model.Outbox, error) {
	p := ToRequesterPayload{
		RequesterChannelID: payload.RequesterChannelID,
		SlackUserID:        payload.SlackUserID,
	}
	marshaled, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access request to requester payload: %w", err)
	}

	outbox := &model.Outbox{
		EventType:     constant.AccessRequestToRequester,
		AggregateType: constant.AccessRequest,
		AggregateID:   ob.AggregateID,
		Payload:       string(marshaled),
		Status:        constant.Pending,
	}
	return outbox, nil
}
