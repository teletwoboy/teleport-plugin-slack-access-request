package policy

import (
	"context"
	"teleport-plugin-slack-access-request/internal/policy/models"
)

type Service interface {
	CreateAccessPolicy(ctx context.Context, policy *models.AccessPolicy) (*models.AccessPolicy, error)
	DeleteAccessPolicyByAccessPolicyID(ctx context.Context, id int32) (*models.AccessPolicy, error)
	GetAccessPoliciesByInputChannelID(ctx context.Context, channelID string) ([]*models.AccessPolicy, error)
	UpdateAccessPolicyMessageTimestamp(ctx context.Context, ap *models.AccessPolicy) error
}

type Repository interface {
	CreateAccessPolicy(ctx context.Context, policy *models.AccessPolicy) (*models.AccessPolicy, error)
	DeleteAccessPolicyByAccessPolicyID(ctx context.Context, id int32) (*models.AccessPolicy, error)
	GetAccessPoliciesByInputChannelID(ctx context.Context, channelID string) ([]*models.AccessPolicy, error)
	UpdateAccessPolicyMessageTimestamp(ctx context.Context, ap *models.AccessPolicy) error
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

func (s *service) DeleteAccessPolicyByAccessPolicyID(ctx context.Context, id int32) (*models.AccessPolicy, error) {
	return s.repo.DeleteAccessPolicyByAccessPolicyID(ctx, id)
}

func (s *service) GetAccessPoliciesByInputChannelID(ctx context.Context, channelID string) ([]*models.AccessPolicy, error) {
	return s.repo.GetAccessPoliciesByInputChannelID(ctx, channelID)
}

func (s *service) UpdateAccessPolicyMessageTimestamp(ctx context.Context, ap *models.AccessPolicy) error {
	return s.repo.UpdateAccessPolicyMessageTimestamp(ctx, ap)
}
