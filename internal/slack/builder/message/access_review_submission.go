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
	reviewer      *slackmodels.User
}

func NewAccessReviewSubmissionBuilder(
	accessRequest *teleportmodels.AccessRequest,
	accessReview *teleportmodels.AccessReview,
	reviewer *slackmodels.User,
) Builder {
	return &AccessReviewSubmissionBuilder{
		accessRequest: accessRequest,
		accessReview:  accessReview,
		reviewer:      reviewer,
	}
}

func (a *AccessReviewSubmissionBuilder) Build() slack.MsgOption {
	text := "*🔐 Access request review completed*\n"
	text += "\n```\n"
	text += fmt.Sprintf("📝 Access Request UUID : %s\n", a.accessRequest.Name)
	text += "\n"
	text += fmt.Sprintf("📝 State          : %s\n", a.accessRequest.State)
	text += fmt.Sprintf("👤 Reviewer       : %s\n", a.reviewer.RealName)
	text += fmt.Sprintf("📝 Review Reason  : %s\n", a.accessReview.Reason)
	text += "```\n"
	return slack.MsgOptionText(text, false)
}

// -- To requestor
type accessReviewToRequestorBuilder struct {
	State, Reviewer, ReviewChannelName, ReviewReason, Requestor, RequestRole string
}

func NewAccessReviewToRequestorBuilder(
	state, reviewer, reviewChannelName, reviewReason, requestor, requestRole string,
) Builder {
	return &accessReviewToRequestorBuilder{
		State:             state,
		Reviewer:          reviewer,
		ReviewChannelName: reviewChannelName,
		ReviewReason:      reviewReason,
		Requestor:         requestor,
		RequestRole:       requestRole,
	}
}

func (a *accessReviewToRequestorBuilder) Build() slack.MsgOption {
	text := "*🔐 Access request review completed*\n"
	text += "\n```\n"
	text += fmt.Sprintf("📝 State           : %s\n", a.State)
	text += fmt.Sprintf("👤 Reviewer        : %s\n", a.Reviewer)
	text += fmt.Sprintf("📡 Review Channel  : %s\n", a.ReviewChannelName)
	text += fmt.Sprintf("✏️ Review Reason   : %s\n", a.ReviewReason)
	text += fmt.Sprintf("👤 Requestor       : %s\n", a.Requestor)
	text += fmt.Sprintf("🎯 Request Role    : %s\n", a.RequestRole)
	text += "```\n"
	return slack.MsgOptionText(text, false)
}
