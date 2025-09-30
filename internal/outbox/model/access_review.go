package model

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
)

type AccessReviewReviewerPayload struct {
	AccessRequest *teleportmodels.AccessRequest
	AccessReview  *teleportmodels.AccessReview
	Requester     *slackmodels.User
	Reviewer      *slackmodels.User
	MessageTs     string
}

func NewOutboxWithAccessReviewReviewer(
	aRequest *teleportmodels.AccessRequest,
	aReview *teleportmodels.AccessReview,
	requester *slackmodels.User,
	reviewer *slackmodels.User,
	messageTs string,
) (*Outbox, error) {
	payload := AccessReviewReviewerPayload{
		AccessRequest: aRequest,
		AccessReview:  aReview,
		Requester:     requester,
		Reviewer:      reviewer,
		MessageTs:     messageTs,
	}
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access review payload: %w", err)
	}

	outbox := &Outbox{
		EventType:   constant.AccessReviewReviewer,
		AggregateID: aReview.AccessReviewID,
		Payload:     string(marshaled),
		Status:      constant.Pending,
	}
	return outbox, nil
}

type AccessReviewRequesterPayload struct {
	AccessRequest *teleportmodels.AccessRequest
	AccessReview  *teleportmodels.AccessReview
	Requester     *slackmodels.User
	Reviewer      *slackmodels.User
}

func NewOutboxWithAccessReviewRequester(
	aRequest *teleportmodels.AccessRequest,
	aReview *teleportmodels.AccessReview,
	requester *slackmodels.User,
	reviewer *slackmodels.User,
) (*Outbox, error) {
	payload := AccessReviewRequesterPayload{
		AccessRequest: aRequest,
		AccessReview:  aReview,
		Requester:     requester,
		Reviewer:      reviewer,
	}
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access review payload: %w", err)
	}

	outbox := &Outbox{
		EventType:   constant.AccessReviewRequester,
		AggregateID: aReview.AccessReviewID,
		Payload:     string(marshaled),
		Status:      constant.Pending,
	}
	return outbox, nil
}
