package seedinit

import (
	"context"
	"database/sql"
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

func (r *PostgresRepository) Create(ctx context.Context) error {
	baseEntity := database.MarkCreate()

	createSeedInitParams := sqlc.CreateSeedInitParams{
		UseYn:      baseEntity.UseYn,
		CreateCode: baseEntity.CreateCode,
		CreateDate: baseEntity.CreateDate,
		Version:    baseEntity.Version,
	}

	err := r.db.Queries.CreateSeedInit(ctx, createSeedInitParams)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) UpdateStaus(ctx context.Context) error {
	baseEntity := database.MarkUpdate()

	updateSeedInitParams := sqlc.UpdateSeedInitStatusParams{
		UpdateCode: sql.NullString{String: baseEntity.UpdateCode, Valid: baseEntity.UpdateCode != ""},
		UpdateDate: sql.NullTime{Time: baseEntity.UpdateDate, Valid: baseEntity.UpdateDate.IsZero()},
	}

	err := r.db.Queries.UpdateSeedInitStatus(ctx, updateSeedInitParams)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) GetStatus(ctx context.Context) (*SeedInit, error) {
	seedInit, err := r.db.Queries.GetSeedInitStatus(ctx)
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
