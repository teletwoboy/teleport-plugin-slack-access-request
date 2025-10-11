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
	"context"
	"fmt"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/config"

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

func (c *Client) AddPinContext(ctx context.Context, channel string, item slack.ItemRef) error {
	return c.api.AddPinContext(ctx, channel, item)
}

func (c *Client) GetConversationsContext(ctx context.Context, params *slack.GetConversationsParameters) (channels []slack.Channel, nextCursor string, err error) {
	return c.api.GetConversationsContext(ctx, params)
}

func (c *Client) GetPermalinkContext(ctx context.Context, params *slack.PermalinkParameters) (string, error) {
	return c.api.GetPermalinkContext(ctx, params)
}

func (c *Client) GetTeamInfoContext(ctx context.Context) (*slack.TeamInfo, error) {
	return c.api.GetTeamInfoContext(ctx)
}

func (c *Client) GetUserInfoContext(ctx context.Context, user string) (*slack.User, error) {
	return c.api.GetUserInfoContext(ctx, user)
}

func (c *Client) GetUsersContext(ctx context.Context, options ...slack.GetUsersOption) ([]slack.User, error) {
	return c.api.GetUsersContext(ctx, options...)
}

func (c *Client) GetUsersInConversationContext(ctx context.Context, params *slack.GetUsersInConversationParameters) ([]string, string, error) {
	return c.api.GetUsersInConversationContext(ctx, params)
}

func (c *Client) OpenViewContext(ctx context.Context, triggerID string, view slack.ModalViewRequest) (*slack.ViewResponse, error) {
	return c.api.OpenViewContext(ctx, triggerID, view)
}

func (c *Client) PostEphemeralContext(ctx context.Context, channelID, userID string, options ...slack.MsgOption) (timestamp string, err error) {
	return c.api.PostEphemeralContext(ctx, channelID, userID, options...)
}

func (c *Client) PostMessageContext(ctx context.Context, channel string, options ...slack.MsgOption) (string, string, error) {
	return c.api.PostMessageContext(ctx, channel, options...)
}

func (c *Client) PushViewContext(ctx context.Context, triggerID string, view slack.ModalViewRequest) (*slack.ViewResponse, error) {
	return c.api.PushViewContext(ctx, triggerID, view)
}

func (c *Client) RemovePinContext(ctx context.Context, channel string, item slack.ItemRef) error {
	return c.api.RemovePinContext(ctx, channel, item)
}

func (c *Client) UpdateMessageContext(ctx context.Context, channelID, timestamp string, options ...slack.MsgOption) (string, string, string, error) {
	return c.api.UpdateMessageContext(ctx, channelID, timestamp, options...)
}

func (c *Client) UpdateViewContext(ctx context.Context, view slack.ModalViewRequest, externalID, hash, viewID string) (*slack.ViewResponse, error) {
	return c.api.UpdateViewContext(ctx, view, externalID, hash, viewID)
}
