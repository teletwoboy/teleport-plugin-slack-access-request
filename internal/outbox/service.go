package outbox

import (
	"context"
	"teleport-plugin-slack-access-request/internal/outbox/constant"
	"teleport-plugin-slack-access-request/internal/outbox/model"
)

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
	status := constant.Processing
	set := constant.Dead

	return s.repo.MarkStatus(ctx, o, status, set)
}

func (s *service) MarkDone(ctx context.Context, o *model.Outbox) error {
	status := constant.Processing
	set := constant.Done

	return s.repo.MarkStatus(ctx, o, status, set)
}

func (s *service) MarkFailed(ctx context.Context, o *model.Outbox, err error) error {
	status := constant.Processing
	set := constant.Failed
	secs := constant.NextTrySecs

	return s.repo.MarkStatusAndNextTry(ctx, o, status, set, err, secs)
}
