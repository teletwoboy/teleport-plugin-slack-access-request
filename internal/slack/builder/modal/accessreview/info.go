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

package accessreview

import (
	"fmt"

	slackmodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/teleport/models"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/util"
)

func BuildAccessReviewText(a *teleportmodels.AccessRequest, s *slackmodels.User) string {
	timezone := s.TimeZone
	text := "```\n"
	text += fmt.Sprintf("👤 Requester          : %s\n", s.RealName)
	text += fmt.Sprintf("💬 Requester Channel  : #%s\n", a.InputChannelName)
	text += fmt.Sprintf("🎯 Request Role       : %s\n", a.Role)
	text += fmt.Sprintf("📝 Request Reason     : %s\n", a.Reason)
	text += fmt.Sprintf("📡 Reviewers Channel  : #%s\n", a.ReviewChannelName)
	text += "\n"
	if a.StartDate.IsZero() {
		text += fmt.Sprintf("⏳ Start Date      : %s\n", util.ARequestStartDateFirstOption)
	} else {
		sD := util.ParseInLocation(a.StartDate, timezone)
		text += fmt.Sprintf("⏳ Start Date      : %s\n", sD.String())
	}
	aD := util.ParseInLocation(a.AccessDuration, timezone)
	rT := util.ParseInLocation(a.RequestTTL, timezone)
	cD := util.ParseInLocation(a.CreateDate, timezone)
	text += fmt.Sprintf("⏰ Access Duration : %s\n", aD.String())
	text += fmt.Sprintf("⏳ Request TTL     : %s\n", rT.String())
	text += "\n"
	text += fmt.Sprintf("📅 Created At      : %s\n", cD.String())
	text += "```"
	return text
}
