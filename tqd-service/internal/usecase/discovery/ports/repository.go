package ports

import (
	"context"
	"tqd/internal/domain/discovery/model"
)

type Repository interface {
	IdentifyParcels(ctx context.Context, req domain.IdentifyRequest) ([]domain.EntityCandidate, error)
	IdentifyRegions(ctx context.Context, req domain.IdentifyRequest) ([]domain.EntityCandidate, error)
	IdentifyAdministrativeUnits(ctx context.Context, req domain.IdentifyRequest) ([]domain.EntityCandidate, error)
	IdentifyPOIs(ctx context.Context, req domain.IdentifyRequest) ([]domain.EntityCandidate, error)

	SearchParcels(ctx context.Context, req domain.SearchRequest) ([]domain.EntityCandidate, error)
	SearchRegions(ctx context.Context, req domain.SearchRequest) ([]domain.EntityCandidate, error)
	SearchPlanningProjects(ctx context.Context, req domain.SearchRequest) ([]domain.EntityCandidate, error)
	SearchAdministrativeUnits(ctx context.Context, req domain.SearchRequest) ([]domain.EntityCandidate, error)
	SearchPOIs(ctx context.Context, req domain.SearchRequest) ([]domain.EntityCandidate, error)

	GetEntity(ctx context.Context, ref domain.EntityRef) (*domain.EntityCandidate, error)
}
