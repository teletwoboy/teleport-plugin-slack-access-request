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

package teleport

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/config"
	"time"

	"github.com/gravitational/teleport/api/client"
	"github.com/gravitational/teleport/api/client/userloginstate"
	"github.com/gravitational/teleport/api/types"
	userloginstatetype "github.com/gravitational/teleport/api/types/userloginstate"
)

type Client struct {
	api API
}

func Init(ctx context.Context) (*Client, error) {
	addr := config.Cfg.Teleport.Addr
	identityPath := config.Cfg.Teleport.IdentityPath
	credentials := client.LoadIdentityFile(identityPath)
	cfg := client.Config{
		Addrs:       []string{addr},
		Credentials: []client.Credentials{credentials},
		DialTimeout: 5 * time.Second,
		Context:     ctx,
	}

	api, err := client.New(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create teleport client: %w", err)
	}
	_, err = api.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping teleport client: %w", err)
	}
	return &Client{api: api}, nil
}

func (c *Client) Close() error {
	return c.api.Close()
}

func (c *Client) CreateAccessRequestV2(ctx context.Context, req types.AccessRequest) (types.AccessRequest, error) {
	return c.api.CreateAccessRequestV2(ctx, req)
}

func (c *Client) UserLoginStateClient() *userloginstate.Client {
	return c.api.UserLoginStateClient()
}

func (c *Client) DeleteUserLoginState(ctx context.Context, name string) error {
	return c.UserLoginStateClient().DeleteUserLoginState(ctx, name)
}

func (c *Client) GetAccessCapabilities(ctx context.Context, req types.AccessCapabilitiesRequest) (*types.AccessCapabilities, error) {
	return c.api.GetAccessCapabilities(ctx, req)
}

func (c *Client) GetAccessRequests(ctx context.Context, filter types.AccessRequestFilter) ([]types.AccessRequest, error) {
	return c.api.GetAccessRequests(ctx, filter)
}

func (c *Client) GetRole(ctx context.Context, name string) (types.Role, error) {
	return c.api.GetRole(ctx, name)
}

func (c *Client) GetUser(ctx context.Context, name string, withSecrets bool) (types.User, error) {
	return c.api.GetUser(ctx, name, withSecrets)
}

func (c *Client) GetUsers(ctx context.Context, withSecrets bool) ([]types.User, error) {
	return c.api.GetUsers(ctx, withSecrets)
}

func (c *Client) GetUserLoginState(ctx context.Context, name string) (*userloginstatetype.UserLoginState, error) {
	return c.api.UserLoginStateClient().GetUserLoginState(ctx, name)
}

func (c *Client) SetAccessRequestState(ctx context.Context, params types.AccessRequestUpdate) error {
	return c.api.SetAccessRequestState(ctx, params)
}

func (c *Client) NewWatcher(ctx context.Context, watch types.Watch) (types.Watcher, error) {
	return c.api.NewWatcher(ctx, watch)
}
