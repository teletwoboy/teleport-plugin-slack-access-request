/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package seedinit

import (
	"context"
	"database/sql"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/database"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/database/sqlc"
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
