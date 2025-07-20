package teleport

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"teleport-plugin-slack-access-request/internal/teleport/models"
	teleporttypes "teleport-plugin-slack-access-request/internal/teleport/types"

	"github.com/gravitational/teleport/api/types"
)

type Service interface {
	CreateAccessRequest(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error)
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	ExistsAccessRequestByName(ctx context.Context, name string) (bool, error)
	FetchAccessRequests(ctx context.Context, builder accessrequest.FilterBuilder) ([]types.AccessRequest, error)
	FetchUsersWithoutSecrets(ctx context.Context) ([]models.User, error)
	FetchUserAccessInfo(ctx context.Context, user models.User) (*teleporttypes.UserAccessInfo, error)
	GetAccessRequestByName(ctx context.Context, name string) (*models.AccessRequest, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	SubmitAccessRequest(ctx context.Context, builder accessrequest.CreateBuilder) (types.AccessRequest, error)
	SubmitAccessRequestState(ctx context.Context, builder accessrequest.UpdateBuilder) error
}

type API interface {
	CreateAccessRequestV2(ctx context.Context, req types.AccessRequest) (types.AccessRequest, error)
	GetAccessCapabilities(ctx context.Context, req types.AccessCapabilitiesRequest) (*types.AccessCapabilities, error)
	GetAccessRequests(ctx context.Context, filter types.AccessRequestFilter) ([]types.AccessRequest, error)
	GetUsers(ctx context.Context, withSecrets bool) ([]types.User, error)
	SetAccessRequestState(ctx context.Context, params types.AccessRequestUpdate) error
}

type Repository interface {
	CreateAccessRequest(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error)
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	ExistsAccessRequestByName(ctx context.Context, name string) (bool, error)
	GetAccessRequestByName(ctx context.Context, name string) (*models.AccessRequest, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type service struct {
	api  API
	repo Repository
}

func NewService(api API, repo Repository) Service {
	return &service{api: api, repo: repo}
}

func (s *service) CreateAccessRequest(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error) {
	return s.repo.CreateAccessRequest(ctx, accessRequest)
}

func (s *service) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed tp create Teleport user: %w", err)
	}
	return createdUser, nil
}

func (s *service) ExistsAccessRequestByName(ctx context.Context, name string) (bool, error) {
	exists, err := s.repo.ExistsAccessRequestByName(ctx, name)
	if err != nil {
		return false, fmt.Errorf("failed to check if access request exists: %w", err)
	}
	return exists, nil
}

func (s *service) FetchAccessRequests(ctx context.Context, builder accessrequest.FilterBuilder) ([]types.AccessRequest, error) {
	accessRequestFilter := builder.Build()
	return s.api.GetAccessRequests(ctx, accessRequestFilter)
}

func (s *service) FetchUsersWithoutSecrets(ctx context.Context) ([]models.User, error) {
	rawUsers, err := s.api.GetUsers(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	humanUsers := filterHumanUsers(rawUsers)
	return convertToUsers(humanUsers), nil
}

func (s *service) FetchUserAccessInfo(ctx context.Context, user models.User) (*teleporttypes.UserAccessInfo, error) {
	req := types.AccessCapabilitiesRequest{
		User:             user.Username,
		RequestableRoles: true,
	}

	resp, err := s.api.GetAccessCapabilities(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get access capabilities: %w", err)
	}

	return &teleporttypes.UserAccessInfo{
		Roles:         resp.RequestableRoles,
		RequireReason: resp.RequireReason,
	}, nil
}

func (s *service) GetAccessRequestByName(ctx context.Context, name string) (*models.AccessRequest, error) {
	accessRequest, err := s.repo.GetAccessRequestByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get access request status: %w", err)
	}
	return accessRequest, nil
}

func (s *service) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	u, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get teleport user by username: %w", err)
	}
	return u, nil
}

func (s *service) SubmitAccessRequestState(ctx context.Context, builder accessrequest.UpdateBuilder) error {
	accessRequestState := builder.Build()
	return s.api.SetAccessRequestState(ctx, accessRequestState)
}

func (s *service) SubmitAccessRequest(ctx context.Context, builder accessrequest.CreateBuilder) (types.AccessRequest, error) {
	accessRequest := builder.Build()
	return s.api.CreateAccessRequestV2(ctx, accessRequest)
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

func convertToUsers(users []types.User) []models.User {
	var result []models.User
	for _, user := range users {
		copiedUser := user
		result = append(result, models.User{
			Username: copiedUser.GetName(),
		})
	}
	return result
}
