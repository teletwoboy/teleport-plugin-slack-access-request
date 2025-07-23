package app

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/database/sqlc"
	"teleport-plugin-slack-access-request/internal/events"
	"teleport-plugin-slack-access-request/internal/seedinit"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/user"
)

type Clients struct {
	Slack    *slack.Client
	Teleport *teleport.Client
}

func NewClients(ctx context.Context) (*Clients, error) {
	slackClt, err := slack.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize slack client: %w", err)
	}

	teleportClt, err := teleport.Init(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize teleport client: %w", err)
	}

	return &Clients{
		Slack:    slackClt,
		Teleport: teleportClt,
	}, nil
}

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

type Services struct {
	Events   events.Service
	SeedInit seedinit.Service
	Slack    slack.Service
	Teleport teleport.Service
	User     user.Service
}

func NewServices(clients *Clients, repos *Repositories) *Services {
	slackSrv := slack.NewService(clients.Slack, repos.Slack)
	teleportSrv := teleport.NewService(clients.Teleport, repos.Teleport)
	userSrv := user.NewService(repos.User, slackSrv, teleportSrv)
	eventHandler := events.NewEventHandler(slackSrv, teleportSrv, userSrv)
	eventSrv := events.NewService(clients.Teleport, eventHandler)
	seedInitSrv := seedinit.NewService(repos.SeedInit, slackSrv, teleportSrv, userSrv)
	return &Services{
		Events:   eventSrv,
		SeedInit: seedInitSrv,
		Slack:    slackSrv,
		Teleport: teleportSrv,
		User:     userSrv,
	}
}
