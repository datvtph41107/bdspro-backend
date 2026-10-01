package provider

import (
	"context"
	"errors"
	"fmt"

	"userv4/modules/seo/dto"
	"userv4/modules/seo/entity"
	"userv4/modules/seo/services"
	tqddto "userv4/modules/tqd/dto"
	qh_dto "userv4/modules/tqd/dto/qh"
)

// LocalCrmSeoProvider implements tqdprovider.CrmSeoProvider against modules/seo
// instead of a remote CRM gRPC client.
type LocalCrmSeoProvider struct {
	svc services.SeoDomainService
}

func NewLocalCrmSeoProvider(svc services.SeoDomainService) *LocalCrmSeoProvider {
	return &LocalCrmSeoProvider{svc: svc}
}

func (p *LocalCrmSeoProvider) GenerateSeoFromParcelIDs(ctx context.Context, parcelIDs []uint64) (*tqddto.AdminGenParcelSeoResult, error) {
	batch, err := p.svc.GenerateBatch(ctx, entity.SEORefTypeParcel, "tqd_parcel", parcelIDs)
	if err != nil {
		return nil, err
	}
	out := &tqddto.AdminGenParcelSeoResult{
		Total:   batch.Total,
		Created: batch.Created,
		Skipped: batch.Skipped,
		Failed:  batch.Failed,
		Errors:  batch.Errors,
		Items:   make([]tqddto.AdminGenParcelSeoItem, 0, len(batch.Items)),
	}
	for _, item := range batch.Items {
		out.Items = append(out.Items, tqddto.AdminGenParcelSeoItem{
			ParcelID:     item.RefID,
			SeoID:        item.SeoID,
			CanonicalURL: item.CanonicalURL,
			Slug:         item.Slug,
			Created:      item.Created,
		})
	}
	return out, nil
}

func (p *LocalCrmSeoProvider) GenerateSeoFromLayerIDs(ctx context.Context, layerIDs []uint64) (*tqddto.AdminGenLayerSeoResult, error) {
	batch, err := p.svc.GenerateBatch(ctx, entity.SEORefTypeMap, "tqd_layer", layerIDs)
	if err != nil {
		return nil, err
	}
	out := &tqddto.AdminGenLayerSeoResult{
		Total:   batch.Total,
		Created: batch.Created,
		Skipped: batch.Skipped,
		Failed:  batch.Failed,
		Errors:  batch.Errors,
		Items:   make([]tqddto.AdminGenLayerSeoItem, 0, len(batch.Items)),
	}
	for _, item := range batch.Items {
		out.Items = append(out.Items, tqddto.AdminGenLayerSeoItem{
			LayerID:      item.RefID,
			SeoID:        item.SeoID,
			CanonicalURL: item.CanonicalURL,
			Slug:         item.Slug,
			Created:      item.Created,
		})
	}
	return out, nil
}

func (p *LocalCrmSeoProvider) GenerateSeoFromPlanningProjectIDs(ctx context.Context, planningProjectIDs []uint64) (*qh_dto.AdminGenPlanningProjectSeoResult, error) {
	batch, err := p.svc.GenerateBatch(ctx, entity.SEORefTypeProject, "tqd_planning_project", planningProjectIDs)
	if err != nil {
		return nil, err
	}
	out := &qh_dto.AdminGenPlanningProjectSeoResult{
		Total:   batch.Total,
		Created: batch.Created,
		Skipped: batch.Skipped,
		Failed:  batch.Failed,
		Errors:  batch.Errors,
		Items:   make([]qh_dto.AdminGenPlanningProjectSeoItem, 0, len(batch.Items)),
	}
	for _, item := range batch.Items {
		out.Items = append(out.Items, qh_dto.AdminGenPlanningProjectSeoItem{
			PlanningProjectID: item.RefID,
			SeoID:             item.SeoID,
			CanonicalURL:      item.CanonicalURL,
			Slug:              item.Slug,
			Created:           item.Created,
		})
	}
	return out, nil
}

func (p *LocalCrmSeoProvider) GenerateSeoFromPlanningDocumentIDs(ctx context.Context, documentIDs []uint64) (*qh_dto.AdminGenPlanningDocumentSeoResult, error) {
	batch, err := p.svc.GenerateBatch(ctx, entity.SEORefTypeDocument, "tqd_planning_document", documentIDs)
	if err != nil {
		return nil, err
	}
	out := &qh_dto.AdminGenPlanningDocumentSeoResult{
		Total:   batch.Total,
		Created: batch.Created,
		Skipped: batch.Skipped,
		Failed:  batch.Failed,
		Errors:  batch.Errors,
		Items:   make([]qh_dto.AdminGenPlanningDocumentSeoItem, 0, len(batch.Items)),
	}
	for _, item := range batch.Items {
		out.Items = append(out.Items, qh_dto.AdminGenPlanningDocumentSeoItem{
			PlanningDocumentID: item.RefID,
			SeoID:              item.SeoID,
			CanonicalURL:       item.CanonicalURL,
			Slug:               item.Slug,
			Created:            item.Created,
		})
	}
	return out, nil
}

func (p *LocalCrmSeoProvider) GenerateSeoFromAdministrativeUnitIdentities(ctx context.Context, identities []uint64) (*tqddto.AdminGenAdministrativeUnitSeoResult, error) {
	batch, err := p.svc.GenerateBatch(ctx, entity.SEORefTypeAdmUnit, "tqd_adm_unit", identities)
	if err != nil {
		return nil, err
	}
	out := &tqddto.AdminGenAdministrativeUnitSeoResult{
		Total:   batch.Total,
		Created: batch.Created,
		Skipped: batch.Skipped,
		Failed:  batch.Failed,
		Errors:  batch.Errors,
		Items:   make([]tqddto.AdminGenAdministrativeUnitSeoItem, 0, len(batch.Items)),
	}
	for _, item := range batch.Items {
		out.Items = append(out.Items, tqddto.AdminGenAdministrativeUnitSeoItem{
			Identity:     item.RefID,
			CompositeID:  item.RefID,
			SeoID:        item.SeoID,
			CanonicalURL: item.CanonicalURL,
			Slug:         item.Slug,
			Created:      item.Created,
		})
	}
	return out, nil
}

func (p *LocalCrmSeoProvider) CountSeoByPlanningProject(ctx context.Context, planningProjectID uint64) (int64, error) {
	if planningProjectID == 0 {
		return 0, errors.New("planningProjectId is required")
	}
	return p.svc.CountByRef(ctx, uint32(entity.SEORefTypeProject), planningProjectID)
}

func (p *LocalCrmSeoProvider) CreateOrGetDefaultSeoByRef(ctx context.Context, refType uint32, refID uint64) (*tqddto.CreateDefaultSeoByRefResult, error) {
	if refType == 0 {
		return nil, errors.New("refType is required")
	}
	if refID == 0 {
		return nil, errors.New("refId is required")
	}
	existing, err := p.svc.GetByRef(ctx, refType, refID)
	if err == nil && existing != nil {
		var existingRefID uint64
		if existing.RefID != nil {
			existingRefID = *existing.RefID
		}
		return &tqddto.CreateDefaultSeoByRefResult{
			Created:      false,
			SeoID:        existing.ID,
			RefType:      uint32(existing.RefType),
			RefID:        existingRefID,
			Slug:         existing.Slug,
			OriginURL:    existing.OriginURL,
			CanonicalURL: existing.CanonicalURL,
			Title:        existing.Title,
			Description:  existing.Description,
			PageStatus:   existing.PageStatus,
			RenderStatus: existing.RenderStatus,
		}, nil
	}
	seo, created, err := p.svc.GenerateFromRef(ctx, dto.GenerateFromRefRequest{
		RefType:   refType,
		RefID:     refID,
		RefSource: fmt.Sprintf("ref-%d", refType),
	})
	if err != nil {
		return nil, err
	}
	if seo == nil {
		return nil, errors.New("seo page is empty")
	}
	var outRefID uint64
	if seo.RefID != nil {
		outRefID = *seo.RefID
	}
	return &tqddto.CreateDefaultSeoByRefResult{
		Created:      created,
		SeoID:        seo.ID,
		RefType:      uint32(seo.RefType),
		RefID:        outRefID,
		Slug:         seo.Slug,
		OriginURL:    seo.OriginURL,
		CanonicalURL: seo.CanonicalURL,
		Title:        seo.Title,
		Description:  seo.Description,
		PageStatus:   seo.PageStatus,
		RenderStatus: seo.RenderStatus,
	}, nil
}
