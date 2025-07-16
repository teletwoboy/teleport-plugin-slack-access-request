package seedinit

import (
	"context"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context) error
	UpdateStaus(ctx context.Context) error
	GetStatus(ctx context.Context) (*SeedInit, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context) error {
	err := s.repo.Create(ctx)
	if err != nil {
		return fmt.Errorf("failed create seedinit: %w", err)
	}
	return nil
}

func (s *Service) UpdateStaus(ctx context.Context) error {
	err := s.repo.UpdateStaus(ctx)
	if err != nil {
		return fmt.Errorf("failed update staus: %w", err)
	}
	return nil
}

func (s *Service) GetStatus(ctx context.Context) (string, error) {
	seedInit, err := s.repo.GetStatus(ctx)
	if err != nil {
		return "", fmt.Errorf("failed get seedinit status: %w", err)
	}
	return seedInit.Status, nil
}
