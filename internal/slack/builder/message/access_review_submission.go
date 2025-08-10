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

package message

import (
	"fmt"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/gravitational/teleport/api/types"
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
	var text string
	if a.accessRequest.State == types.RequestState_APPROVED.String() {
		text = "*🔐 Access Request Review completed*\n"
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
	text = "*🔐 Access Request Review completed*\n"
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
	var text string
	if a.accessRequest.State == types.RequestState_APPROVED.String() {
		text += fmt.Sprintf("*🔐 %s's Access Request APPROVED ⭕️*\n", a.requester.RealName)
		text += "\n```\n"
		text += fmt.Sprintf("📝 Access Request UUID : %s\n", a.accessRequest.Name)
		text += "\n"
		text += fmt.Sprintf("📝 State              : %s\n", a.accessRequest.State)
		text += fmt.Sprintf("✏️ Review Reason      : %s\n", a.accessReview.Reason)
		text += fmt.Sprintf("👤 Reviewer           : %s\n", a.reviewer.RealName)
		text += fmt.Sprintf("📡 Reviewers Channel  : %s\n", a.accessRequest.ReviewChannelName)
		text += fmt.Sprintf("👤 Requestor          : %s\n", a.requester.RealName)
		text += fmt.Sprintf("🎯 Request Role       : %s\n", a.accessRequest.Role)
		text += "\n"
		text += fmt.Sprintf("🧭 Start Date  : %s (UTC)\n", a.accessRequest.StartDate.Format(util.SecondTimeFormat))
		text += fmt.Sprintf("🧭 Role Expiry : %s (UTC)\n", a.accessRequest.AccessDuration.Format(util.SecondTimeFormat))
		text += "\n"
		text += "// --------------------\n"
		text += "If you want to use the requested role, you must log in with an approved request\n"
		text += "\n"
		text += "// 1️⃣ If you are already logged in via CLI\n"
		text += "$ tsh login --request-id=<REQUEST_UUID>\n"
		text += "\n"
		text += "// 2️⃣ If you are not already logged in\n"
		text += "$ tsh login --proxy=<Teleport URL> --user=<Teleport Username> --request-id=<REQUEST_UUID>\n"
		text += "```\n"
		return slack.MsgOptionText(text, false)
	}

	text += fmt.Sprintf("*🔐 %s's Access Request DENIED ❌*\n", a.requester.RealName)
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
