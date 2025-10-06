package outbox

import (
	"context"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer(telemetry.OutboxService)

type Service interface {
	CreateOutbox(ctx context.Context, o *model.Outbox) (*model.Outbox, error)
	ClaimOutboxByOutboxID(ctx context.Context, outboxID int32) (*model.Outbox, error)
	ClaimOutboxes(ctx context.Context, limit int32) ([]*model.Outbox, error)
	MarkDeadBatch(ctx context.Context, limit int32) ([]*model.Outbox, error)
	MarkDead(ctx context.Context, o *model.Outbox) error
	MarkDone(ctx context.Context, o *model.Outbox) error
	MarkFailed(ctx context.Context, o *model.Outbox, err error) error
	Notify(ctx context.Context, ob *model.OutboxNotification) error
}

type Repository interface {
	CreateOutbox(ctx context.Context, o *model.Outbox) (*model.Outbox, error)
	ClaimOutboxByOutboxID(ctx context.Context, id int32, set string, secs float64) (*model.Outbox, error)
	ClaimOutboxes(ctx context.Context, status string, limit int32, set string, secs float64) ([]*model.Outbox, error)
	MarkDeadBatch(ctx context.Context, statuses []string, limit int32, set string) ([]*model.Outbox, error)
	MarkStatus(ctx context.Context, o *model.Outbox, status string, set string) error
	MarkStatusAndNextTry(ctx context.Context, o *model.Outbox, status string, set string, err error, secs float64) error
	Notify(ctx context.Context, ob *model.OutboxNotification) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateOutbox(ctx context.Context, o *model.Outbox) (*model.Outbox, error) {
	return s.repo.CreateOutbox(ctx, o)
}

func (s *service) ClaimOutboxByOutboxID(ctx context.Context, outboxID int32) (*model.Outbox, error) {
	ctx, span := tracer.Start(ctx, "ClaimOutboxByOutboxID",
		trace.WithAttributes(
			attribute.Int64("outbox_id", int64(outboxID)),
		),
	)
	defer span.End()

	set := constant.Processing
	secs := constant.NextTrySecs
	return s.repo.ClaimOutboxByOutboxID(ctx, outboxID, set, secs)
}

func (s *service) ClaimOutboxes(ctx context.Context, limit int32) ([]*model.Outbox, error) {
	ctx, span := tracer.Start(ctx, "ClaimOutboxes")
	defer span.End()

	status := constant.Pending
	set := constant.Processing
	secs := constant.NextTrySecs
	return s.repo.ClaimOutboxes(ctx, status, limit, set, secs)
}

func (s *service) MarkDeadBatch(ctx context.Context, limit int32) ([]*model.Outbox, error) {
	ctx, span := tracer.Start(ctx, "MarkDeadBatch")
	defer span.End()

	statuses := []string{
		constant.Processing,
		constant.Failed,
	}
	set := constant.Dead
	return s.repo.MarkDeadBatch(ctx, statuses, limit, set)
}

func (s *service) MarkDead(ctx context.Context, o *model.Outbox) error {
	ctx, span := tracer.Start(ctx, "MarkDead",
		trace.WithAttributes(
			attribute.Int64("outbox_id", int64(o.OutboxID)),
		),
	)
	defer span.End()

	status := constant.Processing
	set := constant.Dead
	return s.repo.MarkStatus(ctx, o, status, set)
}

func (s *service) MarkDone(ctx context.Context, o *model.Outbox) error {
	ctx, span := tracer.Start(ctx, "MarkDone",
		trace.WithAttributes(
			attribute.Int64("outbox_id", int64(o.OutboxID)),
		),
	)
	defer span.End()

	status := constant.Processing
	set := constant.Done
	return s.repo.MarkStatus(ctx, o, status, set)
}

func (s *service) MarkFailed(ctx context.Context, o *model.Outbox, err error) error {
	ctx, span := tracer.Start(ctx, "MarkFailed",
		trace.WithAttributes(
			attribute.Int64("outbox_id", int64(o.OutboxID)),
		),
	)
	defer span.End()

	status := constant.Processing
	set := constant.Failed
	secs := constant.NextTrySecs
	return s.repo.MarkStatusAndNextTry(ctx, o, status, set, err, secs)
}

func (s *service) Notify(ctx context.Context, ob *model.OutboxNotification) error {
	ctx, span := tracer.Start(ctx, "Notify")
	defer span.End()

	return s.repo.Notify(ctx, ob)
}
