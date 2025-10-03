package accessrequest

import (
	"encoding/json"
	"fmt"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	usermodels "teleport-plugin-slack-access-request/internal/user/models"
)

type SubmissionPayload struct {
	Payload     *viewsubmission.AccessRequestModal
	SlackUserID int32
	UserID      int32
	Username    string
}

func NewOutboxWithSubmission(
	p *viewsubmission.AccessRequestModal,
	slackUser *slackmodels.User,
	teleportUser *teleportmodels.User,
	user *usermodels.User,
	aRequestID int32,
) (*model.Outbox, error) {
	payload := SubmissionPayload{
		Payload:     p,
		SlackUserID: slackUser.SlackUserID,
		UserID:      user.UserID,
		Username:    teleportUser.Username,
	}
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access request creation payload: %w", err)
	}

	outbox := &model.Outbox{
		EventType:   constant.AccessRequestSubmission,
		AggregateID: aRequestID,
		Payload:     string(marshaled),
		Status:      constant.Pending,
	}
	return outbox, nil
}
