package seedinit

import (
	"context"
	"database/sql"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/database/sqlc"
)

type PostgresRepository struct {
	q sqlc.Querier
}

func NewRepository(q sqlc.Querier) *PostgresRepository {
	return &PostgresRepository{q: q}
}

func (r *PostgresRepository) Create(ctx context.Context) error {
	baseEntity := database.MarkCreate()

	createSeedInitParams := sqlc.CreateSeedInitParams{
		UseYn:      baseEntity.UseYn,
		CreateCode: baseEntity.CreateCode,
		CreateDate: baseEntity.CreateDate,
		Version:    baseEntity.Version,
	}

	err := r.q.CreateSeedInit(ctx, createSeedInitParams)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) UpdateStaus(ctx context.Context) error {
	baseEntity := database.MarkUpdate()

	updateSeedInitParams := sqlc.UpdateSeedInitStatusParams{
		UpdateCode: sql.NullString{String: baseEntity.UpdateCode, Valid: baseEntity.UpdateCode != ""},
		UpdateDate: sql.NullTime{Time: baseEntity.UpdateDate, Valid: !baseEntity.UpdateDate.IsZero()},
	}

	err := r.q.UpdateSeedInitStatus(ctx, updateSeedInitParams)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) GetStatus(ctx context.Context) (*SeedInit, error) {
	seedInit, err := r.q.GetSeedInitStatus(ctx)
	if err != nil {
		return nil, err
	}
	return &SeedInit{
		SeedInitId: seedInit.SeedInitID,
		Status:     seedInit.Status,
		UseYn:      seedInit.UseYn,
		CreateCode: seedInit.CreateCode,
		CreateDate: seedInit.CreateDate,
		UpdateCode: seedInit.UpdateCode.String,
		UpdateDate: seedInit.UpdateDate.Time,
		Version:    seedInit.Version,
	}, nil
}
