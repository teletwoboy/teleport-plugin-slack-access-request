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
	"teleport-plugin-slack-access-request/internal/policy/models"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/gravitational/teleport/api/types"
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
	var text string
	if a.accessRequest.State == types.RequestState_APPROVED.String() {
		text = fmt.Sprintf("*🔐 %s's Access Request APPROVED ⭕️*\n", a.requester.RealName)
		text += "\n```\n"
		text += fmt.Sprintf("📝 Access Request UUID : %s\n", a.accessRequest.Name)
		text += "\n"
		text += fmt.Sprintf("📝 State              : %s\n", a.accessRequest.State)
		text += fmt.Sprintf("📝 Review Reason      : %s\n", a.accessReview.Reason)
		text += fmt.Sprintf("📡 Reviewers Channel  : %s\n", a.accessRequest.ReviewChannelName)
		text += fmt.Sprintf("👤 Requester          : %s\n", a.requester.RealName)
		text += fmt.Sprintf("💬 Requester Channel  : #%s\n", a.accessRequest.InputChannelName)
		text += "\n"
		text += fmt.Sprintf("🧭 Start Date  : %s (UTC)\n", a.accessRequest.StartDate.Format(util.SecondTimeFormat))
		text += fmt.Sprintf("🧭 Role Expiry : %s (UTC) \n", a.accessRequest.AccessDuration.Format(util.SecondTimeFormat))
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
	text += fmt.Sprintf("📝 Review Reason      : %s\n", a.accessReview.Reason)
	text += fmt.Sprintf("📡 Reviewers Channel  : %s\n", a.accessRequest.ReviewChannelName)
	text += fmt.Sprintf("👤 Requester          : %s\n", a.requester.RealName)
	text += fmt.Sprintf("💬 Requester Channel  : #%s\n", a.accessRequest.InputChannelName)
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
