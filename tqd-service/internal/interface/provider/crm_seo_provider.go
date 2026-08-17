package provider

import (
	"context"
	"tqd/internal/dto"
	qh_dto "tqd/internal/dto/qh"
)

type CrmSeoProvider interface {
	GenerateSeoFromParcelIDs(ctx context.Context, parcelIDs []uint64) (*dto.AdminGenParcelSeoResult, error)
	GenerateSeoFromLayerIDs(ctx context.Context, layerIDs []uint64) (*dto.AdminGenLayerSeoResult, error)
	GenerateSeoFromPlanningProjectIDs(ctx context.Context, planningProjectIDs []uint64) (*qh_dto.AdminGenPlanningProjectSeoResult, error)
	GenerateSeoFromAdministrativeUnitIdentities(ctx context.Context, identities []uint64) (*dto.AdminGenAdministrativeUnitSeoResult, error)
	CountSeoByPlanningProject(ctx context.Context, planningProjectID uint64) (int64, error)
	CreateOrGetDefaultSeoByRef(ctx context.Context, refType uint32, refID uint64) (*dto.CreateDefaultSeoByRefResult, error)
}
