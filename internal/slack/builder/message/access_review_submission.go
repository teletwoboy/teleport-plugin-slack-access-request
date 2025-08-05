package message

import (
	"fmt"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type accessReviewSubmissionBuilder struct {
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
	return &accessReviewSubmissionBuilder{
		accessRequest: accessRequest,
		accessReview:  accessReview,
		requester:     requester,
		reviewer:      reviewer,
	}
}

func (a *accessReviewSubmissionBuilder) Build() slack.MsgOption {
	text := "*🔐 Access request review completed*\n"
	text += "\n```\n"
	text += fmt.Sprintf("📝 Access Request UUID : %s\n", a.accessRequest.Name)
	text += "\n"
	text += fmt.Sprintf("📝 State              : %s\n", a.accessRequest.State)
	text += fmt.Sprintf("👤 Reviewer           : %s\n", a.reviewer.RealName)
	text += fmt.Sprintf("📝 Review Reason      : %s\n", a.accessReview.Reason)
	text += fmt.Sprintf("👤 Requester          : %s\n", a.requester.RealName)
	text += fmt.Sprintf("💬 Requester Channel  : #%s\n", a.accessRequest.InputChannelName)
	text += "\n"
	text += fmt.Sprintf("🧭 Start Date  : %s (UTC)\n", a.accessRequest.StartDate.Format(util.SecondTimeFormat))
	text += fmt.Sprintf("🧭 Role Expiry : %s (UTC) \n", a.accessRequest.AccessDuration.Format(util.SecondTimeFormat))
	text += "```\n"
	return slack.MsgOptionText(text, false)
}

// -- To requester

type accessReviewToRequestorBuilder struct {
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
	return &accessReviewToRequestorBuilder{
		accessRequest: accessRequest,
		accessReview:  accessReview,
		requester:     requester,
		reviewer:      reviewer,
	}
}

func (a *accessReviewToRequestorBuilder) Build() slack.MsgOption {
	text := "*🔐 Access request review completed*\n"
	text += "\n```\n"
	text += fmt.Sprintf("📝 State              : %s\n", a.accessRequest.State)
	text += fmt.Sprintf("✏️ Review Reason      : %s\n", a.accessReview.Reason)
	text += fmt.Sprintf("👤 Reviewer           : %s\n", a.reviewer.RealName)
	text += fmt.Sprintf("📡 Reviewers Channel  : %s\n", a.accessRequest.ReviewChannelName)
	text += fmt.Sprintf("👤 Requestor          : %s\n", a.requester.RealName)
	text += fmt.Sprintf("🎯 Request Role       : %s\n", a.accessRequest.Role)
	text += "\n"
	text += fmt.Sprintf("🧭 Start Date  : %s (UTC)\n", a.accessRequest.StartDate.Format(util.SecondTimeFormat))
	text += fmt.Sprintf("🧭 Role Expiry : %s (UTC) \n", a.accessRequest.AccessDuration.Format(util.SecondTimeFormat))
	text += "```\n"
	return slack.MsgOptionText(text, false)
}
