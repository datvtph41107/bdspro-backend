package usecase

import (
	_commonDto "common/domain/dto"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"

	"gorm.io/datatypes"
)

type SeoMeasurementValidationError struct {
	err error
}

func (e *SeoMeasurementValidationError) Error() string { return e.err.Error() }
func (e *SeoMeasurementValidationError) Unwrap() error { return e.err }

func measurementValidationError(err error) error {
	if err == nil {
		return nil
	}
	return &SeoMeasurementValidationError{err: err}
}

func measurementValidationErrorf(format string, args ...interface{}) error {
	return measurementValidationError(fmt.Errorf(format, args...))
}

var seoMeasurementSources = map[string]int{
	seo_domain.SeoMeasurementDataSourceSearchConsole: 48,
	seo_domain.SeoMeasurementDataSourceURLInspection: 72,
	seo_domain.SeoMeasurementDataSourceAnalytics:     48,
	seo_domain.SeoMeasurementDataSourceCrawlLog:      24,
	seo_domain.SeoMeasurementDataSourceRankTracker:   168,
}

func normalizeMeasurementSource(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "gsc", "google_search_console", "search-console":
		return seo_domain.SeoMeasurementDataSourceSearchConsole
	case "url-inspection", "google_url_inspection", "inspection":
		return seo_domain.SeoMeasurementDataSourceURLInspection
	case "ga", "ga4", "google_analytics":
		return seo_domain.SeoMeasurementDataSourceAnalytics
	case "crawl", "nginx", "gateway_log":
		return seo_domain.SeoMeasurementDataSourceCrawlLog
	case "rank", "ranking":
		return seo_domain.SeoMeasurementDataSourceRankTracker
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}

func (u *SeoMeasurementUsecase) ImportMeasurement(ctx context.Context, req *dto.SeoMeasurementImportRequest) (*dto.SeoMeasurementImportResponse, error) {
	if req == nil {
		return nil, measurementValidationErrorf("measurement import request is required")
	}
	source := normalizeMeasurementSource(req.Source)
	staleHours, supported := seoMeasurementSources[source]
	if !supported {
		return nil, measurementValidationErrorf("unsupported measurement source %q", req.Source)
	}
	if len(req.Rows) == 0 {
		return nil, measurementValidationErrorf("measurement rows are required")
	}
	if len(req.Rows) > 10000 {
		return nil, measurementValidationErrorf("measurement import is limited to 10000 rows per run")
	}

	rangeStart, err := optionalMeasurementTime(req.RangeStart)
	if err != nil {
		return nil, measurementValidationError(fmt.Errorf("invalid rangeStart: %w", err))
	}
	rangeEnd, err := optionalMeasurementTime(req.RangeEnd)
	if err != nil {
		return nil, measurementValidationError(fmt.Errorf("invalid rangeEnd: %w", err))
	}
	if rangeStart != nil && rangeEnd != nil && rangeStart.After(*rangeEnd) {
		return nil, measurementValidationErrorf("rangeStart must not be after rangeEnd")
	}
	checksum, err := measurementImportChecksum(source, req)
	if err != nil {
		return nil, err
	}
	existing, err := u.seoMeasurementRepo.GetReusableImportRunByChecksum(ctx, source, checksum)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return measurementImportResponseFromRun(existing, true), nil
	}
	metadata := mustJSON(req.Metadata)
	run := &seo_domain.SeoMeasurementImportRun{
		Source:      source,
		Status:      seo_domain.SeoMeasurementImportStatusRunning,
		RequestedBy: strings.TrimSpace(req.RequestedBy),
		RangeStart:  rangeStart,
		RangeEnd:    rangeEnd,
		RowCount:    int64(len(req.Rows)),
		Checksum:    checksum,
		Metadata:    metadata,
		StartedAt:   time.Now().UTC(),
	}
	if err := u.seoMeasurementRepo.CreateImportRun(ctx, run); err != nil {
		// A concurrent request may have inserted the same active checksum after
		// the preflight read. Return that run instead of importing twice.
		existing, lookupErr := u.seoMeasurementRepo.GetReusableImportRunByChecksum(ctx, source, checksum)
		if lookupErr == nil && existing != nil {
			return measurementImportResponseFromRun(existing, true), nil
		}
		return nil, err
	}

	response := &dto.SeoMeasurementImportResponse{
		RunID:    run.ID,
		Source:   source,
		Status:   seo_domain.SeoMeasurementImportStatusRunning,
		RowCount: int64(len(req.Rows)),
		Errors:   []string{},
	}
	now := time.Now().UTC()
	seoIDCache := map[string]*uint64{}
	resolveSeoID := func(canonicalURL string) *uint64 {
		if value, ok := seoIDCache[canonicalURL]; ok {
			return value
		}
		page, lookupErr := u.seoDomainRepo.GetByCanonicalURL(ctx, canonicalURL)
		if lookupErr != nil || page == nil {
			seoIDCache[canonicalURL] = nil
			return nil
		}
		id := page.ID
		seoIDCache[canonicalURL] = &id
		return &id
	}

	var searchRows []seo_domain.SeoSearchPerformanceDaily
	var inspectionRows []seo_domain.SeoURLInspectionSnapshot
	var trafficRows []seo_domain.SeoTrafficDaily
	var rankingRows []seo_domain.SeoRankingSnapshot
	var crawlRows []seo_domain.SeoCrawlEvent

	reject := func(index int, err error) {
		response.RejectedCount++
		if len(response.Errors) < 25 {
			response.Errors = append(response.Errors, fmt.Sprintf("row %d: %v", index+1, err))
		}
	}

	for index, row := range req.Rows {
		canonicalURL := mapString(row, "canonicalUrl", "canonical_url", "url")
		if err := validateCanonicalMeasurementURL(canonicalURL); err != nil {
			reject(index, err)
			continue
		}
		seoDomainID := resolveSeoID(canonicalURL)

		switch source {
		case seo_domain.SeoMeasurementDataSourceSearchConsole:
			date, parseErr := requiredMeasurementDate(mapString(row, "date"))
			if parseErr != nil {
				reject(index, parseErr)
				continue
			}
			clicks := mapFloat(row, "clicks")
			impressions := mapFloat(row, "impressions")
			ctr := mapFloat(row, "ctr")
			if ctr == 0 && impressions > 0 {
				ctr = clicks / impressions
			}
			searchRows = append(searchRows, seo_domain.SeoSearchPerformanceDaily{
				SeoDomainID: seoDomainID, CanonicalURL: canonicalURL,
				Query: mapString(row, "query"), Device: mapStringDefault(row, "desktop", "device"),
				Country: mapStringDefault(row, "vn", "country"), SearchType: mapStringDefault(row, "web", "searchType", "search_type"),
				Date: date, Clicks: nonNegative(clicks), Impressions: nonNegative(impressions), CTR: clamp(ctr, 0, 1),
				Position: nonNegative(mapFloat(row, "position")), DataSource: source, SourceRunID: &run.ID, ImportedAt: now,
			})

		case seo_domain.SeoMeasurementDataSourceURLInspection:
			inspectedAt, parseErr := measurementTimeDefault(mapString(row, "inspectedAt", "inspected_at"), now)
			if parseErr != nil {
				reject(index, parseErr)
				continue
			}
			lastCrawl, parseErr := optionalMeasurementTime(mapString(row, "lastCrawlTime", "last_crawl_time"))
			if parseErr != nil {
				reject(index, parseErr)
				continue
			}
			raw, _ := json.Marshal(row)
			inspectionRows = append(inspectionRows, seo_domain.SeoURLInspectionSnapshot{
				SeoDomainID: seoDomainID, CanonicalURL: canonicalURL,
				InspectionStatus: mapStringDefault(row, "complete", "inspectionStatus", "inspection_status"),
				IndexStatus:      mapStringDefault(row, seo_domain.SeoIndexStatusUnknown, "indexStatus", "index_status"),
				CoverageState:    mapString(row, "coverageState", "coverage_state"), RobotsTxtState: mapString(row, "robotsTxtState", "robots_txt_state"),
				IndexingState: mapString(row, "indexingState", "indexing_state"), GoogleCanonical: mapString(row, "googleCanonical", "google_canonical"),
				UserCanonical: mapStringDefault(row, canonicalURL, "userCanonical", "user_canonical"), LastCrawlTime: lastCrawl,
				PageFetchState: mapString(row, "pageFetchState", "page_fetch_state"), Verdict: mapString(row, "verdict"),
				RawJSON: datatypes.JSON(raw), DataSource: source, SourceRunID: &run.ID, InspectedAt: inspectedAt,
			})

		case seo_domain.SeoMeasurementDataSourceAnalytics:
			date, parseErr := requiredMeasurementDate(mapString(row, "date"))
			if parseErr != nil {
				reject(index, parseErr)
				continue
			}
			trafficRows = append(trafficRows, seo_domain.SeoTrafficDaily{
				SeoDomainID: seoDomainID, CanonicalURL: canonicalURL, Date: date,
				Sessions: nonNegative(mapFloat(row, "sessions")), Users: nonNegative(mapFloat(row, "users")),
				ActiveUsers: nonNegative(mapFloat(row, "activeUsers", "active_users")), PageViews: nonNegative(mapFloat(row, "pageViews", "page_views")),
				EngagementSeconds: nonNegative(mapFloat(row, "engagementSeconds", "engagement_seconds")), Conversions: nonNegative(mapFloat(row, "conversions")),
				Source: mapStringDefault(row, "google", "trafficSource", "source"), Medium: mapStringDefault(row, "organic", "medium"),
				Device: mapStringDefault(row, "desktop", "device"), Country: mapStringDefault(row, "vn", "country"),
				DataSource: source, SourceRunID: &run.ID, ImportedAt: now,
			})

		case seo_domain.SeoMeasurementDataSourceRankTracker:
			checkedAt, parseErr := measurementTimeDefault(mapString(row, "checkedAt", "checked_at"), now)
			if parseErr != nil {
				reject(index, parseErr)
				continue
			}
			keyword := mapString(row, "keyword", "query")
			if strings.TrimSpace(keyword) == "" {
				reject(index, fmt.Errorf("keyword is required"))
				continue
			}
			raw, _ := json.Marshal(row)
			position := nonNegative(mapFloat(row, "position"))
			previous := nonNegative(mapFloat(row, "previousPosition", "previous_position"))
			delta := mapFloat(row, "positionDelta", "position_delta")
			if delta == 0 && previous > 0 && position > 0 {
				delta = previous - position
			}
			rankingRows = append(rankingRows, seo_domain.SeoRankingSnapshot{
				SeoDomainID: seoDomainID, Keyword: keyword, CanonicalURL: canonicalURL, Position: position,
				PreviousPosition: previous, PositionDelta: delta, Device: mapStringDefault(row, "desktop", "device"),
				Location: mapStringDefault(row, "vn", "location"), Language: mapStringDefault(row, "vi", "language"),
				DataSource: source, SourceRunID: &run.ID, CheckedAt: checkedAt, RawJSON: datatypes.JSON(raw),
			})

		case seo_domain.SeoMeasurementDataSourceCrawlLog:
			crawledAt, parseErr := measurementTimeDefault(mapString(row, "crawledAt", "crawled_at", "timestamp"), now)
			if parseErr != nil {
				reject(index, parseErr)
				continue
			}
			statusCode := int(mapFloat(row, "statusCode", "status_code"))
			if statusCode < 100 || statusCode > 599 {
				reject(index, fmt.Errorf("statusCode must be between 100 and 599"))
				continue
			}
			raw, _ := json.Marshal(row)
			crawlRows = append(crawlRows, seo_domain.SeoCrawlEvent{
				SeoDomainID: seoDomainID, CanonicalURL: canonicalURL, BotName: mapStringDefault(row, "unknown", "botName", "bot_name"),
				UserAgent: mapString(row, "userAgent", "user_agent"), StatusCode: statusCode,
				Method: mapStringDefault(row, "GET", "method"), ResponseTimeMs: int64(nonNegative(mapFloat(row, "responseTimeMs", "response_time_ms"))),
				Referer: mapString(row, "referer"), IP: mapString(row, "ip"), CrawledAt: crawledAt,
				RawJSON: datatypes.JSON(raw), DataSource: source, SourceRunID: &run.ID,
			})
		}
		response.AcceptedCount++
	}

	persistErr := func() error {
		switch source {
		case seo_domain.SeoMeasurementDataSourceSearchConsole:
			return u.seoMeasurementRepo.UpsertSearchPerformanceDaily(ctx, searchRows)
		case seo_domain.SeoMeasurementDataSourceURLInspection:
			return u.seoMeasurementRepo.UpsertURLInspectionSnapshots(ctx, inspectionRows)
		case seo_domain.SeoMeasurementDataSourceAnalytics:
			return u.seoMeasurementRepo.UpsertTrafficDaily(ctx, trafficRows)
		case seo_domain.SeoMeasurementDataSourceRankTracker:
			return u.seoMeasurementRepo.CreateRankingSnapshots(ctx, rankingRows)
		case seo_domain.SeoMeasurementDataSourceCrawlLog:
			return u.seoMeasurementRepo.CreateCrawlEvents(ctx, crawlRows)
		default:
			return fmt.Errorf("unsupported source %s", source)
		}
	}()

	finished := time.Now().UTC()
	run.FinishedAt = &finished
	run.AcceptedCount = response.AcceptedCount
	run.RejectedCount = response.RejectedCount
	if persistErr != nil {
		run.Status = seo_domain.SeoMeasurementImportStatusFailed
		run.ErrorMessage = persistErr.Error()
		response.Status = run.Status
		_ = u.seoMeasurementRepo.CompleteImportRun(ctx, run)
		_ = u.seoMeasurementRepo.UpsertSourceState(ctx, &seo_domain.SeoMeasurementSourceState{
			Source: source, Status: run.Status, LastRunID: &run.ID, LastAttemptAt: &finished,
			LastRowCount: 0, StaleAfterHours: staleHours, ConsecutiveFailures: 1, LastError: persistErr.Error(),
		})
		return response, persistErr
	}

	if response.RejectedCount > 0 && len(response.Errors) > 0 {
		run.ErrorMessage = strings.Join(response.Errors, "; ")
		if len(run.ErrorMessage) > 4000 {
			run.ErrorMessage = run.ErrorMessage[:4000]
		}
	}
	switch {
	case response.AcceptedCount == 0:
		run.Status = seo_domain.SeoMeasurementImportStatusFailed
		run.ErrorMessage = "all rows were rejected"
	case response.RejectedCount > 0:
		run.Status = seo_domain.SeoMeasurementImportStatusPartial
	default:
		run.Status = seo_domain.SeoMeasurementImportStatusSucceeded
	}
	response.Status = run.Status
	if err := u.seoMeasurementRepo.CompleteImportRun(ctx, run); err != nil {
		return response, err
	}
	state := &seo_domain.SeoMeasurementSourceState{
		Source: source, Status: run.Status, LastRunID: &run.ID, LastAttemptAt: &finished,
		LastRowCount: response.AcceptedCount, StaleAfterHours: staleHours,
	}
	if response.AcceptedCount > 0 {
		state.LastSuccessAt = &finished
		if run.Status == seo_domain.SeoMeasurementImportStatusPartial {
			state.LastError = run.ErrorMessage
		}
	} else {
		state.ConsecutiveFailures = 1
		state.LastError = run.ErrorMessage
	}
	if err := u.seoMeasurementRepo.UpsertSourceState(ctx, state); err != nil {
		return response, err
	}
	return response, nil
}

func (u *SeoMeasurementUsecase) GetSourceHealth(ctx context.Context) ([]dto.SeoMeasurementSourceHealth, error) {
	rows, err := u.seoMeasurementRepo.ListSourceHealth(ctx)
	if err != nil {
		return nil, err
	}
	bySource := make(map[string]dto.SeoMeasurementSourceHealth, len(rows))
	for _, row := range rows {
		bySource[row.Source] = row
	}
	out := make([]dto.SeoMeasurementSourceHealth, 0, len(seoMeasurementSources))
	for _, source := range []string{
		seo_domain.SeoMeasurementDataSourceSearchConsole,
		seo_domain.SeoMeasurementDataSourceURLInspection,
		seo_domain.SeoMeasurementDataSourceAnalytics,
		seo_domain.SeoMeasurementDataSourceCrawlLog,
		seo_domain.SeoMeasurementDataSourceRankTracker,
	} {
		if row, ok := bySource[source]; ok {
			out = append(out, row)
			continue
		}
		out = append(out, dto.SeoMeasurementSourceHealth{Source: source, Status: "never_imported", StaleAfterHours: seoMeasurementSources[source]})
	}
	return out, nil
}

func (u *SeoMeasurementUsecase) ListImportRuns(ctx context.Context, req *dto.SeoMeasurementImportRunListRequest) ([]dto.SeoMeasurementImportRunRow, int64, error) {
	return u.seoMeasurementRepo.ListImportRuns(ctx, req)
}

func (u *SeoMeasurementUsecase) ListPageScores(ctx context.Context, req *dto.SeoPageScoreListRequest) ([]dto.SeoPageScoreRow, int64, error) {
	return u.seoMeasurementRepo.ListPageScores(ctx, req)
}

func (u *SeoMeasurementUsecase) RunScoring(ctx context.Context, req *dto.SeoMeasurementScoreRequest) (*dto.SeoMeasurementScoreResponse, error) {
	if req == nil {
		req = &dto.SeoMeasurementScoreRequest{Range: "28d", Limit: 1000}
	}
	if strings.TrimSpace(req.Range) == "" {
		req.Range = "28d"
	}
	limit := req.Limit
	if limit == 0 {
		limit = 1000
	}
	if limit > 5000 {
		limit = 5000
	}
	pages, _, err := u.seoMeasurementRepo.ListPagePerformance(ctx, &dto.SeoPagePerformanceListRequest{
		Pagable: _commonDto.Pagable{Page: 0, Size: limit}, Range: req.Range,
	})
	if err != nil {
		return nil, err
	}
	health, err := u.GetSourceHealth(ctx)
	if err != nil {
		return nil, err
	}
	freshness := map[string]interface{}{}
	staleSourceCount := 0
	for _, item := range health {
		freshness[item.Source] = map[string]interface{}{"status": item.Status, "fresh": item.Fresh, "ageHours": item.AgeHours, "lastSuccessAt": item.LastSuccessAt}
		if !item.Fresh {
			staleSourceCount++
		}
	}
	freshnessJSON := mustJSON(freshness)
	measuredAt := time.Now().UTC()
	snapshots := make([]seo_domain.SeoPageScoreSnapshot, 0, len(pages))
	result := &dto.SeoMeasurementScoreResponse{Scanned: uint32(len(pages))}
	for _, page := range pages {
		snapshot := scoreMeasurementPage(page, req.Range, measuredAt, freshnessJSON, staleSourceCount)
		snapshots = append(snapshots, snapshot)
		switch snapshot.Grade {
		case "A":
			result.A++
		case "B":
			result.B++
		case "C":
			result.C++
		case "D":
			result.D++
		default:
			result.F++
		}
	}
	if err := u.seoMeasurementRepo.UpsertPageScores(ctx, snapshots); err != nil {
		return nil, err
	}
	result.Persisted = uint32(len(snapshots))
	return result, nil
}

func scoreMeasurementPage(page dto.SeoPagePerformanceRow, rangeKey string, measuredAt time.Time, freshness datatypes.JSON, staleSourceCount int) seo_domain.SeoPageScoreSnapshot {
	blockers := []string{}
	warnings := []string{}
	indexScore := 30.0
	if !page.IsIndex {
		indexScore = 100
	} else {
		switch page.IndexStatus {
		case seo_domain.SeoIndexStatusIndexed:
			indexScore = 100
		case seo_domain.SeoIndexStatusCrawledNotIndexed:
			indexScore = 35
		case seo_domain.SeoIndexStatusNotIndexed:
			indexScore = 15
		case "blocked":
			indexScore = 0
		}
	}
	if page.IsIndex && page.IndexStatus != seo_domain.SeoIndexStatusIndexed {
		blockers = append(blockers, "index_not_confirmed")
	}

	visibility := 20.0
	if page.Impressions > 0 {
		switch {
		case page.Position > 0 && page.Position <= 3:
			visibility = 100
		case page.Position > 0 && page.Position <= 10:
			visibility = 85
		case page.Position > 0 && page.Position <= 20:
			visibility = 65
		case page.Position > 0 && page.Position <= 50:
			visibility = 40
		default:
			visibility = 25
		}
		if page.CTR >= 0.05 {
			visibility = math.Min(100, visibility+10)
		} else if page.Impressions >= 100 && page.CTR < 0.01 {
			visibility = math.Max(0, visibility-20)
			warnings = append(warnings, "low_ctr")
		}
	} else {
		warnings = append(warnings, "no_search_evidence")
	}

	engagement := math.Min(100, 20+math.Log10(page.Sessions+1)*30)
	conversion := 20.0
	if page.Sessions > 0 {
		conversion = math.Min(100, (page.Conversions/page.Sessions)*2000)
	}
	if page.Sessions >= 100 && page.Conversions == 0 {
		warnings = append(warnings, "traffic_without_conversion")
	}

	technical := 0.0
	if page.RenderStatus == seo_domain.SeoRenderStatusSuccess {
		technical += 55
	} else {
		blockers = append(blockers, "render_not_success")
	}
	if page.IsSitemap || !page.IsIndex {
		technical += 20
	} else {
		warnings = append(warnings, "missing_from_sitemap")
	}
	if page.LastCrawlTime != "" {
		technical += 15
	}
	if page.QualityScore >= 70 {
		technical += 10
	} else if page.QualityScore < 50 {
		warnings = append(warnings, "low_content_quality")
	}

	freshnessPenalty := math.Min(20, float64(staleSourceCount)*4)
	if staleSourceCount > 0 {
		warnings = append(warnings, "stale_measurement_sources")
	}
	issuePenalty := math.Min(40, float64(page.OpenIssueCount)*8) + freshnessPenalty
	overall := indexScore*0.25 + visibility*0.25 + engagement*0.15 + conversion*0.10 + technical*0.25 - issuePenalty
	overall = math.Round(clamp(overall, 0, 100)*100) / 100
	grade := "F"
	switch {
	case overall >= 85:
		grade = "A"
	case overall >= 70:
		grade = "B"
	case overall >= 55:
		grade = "C"
	case overall >= 40:
		grade = "D"
	}
	seoID := page.SeoDomainID
	return seo_domain.SeoPageScoreSnapshot{
		SeoDomainID: &seoID, CanonicalURL: page.CanonicalURL, RangeKey: rangeKey, MeasuredAt: measuredAt,
		IndexScore: indexScore, VisibilityScore: visibility, EngagementScore: engagement, ConversionScore: conversion,
		TechnicalScore: technical, IssuePenalty: issuePenalty, OverallScore: overall, Grade: grade,
		Blockers: mustJSON(blockers), Warnings: mustJSON(warnings), SourceFreshness: freshness,
	}
}

func measurementImportChecksum(source string, req *dto.SeoMeasurementImportRequest) (string, error) {
	provided := strings.ToLower(strings.TrimSpace(req.Checksum))
	if len(provided) > 128 {
		return "", measurementValidationErrorf("checksum is limited to 128 characters")
	}
	if provided != "" {
		return provided, nil
	}
	payload, err := json.Marshal(struct {
		Source     string                   `json:"source"`
		RangeStart string                   `json:"rangeStart"`
		RangeEnd   string                   `json:"rangeEnd"`
		Rows       []map[string]interface{} `json:"rows"`
	}{
		Source: source, RangeStart: strings.TrimSpace(req.RangeStart), RangeEnd: strings.TrimSpace(req.RangeEnd), Rows: req.Rows,
	})
	if err != nil {
		return "", fmt.Errorf("cannot calculate import checksum: %w", err)
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func measurementImportResponseFromRun(run *seo_domain.SeoMeasurementImportRun, idempotent bool) *dto.SeoMeasurementImportResponse {
	if run == nil {
		return nil
	}
	errors := []string{}
	if strings.TrimSpace(run.ErrorMessage) != "" {
		errors = append(errors, run.ErrorMessage)
	}
	return &dto.SeoMeasurementImportResponse{
		RunID: run.ID, Source: run.Source, Status: run.Status, RowCount: run.RowCount,
		AcceptedCount: run.AcceptedCount, RejectedCount: run.RejectedCount, Errors: errors, Idempotent: idempotent,
	}
}

func validateCanonicalMeasurementURL(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("canonicalUrl is required")
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("canonicalUrl must be an absolute URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("canonicalUrl scheme must be http or https")
	}
	return nil
}

func optionalMeasurementTime(value string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	parsed, err := parseMeasurementTime(value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func measurementTimeDefault(value string, fallback time.Time) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return fallback, nil
	}
	return parseMeasurementTime(value)
}

func requiredMeasurementDate(value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Time{}, fmt.Errorf("date is required")
	}
	parsed, err := parseMeasurementTime(value)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC), nil
}

func parseMeasurementTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported time format %q", value)
}

func mapString(row map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if value, ok := row[key]; ok && value != nil {
			return strings.TrimSpace(fmt.Sprint(value))
		}
	}
	return ""
}

func mapStringDefault(row map[string]interface{}, fallback string, keys ...string) string {
	if value := mapString(row, keys...); value != "" {
		return value
	}
	return fallback
}

func mapFloat(row map[string]interface{}, keys ...string) float64 {
	for _, key := range keys {
		value, ok := row[key]
		if !ok || value == nil {
			continue
		}
		switch typed := value.(type) {
		case float64:
			return typed
		case float32:
			return float64(typed)
		case int:
			return float64(typed)
		case int64:
			return float64(typed)
		case json.Number:
			n, _ := typed.Float64()
			return n
		default:
			n, _ := strconv.ParseFloat(strings.TrimSpace(fmt.Sprint(value)), 64)
			return n
		}
	}
	return 0
}

func nonNegative(value float64) float64 {
	if value < 0 {
		return 0
	}
	return value
}
func clamp(value, minimum, maximum float64) float64 {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}
