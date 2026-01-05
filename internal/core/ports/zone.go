package ports

import (
	"context"

	"app/internal/core/domain"
)

type ZoneRepository interface {
	Save(ctx context.Context, zone *domain.Zone) error
	GetByID(ctx context.Context, id string) (*domain.Zone, error)
	List(ctx context.Context, search string) ([]*domain.Zone, error)
	Update(ctx context.Context, zone *domain.Zone) error
	Delete(ctx context.Context, id string) error
}

type ZoneService interface {
	CreateZone(ctx context.Context, req *domain.CreateZoneRequest) (*domain.Zone, error)
	GetZone(ctx context.Context, id string) (*domain.Zone, error)
	ListZones(ctx context.Context, search string) ([]*domain.Zone, error)
	UpdateZone(ctx context.Context, id string, req *domain.UpdateZoneRequest) (*domain.Zone, error)
	DeleteZone(ctx context.Context, id string) error
}
