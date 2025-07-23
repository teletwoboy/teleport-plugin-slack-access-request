package message

import (
	"fmt"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"

	"github.com/slack-go/slack"
)

type AccessReviewSubmissionBuilder struct {
	accessRequest *teleportmodels.AccessRequest
	accessReview  *teleportmodels.AccessReview
	requester     *slackmodels.User
	reviewer      *slackmodels.User
}

func NewAccessReviewSubmissionBuilder(
	accessRequest *teleportmodels.AccessRequest,
	accessReview *teleportmodels.AccessReview,
	requester *slackmodels.User,
	reviewer *slackmodels.User,
) Builder {
	return &AccessReviewSubmissionBuilder{
		accessRequest: accessRequest,
		accessReview:  accessReview,
		requester:     requester,
		reviewer:      reviewer,
	}
}

func (a *AccessReviewSubmissionBuilder) Build() slack.MsgOption {
	text := "*🔐 Access request review completed*\n"
	text += "\n```\n"
	text += fmt.Sprintf("📝 Access Request UUID : %s\n", a.accessRequest.Name)
	text += "\n"
	text += fmt.Sprintf("📝 State              : %s\n", a.accessRequest.State)
	text += fmt.Sprintf("👤 Reviewer           : %s\n", a.reviewer.RealName)
	text += fmt.Sprintf("📝 Review Reason      : %s\n", a.accessReview.Reason)
	text += fmt.Sprintf("👤 Requester          : %s\n", a.requester.RealName)
	text += fmt.Sprintf("💬 Requester Channel  : #%s\n", a.accessRequest.InputChannelName)
	text += "```\n"
	return slack.MsgOptionText(text, false)
}

// -- To requester

type AccessReviewToRequestorBuilder struct {
	accessRequest *teleportmodels.AccessRequest
	accessReview  *teleportmodels.AccessReview
	requester     *slackmodels.User
	reviewer      *slackmodels.User
}

func NewAccessReviewToRequestorBuilder(
	accessRequest *teleportmodels.AccessRequest,
	accessReview *teleportmodels.AccessReview,
	requester *slackmodels.User,
	reviewer *slackmodels.User,
) Builder {
	return &AccessReviewToRequestorBuilder{
		accessRequest: accessRequest,
		accessReview:  accessReview,
		requester:     requester,
		reviewer:      reviewer,
	}
}

func (a *AccessReviewToRequestorBuilder) Build() slack.MsgOption {
	text := "*🔐 Access request review completed*\n"
	text += "\n```\n"
	text += fmt.Sprintf("📝 State              : %s\n", a.accessRequest.State)
	text += fmt.Sprintf("✏️ Review Reason      : %s\n", a.accessReview.Reason)
	text += fmt.Sprintf("👤 Reviewer           : %s\n", a.reviewer.RealName)
	text += fmt.Sprintf("📡 Reviewers Channel  : %s\n", a.accessRequest.ReviewChannelName)
	text += fmt.Sprintf("👤 Requestor          : %s\n", a.requester.RealName)
	text += fmt.Sprintf("🎯 Request Role       : %s\n", a.accessRequest.Role)
	text += "```\n"
	return slack.MsgOptionText(text, false)
}
