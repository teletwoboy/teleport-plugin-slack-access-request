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
	OpenModal              = "open_modal"
	OpenModalAccessRequest = "open_access_request_modal"
	OpenModalAccessReview  = "open_access_review_modal"
	OpenModalAccessPolicy  = "open_access_policy_modal"

	ARequest                              = "access_request"
	ARequestRoleSelection                 = "select_access_request_role"
	ARequestChannelSelection              = "select_access_request_channel"
	ARequestStartDateOptionSelection      = "select_access_request_start_date_option"
	ARequestStartDateSelection            = "select_access_request_start_date"
	ARequestStartTimeSelection            = "select_access_request_start_time"
	ARequestAccessDurationOptionSelection = "select_access_request_access_duration_option"
	ARequestAccessDurationDateSelection   = "select_access_request_access_duration_date"
	ARequestAccessDurationTimeSelection   = "select_access_request_access_duration_time"
	ARequestRequestTTLOptionSelection     = "select_access_request_request_ttl_option"
	ARequestRequestTTLDateSelection       = "select_access_request_request_ttl_date"
	ARequestRequestTTLTimeSelection       = "select_access_request_request_ttl_time"
	ARequestModalSubmission               = "submit_access_request_modal"

	AReview                = "access_review"
	AReviewModalSubmission = "submit_access_review_modal"

	APolicy                   = "access_policy"
	APolicyChannelSelection   = "select_access_policy_channel"
	APolicyRoleSelection      = "select_access_policy_role"
	APolicyUserSelection      = "select_access_policy_user"
	APolicyStartDateSelection = "select_access_policy_start_date"
	APolicyStartTimeSelection = "select_access_policy_start_time"
	APolicyEndDateSelection   = "select_access_policy_end_date"
	APolicyEndTimeSelection   = "select_access_policy_end_time"
	APolicyEffectSelection    = "select_access_policy_effect"
	APolicyModalSubmission    = "submit_access_policy_modal"

	SlackService    = "slack_service"
	TeleportService = "teleport_service"
	UserService     = "user_service"
	PolicyService   = "policy_service"
	WorkerService   = "worker_service"
	OutboxService   = "outbox_service"

	WorkerAccessRequest                      = "worker_access_request"
	WorkerAccessRequestSubmission            = "worker_access_request_submission"
	WorkerAccessRequestJudgement             = "worker_access_request_judgement"
	WorkerAccessRequestAutoReview            = "worker_access_request_auto_review"
	WorkerAccessRequestAutoReviewToRequester = "worker_access_request_auto_review_to_requester"
	WorkerAccessRequestAutoReviewToReviewer  = "worker_access_request_auto_review_to_reviewer"
	WorkerAccessRequestToRequester           = "worker_access_request_to_requester"
	WorkerAccessRequestToReviewer            = "worker_access_request_to_reviewer"

	WorkerAccessReview          = "worker_access_review"
	WorkerAccessReviewReviewer  = "worker_access_review_reviewer"
	WorkerAccessReviewRequester = "worker_access_review_requester"

	WorkerAccessPolicy         = "worker_access_policy"
	WorkerAccessPolicyCreation = "worker_access_policy_creation"
	WorkerAccessPolicyDeletion = "worker_access_policy_deletion"
)
