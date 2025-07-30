package policy

import (
	"context"
	"database/sql"
	"fmt"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/database/sqlc"
	"teleport-plugin-slack-access-request/internal/policy/models"
)

type PostgresRepository struct {
	q sqlc.Querier
}

func NewRepository(q sqlc.Querier) *PostgresRepository {
	return &PostgresRepository{q: q}
}

func (r *PostgresRepository) CreateAccessPolicy(ctx context.Context, ap *models.AccessPolicy) (*models.AccessPolicy, error) {
	baseEntity := database.MarkCreate()

	createAccessPolicyParams := sqlc.CreateAccessPolicyParams{
		UserID:            ap.UserID,
		InputChannelID:    ap.InputChannelID,
		InputChannelName:  sql.NullString{String: ap.InputChannelName, Valid: ap.InputChannelName != ""},
		Title:             ap.Title,
		Reason:            ap.Reason,
		TimeZone:          ap.TimeZone,
		StartDate:         ap.StartDate,
		EndDate:           ap.EndDate,
		Effect:            ap.Effect,
		TargetChannelID:   ap.TargetChannelID,
		TargetChannelName: ap.TargetChannelName,
		TargetRole:        ap.TargetRole,
		TargetRoleName:    ap.TargetRole,
		TargetSlackID:     ap.TargetSlackID,
		TargetRealName:    ap.TargetRealName,
		UseYn:             baseEntity.UseYn,
		CreateCode:        baseEntity.CreateCode,
		CreateDate:        baseEntity.CreateDate,
		Version:           baseEntity.Version,
	}

	createdAccessPolicy, err := r.q.CreateAccessPolicy(ctx, createAccessPolicyParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create access policy in DB: %w", err)
	}
	return &models.AccessPolicy{
		AccessPolicyID:    createdAccessPolicy.AccessPolicyID,
		UserID:            createdAccessPolicy.UserID,
		InputChannelID:    createdAccessPolicy.InputChannelID,
		InputChannelName:  createdAccessPolicy.InputChannelName.String,
		Title:             createdAccessPolicy.Title,
		Reason:            createdAccessPolicy.Reason,
		TimeZone:          createdAccessPolicy.TimeZone,
		StartDate:         createdAccessPolicy.StartDate,
		EndDate:           createdAccessPolicy.EndDate,
		Effect:            createdAccessPolicy.Effect,
		TargetChannelID:   createdAccessPolicy.TargetChannelID,
		TargetChannelName: createdAccessPolicy.TargetChannelName,
		TargetRole:        createdAccessPolicy.TargetRole,
		TargetRoleName:    createdAccessPolicy.TargetRole,
		TargetSlackID:     createdAccessPolicy.TargetSlackID,
		TargetRealName:    createdAccessPolicy.TargetRealName,
		UseYn:             createdAccessPolicy.UseYn,
		CreateCode:        createdAccessPolicy.CreateCode,
		CreateDate:        createdAccessPolicy.CreateDate,
		Version:           createdAccessPolicy.Version,
	}, nil
}
