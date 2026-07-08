package dto

type SeoOperationsDashboardRequest struct {
	Range string `json:"range" form:"range"`
	Limit uint32 `json:"limit" form:"limit"`
}

type SeoOperationKPI struct {
	Label  string `json:"label"`
	Value  string `json:"value"`
	Change string `json:"change"`
	Tone   string `json:"tone"`
}

type SeoOperationTrendPoint struct {
	Label       string `json:"label"`
	Clicks      uint64 `json:"clicks"`
	Impressions uint64 `json:"impressions"`
}

type SeoOperationOpportunity struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Issue       string `json:"issue"`
	Impact      string `json:"impact"`
	PageCount   uint32 `json:"pageCount"`
	Owner       string `json:"owner"`
	Tone        string `json:"tone"`
}

type SeoOperationTopQuery struct {
	Query       string  `json:"query"`
	Clicks      uint64  `json:"clicks"`
	Impressions uint64  `json:"impressions"`
	CTR         float64 `json:"ctr"`
	AvgPosition float64 `json:"avgPosition"`
}

type SeoOperationEvent struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	Tone        string `json:"tone"`
}

type SeoRequirementItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Tone        string `json:"tone"`
	Owner       string `json:"owner"`
	Coverage    uint32 `json:"coverage"`
	Hint        string `json:"hint"`
}

type SeoRequirementGroup struct {
	ID          string               `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Items       []SeoRequirementItem `json:"items"`
}

type SeoRenderQueueItem struct {
	ID          string `json:"id"`
	SeoDomainID uint64 `json:"seoDomainId"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Worker      string `json:"worker"`
	Duration    string `json:"duration"`
	Priority    string `json:"priority"`
}

type SeoIndexingRow struct {
	URL       string `json:"url"`
	Canonical string `json:"canonical"`
	Robots    string `json:"robots"`
	Sitemap   string `json:"sitemap"`
	State     string `json:"state"`
	LastPing  string `json:"lastPing"`
}

type SeoSourceHealthItem struct {
	Source        string `json:"source"`
	Label         string `json:"label"`
	Service       string `json:"service"`
	Status        string `json:"status"`
	TotalPages    uint32 `json:"totalPages"`
	NeedAttention uint32 `json:"needAttention"`
}

type SeoOperationsDashboardResponse struct {
	KPIs          []SeoOperationKPI         `json:"kpis"`
	Trend         []SeoOperationTrendPoint  `json:"trend"`
	Opportunities []SeoOperationOpportunity `json:"opportunities"`
	TopQueries    []SeoOperationTopQuery    `json:"topQueries"`
	SystemEvents  []SeoOperationEvent       `json:"systemEvents"`
	Requirements  []SeoRequirementGroup     `json:"requirements"`
	RenderQueue   []SeoRenderQueueItem      `json:"renderQueue"`
	IndexingRows  []SeoIndexingRow          `json:"indexingRows"`
	SourceHealth  []SeoSourceHealthItem     `json:"sourceHealth"`
}
