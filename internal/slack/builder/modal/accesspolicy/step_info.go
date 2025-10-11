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

package accesspolicy

import "github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack/payload/blockactions/accesspolicy"

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
