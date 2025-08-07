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

	"github.com/slack-go/slack"
)

type accessRequestSubmissionBuilder struct {
	accessRequest *teleportmodels.AccessRequest
	slackUser     *slackmodels.User
}

func NewAccessRequestSubmissionBuilder(a *teleportmodels.AccessRequest, s *slackmodels.User) Builder {
	return &accessRequestSubmissionBuilder{
		accessRequest: a,
		slackUser:     s,
	}
}

func (a *accessRequestSubmissionBuilder) Build() slack.MsgOption {
	text := "*🔐 Successfully submitted Access Request*\n"
	text += "\n```\n"
	text += fmt.Sprintf("👤 Requester          : %s\n", a.slackUser.RealName)
	text += fmt.Sprintf("🎯 Request Role       : %s\n", a.accessRequest.Role)
	text += fmt.Sprintf("📝 Request Reason     : %s\n", a.accessRequest.Reason)
	text += fmt.Sprintf("📡 Reviewers Channel  : #%s\n", a.accessRequest.ReviewChannelName)
	text += "\n"
	text += fmt.Sprintf("📅 Created At         : %s (UTC)", a.accessRequest.CreateDate.Format(util.SecondTimeFormat))
	text += "```\n"
	return slack.MsgOptionText(text, false)
}

// ------------------------------------------------------------------------
// -- To reviewers
type accessRequestToReviewersBuilder struct {
	accessRequest *teleportmodels.AccessRequest
	slackUser     *slackmodels.User
}

func NewAccessRequestToReviewersBuilder(a *teleportmodels.AccessRequest, s *slackmodels.User) Builder {
	return &accessRequestToReviewersBuilder{
		accessRequest: a,
		slackUser:     s,
	}
}

func (a *accessRequestToReviewersBuilder) Build() slack.MsgOption {
	text := "*🔐 Someone submitted Access Request*\n"
	text += "\n```\n"
	text += fmt.Sprintf("👤 Requester          : %s\n", a.slackUser.RealName)
	text += fmt.Sprintf("💬 Requester Channel  : #%s\n", a.accessRequest.InputChannelName)
	text += fmt.Sprintf("🎯 Request Role       : %s\n", a.accessRequest.Role)
	text += fmt.Sprintf("📝 Request Reason     : %s\n", a.accessRequest.Reason)
	text += fmt.Sprintf("📡 Reviewers Channel  : #%s\n", a.accessRequest.ReviewChannelName)
	text += fmt.Sprintf("⏳ Request Expiry     : %s (UTC)\n", a.accessRequest.Expires.Format(util.SecondTimeFormat))
	text += fmt.Sprintf("⏰ Role Expiry        : %s (UTC)\n", a.accessRequest.AccessDuration.Format(util.SecondTimeFormat))
	text += "\n"
	text += fmt.Sprintf("📅 Created At         : %s (UTC)", a.accessRequest.CreateDate.Format(util.SecondTimeFormat))
	text += "```"
	text += "\n👉 Click the button below to review this request."

	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", text, false, false),
			nil,
			nil,
		),
		slack.NewActionBlock(
			"access_request_actions",
			slack.NewButtonBlockElement(
				"open_access_request_review_modal",
				a.accessRequest.Name,
				slack.NewTextBlockObject("plain_text", "Review Request", false, false),
			).WithStyle("primary"),
		),
	}

	return slack.MsgOptionBlocks(blocks...)
}
