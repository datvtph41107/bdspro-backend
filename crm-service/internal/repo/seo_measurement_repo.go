package repo

import (
	"context"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
)

// SeoMeasurementRepo is the repository boundary for SEO Measurement & Feedback.
//
// This is intentionally separate from SeoDomainRepo:
//   - SeoDomainRepo owns publishing/lifecycle/render state.
//   - SeoMeasurementRepo owns performance, index, crawl, ranking and issues.
//
// This separation keeps SEO page creation independent from analytics providers.
type SeoMeasurementRepo interface {
	UpsertSearchPerformanceDaily(ctx context.Context, items []seo_domain.SeoSearchPerformanceDaily) error
	UpsertURLInspectionSnapshots(ctx context.Context, items []seo_domain.SeoURLInspectionSnapshot) error
	UpsertTrafficDaily(ctx context.Context, items []seo_domain.SeoTrafficDaily) error
	UpsertKeywordTargets(ctx context.Context, items []seo_domain.SeoKeywordTarget) error
	CreateRankingSnapshots(ctx context.Context, items []seo_domain.SeoRankingSnapshot) error
	CreateCrawlEvents(ctx context.Context, items []seo_domain.SeoCrawlEvent) error

	GetOverview(ctx context.Context, req *dto.SeoMeasurementRangeRequest) (*dto.SeoMeasurementOverviewResponse, error)
	ListPagePerformance(ctx context.Context, req *dto.SeoPagePerformanceListRequest) ([]dto.SeoPagePerformanceRow, int64, error)
	ListQueryPerformance(ctx context.Context, req *dto.SeoQueryPerformanceListRequest) ([]dto.SeoQueryPerformanceRow, int64, error)
	ListIndexCoverage(ctx context.Context, req *dto.SeoIndexCoverageListRequest) ([]dto.SeoIndexCoverageRow, int64, error)
	ListTraffic(ctx context.Context, req *dto.SeoTrafficListRequest) ([]dto.SeoTrafficRow, int64, error)
	ListRanking(ctx context.Context, req *dto.SeoRankingListRequest) ([]dto.SeoRankingRow, int64, error)
	ListCrawlEvents(ctx context.Context, req *dto.SeoCrawlEventListRequest) ([]dto.SeoCrawlEventRow, int64, error)
	ListIssues(ctx context.Context, req *dto.SeoIssueListRequest) ([]dto.SeoIssueRow, int64, error)

	CreateIssues(ctx context.Context, items []seo_domain.SeoIssue) error
	ResolveIssue(ctx context.Context, id uint64, status string) error
	ReconcileIssues(ctx context.Context, canonicalURL string, items []seo_domain.SeoIssue) (created uint32, resolved uint32, err error)

	GetReusableImportRunByChecksum(ctx context.Context, source string, checksum string) (*seo_domain.SeoMeasurementImportRun, error)
	CreateImportRun(ctx context.Context, run *seo_domain.SeoMeasurementImportRun) error
	CompleteImportRun(ctx context.Context, run *seo_domain.SeoMeasurementImportRun) error
	UpsertSourceState(ctx context.Context, state *seo_domain.SeoMeasurementSourceState) error
	ListSourceHealth(ctx context.Context) ([]dto.SeoMeasurementSourceHealth, error)
	ListImportRuns(ctx context.Context, req *dto.SeoMeasurementImportRunListRequest) ([]dto.SeoMeasurementImportRunRow, int64, error)
	UpsertPageScores(ctx context.Context, items []seo_domain.SeoPageScoreSnapshot) error
	ListPageScores(ctx context.Context, req *dto.SeoPageScoreListRequest) ([]dto.SeoPageScoreRow, int64, error)
}
