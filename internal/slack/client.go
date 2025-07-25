// Package slack provides business logic for interacting with Slack,
// such as sending and receiving messages, and retrieving channel lists.
package slack

import (
	"fmt"
	"teleport-plugin-slack-access-request/internal/config"

	"github.com/slack-go/slack"
)

type Client struct {
	api API
}

// Init initializes and returns a new Slack client using token
func Init() (*Client, error) {
	token := config.Cfg.Slack.Token
	api := slack.New(token)

	_, err := api.AuthTest()
	if err != nil {
		return nil, fmt.Errorf("failed to perform slack auth test: %w", err)
	}
	return &Client{api: api}, nil
}

func (c *Client) GetConversations(params *slack.GetConversationsParameters) ([]slack.Channel, string, error) {
	return c.api.GetConversations(params)
}

func (c *Client) GetTeamInfo() (*slack.TeamInfo, error) {
	return c.api.GetTeamInfo()
}

func (c *Client) GetUsers(options ...slack.GetUsersOption) ([]slack.User, error) {
	return c.api.GetUsers(options...)
}

func (c *Client) GetUsersInConversation(params *slack.GetUsersInConversationParameters) ([]string, string, error) {
	return c.api.GetUsersInConversation(params)
}

func (c *Client) OpenView(triggerID string, view slack.ModalViewRequest) (*slack.ViewResponse, error) {
	return c.api.OpenView(triggerID, view)
}

func (c *Client) PostMessage(channel string, options ...slack.MsgOption) (string, string, error) {
	return c.api.PostMessage(channel, options...)
}

func (c *Client) PushView(triggerID string, view slack.ModalViewRequest) (*slack.ViewResponse, error) {
	return c.api.PushView(triggerID, view)
}
