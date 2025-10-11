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

package constant

const (
	AccessRequest                      = "access_request"
	AccessRequestSubmission            = "access_request_submission"
	AccessRequestJudgement             = "access_request_judgement"
	AccessRequestAutoReview            = "access_request_auto+review"
	AccessRequestAutoReviewToRequester = "access_request_auto_review_to_requester"
	AccessRequestAutoReviewToReviewer  = "access_request_auto_review_to_reviewer"
	AccessRequestToRequester           = "access_request_to_requester"
	AccessRequestToReviewer            = "access_request_to_reviewer"

	AccessReview          = "access_review"
	AccessReviewReviewer  = "access_review_reviewer"
	AccessReviewRequester = "access_review_requester"

	AccessPolicy         = "access_policy"
	AccessPolicyCreation = "access_policy_creation"
	AccessPolicyDeletion = "access_policy_deletion"
)
