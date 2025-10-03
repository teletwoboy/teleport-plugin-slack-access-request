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

package telemetry

const (
	OpenModal              = "open-modal"
	OpenModalAccessRequest = "open-access-request-modal"
	OpenModalAccessReview  = "open-access-review-modal"
	OpenModalAccessPolicy  = "open-access-policy-modal"

	ARequest                              = "access-request"
	ARequestRoleSelection                 = "select-access-request-role"
	ARequestChannelSelection              = "select-access-request-channel"
	ARequestStartDateOptionSelection      = "select-access-request-start-date-option"
	ARequestStartDateSelection            = "select-access-request-start-date"
	ARequestStartTimeSelection            = "select-access-request-start-time"
	ARequestAccessDurationOptionSelection = "select-access-request-access-duration-option"
	ARequestAccessDurationDateSelection   = "select-access-request-access-duration-date"
	ARequestAccessDurationTimeSelection   = "select-access-request-access-duration-time"
	ARequestRequestTTLOptionSelection     = "select-access-request-request-ttl-option"
	ARequestRequestTTLDateSelection       = "select-access-request-request-ttl-date"
	ARequestRequestTTLTimeSelection       = "select-access-request-request-ttl-time"
	ARequestModalSubmission               = "submit-access-request-modal"

	AReview                = "access-review"
	AReviewModalSubmission = "submit-access-review-modal"

	APolicy                   = "access-policy"
	APolicyChannelSelection   = "select-access-policy-channel"
	APolicyRoleSelection      = "select-access-policy-role"
	APolicyUserSelection      = "select-access-policy-user"
	APolicyStartDateSelection = "select-access-policy-start-date"
	APolicyStartTimeSelection = "select-access-policy-start-time"
	APolicyEndDateSelection   = "select-access-policy-end-date"
	APolicyEndTimeSelection   = "select-access-policy-end-time"
	APolicyEffectSelection    = "select-access-policy-effect"
	APolicyModalSubmission    = "submit-access-policy-modal"

	SlackService    = "slack-service"
	TeleportService = "teleport-service"
	UserService     = "user-service"
	PolicyService   = "policy-service"
	WorkerService   = "worker-service"
	OutboxService   = "outbox-service"

	WorkerAccessRequest                      = "worker-access-request"
	WorkerAccessRequestSubmission            = "worker-access-request-submission"
	WorkerAccessRequestJudgement             = "worker-access-request-judgement"
	WorkerAccessRequestAutoReview            = "worker-access-request-auto-review"
	WorkerAccessRequestAutoReviewToRequester = "worker-access-request-auto-review-to-requester"
	WorkerAccessRequestAutoReviewToReviewer  = "worker-access-request-auto-review-to-reviewer"
	WorkerAccessRequestToRequester           = "worker-access-request-to-requester"
	WorkerAccessRequestToReviewer            = "worker-access-request-to-reviewer"

	WorkerAccessReview          = "worker-access-review"
	WorkerAccessReviewReviewer  = "worker-access-review-reviewer"
	WorkerAccessReviewRequester = "worker-access-review-requester"

	WorkerAccessPolicy         = "worker-access-policy"
	WorkerAccessPolicyCreation = "worker-access-policy-creation"
	WorkerAccessPolicyDeletion = "worker-access-policy-deletion"
)
