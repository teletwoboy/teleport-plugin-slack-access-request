package outbox

import (
	"context"
	"teleport-plugin-slack-access-request/internal/outbox/model"
)

type Service interface {
	CreateOutbox(ctx context.Context, o *model.Outbox) error
	ClaimNextOutbox(ctx context.Context) (*model.Outbox, error)
	MarkFailed(ctx context.Context, o *model.Outbox, err error) error
}

type Repository interface {
	CreateOutbox(ctx context.Context, o *model.Outbox) error
	ClaimNextOutbox(ctx context.Context) (*model.Outbox, error)
	MarkFailed(ctx context.Context, o *model.Outbox, err error) error
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
	return s.repo.ClaimNextOutbox(ctx)
}

func (s *service) MarkFailed(ctx context.Context, o *model.Outbox, err error) error {
	return s.repo.MarkFailed(ctx, o, err)
}
