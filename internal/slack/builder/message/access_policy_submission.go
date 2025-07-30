package message

import (
	"fmt"
	policymodels "teleport-plugin-slack-access-request/internal/policy/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"

	"github.com/slack-go/slack"
)

type accessPolicySubmission struct {
	accessPolicy *policymodels.AccessPolicy
	payload      *viewsubmission.AccessPolicyModal
}

func NewAccessPolicySubmissionBuilder(a *policymodels.AccessPolicy, p *viewsubmission.AccessPolicyModal) Builder {
	return &accessPolicySubmission{
		accessPolicy: a,
		payload:      p,
	}
}

func (a *accessPolicySubmission) Build() slack.MsgOption {
	text := fmt.Sprintf(
		"```\n"+
			"🙋 Requester         : %s\n"+
			"💬 Requester Channel : #%s\n"+
			"\n"+
			"📥 Target Channel    : %s\n"+
			"🏷️ Target Role       : %s\n"+
			"👤 Target User       : %s\n"+
			"\n"+
			"🌍 Time Zone         : %s\n"+
			"🕐 Start Date        : %s\n"+
			"🕐 End Date          : %s\n"+
			"⚙️ Effect            : %s\n"+
			"\n"+
			"📅 Created At        : %s"+
			"\n```",
		a.payload.RequesterRealName,
		a.accessPolicy.InputChannelName,
		a.accessPolicy.TargetChannelName,
		a.accessPolicy.TargetRoleName,
		a.accessPolicy.TargetRealName,
		a.accessPolicy.TimeZone,
		a.payload.SelectedStartDate,
		a.payload.SelectedEndDate,
		a.accessPolicy.Effect,
		a.accessPolicy.CreateDate,
	)
	return slack.MsgOptionText(text, false)
}
