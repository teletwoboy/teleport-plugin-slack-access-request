/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package model

import (
	"encoding/json"
	"fmt"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox/constant"
	policymodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/policy/models"
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
		EventType:     constant.AccessPolicyCreation,
		AggregateType: constant.AccessPolicy,
		AggregateID:   a.AccessPolicyID,
		Payload:       string(marshaled),
		Status:        constant.Pending,
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
		EventType:     constant.AccessPolicyDeletion,
		AggregateType: constant.AccessPolicy,
		AggregateID:   a.AccessPolicyID,
		Payload:       string(marshaled),
		Status:        constant.Pending,
	}
	return outbox, nil
}
