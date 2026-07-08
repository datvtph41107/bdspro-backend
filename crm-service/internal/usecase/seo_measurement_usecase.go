package usecase

import (
	_commonDto "common/domain/dto"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/repo"

	"gorm.io/datatypes"
)

// SeoMeasurementUsecase owns the SEO Measurement & Feedback Loop.
//
// Naming note:
//   - This is intentionally NOT called Intelligence.
//   - The purpose is explicit and operational: measure SEO results, detect
//     feedback, and guide admin actions.
//
// This usecase is mock-first:
//   - It generates realistic mock data so admin/frontend can be built now.
//   - Later, Google Search Console, Google Analytics, crawl logs and rank
//     tracker connectors only need to write into the same tables.
type SeoMeasurementUsecase struct {
	seoMeasurementRepo repo.SeoMeasurementRepo
	seoDomainRepo      repo.SeoDomainRepo
}

func NewSeoMeasurementUsecase(
	seoMeasurementRepo repo.SeoMeasurementRepo,
	seoDomainRepo repo.SeoDomainRepo,
) *SeoMeasurementUsecase {
	return &SeoMeasurementUsecase{
		seoMeasurementRepo: seoMeasurementRepo,
		seoDomainRepo:      seoDomainRepo,
	}
}

func (u *SeoMeasurementUsecase) GetOverview(ctx context.Context, req *dto.SeoMeasurementRangeRequest) (*dto.SeoMeasurementOverviewResponse, error) {
	if req == nil {
		req = &dto.SeoMeasurementRangeRequest{Range: "28d", Limit: 50}
	}
	return u.seoMeasurementRepo.GetOverview(ctx, req)
}

func (u *SeoMeasurementUsecase) ListPagePerformance(ctx context.Context, req *dto.SeoPagePerformanceListRequest) ([]dto.SeoPagePerformanceRow, int64, error) {
	if req == nil {
		req = &dto.SeoPagePerformanceListRequest{}
	}
	req.Normalize()
	return u.seoMeasurementRepo.ListPagePerformance(ctx, req)
}

func (u *SeoMeasurementUsecase) ListQueryPerformance(ctx context.Context, req *dto.SeoQueryPerformanceListRequest) ([]dto.SeoQueryPerformanceRow, int64, error) {
	if req == nil {
		req = &dto.SeoQueryPerformanceListRequest{}
	}
	req.Normalize()
	return u.seoMeasurementRepo.ListQueryPerformance(ctx, req)
}

func (u *SeoMeasurementUsecase) ListIndexCoverage(ctx context.Context, req *dto.SeoIndexCoverageListRequest) ([]dto.SeoIndexCoverageRow, int64, error) {
	if req == nil {
		req = &dto.SeoIndexCoverageListRequest{}
	}
	req.Normalize()
	return u.seoMeasurementRepo.ListIndexCoverage(ctx, req)
}

func (u *SeoMeasurementUsecase) ListTraffic(ctx context.Context, req *dto.SeoTrafficListRequest) ([]dto.SeoTrafficRow, int64, error) {
	if req == nil {
		req = &dto.SeoTrafficListRequest{}
	}
	req.Normalize()
	return u.seoMeasurementRepo.ListTraffic(ctx, req)
}

func (u *SeoMeasurementUsecase) ListRanking(ctx context.Context, req *dto.SeoRankingListRequest) ([]dto.SeoRankingRow, int64, error) {
	if req == nil {
		req = &dto.SeoRankingListRequest{}
	}
	req.Normalize()
	return u.seoMeasurementRepo.ListRanking(ctx, req)
}

func (u *SeoMeasurementUsecase) ListCrawlEvents(ctx context.Context, req *dto.SeoCrawlEventListRequest) ([]dto.SeoCrawlEventRow, int64, error) {
	if req == nil {
		req = &dto.SeoCrawlEventListRequest{}
	}
	req.Normalize()
	return u.seoMeasurementRepo.ListCrawlEvents(ctx, req)
}

func (u *SeoMeasurementUsecase) ListIssues(ctx context.Context, req *dto.SeoIssueListRequest) ([]dto.SeoIssueRow, int64, error) {
	if req == nil {
		req = &dto.SeoIssueListRequest{}
	}
	req.Normalize()
	return u.seoMeasurementRepo.ListIssues(ctx, req)
}

func (u *SeoMeasurementUsecase) ResolveIssue(ctx context.Context, req *dto.SeoResolveIssueRequest) error {
	if req == nil || req.ID == 0 {
		return fmt.Errorf("seo issue id is required")
	}
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = seo_domain.SeoIssueStatusResolved
	}
	return u.seoMeasurementRepo.ResolveIssue(ctx, req.ID, status)
}

func (u *SeoMeasurementUsecase) GenerateMockMeasurement(ctx context.Context, req *dto.SeoMeasurementGenerateMockRequest) (*dto.SeoMeasurementGenerateMockResponse, error) {
	if req == nil {
		req = &dto.SeoMeasurementGenerateMockRequest{Range: "28d", Limit: 200}
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 200
	}
	if limit > 1000 {
		limit = 1000
	}

	pages, _, err := u.seoDomainRepo.GetList(ctx, &dto.SeoDomainListRequest{Pagable: _commonDto.Pagable{Page: 0, Size: limit}})
	if err != nil {
		return nil, err
	}

	from, to := seoMeasurementDateRange(req.Range)
	var searchRows []seo_domain.SeoSearchPerformanceDaily
	var inspectionRows []seo_domain.SeoURLInspectionSnapshot
	var trafficRows []seo_domain.SeoTrafficDaily
	var keywordTargets []seo_domain.SeoKeywordTarget
	var rankingRows []seo_domain.SeoRankingSnapshot
	var crawlRows []seo_domain.SeoCrawlEvent

	for i := range pages {
		page := &pages[i]
		if strings.TrimSpace(page.CanonicalURL) == "" {
			continue
		}

		queries := buildMockQueries(page)
		for _, query := range queries {
			keywordTargets = append(keywordTargets, seo_domain.SeoKeywordTarget{
				SeoDomainID: page.ID,
				Keyword:     query,
				Intent:      mockIntentForPage(page),
				Priority:    mockPriorityForPage(page),
				TargetURL:   page.CanonicalURL,
				Status:      "active",
				Note:        "mock-first target; replace with real keyword strategy later",
			})
		}

		for day := from; !day.After(to); day = day.AddDate(0, 0, 1) {
			for qIndex, query := range queries {
				impressions := mockImpressions(page, qIndex, day)
				ctr := mockCTR(page, qIndex)
				clicks := math.Round(impressions*ctr*100) / 100
				position := mockPosition(page, qIndex)
				searchRows = append(searchRows, seo_domain.SeoSearchPerformanceDaily{
					SeoDomainID:  &page.ID,
					CanonicalURL: page.CanonicalURL,
					Query:        query,
					Device:       "desktop",
					Country:      "vn",
					SearchType:   "web",
					Date:         day,
					Clicks:       clicks,
					Impressions:  impressions,
					CTR:          ctr,
					Position:     position,
					DataSource:   seo_domain.SeoMeasurementDataSourceMock,
					ImportedAt:   time.Now(),
				})
			}

			sessions := math.Round(mockDailyClicks(page)*1.25*100) / 100
			conversions := 0.0
			if sessions > 20 && page.QualityScore >= 70 {
				conversions = math.Round(sessions*0.025*100) / 100
			}
			trafficRows = append(trafficRows, seo_domain.SeoTrafficDaily{
				SeoDomainID:       &page.ID,
				CanonicalURL:      page.CanonicalURL,
				Date:              day,
				Sessions:          sessions,
				Users:             math.Round(sessions*0.82*100) / 100,
				ActiveUsers:       math.Round(sessions*0.7*100) / 100,
				PageViews:         math.Round(sessions*1.35*100) / 100,
				EngagementSeconds: sessions * 72,
				Conversions:       conversions,
				Source:            "google",
				Medium:            "organic",
				Device:            "desktop",
				Country:           "vn",
				DataSource:        seo_domain.SeoMeasurementDataSourceMock,
				ImportedAt:        time.Now(),
			})
		}

		indexed := page.Published && page.IsIndex && strings.EqualFold(page.RenderStatus, seo_domain.SeoRenderStatusSuccess)
		indexStatus := seo_domain.SeoIndexStatusNotIndexed
		coverageState := "URL is not on Google"
		verdict := "neutral"
		lastCrawl := time.Now().Add(-24 * time.Hour)
		if indexed {
			indexStatus = seo_domain.SeoIndexStatusIndexed
			coverageState = "Submitted and indexed"
			verdict = "pass"
		} else if page.Published {
			indexStatus = seo_domain.SeoIndexStatusCrawledNotIndexed
			coverageState = "Crawled - currently not indexed"
			verdict = "partial"
		}
		inspectionRows = append(inspectionRows, seo_domain.SeoURLInspectionSnapshot{
			SeoDomainID:      &page.ID,
			CanonicalURL:     page.CanonicalURL,
			InspectionStatus: "complete",
			IndexStatus:      indexStatus,
			CoverageState:    coverageState,
			RobotsTxtState:   robotsStateFromPage(page),
			IndexingState:    indexStatus,
			GoogleCanonical:  mockGoogleCanonical(page),
			UserCanonical:    page.CanonicalURL,
			LastCrawlTime:    &lastCrawl,
			PageFetchState:   mockFetchState(page),
			Verdict:          verdict,
			RawJSON:          mustJSON(map[string]any{"mock": true, "pageStatus": page.PageStatus, "renderStatus": page.RenderStatus}),
			DataSource:       seo_domain.SeoMeasurementDataSourceMock,
			InspectedAt:      time.Now(),
		})

		if page.Published {
			statusCode := 200
			if strings.EqualFold(page.RenderStatus, seo_domain.SeoRenderStatusFailed) {
				statusCode = 500
			}
			crawlRows = append(crawlRows, seo_domain.SeoCrawlEvent{
				SeoDomainID:    &page.ID,
				CanonicalURL:   page.CanonicalURL,
				BotName:        "Googlebot",
				UserAgent:      "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
				StatusCode:     statusCode,
				Method:         "GET",
				ResponseTimeMs: int64(120 + (page.ID % 80)),
				IP:             "66.249.66.1",
				CrawledAt:      time.Now().Add(-12 * time.Hour),
				RawJSON:        mustJSON(map[string]any{"mock": true}),
			})
		}

		for qIndex, query := range queries {
			position := mockPosition(page, qIndex)
			rankingRows = append(rankingRows, seo_domain.SeoRankingSnapshot{
				SeoDomainID:      &page.ID,
				Keyword:          query,
				CanonicalURL:     page.CanonicalURL,
				Position:         position,
				PreviousPosition: position + float64((page.ID%5)-2),
				PositionDelta:    float64((page.ID % 5) - 2),
				Device:           "desktop",
				Location:         "Vietnam",
				Language:         "vi",
				DataSource:       seo_domain.SeoMeasurementDataSourceMock,
				CheckedAt:        time.Now(),
				RawJSON:          mustJSON(map[string]any{"mock": true}),
			})
		}
	}

	if err := u.seoMeasurementRepo.UpsertKeywordTargets(ctx, keywordTargets); err != nil {
		return nil, err
	}
	if err := u.seoMeasurementRepo.UpsertSearchPerformanceDaily(ctx, searchRows); err != nil {
		return nil, err
	}
	if err := u.seoMeasurementRepo.UpsertURLInspectionSnapshots(ctx, inspectionRows); err != nil {
		return nil, err
	}
	if err := u.seoMeasurementRepo.UpsertTrafficDaily(ctx, trafficRows); err != nil {
		return nil, err
	}
	if err := u.seoMeasurementRepo.CreateRankingSnapshots(ctx, rankingRows); err != nil {
		return nil, err
	}
	if err := u.seoMeasurementRepo.CreateCrawlEvents(ctx, crawlRows); err != nil {
		return nil, err
	}

	return &dto.SeoMeasurementGenerateMockResponse{
		Pages:          uint32(len(pages)),
		SearchRows:     uint32(len(searchRows)),
		InspectionRows: uint32(len(inspectionRows)),
		TrafficRows:    uint32(len(trafficRows)),
		RankingRows:    uint32(len(rankingRows)),
		CrawlRows:      uint32(len(crawlRows)),
		KeywordTargets: uint32(len(keywordTargets)),
		DataSource:     seo_domain.SeoMeasurementDataSourceMock,
		Messages: []string{
			"Mock measurement data generated. Replace provider later with Google Search Console, GA, crawl logs and rank tracker.",
		},
	}, nil
}

func (u *SeoMeasurementUsecase) RunFeedbackDetector(ctx context.Context, req *dto.SeoMeasurementDetectorRequest) (*dto.SeoMeasurementDetectorResponse, error) {
	if req == nil {
		req = &dto.SeoMeasurementDetectorRequest{Range: "28d", Limit: 500}
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 500
	}
	if limit > 1000 {
		limit = 1000
	}

	rows, _, err := u.seoMeasurementRepo.ListPagePerformance(ctx, &dto.SeoPagePerformanceListRequest{Pagable: _commonDto.Pagable{Page: 0, Size: limit}, Range: req.Range})
	if err != nil {
		return nil, err
	}

	issues := make([]seo_domain.SeoIssue, 0)
	for i := range rows {
		row := rows[i]
		issues = append(issues, detectIssuesForPage(row)...)
	}
	if err := u.seoMeasurementRepo.CreateIssues(ctx, issues); err != nil {
		return nil, err
	}
	return &dto.SeoMeasurementDetectorResponse{Scanned: uint32(len(rows)), Created: uint32(len(issues)), Errors: []string{}}, nil
}

func detectIssuesForPage(row dto.SeoPagePerformanceRow) []seo_domain.SeoIssue {
	var issues []seo_domain.SeoIssue
	seoID := row.SeoDomainID
	seoIDPtr := &seoID
	if row.PageStatus == seo_domain.SeoPageStatusPublished && row.RenderStatus != seo_domain.SeoRenderStatusSuccess {
		issues = append(issues, buildIssue(seoIDPtr, row.CanonicalURL, seo_domain.SeoIssueTypeNotRendered, seo_domain.SeoIssueSeverityHigh, "Published page has not rendered successfully", "Page is published but rendered HTML is not successful.", "Run GenerateSeoPage, check template/source data and retry."))
	}
	if row.PageStatus == seo_domain.SeoPageStatusPublished && row.RenderStatus == seo_domain.SeoRenderStatusSuccess && !row.IsSitemap && row.IsIndex {
		issues = append(issues, buildIssue(seoIDPtr, row.CanonicalURL, seo_domain.SeoIssueTypeNotInSitemap, seo_domain.SeoIssueSeverityMedium, "Indexable page is not in sitemap", "Page can be indexed but is not included in sitemap.", "Check sitemap gate: is_site_map, quality_score, source_status and canonical URL."))
	}
	if row.IndexStatus == seo_domain.SeoIndexStatusCrawledNotIndexed || row.IndexStatus == seo_domain.SeoIndexStatusNotIndexed {
		issues = append(issues, buildIssue(seoIDPtr, row.CanonicalURL, seo_domain.SeoIssueTypeCrawledNotIndexed, seo_domain.SeoIssueSeverityHigh, "URL is crawled but not indexed", "Google-like inspection state says the page is not indexed.", "Improve content uniqueness, canonical consistency, internal links and page quality."))
	}
	if row.Impressions >= 100 && row.CTR < 0.01 {
		issues = append(issues, buildIssue(seoIDPtr, row.CanonicalURL, seo_domain.SeoIssueTypeHighImpressionLowCTR, seo_domain.SeoIssueSeverityHigh, "High impressions but low CTR", "The page is shown in search but users are not clicking enough.", "Rewrite title/meta description to match user intent and improve snippet attractiveness."))
	}
	if row.Position > 0 && row.Position <= 5 && row.Impressions >= 50 && row.CTR < 0.03 {
		issues = append(issues, buildIssue(seoIDPtr, row.CanonicalURL, seo_domain.SeoIssueTypeGoodPositionLowCTR, seo_domain.SeoIssueSeverityMedium, "Good position but low CTR", "The page ranks well but the snippet underperforms.", "Improve title, meta description, structured data and intent match."))
	}
	if row.Sessions >= 100 && row.Conversions <= 0 {
		issues = append(issues, buildIssue(seoIDPtr, row.CanonicalURL, seo_domain.SeoIssueTypeTrafficNoConversion, seo_domain.SeoIssueSeverityMedium, "Traffic without conversion", "The page receives traffic but does not convert.", "Add clearer CTA, lead form, phone/action links or internal links to conversion pages."))
	}
	return issues
}

func buildIssue(seoDomainID *uint64, canonicalURL, issueType, severity, title, message, recommendation string) seo_domain.SeoIssue {
	return seo_domain.SeoIssue{
		SeoDomainID:    seoDomainID,
		CanonicalURL:   canonicalURL,
		IssueType:      issueType,
		Severity:       severity,
		Status:         seo_domain.SeoIssueStatusOpen,
		Title:          title,
		Message:        message,
		Evidence:       mustJSON(map[string]any{"rule": issueType, "detectedBy": "seo_measurement_feedback_loop"}),
		Recommendation: recommendation,
		DetectedAt:     time.Now(),
	}
}

func buildMockQueries(page *seo_domain.SeoDomain) []string {
	base := strings.TrimSpace(page.Title)
	if base == "" {
		base = strings.Trim(strings.TrimSpace(page.Slug), "/")
	}
	if base == "" {
		base = "seo page"
	}
	ref := strings.TrimSpace(page.RefLabel)
	queries := []string{base}
	if ref != "" && !strings.EqualFold(ref, base) {
		queries = append(queries, ref)
	}
	queries = append(queries, fmt.Sprintf("%s quy hoạch", base))
	return uniqueStrings(queries)
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		key := strings.ToLower(value)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, value)
	}
	return out
}

func mockIntentForPage(page *seo_domain.SeoDomain) string {
	if page == nil {
		return "informational"
	}
	schema := strings.ToLower(page.GetSchema())
	if strings.Contains(schema, "parcel") {
		return "verification"
	}
	if strings.Contains(schema, "project") {
		return "commercial_investigation"
	}
	if strings.Contains(schema, "region") || strings.Contains(schema, "unit") {
		return "local_information"
	}
	return "informational"
}

func mockPriorityForPage(page *seo_domain.SeoDomain) int {
	if page == nil {
		return 3
	}
	if page.QualityScore >= 85 {
		return 0
	}
	if page.QualityScore >= 70 {
		return 1
	}
	if page.QualityScore >= 50 {
		return 2
	}
	return 3
}

func mockImpressions(page *seo_domain.SeoDomain, queryIndex int, day time.Time) float64 {
	base := 8.0 + float64(page.QualityScore)/2
	if page.Published {
		base += 25
	}
	if page.IsSiteMap {
		base += 15
	}
	if page.RenderStatus == seo_domain.SeoRenderStatusSuccess {
		base += 20
	}
	base = base / float64(queryIndex+1)
	seasonality := 1 + float64(day.Day()%7)/20
	return math.Round(base*seasonality*100) / 100
}

func mockCTR(page *seo_domain.SeoDomain, queryIndex int) float64 {
	ctr := 0.012 + float64(page.QualityScore)/5000
	if page.QualityScore >= 80 {
		ctr += 0.015
	}
	if queryIndex > 0 {
		ctr = ctr * 0.72
	}
	if ctr > 0.12 {
		ctr = 0.12
	}
	return math.Round(ctr*10000) / 10000
}

func mockDailyClicks(page *seo_domain.SeoDomain) float64 {
	return mockImpressions(page, 0, time.Now()) * mockCTR(page, 0)
}

func mockPosition(page *seo_domain.SeoDomain, queryIndex int) float64 {
	position := 28.0 - float64(page.QualityScore)/4 + float64(queryIndex*4)
	if page.RenderStatus == seo_domain.SeoRenderStatusSuccess {
		position -= 3
	}
	if page.IsSiteMap {
		position -= 2
	}
	if position < 1.2 {
		position = 1.2
	}
	return math.Round(position*100) / 100
}

func robotsStateFromPage(page *seo_domain.SeoDomain) string {
	if page == nil || !page.IsIndex {
		return "noindex"
	}
	return "allowed"
}

func mockFetchState(page *seo_domain.SeoDomain) string {
	if page == nil {
		return "not_found"
	}
	if page.RenderStatus == seo_domain.SeoRenderStatusFailed {
		return "server_error"
	}
	if page.Published {
		return "successful"
	}
	return "not_found"
}

func mockGoogleCanonical(page *seo_domain.SeoDomain) string {
	if page == nil {
		return ""
	}
	if page.QualityScore < 35 && page.OriginURL != "" {
		return page.OriginURL
	}
	return page.CanonicalURL
}

func mustJSON(value any) datatypes.JSON {
	raw, err := json.Marshal(value)
	if err != nil {
		return datatypes.JSON([]byte(`{}`))
	}
	return datatypes.JSON(raw)
}

func seoMeasurementDateRange(rangeValue string) (time.Time, time.Time) {
	now := time.Now()
	days := 28
	switch strings.TrimSpace(rangeValue) {
	case "7d":
		days = 7
	case "14d":
		days = 14
	case "90d":
		days = 90
	}
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	from := to.AddDate(0, 0, -(days - 1))
	return from, to
}
