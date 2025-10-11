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

package container

import (
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/outbox"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/policy"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/seedinit"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/slack"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/teleport"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/user"
)

type Services struct {
	Policy   policy.Service
	SeedInit seedinit.Service
	Slack    slack.Service
	Teleport teleport.Service
	User     user.Service
	Outbox   outbox.Service
}

func NewServices(clients *Clients, repos *Repositories) *Services {
	policySrv := policy.NewService(repos.Policy)
	slackSrv := slack.NewService(clients.Slack, repos.Slack)
	teleportSrv := teleport.NewService(clients.Teleport, repos.Teleport)
	userSrv := user.NewService(repos.User, slackSrv, teleportSrv)
	seedInitSrv := seedinit.NewService(repos.SeedInit, slackSrv, teleportSrv, userSrv)
	outboxSrv := outbox.NewService(repos.Outbox)
	return &Services{
		Policy:   policySrv,
		SeedInit: seedInitSrv,
		Slack:    slackSrv,
		Teleport: teleportSrv,
		User:     userSrv,
		Outbox:   outboxSrv,
	}
}
