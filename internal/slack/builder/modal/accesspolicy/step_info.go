package accesspolicy

import "teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"

func BuildSummaryInfoText(p *accesspolicy.EffectSelect) string {
	text := "```\n"
	text += "🙋 Requester         : " + p.RequesterRealName + "\n"
	text += "💬 Requester Channel : " + p.RequesterChannelName + "\n"
	text += "\n"
	text += "📥 Target Channel    : " + p.SelectedChannelName + "\n"
	text += "🏷️ Target Role       : " + p.SelectedRoleName + "\n"
	text += "👤 Target User       : " + p.SelectedRealName + "\n"
	text += "\n"
	text += "🕐 Start Date        : " + p.SelectedStartDate.String() + "\n"
	text += "🕐 End Date          : " + p.SelectedEndDate.String() + "\n"
	text += "⚙️ Effect            : " + p.Effect + "\n"
	text += "```"
	return text
}
