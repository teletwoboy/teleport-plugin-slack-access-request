package user

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
)

type Repository interface {
	CreateUser(ctx context.Context, user User) (*User, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, user User) (*User, error) {
	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed tp create user: %w", err)
	}
	return createdUser, nil
}

func (s *Service) MapUsersByUsername(slackUsers []slack.User, teleportUsers []teleport.User) []User {
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
