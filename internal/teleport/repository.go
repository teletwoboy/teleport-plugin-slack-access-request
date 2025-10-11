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

package teleport

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/database"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/database/sqlc"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/teleport/models"
)

type PostgresRepository struct {
	q sqlc.Querier
}

func NewRepository(q sqlc.Querier) *PostgresRepository {
	return &PostgresRepository{q: q}
}

func (r *PostgresRepository) CreateAccessRequest(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error) {
	baseEntity := database.MarkCreate()

	createAccessRequestParams := sqlc.CreateAccessRequestParams{
		RequesterUserID:   accessRequest.RequesterUserID,
		InputChannelID:    accessRequest.InputChannelID,
		InputChannelName:  accessRequest.InputChannelName,
		Role:              accessRequest.Role,
		Reason:            sql.NullString{String: accessRequest.Reason, Valid: accessRequest.Reason != ""},
		ReviewChannelID:   accessRequest.ReviewChannelID,
		ReviewChannelName: accessRequest.ReviewChannelName,
		State:             accessRequest.State,
		UseYn:             baseEntity.UseYn,
		CreateCode:        baseEntity.CreateCode,
		CreateDate:        baseEntity.CreateDate,
		Version:           baseEntity.Version,
	}

	row, err := r.q.CreateAccessRequest(ctx, createAccessRequestParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create access request in DB: %w", err)
	}

	return &models.AccessRequest{
		AccessRequestID:   row.AccessRequestID,
		RequesterUserID:   row.RequesterUserID,
		InputChannelID:    row.InputChannelID,
		InputChannelName:  row.InputChannelName,
		Role:              row.Role,
		Reason:            row.Reason.String,
		ReviewChannelID:   row.ReviewChannelID,
		ReviewChannelName: row.ReviewChannelName,
		State:             row.State,
		UseYn:             row.UseYn,
		CreateCode:        row.CreateCode,
		CreateDate:        row.CreateDate,
		Version:           row.Version,
	}, nil
}

func (r *PostgresRepository) CreateAccessReview(ctx context.Context, accessReview *models.AccessReview) (*models.AccessReview, error) {
	baseEntity := database.MarkCreate()

	createAccessReviewParams := sqlc.CreateAccessReviewParams{
		AccessRequestID: accessReview.AccessRequestID,
		ReviewerUserID:  accessReview.ReviewerUserID,
		Reason:          sql.NullString{String: accessReview.Reason, Valid: accessReview.Reason != ""},
		Decision:        accessReview.Decision,
		UseYn:           baseEntity.UseYn,
		CreateCode:      baseEntity.CreateCode,
		CreateDate:      baseEntity.CreateDate,
	}

	createdAccessReview, err := r.q.CreateAccessReview(ctx, createAccessReviewParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create access review in DB: %w", err)
	}
	return &models.AccessReview{
		AccessReviewID:  createdAccessReview.AccessReviewID,
		AccessRequestID: createdAccessReview.AccessRequestID,
		ReviewerUserID:  createdAccessReview.ReviewerUserID,
		Reason:          createdAccessReview.Reason.String,
		Decision:        createdAccessReview.Decision,
		UseYn:           createdAccessReview.UseYn,
		CreateCode:      createdAccessReview.CreateCode,
		CreateDate:      createdAccessReview.CreateDate,
		Version:         createdAccessReview.Version,
	}, nil
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
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
	return &models.User{
		TeleportUserID: createdTeleportUser.TeleportUserID,
		Username:       createdTeleportUser.Username,
		UseYn:          createdTeleportUser.UseYn,
		CreateCode:     createdTeleportUser.CreateCode,
		CreateDate:     createdTeleportUser.CreateDate,
		Version:        createdTeleportUser.Version,
	}, nil
}

func (r *PostgresRepository) DeleteUser(ctx context.Context, user *models.User) (*models.User, error) {
	baseEntity := database.MarkDelete()

	DeleteTeleportUserUseYnByUsernameParams := sqlc.DeleteTeleportUserUseYnByUsernameParams{
		Username:   user.Username,
		UseYn:      baseEntity.UseYn,
		DeleteCode: sql.NullString{String: baseEntity.DeleteCode, Valid: baseEntity.DeleteCode != ""},
		DeleteDate: sql.NullTime{Time: baseEntity.DeleteDate, Valid: !baseEntity.DeleteDate.IsZero()},
	}

	deletedTeleportUser, err := r.q.DeleteTeleportUserUseYnByUsername(ctx, DeleteTeleportUserUseYnByUsernameParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create teleport user in DB: %w", err)
	}
	return &models.User{
		TeleportUserID: deletedTeleportUser.TeleportUserID,
		Username:       deletedTeleportUser.Username,
		UseYn:          deletedTeleportUser.UseYn,
		CreateCode:     deletedTeleportUser.CreateCode,
		CreateDate:     deletedTeleportUser.CreateDate,
		Version:        deletedTeleportUser.Version,
	}, nil
}

func (r *PostgresRepository) ExistsUserByUsername(ctx context.Context, username string) (bool, error) {
	exists, err := r.q.ExistsTeleportUserByUsername(ctx, username)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return exists, nil
}

func (r *PostgresRepository) GetAccessRequestByAccessRequestID(ctx context.Context, accessRequestID int32) (*models.AccessRequest, error) {
	row, err := r.q.GetAccessRequestByAccessRequestID(ctx, accessRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get access request by access request id in DB: %w", err)
	}
	return &models.AccessRequest{
		AccessRequestID:   row.AccessRequestID,
		RequesterUserID:   row.RequesterUserID,
		Name:              row.Name.String,
		InputChannelID:    row.InputChannelID,
		InputChannelName:  row.InputChannelName,
		Role:              row.Role,
		Reason:            row.Reason.String,
		ReviewChannelID:   row.ReviewChannelID,
		ReviewChannelName: row.ReviewChannelName,
		State:             row.State,
		StartDate:         row.StartDate.Time,
		AccessDuration:    row.AccessDuration.Time,
		RequestTTL:        row.RequestTtl.Time,
		UseYn:             row.UseYn,
		CreateCode:        row.CreateCode,
		CreateDate:        row.CreateDate,
		UpdateCode:        row.UpdateCode.String,
		UpdateDate:        row.UpdateDate.Time,
		Version:           row.Version,
	}, nil
}

func (r *PostgresRepository) GetAccessRequestByName(ctx context.Context, name string) (*models.AccessRequest, error) {
	accessRequest, err := r.q.GetAccessRequestByName(ctx, sql.NullString{String: name, Valid: name != ""})
	if err != nil {
		return nil, fmt.Errorf("failed to get access request state in DB: %w", err)
	}
	return &models.AccessRequest{
		AccessRequestID:   accessRequest.AccessRequestID,
		RequesterUserID:   accessRequest.RequesterUserID,
		Name:              accessRequest.Name.String,
		InputChannelID:    accessRequest.InputChannelID,
		InputChannelName:  accessRequest.InputChannelName,
		Role:              accessRequest.Role,
		Reason:            accessRequest.Reason.String,
		ReviewChannelID:   accessRequest.ReviewChannelID,
		ReviewChannelName: accessRequest.ReviewChannelName,
		State:             accessRequest.State,
		StartDate:         accessRequest.StartDate.Time,
		AccessDuration:    accessRequest.AccessDuration.Time,
		RequestTTL:        accessRequest.RequestTtl.Time,
		UseYn:             accessRequest.UseYn,
		CreateCode:        accessRequest.CreateCode,
		CreateDate:        accessRequest.CreateDate,
		UpdateCode:        accessRequest.UpdateCode.String,
		UpdateDate:        accessRequest.UpdateDate.Time,
		Version:           accessRequest.Version,
	}, nil
}

func (r *PostgresRepository) GetAccessRequestStateByName(ctx context.Context, name string) (string, error) {
	state, err := r.q.GetAccessRequestStateByName(ctx, sql.NullString{
		String: name,
		Valid:  name != "",
	})
	if err != nil {
		return "", fmt.Errorf("failed to get access request state in DB: %w", err)
	}
	return state, nil
}

func (r *PostgresRepository) GetAccessReviewByAccessReviewID(ctx context.Context, accessReviewID int32) (*models.AccessReview, error) {
	row, err := r.q.GetAccessReviewByAccessReviewID(ctx, accessReviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to get access review in DB: %w", err)
	}
	return &models.AccessReview{
		AccessReviewID:  row.AccessReviewID,
		AccessRequestID: row.AccessRequestID,
		ReviewerUserID:  row.ReviewerUserID,
		Reason:          row.Reason.String,
		Decision:        row.Decision,
		UseYn:           row.UseYn,
		CreateCode:      row.CreateCode,
		CreateDate:      row.CreateDate,
		UpdateCode:      row.UpdateCode.String,
		UpdateDate:      row.UpdateDate.Time,
		Version:         row.Version,
	}, nil
}

func (r *PostgresRepository) GetUserByTeleportUserID(ctx context.Context, id int32) (*models.User, error) {
	row, err := r.q.GetTeleportUserByTeleportUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by teleport user id from DB: %w", err)
	}
	return &models.User{
		TeleportUserID: row.TeleportUserID,
		Username:       row.Username,
		UseYn:          row.UseYn,
		CreateCode:     row.CreateCode,
		CreateDate:     row.CreateDate,
		UpdateCode:     row.UpdateCode.String,
		UpdateDate:     row.UpdateDate.Time,
		Version:        row.Version,
	}, nil
}

func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	row, err := r.q.GetTeleportUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get teleport user by username in DB, %s: %w", username, err)
	}
	return &models.User{
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
	}, err
}

func (r *PostgresRepository) UpdateAccessRequestByAccessRequestID(ctx context.Context, ar *models.AccessRequest) error {
	baseEntity := database.MarkUpdate()

	updateAccessRequestByAccessRequestIDParams := sqlc.UpdateAccessRequestByAccessRequestIDParams{
		AccessRequestID: ar.AccessRequestID,
		Name:            sql.NullString{String: ar.Name, Valid: ar.Name != ""},
		StartDate:       sql.NullTime{Time: ar.StartDate, Valid: !ar.StartDate.IsZero()},
		AccessDuration:  sql.NullTime{Time: ar.AccessDuration, Valid: !ar.AccessDuration.IsZero()},
		RequestTtl:      sql.NullTime{Time: ar.RequestTTL, Valid: !ar.RequestTTL.IsZero()},
		UpdateCode:      sql.NullString{String: baseEntity.UpdateCode, Valid: baseEntity.UpdateCode != ""},
		UpdateDate:      sql.NullTime{Time: ar.UpdateDate, Valid: !ar.UpdateDate.IsZero()},
	}

	if err := r.q.UpdateAccessRequestByAccessRequestID(ctx, updateAccessRequestByAccessRequestIDParams); err != nil {
		return fmt.Errorf("failed to update access request by access request id in DB: %w", err)
	}
	return nil
}

func (r *PostgresRepository) UpdateAccessRequestStateByName(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error) {
	baseEntity := database.MarkUpdate()

	updateAccessRequestStateByNameParams := sqlc.UpdateAccessRequestStateByNameParams{
		State:      accessRequest.State,
		StartDate:  sql.NullTime{Time: accessRequest.StartDate, Valid: !accessRequest.StartDate.IsZero()},
		RequestTtl: sql.NullTime{Time: accessRequest.RequestTTL, Valid: !accessRequest.RequestTTL.IsZero()},
		UpdateCode: sql.NullString{String: baseEntity.UpdateCode, Valid: baseEntity.UpdateCode != ""},
		UpdateDate: sql.NullTime{Time: baseEntity.UpdateDate, Valid: !baseEntity.UpdateDate.IsZero()},
		Name:       sql.NullString{String: accessRequest.Name, Valid: accessRequest.Name != ""},
	}

	row, err := r.q.UpdateAccessRequestStateByName(ctx, updateAccessRequestStateByNameParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update access request state in DB: %w", err)
	}
	return &models.AccessRequest{
		AccessRequestID:   row.AccessRequestID,
		RequesterUserID:   row.RequesterUserID,
		Name:              row.Name.String,
		InputChannelID:    row.InputChannelID,
		InputChannelName:  row.InputChannelName,
		Role:              row.Role,
		Reason:            row.Reason.String,
		ReviewChannelID:   row.ReviewChannelID,
		ReviewChannelName: row.ReviewChannelName,
		State:             row.State,
		StartDate:         row.StartDate.Time,
		AccessDuration:    row.AccessDuration.Time,
		RequestTTL:        row.RequestTtl.Time,
		CreateCode:        row.CreateCode,
		CreateDate:        row.CreateDate,
		UpdateCode:        row.UpdateCode.String,
		UpdateDate:        row.UpdateDate.Time,
	}, nil
}
