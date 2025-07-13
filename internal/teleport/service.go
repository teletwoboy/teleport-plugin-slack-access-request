package teleport

import (
	"context"
	"fmt"

	"github.com/gravitational/teleport/api/types"
)

type API interface {
	GetUsers(ctx context.Context, withSecrets bool) ([]types.User, error)
	GetAccessCapabilities(ctx context.Context, req types.AccessCapabilitiesRequest) (*types.AccessCapabilities, error)
}

type Service struct {
	api API
}

func (s *Service) GetUsersWithoutSecrets(ctx context.Context) ([]User, error) {
	rawUsers, err := s.api.GetUsers(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	humanUsers := filterHumanUsers(rawUsers)
	return convertToUsers(humanUsers), nil
}

func (s *Service) GetAccessRequestableRoles(ctx context.Context, user User) ([]string, error) {
	req := types.AccessCapabilitiesRequest{
		User:             user.Username,
		RequestableRoles: true,
	}

	resp, err := s.api.GetAccessCapabilities(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get access capabilities: %w", err)
	}

	return resp.RequestableRoles, nil
}

func filterHumanUsers(users []types.User) []types.User {
	var humanUsers []types.User
	for _, user := range users {
		if !user.IsBot() {
			humanUsers = append(humanUsers, user)
		}
	}
	return humanUsers
}

func convertToUsers(users []types.User) []User {
	var result []User
	for _, user := range users {
		result = append(result, User{
			Username: user.GetName(),
		})
	}
	return result
}
