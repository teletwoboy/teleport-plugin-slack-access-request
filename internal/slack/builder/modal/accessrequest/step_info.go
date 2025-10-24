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

package accessrequest

import (
	"time"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accessrequest"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util"
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
	text += util.ARequestStartDateSecondInfo + ttl.Format(util.SlackTimeFormat)
	text += "\n```"
	return text
}

func BuildAccessDurationInfoText(ttl time.Time) string {
	text := "```\n"
	text += "1️⃣ Default – Ends at " + ttl.Format(util.SlackTimeFormat) + "\n"
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
	text += util.ARequestAccessDurationSecondInfo + ttl.Format(util.SlackTimeFormat)
	text += "\n```"
	return text
}

func BuildRequestTTLInfoText(ttl time.Time) string {
	text := "```\n"
	text += "1️⃣ Default – Expires at " + ttl.Format(util.SlackTimeFormat) + "\n"
	text += "\n"
	text += "2️⃣ Select DateTime – Set how long the request remains valid"
	text += "\n```"
	return text
}

func BuildRequestTTLSecondOptInfoText(requestTTLOpt string, ttl time.Time) string {
	text := "```\n"
	text += requestTTLOpt + "\n"
	text += "\n"
	text += util.ARequestRequestTTLSecondInfo + ttl.Format(util.SlackTimeFormat)
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
