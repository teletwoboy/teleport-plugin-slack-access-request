package teleport

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/database/sqlc"
)

type PostgresRepository struct {
	q sqlc.Querier
}

func NewRepository(q sqlc.Querier) *PostgresRepository {
	return &PostgresRepository{q: q}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user User) (*User, error) {
	baseEntity := database.MarkCreate()

	createTeleportUserParams := sqlc.CreateTeleportUserParams{
		Username:   user.Username,
		UseYn:      baseEntity.UseYn,
		CreateCode: baseEntity.CreateCode,
		CreateDate: baseEntity.CreateDate,
		Version:    baseEntity.Version,
	}

	createdTeleportUser, err := r.q.CreateTeleportUser(ctx, createTeleportUserParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create teleport user in DB: %w", err)
	}
	return &User{
		TeleportUserID: createdTeleportUser.TeleportUserID,
		Username:       createdTeleportUser.Username,
		UseYn:          createdTeleportUser.UseYn,
		CreateCode:     createdTeleportUser.CreateCode,
		CreateDate:     createdTeleportUser.CreateDate,
		Version:        createdTeleportUser.Version,
	}, nil
}

func (r *PostgresRepository) GetUserByTeleportUserID(ctx context.Context, id int32) (*User, error) {
	row, err := r.q.GetTeleportUserByTeleportUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get teleport user by teleport user id %d: %w", id, err)
	}
	return &User{
		TeleportUserID: row.TeleportUserID,
		Username:       row.Username,
		UseYn:          row.UseYn,
		CreateCode:     row.CreateCode,
		CreateDate:     row.CreateDate,
		UpdateCode:     row.UpdateCode.String,
		UpdateDate:     row.UpdateDate.Time,
		DeleteCode:     row.DeleteCode.String,
		DeleteDate:     row.DeleteDate.Time,
		Version:        row.Version,
	}, nil
}
