package outbox

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"teleport-plugin-slack-access-request/internal/database"
	"teleport-plugin-slack-access-request/internal/database/sqlc"
	"teleport-plugin-slack-access-request/internal/outbox/model"
	"time"
)

type PostgresRepository struct {
	q sqlc.Querier
}

func NewRepository(q sqlc.Querier) *PostgresRepository {
	return &PostgresRepository{q: q}
}

func (r *PostgresRepository) CreateOutbox(ctx context.Context, o *model.Outbox) (*model.Outbox, error) {
	baseEntity := database.MarkCreate()

	createOutboxParams := sqlc.CreateOutboxParams{
		EventType:     o.EventType,
		AggregateType: o.AggregateType,
		AggregateID:   o.AggregateID,
		Payload:       o.Payload,
		Status:        o.Status,
		UseYn:         baseEntity.UseYn,
		CreateCode:    baseEntity.CreateCode,
		CreateDate:    baseEntity.CreateDate,
		Version:       baseEntity.Version,
	}

	id, err := r.q.CreateOutbox(ctx, createOutboxParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create outbox: %w", err)
	}
	return &model.Outbox{OutboxID: id}, nil
}

func (r *PostgresRepository) ClaimOutboxByOutboxID(ctx context.Context, id int32, set string, secs float64) (*model.Outbox, error) {
	baseEntity := database.MarkUpdate()

	params := sqlc.ClaimOutboxByOutboxIDParams{
		OutboxID:   id,
		Status:     set,
		Secs:       secs,
		UpdateCode: sql.NullString{String: baseEntity.UpdateCode, Valid: baseEntity.UpdateCode != ""},
		UpdateDate: sql.NullTime{Time: time.Now(), Valid: !baseEntity.UpdateDate.IsZero()},
	}

	row, err := r.q.ClaimOutboxByOutboxID(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to claim outbox by outbox id: %w", err)
	}
	return &model.Outbox{
		OutboxID:      row.OutboxID,
		EventType:     row.EventType,
		AggregateType: row.AggregateType,
		AggregateID:   row.AggregateID,
		Payload:       row.Payload,
		Status:        row.Status,
		Attempts:      row.Attempts,
		NextTryAt:     row.NextTryAt.Time,
		LastError:     row.LastError.String,
		UseYn:         row.UseYn,
		CreateCode:    row.CreateCode,
		CreateDate:    row.CreateDate,
		UpdateCode:    row.UpdateCode.String,
		UpdateDate:    row.UpdateDate.Time,
		Version:       row.Version,
	}, nil
}

func (r *PostgresRepository) ClaimOutboxes(ctx context.Context, status string, limit int32, set string, secs float64) ([]*model.Outbox, error) {
	var result []*model.Outbox

	baseEntity := database.MarkUpdate()
	params := sqlc.ClaimOutboxesParams{
		Status:     status,
		Limit:      limit,
		Status_2:   set,
		Secs:       secs,
		UpdateCode: sql.NullString{String: baseEntity.UpdateCode, Valid: baseEntity.UpdateCode != ""},
		UpdateDate: sql.NullTime{Time: time.Now(), Valid: !baseEntity.UpdateDate.IsZero()},
	}
	rows, err := r.q.ClaimOutboxes(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to claim outboxes: %w", err)
	}

	for _, row := range rows {
		copiedRow := row
		result = append(result, &model.Outbox{
			OutboxID:      copiedRow.OutboxID,
			EventType:     copiedRow.EventType,
			AggregateType: copiedRow.AggregateType,
			AggregateID:   copiedRow.AggregateID,
			Payload:       copiedRow.Payload,
			Status:        copiedRow.Status,
			Attempts:      copiedRow.Attempts,
			NextTryAt:     copiedRow.NextTryAt.Time,
			LastError:     copiedRow.LastError.String,
			UseYn:         copiedRow.UseYn,
			CreateCode:    copiedRow.CreateCode,
			CreateDate:    copiedRow.CreateDate,
			UpdateCode:    copiedRow.UpdateCode.String,
			UpdateDate:    copiedRow.UpdateDate.Time,
			Version:       copiedRow.Version,
		})
	}
	return result, nil
}

func (r *PostgresRepository) MarkDeadBatch(ctx context.Context, statuses []string, limit int32, set string) ([]*model.Outbox, error) {
	var result []*model.Outbox

	baseEntity := database.MarkUpdate()
	params := sqlc.MarKDeadBatchParams{
		Column1:    statuses,
		Limit:      limit,
		Status:     set,
		UpdateCode: sql.NullString{String: baseEntity.UpdateCode, Valid: baseEntity.UpdateCode != ""},
		UpdateDate: sql.NullTime{Time: time.Now(), Valid: !baseEntity.UpdateDate.IsZero()},
	}
	rows, err := r.q.MarKDeadBatch(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to mark outbox dead batch: %w", err)
	}
	for _, row := range rows {
		copiedRow := row
		result = append(result, &model.Outbox{
			OutboxID:  copiedRow.OutboxID,
			Attempts:  copiedRow.Attempts,
			LastError: copiedRow.LastError.String,
		})
	}
	return result, nil
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

func (r *PostgresRepository) Notify(ctx context.Context, ob *model.OutboxNotification) error {
	notifyParams := sqlc.NotifyParams{
		PgNotify:   ob.Channel,
		PgNotify_2: ob.Payload,
	}
	return r.q.Notify(ctx, notifyParams)
}
