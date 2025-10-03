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

func (r *PostgresRepository) ClaimNextOutbox(ctx context.Context, set string, secs float64, statuses []string, limit int32) (*model.Outbox, error) {
	baseEntity := database.MarkUpdate()

	claimNextOutboxParams := sqlc.ClaimNextOutboxParams{
		Status:     set,
		Secs:       secs,
		UpdateCode: sql.NullString{String: baseEntity.UpdateCode, Valid: baseEntity.UpdateCode != ""},
		UpdateDate: sql.NullTime{Time: baseEntity.UpdateDate, Valid: !baseEntity.UpdateDate.IsZero()},
		Column5:    statuses,
		Limit:      limit,
	}

	ob, err := r.q.ClaimNextOutbox(ctx, claimNextOutboxParams)
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

func (r *PostgresRepository) MarkStatus(ctx context.Context, o *model.Outbox, status, set string) error {
	baseEntity := database.MarkUpdate()

	markStatusParams := sqlc.MarkStatusParams{
		OutboxID:   o.OutboxID,
		Status:     status,
		Status_2:   set,
		UpdateCode: sql.NullString{String: baseEntity.UpdateCode, Valid: baseEntity.UpdateCode != ""},
		UpdateDate: sql.NullTime{Time: baseEntity.UpdateDate, Valid: !baseEntity.UpdateDate.IsZero()},
	}

	if err := r.q.MarkStatus(ctx, markStatusParams); err != nil {
		return fmt.Errorf("failed to mark outbox status: %w, outbox_id: %d", err, o.OutboxID)
	}
	return nil
}

func (r *PostgresRepository) MarkStatusAndNextTry(ctx context.Context, o *model.Outbox, status, set string, err error, secs float64) error {
	baseEntity := database.MarkUpdate()

	markStatusAndNextTryParams := sqlc.MarkStatusAndNextTryParams{
		OutboxID:   0,
		Status:     status,
		Status_2:   set,
		LastError:  sql.NullString{String: err.Error(), Valid: err.Error() != ""},
		Secs:       secs,
		UpdateCode: sql.NullString{String: baseEntity.UpdateCode, Valid: baseEntity.UpdateCode != ""},
		UpdateDate: sql.NullTime{Time: baseEntity.UpdateDate, Valid: !baseEntity.UpdateDate.IsZero()},
	}

	if err := r.q.MarkStatusAndNextTry(ctx, markStatusAndNextTryParams); err != nil {
		return fmt.Errorf("failed to mark outbox status and next try: %w, outbox_id: %d", err, o.OutboxID)
	}
	return nil
}
