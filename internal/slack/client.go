// Package slack provides business logic for interacting with Slack,
// such as sending and receiving messages, and retrieving channel lists.
package slack

import (
	"fmt"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/config"

	"github.com/slack-go/slack"
)

// API interface will later include methods like PostMessage from the Slack client
type API any

type Client struct {
	api API
}

// Init initializes and returns a new slack client using token
func Init() (*Client, error) {
	token := config.Cfg.Slack.Token
	api := slack.New(token)

	_, err := api.AuthTest()
	if err != nil {
		return nil, fmt.Errorf("failed to test slack auth: %w", err)
	}
	slog.Info("succeeded slack auth test")

	return &Client{api: api}, nil
}
