package seedinit

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
	"teleport-plugin-slack-access-request/internal/user"
	"teleport-plugin-slack-access-request/internal/user/models"
)

type Service interface {
	create(ctx context.Context) error
	updateStaus(ctx context.Context) error
	getStatus(ctx context.Context) (string, error)
	Init(ctx context.Context, db *database.DB, sClt *slack.Client, tClt *teleport.Client) error
}

type Repository interface {
	Create(ctx context.Context) error
	UpdateStaus(ctx context.Context) error
	GetStatus(ctx context.Context) (*SeedInit, error)
}

type service struct {
	repo        Repository
	slackSrv    slack.Service
	teleportSrv teleport.Service
	userSrv     user.Service
}

func NewService(r Repository, s slack.Service, t teleport.Service, u user.Service) Service {
	return &service{
		repo:        r,
		slackSrv:    s,
		teleportSrv: t,
		userSrv:     u,
	}
}

func (s *service) Init(ctx context.Context, db *database.DB, sClt *slack.Client, tClt *teleport.Client) error {
	if err := s.create(ctx); err != nil {
		return fmt.Errorf("failed to create seedinit: %w", err)
	}

	status, err := s.getStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	if status != "initialized" {
		slog.Info("started to initialize seed")

		users, err := s.userSrv.FetchUsers(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch users: %w", err)
		}
		slog.Info("total users - before committed: ", "count", len(users))

		tx, err := db.Conn.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		committed := false
		defer func(tx *sql.Tx) {
			if !committed {
				if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
					slog.Error("failed to rollback transaction", "err", err)
				}
			}
		}(tx)

		qtx := db.Queries.WithTx(tx)

		slackTxRepo := slack.NewRepository(qtx)
		teleportTxRepo := teleport.NewRepository(qtx)
		userTxRepo := user.NewRepository(qtx)
		seedInitTxRepo := NewRepository(qtx)

		slackTxSrv := slack.NewService(sClt, slackTxRepo)
		teleportTxSrv := teleport.NewService(tClt, teleportTxRepo)
		userTxSrv := user.NewService(userTxRepo, slackTxSrv, teleportTxSrv)
		seedInitTxSrv := NewService(seedInitTxRepo, slackTxSrv, teleportTxSrv, userTxSrv)

		var createdUsers []models.User
		for _, u := range users {
			copiedUser := u
			createdSlackUser, err := slackTxSrv.CreateUser(ctx, *copiedUser.SlackUser)
			if err != nil {
				return fmt.Errorf("failed to create slack user: %w", err)
			}

			createdTeleportUser, err := teleportTxSrv.CreateUser(ctx, *copiedUser.TeleportUser)
			if err != nil {
				return fmt.Errorf("failed to create teleport user: %w", err)
			}

			copiedUser.SlackUser = createdSlackUser
			copiedUser.TeleportUser = createdTeleportUser

			createdUser, err := userTxSrv.CreateUser(ctx, copiedUser)
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}

			createdUsers = append(createdUsers, *createdUser)
		}

		if err := seedInitTxSrv.updateStaus(ctx); err != nil {
			return fmt.Errorf("failed to update seedinit staus: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		committed = true
		slog.Info("total users - after committed: ", "count", len(createdUsers))
		slog.Info("successfully initialized seed")
	} else {
		slog.Info("already initialized seed")
	}
	return nil
}

func (s *service) create(ctx context.Context) error {
	err := s.repo.Create(ctx)
	if err != nil {
		return fmt.Errorf("failed create seedinit: %w", err)
	}
	return nil
}

func (s *service) updateStaus(ctx context.Context) error {
	err := s.repo.UpdateStaus(ctx)
	if err != nil {
		return fmt.Errorf("failed update staus: %w", err)
	}
	return nil
}

func (s *service) getStatus(ctx context.Context) (string, error) {
	seedInit, err := s.repo.GetStatus(ctx)
	if err != nil {
		return "", fmt.Errorf("failed get seedinit status: %w", err)
	}
	return seedInit.Status, nil
}
