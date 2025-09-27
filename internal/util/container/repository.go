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
	"teleport-plugin-slack-access-request/internal/database/sqlc"
	"teleport-plugin-slack-access-request/internal/outbox"
	"teleport-plugin-slack-access-request/internal/policy"
	"teleport-plugin-slack-access-request/internal/seedinit"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/user"
)

type Repositories struct {
	Policy   *policy.PostgresRepository
	SeedInit *seedinit.PostgresRepository
	Slack    *slack.PostgresRepository
	Teleport *teleport.PostgresRepository
	User     *user.PostgresRepository
	Outbox   *outbox.PostgresRepository
}

func NewRepositories(q sqlc.Querier) *Repositories {
	return &Repositories{
		Policy:   policy.NewRepository(q),
		SeedInit: seedinit.NewRepository(q),
		Slack:    slack.NewRepository(q),
		Teleport: teleport.NewRepository(q),
		User:     user.NewRepository(q),
		Outbox:   outbox.NewRepository(q),
	}
}
