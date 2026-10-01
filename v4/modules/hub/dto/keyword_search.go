package dto

// KeywordSearchItem is one top-keyword result row.
type KeywordSearchItem struct {
	Keyword string `json:"keyword"`
	Source  string `json:"source"`
	Count   uint64 `json:"count"`
}

// GetKeywordSearchStatsQuery binds query params for GET /v2/hub/keyword-search-stats.
type GetKeywordSearchStatsQuery struct {
	LastDates int32  `form:"lastDates"`
	Limit     int32  `form:"limit"`
	Source    string `form:"source"`
}

// GetKeywordSearchStatsResponse is the dashboard response.
type GetKeywordSearchStatsResponse struct {
	Data          []KeywordSearchItem `json:"data"`
	TotalSearches uint64              `json:"totalSearches"`
	LastDates     int32               `json:"lastDates"`
	Limit         int32               `json:"limit"`
}
