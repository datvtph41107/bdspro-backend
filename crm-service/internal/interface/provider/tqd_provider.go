package provider

import (
	"context"
	"crm/internal/dto"
)

type TqdProvider interface {
	GetAdministrativeUnitProjection(ctx context.Context, identity string) (*dto.TqdAdministrativeUnitProjection, error)
	GetPlanningProjectProjection(ctx context.Context, identity string) (*dto.TqdPlanningProjectProjection, error)
	GetParcelSeoSource(ctx context.Context, parcelID uint64) (*dto.ParcelSeoSource, error)
	GetParcelSeoSourcesForGenerate(ctx context.Context, limit uint32) ([]dto.ParcelSeoSource, error)
	UpdateParcelSeoID(ctx context.Context, parcelID uint64, seoID uint64) error
}
