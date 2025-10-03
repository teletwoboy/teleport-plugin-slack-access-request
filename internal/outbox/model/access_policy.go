package model

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	policymodels "teleport-plugin-slack-access-request/internal/policy/models"
)

type AccessPolicyCreationPayload struct {
	AccessPolicy      *policymodels.AccessPolicy
	RequesterRealName string
}

func NewOutboxWithAccessPolicyCreation(a *policymodels.AccessPolicy, r string) (*Outbox, error) {
	payload := AccessPolicyCreationPayload{
		AccessPolicy:      a,
		RequesterRealName: r,
	}
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access policy payload: %w", err)
	}

	outbox := &Outbox{
		EventType:   constant.AccessPolicyCreation,
		AggregateID: a.AccessPolicyID,
		Payload:     string(marshaled),
		Status:      constant.Pending,
	}
	return outbox, nil
}

type AccessPolicyDeletion struct {
	MessageTimestamp string
	InputChannelID   string
}

func NewOutboxWithAccessPolicyDeletion(a *policymodels.AccessPolicy) (*Outbox, error) {
	payload := AccessPolicyDeletion{
		MessageTimestamp: a.MessageTimestamp,
		InputChannelID:   a.InputChannelID,
	}
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access policy deletion payload: %w", err)
	}

	outbox := &Outbox{
		EventType:   constant.AccessPolicyDeletion,
		AggregateID: a.AccessPolicyID,
		Payload:     string(marshaled),
		Status:      constant.Pending,
	}
	return outbox, nil
}
