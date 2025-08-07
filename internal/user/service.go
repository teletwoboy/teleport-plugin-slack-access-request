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

package user

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	"teleport-plugin-slack-access-request/internal/teleport"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	usermodels "teleport-plugin-slack-access-request/internal/user/models"
	"teleport-plugin-slack-access-request/internal/util"
)

type Service interface {
	CreateUser(ctx context.Context, user *usermodels.User) (*usermodels.User, error)
	DeleteUser(ctx context.Context, user *usermodels.User) (*usermodels.User, error)
	FetchUsers(ctx context.Context) ([]usermodels.User, error)
	MapUsersByUsername(slackUsers []slackmodels.User, teleportUsers []teleportmodels.User) []usermodels.User
	GetUserBySlackUserID(ctx context.Context, id int32) (*usermodels.User, error)
	GetUserByTeleportUserID(ctx context.Context, id int32) (*usermodels.User, error)
}

type Repository interface {
	CreateUser(ctx context.Context, user *usermodels.User) (*usermodels.User, error)
	DeleteUser(ctx context.Context, user *usermodels.User) (*usermodels.User, error)
	GetUserBySlackUserID(ctx context.Context, id int32) (*usermodels.User, error)
	GetUserByTeleportUserID(ctx context.Context, id int32) (*usermodels.User, error)
}

type service struct {
	repo        Repository
	slackSrv    slack.Service
	teleportSrv teleport.Service
}

func NewService(r Repository, s slack.Service, t teleport.Service) Service {
	return &service{
		repo:        r,
		slackSrv:    s,
		teleportSrv: t,
	}
}

func (s *service) CreateUser(ctx context.Context, user *usermodels.User) (*usermodels.User, error) {
	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return createdUser, nil
}

func (s *service) DeleteUser(ctx context.Context, user *usermodels.User) (*usermodels.User, error) {
	deletedUser, err := s.repo.DeleteUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}
	return deletedUser, nil
}

func (s *service) FetchUsers(ctx context.Context) ([]usermodels.User, error) {
	sUsers, err := s.slackSrv.FetchUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch slack users: %w", err)
	}

	tUsers, err := s.teleportSrv.FetchUsersWithoutSecrets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teleport users: %w", err)
	}
	return s.MapUsersByUsername(sUsers, tUsers), nil
}

func (s *service) GetUserBySlackUserID(ctx context.Context, id int32) (*usermodels.User, error) {
	user, err := s.repo.GetUserBySlackUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed tp get user by slack id: %w", err)
	}
	return user, nil
}

func (s *service) GetUserByTeleportUserID(ctx context.Context, id int32) (*usermodels.User, error) {
	user, err := s.repo.GetUserByTeleportUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed tp get user by slack id: %w", err)
	}
	return user, nil
}

func (s *service) MapUsersByUsername(slackUsers []slackmodels.User, teleportUsers []teleportmodels.User) []usermodels.User {
	var users []usermodels.User
	for _, teleportUser := range teleportUsers {
		for _, slackUser := range slackUsers {
			copiedTU := teleportUser
			copiedSU := slackUser
			if util.MatchesIdentifier(copiedTU.Username, copiedSU.Email) {
				users = append(users, usermodels.User{
					TeleportUser: &copiedTU,
					SlackUser:    &copiedSU,
				})
			}
		}
	}
	return users
}
