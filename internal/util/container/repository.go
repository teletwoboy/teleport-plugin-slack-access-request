package container

import (
	"teleport-plugin-slack-access-request/internal/database/sqlc"
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
}

func NewRepositories(q sqlc.Querier) *Repositories {
	return &Repositories{
		Policy:   policy.NewRepository(q),
		SeedInit: seedinit.NewRepository(q),
		Slack:    slack.NewRepository(q),
		Teleport: teleport.NewRepository(q),
		User:     user.NewRepository(q),
	}
}
