package teleport

import (
	"context"
	"fmt"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/config"
	"time"

	"github.com/gravitational/teleport/api/client"
)

// API interface will later include methods like GetUsers from the Teleport client
type API any

type Client struct {
	api API
}

func Init() (*Client, error) {
	ctx := context.Background()

	authAddr := config.Cfg.Teleport.AuthAddr
	identityPath := config.Cfg.Teleport.IdentityPath
	credentials := client.LoadIdentityFile(identityPath)

	cfg := client.Config{
		Addrs:       []string{authAddr},
		Credentials: []client.Credentials{credentials},
		DialTimeout: 5 * time.Second,
		Context:     ctx,
	}

	api, err := client.New(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create teleport client: %w", err)
	}
	slog.Info("successfully created teleport client")

	_, err = api.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping teleport client: %w", err)
	}
	slog.Info("successfully pinged to teleport server")

	return &Client{api: api}, nil
}
