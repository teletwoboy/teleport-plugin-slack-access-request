package outbox

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"teleport-plugin-slack-access-request/internal/metric/telemetry"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
)

var tracer = otel.Tracer(telemetry.OutboxService)

type Service interface {
	CreateOutbox(ctx context.Context, o *model.Outbox) error
	ClaimNextOutbox(ctx context.Context) (*model.Outbox, error)
	MarkDead(ctx context.Context, o *model.Outbox) error
	MarkDone(ctx context.Context, o *model.Outbox) error
	MarkFailed(ctx context.Context, o *model.Outbox, err error) error
}

type Repository interface {
	CreateOutbox(ctx context.Context, o *model.Outbox) error
	ClaimNextOutbox(ctx context.Context, set string, secs float64, statuses []string, limit int32) (*model.Outbox, error)
	MarkStatus(ctx context.Context, o *model.Outbox, status string, set string) error
	MarkStatusAndNextTry(ctx context.Context, o *model.Outbox, status string, set string, err error, secs float64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateOutbox(ctx context.Context, o *model.Outbox) error {
	return s.repo.CreateOutbox(ctx, o)
}

func (s *service) ClaimNextOutbox(ctx context.Context) (*model.Outbox, error) {
	set := constant.Processing
	secs := constant.NextTrySecs
	statuses := []string{
		constant.Pending,
		constant.Failed,
	}
	limit := constant.PollSize

	return s.repo.ClaimNextOutbox(ctx, set, secs, statuses, limit)
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
