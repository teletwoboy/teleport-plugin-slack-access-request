package accessrequest

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	policymodels "teleport-plugin-slack-access-request/internal/policy/models"
)

type AutoReviewPayload struct {
	AccessPolicyID     int32
	RequesterID        string
	RequesterChannelID string
	SelectedChannelID  string
	SlackUserID        int32
	UserID             int32
	Username           string
}

func NewOutboxWithAutoReview(
	policy *policymodels.AccessPolicy,
	ob *model.Outbox,
	payload JudgementPayload,
) (*model.Outbox, error) {
	p := AutoReviewPayload{
		AccessPolicyID:     policy.AccessPolicyID,
		RequesterID:        payload.RequesterID,
		RequesterChannelID: payload.RequesterChannelID,
		SelectedChannelID:  payload.SelectedChannelID,
		SlackUserID:        payload.SlackUserID,
		UserID:             payload.UserID,
		Username:           payload.Username,
	}
	marshaled, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access request auto review payload: %w", err)
	}

	outbox := &model.Outbox{
		EventType:     constant.AccessRequestAutoReview,
		AggregateType: constant.AccessRequest,
		AggregateID:   ob.AggregateID,
		Payload:       string(marshaled),
		Status:        constant.Pending,
	}
	return outbox, nil
}
