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

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/metric/telemetry"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/policy/models"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer(telemetry.PolicyService)

type Service interface {
	CreateAccessPolicy(ctx context.Context, policy *models.AccessPolicy) (*models.AccessPolicy, error)
	DeleteAccessPolicyByAccessPolicyID(ctx context.Context, id int32) error
	DeleteAccessPolicyByUserID(ctx context.Context, id int32) ([]*models.AccessPolicy, error)
	GetAccessPoliciesByAccessPolicyID(ctx context.Context, accessPolicyID int32) (*models.AccessPolicy, error)
	GetAccessPoliciesByInputChannelID(ctx context.Context, channelID string) ([]*models.AccessPolicy, error)
	UpdateAccessPolicyMsgTs(ctx context.Context, ap *models.AccessPolicy) error
}

type Repository interface {
	CreateAccessPolicy(ctx context.Context, policy *models.AccessPolicy) (*models.AccessPolicy, error)
	DeleteAccessPolicyByAccessPolicyID(ctx context.Context, id int32) error
	DeleteAccessPolicyByUserID(ctx context.Context, id int32) ([]*models.AccessPolicy, error)
	GetAccessPoliciesByAccessPolicyID(ctx context.Context, accessPolicyID int32) (*models.AccessPolicy, error)
	GetAccessPoliciesByInputChannelID(ctx context.Context, channelID string) ([]*models.AccessPolicy, error)
	UpdateAccessPolicyMsgTs(ctx context.Context, ap *models.AccessPolicy) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateAccessPolicy(ctx context.Context, policy *models.AccessPolicy) (*models.AccessPolicy, error) {
	ctx, span := tracer.Start(ctx, "CreateAccessPolicy")
	defer span.End()
	return s.repo.CreateAccessPolicy(ctx, policy)
}

func (s *service) DeleteAccessPolicyByAccessPolicyID(ctx context.Context, id int32) error {
	ctx, span := tracer.Start(ctx, "DeleteAccessPolicyByAccessPolicyID",
		trace.WithAttributes(
			attribute.Int64("accessPolicyID", int64(id)),
		),
	)
	defer span.End()
	return s.repo.DeleteAccessPolicyByAccessPolicyID(ctx, id)
}

func (s *service) DeleteAccessPolicyByUserID(ctx context.Context, id int32) ([]*models.AccessPolicy, error) {
	ctx, span := tracer.Start(ctx, "DeleteAccessPolicyByUserID",
		trace.WithAttributes(
			attribute.Int64("userID", int64(id)),
		),
	)
	defer span.End()
	return s.repo.DeleteAccessPolicyByUserID(ctx, id)
}

func (s *service) GetAccessPoliciesByAccessPolicyID(ctx context.Context, accessPolicyID int32) (*models.AccessPolicy, error) {
	ctx, span := tracer.Start(ctx, "GetAccessPoliciesByAccessPolicyID",
		trace.WithAttributes(
			attribute.Int64("accessPolicyID", int64(accessPolicyID)),
		),
	)
	defer span.End()
	return s.repo.GetAccessPoliciesByAccessPolicyID(ctx, accessPolicyID)
}

func (s *service) GetAccessPoliciesByInputChannelID(ctx context.Context, channelID string) ([]*models.AccessPolicy, error) {
	ctx, span := tracer.Start(ctx, "GetAccessPoliciesByInputChannelID",
		trace.WithAttributes(
			attribute.String("inputChannelID", channelID),
		),
	)
	defer span.End()
	return s.repo.GetAccessPoliciesByInputChannelID(ctx, channelID)
}

func (s *service) UpdateAccessPolicyMsgTs(ctx context.Context, ap *models.AccessPolicy) error {
	ctx, span := tracer.Start(ctx, "UpdateAccessPolicyMsgTs",
		trace.WithAttributes(
			attribute.Int64("accessPolicyID", int64(ap.AccessPolicyID)),
		),
	)
	defer span.End()
	return s.repo.UpdateAccessPolicyMsgTs(ctx, ap)
}
