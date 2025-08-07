/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package policy

import (
	"context"
	"teleport-plugin-slack-access-request/internal/policy/models"
)

type Service interface {
	CreateAccessPolicy(ctx context.Context, policy *models.AccessPolicy) (*models.AccessPolicy, error)
	DeleteAccessPolicyByAccessPolicyID(ctx context.Context, id int32) (*models.AccessPolicy, error)
	DeleteAccessPolicyByUserID(ctx context.Context, id int32) ([]*models.AccessPolicy, error)
	GetAccessPoliciesByInputChannelID(ctx context.Context, channelID string) ([]*models.AccessPolicy, error)
	UpdateAccessPolicyMessageTimestamp(ctx context.Context, ap *models.AccessPolicy) error
}

type Repository interface {
	CreateAccessPolicy(ctx context.Context, policy *models.AccessPolicy) (*models.AccessPolicy, error)
	DeleteAccessPolicyByAccessPolicyID(ctx context.Context, id int32) (*models.AccessPolicy, error)
	DeleteAccessPolicyByUserID(ctx context.Context, id int32) ([]*models.AccessPolicy, error)
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

func (s *service) DeleteAccessPolicyByUserID(ctx context.Context, id int32) ([]*models.AccessPolicy, error) {
	return s.repo.DeleteAccessPolicyByUserID(ctx, id)
}

func (s *service) GetAccessPoliciesByInputChannelID(ctx context.Context, channelID string) ([]*models.AccessPolicy, error) {
	return s.repo.GetAccessPoliciesByInputChannelID(ctx, channelID)
}

func (s *service) UpdateAccessPolicyMessageTimestamp(ctx context.Context, ap *models.AccessPolicy) error {
	return s.repo.UpdateAccessPolicyMessageTimestamp(ctx, ap)
}
