package teleport

import (
	"context"
	"fmt"

	"github.com/gravitational/teleport/api/types"
)

type Service interface {
	CreateUser(ctx context.Context, user User) (*User, error)
	FetchUsersWithoutSecrets(ctx context.Context) ([]User, error)
	FetchUserAccessInfo(ctx context.Context, user User) (*UserAccessInfo, error)
	GetUserByTeleportUserID(ctx context.Context, id int32) (*User, error)
}

type API interface {
	GetUsers(ctx context.Context, withSecrets bool) ([]types.User, error)
	GetAccessCapabilities(ctx context.Context, req types.AccessCapabilitiesRequest) (*types.AccessCapabilities, error)
}

type Repository interface {
	CreateUser(ctx context.Context, user User) (*User, error)
	GetUserByTeleportUserID(ctx context.Context, id int32) (*User, error)
}

type service struct {
	api  API
	repo Repository
}

func NewService(api API, repo Repository) Service {
	return &service{api: api, repo: repo}
}

func (s *service) CreateUser(ctx context.Context, user User) (*User, error) {
	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed tp create Teleport user: %w", err)
	}
	return createdUser, nil
}

func (s *service) FetchUsersWithoutSecrets(ctx context.Context) ([]User, error) {
	rawUsers, err := s.api.GetUsers(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	humanUsers := filterHumanUsers(rawUsers)
	return convertToUsers(humanUsers), nil
}

func (s *service) FetchUserAccessInfo(ctx context.Context, user User) (*UserAccessInfo, error) {
	req := types.AccessCapabilitiesRequest{
		User:             user.Username,
		RequestableRoles: true,
	}

	resp, err := s.api.GetAccessCapabilities(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get access capabilities: %w", err)
	}

	return &UserAccessInfo{
		Roles:         resp.RequestableRoles,
		RequireReason: resp.RequireReason,
	}, nil
}

func (s *service) GetUserByTeleportUserID(ctx context.Context, id int32) (*User, error) {
	u, err := s.repo.GetUserByTeleportUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get teleport user by teleport userID: %w", err)
	}
	return u, nil
}

func filterHumanUsers(users []types.User) []types.User {
	var humanUsers []types.User
	for _, user := range users {
		copiedUser := user
		if !copiedUser.IsBot() {
			humanUsers = append(humanUsers, copiedUser)
		}
	}
	return humanUsers
}

func convertToUsers(users []types.User) []User {
	var result []User
	for _, user := range users {
		copiedUser := user
		result = append(result, User{
			Username: copiedUser.GetName(),
		})
	}
	return result
}
