package modal

const (
	// access policy step section
	firstStepSection  = "*Step 1 of 6 - Target Channel*"
	secondStepSection = "*Step 2 of 6 - Target Role*"
	thirdStepSection  = "*Step 3 of 6 - Target User*"
	fourthStepSection = "*Step 4 of 6 - Duration*"
	fifthStepSection  = "*Step 5 of 6 - Effect*"
	sixthStepSection  = "*Step 6 of 6 - Summary*"

	// block
	plainText    = "plain_text"
	StaticSelect = "static_select"
	SelectOne    = "Select One"
	Markdown     = "mrkdwn"

	// access policy modal
	accessPolicyTitle          = "Access AutoReview Policy"
	accessPolicyCallBackID     = "access_policy_modal"
	accessPolicyAllOption      = "* (all)"
	accessPolicyAllOptionValue = "*"

	// access policy first step action block
	channelActionBlockID       = "channel_block"
	channelOptionBlockActionID = "access_policy_channel_select"
	// access policy second step action block
	roleActionBlockID       = "role_block"
	roleOptionBlockActionID = "access_policy_role_select"
)
