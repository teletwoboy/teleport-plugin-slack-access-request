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

package util

import "time"

const (
	SlackTimeout = 3 * time.Second

	// MinuteTimeFormat is the time format used for displaying date and time as "YYYY-MM-DD HH:MM".
	MinuteTimeFormat = "2006-01-02 15:04"
	SlackTimeFormat  = "2006-01-02 15:04:05 -0700 MST"

	// PlainText is the Slack block type for plain text.
	PlainText = "plain_text"
	// StaticSelect is the Slack block type for static select menus.
	StaticSelect = "static_select"
	// SelectOne is the placeholder text for selecting one item.
	SelectOne = "Select One"
	// Markdown is the Slack block type for markdown-formatted text.
	Markdown = "mrkdwn"
	// Back is the label text for the back button in Slack modals.
	Back = "Back"
	// Submit is the label text for the submit button in Slack modals.
	Submit = "Submit"

	Close = "Close"

	// access request --------------------------------------------------------------

	ARequestFirstStepSection  = "*Step 6 of 1 - Select %s's Requestable Role*"
	ARequestSecondStepSection = "*Step 6 of 2 - Select Reviewers Channel*"
	ARequestThirdStepSection  = "*Step 6 of 3 - Start Date*"
	ARequestFourthStepSection = "*Step 6 of 4 - Access Duration*"
	ARequestFifthStepSection  = "*Step 6 of 5 - Request TTL*"
	ARequestSixthStepSection  = "*Step 6 of 6 - Summary*"

	ARequestTitle      = "Access Request"
	ARequestCallBackID = "access_request_modal"

	ARequestRoleActionBlockID       = "role_block"
	ARequestRoleOptionBlockActionID = "access_request_role_select"

	ARequestChannelActionBlockID       = "channel_block"
	ARequestChannelOptionBlockActionID = "access_request_channel_select"

	ARequestStartDateOptionActionBlockID       = "start_date_option_block"
	ARequestStartDateOptionOptionBlockActionID = "access_request_start_date_option_select"
	ARequestStartDateFirstOption               = "Immediately"
	ARequestStartDateSecondOption              = "Select DateTime"
	ARequestStartDateSecondInfo                = "💡 You can select a time before "
	ARequestStartDateTimeBlockID               = "start_date_time_block"
	ARequestStartDateBlockActionID             = "access_request_start_date_select"
	ARequestStartTimeBlockActionID             = "access_request_start_time_select"

	ARequestAccessDurationOptionActionBlockID = "access_duration_option_block"
	ARequestAccessDurationOptionBlockActionID = "access_request_access_duration_option_select"
	ARequestAccessDurationFirstOption         = "Default"
	ARequestAccessDurationSecondOption        = "Select DateTime"
	ARequestAccessDurationSecondInfo          = "💡 You can select a time before "
	ARequestAccessDurationDateTimeBlockID     = "access_duration_date_time_block"
	ARequestAccessDurationDateBlockActionID   = "access_request_access_duration_date_select"
	ARequestAccessDurationTimeBlockActionID   = "access_request_access_duration_time_select"

	ARequestRequestTTLOptionActionBlockID = "request_ttl_option_block"
	ARequestRequestTTLOptionBlockActionID = "access_request_request_ttl_option_select"
	ARequestRequestTTLFirstOption         = "Default"
	ARequestRequestTTLSecondOption        = "Select DateTime"
	ARequestRequestTTLSecondInfo          = "💡 You can select a time before "
	ARequestRequestTTLDateTimeBlockID     = "request_ttl_date_time_block"
	ARequestRequestTTLDateBlockActionID   = "access_request_request_ttl_date_select"
	ARequestRequestTTLTimeBlockActionID   = "access_request_request_ttl_time_select"

	ARequestReasonBlockID           = "reason_block"
	ARequestReasonBlockText         = "Request Reason"
	ARequestReasonElemBlockTest     = "Enter the reason"
	ARequestReasonElemBlockActionID = "access_request_reason_input"

	// access policy --------------------------------------------------------------

	// APolicyFirstStepSection is the section title for step 1 (Target Channel) in the access policy modal.
	APolicyFirstStepSection = "*Step 1 of 6 - Target Channel*"
	// APolicySecondStepSection is the section title for step 2 (Target Role) in the access policy modal.
	APolicySecondStepSection = "*Step 2 of 6 - Target Role*"
	// APolicyThirdStepSection is the section title for step 3 (Target User) in the access policy modal.
	APolicyThirdStepSection = "*Step 3 of 6 - Target User*"
	// APolicyFourthStepSection is the section title for step 4 (Duration) in the access policy modal.
	APolicyFourthStepSection = "*Step 4 of 6 - Duration*"

	APolicyFourthStepCautionSection = "```💡 Caution: The selected time will be converted to UTC```"
	// APolicyFourthStepFirstSubSection is the subsection title for selecting start date/time.
	APolicyFourthStepFirstSubSection = "4-1. Start Date/Time"
	// APolicyFourthStepSecondSubSection is the subsection title for selecting end date/time.
	APolicyFourthStepSecondSubSection = "4-2. End Date/Time"
	// APolicyFifthStepSection is the section title for step 5 (Effect) in the access policy modal.
	APolicyFifthStepSection = "*Step 5 of 6 - Effect*"
	// APolicySixthStepSection is the section title for step 6 (Summary) in the access policy modal.
	APolicySixthStepSection = "*Step 6 of 6 - Summary*"

	// APolicyTitle is the title text for the access policy Slack modal.
	APolicyTitle = "Access AutoReview Policy"
	// APolicyCallBackID is the callback ID for the access policy Slack modal.
	APolicyCallBackID = "access_policy_modal"
	// APolicyAllOption is the display label for the "all" option in dropdowns.
	APolicyAllOption = "* (all)"
	// APolicyAllOptionValue is the internal value for the "all" option in dropdowns.
	APolicyAllOptionValue = "*"

	// APolicyChanActionBlockID is the block ID for selecting a channel in step 1.
	APolicyChanActionBlockID = "channel_block"
	// APolicyChanOptionBlockActionID is the action ID for the channel select menu in step 1.
	APolicyChanOptionBlockActionID = "access_policy_channel_select"

	// APolicyRoleActionBlockID is the block ID for selecting a role in step 2.
	APolicyRoleActionBlockID = "role_block"
	// APolicyRoleOptionBlockActionID is the action ID for the role select menu in step 2.
	APolicyRoleOptionBlockActionID = "access_policy_role_select"

	// APolicyUserActionBlockID is the block ID for selecting a user in step 3.
	APolicyUserActionBlockID = "user_block"
	// APolicyUserOptionBlockActionID is the action ID for the user select menu in step 3.
	APolicyUserOptionBlockActionID = "access_policy_user_select"

	// APolicyStartDateTimeBlockID is the block ID for the start date and time in step 4.
	APolicyStartDateTimeBlockID = "start_date_time_block"
	// APolicyStartDateBlockActionID is the action ID for the start date picker.
	APolicyStartDateBlockActionID = "access_policy_start_date_select"
	// APolicyStartTimeBlockActionID is the action ID for the start time picker.
	APolicyStartTimeBlockActionID = "access_policy_start_time_select"
	// APolicyEndDateTimeBlockID is the block ID for the end date and time in step 4.
	APolicyEndDateTimeBlockID = "end_date_time_block"
	// APolicyEndDateBlockActionID is the action ID for the end date picker.
	APolicyEndDateBlockActionID = "access_policy_end_date_select"
	// APolicyEndTimeBlockActionID is the action ID for the end time picker.
	APolicyEndTimeBlockActionID = "access_policy_end_time_select"

	// APolicyEffectBlockID is the block ID for selecting the access effect (Allow/Deny) in step 5.
	APolicyEffectBlockID = "effect_block"
	// APolicyAllowButtonBlockActionID is the action ID for the "Allow" button.
	APolicyAllowButtonBlockActionID = "access_policy_allow_select"
	// APolicyAllowButtonValue is the value associated with the "Allow" button.
	APolicyAllowButtonValue = "allow"
	// APolicyAllowButtonText is the display text for the "Allow" button.
	APolicyAllowButtonText = "✅ Allow"
	// APolicyDenyButtonBlockActionID is the action ID for the "Deny" button.
	APolicyDenyButtonBlockActionID = "access_policy_deny_select"
	// APolicyDenyButtonValue is the value associated with the "Deny" button.
	APolicyDenyButtonValue = "deny"
	// APolicyDenyButtonText is the display text for the "Deny" button.
	APolicyDenyButtonText = "⛔ Deny"

	// APolicyTitleBlockID is the block ID for entering the policy title in step 6.
	APolicyTitleBlockID = "title_block"
	// APolicyTitleBlockText is the label text for the title input field.
	APolicyTitleBlockText = "Policy Title"
	// APolicyTitleElemBlockText is the placeholder text inside the title input field.
	APolicyTitleElemBlockText = "Enter the title"
	// APolicyTitleElemBlockActionID is the action ID for the title input field.
	APolicyTitleElemBlockActionID = "access_policy_title_input"
	// APolicyReasonBlockID is the block ID for entering the policy reason in step 6.
	APolicyReasonBlockID = "reason_block"
	// APolicyReasonBlockText is the label text for the reason input field.
	APolicyReasonBlockText = "Policy Reason"
	// APolicyReasonElemBlockText is the placeholder text inside the reason input field.
	APolicyReasonElemBlockText = "Enter the reason"
	// APolicyReasonElemBlockActionID is the action ID for the reason input field.
	APolicyReasonElemBlockActionID = "access_policy_reason_input"
)
