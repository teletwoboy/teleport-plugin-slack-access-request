package outbox

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/database/sqlc"
	"teleport-plugin-slack-access-request/internal/outbox/model"
)

type PostgresRepository struct {
	q sqlc.Querier
}

func NewRepository(q sqlc.Querier) *PostgresRepository {
	return &PostgresRepository{q: q}
}

func (r *PostgresRepository) CreateOutbox(ctx context.Context, o *model.Outbox) error {
	baseEntity := database.MarkCreate()

	createOutboxParams := sqlc.CreateOutboxParams{
		EventType:   o.EventType,
		AggregateID: o.AggregateID,
		Payload:     o.Payload,
		Status:      o.Status,
		UseYn:       baseEntity.UseYn,
		CreateCode:  baseEntity.CreateCode,
		CreateDate:  baseEntity.CreateDate,
		Version:     baseEntity.Version,
	}

	err := r.q.CreateOutbox(ctx, createOutboxParams)
	if err != nil {
		return fmt.Errorf("failed to create outbox: %w", err)
	}
	return nil
}

func (r *PostgresRepository) ClaimNextOutbox(ctx context.Context) (*model.Outbox, error) {
	ob, err := r.q.ClaimNextOutbox(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to claim next outbox: %w", err)
	}
	return &model.Outbox{
		OutboxID:    ob.OutboxID,
		EventType:   ob.EventType,
		AggregateID: ob.AggregateID,
		Payload:     ob.Payload,
		Status:      ob.Status,
		Attempts:    ob.Attempts,
		ApiAttempts: ob.ApiAttempts,
		DBAttempts:  ob.DbAttempts,
		NextTryAt:   ob.NextTryAt.Time,
		LastError:   ob.LastError.String,
		UseYn:       ob.UseYn,
		CreateCode:  ob.CreateCode,
		CreateDate:  ob.CreateDate,
		UpdateCode:  ob.UpdateCode.String,
		UpdateDate:  ob.UpdateDate.Time,
		DeleteCode:  ob.DeleteCode.String,
		DeleteDate:  ob.DeleteDate.Time,
		Version:     ob.Version,
	}, nil
}

func (r *PostgresRepository) MarkFailed(ctx context.Context, o *model.Outbox, err error) error {
	markFailedParams := sqlc.MarkFailedParams{
		OutboxID:  o.OutboxID,
		LastError: sql.NullString{String: err.Error(), Valid: err.Error() != ""},
	}
	if err := r.q.MarkFailed(ctx, markFailedParams); err != nil {
		return fmt.Errorf("failed to mark failed outbox: %w, outbox_id: %d", err, o.OutboxID)
	}
	return nil
}
