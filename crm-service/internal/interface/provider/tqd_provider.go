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
	GetLayerSeoSource(ctx context.Context, layerID uint64) (*dto.LayerSeoSource, error)
	UpdateLayerSeoID(ctx context.Context, layerID uint64, seoID uint64) error
	GetPlanningProjectSeoSource(ctx context.Context, planningProjectID uint64) (*dto.PlanningProjectSeoSource, error)
	UpdatePlanningProjectSeoID(ctx context.Context, planningProjectID uint64, seoID uint64) error
	GetAdministrativeUnitSeoSource(ctx context.Context, compositeID uint64) (*dto.AdministrativeUnitSeoSource, error)
	UpdateAdministrativeUnitSeoID(ctx context.Context, compositeID uint64, seoID uint64) error
}
