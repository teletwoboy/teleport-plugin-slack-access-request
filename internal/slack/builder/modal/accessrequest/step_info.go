package accessrequest

import (
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
