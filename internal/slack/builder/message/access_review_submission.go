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
	slackmodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/teleport/models"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util"

	"github.com/slack-go/slack"
)

type accessReviewSubmissionBuilder struct{}

func NewAccessReviewSubmissionBuilder() Builder {
	return &accessReviewSubmissionBuilder{}
}

func (a *accessReviewSubmissionBuilder) Build() slack.MsgOption {
	text := BuildAccessReviewSubmissionText()
	return slack.MsgOptionText(text, false)
}

// ------------------------------------------------------------------------

type toReviewersUpdateBuilder struct {
	accessRequest *teleportmodels.AccessRequest
	requester     *slackmodels.User
	reviewer      *slackmodels.User
}

func NewToReviewersUpdateBuilder(a *teleportmodels.AccessRequest, requester, reviewer *slackmodels.User) Builder {
	return &toReviewersUpdateBuilder{
		accessRequest: a,
		requester:     requester,
		reviewer:      reviewer,
	}
}

func (t *toReviewersUpdateBuilder) Build() slack.MsgOption {
	text := BuildToReviewersUpdateText(t.accessRequest, t.requester, t.reviewer)
	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject(util.Markdown, text, false, false),
			nil,
			nil,
		),
	}
	return slack.MsgOptionBlocks(blocks...)
}

// -- To reviewers

type accessReviewToReviewersBuilder struct {
	accessRequest *teleportmodels.AccessRequest
	accessReview  *teleportmodels.AccessReview
	requester     *slackmodels.User
	reviewer      *slackmodels.User
	permalink     string
}

func NewAccessReviewToReviewersBuilder(
	accessRequest *teleportmodels.AccessRequest,
	accessReview *teleportmodels.AccessReview,
	requester *slackmodels.User,
	reviewer *slackmodels.User,
	permalink string,
) Builder {
	return &accessReviewToReviewersBuilder{
		accessRequest: accessRequest,
		accessReview:  accessReview,
		requester:     requester,
		reviewer:      reviewer,
		permalink:     permalink,
	}
}

func (a *accessReviewToReviewersBuilder) Build() slack.MsgOption {
	text := BuildAccessReviewToReviewersText(a.accessRequest, a.accessReview, a.requester, a.reviewer, a.permalink)
	return slack.MsgOptionText(text, false)
}

// -- To requester

type accessReviewToRequesterBuilder struct {
	accessRequest *teleportmodels.AccessRequest
	accessReview  *teleportmodels.AccessReview
	requester     *slackmodels.User
	reviewer      *slackmodels.User
}

func NewAccessReviewToRequesterBuilder(
	accessRequest *teleportmodels.AccessRequest,
	accessReview *teleportmodels.AccessReview,
	requester *slackmodels.User,
	reviewer *slackmodels.User,
) Builder {
	return &accessReviewToRequesterBuilder{
		accessRequest: accessRequest,
		accessReview:  accessReview,
		requester:     requester,
		reviewer:      reviewer,
	}
}

func (a *accessReviewToRequesterBuilder) Build() slack.MsgOption {
	text := BuildAccessReviewToRequesterText(a.accessRequest, a.accessReview, a.requester, a.reviewer)
	return slack.MsgOptionText(text, false)
}
