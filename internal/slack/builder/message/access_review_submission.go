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
	text := builder.BuildAccessReviewSubmissionText(a.accessRequest, a.accessReview, a.requester, a.reviewer)
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
	text := builder.BuildAccessReviewToRequesterText(a.accessRequest, a.accessReview, a.requester, a.reviewer)
	return slack.MsgOptionText(text, false)
}
