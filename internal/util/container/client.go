package container

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
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
