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
	"tqd/env"
	"tqd/internal/dto"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/interface/provider"
	"time"

	"google.golang.org/grpc"
)

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

func (c *CrmSeoClient) GenerateSeoFromAdministrativeUnitIdentities(ctx context.Context, identities []string) (*dto.AdminGenAdministrativeUnitSeoResult, error) {
	if c == nil || c.client == nil {
		return nil, errors.New("crm seo client chưa sẵn sàng")
	}
	if len(identities) == 0 {
		return nil, errors.New("identities không được để trống")
	}

	resp, err := c.client.GenerateSeoFromAdministrativeUnits(ctx, &crmpb.GenerateSeoFromAdministrativeUnitsRequest{
		Identities: identities,
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
		if item == nil {
			continue
		}
		identity := strings.TrimSpace(item.RefUrl)
		result.Items = append(result.Items, dto.AdminGenAdministrativeUnitSeoItem{
			Identity:     identity,
			SeoID:        item.Id,
			CanonicalURL: item.CanonicalUrl,
			Slug:         item.Slug,
		})
	}
	return result, nil
}
