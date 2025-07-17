package user

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
)

type Service interface {
	CreateUser(ctx context.Context, user User) (*User, error)
	FetchUsers(ctx context.Context) ([]User, error)
	MapUsersByUsername(slackUsers []slack.User, teleportUsers []teleport.User) []User
	GetUserBySlackUserID(ctx context.Context, id int32) (*User, error)
}

type Repository interface {
	CreateUser(ctx context.Context, user User) (*User, error)
	GetUserBySlackUserID(ctx context.Context, id int32) (*User, error)
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

func (s *service) CreateUser(ctx context.Context, user User) (*User, error) {
	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed tp create user: %w", err)
	}
	return createdUser, nil
}

func (s *service) FetchUsers(ctx context.Context) ([]User, error) {
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

func (s *service) GetUserBySlackUserID(ctx context.Context, id int32) (*User, error) {
	user, err := s.repo.GetUserBySlackUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed tp get user by slack id: %w", err)
	}
	return user, nil
}

func (s *service) MapUsersByUsername(slackUsers []slack.User, teleportUsers []teleport.User) []User {
	var users []User
	for _, teleportUser := range teleportUsers {
		for _, slackUser := range slackUsers {
			copiedTU := teleportUser
			copiedSU := slackUser
			if copiedTU.Username == copiedSU.Email {
				users = append(users, User{
					TeleportUser: &copiedTU,
					SlackUser:    &copiedSU,
				})
			}
		}
	}

	return users
}
