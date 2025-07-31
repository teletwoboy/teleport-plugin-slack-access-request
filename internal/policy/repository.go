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

func (r *PostgresRepository) DeleteAccessPolicyByAccessPolicyID(ctx context.Context, id int32) (*models.AccessPolicy, error) {
	baseEntity := database.MarkDelete()

	params := sqlc.DeleteAccessPolicyByAccessPolicyIDParams{
		AccessPolicyID: id,
		UseYn:          baseEntity.UseYn,
		DeleteCode:     sql.NullString{String: baseEntity.DeleteCode, Valid: baseEntity.DeleteCode != ""},
		DeleteDate:     sql.NullTime{Time: baseEntity.DeleteDate, Valid: !baseEntity.DeleteDate.IsZero()},
	}

	deletedAccessPolicy, err := r.q.DeleteAccessPolicyByAccessPolicyID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to delete access policy in DB: %w", err)
	}
	return &models.AccessPolicy{
		AccessPolicyID:    deletedAccessPolicy.AccessPolicyID,
		UserID:            deletedAccessPolicy.UserID,
		InputChannelID:    deletedAccessPolicy.InputChannelID,
		InputChannelName:  deletedAccessPolicy.InputChannelName.String,
		Title:             deletedAccessPolicy.Title,
		Reason:            deletedAccessPolicy.Reason,
		StartDate:         deletedAccessPolicy.StartDate,
		EndDate:           deletedAccessPolicy.EndDate,
		Effect:            deletedAccessPolicy.Effect,
		TargetChannelID:   deletedAccessPolicy.TargetChannelID,
		TargetChannelName: deletedAccessPolicy.TargetChannelName,
		TargetRole:        deletedAccessPolicy.TargetRole,
		TargetRoleName:    deletedAccessPolicy.TargetRole,
		TargetSlackID:     deletedAccessPolicy.TargetSlackID,
		TargetRealName:    deletedAccessPolicy.TargetRealName,
		MessageTimestamp:  deletedAccessPolicy.MessageTimestamp.String,
		UseYn:             deletedAccessPolicy.UseYn,
		CreateCode:        deletedAccessPolicy.CreateCode,
		CreateDate:        deletedAccessPolicy.CreateDate,
		UpdateCode:        deletedAccessPolicy.UpdateCode.String,
		UpdateDate:        deletedAccessPolicy.UpdateDate.Time,
		DeleteCode:        deletedAccessPolicy.DeleteCode.String,
		DeleteDate:        deletedAccessPolicy.DeleteDate.Time,
		Version:           deletedAccessPolicy.Version,
	}, nil
}

func (r *PostgresRepository) GetAccessPoliciesByInputChannelID(ctx context.Context, channelID string) ([]*models.AccessPolicy, error) {
	rows, err := r.q.GetAccessPoliciesByInputChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get access policies by input channel id in DB: %w", err)
	}
	var result []*models.AccessPolicy
	for _, r := range rows {
		copiedRow := r
		result = append(result, &models.AccessPolicy{
			AccessPolicyID:    copiedRow.AccessPolicyID,
			UserID:            copiedRow.UserID,
			InputChannelID:    copiedRow.InputChannelID,
			InputChannelName:  copiedRow.InputChannelName.String,
			Title:             copiedRow.Title,
			Reason:            copiedRow.Reason,
			StartDate:         copiedRow.StartDate,
			EndDate:           copiedRow.EndDate,
			Effect:            copiedRow.Effect,
			TargetChannelID:   copiedRow.TargetChannelID,
			TargetChannelName: copiedRow.TargetChannelName,
			TargetRole:        copiedRow.TargetRole,
			TargetRoleName:    copiedRow.TargetRole,
			TargetSlackID:     copiedRow.TargetSlackID,
			TargetRealName:    copiedRow.TargetRealName,
			MessageTimestamp:  copiedRow.MessageTimestamp.String,
			UseYn:             copiedRow.UseYn,
			CreateCode:        copiedRow.CreateCode,
			CreateDate:        copiedRow.CreateDate,
			UpdateCode:        copiedRow.UpdateCode.String,
			UpdateDate:        copiedRow.UpdateDate.Time,
			Version:           copiedRow.Version,
		})
	}
	return result, nil
}

func (r *PostgresRepository) UpdateAccessPolicyMessageTimestamp(ctx context.Context, ap *models.AccessPolicy) error {
	baseEntity := database.MarkUpdate()

	params := sqlc.UpdateAccessPolicyMessageTimestampParams{
		AccessPolicyID:   ap.AccessPolicyID,
		MessageTimestamp: sql.NullString{String: ap.MessageTimestamp, Valid: ap.MessageTimestamp != ""},
		UpdateCode:       sql.NullString{String: baseEntity.UpdateCode, Valid: baseEntity.UpdateCode != ""},
		UpdateDate:       sql.NullTime{Time: baseEntity.UpdateDate, Valid: !baseEntity.UpdateDate.IsZero()},
	}

	err := r.q.UpdateAccessPolicyMessageTimestamp(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update access policy message timestamp in DB: %w", err)
	}
	return nil
}
