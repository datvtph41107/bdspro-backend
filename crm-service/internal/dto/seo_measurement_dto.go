package dto

import _dto "common/domain/dto"

// SeoMeasurementRangeRequest is used by dashboard and detector APIs.
// Range accepts: 7d, 14d, 28d, 90d. Empty defaults to 28d.
type SeoMeasurementRangeRequest struct {
	Range string `json:"range" form:"range"`
	Limit uint32 `json:"limit" form:"limit"`
}

type SeoMeasurementGenerateMockRequest struct {
	Range     string `json:"range" form:"range"`
	Limit     uint32 `json:"limit" form:"limit"`
	Overwrite bool   `json:"overwrite" form:"overwrite"`
}

type SeoMeasurementGenerateMockResponse struct {
	Pages          uint32   `json:"pages"`
	SearchRows     uint32   `json:"searchRows"`
	InspectionRows uint32   `json:"inspectionRows"`
	TrafficRows    uint32   `json:"trafficRows"`
	RankingRows    uint32   `json:"rankingRows"`
	CrawlRows      uint32   `json:"crawlRows"`
	KeywordTargets uint32   `json:"keywordTargets"`
	DataSource     string   `json:"dataSource"`
	Messages       []string `json:"messages"`
}

type SeoMeasurementDetectorRequest struct {
	Range string `json:"range" form:"range"`
	Limit uint32 `json:"limit" form:"limit"`
}

type SeoMeasurementDetectorResponse struct {
	Scanned uint32   `json:"scanned"`
	Created uint32   `json:"created"`
	Skipped uint32   `json:"skipped"`
	Errors  []string `json:"errors"`
}

type SeoMeasurementKPI struct {
	Key      string  `json:"key"`
	Label    string  `json:"label"`
	Value    float64 `json:"value"`
	Unit     string  `json:"unit"`
	Delta    float64 `json:"delta"`
	Severity string  `json:"severity"`
}

type SeoMeasurementTrendPoint struct {
	Date        string  `json:"date"`
	Clicks      float64 `json:"clicks"`
	Impressions float64 `json:"impressions"`
	CTR         float64 `json:"ctr"`
	Position    float64 `json:"position"`
	Sessions    float64 `json:"sessions"`
	Conversions float64 `json:"conversions"`
}

type SeoIssueSummary struct {
	IssueType string `json:"issueType"`
	Severity  string `json:"severity"`
	OpenCount uint32 `json:"openCount"`
}

type SeoMeasurementOverviewResponse struct {
	KPIs          []SeoMeasurementKPI        `json:"kpis"`
	Trend         []SeoMeasurementTrendPoint `json:"trend"`
	IssueSummary  []SeoIssueSummary          `json:"issueSummary"`
	TopPages      []SeoPagePerformanceRow    `json:"topPages"`
	NeedAttention []SeoIssueRow              `json:"needAttention"`
}

type SeoPagePerformanceListRequest struct {
	_dto.Pagable
	Keyword      string `json:"keyword" form:"keyword"`
	PageStatus   string `json:"pageStatus" form:"pageStatus"`
	RenderStatus string `json:"renderStatus" form:"renderStatus"`
	IndexStatus  string `json:"indexStatus" form:"indexStatus"`
	IssueStatus  string `json:"issueStatus" form:"issueStatus"`
	Range        string `json:"range" form:"range"`
}

type SeoPagePerformanceRow struct {
	SeoDomainID    uint64  `json:"seoDomainId"`
	CanonicalURL   string  `json:"canonicalUrl"`
	Slug           string  `json:"slug"`
	Title          string  `json:"title"`
	PageStatus     string  `json:"pageStatus"`
	RenderStatus   string  `json:"renderStatus"`
	IsSitemap      bool    `json:"isSitemap"`
	IsIndex        bool    `json:"isIndex"`
	QualityScore   uint32  `json:"qualityScore"`
	Clicks         float64 `json:"clicks"`
	Impressions    float64 `json:"impressions"`
	CTR            float64 `json:"ctr"`
	Position       float64 `json:"position"`
	Sessions       float64 `json:"sessions"`
	Conversions    float64 `json:"conversions"`
	IndexStatus    string  `json:"indexStatus"`
	LastCrawlTime  string  `json:"lastCrawlTime"`
	OpenIssueCount uint32  `json:"openIssueCount"`
}

type SeoQueryPerformanceListRequest struct {
	_dto.Pagable
	Query        string `json:"query" form:"query"`
	CanonicalURL string `json:"canonicalUrl" form:"canonicalUrl"`
	Device       string `json:"device" form:"device"`
	Country      string `json:"country" form:"country"`
	Range        string `json:"range" form:"range"`
}

type SeoQueryPerformanceRow struct {
	Query        string  `json:"query"`
	CanonicalURL string  `json:"canonicalUrl"`
	Device       string  `json:"device"`
	Country      string  `json:"country"`
	Clicks       float64 `json:"clicks"`
	Impressions  float64 `json:"impressions"`
	CTR          float64 `json:"ctr"`
	Position     float64 `json:"position"`
}

type SeoIndexCoverageListRequest struct {
	_dto.Pagable
	CanonicalURL string `json:"canonicalUrl" form:"canonicalUrl"`
	IndexStatus  string `json:"indexStatus" form:"indexStatus"`
	Verdict      string `json:"verdict" form:"verdict"`
	Range        string `json:"range" form:"range"`
}

type SeoIndexCoverageRow struct {
	SeoDomainID      uint64 `json:"seoDomainId"`
	CanonicalURL     string `json:"canonicalUrl"`
	InspectionStatus string `json:"inspectionStatus"`
	IndexStatus      string `json:"indexStatus"`
	CoverageState    string `json:"coverageState"`
	RobotsTxtState   string `json:"robotsTxtState"`
	GoogleCanonical  string `json:"googleCanonical"`
	UserCanonical    string `json:"userCanonical"`
	LastCrawlTime    string `json:"lastCrawlTime"`
	Verdict          string `json:"verdict"`
	InspectedAt      string `json:"inspectedAt"`
}

type SeoTrafficListRequest struct {
	_dto.Pagable
	CanonicalURL string `json:"canonicalUrl" form:"canonicalUrl"`
	Source       string `json:"source" form:"source"`
	Medium       string `json:"medium" form:"medium"`
	Device       string `json:"device" form:"device"`
	Country      string `json:"country" form:"country"`
	Range        string `json:"range" form:"range"`
}

type SeoTrafficRow struct {
	CanonicalURL      string  `json:"canonicalUrl"`
	Date              string  `json:"date"`
	Sessions          float64 `json:"sessions"`
	Users             float64 `json:"users"`
	ActiveUsers       float64 `json:"activeUsers"`
	PageViews         float64 `json:"pageViews"`
	EngagementSeconds float64 `json:"engagementSeconds"`
	Conversions       float64 `json:"conversions"`
	Source            string  `json:"source"`
	Medium            string  `json:"medium"`
	Device            string  `json:"device"`
	Country           string  `json:"country"`
}

type SeoRankingListRequest struct {
	_dto.Pagable
	Keyword      string `json:"keyword" form:"keyword"`
	CanonicalURL string `json:"canonicalUrl" form:"canonicalUrl"`
	Device       string `json:"device" form:"device"`
	Range        string `json:"range" form:"range"`
}

type SeoRankingRow struct {
	Keyword          string  `json:"keyword"`
	CanonicalURL     string  `json:"canonicalUrl"`
	Position         float64 `json:"position"`
	PreviousPosition float64 `json:"previousPosition"`
	PositionDelta    float64 `json:"positionDelta"`
	Device           string  `json:"device"`
	Location         string  `json:"location"`
	Language         string  `json:"language"`
	CheckedAt        string  `json:"checkedAt"`
}

type SeoCrawlEventListRequest struct {
	_dto.Pagable
	CanonicalURL string `json:"canonicalUrl" form:"canonicalUrl"`
	BotName      string `json:"botName" form:"botName"`
	StatusCode   int32  `json:"statusCode" form:"statusCode"`
	Range        string `json:"range" form:"range"`
}

type SeoCrawlEventRow struct {
	ID             uint64 `json:"id"`
	SeoDomainID    uint64 `json:"seoDomainId"`
	CanonicalURL   string `json:"canonicalUrl"`
	BotName        string `json:"botName"`
	UserAgent      string `json:"userAgent"`
	StatusCode     int32  `json:"statusCode"`
	Method         string `json:"method"`
	ResponseTimeMs int64  `json:"responseTimeMs"`
	CrawledAt      string `json:"crawledAt"`
}

type SeoIssueListRequest struct {
	_dto.Pagable
	CanonicalURL string `json:"canonicalUrl" form:"canonicalUrl"`
	IssueType    string `json:"issueType" form:"issueType"`
	Severity     string `json:"severity" form:"severity"`
	Status       string `json:"status" form:"status"`
	Range        string `json:"range" form:"range"`
}

type SeoIssueRow struct {
	ID             uint64 `json:"id"`
	SeoDomainID    uint64 `json:"seoDomainId"`
	CanonicalURL   string `json:"canonicalUrl"`
	IssueType      string `json:"issueType"`
	Severity       string `json:"severity"`
	Status         string `json:"status"`
	Title          string `json:"title"`
	Message        string `json:"message"`
	Evidence       string `json:"evidence"`
	Recommendation string `json:"recommendation"`
	DetectedAt     string `json:"detectedAt"`
	ResolvedAt     string `json:"resolvedAt"`
}

type SeoResolveIssueRequest struct {
	ID     uint64 `json:"id"`
	Status string `json:"status"`
	Note   string `json:"note"`
}
