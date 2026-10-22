package provider

import (
	"context"
	"mdn4/modules/tqd/dto"
	qh_dto "mdn4/modules/tqd/dto/qh"
)

type CrmSeoProvider interface {
	GenerateSeoFromParcelIDs(ctx context.Context, parcelIDs []uint64) (*dto.AdminGenParcelSeoResult, error)
	GenerateSeoFromLayerIDs(ctx context.Context, layerIDs []uint64) (*dto.AdminGenLayerSeoResult, error)
	GenerateSeoFromPlanningProjectIDs(ctx context.Context, planningProjectIDs []uint64) (*qh_dto.AdminGenPlanningProjectSeoResult, error)
	GenerateSeoFromPlanningDocumentIDs(ctx context.Context, documentIDs []uint64) (*qh_dto.AdminGenPlanningDocumentSeoResult, error)
	GenerateSeoFromAdministrativeUnitIdentities(ctx context.Context, identities []uint64) (*dto.AdminGenAdministrativeUnitSeoResult, error)
	CountSeoByPlanningProject(ctx context.Context, planningProjectID uint64) (int64, error)
	CreateOrGetDefaultSeoByRef(ctx context.Context, refType uint32, refID uint64) (*dto.CreateDefaultSeoByRefResult, error)
}
