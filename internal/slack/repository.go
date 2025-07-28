package slack

import (
	"context"
	"database/sql"
	"fmt"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/database/sqlc"
	"teleport-plugin-slack-access-request/internal/slack/models"
)

type PostgresRepository struct {
	q sqlc.Querier
}

func NewRepository(q sqlc.Querier) *PostgresRepository {
	return &PostgresRepository{q: q}
}

// CreateUser creates a new Slack user in the database.
// This operation executes a single INSERT statement and does not require an explicit transaction.
func (r *PostgresRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	baseEntity := database.MarkCreate()

	createSlackUserParams := sqlc.CreateSlackUserParams{
		ID:         user.ID,
		Name:       user.Name,
		RealName:   sql.NullString{String: user.RealName, Valid: user.RealName != ""},
		Email:      user.Email,
		UseYn:      baseEntity.UseYn,
		CreateCode: baseEntity.CreateCode,
		CreateDate: baseEntity.CreateDate,
		Version:    baseEntity.Version,
	}

	createdSlackUser, err := r.q.CreateSlackUser(ctx, createSlackUserParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create slack user in DB: %w", err)
	}
	return &models.User{
		SlackUserID: createdSlackUser.SlackUserID,
		ID:          createdSlackUser.ID,
		Name:        createdSlackUser.Name,
		RealName:    createdSlackUser.RealName.String,
		Email:       createdSlackUser.Email,
		UseYn:       createdSlackUser.UseYn,
		CreateCode:  createdSlackUser.CreateCode,
		CreateDate:  createdSlackUser.CreateDate,
		Version:     createdSlackUser.Version,
	}, nil
}

func (r *PostgresRepository) DeleteUser(ctx context.Context, user *models.User) (*models.User, error) {
	baseEntity := database.MarkDelete()

	deleteSlackUserByNameParams := sqlc.DeleteSlackUserByNameParams{
		Name:       user.Name,
		UseYn:      baseEntity.UseYn,
		DeleteCode: sql.NullString{String: baseEntity.DeleteCode, Valid: baseEntity.DeleteCode != ""},
		DeleteDate: sql.NullTime{Time: baseEntity.DeleteDate, Valid: !baseEntity.DeleteDate.IsZero()},
	}

	deletedSlackUser, err := r.q.DeleteSlackUserByName(ctx, deleteSlackUserByNameParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create slack user in DB: %w", err)
	}
	return &models.User{
		SlackUserID: deletedSlackUser.SlackUserID,
		ID:          deletedSlackUser.ID,
		Name:        deletedSlackUser.Name,
		RealName:    deletedSlackUser.RealName.String,
		Email:       deletedSlackUser.Email,
		UseYn:       deletedSlackUser.UseYn,
		CreateCode:  deletedSlackUser.CreateCode,
		CreateDate:  deletedSlackUser.CreateDate,
		Version:     deletedSlackUser.Version,
	}, nil
}

func (r *PostgresRepository) ExistsUserByID(ctx context.Context, id string) (bool, error) {
	exists, err := r.q.ExistsSlackUserByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return exists, nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	row, err := r.q.GetSlackUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &models.User{
		SlackUserID: row.SlackUserID,
		ID:          row.ID,
		Name:        row.Name,
		RealName:    row.RealName.String,
		Email:       row.Email,
		UseYn:       row.UseYn,
		CreateCode:  row.CreateCode,
		CreateDate:  row.CreateDate,
		UpdateCode:  row.UpdateCode.String,
		UpdateDate:  row.UpdateDate.Time,
		Version:     row.Version,
	}, nil
}

func (r *PostgresRepository) GetUserByName(ctx context.Context, name string) (*models.User, error) {
	row, err := r.q.GetSlackUserByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return &models.User{
		SlackUserID: row.SlackUserID,
		ID:          row.ID,
		Name:        row.Name,
		RealName:    row.RealName.String,
		Email:       row.Email,
		UseYn:       row.UseYn,
		CreateCode:  row.CreateCode,
		CreateDate:  row.CreateDate,
		UpdateCode:  row.UpdateCode.String,
		UpdateDate:  row.UpdateDate.Time,
		Version:     row.Version,
	}, nil
}

func (r *PostgresRepository) GetUserBySlackUserID(ctx context.Context, id int32) (*models.User, error) {
	row, err := r.q.GetSlackUserBySlackUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by Slack user ID: %w", err)
	}
	return &models.User{
		SlackUserID: row.SlackUserID,
		ID:          row.ID,
		Name:        row.Name,
		RealName:    row.RealName.String,
		Email:       row.Email,
		UseYn:       row.UseYn,
		CreateCode:  row.CreateCode,
		CreateDate:  row.CreateDate,
		UpdateCode:  row.UpdateCode.String,
		UpdateDate:  row.UpdateDate.Time,
		Version:     row.Version,
	}, nil
}
