package handler

import (
	_dto "common/domain/dto"
	"context"
	"crm/internal/dto"

	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
)

// GetSeoMeasurementOverview returns the measurement/feedback dashboard.
//
// This is the admin-facing overview for the SEO Measurement & Feedback Loop.
// It combines search performance, index coverage, traffic, issues and top pages.
func (h *SeoDomainHandler) GetSeoMeasurementOverview(ctx context.Context, req *crmpb.GetSeoMeasurementOverviewRequest) (*crmpb.SeoMeasurementOverviewResponse, error) {
	if req == nil {
		req = &crmpb.GetSeoMeasurementOverviewRequest{}
	}
	result, err := h.seoMeasurementUsecase.GetOverview(ctx, &dto.SeoMeasurementRangeRequest{Range: req.GetRange(), Limit: req.GetLimit()})
	if err != nil {
		return nil, err
	}
	return seoMeasurementOverviewToPb(result), nil
}

func (h *SeoDomainHandler) GenerateSeoMockMeasurement(ctx context.Context, req *crmpb.GenerateSeoMockMeasurementRequest) (*crmpb.GenerateSeoMockMeasurementResponse, error) {
	if req == nil {
		req = &crmpb.GenerateSeoMockMeasurementRequest{}
	}
	result, err := h.seoMeasurementUsecase.GenerateMockMeasurement(ctx, &dto.SeoMeasurementGenerateMockRequest{Range: req.GetRange(), Limit: req.GetLimit(), Overwrite: req.GetOverwrite()})
	if err != nil {
		return nil, err
	}
	return &crmpb.GenerateSeoMockMeasurementResponse{
		Pages:          result.Pages,
		SearchRows:     result.SearchRows,
		InspectionRows: result.InspectionRows,
		TrafficRows:    result.TrafficRows,
		RankingRows:    result.RankingRows,
		CrawlRows:      result.CrawlRows,
		KeywordTargets: result.KeywordTargets,
		DataSource:     result.DataSource,
		Messages:       result.Messages,
	}, nil
}

func (h *SeoDomainHandler) RunSeoFeedbackDetector(ctx context.Context, req *crmpb.RunSeoFeedbackDetectorRequest) (*crmpb.RunSeoFeedbackDetectorResponse, error) {
	if req == nil {
		req = &crmpb.RunSeoFeedbackDetectorRequest{}
	}
	result, err := h.seoMeasurementUsecase.RunFeedbackDetector(ctx, &dto.SeoMeasurementDetectorRequest{Range: req.GetRange(), Limit: req.GetLimit()})
	if err != nil {
		return nil, err
	}
	return &crmpb.RunSeoFeedbackDetectorResponse{Scanned: result.Scanned, Created: result.Created, Skipped: result.Skipped, Errors: result.Errors}, nil
}

func (h *SeoDomainHandler) GetSeoPagePerformanceList(ctx context.Context, req *crmpb.GetSeoPagePerformanceListRequest) (*crmpb.GetSeoPagePerformanceListResponse, error) {
	if req == nil {
		req = &crmpb.GetSeoPagePerformanceListRequest{}
	}
	items, total, err := h.seoMeasurementUsecase.ListPagePerformance(ctx, &dto.SeoPagePerformanceListRequest{
		Pagable:      _dto.Pagable{Page: req.GetPage(), Size: req.GetSize(), Sort: req.GetSort()},
		Keyword:      req.GetKeyword(),
		PageStatus:   req.GetPageStatus(),
		RenderStatus: req.GetRenderStatus(),
		IndexStatus:  req.GetIndexStatus(),
		IssueStatus:  req.GetIssueStatus(),
		Range:        req.GetRange(),
	})
	if err != nil {
		return nil, err
	}
	return &crmpb.GetSeoPagePerformanceListResponse{Data: seoPagePerformanceRowsToPb(items), Total: uint64(total)}, nil
}

func (h *SeoDomainHandler) GetSeoQueryPerformanceList(ctx context.Context, req *crmpb.GetSeoQueryPerformanceListRequest) (*crmpb.GetSeoQueryPerformanceListResponse, error) {
	if req == nil {
		req = &crmpb.GetSeoQueryPerformanceListRequest{}
	}
	items, total, err := h.seoMeasurementUsecase.ListQueryPerformance(ctx, &dto.SeoQueryPerformanceListRequest{
		Pagable:      _dto.Pagable{Page: req.GetPage(), Size: req.GetSize(), Sort: req.GetSort()},
		Query:        req.GetQuery(),
		CanonicalURL: req.GetCanonicalUrl(),
		Device:       req.GetDevice(),
		Country:      req.GetCountry(),
		Range:        req.GetRange(),
	})
	if err != nil {
		return nil, err
	}
	return &crmpb.GetSeoQueryPerformanceListResponse{Data: seoQueryPerformanceRowsToPb(items), Total: uint64(total)}, nil
}

func (h *SeoDomainHandler) GetSeoIndexCoverageList(ctx context.Context, req *crmpb.GetSeoIndexCoverageListRequest) (*crmpb.GetSeoIndexCoverageListResponse, error) {
	if req == nil {
		req = &crmpb.GetSeoIndexCoverageListRequest{}
	}
	items, total, err := h.seoMeasurementUsecase.ListIndexCoverage(ctx, &dto.SeoIndexCoverageListRequest{
		Pagable:      _dto.Pagable{Page: req.GetPage(), Size: req.GetSize(), Sort: req.GetSort()},
		CanonicalURL: req.GetCanonicalUrl(),
		IndexStatus:  req.GetIndexStatus(),
		Verdict:      req.GetVerdict(),
		Range:        req.GetRange(),
	})
	if err != nil {
		return nil, err
	}
	return &crmpb.GetSeoIndexCoverageListResponse{Data: seoIndexCoverageRowsToPb(items), Total: uint64(total)}, nil
}

func (h *SeoDomainHandler) GetSeoTrafficList(ctx context.Context, req *crmpb.GetSeoTrafficListRequest) (*crmpb.GetSeoTrafficListResponse, error) {
	if req == nil {
		req = &crmpb.GetSeoTrafficListRequest{}
	}
	items, total, err := h.seoMeasurementUsecase.ListTraffic(ctx, &dto.SeoTrafficListRequest{
		Pagable:      _dto.Pagable{Page: req.GetPage(), Size: req.GetSize(), Sort: req.GetSort()},
		CanonicalURL: req.GetCanonicalUrl(),
		Source:       req.GetSource(),
		Medium:       req.GetMedium(),
		Device:       req.GetDevice(),
		Country:      req.GetCountry(),
		Range:        req.GetRange(),
	})
	if err != nil {
		return nil, err
	}
	return &crmpb.GetSeoTrafficListResponse{Data: seoTrafficRowsToPb(items), Total: uint64(total)}, nil
}

func (h *SeoDomainHandler) GetSeoRankingList(ctx context.Context, req *crmpb.GetSeoRankingListRequest) (*crmpb.GetSeoRankingListResponse, error) {
	if req == nil {
		req = &crmpb.GetSeoRankingListRequest{}
	}
	items, total, err := h.seoMeasurementUsecase.ListRanking(ctx, &dto.SeoRankingListRequest{
		Pagable:      _dto.Pagable{Page: req.GetPage(), Size: req.GetSize(), Sort: req.GetSort()},
		Keyword:      req.GetKeyword(),
		CanonicalURL: req.GetCanonicalUrl(),
		Device:       req.GetDevice(),
		Range:        req.GetRange(),
	})
	if err != nil {
		return nil, err
	}
	return &crmpb.GetSeoRankingListResponse{Data: seoRankingRowsToPb(items), Total: uint64(total)}, nil
}

func (h *SeoDomainHandler) GetSeoCrawlEventList(ctx context.Context, req *crmpb.GetSeoCrawlEventListRequest) (*crmpb.GetSeoCrawlEventListResponse, error) {
	if req == nil {
		req = &crmpb.GetSeoCrawlEventListRequest{}
	}
	items, total, err := h.seoMeasurementUsecase.ListCrawlEvents(ctx, &dto.SeoCrawlEventListRequest{
		Pagable:      _dto.Pagable{Page: req.GetPage(), Size: req.GetSize(), Sort: req.GetSort()},
		CanonicalURL: req.GetCanonicalUrl(),
		BotName:      req.GetBotName(),
		StatusCode:   req.GetStatusCode(),
		Range:        req.GetRange(),
	})
	if err != nil {
		return nil, err
	}
	return &crmpb.GetSeoCrawlEventListResponse{Data: seoCrawlEventRowsToPb(items), Total: uint64(total)}, nil
}

func (h *SeoDomainHandler) GetSeoIssueList(ctx context.Context, req *crmpb.GetSeoIssueListRequest) (*crmpb.GetSeoIssueListResponse, error) {
	if req == nil {
		req = &crmpb.GetSeoIssueListRequest{}
	}
	items, total, err := h.seoMeasurementUsecase.ListIssues(ctx, &dto.SeoIssueListRequest{
		Pagable:      _dto.Pagable{Page: req.GetPage(), Size: req.GetSize(), Sort: req.GetSort()},
		CanonicalURL: req.GetCanonicalUrl(),
		IssueType:    req.GetIssueType(),
		Severity:     req.GetSeverity(),
		Status:       req.GetStatus(),
		Range:        req.GetRange(),
	})
	if err != nil {
		return nil, err
	}
	return &crmpb.GetSeoIssueListResponse{Data: seoIssueRowsToPb(items), Total: uint64(total)}, nil
}

func (h *SeoDomainHandler) ResolveSeoIssue(ctx context.Context, req *crmpb.ResolveSeoIssueRequest) (*sharepb.SubmitResponse, error) {
	if req == nil {
		req = &crmpb.ResolveSeoIssueRequest{}
	}
	if err := h.seoMeasurementUsecase.ResolveIssue(ctx, &dto.SeoResolveIssueRequest{ID: req.GetId(), Status: req.GetStatus(), Note: req.GetNote()}); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.GetId(), Message: "SEO issue updated"}, nil
}

func seoMeasurementOverviewToPb(item *dto.SeoMeasurementOverviewResponse) *crmpb.SeoMeasurementOverviewResponse {
	if item == nil {
		return &crmpb.SeoMeasurementOverviewResponse{}
	}
	return &crmpb.SeoMeasurementOverviewResponse{
		Kpis:          seoMeasurementKPIsToPb(item.KPIs),
		Trend:         seoMeasurementTrendToPb(item.Trend),
		IssueSummary:  seoIssueSummaryToPb(item.IssueSummary),
		TopPages:      seoPagePerformanceRowsToPb(item.TopPages),
		NeedAttention: seoIssueRowsToPb(item.NeedAttention),
	}
}

func seoMeasurementKPIsToPb(items []dto.SeoMeasurementKPI) []*crmpb.SeoMeasurementKPI {
	out := make([]*crmpb.SeoMeasurementKPI, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoMeasurementKPI{Key: item.Key, Label: item.Label, Value: item.Value, Unit: item.Unit, Delta: item.Delta, Severity: item.Severity})
	}
	return out
}

func seoMeasurementTrendToPb(items []dto.SeoMeasurementTrendPoint) []*crmpb.SeoMeasurementTrendPoint {
	out := make([]*crmpb.SeoMeasurementTrendPoint, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoMeasurementTrendPoint{Date: item.Date, Clicks: item.Clicks, Impressions: item.Impressions, Ctr: item.CTR, Position: item.Position, Sessions: item.Sessions, Conversions: item.Conversions})
	}
	return out
}

func seoIssueSummaryToPb(items []dto.SeoIssueSummary) []*crmpb.SeoIssueSummary {
	out := make([]*crmpb.SeoIssueSummary, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoIssueSummary{IssueType: item.IssueType, Severity: item.Severity, OpenCount: item.OpenCount})
	}
	return out
}

func seoPagePerformanceRowsToPb(items []dto.SeoPagePerformanceRow) []*crmpb.SeoPagePerformanceRow {
	out := make([]*crmpb.SeoPagePerformanceRow, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoPagePerformanceRow{
			SeoDomainId:    item.SeoDomainID,
			CanonicalUrl:   item.CanonicalURL,
			Slug:           item.Slug,
			Title:          item.Title,
			PageStatus:     item.PageStatus,
			RenderStatus:   item.RenderStatus,
			IsSiteMap:      item.IsSitemap,
			IsIndex:        item.IsIndex,
			QualityScore:   item.QualityScore,
			Clicks:         item.Clicks,
			Impressions:    item.Impressions,
			Ctr:            item.CTR,
			Position:       item.Position,
			Sessions:       item.Sessions,
			Conversions:    item.Conversions,
			IndexStatus:    item.IndexStatus,
			LastCrawlTime:  item.LastCrawlTime,
			OpenIssueCount: item.OpenIssueCount,
		})
	}
	return out
}

func seoQueryPerformanceRowsToPb(items []dto.SeoQueryPerformanceRow) []*crmpb.SeoQueryPerformanceRow {
	out := make([]*crmpb.SeoQueryPerformanceRow, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoQueryPerformanceRow{Query: item.Query, CanonicalUrl: item.CanonicalURL, Device: item.Device, Country: item.Country, Clicks: item.Clicks, Impressions: item.Impressions, Ctr: item.CTR, Position: item.Position})
	}
	return out
}

func seoIndexCoverageRowsToPb(items []dto.SeoIndexCoverageRow) []*crmpb.SeoIndexCoverageRow {
	out := make([]*crmpb.SeoIndexCoverageRow, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoIndexCoverageRow{SeoDomainId: item.SeoDomainID, CanonicalUrl: item.CanonicalURL, InspectionStatus: item.InspectionStatus, IndexStatus: item.IndexStatus, CoverageState: item.CoverageState, RobotsTxtState: item.RobotsTxtState, GoogleCanonical: item.GoogleCanonical, UserCanonical: item.UserCanonical, LastCrawlTime: item.LastCrawlTime, Verdict: item.Verdict, InspectedAt: item.InspectedAt})
	}
	return out
}

func seoTrafficRowsToPb(items []dto.SeoTrafficRow) []*crmpb.SeoTrafficRow {
	out := make([]*crmpb.SeoTrafficRow, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoTrafficRow{CanonicalUrl: item.CanonicalURL, Date: item.Date, Sessions: item.Sessions, Users: item.Users, ActiveUsers: item.ActiveUsers, PageViews: item.PageViews, EngagementSeconds: item.EngagementSeconds, Conversions: item.Conversions, Source: item.Source, Medium: item.Medium, Device: item.Device, Country: item.Country})
	}
	return out
}

func seoRankingRowsToPb(items []dto.SeoRankingRow) []*crmpb.SeoRankingRow {
	out := make([]*crmpb.SeoRankingRow, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoRankingRow{Keyword: item.Keyword, CanonicalUrl: item.CanonicalURL, Position: item.Position, PreviousPosition: item.PreviousPosition, PositionDelta: item.PositionDelta, Device: item.Device, Location: item.Location, Language: item.Language, CheckedAt: item.CheckedAt})
	}
	return out
}

func seoCrawlEventRowsToPb(items []dto.SeoCrawlEventRow) []*crmpb.SeoCrawlEventRow {
	out := make([]*crmpb.SeoCrawlEventRow, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoCrawlEventRow{Id: item.ID, SeoDomainId: item.SeoDomainID, CanonicalUrl: item.CanonicalURL, BotName: item.BotName, UserAgent: item.UserAgent, StatusCode: item.StatusCode, Method: item.Method, ResponseTimeMs: item.ResponseTimeMs, CrawledAt: item.CrawledAt})
	}
	return out
}

func seoIssueRowsToPb(items []dto.SeoIssueRow) []*crmpb.SeoIssueRow {
	out := make([]*crmpb.SeoIssueRow, 0, len(items))
	for _, item := range items {
		out = append(out, &crmpb.SeoIssueRow{Id: item.ID, SeoDomainId: item.SeoDomainID, CanonicalUrl: item.CanonicalURL, IssueType: item.IssueType, Severity: item.Severity, Status: item.Status, Title: item.Title, Message: item.Message, Evidence: item.Evidence, Recommendation: item.Recommendation, DetectedAt: item.DetectedAt, ResolvedAt: item.ResolvedAt})
	}
	return out
}
