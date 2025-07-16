package teleport

import (
	"context"
	"fmt"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/database/sqlc"
)

type PostgresRepository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func NewRepositoryWithTx(qtx *sqlc.Queries) *PostgresRepository {
	return &PostgresRepository{
		db: &database.DB{
			Conn:    nil,
			Queries: qtx,
		},
	}
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

	createdTeleportUser, err := r.db.Queries.CreateTeleportUser(ctx, createTeleportUserParams)
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
