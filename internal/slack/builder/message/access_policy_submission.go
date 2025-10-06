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
	policymodels "teleport-plugin-slack-access-request/internal/policy/models"

	"github.com/slack-go/slack"
)

type accessPolicySubmissionBuilder struct{}

func NewAccessPolicySubmissionBuilder() Builder {
	return &accessPolicySubmissionBuilder{}
}

func (a *accessPolicySubmissionBuilder) Build() slack.MsgOption {
	text := BuildAccessPolicySubmissionText()
	return slack.MsgOptionText(text, false)
}

// ------------------------------------------------------------------------

type accessPolicyToReviewersBuilder struct {
	accessPolicy      *policymodels.AccessPolicy
	requesterRealName string
}

func NewAccessPolicyToReviewersBuilder(a *policymodels.AccessPolicy, r string) Builder {
	return &accessPolicyToReviewersBuilder{
		accessPolicy:      a,
		requesterRealName: r,
	}
}

func (a *accessPolicyToReviewersBuilder) Build() slack.MsgOption {
	text := BuildAccessPolicyToReviewersText(a.accessPolicy, a.requesterRealName)
	return slack.MsgOptionText(text, false)
}
