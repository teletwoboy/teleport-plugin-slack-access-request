package container

import (
	"teleport-plugin-slack-access-request/internal/database/sqlc"
	"teleport-plugin-slack-access-request/internal/seedinit"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/user"
)

type Repositories struct {
	SeedInit *seedinit.PostgresRepository
	Slack    *slack.PostgresRepository
	Teleport *teleport.PostgresRepository
	User     *user.PostgresRepository
}

func NewRepositories(q sqlc.Querier) *Repositories {
	return &Repositories{
		SeedInit: seedinit.NewRepository(q),
		Slack:    slack.NewRepository(q),
		Teleport: teleport.NewRepository(q),
		User:     user.NewRepository(q),
	}
}
