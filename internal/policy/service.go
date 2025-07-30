package policy

import (
	"context"
	"teleport-plugin-slack-access-request/internal/policy/models"
)

type Service interface {
	CreateAccessPolicy(ctx context.Context, policy *models.AccessPolicy) (*models.AccessPolicy, error)
}

type Repository interface {
	CreateAccessPolicy(ctx context.Context, policy *models.AccessPolicy) (*models.AccessPolicy, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateAccessPolicy(ctx context.Context, policy *models.AccessPolicy) (*models.AccessPolicy, error) {
	return s.repo.CreateAccessPolicy(ctx, policy)
}
