package postgre

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	impl "common/db"
	_commonDto "common/domain/dto"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/repo"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SeoMeasurementPostgres struct {
	db *gorm.DB
}

func NewSeoMeasurementPostgres(db *gorm.DB) repo.SeoMeasurementRepo {
	return &SeoMeasurementPostgres{db: db}
}

func (r *SeoMeasurementPostgres) UpsertSearchPerformanceDaily(ctx context.Context, items []seo_domain.SeoSearchPerformanceDaily) error {
	if len(items) == 0 {
		return nil
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "canonical_url"}, {Name: "query"}, {Name: "device"}, {Name: "country"}, {Name: "search_type"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"seo_domain_id", "clicks", "impressions", "ctr", "position", "data_source", "source_run_id", "imported_at", "updated_at"}),
	}).Create(&items).Error
}

func (r *SeoMeasurementPostgres) UpsertURLInspectionSnapshots(ctx context.Context, items []seo_domain.SeoURLInspectionSnapshot) error {
	if len(items) == 0 {
		return nil
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "canonical_url"}, {Name: "data_source"}},
		DoUpdates: clause.AssignmentColumns([]string{"seo_domain_id", "inspection_status", "index_status", "coverage_state", "robots_txt_state", "indexing_state", "google_canonical", "user_canonical", "last_crawl_time", "page_fetch_state", "verdict", "raw_json", "data_source", "source_run_id", "inspected_at", "updated_at"}),
	}).Create(&items).Error
}

func (r *SeoMeasurementPostgres) UpsertTrafficDaily(ctx context.Context, items []seo_domain.SeoTrafficDaily) error {
	if len(items) == 0 {
		return nil
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "canonical_url"}, {Name: "date"}, {Name: "source"}, {Name: "medium"}, {Name: "device"}, {Name: "country"}},
		DoUpdates: clause.AssignmentColumns([]string{"seo_domain_id", "sessions", "users", "active_users", "page_views", "engagement_seconds", "conversions", "data_source", "source_run_id", "imported_at", "updated_at"}),
	}).Create(&items).Error
}

func (r *SeoMeasurementPostgres) UpsertKeywordTargets(ctx context.Context, items []seo_domain.SeoKeywordTarget) error {
	if len(items) == 0 {
		return nil
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "seo_domain_id"}, {Name: "keyword"}},
		DoUpdates: clause.AssignmentColumns([]string{"intent", "priority", "target_url", "status", "note", "updated_at"}),
	}).Create(&items).Error
}

func (r *SeoMeasurementPostgres) CreateRankingSnapshots(ctx context.Context, items []seo_domain.SeoRankingSnapshot) error {
	if len(items) == 0 {
		return nil
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).Create(&items).Error
}

func (r *SeoMeasurementPostgres) CreateCrawlEvents(ctx context.Context, items []seo_domain.SeoCrawlEvent) error {
	if len(items) == 0 {
		return nil
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).Create(&items).Error
}

func (r *SeoMeasurementPostgres) GetOverview(ctx context.Context, req *dto.SeoMeasurementRangeRequest) (*dto.SeoMeasurementOverviewResponse, error) {
	from, to := measurementDateRange(reqRange(req))

	var kpi struct {
		Clicks      float64
		Impressions float64
		CTR         float64
		Position    float64
		Sessions    float64
		Conversions float64
		Indexed     float64
		OpenIssues  float64
	}

	db := impl.GetDB(ctx, r.db).WithContext(ctx)
	_ = db.Raw(`
		SELECT
			COALESCE(SUM(sp.clicks), 0) AS clicks,
			COALESCE(SUM(sp.impressions), 0) AS impressions,
			CASE WHEN COALESCE(SUM(sp.impressions), 0) > 0 THEN COALESCE(SUM(sp.clicks), 0) / COALESCE(SUM(sp.impressions), 0) ELSE 0 END AS ctr,
			COALESCE(AVG(NULLIF(sp.position, 0)), 0) AS position,
			COALESCE((SELECT SUM(sessions) FROM seo_traffic_daily WHERE date >= ? AND date <= ?), 0) AS sessions,
			COALESCE((SELECT SUM(conversions) FROM seo_traffic_daily WHERE date >= ? AND date <= ?), 0) AS conversions,
			COALESCE((SELECT COUNT(*) FROM seo_url_inspection_snapshot WHERE index_status = 'indexed'), 0) AS indexed,
			COALESCE((SELECT COUNT(*) FROM seo_issue WHERE status = 'open'), 0) AS open_issues
		FROM seo_search_performance_daily sp
		WHERE sp.date >= ? AND sp.date <= ?
	`, from, to, from, to, from, to).Scan(&kpi).Error

	trend, err := r.loadMeasurementTrend(ctx, from, to)
	if err != nil {
		return nil, err
	}
	issues, err := r.loadIssueSummary(ctx)
	if err != nil {
		return nil, err
	}
	topPages, _, err := r.ListPagePerformance(ctx, &dto.SeoPagePerformanceListRequest{Pagable: _commonDto.Pagable{Page: 0, Size: 10}, Range: reqRange(req)})
	if err != nil {
		return nil, err
	}
	needAttention, _, err := r.ListIssues(ctx, &dto.SeoIssueListRequest{Pagable: _commonDto.Pagable{Page: 0, Size: 10}, Status: seo_domain.SeoIssueStatusOpen, Range: reqRange(req)})
	if err != nil {
		return nil, err
	}

	return &dto.SeoMeasurementOverviewResponse{
		KPIs: []dto.SeoMeasurementKPI{
			{Key: "clicks", Label: "Clicks", Value: kpi.Clicks, Unit: "count"},
			{Key: "impressions", Label: "Impressions", Value: kpi.Impressions, Unit: "count"},
			{Key: "ctr", Label: "CTR", Value: kpi.CTR, Unit: "ratio", Severity: severityForCTR(kpi.CTR)},
			{Key: "avg_position", Label: "Average Position", Value: kpi.Position, Unit: "position"},
			{Key: "sessions", Label: "Organic Sessions", Value: kpi.Sessions, Unit: "count"},
			{Key: "conversions", Label: "Conversions", Value: kpi.Conversions, Unit: "count"},
			{Key: "indexed_pages", Label: "Indexed Pages", Value: kpi.Indexed, Unit: "count"},
			{Key: "open_issues", Label: "Open Issues", Value: kpi.OpenIssues, Unit: "count", Severity: severityForOpenIssues(kpi.OpenIssues)},
		},
		Trend:         trend,
		IssueSummary:  issues,
		TopPages:      topPages,
		NeedAttention: needAttention,
	}, nil
}

func (r *SeoMeasurementPostgres) ListPagePerformance(ctx context.Context, req *dto.SeoPagePerformanceListRequest) ([]dto.SeoPagePerformanceRow, int64, error) {
	if req == nil {
		req = &dto.SeoPagePerformanceListRequest{}
	}
	req.Normalize()
	from, to := measurementDateRange(req.Range)

	var total int64
	base := impl.GetDB(ctx, r.db).WithContext(ctx).Table("seo_domain sd")
	base = applyPageFilters(base, req)
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []dto.SeoPagePerformanceRow
	err := applyPageFilters(impl.GetDB(ctx, r.db).WithContext(ctx).Table("seo_domain sd"), req).
		Select(`
			sd.id AS seo_domain_id,
			sd.canonical_url,
			sd.slug,
			sd.title,
			sd.page_status,
			sd.render_status,
			sd.is_site_map AS is_sitemap,
			sd.is_index,
			sd.quality_score,
			COALESCE(sp.clicks, 0) AS clicks,
			COALESCE(sp.impressions, 0) AS impressions,
			CASE WHEN COALESCE(sp.impressions, 0) > 0 THEN COALESCE(sp.clicks, 0) / COALESCE(sp.impressions, 0) ELSE 0 END AS ctr,
			COALESCE(sp.position, 0) AS position,
			COALESCE(td.sessions, 0) AS sessions,
			COALESCE(td.conversions, 0) AS conversions,
			COALESCE(ui.index_status, 'unknown') AS index_status,
			COALESCE(TO_CHAR(ui.last_crawl_time, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '') AS last_crawl_time,
			COALESCE(issue.open_issue_count, 0) AS open_issue_count
		`).
		Joins(`LEFT JOIN (
			SELECT canonical_url, SUM(clicks) AS clicks, SUM(impressions) AS impressions, AVG(NULLIF(position,0)) AS position
			FROM seo_search_performance_daily
			WHERE date >= ? AND date <= ?
			GROUP BY canonical_url
		) sp ON sp.canonical_url = sd.canonical_url`, from, to).
		Joins(`LEFT JOIN (
			SELECT canonical_url, SUM(sessions) AS sessions, SUM(conversions) AS conversions
			FROM seo_traffic_daily
			WHERE date >= ? AND date <= ?
			GROUP BY canonical_url
		) td ON td.canonical_url = sd.canonical_url`, from, to).
		Joins(`LEFT JOIN (
			SELECT DISTINCT ON (canonical_url) canonical_url, index_status, last_crawl_time
			FROM seo_url_inspection_snapshot
			ORDER BY canonical_url, inspected_at DESC
		) ui ON ui.canonical_url = sd.canonical_url`).
		Joins(`LEFT JOIN (
			SELECT canonical_url, COUNT(*) AS open_issue_count
			FROM seo_issue WHERE status = 'open'
			GROUP BY canonical_url
		) issue ON issue.canonical_url = sd.canonical_url`).
		Order("COALESCE(sp.impressions, 0) DESC, sd.updated_at DESC").
		Offset(req.GetOffset()).Limit(req.GetLimitAdmin()).Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *SeoMeasurementPostgres) ListQueryPerformance(ctx context.Context, req *dto.SeoQueryPerformanceListRequest) ([]dto.SeoQueryPerformanceRow, int64, error) {
	if req == nil {
		req = &dto.SeoQueryPerformanceListRequest{}
	}
	req.Normalize()
	from, to := measurementDateRange(req.Range)
	query := impl.GetDB(ctx, r.db).WithContext(ctx).Table("seo_search_performance_daily").Where("date >= ? AND date <= ?", from, to)
	if strings.TrimSpace(req.Query) != "" {
		query = query.Where("LOWER(query) LIKE ?", "%"+strings.ToLower(strings.TrimSpace(req.Query))+"%")
	}
	if strings.TrimSpace(req.CanonicalURL) != "" {
		query = query.Where("canonical_url = ?", strings.TrimSpace(req.CanonicalURL))
	}
	if strings.TrimSpace(req.Device) != "" {
		query = query.Where("device = ?", strings.TrimSpace(req.Device))
	}
	if strings.TrimSpace(req.Country) != "" {
		query = query.Where("country = ?", strings.TrimSpace(req.Country))
	}
	var grouped []dto.SeoQueryPerformanceRow
	err := query.Select(`query, canonical_url, device, country, SUM(clicks) AS clicks, SUM(impressions) AS impressions, CASE WHEN SUM(impressions) > 0 THEN SUM(clicks)/SUM(impressions) ELSE 0 END AS ctr, AVG(NULLIF(position,0)) AS position`).Group("query, canonical_url, device, country").Order("impressions DESC").Offset(req.GetOffset()).Limit(req.GetLimitAdmin()).Scan(&grouped).Error
	if err != nil {
		return nil, 0, err
	}
	var countRows []struct{ Count int64 }
	countQuery := impl.GetDB(ctx, r.db).WithContext(ctx).Table("seo_search_performance_daily").Where("date >= ? AND date <= ?", from, to)
	if strings.TrimSpace(req.Query) != "" {
		countQuery = countQuery.Where("LOWER(query) LIKE ?", "%"+strings.ToLower(strings.TrimSpace(req.Query))+"%")
	}
	if strings.TrimSpace(req.CanonicalURL) != "" {
		countQuery = countQuery.Where("canonical_url = ?", strings.TrimSpace(req.CanonicalURL))
	}
	if strings.TrimSpace(req.Device) != "" {
		countQuery = countQuery.Where("device = ?", strings.TrimSpace(req.Device))
	}
	if strings.TrimSpace(req.Country) != "" {
		countQuery = countQuery.Where("country = ?", strings.TrimSpace(req.Country))
	}
	_ = countQuery.Select("COUNT(*) AS count").Group("query, canonical_url, device, country").Scan(&countRows).Error
	return grouped, int64(len(countRows)), nil
}

func (r *SeoMeasurementPostgres) ListIndexCoverage(ctx context.Context, req *dto.SeoIndexCoverageListRequest) ([]dto.SeoIndexCoverageRow, int64, error) {
	if req == nil {
		req = &dto.SeoIndexCoverageListRequest{}
	}
	req.Normalize()
	query := impl.GetDB(ctx, r.db).WithContext(ctx).Table("seo_url_inspection_snapshot")
	if strings.TrimSpace(req.CanonicalURL) != "" {
		query = query.Where("canonical_url = ?", strings.TrimSpace(req.CanonicalURL))
	}
	if strings.TrimSpace(req.IndexStatus) != "" {
		query = query.Where("index_status = ?", strings.TrimSpace(req.IndexStatus))
	}
	if strings.TrimSpace(req.Verdict) != "" {
		query = query.Where("verdict = ?", strings.TrimSpace(req.Verdict))
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []dto.SeoIndexCoverageRow
	err := query.Select(`COALESCE(seo_domain_id,0) AS seo_domain_id, canonical_url, inspection_status, index_status, coverage_state, robots_txt_state, google_canonical, user_canonical, COALESCE(TO_CHAR(last_crawl_time, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '') AS last_crawl_time, verdict, COALESCE(TO_CHAR(inspected_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '') AS inspected_at`).Order("inspected_at DESC").Offset(req.GetOffset()).Limit(req.GetLimitAdmin()).Scan(&rows).Error
	return rows, total, err
}

func (r *SeoMeasurementPostgres) ListTraffic(ctx context.Context, req *dto.SeoTrafficListRequest) ([]dto.SeoTrafficRow, int64, error) {
	if req == nil {
		req = &dto.SeoTrafficListRequest{}
	}
	req.Normalize()
	from, to := measurementDateRange(req.Range)
	query := impl.GetDB(ctx, r.db).WithContext(ctx).Table("seo_traffic_daily").Where("date >= ? AND date <= ?", from, to)
	if strings.TrimSpace(req.CanonicalURL) != "" {
		query = query.Where("canonical_url = ?", strings.TrimSpace(req.CanonicalURL))
	}
	if strings.TrimSpace(req.Source) != "" {
		query = query.Where("source = ?", strings.TrimSpace(req.Source))
	}
	if strings.TrimSpace(req.Medium) != "" {
		query = query.Where("medium = ?", strings.TrimSpace(req.Medium))
	}
	if strings.TrimSpace(req.Device) != "" {
		query = query.Where("device = ?", strings.TrimSpace(req.Device))
	}
	if strings.TrimSpace(req.Country) != "" {
		query = query.Where("country = ?", strings.TrimSpace(req.Country))
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []dto.SeoTrafficRow
	err := query.Select(`canonical_url, TO_CHAR(date, 'YYYY-MM-DD') AS date, sessions, users, active_users, page_views, engagement_seconds, conversions, source, medium, device, country`).Order("date DESC, sessions DESC").Offset(req.GetOffset()).Limit(req.GetLimitAdmin()).Scan(&rows).Error
	return rows, total, err
}

func (r *SeoMeasurementPostgres) ListRanking(ctx context.Context, req *dto.SeoRankingListRequest) ([]dto.SeoRankingRow, int64, error) {
	if req == nil {
		req = &dto.SeoRankingListRequest{}
	}
	req.Normalize()
	from, to := measurementDateRange(req.Range)
	query := impl.GetDB(ctx, r.db).WithContext(ctx).Table("seo_ranking_snapshot").Where("checked_at >= ? AND checked_at < ?", from, to.AddDate(0, 0, 1))
	if strings.TrimSpace(req.Keyword) != "" {
		query = query.Where("LOWER(keyword) LIKE ?", "%"+strings.ToLower(strings.TrimSpace(req.Keyword))+"%")
	}
	if strings.TrimSpace(req.CanonicalURL) != "" {
		query = query.Where("canonical_url = ?", strings.TrimSpace(req.CanonicalURL))
	}
	if strings.TrimSpace(req.Device) != "" {
		query = query.Where("device = ?", strings.TrimSpace(req.Device))
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []dto.SeoRankingRow
	err := query.Select(`keyword, canonical_url, position, previous_position, position_delta, device, location, language, COALESCE(TO_CHAR(checked_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '') AS checked_at`).Order("checked_at DESC, position ASC").Offset(req.GetOffset()).Limit(req.GetLimitAdmin()).Scan(&rows).Error
	return rows, total, err
}

func (r *SeoMeasurementPostgres) ListCrawlEvents(ctx context.Context, req *dto.SeoCrawlEventListRequest) ([]dto.SeoCrawlEventRow, int64, error) {
	if req == nil {
		req = &dto.SeoCrawlEventListRequest{}
	}
	req.Normalize()
	from, to := measurementDateRange(req.Range)
	query := impl.GetDB(ctx, r.db).WithContext(ctx).Table("seo_crawl_event").Where("crawled_at >= ? AND crawled_at < ?", from, to.AddDate(0, 0, 1))
	if strings.TrimSpace(req.CanonicalURL) != "" {
		query = query.Where("canonical_url = ?", strings.TrimSpace(req.CanonicalURL))
	}
	if strings.TrimSpace(req.BotName) != "" {
		query = query.Where("bot_name = ?", strings.TrimSpace(req.BotName))
	}
	if req.StatusCode > 0 {
		query = query.Where("status_code = ?", req.StatusCode)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []dto.SeoCrawlEventRow
	err := query.Select(`id, COALESCE(seo_domain_id,0) AS seo_domain_id, canonical_url, bot_name, user_agent, status_code, method, response_time_ms, COALESCE(TO_CHAR(crawled_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '') AS crawled_at`).Order("crawled_at DESC").Offset(req.GetOffset()).Limit(req.GetLimitAdmin()).Scan(&rows).Error
	return rows, total, err
}

func (r *SeoMeasurementPostgres) ListIssues(ctx context.Context, req *dto.SeoIssueListRequest) ([]dto.SeoIssueRow, int64, error) {
	if req == nil {
		req = &dto.SeoIssueListRequest{}
	}
	req.Normalize()
	query := impl.GetDB(ctx, r.db).WithContext(ctx).Table("seo_issue")
	if strings.TrimSpace(req.CanonicalURL) != "" {
		query = query.Where("canonical_url = ?", strings.TrimSpace(req.CanonicalURL))
	}
	if strings.TrimSpace(req.IssueType) != "" {
		query = query.Where("issue_type = ?", strings.TrimSpace(req.IssueType))
	}
	if strings.TrimSpace(req.Severity) != "" {
		query = query.Where("severity = ?", strings.TrimSpace(req.Severity))
	}
	if strings.TrimSpace(req.Status) != "" {
		query = query.Where("status = ?", strings.TrimSpace(req.Status))
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []dto.SeoIssueRow
	err := query.Select(`id, COALESCE(seo_domain_id,0) AS seo_domain_id, canonical_url, issue_type, severity, status, title, message, evidence::text AS evidence, recommendation, COALESCE(TO_CHAR(detected_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '') AS detected_at, COALESCE(TO_CHAR(resolved_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '') AS resolved_at`).Order("CASE severity WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 ELSE 5 END, detected_at DESC").Offset(req.GetOffset()).Limit(req.GetLimitAdmin()).Scan(&rows).Error
	return rows, total, err
}

func (r *SeoMeasurementPostgres) CreateIssues(ctx context.Context, items []seo_domain.SeoIssue) error {
	if len(items) == 0 {
		return nil
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "canonical_url"}, {Name: "issue_type"}, {Name: "status"}},
		DoNothing: true,
	}).Create(&items).Error
}

func (r *SeoMeasurementPostgres) ResolveIssue(ctx context.Context, id uint64, status string) error {
	status = strings.TrimSpace(status)
	if status == "" {
		status = seo_domain.SeoIssueStatusResolved
	}
	now := time.Now()
	return impl.GetDB(ctx, r.db).WithContext(ctx).Model(&seo_domain.SeoIssue{}).Where("id = ?", id).Updates(map[string]interface{}{"status": status, "resolved_at": &now, "updated_at": now}).Error
}

func (r *SeoMeasurementPostgres) loadMeasurementTrend(ctx context.Context, from time.Time, to time.Time) ([]dto.SeoMeasurementTrendPoint, error) {
	var rows []dto.SeoMeasurementTrendPoint
	err := impl.GetDB(ctx, r.db).WithContext(ctx).Raw(`
		SELECT
			TO_CHAR(d.day, 'YYYY-MM-DD') AS date,
			COALESCE(sp.clicks, 0) AS clicks,
			COALESCE(sp.impressions, 0) AS impressions,
			CASE WHEN COALESCE(sp.impressions, 0) > 0 THEN COALESCE(sp.clicks, 0)/COALESCE(sp.impressions, 0) ELSE 0 END AS ctr,
			COALESCE(sp.position, 0) AS position,
			COALESCE(td.sessions, 0) AS sessions,
			COALESCE(td.conversions, 0) AS conversions
		FROM generate_series(?::date, ?::date, interval '1 day') d(day)
		LEFT JOIN (
			SELECT date, SUM(clicks) AS clicks, SUM(impressions) AS impressions, AVG(NULLIF(position,0)) AS position
			FROM seo_search_performance_daily GROUP BY date
		) sp ON sp.date = d.day
		LEFT JOIN (
			SELECT date, SUM(sessions) AS sessions, SUM(conversions) AS conversions
			FROM seo_traffic_daily GROUP BY date
		) td ON td.date = d.day
		ORDER BY d.day ASC
	`, from.Format("2006-01-02"), to.Format("2006-01-02")).Scan(&rows).Error
	return rows, err
}

func (r *SeoMeasurementPostgres) loadIssueSummary(ctx context.Context) ([]dto.SeoIssueSummary, error) {
	var rows []dto.SeoIssueSummary
	err := impl.GetDB(ctx, r.db).WithContext(ctx).Table("seo_issue").Select("issue_type, severity, COUNT(*) AS open_count").Where("status = ?", seo_domain.SeoIssueStatusOpen).Group("issue_type, severity").Order("open_count DESC").Scan(&rows).Error
	return rows, err
}

func applyPageFilters(query *gorm.DB, req *dto.SeoPagePerformanceListRequest) *gorm.DB {
	if req == nil {
		return query
	}
	if strings.TrimSpace(req.Keyword) != "" {
		kw := "%" + strings.ToLower(strings.TrimSpace(req.Keyword)) + "%"
		query = query.Where("LOWER(sd.title) LIKE ? OR LOWER(sd.slug) LIKE ? OR LOWER(sd.canonical_url) LIKE ?", kw, kw, kw)
	}
	if strings.TrimSpace(req.PageStatus) != "" {
		query = query.Where("sd.page_status = ?", strings.TrimSpace(req.PageStatus))
	}
	if strings.TrimSpace(req.RenderStatus) != "" {
		query = query.Where("sd.render_status = ?", strings.TrimSpace(req.RenderStatus))
	}
	return query
}

func measurementDateRange(rangeValue string) (time.Time, time.Time) {
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

func reqRange(req *dto.SeoMeasurementRangeRequest) string {
	if req == nil {
		return "28d"
	}
	return req.Range
}

func severityForCTR(ctr float64) string {
	if ctr < 0.01 {
		return seo_domain.SeoIssueSeverityHigh
	}
	if ctr < 0.03 {
		return seo_domain.SeoIssueSeverityMedium
	}
	return seo_domain.SeoIssueSeverityInfo
}

func severityForOpenIssues(open float64) string {
	if open >= 20 {
		return seo_domain.SeoIssueSeverityHigh
	}
	if open > 0 {
		return seo_domain.SeoIssueSeverityMedium
	}
	return seo_domain.SeoIssueSeverityInfo
}

func (r *SeoMeasurementPostgres) ReconcileIssues(ctx context.Context, canonicalURL string, items []seo_domain.SeoIssue) (uint32, uint32, error) {
	db := impl.GetDB(ctx, r.db).WithContext(ctx)
	managedTypes := []string{
		seo_domain.SeoIssueTypeNotRendered,
		seo_domain.SeoIssueTypeNotInSitemap,
		seo_domain.SeoIssueTypeCrawledNotIndexed,
		seo_domain.SeoIssueTypeHighImpressionLowCTR,
		seo_domain.SeoIssueTypeGoodPositionLowCTR,
		seo_domain.SeoIssueTypeTrafficNoConversion,
	}
	var created uint32
	var resolved uint32
	err := db.Transaction(func(tx *gorm.DB) error {
		var existing []seo_domain.SeoIssue
		if err := tx.Where("canonical_url = ? AND status = ? AND issue_type IN ?", canonicalURL, seo_domain.SeoIssueStatusOpen, managedTypes).Find(&existing).Error; err != nil {
			return err
		}
		existingTypes := make(map[string]bool, len(existing))
		for _, issue := range existing {
			existingTypes[issue.IssueType] = true
		}
		activeTypes := make([]string, 0, len(items))
		for i := range items {
			item := &items[i]
			activeTypes = append(activeTypes, item.IssueType)
			if !existingTypes[item.IssueType] {
				created++
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "canonical_url"}, {Name: "issue_type"}, {Name: "status"}},
				DoUpdates: clause.AssignmentColumns([]string{"seo_domain_id", "severity", "title", "message", "evidence", "recommendation", "detected_at", "resolved_at", "updated_at"}),
			}).Create(item).Error; err != nil {
				return err
			}
		}
		query := tx.Model(&seo_domain.SeoIssue{}).
			Where("canonical_url = ? AND status = ? AND issue_type IN ?", canonicalURL, seo_domain.SeoIssueStatusOpen, managedTypes)
		if len(activeTypes) > 0 {
			query = query.Where("issue_type NOT IN ?", activeTypes)
		}
		now := time.Now()
		result := query.Updates(map[string]interface{}{
			"status":      seo_domain.SeoIssueStatusResolved,
			"resolved_at": &now,
			"updated_at":  now,
		})
		if result.Error != nil {
			return result.Error
		}
		resolved = uint32(result.RowsAffected)
		return nil
	})
	return created, resolved, err
}

func (r *SeoMeasurementPostgres) GetReusableImportRunByChecksum(ctx context.Context, source string, checksum string) (*seo_domain.SeoMeasurementImportRun, error) {
	var run seo_domain.SeoMeasurementImportRun
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Where("source = ? AND checksum = ? AND status IN ?", strings.TrimSpace(source), strings.TrimSpace(checksum), []string{
			seo_domain.SeoMeasurementImportStatusRunning,
			seo_domain.SeoMeasurementImportStatusSucceeded,
			seo_domain.SeoMeasurementImportStatusPartial,
		}).
		Order("id DESC").
		First(&run).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	if run.Status == seo_domain.SeoMeasurementImportStatusRunning && time.Since(run.StartedAt) > 2*time.Hour {
		now := time.Now().UTC()
		if err := impl.GetDB(ctx, r.db).WithContext(ctx).Model(&seo_domain.SeoMeasurementImportRun{}).
			Where("id = ? AND status = ?", run.ID, seo_domain.SeoMeasurementImportStatusRunning).
			Updates(map[string]interface{}{
				"status": seo_domain.SeoMeasurementImportStatusFailed, "error_message": "expired stale running import",
				"finished_at": now, "updated_at": now,
			}).Error; err != nil {
			return nil, err
		}
		return nil, nil
	}
	return &run, nil
}

func (r *SeoMeasurementPostgres) CreateImportRun(ctx context.Context, run *seo_domain.SeoMeasurementImportRun) error {
	if run == nil {
		return nil
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).Create(run).Error
}

func (r *SeoMeasurementPostgres) CompleteImportRun(ctx context.Context, run *seo_domain.SeoMeasurementImportRun) error {
	if run == nil || run.ID == 0 {
		return nil
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).Model(&seo_domain.SeoMeasurementImportRun{}).
		Where("id = ?", run.ID).
		Updates(map[string]interface{}{
			"status":         run.Status,
			"row_count":      run.RowCount,
			"accepted_count": run.AcceptedCount,
			"rejected_count": run.RejectedCount,
			"checksum":       run.Checksum,
			"error_message":  run.ErrorMessage,
			"metadata":       run.Metadata,
			"finished_at":    run.FinishedAt,
			"updated_at":     time.Now(),
		}).Error
}

func (r *SeoMeasurementPostgres) UpsertSourceState(ctx context.Context, state *seo_domain.SeoMeasurementSourceState) error {
	if state == nil || strings.TrimSpace(state.Source) == "" {
		return nil
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "source"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"status":               gorm.Expr("EXCLUDED.status"),
			"last_run_id":          gorm.Expr("EXCLUDED.last_run_id"),
			"last_attempt_at":      gorm.Expr("EXCLUDED.last_attempt_at"),
			"last_success_at":      gorm.Expr("COALESCE(EXCLUDED.last_success_at, seo_measurement_source_state.last_success_at)"),
			"last_row_count":       gorm.Expr("CASE WHEN EXCLUDED.last_success_at IS NULL THEN seo_measurement_source_state.last_row_count ELSE EXCLUDED.last_row_count END"),
			"stale_after_hours":    gorm.Expr("EXCLUDED.stale_after_hours"),
			"consecutive_failures": gorm.Expr("CASE WHEN EXCLUDED.last_success_at IS NULL THEN seo_measurement_source_state.consecutive_failures + 1 ELSE 0 END"),
			"last_error":           gorm.Expr("CASE WHEN EXCLUDED.status = 'partial' THEN EXCLUDED.last_error WHEN EXCLUDED.last_success_at IS NULL THEN EXCLUDED.last_error ELSE '' END"),
			"metadata":             gorm.Expr("EXCLUDED.metadata"),
			"updated_at":           gorm.Expr("NOW()"),
		}),
	}).Create(state).Error
}

func (r *SeoMeasurementPostgres) ListSourceHealth(ctx context.Context) ([]dto.SeoMeasurementSourceHealth, error) {
	var rows []struct {
		Source              string
		Status              string
		LastRunID           uint64
		LastAttemptAt       *time.Time
		LastSuccessAt       *time.Time
		LastRowCount        int64
		StaleAfterHours     int
		ConsecutiveFailures int
		LastError           string
	}
	if err := impl.GetDB(ctx, r.db).WithContext(ctx).Table("seo_measurement_source_state").Order("source ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	now := time.Now()
	out := make([]dto.SeoMeasurementSourceHealth, 0, len(rows))
	for _, row := range rows {
		item := dto.SeoMeasurementSourceHealth{
			Source:              row.Source,
			Status:              row.Status,
			LastRunID:           row.LastRunID,
			LastRowCount:        row.LastRowCount,
			StaleAfterHours:     row.StaleAfterHours,
			ConsecutiveFailures: row.ConsecutiveFailures,
			LastError:           row.LastError,
		}
		if row.LastAttemptAt != nil {
			item.LastAttemptAt = row.LastAttemptAt.UTC().Format(time.RFC3339)
		}
		if row.LastSuccessAt != nil {
			item.LastSuccessAt = row.LastSuccessAt.UTC().Format(time.RFC3339)
			item.AgeHours = int64(now.Sub(*row.LastSuccessAt).Hours())
			item.Fresh = item.AgeHours <= int64(row.StaleAfterHours)
		}
		if row.LastSuccessAt == nil {
			item.Status = "never_imported"
		}
		if row.LastSuccessAt != nil && !item.Fresh && item.Status != seo_domain.SeoMeasurementImportStatusFailed && item.Status != seo_domain.SeoMeasurementImportStatusRunning {
			item.Status = "stale"
		}
		out = append(out, item)
	}
	return out, nil
}

func (r *SeoMeasurementPostgres) ListImportRuns(ctx context.Context, req *dto.SeoMeasurementImportRunListRequest) ([]dto.SeoMeasurementImportRunRow, int64, error) {
	if req == nil {
		req = &dto.SeoMeasurementImportRunListRequest{}
	}
	size := int(req.Size)
	if size <= 0 {
		size = 20
	}
	if size > 200 {
		size = 200
	}
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * size
	query := impl.GetDB(ctx, r.db).WithContext(ctx).Table("seo_measurement_import_run")
	if strings.TrimSpace(req.Source) != "" {
		query = query.Where("source = ?", strings.TrimSpace(req.Source))
	}
	if strings.TrimSpace(req.Status) != "" {
		query = query.Where("status = ?", strings.TrimSpace(req.Status))
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	rows := make([]dto.SeoMeasurementImportRunRow, 0)
	err := query.Select(`id, source, status, requested_by, row_count, accepted_count, rejected_count, checksum, error_message,
		COALESCE(TO_CHAR(started_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '') AS started_at,
		COALESCE(TO_CHAR(finished_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '') AS finished_at`).
		Order("started_at DESC").Offset(offset).Limit(size).Scan(&rows).Error
	return rows, total, err
}

func (r *SeoMeasurementPostgres) UpsertPageScores(ctx context.Context, items []seo_domain.SeoPageScoreSnapshot) error {
	if len(items) == 0 {
		return nil
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).Create(&items).Error
}

func (r *SeoMeasurementPostgres) ListPageScores(ctx context.Context, req *dto.SeoPageScoreListRequest) ([]dto.SeoPageScoreRow, int64, error) {
	if req == nil {
		req = &dto.SeoPageScoreListRequest{}
	}
	size := int(req.Size)
	if size <= 0 {
		size = 20
	}
	if size > 200 {
		size = 200
	}
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * size
	latest := impl.GetDB(ctx, r.db).WithContext(ctx).Raw(`
		SELECT DISTINCT ON (canonical_url)
			COALESCE(seo_domain_id,0) AS seo_domain_id, canonical_url, range_key,
			TO_CHAR(measured_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS measured_at,
			index_score, visibility_score, engagement_score, conversion_score, technical_score,
			issue_penalty, overall_score, grade, blockers, warnings, source_freshness
		FROM seo_page_score_snapshot
		ORDER BY canonical_url, measured_at DESC
	`)
	base := r.db.Table("(?) scores", latest)
	if strings.TrimSpace(req.CanonicalURL) != "" {
		base = base.Where("canonical_url ILIKE ?", "%"+strings.TrimSpace(req.CanonicalURL)+"%")
	}
	if strings.TrimSpace(req.Grade) != "" {
		base = base.Where("grade = ?", strings.ToUpper(strings.TrimSpace(req.Grade)))
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rawRows []struct {
		SeoDomainID                                                                                               uint64
		CanonicalURL, RangeKey, MeasuredAt                                                                        string
		IndexScore, VisibilityScore, EngagementScore, ConversionScore, TechnicalScore, IssuePenalty, OverallScore float64
		Grade                                                                                                     string
		Blockers, Warnings, SourceFreshness                                                                       []byte
	}
	if err := base.Order("overall_score ASC, canonical_url ASC").Offset(offset).Limit(size).Scan(&rawRows).Error; err != nil {
		return nil, 0, err
	}
	rows := make([]dto.SeoPageScoreRow, 0, len(rawRows))
	for _, row := range rawRows {
		rows = append(rows, dto.SeoPageScoreRow{
			SeoDomainID: row.SeoDomainID, CanonicalURL: row.CanonicalURL, RangeKey: row.RangeKey, MeasuredAt: row.MeasuredAt,
			IndexScore: row.IndexScore, VisibilityScore: row.VisibilityScore, EngagementScore: row.EngagementScore,
			ConversionScore: row.ConversionScore, TechnicalScore: row.TechnicalScore, IssuePenalty: row.IssuePenalty,
			OverallScore: row.OverallScore, Grade: row.Grade,
			Blockers: json.RawMessage(row.Blockers), Warnings: json.RawMessage(row.Warnings), SourceFreshness: json.RawMessage(row.SourceFreshness),
		})
	}
	return rows, total, nil
}
