package teleport

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/config"
	"time"

	"github.com/gravitational/teleport/api/client"
	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/api/types/events"
)

type Client struct {
	api API
}

func Init(ctx context.Context) (*Client, error) {
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
	_, err = api.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping teleport client: %w", err)
	}
	return &Client{api: api}, nil
}

func (c *Client) CreateAccessRequestV2(ctx context.Context, req types.AccessRequest) (types.AccessRequest, error) {
	return c.api.CreateAccessRequestV2(ctx, req)
}

func (c *Client) GetAccessCapabilities(ctx context.Context, req types.AccessCapabilitiesRequest) (*types.AccessCapabilities, error) {
	return c.api.GetAccessCapabilities(ctx, req)
}

func (c *Client) GetAccessRequests(ctx context.Context, filter types.AccessRequestFilter) ([]types.AccessRequest, error) {
	return c.api.GetAccessRequests(ctx, filter)
}

func (c *Client) GetUsers(ctx context.Context, withSecrets bool) ([]types.User, error) {
	return c.api.GetUsers(ctx, withSecrets)
}

func (c *Client) SetAccessRequestState(ctx context.Context, params types.AccessRequestUpdate) error {
	return c.api.SetAccessRequestState(ctx, params)
}

func (c *Client) SearchEvents(ctx context.Context, fromUTC, toUTC time.Time, namespace string, eventTypes []string, limit int, order types.EventOrder, startKey string) ([]events.AuditEvent, string, error) {
	return c.api.SearchEvents(ctx, fromUTC, toUTC, namespace, eventTypes, limit, order, startKey)
}
