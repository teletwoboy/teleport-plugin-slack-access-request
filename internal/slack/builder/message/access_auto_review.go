package message

import (
	"fmt"
	"teleport-plugin-slack-access-request/internal/policy/models"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type autoReviewToRequesterBuilder struct {
	accessRequest *teleportmodels.AccessRequest
	accessReview  *teleportmodels.AccessReview
	requester     *slackmodels.User
	policy        *models.AccessPolicy
}

func NewAutoReviewToRequesterBuilder(
	aRequest *teleportmodels.AccessRequest,
	aReview *teleportmodels.AccessReview,
	requester *slackmodels.User,
	policy *models.AccessPolicy,
) Builder {
	return &autoReviewToRequesterBuilder{
		accessRequest: aRequest,
		accessReview:  aReview,
		requester:     requester,
		policy:        policy,
	}
}

func (a *autoReviewToRequesterBuilder) Build() slack.MsgOption {
	text := "*🔐 Access request Auto Reviewed*\n"
	text += "\n```\n"
	text += fmt.Sprintf("📝 State              : %s\n", a.accessRequest.State)
	text += fmt.Sprintf("📝 Review Reason      : %s\n", a.accessReview.Reason)
	text += fmt.Sprintf("👤 Requester          : %s\n", a.requester.RealName)
	text += fmt.Sprintf("💬 Requester Channel  : #%s\n", a.accessRequest.InputChannelName)
	text += "\n"
	text += fmt.Sprintf("🧭 Start Date  : %s (UTC)\n", a.accessRequest.StartDate.Format(util.SecondTimeFormat))
	text += fmt.Sprintf("🧭 Role Expiry : %s (UTC) \n", a.accessRequest.AccessDuration.Format(util.SecondTimeFormat))
	text += "```\n"
	return slack.MsgOptionText(text, false)
}

type autoReviewToReviewersBuilder struct {
	accessRequest *teleportmodels.AccessRequest
	accessReview  *teleportmodels.AccessReview
	requester     *slackmodels.User
	policy        *models.AccessPolicy
}

func NewAutoReviewToReviewersBuilder(
	aRequest *teleportmodels.AccessRequest,
	aReview *teleportmodels.AccessReview,
	reviewer *slackmodels.User,
	policy *models.AccessPolicy,
) Builder {
	return &autoReviewToReviewersBuilder{
		accessRequest: aRequest,
		accessReview:  aReview,
		requester:     reviewer,
		policy:        policy,
	}
}

func (a *autoReviewToReviewersBuilder) Build() slack.MsgOption {
	text := "*🔐 Access request Auto Reviewed*\n"
	text += "\n```\n"
	text += fmt.Sprintf("📝 Access Request UUID : %s\n", a.accessRequest.Name)
	text += "\n"
	text += fmt.Sprintf("🏷️ Used Policy Title  : %s\n", a.policy.Title)
	text += fmt.Sprintf("⚡️ Used Policy Effect : %s\n", a.policy.Effect)
	text += fmt.Sprintf("📝 Request State      : %s\n", a.accessRequest.State)
	text += fmt.Sprintf("✏️ Review Reason      : %s\n", a.accessReview.Reason)
	text += fmt.Sprintf("👤 Requester          : %s\n", a.requester.RealName)
	text += fmt.Sprintf("💬 Requester Channel  : #%s\n", a.accessRequest.InputChannelName)
	text += "\n"
	text += fmt.Sprintf("🧭 Start Date  : %s (UTC) \n", a.accessRequest.StartDate.Format(util.SecondTimeFormat))
	text += fmt.Sprintf("🧭 Role Expiry : %s (UTC) \n", a.accessRequest.AccessDuration.Format(util.SecondTimeFormat))
	text += "```\n"
	return slack.MsgOptionText(text, false)
}
