package user

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/database/sqlc"
	"teleport-plugin-slack-access-request/internal/slack"
	"teleport-plugin-slack-access-request/internal/teleport"
)

type PostgresRepository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user User) (*User, error) {
	baseEntity := database.PrePersist()

	createUserParams := sqlc.CreateUserParams{
		TeleportUserID: user.TeleportUser.TeleportUserID,
		SlackUserID:    user.SlackUser.SlackUserID,
		UseYn:          baseEntity.UseYn,
		CreateCode:     baseEntity.CreateCode,
		CreateDate:     baseEntity.CreateDate,
		Version:        baseEntity.Version,
	}

	createdUser, err := r.db.Queries.CreateUser(ctx, createUserParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in DB: %w", err)
	}

	return &User{
		UserID:       createdUser.UserID,
		TeleportUser: &teleport.User{TeleportUserID: createdUser.TeleportUserID},
		SlackUser:    &slack.User{SlackUserID: createdUser.SlackUserID},
		UseYn:        createdUser.UseYn,
		CreateCode:   createdUser.CreateCode,
		CreateDate:   createdUser.CreateDate,
		Version:      createdUser.Version,
	}, nil
}
