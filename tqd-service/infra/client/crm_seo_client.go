package client

import (
	"common"
	_middleware "common/middleware"
	"context"
	"errors"
	"fmt"
	"log"
	crmpb "pb/types/crm"
	"strings"
	"time"
	"tqd/env"
	"tqd/internal/dto"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/interface/provider"

	"google.golang.org/grpc"
)

const seoRefTypePlanningProject uint32 = 400

type CrmSeoClient struct {
	client crmpb.SeoDomainServiceClient
}

// @bind: tqd/internal/interface/provider.CrmSeoProvider
func NewCrmSeoClient() provider.CrmSeoProvider {
	c := &CrmSeoClient{}
	name := "crm"
	serviceURL := common.GetPrefixProtobufUrl(env.ENV_RUNTIME, name)

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:8206", serviceURL),
		grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(5*time.Second),
		grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(_middleware.StreamClientInterceptor),
	)
	if err != nil {
		log.Printf("failed to dial crm seo service: %v", err)
		return c
	}

	monitorCtx, cancel := context.WithCancel(context.Background())
	_ = cancel
	go _middleware.MonitorConnection(monitorCtx, name, conn)

	c.client = crmpb.NewSeoDomainServiceClient(conn)
	return c
}

func (c *CrmSeoClient) GenerateSeoFromParcelIDs(ctx context.Context, parcelIDs []uint64) (*dto.AdminGenParcelSeoResult, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("crm seo client chưa sẵn sàng")
	}
	if len(parcelIDs) == 0 {
		return nil, errors.New("parcelIds không được để trống")
	}

	resp, err := c.client.GenerateSeoFromParcels(ctx, &crmpb.GenerateSeoFromParcelsRequest{
		ParcelIds: parcelIDs,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &dto.AdminGenParcelSeoResult{}, nil
	}

	result := &dto.AdminGenParcelSeoResult{
		Total:   resp.Total,
		Created: resp.Created,
		Skipped: resp.Skipped,
		Failed:  resp.Failed,
		Errors:  resp.Errors,
		Items:   make([]dto.AdminGenParcelSeoItem, 0, len(resp.Data)),
	}
	for _, item := range resp.Data {
		if item == nil {
			continue
		}
		var refID uint64
		if item.RefId != nil {
			refID = *item.RefId
		}
		result.Items = append(result.Items, dto.AdminGenParcelSeoItem{
			ParcelID:     refID,
			SeoID:        item.Id,
			CanonicalURL: item.CanonicalUrl,
			Slug:         item.Slug,
		})
	}
	return result, nil
}

func (c *CrmSeoClient) GenerateSeoFromLayerIDs(ctx context.Context, layerIDs []uint64) (*dto.AdminGenLayerSeoResult, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("crm seo client chưa sẵn sàng")
	}
	if len(layerIDs) == 0 {
		return nil, errors.New("layerIds không được để trống")
	}

	resp, err := c.client.GenerateSeoFromLayers(ctx, &crmpb.GenerateSeoFromLayersRequest{
		LayerIds: layerIDs,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &dto.AdminGenLayerSeoResult{}, nil
	}

	result := &dto.AdminGenLayerSeoResult{
		Total:   resp.Total,
		Created: resp.Created,
		Skipped: resp.Skipped,
		Failed:  resp.Failed,
		Errors:  resp.Errors,
		Items:   make([]dto.AdminGenLayerSeoItem, 0, len(resp.Data)),
	}
	for _, item := range resp.Data {
		if item == nil {
			continue
		}
		var refID uint64
		if item.RefId != nil {
			refID = *item.RefId
		}
		result.Items = append(result.Items, dto.AdminGenLayerSeoItem{
			LayerID:      refID,
			SeoID:        item.Id,
			CanonicalURL: item.CanonicalUrl,
			Slug:         item.Slug,
		})
	}
	return result, nil
}

func (c *CrmSeoClient) GenerateSeoFromPlanningProjectIDs(ctx context.Context, planningProjectIDs []uint64) (*qh_dto.AdminGenPlanningProjectSeoResult, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("crm seo client chưa sẵn sàng")
	}
	if len(planningProjectIDs) == 0 {
		return nil, errors.New("planningProjectIds không được để trống")
	}

	resp, err := c.client.GenerateSeoFromPlanningProjects(ctx, &crmpb.GenerateSeoFromPlanningProjectsRequest{
		PlanningProjectIds: planningProjectIDs,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &qh_dto.AdminGenPlanningProjectSeoResult{}, nil
	}

	result := &qh_dto.AdminGenPlanningProjectSeoResult{
		Total:   resp.Total,
		Created: resp.Created,
		Skipped: resp.Skipped,
		Failed:  resp.Failed,
		Errors:  resp.Errors,
		Items:   make([]qh_dto.AdminGenPlanningProjectSeoItem, 0, len(resp.Data)),
	}
	for _, item := range resp.Data {
		if item == nil {
			continue
		}
		var refID uint64
		if item.RefId != nil {
			refID = *item.RefId
		}
		result.Items = append(result.Items, qh_dto.AdminGenPlanningProjectSeoItem{
			PlanningProjectID: refID,
			SeoID:             item.Id,
			CanonicalURL:      item.CanonicalUrl,
			Slug:              item.Slug,
		})
	}
	return result, nil
}

func (c *CrmSeoClient) GenerateSeoFromAdministrativeUnitIdentities(ctx context.Context, compositeIDs []uint64) (*dto.AdminGenAdministrativeUnitSeoResult, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("crm seo client chưa sẵn sàng")
	}
	if len(compositeIDs) == 0 {
		return nil, errors.New("compositeIds không được để trống")
	}

	resp, err := c.client.GenerateSeoFromAdministrativeUnits(ctx, &crmpb.GenerateSeoFromAdministrativeUnitsRequest{
		CompositeIds: compositeIDs,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &dto.AdminGenAdministrativeUnitSeoResult{}, nil
	}

	result := &dto.AdminGenAdministrativeUnitSeoResult{
		Total:   resp.Total,
		Created: resp.Created,
		Skipped: resp.Skipped,
		Failed:  resp.Failed,
		Errors:  resp.Errors,
		Items:   make([]dto.AdminGenAdministrativeUnitSeoItem, 0, len(resp.Data)),
	}
	for _, item := range resp.Data {
		if item == nil || item.RefId == nil {
			continue
		}
		identity := *item.RefId
		result.Items = append(result.Items, dto.AdminGenAdministrativeUnitSeoItem{
			Identity:     identity,
			SeoID:        item.Id,
			CanonicalURL: item.CanonicalUrl,
			Slug:         item.Slug,
		})
	}
	return result, nil
}

func (c *CrmSeoClient) CountSeoByPlanningProject(ctx context.Context, planningProjectID uint64) (int64, error) {
	if c == nil || c.client == nil {
		return 0, errors.New("crm seo client chưa sẵn sàng")
	}
	if planningProjectID == 0 {
		return 0, errors.New("planningProjectId là bắt buộc")
	}

	resp, err := c.client.CountSeoRefs(ctx, &crmpb.CountSeoRefsRequest{
		RefType: seoRefTypePlanningProject,
		RefId:   planningProjectID,
	})
	if err != nil {
		return 0, err
	}
	if resp == nil {
		return 0, nil
	}
	return resp.Total, nil
}

func (c *CrmSeoClient) CreateOrGetDefaultSeoByRef(ctx context.Context, refType uint32, refID uint64) (*dto.CreateDefaultSeoByRefResult, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("crm seo client chưa sẵn sàng")
	}
	if refType == 0 {
		return nil, errors.New("refType là bắt buộc")
	}
	if refID == 0 {
		return nil, errors.New("refId là bắt buộc")
	}

	existing, err := c.client.GetSeoRefs(ctx, &crmpb.GetSeoRefsRequest{
		RefType: refType,
		RefId:   refID,
		Page:    1,
		Size:    1,
	})
	if err != nil {
		return nil, err
	}
	if existing != nil && len(existing.Data) > 0 && existing.Data[0] != nil {
		item := existing.Data[0]
		return &dto.CreateDefaultSeoByRefResult{
			Created:      false,
			SeoID:        item.Id,
			RefType:      item.RefType,
			RefID:        item.GetRefId(),
			Slug:         item.Slug,
			OriginURL:    item.OriginUrl,
			CanonicalURL: item.CanonicalUrl,
			Title:        item.Title,
			Description:  item.Description,
			PageStatus:   item.PageStatus,
			RenderStatus: item.RenderStatus,
		}, nil
	}

	snapshot, err := c.client.ResolveSeoRefValue(ctx, &crmpb.ResolveSeoRefValueRequest{
		RefType: refType,
		RefId:   refID,
	})
	if err != nil {
		return nil, err
	}

	slug := strings.TrimSpace(snapshot.GetSuggestedSlug())
	if slug == "" {
		slug = fmt.Sprintf("public-%d-%d", refType, refID)
	}
	title := firstNonEmpty(snapshot.GetSuggestedTitle(), snapshot.GetSuggestedMetaTitle(), snapshot.GetLabel(), slug)
	description := firstNonEmpty(snapshot.GetSuggestedMetaDescription(), snapshot.GetDescription())
	originURL := firstNonEmpty(snapshot.GetPublicUrl(), snapshot.GetAdminUrl(), "/seo/"+slug)
	canonicalURL := firstNonEmpty(snapshot.GetPublicUrl(), snapshot.GetAdminUrl(), originURL)
	refSnapshotJSON := firstNonEmpty(snapshot.GetDataJson(), "{}")
	sourceStatus := firstNonEmpty(snapshot.GetSourceStatus(), "linked")

	created, err := c.client.CreateSeoDomain(ctx, &crmpb.CreateSeoDomainRequest{
		Slug:            slug,
		OriginUrl:       originURL,
		CanonicalUrl:    &canonicalURL,
		RefType:         refType,
		RefId:           &refID,
		Scope:           "public",
		Title:           title,
		Description:     description,
		Summary:         description,
		RefSource:       firstNonEmpty(snapshot.GetResolverKey(), snapshot.GetSourceService()),
		RefLabel:        snapshot.GetLabel(),
		RefUrl:          canonicalURL,
		SourceStatus:    sourceStatus,
		RefSnapshotJson: &refSnapshotJSON,
		RefHash:         snapshot.GetHash(),
		PageStatus:      "draft",
		RenderStatus:    "none",
		IsIndex:         true,
	})
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, errors.New("crm trả về seo rỗng")
	}

	return &dto.CreateDefaultSeoByRefResult{
		Created:      true,
		SeoID:        created.Id,
		RefType:      created.RefType,
		RefID:        created.GetRefId(),
		Slug:         created.Slug,
		OriginURL:    created.OriginUrl,
		CanonicalURL: created.CanonicalUrl,
		Title:        created.Title,
		Description:  created.Description,
		PageStatus:   created.PageStatus,
		RenderStatus: created.RenderStatus,
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
