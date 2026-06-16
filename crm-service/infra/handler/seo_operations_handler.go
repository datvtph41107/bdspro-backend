package handler

import (
	"context"
	"crm/internal/dto"

	crmpb "pb/types/crm"
)

func (h *SeoDomainHandler) GetSeoOperationsDashboard(ctx context.Context, req *crmpb.GetSeoOperationsDashboardRequest) (*crmpb.SeoOperationsDashboardResponse, error) {
	result, err := h.seoDomainUsecase.GetSeoOperationsDashboard(ctx, &dto.SeoOperationsDashboardRequest{
		Range: req.Range,
		Limit: req.Limit,
	})
	if err != nil {
		return nil, err
	}
	return seoOperationsDashboardToPb(result), nil
}

func seoOperationsDashboardToPb(in *dto.SeoOperationsDashboardResponse) *crmpb.SeoOperationsDashboardResponse {
	if in == nil {
		return &crmpb.SeoOperationsDashboardResponse{}
	}
	return &crmpb.SeoOperationsDashboardResponse{
		Kpis:          seoOperationKPIsToPb(in.KPIs),
		Trend:         seoOperationTrendToPb(in.Trend),
		Opportunities: seoOperationOpportunitiesToPb(in.Opportunities),
		TopQueries:    seoOperationTopQueriesToPb(in.TopQueries),
		SystemEvents:  seoOperationEventsToPb(in.SystemEvents),
		Requirements:  seoRequirementGroupsToPb(in.Requirements),
		RenderQueue:   seoRenderQueueToPb(in.RenderQueue),
		IndexingRows:  seoIndexingRowsToPb(in.IndexingRows),
		SourceHealth:  seoSourceHealthToPb(in.SourceHealth),
	}
}

func seoOperationKPIsToPb(items []dto.SeoOperationKPI) []*crmpb.SeoOperationKpiResponse {
	out := make([]*crmpb.SeoOperationKpiResponse, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoOperationKpiResponse{Label: item.Label, Value: item.Value, Change: item.Change, Tone: item.Tone})
	}
	return out
}

func seoOperationTrendToPb(items []dto.SeoOperationTrendPoint) []*crmpb.SeoOperationTrendPointResponse {
	out := make([]*crmpb.SeoOperationTrendPointResponse, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoOperationTrendPointResponse{Label: item.Label, Clicks: item.Clicks, Impressions: item.Impressions})
	}
	return out
}

func seoOperationOpportunitiesToPb(items []dto.SeoOperationOpportunity) []*crmpb.SeoOperationOpportunityResponse {
	out := make([]*crmpb.SeoOperationOpportunityResponse, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoOperationOpportunityResponse{
			Id: item.ID, Title: item.Title, Description: item.Description, Issue: item.Issue,
			Impact: item.Impact, PageCount: item.PageCount, Owner: item.Owner, Tone: item.Tone,
		})
	}
	return out
}

func seoOperationTopQueriesToPb(items []dto.SeoOperationTopQuery) []*crmpb.SeoOperationTopQueryResponse {
	out := make([]*crmpb.SeoOperationTopQueryResponse, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoOperationTopQueryResponse{
			Query: item.Query, Clicks: item.Clicks, Impressions: item.Impressions, Ctr: item.CTR, AvgPosition: item.AvgPosition,
		})
	}
	return out
}

func seoOperationEventsToPb(items []dto.SeoOperationEvent) []*crmpb.SeoOperationEventResponse {
	out := make([]*crmpb.SeoOperationEventResponse, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoOperationEventResponse{Id: item.ID, Title: item.Title, Description: item.Description, CreatedAt: item.CreatedAt, Tone: item.Tone})
	}
	return out
}

func seoRequirementGroupsToPb(items []dto.SeoRequirementGroup) []*crmpb.SeoRequirementGroupResponse {
	out := make([]*crmpb.SeoRequirementGroupResponse, 0, len(items))
	for _, item := range items {
		group := &crmpb.SeoRequirementGroupResponse{Id: item.ID, Title: item.Title, Description: item.Description}
		for _, req := range item.Items {
			group.Items = append(group.Items, &crmpb.SeoRequirementItemResponse{
				Id: req.ID, Title: req.Title, Description: req.Description, Status: req.Status, Tone: req.Tone,
				Owner: req.Owner, Coverage: req.Coverage, Hint: req.Hint,
			})
		}
		out = append(out, group)
	}
	return out
}

func seoRenderQueueToPb(items []dto.SeoRenderQueueItem) []*crmpb.SeoRenderQueueItemResponse {
	out := make([]*crmpb.SeoRenderQueueItemResponse, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoRenderQueueItemResponse{
			Id: item.ID, SeoDomainId: item.SeoDomainID, Title: item.Title, Type: item.Type,
			Status: item.Status, Worker: item.Worker, Duration: item.Duration, Priority: item.Priority,
		})
	}
	return out
}

func seoIndexingRowsToPb(items []dto.SeoIndexingRow) []*crmpb.SeoIndexingRowResponse {
	out := make([]*crmpb.SeoIndexingRowResponse, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoIndexingRowResponse{
			Url: item.URL, Canonical: item.Canonical, Robots: item.Robots, Sitemap: item.Sitemap, State: item.State, LastPing: item.LastPing,
		})
	}
	return out
}

func seoSourceHealthToPb(items []dto.SeoSourceHealthItem) []*crmpb.SeoSourceHealthResponse {
	out := make([]*crmpb.SeoSourceHealthResponse, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoSourceHealthResponse{
			Source: item.Source, Label: item.Label, Service: item.Service, Status: item.Status,
			TotalPages: item.TotalPages, NeedAttention: item.NeedAttention,
		})
	}
	return out
}
