package user

import (
	"context"
	"database/sql"
	"fmt"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/database/sqlc"
	slackmodels "teleport-plugin-slack-access-request/internal/slack/models"
	teleportmodels "teleport-plugin-slack-access-request/internal/teleport/models"
	usertmodels "teleport-plugin-slack-access-request/internal/user/models"
)

type PostgresRepository struct {
	q sqlc.Querier
}

func NewRepository(q sqlc.Querier) *PostgresRepository {
	return &PostgresRepository{q: q}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user *usertmodels.User) (*usertmodels.User, error) {
	baseEntity := database.MarkCreate()

	createUserParams := sqlc.CreateUserParams{
		TeleportUserID: user.TeleportUser.TeleportUserID,
		SlackUserID:    user.SlackUser.SlackUserID,
		UseYn:          baseEntity.UseYn,
		CreateCode:     baseEntity.CreateCode,
		CreateDate:     baseEntity.CreateDate,
		Version:        baseEntity.Version,
	}

	createdUser, err := r.q.CreateUser(ctx, createUserParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in DB: %w", err)
	}
	return &usertmodels.User{
		UserID: createdUser.UserID,
		TeleportUser: &teleportmodels.User{
			TeleportUserID: createdUser.TeleportUserID,
			Username:       user.TeleportUser.Username,
			UseYn:          user.TeleportUser.UseYn,
			CreateCode:     user.TeleportUser.CreateCode,
			CreateDate:     user.TeleportUser.CreateDate,
			Version:        user.TeleportUser.Version,
		},
		SlackUser: &slackmodels.User{
			SlackUserID: createdUser.SlackUserID,
			ID:          user.SlackUser.ID,
			Name:        user.SlackUser.Name,
			RealName:    user.SlackUser.RealName,
			Email:       user.SlackUser.Email,
			UseYn:       user.SlackUser.UseYn,
			CreateCode:  user.SlackUser.CreateCode,
			CreateDate:  user.SlackUser.CreateDate,
			Version:     user.SlackUser.Version,
		},
		UseYn:      createdUser.UseYn,
		CreateCode: createdUser.CreateCode,
		CreateDate: createdUser.CreateDate,
		Version:    createdUser.Version,
	}, nil
}

func (r *PostgresRepository) DeleteUser(ctx context.Context, user *usertmodels.User) (*usertmodels.User, error) {
	baseEntity := database.MarkCreate()

	deleteUserParams := sqlc.DeleteUserByTeleportAndSlackIDParams{
		TeleportUserID: user.TeleportUser.TeleportUserID,
		SlackUserID:    user.SlackUser.SlackUserID,
		DeleteCode:     sql.NullString{String: baseEntity.DeleteCode, Valid: baseEntity.DeleteCode != ""},
		DeleteDate:     sql.NullTime{Time: baseEntity.DeleteDate, Valid: !baseEntity.DeleteDate.IsZero()},
	}

	createdUser, err := r.q.DeleteUserByTeleportAndSlackID(ctx, deleteUserParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in DB: %w", err)
	}
	return &usertmodels.User{
		UserID: createdUser.UserID,
		TeleportUser: &teleportmodels.User{
			TeleportUserID: createdUser.TeleportUserID,
			Username:       user.TeleportUser.Username,
			UseYn:          user.TeleportUser.UseYn,
			CreateCode:     user.TeleportUser.CreateCode,
			CreateDate:     user.TeleportUser.CreateDate,
			Version:        user.TeleportUser.Version,
		},
		SlackUser: &slackmodels.User{
			SlackUserID: createdUser.SlackUserID,
			ID:          user.SlackUser.ID,
			Name:        user.SlackUser.Name,
			RealName:    user.SlackUser.RealName,
			Email:       user.SlackUser.Email,
			UseYn:       user.SlackUser.UseYn,
			CreateCode:  user.SlackUser.CreateCode,
			CreateDate:  user.SlackUser.CreateDate,
			Version:     user.SlackUser.Version,
		},
		UseYn:      createdUser.UseYn,
		CreateCode: createdUser.CreateCode,
		CreateDate: createdUser.CreateDate,
		Version:    createdUser.Version,
	}, nil
}

func (r *PostgresRepository) GetUserBySlackUserID(ctx context.Context, id int32) (*usertmodels.User, error) {
	row, err := r.q.GetUserBySlackUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by slack user id %d: %w", id, err)
	}
	return &usertmodels.User{
		UserID:       row.UserID,
		TeleportUser: &teleportmodels.User{TeleportUserID: row.TeleportUserID},
		SlackUser:    &slackmodels.User{SlackUserID: row.SlackUserID},
		UseYn:        row.UseYn,
		CreateCode:   row.CreateCode,
		CreateDate:   row.CreateDate,
		UpdateCode:   row.UpdateCode.String,
		UpdateDate:   row.UpdateDate.Time,
		DeleteCode:   row.DeleteCode.String,
		DeleteDate:   row.DeleteDate.Time,
		Version:      row.Version,
	}, nil
}
