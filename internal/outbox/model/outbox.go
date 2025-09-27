package model

import (
	"encoding/json"
	"fmt"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"time"
)

const (
	AccessReview = "access-review"
)

type Outbox struct {
	OutboxID    int32
	EventType   string
	AggregateID int32
	Payload     string
	Status      string
	Attempts    int32
	ApiAttempts int32
	DBAttempts  int32
	NextTryAt   time.Time
	LastError   string
	UseYn       bool
	CreateCode  string
	CreateDate  time.Time
	UpdateCode  string
	UpdateDate  time.Time
	DeleteCode  string
	DeleteDate  time.Time
	Version     int64
}

type AccessReviewPayload struct {
	AccessRequest *teleportmodels.AccessRequest
	AccessReview  *teleportmodels.AccessReview
	Requester     *slackmodels.User
	Reviewer      *slackmodels.User
	MessageTs     string
}

func NewOutboxWithAccessReview(
	aRequest *teleportmodels.AccessRequest,
	aReview *teleportmodels.AccessReview,
	requester *slackmodels.User,
	reviewer *slackmodels.User,
	messageTs string,
) (*Outbox, error) {
	aReviewPayload := AccessReviewPayload{
		AccessRequest: aRequest,
		AccessReview:  aReview,
		Requester:     requester,
		Reviewer:      reviewer,
		MessageTs:     messageTs,
	}
	marshaled, err := json.Marshal(aReviewPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access review payload: %w", err)
	}

	outbox := &Outbox{
		EventType:   AccessReview,
		AggregateID: aReview.AccessReviewID,
		Payload:     string(marshaled),
		Status:      "pending",
	}
	return outbox, nil
}
