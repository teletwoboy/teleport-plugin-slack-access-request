package accessrequest

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	policymodels "teleport-plugin-slack-access-request/internal/policy/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
)

type AutoReviewToRequesterPayload struct {
	AccessPolicyID  int32
	AccessRequestID int32
	AccessReviewID  int32
	SlackUserID     int32
}

func NewOutboxWithAutoReviewToRequester(
	accessPolicy *policymodels.AccessPolicy,
	accessRequest *teleportmodels.AccessRequest,
	accessReview *teleportmodels.AccessReview,
	slackUserID int32,
) (*model.Outbox, error) {
	p := AutoReviewToRequesterPayload{
		AccessPolicyID:  accessPolicy.AccessPolicyID,
		AccessRequestID: accessRequest.AccessRequestID,
		AccessReviewID:  accessReview.AccessReviewID,
		SlackUserID:     slackUserID,
	}
	marshaled, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access request creation payload: %w", err)
	}

	outbox := model.Outbox{
		EventType:     constant.AccessRequestAutoReviewToRequester,
		AggregateType: constant.AccessRequest,
		AggregateID:   accessRequest.AccessRequestID,
		Payload:       string(marshaled),
		Status:        constant.Pending,
	}
	return &outbox, nil
}
