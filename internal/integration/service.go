package integration

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/database/sqlc"
	"teleport-plugin-slack-access-request/internal/seedinit"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/user"
)

type Services struct {
	SeedInitSrv *seedinit.Service
	SlackSrv    *slack.Service
	TeleportSrv *teleport.Service
	UserSrv     *user.Service
}

type Service struct {
	Services *Services
}

func NewServiceWithTx(qtx *sqlc.Queries, slackClt *slack.Client, teleportClt *teleport.Client) *Service {
	seedInitTxRepo := seedinit.NewRepositoryWithTx(qtx)
	slackTxRepo := slack.NewRepositoryWithTx(qtx)
	teleportTxRepo := teleport.NewRepositoryWithTx(qtx)
	userTxRepo := user.NewRepositoryWithTx(qtx)

	seedInitSrv := seedinit.NewService(seedInitTxRepo)
	slackSrv := slack.NewService(slackClt, slackTxRepo)
	teleportSrv := teleport.NewService(teleportClt, teleportTxRepo)
	userSrv := user.NewService(userTxRepo)
	return &Service{
		Services: &Services{
			SeedInitSrv: seedInitSrv,
			SlackSrv:    slackSrv,
			TeleportSrv: teleportSrv,
			UserSrv:     userSrv,
		},
	}
}

func (s *Service) FetchAndMapUsers(ctx context.Context) ([]user.User, error) {
	sUsers, err := s.Services.SlackSrv.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch slack users: %w", err)
	}

	tUsers, err := s.Services.TeleportSrv.GetUsersWithoutSecrets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teleport users: %w", err)
	}
	return s.Services.UserSrv.MapUsersByUsername(sUsers, tUsers), nil
}

func (s *Service) ProvisionUsers(ctx context.Context, users []user.User) error {
	for _, u := range users {
		copiedUser := u
		createdSlackUser, err := s.Services.SlackSrv.CreateUser(ctx, *copiedUser.SlackUser)
		if err != nil {
			return fmt.Errorf("failed to create slack user: %w", err)
		}

		createdTeleportUser, err := s.Services.TeleportSrv.CreateUser(ctx, *copiedUser.TeleportUser)
		if err != nil {
			return fmt.Errorf("failed to create teleport user: %w", err)
		}

		copiedUser.SlackUser = createdSlackUser
		copiedUser.TeleportUser = createdTeleportUser

		_, err = s.Services.UserSrv.CreateUser(ctx, copiedUser)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	err := s.Services.SeedInitSrv.UpdateStaus(ctx)
	if err != nil {
		return fmt.Errorf("failed to update seedinit status: %w", err)
	}
	return nil
}
