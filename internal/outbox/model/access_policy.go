package model

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	policymodels "teleport-plugin-slack-access-request/internal/policy/models"
)

type AccessPolicyPayload struct {
	AccessPolicy      *policymodels.AccessPolicy
	RequesterRealName string
}

func NewOutboxWithAccessPolicy(a *policymodels.AccessPolicy, r string) (*Outbox, error) {
	payload := AccessPolicyPayload{
		AccessPolicy:      a,
		RequesterRealName: r,
	}
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access review payload: %w", err)
	}

	outbox := &Outbox{
		EventType:   constant.AccessPolicy,
		AggregateID: a.AccessPolicyID,
		Payload:     string(marshaled),
		Status:      constant.Pending,
	}
	return outbox, nil
}
