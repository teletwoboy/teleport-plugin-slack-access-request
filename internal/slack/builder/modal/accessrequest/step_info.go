package accessrequest

import (
	"teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"teleport-plugin-slack-access-request/internal/util"
	"time"
)

func BuildStartDateInfoText() string {
	text := "```\n"
	text += "1️⃣ Immediately – Begins right after approval" + "\n"
	text += "\n"
	text += "2️⃣ Select DateTime – Set when access starts"
	text += "\n```"
	return text
}

func BuildStartDateFirstOptInfoText(startDateOpt string) string {
	text := "```\n"
	text += startDateOpt
	text += "\n```"
	return text
}

func BuildStartDateSecondOptInfoText(startDateOpt string, ttl time.Time) string {
	text := "```\n"
	text += startDateOpt + "\n"
	text += "\n"
	text += util.ARequestStartDateSecondInfo + ttl.String()
	text += "\n```"
	return text
}

func BuildAccessDurationInfoText(ttl time.Time) string {
	text := "```\n"
	text += "1️⃣ Default – Ends at " + ttl.String() + "\n"
	text += "\n"
	text += "2️⃣ Select DateTime – Set how long access lasts"
	text += "\n```"
	return text
}

func BuildAccessDurationFirstOptInfoText(accessDurationOpt string) string {
	text := "```\n"
	text += accessDurationOpt
	text += "\n```"
	return text
}

func BuildAccessDurationSecondOptInfoText(accessDurationOpt string, ttl time.Time) string {
	text := "```\n"
	text += accessDurationOpt + "\n"
	text += "\n"
	text += util.ARequestAccessDurationSecondInfo + ttl.String()
	text += "\n```"
	return text
}

func BuildRequestTTLInfoText(ttl time.Time) string {
	text := "```\n"
	text += "1️⃣ Default – Expires at " + ttl.String() + "\n"
	text += "\n"
	text += "2️⃣ Select DateTime – Set how long the request remains valid"
	text += "\n```"
	return text
}

func BuildRequestTTLSecondOptInfoText(requestTTLOpt string, ttl time.Time) string {
	text := "```\n"
	text += requestTTLOpt + "\n"
	text += "\n"
	text += util.ARequestRequestTTLSecondInfo + ttl.String()
	text += "\n```"
	return text
}

func BuildSummaryInfoText(p *accessrequest.RequestTTLTimeSelect) string {
	text := "```\n"
	text += "Requester         : " + p.RequesterRealName + "\n"
	text += "Requester Channel : " + p.RequesterChannelName + "\n"
	text += "\n"
	text += "Requested Role    : " + p.SelectedRole + "\n"
	text += "Reviewers Channel : " + p.SelectedChannelName + "\n"
	text += "\n"
	text += "Start Date Option      : " + p.SelectedStartDateOptionName + "\n"
	if p.SelectedStartDateOptionID == util.ARequestStartDateSecondOption {
		text += "Start Date - Date      : " + p.SelectedStartDate + "\n"
		text += "Start Date - Time      : " + p.SelectedStartTime + "\n"
		text += "\n"
	}
	text += "Access Duration Option : " + p.SelectedAccessDurationOptionName + "\n"
	if p.SelectedAccessDurationOptionID == util.ARequestAccessDurationSecondOption {
		text += "Access Duration - Date : " + p.SelectedAccessDurationDate + "\n"
		text += "Access Duration - Time : " + p.SelectedAccessDurationTime + "\n"
		text += "\n"
	}
	text += "Request TTL Option     : " + p.SelectedRequestTTLOptionName + "\n"
	if p.SelectedRequestTTLOptionID == util.ARequestRequestTTLSecondOption {
		text += "Request TTL - Date     : " + p.SelectedRequestTTLDate + "\n"
		text += "Request TTL - Time     : " + p.RequestTTLTime + "\n"
	}
	text += "```"
	return text
}
