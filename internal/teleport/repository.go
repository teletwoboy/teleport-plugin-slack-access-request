package teleport

import (
	"context"
	"database/sql"
	"fmt"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/database/sqlc"
	"teleport-plugin-slack-access-request/internal/teleport/models"
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
		Name:              accessRequest.Name,
		InputChannelID:    accessRequest.InputChannelID,
		InputChannelName:  sql.NullString{String: accessRequest.InputChannelName, Valid: accessRequest.InputChannelName != ""},
		Role:              accessRequest.Role,
		Reason:            sql.NullString{String: accessRequest.Reason, Valid: accessRequest.Reason != ""},
		ReviewChannelID:   accessRequest.ReviewChannelID,
		ReviewChannelName: sql.NullString{String: accessRequest.ReviewChannelName, Valid: accessRequest.ReviewChannelName != ""},
		Status:            accessRequest.Status,
		Expires:           accessRequest.Expires,
		SessionTtl:        accessRequest.SessionTTL,
		AccessDuration:    accessRequest.AccessDuration,
		StartDate:         sql.NullTime{Time: accessRequest.StartDate, Valid: !accessRequest.StartDate.IsZero()},
		ExpiryDate:        sql.NullTime{Time: accessRequest.ExpiryDate, Valid: !accessRequest.ExpiryDate.IsZero()},
		UseYn:             baseEntity.UseYn,
		CreateCode:        baseEntity.CreateCode,
		CreateDate:        baseEntity.CreateDate,
		Version:           baseEntity.Version,
	}

	createdAccessRequest, err := r.q.CreateAccessRequest(ctx, createAccessRequestParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create access request in DB: %w", err)
	}

	return &models.AccessRequest{
		AccessRequestID:   createdAccessRequest.AccessRequestID,
		RequesterUserID:   createdAccessRequest.RequesterUserID,
		Name:              createdAccessRequest.Name,
		InputChannelID:    createdAccessRequest.InputChannelID,
		InputChannelName:  createdAccessRequest.InputChannelName.String,
		Role:              createdAccessRequest.Role,
		Reason:            createdAccessRequest.Reason.String,
		ReviewChannelID:   createdAccessRequest.ReviewChannelID,
		ReviewChannelName: createdAccessRequest.ReviewChannelName.String,
		Status:            createdAccessRequest.Status,
		Expires:           createdAccessRequest.Expires,
		SessionTTL:        createdAccessRequest.SessionTtl,
		AccessDuration:    createdAccessRequest.AccessDuration,
		StartDate:         createdAccessRequest.StartDate.Time,
		ExpiryDate:        createdAccessRequest.ExpiryDate.Time,
		UseYn:             createdAccessRequest.UseYn,
		CreateCode:        createdAccessRequest.CreateCode,
		CreateDate:        createdAccessRequest.CreateDate,
		Version:           createdAccessRequest.Version,
	}, nil
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
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

func (r *PostgresRepository) ExistsAccessRequestByName(ctx context.Context, name string) (bool, error) {
	exists, err := r.q.ExistsAccessRequestByName(ctx, name)
	if err != nil {
		return false, fmt.Errorf("failed to check if access request exists in DB: %w", err)
	}
	return exists, nil
}

func (r *PostgresRepository) GetAccessRequestByName(ctx context.Context, name string) (*models.AccessRequest, error) {
	accessRequest, err := r.q.GetAccessRequestByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get access request status in DB: %w", err)
	}
	return &models.AccessRequest{
		AccessRequestID:   accessRequest.AccessRequestID,
		RequesterUserID:   accessRequest.RequesterUserID,
		Name:              accessRequest.Name,
		InputChannelID:    accessRequest.InputChannelID,
		InputChannelName:  accessRequest.InputChannelName.String,
		Role:              accessRequest.Role,
		Reason:            accessRequest.Reason.String,
		ReviewChannelID:   accessRequest.ReviewChannelID,
		ReviewChannelName: accessRequest.ReviewChannelName.String,
		Status:            accessRequest.Status,
		Expires:           accessRequest.Expires,
		SessionTTL:        accessRequest.SessionTtl,
		AccessDuration:    accessRequest.AccessDuration,
		StartDate:         accessRequest.StartDate.Time,
		ExpiryDate:        accessRequest.ExpiryDate.Time,
		UseYn:             accessRequest.UseYn,
		CreateCode:        accessRequest.CreateCode,
		CreateDate:        accessRequest.CreateDate,
		UpdateCode:        accessRequest.UpdateCode.String,
		UpdateDate:        accessRequest.UpdateDate.Time,
		Version:           accessRequest.Version,
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

func (r *PostgresRepository) UpdateAccessRequestStatusByName(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error) {
	baseEntity := database.MarkUpdate()

	updateAccessRequestStatusByNameParams := sqlc.UpdateAccessRequestStatusByNameParams{
		Status:         accessRequest.Status,
		Expires:        accessRequest.Expires,
		SessionTtl:     accessRequest.SessionTTL,
		AccessDuration: accessRequest.AccessDuration,
		StartDate:      sql.NullTime{Time: accessRequest.StartDate, Valid: !accessRequest.StartDate.IsZero()},
		ExpiryDate:     sql.NullTime{Time: accessRequest.ExpiryDate, Valid: !accessRequest.ExpiryDate.IsZero()},
		UpdateCode:     sql.NullString{String: baseEntity.UpdateCode, Valid: baseEntity.UpdateCode != ""},
		UpdateDate:     sql.NullTime{Time: baseEntity.UpdateDate, Valid: !baseEntity.UpdateDate.IsZero()},
		Name:           accessRequest.Name,
	}

	row, err := r.q.UpdateAccessRequestStatusByName(ctx, updateAccessRequestStatusByNameParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update access request status in DB: %w", err)
	}
	return &models.AccessRequest{
		AccessRequestID:   row.AccessRequestID,
		RequesterUserID:   row.RequesterUserID,
		Name:              row.Name,
		InputChannelID:    row.InputChannelID,
		InputChannelName:  row.InputChannelName.String,
		Role:              row.Role,
		Reason:            row.Reason.String,
		ReviewChannelID:   row.ReviewChannelID,
		ReviewChannelName: row.ReviewChannelName.String,
		Status:            row.Status,
		Expires:           row.Expires,
		SessionTTL:        row.SessionTtl,
		AccessDuration:    row.AccessDuration,
		StartDate:         row.StartDate.Time,
		ExpiryDate:        row.ExpiryDate.Time,
		CreateCode:        row.CreateCode,
		CreateDate:        row.CreateDate,
		UpdateCode:        row.UpdateCode.String,
		UpdateDate:        row.UpdateDate.Time,
	}, nil
}
