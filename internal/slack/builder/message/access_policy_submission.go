package message

import (
	"fmt"
	"github.com/slack-go/slack"
	policymodels "teleport-plugin-slack-access-request/internal/policy/models"
	"teleport-plugin-slack-access-request/internal/slack/payload/viewsubmission"
	"teleport-plugin-slack-access-request/internal/util"
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
			"🕐 Start Date        : %s (UTC)\n"+
			"🕐 End Date          : %s (UTC)\n"+
			"⚙️ Effect            : %s\n"+
			"\n"+
			"📅 Created At        : %s (UTC)"+
			"\n```",
		a.payload.RequesterRealName,
		a.accessPolicy.InputChannelName,
		a.accessPolicy.TargetChannelName,
		a.accessPolicy.TargetRoleName,
		a.accessPolicy.TargetRealName,
		a.payload.SelectedStartDate.Format(util.SecondTimeFormat),
		a.payload.SelectedEndDate.Format(util.SecondTimeFormat),
		a.accessPolicy.Effect,
		a.accessPolicy.CreateDate.Format(util.SecondTimeFormat),
	)
	return slack.MsgOptionText(text, false)
}
