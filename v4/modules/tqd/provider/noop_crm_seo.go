package provider

import (
	"context"
	"fmt"

	"mdn4/modules/tqd/dto"
	qh_dto "mdn4/modules/tqd/dto/qh"
)

// NoopCrmSeoProvider satisfies CrmSeoProvider without a CRM RPC client.
// Client map reads do not need SEO generation; admin SEO endpoints return a clear error.
type NoopCrmSeoProvider struct{}

// NewNoopCrmSeoProvider is a fallback; the HTTP graph binds seo.LocalCrmSeoProvider.
// @wire:skip
func NewNoopCrmSeoProvider() CrmSeoProvider {
	return NoopCrmSeoProvider{}
}

func (NoopCrmSeoProvider) GenerateSeoFromParcelIDs(context.Context, []uint64) (*dto.AdminGenParcelSeoResult, error) {
	return nil, fmt.Errorf("crm seo provider not configured")
}

func (NoopCrmSeoProvider) GenerateSeoFromLayerIDs(context.Context, []uint64) (*dto.AdminGenLayerSeoResult, error) {
	return nil, fmt.Errorf("crm seo provider not configured")
}

func (NoopCrmSeoProvider) GenerateSeoFromPlanningProjectIDs(context.Context, []uint64) (*qh_dto.AdminGenPlanningProjectSeoResult, error) {
	return nil, fmt.Errorf("crm seo provider not configured")
}

func (NoopCrmSeoProvider) GenerateSeoFromPlanningDocumentIDs(context.Context, []uint64) (*qh_dto.AdminGenPlanningDocumentSeoResult, error) {
	return nil, fmt.Errorf("crm seo provider not configured")
}

func (NoopCrmSeoProvider) GenerateSeoFromAdministrativeUnitIdentities(context.Context, []uint64) (*dto.AdminGenAdministrativeUnitSeoResult, error) {
	return nil, fmt.Errorf("crm seo provider not configured")
}

func (NoopCrmSeoProvider) CountSeoByPlanningProject(context.Context, uint64) (int64, error) {
	return 0, nil
}

func (NoopCrmSeoProvider) CreateOrGetDefaultSeoByRef(context.Context, uint32, uint64) (*dto.CreateDefaultSeoByRefResult, error) {
	return nil, fmt.Errorf("crm seo provider not configured")
}
