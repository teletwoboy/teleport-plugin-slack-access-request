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
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/policy/models"
	slackmodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/teleport/models"

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
	text := BuildAutoReviewToRequesterText(a.accessRequest, a.accessReview, a.requester)
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
	text := BuildAutoReviewToReviewersText(a.accessRequest, a.accessReview, a.requester, a.policy)
	return slack.MsgOptionText(text, false)
}
