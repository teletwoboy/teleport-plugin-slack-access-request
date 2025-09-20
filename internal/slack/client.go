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

func (c *Client) AddPin(channel string, item slack.ItemRef) error {
	return c.api.AddPin(channel, item)
}

func (c *Client) GetConversations(params *slack.GetConversationsParameters) ([]slack.Channel, string, error) {
	return c.api.GetConversations(params)
}

func (c *Client) GetPermalink(params *slack.PermalinkParameters) (string, error) {
	return c.api.GetPermalink(params)
}

func (c *Client) GetTeamInfo() (*slack.TeamInfo, error) {
	return c.api.GetTeamInfo()
}

func (c *Client) GetUserInfo(user string) (*slack.User, error) {
	return c.api.GetUserInfo(user)
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

func (c *Client) RemovePin(channel string, item slack.ItemRef) error {
	return c.api.RemovePin(channel, item)
}

func (c *Client) UpdateMessage(channelID, timestamp string, options ...slack.MsgOption) (string, string, string, error) {
	return c.api.UpdateMessage(channelID, timestamp, options...)
}

func (c *Client) UpdateView(view slack.ModalViewRequest, externalID, hash, viewID string) (*slack.ViewResponse, error) {
	return c.api.UpdateView(view, externalID, hash, viewID)
}
