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

package message

import (
	"fmt"

	"github.com/slack-go/slack"
)

type successInitSeed struct{}

func NewSuccessInitSeed() Builder {
	return &successInitSeed{}
}

func (s *successInitSeed) Build() slack.MsgOption {
	text := "*🤗 Hello, We are Teletwoboy!*\n"
	text += "\n```\n"
	text += "By extending Teleport OSS with this plugin, you gain:\n"
	text += "1️⃣ Time Efficiency – Faster, simpler workflows\n"
	text += "2️⃣ Cost Reduction – Save with open-source\n"
	text += "3️⃣ Ease of Use – Slack-based, no extra learning curve\n"
	text += "4️⃣ Auditability – Clear records & compliance\n"
	text += "5️⃣ Automation – Auto approval/review flows\n"
	text += "6️⃣ Security – Strong access control\n"
	text += "7️⃣ Scalability – Easy integration & extension\n"
	text += "8️⃣ User-Friendly – Anyone can use it easily\n\n"
	text += "🔑 Features:\n"
	text += "/access-request → Slack-based request-review flow\n"
	text += "/access-policy  → ABAC-based auto-review policy\n\n"
	text += "⚠️ Note\n"
	text += "Time is shown in local timezone in modals, and in UTC in messages.\n"
	text += "This ensures readability for users and consistency in records.\n\n"
	text += "📌 GitHub: https://github.com/teletwoboy\n"
	text += "```\n"
	return slack.MsgOptionText(text, false)
}

type successCreateUser struct {
	realName string
	username string
}

func NewSuccessCreateUser(realName, username string) Builder {
	return &successCreateUser{realName: realName, username: username}
}

func (u *successCreateUser) Build() slack.MsgOption {
	text := "*🤗 Successfully Added User*\n"
	text += "\n```\n"
	text += fmt.Sprintf("Slack Name        : %s\n", u.realName)
	text += fmt.Sprintf("Teleport Username : %s\n", u.username)
	text += "```\n"
	return slack.MsgOptionText(text, false)
}

type successDeleteUser struct {
	realName string
	username string
}

func NewSuccessDeleteUser(realName, username string) Builder {
	return &successDeleteUser{realName: realName, username: username}
}

func (u *successDeleteUser) Build() slack.MsgOption {
	text := "*🤗 Successfully Deleted User*\n"
	text += "\n```\n"
	text += fmt.Sprintf("Deleted Slack Name        : %s\n", u.realName)
	text += fmt.Sprintf("Deleted Teleport Username : %s\n", u.username)
	text += "```\n"
	return slack.MsgOptionText(text, false)
}
