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
		TargetRoleName:    ap.TargetRoleName,
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
		TargetRoleName:    createdAccessPolicy.TargetRoleName,
		TargetSlackID:     createdAccessPolicy.TargetSlackID,
		TargetRealName:    createdAccessPolicy.TargetRealName,
		UseYn:             createdAccessPolicy.UseYn,
		CreateCode:        createdAccessPolicy.CreateCode,
		CreateDate:        createdAccessPolicy.CreateDate,
		Version:           createdAccessPolicy.Version,
	}, nil
}

func (r *PostgresRepository) DeleteAccessPolicyByAccessPolicyID(ctx context.Context, id int32) error {
	baseEntity := database.MarkDelete()

	params := sqlc.DeleteAccessPolicyByAccessPolicyIDParams{
		AccessPolicyID: id,
		UseYn:          baseEntity.UseYn,
		DeleteCode:     sql.NullString{String: baseEntity.DeleteCode, Valid: baseEntity.DeleteCode != ""},
		DeleteDate:     sql.NullTime{Time: baseEntity.DeleteDate, Valid: !baseEntity.DeleteDate.IsZero()},
	}

	if err := r.q.DeleteAccessPolicyByAccessPolicyID(ctx, params); err != nil {
		return fmt.Errorf("failed to delete access policy in DB: %w", err)
	}
	return nil
}

func (r *PostgresRepository) DeleteAccessPolicyByUserID(ctx context.Context, id int32) ([]*models.AccessPolicy, error) {
	baseEntity := database.MarkDelete()

	params := sqlc.DeleteAccessPolicyByUserIDParams{
		UserID:     id,
		UseYn:      baseEntity.UseYn,
		DeleteCode: sql.NullString{String: baseEntity.DeleteCode, Valid: baseEntity.DeleteCode != ""},
		DeleteDate: sql.NullTime{Time: baseEntity.DeleteDate, Valid: !baseEntity.DeleteDate.IsZero()},
	}

	deletedAccessPolicies, err := r.q.DeleteAccessPolicyByUserID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to delete access policies in DB: %w", err)
	}

	var result []*models.AccessPolicy
	for _, deletedPolicy := range deletedAccessPolicies {
		copiedPolicy := deletedPolicy
		result = append(result, &models.AccessPolicy{
			AccessPolicyID:    copiedPolicy.AccessPolicyID,
			UserID:            copiedPolicy.UserID,
			InputChannelID:    copiedPolicy.InputChannelID,
			InputChannelName:  copiedPolicy.InputChannelName.String,
			Title:             copiedPolicy.Title,
			Reason:            copiedPolicy.Reason,
			StartDate:         copiedPolicy.StartDate,
			EndDate:           copiedPolicy.EndDate,
			Effect:            copiedPolicy.Effect,
			TargetChannelID:   copiedPolicy.TargetChannelID,
			TargetChannelName: copiedPolicy.TargetChannelName,
			TargetRole:        copiedPolicy.TargetRole,
			TargetRoleName:    copiedPolicy.TargetRoleName,
			TargetSlackID:     copiedPolicy.TargetSlackID,
			TargetRealName:    copiedPolicy.TargetRealName,
			MessageTimestamp:  copiedPolicy.MessageTimestamp.String,
			UseYn:             copiedPolicy.UseYn,
			CreateCode:        copiedPolicy.CreateCode,
			CreateDate:        copiedPolicy.CreateDate,
			UpdateCode:        copiedPolicy.UpdateCode.String,
			UpdateDate:        copiedPolicy.UpdateDate.Time,
			DeleteCode:        copiedPolicy.DeleteCode.String,
			DeleteDate:        copiedPolicy.DeleteDate.Time,
			Version:           copiedPolicy.Version,
		})
	}

	return result, nil
}

func (r *PostgresRepository) GetAccessPoliciesByAccessPolicyID(ctx context.Context, accessPolicyID int32) (*models.AccessPolicy, error) {
	row, err := r.q.GetAccessPoliciesByAccessPolicyID(ctx, accessPolicyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get access policies by access policy id in DB: %w", err)
	}
	return &models.AccessPolicy{
		AccessPolicyID:    row.AccessPolicyID,
		UserID:            row.UserID,
		InputChannelID:    row.InputChannelID,
		InputChannelName:  row.InputChannelName.String,
		Title:             row.Title,
		Reason:            row.Reason,
		StartDate:         row.StartDate,
		EndDate:           row.EndDate,
		Effect:            row.Effect,
		TargetChannelID:   row.TargetChannelID,
		TargetChannelName: row.TargetChannelName,
		TargetRole:        row.TargetRole,
		TargetRoleName:    row.TargetRoleName,
		TargetSlackID:     row.TargetSlackID,
		TargetRealName:    row.TargetRealName,
		MessageTimestamp:  row.MessageTimestamp.String,
		UseYn:             row.UseYn,
		CreateCode:        row.CreateCode,
		CreateDate:        row.CreateDate,
		UpdateCode:        row.UpdateCode.String,
		UpdateDate:        row.UpdateDate.Time,
		Version:           row.Version,
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
			TargetRoleName:    copiedRow.TargetRoleName,
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

func (r *PostgresRepository) UpdateAccessPolicyMsgTs(ctx context.Context, ap *models.AccessPolicy) error {
	baseEntity := database.MarkUpdate()

	params := sqlc.UpdateAccessPolicyMessageTimestampParams{
		AccessPolicyID:   ap.AccessPolicyID,
		MessageTimestamp: sql.NullString{String: ap.MessageTimestamp, Valid: ap.MessageTimestamp != ""},
		UpdateCode:       sql.NullString{String: baseEntity.UpdateCode, Valid: baseEntity.UpdateCode != ""},
		UpdateDate:       sql.NullTime{Time: baseEntity.UpdateDate, Valid: !baseEntity.UpdateDate.IsZero()},
	}

	if err := r.q.UpdateAccessPolicyMessageTimestamp(ctx, params); err != nil {
		return fmt.Errorf("failed to update access policy message timestamp in DB: %w", err)
	}
	return nil
}
