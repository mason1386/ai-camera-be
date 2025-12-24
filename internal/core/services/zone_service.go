package services

import (
	"context"

	"app/internal/core/domain"
	"app/internal/core/ports"
	"app/pkg/logger"
	"go.uber.org/zap"
)

type ZoneService struct {
	repo ports.ZoneRepository
}

func NewZoneService(repo ports.ZoneRepository) ports.ZoneService {
	return &ZoneService{
		repo: repo,
	}
}

func (s *ZoneService) CreateZone(ctx context.Context, req *domain.CreateZoneRequest) (*domain.Zone, error) {
	zone := &domain.Zone{
		Name:        req.Name,
		Description: req.Description,
	}

	err := s.repo.Save(ctx, zone)
	if err != nil {
		logger.Error("Failed to create zone", zap.Error(err))
		return nil, err
	}

	return zone, nil
}

func (s *ZoneService) GetZone(ctx context.Context, id string) (*domain.Zone, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ZoneService) ListZones(ctx context.Context) ([]*domain.Zone, error) {
	return s.repo.List(ctx)
}

func (s *ZoneService) UpdateZone(ctx context.Context, id string, req *domain.UpdateZoneRequest) (*domain.Zone, error) {
	zone, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if zone == nil {
		return nil, nil // Or custom error NotFound
	}

	// Update fields if provided (naive approach, usually check for empty string/nil)
	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.Description != "" {
		zone.Description = req.Description
	}

	if err := s.repo.Update(ctx, zone); err != nil {
		return nil, err
	}
	return zone, nil
}

func (s *ZoneService) DeleteZone(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
