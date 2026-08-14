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
	GenerateSeoFromAdministrativeUnitIdentities(ctx context.Context, identities []string) (*dto.AdminGenAdministrativeUnitSeoResult, error)
}
