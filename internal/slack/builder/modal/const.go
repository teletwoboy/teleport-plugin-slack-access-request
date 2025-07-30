package modal

const (
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

	// APFirstStepSection is the section title for step 1 (Target Channel) in the access policy modal.
	APFirstStepSection = "*Step 1 of 6 - Target Channel*"
	// APSecondStepSection is the section title for step 2 (Target Role) in the access policy modal.
	APSecondStepSection = "*Step 2 of 6 - Target Role*"
	// APThirdStepSection is the section title for step 3 (Target User) in the access policy modal.
	APThirdStepSection = "*Step 3 of 6 - Target User*"
	// APFourthStepSection is the section title for step 4 (Duration) in the access policy modal.
	APFourthStepSection = "*Step 4 of 6 - Duration*"
	// APFourthStepFirstSubSection is the subsection title for selecting time zone
	APFourthStepFirstSubSection = "4-1. Time Zone"
	// APFourthStepSecondSubSection is the subsection title for selecting start date/time.
	APFourthStepSecondSubSection = "4-2. Start Date/Time"
	// APFourthStepThirdSubSection is the subsection title for selecting end date/time.
	APFourthStepThirdSubSection = "4-3. End Date/Time"
	// APFifthStepSection is the section title for step 5 (Effect) in the access policy modal.
	APFifthStepSection = "*Step 5 of 6 - Effect*"
	// APSixthStepSection is the section title for step 6 (Summary) in the access policy modal.
	APSixthStepSection = "*Step 6 of 6 - Summary*"

	// APTitle is the title text for the access policy Slack modal.
	APTitle = "Access AutoReview Policy"
	// APCallBackID is the callback ID for the access policy Slack modal.
	APCallBackID = "access_policy_modal"
	// APAllOption is the display label for the "all" option in dropdowns.
	APAllOption = "* (all)"
	// APAllOptionValue is the internal value for the "all" option in dropdowns.
	APAllOptionValue = "*"

	// APChannelActionBlockID is the block ID for selecting a channel in step 1.
	APChannelActionBlockID = "channel_block"
	// APChannelOptionBlockActionID is the action ID for the channel select menu in step 1.
	APChannelOptionBlockActionID = "access_policy_channel_select"

	// APRoleActionBlockID is the block ID for selecting a role in step 2.
	APRoleActionBlockID = "role_block"
	// APRoleOptionBlockActionID is the action ID for the role select menu in step 2.
	APRoleOptionBlockActionID = "access_policy_role_select"

	// APUserActionBlockID is the block ID for selecting a user in step 3.
	APUserActionBlockID = "user_block"
	// APUserOptionBlockActionID is the action ID for the user select menu in step 3.
	APUserOptionBlockActionID = "access_policy_user_select"

	// APTimeZoneActionBlockID is the block ID for selecting a time zone in step 4.
	APTimeZoneActionBlockID = "time_zone_block"
	// APTimeZoneOptionBlockActionID is the action ID for the time zone select menu in step 4.
	APTimeZoneOptionBlockActionID = "access_policy_time_zone_select"
	// APStartDateTimeBlockID is the block ID for the start date and time in step 4.
	APStartDateTimeBlockID = "start_date_time_block"
	// APStartDateBlockActionID is the action ID for the start date picker.
	APStartDateBlockActionID = "access_policy_start_date_select"
	// APStartTimeBlockActionID is the action ID for the start time picker.
	APStartTimeBlockActionID = "access_policy_start_time_select"
	// APEndDateTimeBlockID is the block ID for the end date and time in step 4.
	APEndDateTimeBlockID = "end_date_time_block"
	// APEndDateBlockActionID is the action ID for the end date picker.
	APEndDateBlockActionID = "access_policy_end_date_select"
	// APEndTimeBlockActionID is the action ID for the end time picker.
	APEndTimeBlockActionID = "access_policy_end_time_select"

	// APEffectBlockID is the block ID for selecting the access effect (Allow/Deny) in step 5.
	APEffectBlockID = "effect_block"
	// APAllowButtonBlockActionID is the action ID for the "Allow" button.
	APAllowButtonBlockActionID = "access_policy_allow_select"
	// APAllowButtonValue is the value associated with the "Allow" button.
	APAllowButtonValue = "allow"
	// APAllowButtonText is the display text for the "Allow" button.
	APAllowButtonText = "✅ Allow"
	// APDenyButtonBlockActionID is the action ID for the "Deny" button.
	APDenyButtonBlockActionID = "access_policy_deny_select"
	// APDenyButtonValue is the value associated with the "Deny" button.
	APDenyButtonValue = "deny"
	// APDenyButtonText is the display text for the "Deny" button.
	APDenyButtonText = "⛔ Deny"

	// APTitleBlockID is the block ID for entering the policy title in step 6.
	APTitleBlockID = "title_block"
	// APTitleBlockText is the label text for the title input field.
	APTitleBlockText = "Policy Title"
	// APTitleElemBlockText is the placeholder text inside the title input field.
	APTitleElemBlockText = "Enter the title"
	// APTitleElemBlockActionID is the action ID for the title input field.
	APTitleElemBlockActionID = "access_policy_title_input"
	// APReasonBlockID is the block ID for entering the policy reason in step 6.
	APReasonBlockID = "reason_block"
	// APReasonBlockText is the label text for the reason input field.
	APReasonBlockText = "Policy Reason"
	// APReasonElemBlockText is the placeholder text inside the reason input field.
	APReasonElemBlockText = "Enter the reason"
	// APReasonElemBlockActionID is the action ID for the reason input field.
	APReasonElemBlockActionID = "access_policy_reason_input"
)
