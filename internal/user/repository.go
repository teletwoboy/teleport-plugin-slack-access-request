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
	q sqlc.Querier
}

func NewRepository(q sqlc.Querier) *PostgresRepository {
	return &PostgresRepository{q: q}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user User) (*User, error) {
	baseEntity := database.MarkCreate()

	createUserParams := sqlc.CreateUserParams{
		TeleportUserID: user.TeleportUser.TeleportUserID,
		SlackUserID:    user.SlackUser.SlackUserID,
		UseYn:          baseEntity.UseYn,
		CreateCode:     baseEntity.CreateCode,
		CreateDate:     baseEntity.CreateDate,
		Version:        baseEntity.Version,
	}

	createdUser, err := r.q.CreateUser(ctx, createUserParams)
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

func (r *PostgresRepository) GetUserBySlackUserID(ctx context.Context, id int32) (*User, error) {
	row, err := r.q.GetUserBySlackUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by slack user id %d: %w", id, err)
	}
	return &User{
		UserID:       row.UserID,
		TeleportUser: &teleport.User{TeleportUserID: row.TeleportUserID},
		SlackUser:    &slack.User{SlackUserID: row.SlackUserID},
		UseYn:        row.UseYn,
		CreateCode:   row.CreateCode,
		CreateDate:   row.CreateDate,
		UpdateCode:   row.UpdateCode.String,
		UpdateDate:   row.UpdateDate.Time,
		DeleteCode:   row.DeleteCode.String,
		DeleteDate:   row.DeleteDate.Time,
		Version:      row.Version,
	}, nil
}
