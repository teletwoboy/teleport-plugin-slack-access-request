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
		SessionTtl:        accessRequest.SessionTtl,
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
		SessionTtl:        createdAccessRequest.SessionTtl,
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

func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	row, err := r.q.GetTeleportUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get teleport user by username %s: %w", username, err)
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
