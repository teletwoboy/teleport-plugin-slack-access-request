package integration

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
)

type Service struct {
	SlackClt    *slack.Client
	TeleportClt *teleport.Client
}

func NewService(slackClt *slack.Client, teleportClt *teleport.Client) *Service {
	return &Service{
		SlackClt:    slackClt,
		TeleportClt: teleportClt,
	}
}

func (s *Service) SyncUsers(ctx context.Context, db *database.DB) error {
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer func(tx *sql.Tx) {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("failed to rollback transaction", "err", err)
		}
	}(tx)

	qtx := db.Queries.WithTx(tx)
	db.Queries = qtx

	slackRepo := slack.NewRepository(db)
	teleportRepo := teleport.NewRepository(db)
	userRepo := user.NewRepository(db)

	slackSrv := slack.NewService(s.SlackClt, slackRepo)
	teleportSrv := teleport.NewService(s.TeleportClt, teleportRepo)
	userSrv := user.NewService(userRepo)

	sUsers, err := slackSrv.GetUsers()
	if err != nil {
		return fmt.Errorf("error fetching slack users: %w", err)
	}

	tUsers, err := teleportSrv.GetUsersWithoutSecrets(ctx)
	if err != nil {
		return fmt.Errorf("error getting teleport users: %w", err)
	}

	users := userSrv.MapUsersByUsername(sUsers, tUsers)

	for _, u := range users {
		copiedUser := u
		createdSlackUser, err := slackSrv.CreateUser(ctx, *copiedUser.SlackUser)
		if err != nil {
			return fmt.Errorf("error creating slack user: %w", err)
		}

		createdTeleportUser, err := teleportSrv.CreateUser(ctx, *copiedUser.TeleportUser)
		if err != nil {
			return fmt.Errorf("error creating teleport user: %w", err)
		}

		copiedUser.SlackUser = createdSlackUser
		copiedUser.TeleportUser = createdTeleportUser

		_, err = userSrv.CreateUser(ctx, copiedUser)
		if err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
