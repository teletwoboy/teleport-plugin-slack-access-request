package accessrequest

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
)

type JudgementPayload struct {
	RequesterID        string
	RequesterChannelID string
	SelectedChannelID  string
	SlackUserID        int32
	UserID             int32
	Username           string
}

func NewOutboxWithJudgement(
	ob *model.Outbox,
	p *viewsubmission.AccessRequestModal,
	slackUserID, userID int32,
	username string,
) (*model.Outbox, error) {
	payload := JudgementPayload{
		RequesterID:        p.RequesterID,
		RequesterChannelID: p.RequesterChannelID,
		SelectedChannelID:  p.SelectedChannelID,
		SlackUserID:        slackUserID,
		UserID:             userID,
		Username:           username,
	}
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access request auto review judgement payload: %w", err)
	}

	outbox := &model.Outbox{
		EventType:   constant.AccessRequestJudgement,
		AggregateID: ob.AggregateID,
		Payload:     string(marshaled),
		Status:      constant.Pending,
	}
	return outbox, nil
}
