package message

import (
	"fmt"
	policymodels "teleport-plugin-slack-access-request/internal/policy/models"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	"teleport-plugin-slack-access-request/internal/util"

	"github.com/gravitational/teleport/api/types"
)

func BuildAccessRequestSubmissionText(a *teleportmodels.AccessRequest, s *slackmodels.User) string {
	text := "*🔐 Successfully submitted Access Request*\n"
	text += "\n```\n"
	text += fmt.Sprintf("👤 Requester          : %s\n", s.RealName)
	text += fmt.Sprintf("🎯 Request Role       : %s\n", a.Role)
	text += fmt.Sprintf("📝 Request Reason     : %s\n", a.Reason)
	text += fmt.Sprintf("📡 Reviewers Channel  : #%s\n", a.ReviewChannelName)
	text += "\n"
	text += fmt.Sprintf("📅 Created At         : %s (UTC)", a.CreateDate.String())
	text += "```\n"
	return text
}

func BuildAccessRequestToReviewersText(a *teleportmodels.AccessRequest, s *slackmodels.User) string {
	text := "*🔐 Someone submitted Access Request*\n"
	text += "\n```\n"
	text += fmt.Sprintf("👤 Requester          : %s\n", s.RealName)
	text += fmt.Sprintf("💬 Requester Channel  : #%s\n", a.InputChannelName)
	text += fmt.Sprintf("🎯 Request Role       : %s\n", a.Role)
	text += fmt.Sprintf("📝 Request Reason     : %s\n", a.Reason)
	text += fmt.Sprintf("📡 Reviewers Channel  : #%s\n", a.ReviewChannelName)
	text += "\n"
	if a.StartDate.IsZero() {
		text += fmt.Sprintf("🧭 Start Date      : %s\n", util.ARequestStartDateFirstOption)
	} else {
		text += fmt.Sprintf("🧭 Start Date      : %s (UTC)\n", a.StartDate.String())
	}
	text += fmt.Sprintf("🧭 Access Duration : %s (UTC)\n", a.AccessDuration.String())
	text += fmt.Sprintf("🧭 Request TTL     : %s (UTC)\n", a.RequestTTL.String())
	text += "\n"
	text += fmt.Sprintf("📅 Created At      : %s (UTC)", a.CreateDate.String())
	text += "```"
	text += "\n👉 Click the button below to review this request."
	return text
}

func BuildToReviewersUpdateText(a *teleportmodels.AccessRequest, requester, reviewer *slackmodels.User) string {
	text := "*🔐 Someone submitted Access Request*\n"
	text += "\n```\n"
	text += fmt.Sprintf("👤 Requester          : %s\n", requester.RealName)
	text += fmt.Sprintf("💬 Requester Channel  : #%s\n", a.InputChannelName)
	text += fmt.Sprintf("🎯 Request Role       : %s\n", a.Role)
	text += fmt.Sprintf("📝 Request Reason     : %s\n", a.Reason)
	text += fmt.Sprintf("📡 Reviewers Channel  : #%s\n", a.ReviewChannelName)
	text += "\n"
	if a.StartDate.IsZero() {
		text += fmt.Sprintf("🧭 Start Date      : %s\n", util.ARequestStartDateFirstOption)
	} else {
		text += fmt.Sprintf("🧭 Start Date      : %s (UTC)\n", a.StartDate.String())
	}
	text += fmt.Sprintf("🧭 Access Duration : %s (UTC)\n", a.AccessDuration.String())
	text += fmt.Sprintf("🧭 Request TTL     : %s (UTC)\n", a.RequestTTL.String())
	text += "\n"
	text += fmt.Sprintf("📅 Created At      : %s (UTC)", a.CreateDate.String())
	text += "```"
	text += "\n"
	text += fmt.Sprintf("👉 *Reviewed by <@%s>*", reviewer.RealName)
	return text
}

func BuildAccessReviewSubmissionText(
	aRequest *teleportmodels.AccessRequest,
	aReview *teleportmodels.AccessReview,
	requester, reviewer *slackmodels.User,
	permalink string,
) string {
	var text string
	if aRequest.State == types.RequestState_APPROVED.String() {
		text = "*🔐 Access Request Review completed*\n"
		text += "\n```\n"
		text += fmt.Sprintf("📝 Access Request UUID : %s\n", aRequest.Name)
		text += "\n"
		text += fmt.Sprintf("📝 State              : %s\n", aRequest.State)
		text += fmt.Sprintf("👤 Reviewer           : %s\n", reviewer.RealName)
		text += fmt.Sprintf("📝 Review Reason      : %s\n", aReview.Reason)
		text += fmt.Sprintf("👤 Requester          : %s\n", requester.RealName)
		text += fmt.Sprintf("💬 Requester Channel  : #%s\n", aRequest.InputChannelName)
		text += fmt.Sprintf("🎯 Request Role       : %s\n", aRequest.Role)
		text += "\n"
		text += fmt.Sprintf("🧭 Start Date      : %s (UTC)\n", aRequest.StartDate.String())
		text += fmt.Sprintf("🧭 Access Duration : %s (UTC)\n", aRequest.AccessDuration.String())
		text += "\n"
		text += fmt.Sprintf("🔗 View Request : %s", permalink)
		text += "```\n"
		return text
	}
	text = "*🔐 Access Request Review completed*\n"
	text += "\n```\n"
	text += fmt.Sprintf("📝 Access Request UUID : %s\n", aRequest.Name)
	text += "\n"
	text += fmt.Sprintf("📝 State              : %s\n", aRequest.State)
	text += fmt.Sprintf("👤 Reviewer           : %s\n", reviewer.RealName)
	text += fmt.Sprintf("📝 Review Reason      : %s\n", aReview.Reason)
	text += fmt.Sprintf("👤 Requester          : %s\n", requester.RealName)
	text += fmt.Sprintf("💬 Requester Channel  : #%s\n", aRequest.InputChannelName)
	text += fmt.Sprintf("🎯 Request Role       : %s\n", aRequest.Role)
	text += "\n"
	text += fmt.Sprintf("🔗 View Request : %s", permalink)
	text += "```\n"
	return text
}

func BuildAccessReviewToRequesterText(
	aRequest *teleportmodels.AccessRequest,
	aReview *teleportmodels.AccessReview,
	requester, reviewer *slackmodels.User,
) string {
	var text string
	if aRequest.State == types.RequestState_APPROVED.String() {
		text += fmt.Sprintf("*🔐 %s's Access Request APPROVED ⭕️*\n", requester.RealName)
		text += "\n```\n"
		text += fmt.Sprintf("📝 Access Request UUID : %s\n", aRequest.Name)
		text += "\n"
		text += fmt.Sprintf("📝 State              : %s\n", aRequest.State)
		text += fmt.Sprintf("✏️ Review Reason      : %s\n", aReview.Reason)
		text += fmt.Sprintf("👤 Reviewer           : %s\n", reviewer.RealName)
		text += fmt.Sprintf("📡 Reviewers Channel  : %s\n", aRequest.ReviewChannelName)
		text += fmt.Sprintf("👤 Requestor          : %s\n", requester.RealName)
		text += fmt.Sprintf("🎯 Request Role       : %s\n", aRequest.Role)
		text += "\n"
		text += fmt.Sprintf("🧭 Start Date      : %s (UTC)\n", aRequest.StartDate.String())
		text += fmt.Sprintf("🧭 Access Duration : %s (UTC)\n", aRequest.AccessDuration.String())
		text += "\n"
		text += "// --------------------\n"
		text += "If you want to use the requested role, you must log in with an approved request\n"
		text += "\n"
		text += "// 1️⃣ If you are already logged in via CLI\n"
		text += "$ tsh login --request-id=<REQUEST_UUID>\n"
		text += "\n"
		text += "// 2️⃣ If you are not already logged in\n"
		text += "$ tsh login --proxy=<Teleport URL> --user=<Teleport Username> --request-id=<REQUEST_UUID>\n"
		text += "```\n"
		return text
	}
	text += fmt.Sprintf("*🔐 %s's Access Request DENIED ❌*\n", requester.RealName)
	text += "\n```\n"
	text += fmt.Sprintf("📝 State              : %s\n", aRequest.State)
	text += fmt.Sprintf("✏️ Review Reason      : %s\n", aReview.Reason)
	text += fmt.Sprintf("👤 Reviewer           : %s\n", reviewer.RealName)
	text += fmt.Sprintf("📡 Reviewers Channel  : %s\n", aRequest.ReviewChannelName)
	text += fmt.Sprintf("👤 Requestor          : %s\n", requester.RealName)
	text += fmt.Sprintf("🎯 Request Role       : %s\n", aRequest.Role)
	text += "```\n"
	return text
}

func BuildAutoReviewToRequesterText(
	aRequest *teleportmodels.AccessRequest,
	aReview *teleportmodels.AccessReview,
	requester *slackmodels.User,
) string {
	var text string
	if aRequest.State == types.RequestState_APPROVED.String() {
		text = fmt.Sprintf("*🔐 %s's Access Request APPROVED ⭕️*\n", requester.RealName)
		text += "\n```\n"
		text += fmt.Sprintf("📝 Access Request UUID : %s\n", aRequest.Name)
		text += "\n"
		text += fmt.Sprintf("📝 State              : %s\n", aRequest.State)
		text += fmt.Sprintf("📝 Review Reason      : %s\n", aReview.Reason)
		text += fmt.Sprintf("📡 Reviewers Channel  : %s\n", aRequest.ReviewChannelName)
		text += fmt.Sprintf("👤 Requester          : %s\n", requester.RealName)
		text += fmt.Sprintf("🎯 Request Role       : %s\n", aRequest.Role)
		text += "\n"
		text += fmt.Sprintf("🧭 Start Date      : %s (UTC)\n", aRequest.StartDate.String())
		text += fmt.Sprintf("🧭 Access Duration : %s (UTC)\n", aRequest.AccessDuration.String())
		text += "\n"
		text += "// --------------------\n"
		text += "If you want to use the requested role, you must log in with an approved request\n"
		text += "\n"
		text += "// 1️⃣ If you are already logged in via CLI\n"
		text += "$ tsh login --request-id=<REQUEST_UUID>\n"
		text += "\n"
		text += "// 2️⃣ If you are not already logged in\n"
		text += "$ tsh login --proxy=<Teleport URL> --user=<Teleport Username> --request-id=<REQUEST_UUID>\n"
		text += "```\n"
		return text
	}
	text = fmt.Sprintf("*🔐 %s's Access Request DENIED ❌*\n", requester.RealName)
	text += "\n```\n"
	text += fmt.Sprintf("📝 State              : %s\n", aRequest.State)
	text += fmt.Sprintf("📝 Review Reason      : %s\n", aReview.Reason)
	text += fmt.Sprintf("📡 Reviewers Channel  : %s\n", aRequest.ReviewChannelName)
	text += fmt.Sprintf("👤 Requester          : %s\n", requester.RealName)
	text += fmt.Sprintf("🎯 Request Role       : %s\n", aRequest.Role)
	text += "```\n"
	return text
}

func BuildAutoReviewToReviewersText(
	aRequest *teleportmodels.AccessRequest,
	aReview *teleportmodels.AccessReview,
	requester *slackmodels.User,
	policy *policymodels.AccessPolicy,
) string {
	var text string
	if aRequest.State == types.RequestState_APPROVED.String() {
		text = "*🔐 Access request Auto Reviewed*\n"
		text += "\n```\n"
		text += fmt.Sprintf("📝 Access Request UUID : %s\n", aRequest.Name)
		text += "\n"
		text += fmt.Sprintf("🏷️ Used Policy Title  : %s\n", policy.Title)
		text += fmt.Sprintf("⚡️ Used Policy Effect : %s\n", policy.Effect)
		text += fmt.Sprintf("📝 Request State      : %s\n", aRequest.State)
		text += fmt.Sprintf("✏️ Review Reason      : %s\n", aReview.Reason)
		text += fmt.Sprintf("👤 Requester          : %s\n", requester.RealName)
		text += fmt.Sprintf("💬 Requester Channel  : #%s\n", aRequest.InputChannelName)
		text += fmt.Sprintf("🎯 Request Role       : %s\n", aRequest.Role)
		text += "\n"
		text += fmt.Sprintf("🧭 Start Date      : %s (UTC) \n", aRequest.StartDate.String())
		text += fmt.Sprintf("🧭 Access Duration : %s (UTC) \n", aRequest.AccessDuration.String())
		text += "```\n"
		return text
	}
	text = "*🔐 Access request Auto Reviewed*\n"
	text += "\n```\n"
	text += fmt.Sprintf("📝 Access Request UUID : %s\n", aRequest.Name)
	text += "\n"
	text += fmt.Sprintf("🏷️ Used Policy Title  : %s\n", policy.Title)
	text += fmt.Sprintf("⚡️ Used Policy Effect : %s\n", policy.Effect)
	text += fmt.Sprintf("📝 State              : %s\n", aRequest.State)
	text += fmt.Sprintf("📝 Review Reason      : %s\n", aReview.Reason)
	text += fmt.Sprintf("👤 Requester          : %s\n", requester.RealName)
	text += fmt.Sprintf("💬 Requester Channel  : #%s\n", aRequest.InputChannelName)
	text += fmt.Sprintf("🎯 Request Role       : %s\n", aRequest.Role)
	text += "```\n"
	return text
}

func BuildAccessPolicySubmissionText(a *policymodels.AccessPolicy, p *viewsubmission.AccessPolicyModal) string {
	text := "```\n"
	text += fmt.Sprintf("🙋 Requester         : %s\n", p.RequesterRealName)
	text += fmt.Sprintf("💬 Requester Channel : #%s\n", a.InputChannelName)
	text += "\n"
	text += fmt.Sprintf("📥 Target Channel    : %s\n", a.TargetChannelName)
	text += fmt.Sprintf("🏷️ Target Role       : %s\n", a.TargetRoleName)
	text += fmt.Sprintf("👤 Target User       : %s\n", a.TargetRealName)
	text += "\n"
	text += fmt.Sprintf("🕐 Start Date        : %s (UTC)\n", a.StartDate.String())
	text += fmt.Sprintf("🕐 End Date          : %s (UTC)\n", a.EndDate.String())
	text += fmt.Sprintf("⚙️ Effect            : %s\n", a.Effect)
	text += "\n"
	text += fmt.Sprintf("📅 Created At        : %s (UTC)", a.CreateDate.String())
	text += "\n```"
	return text
}
