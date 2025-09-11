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
	"github.com/slack-go/slack"
	"teleport-plugin-slack-access-request/internal/slack/builder"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util"
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
	text := builder.BuildAccessRequestSubmissionText(a.accessRequest, a.slackUser)
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
	text := builder.BuildAccessRequestToReviewersText(a.accessRequest, a.slackUser)
	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject(util.Markdown, text, false, false),
			nil,
			nil,
		),
		slack.NewActionBlock(
			"access_request_actions",
			slack.NewButtonBlockElement(
				"open_access_request_review_modal",
				a.accessRequest.Name,
				slack.NewTextBlockObject(util.PlainText, "Review Request", false, false),
			).WithStyle("primary"),
		),
	}

	return slack.MsgOptionBlocks(blocks...)
}
