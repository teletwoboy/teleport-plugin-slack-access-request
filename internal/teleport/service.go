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
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"teleport-plugin-slack-access-request/internal/teleport/models"
	teleporttypes "teleport-plugin-slack-access-request/internal/teleport/types"

	"github.com/gravitational/teleport/api/client/userloginstate"
	"github.com/gravitational/teleport/api/types"
	userloginstatetype "github.com/gravitational/teleport/api/types/userloginstate"
)

type Service interface {
	Close() error
	CreateAccessRequest(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error)
	CreateAccessReview(ctx context.Context, accessReview *models.AccessReview) (*models.AccessReview, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUserLoginState(ctx context.Context, name string) error
	ExistsAccessRequestByName(ctx context.Context, name string) (bool, error)
	ExistsUserByUsername(ctx context.Context, username string) (bool, error)
	FetchAccessRequests(ctx context.Context, builder accessrequest.FilterBuilder) ([]types.AccessRequest, error)
	FetchRoles(ctx context.Context, users []models.User) (map[string]struct{}, error)
	FetchUsersWithoutSecrets(ctx context.Context) ([]models.User, error)
	FetchUserWithoutSecrets(ctx context.Context, user *models.User) (types.User, error)
	FetchUserAccessInfo(ctx context.Context, user *models.User) (*teleporttypes.UserAccessInfo, error)
	GetAccessRequestByName(ctx context.Context, name string) (*models.AccessRequest, error)
	GetAccessRequestStateByName(ctx context.Context, name string) (string, error)
	GetUserByTeleportUserID(ctx context.Context, id int32) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserLoginState(ctx context.Context, name string) (*userloginstatetype.UserLoginState, error)
	NewWatcher(ctx context.Context, watch types.Watch) (types.Watcher, error)
	SubmitAccessRequest(ctx context.Context, builder accessrequest.CreateBuilder) (types.AccessRequest, error)
	SubmitAccessRequestState(ctx context.Context, builder accessrequest.UpdateBuilder) error
	UpdateAccessRequestStateByName(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error)
}

type API interface {
	Close() error
	CreateAccessRequestV2(ctx context.Context, req types.AccessRequest) (types.AccessRequest, error)
	GetAccessCapabilities(ctx context.Context, req types.AccessCapabilitiesRequest) (*types.AccessCapabilities, error)
	GetAccessRequests(ctx context.Context, filter types.AccessRequestFilter) ([]types.AccessRequest, error)
	GetUser(ctx context.Context, name string, withSecrets bool) (types.User, error)
	GetUsers(ctx context.Context, withSecrets bool) ([]types.User, error)
	NewWatcher(ctx context.Context, watch types.Watch) (types.Watcher, error)
	SetAccessRequestState(ctx context.Context, params types.AccessRequestUpdate) error
	UserLoginStateClient() *userloginstate.Client
}

type Repository interface {
	CreateAccessRequest(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error)
	CreateAccessReview(ctx context.Context, accessReview *models.AccessReview) (*models.AccessReview, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, user *models.User) (*models.User, error)
	ExistsAccessRequestByName(ctx context.Context, name string) (bool, error)
	ExistsUserByUsername(ctx context.Context, username string) (bool, error)
	GetAccessRequestByName(ctx context.Context, name string) (*models.AccessRequest, error)
	GetAccessRequestStateByName(ctx context.Context, name string) (string, error)
	GetUserByTeleportUserID(ctx context.Context, id int32) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateAccessRequestStateByName(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error)
}

type service struct {
	api  API
	repo Repository
}

func NewService(api API, repo Repository) Service {
	return &service{api: api, repo: repo}
}

func (s *service) Close() error {
	return s.api.Close()
}

func (s *service) CreateAccessRequest(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error) {
	return s.repo.CreateAccessRequest(ctx, accessRequest)
}

func (s *service) CreateAccessReview(ctx context.Context, accessReview *models.AccessReview) (*models.AccessReview, error) {
	return s.repo.CreateAccessReview(ctx, accessReview)
}

func (s *service) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create Teleport user: %w", err)
	}
	return createdUser, nil
}

func (s *service) DeleteUser(ctx context.Context, user *models.User) (*models.User, error) {
	Deleted, err := s.repo.DeleteUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to delete Teleport user: %w", err)
	}
	return Deleted, nil
}

func (s *service) DeleteUserLoginState(ctx context.Context, name string) error {
	err := s.api.UserLoginStateClient().DeleteUserLoginState(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to delete state Teleport user: %w", err)
	}
	return nil
}

func (s *service) ExistsAccessRequestByName(ctx context.Context, name string) (bool, error) {
	exists, err := s.repo.ExistsAccessRequestByName(ctx, name)
	if err != nil {
		return false, fmt.Errorf("failed to check if access request exists: %w", err)
	}
	return exists, nil
}

func (s *service) ExistsUserByUsername(ctx context.Context, username string) (bool, error) {
	exists, err := s.repo.ExistsUserByUsername(ctx, username)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return exists, nil
}

func (s *service) FetchAccessRequests(ctx context.Context, builder accessrequest.FilterBuilder) ([]types.AccessRequest, error) {
	accessRequestFilter := builder.Build()
	return s.api.GetAccessRequests(ctx, accessRequestFilter)
}

func (s *service) FetchRoles(ctx context.Context, users []models.User) (map[string]struct{}, error) {
	roles := make(map[string]struct{})
	for _, u := range users {
		copiedUser := u
		user, err := s.api.GetUser(ctx, copiedUser.Username, false)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch role: %w", err)
		}
		for _, r := range user.GetRoles() {
			copiedRole := r
			roles[copiedRole] = struct{}{}
		}
	}
	return roles, nil
}

func (s *service) FetchUsersWithoutSecrets(ctx context.Context) ([]models.User, error) {
	rawUsers, err := s.api.GetUsers(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	humanUsers := filterHumanUsers(rawUsers)
	return convertToUsers(humanUsers), nil
}

func (s *service) FetchUserWithoutSecrets(ctx context.Context, user *models.User) (types.User, error) {
	rawUser, err := s.api.GetUser(ctx, user.Username, false)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	return rawUser, nil
}

func (s *service) FetchUserAccessInfo(ctx context.Context, user *models.User) (*teleporttypes.UserAccessInfo, error) {
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
		return nil, fmt.Errorf("failed to get access request state: %w", err)
	}
	return accessRequest, nil
}

func (s *service) GetAccessRequestStateByName(ctx context.Context, name string) (string, error) {
	accessRequest, err := s.repo.GetAccessRequestByName(ctx, name)
	if err != nil {
		return "", fmt.Errorf("failed to get access request state: %w", err)
	}
	return accessRequest.State, nil
}

func (s *service) GetUserByTeleportUserID(ctx context.Context, id int32) (*models.User, error) {
	u, err := s.repo.GetUserByTeleportUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get teleport user by telport user id: %w", err)
	}
	return u, nil
}

func (s *service) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	u, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get teleport user by username: %w", err)
	}
	return u, nil
}

func (s *service) GetUserLoginState(ctx context.Context, name string) (*userloginstatetype.UserLoginState, error) {
	state, err := s.api.UserLoginStateClient().GetUserLoginState(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get teleport user state by telport username: %w", err)
	}
	return state, nil
}

func (s *service) SubmitAccessRequestState(ctx context.Context, builder accessrequest.UpdateBuilder) error {
	accessRequestState := builder.Build()
	return s.api.SetAccessRequestState(ctx, accessRequestState)
}

func (s *service) SubmitAccessRequest(ctx context.Context, builder accessrequest.CreateBuilder) (types.AccessRequest, error) {
	accessRequest := builder.Build()
	return s.api.CreateAccessRequestV2(ctx, accessRequest)
}

func (s *service) UpdateAccessRequestStateByName(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error) {
	return s.repo.UpdateAccessRequestStateByName(ctx, accessRequest)
}

func (s *service) NewWatcher(ctx context.Context, watch types.Watch) (types.Watcher, error) {
	return s.api.NewWatcher(ctx, watch)
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
