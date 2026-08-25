package dto

import "strings"

type TrackKeywordSearchRequest struct {
	Keyword string `json:"keyword"`
	Source  string `json:"source"`
	Count   int64  `json:"count"`
}

type GetKeywordSearchStatsRequest struct {
	LastDates uint32 `json:"lastDates" form:"lastDates"`
	Limit     uint32 `json:"limit" form:"limit"`
	Source    string `json:"source" form:"source"`
}

func NormalizeKeywordSearchLimit(limit uint32) uint32 {
	if limit == 0 {
		return 50
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func (r *GetKeywordSearchStatsRequest) Normalize() {
	if r == nil {
		return
	}
	r.LastDates = NormalizeLastDates(r.LastDates)
	r.Limit = NormalizeKeywordSearchLimit(r.Limit)
	r.Source = strings.ToLower(strings.TrimSpace(r.Source))
}

type KeywordSearchItem struct {
	Keyword string `json:"keyword"`
	Source  string `json:"source"`
	Count   uint64 `json:"count"`
}

type GetKeywordSearchStatsResponse struct {
	Data          []KeywordSearchItem `json:"data"`
	TotalSearches uint64              `json:"totalSearches"`
	LastDates     uint32              `json:"lastDates"`
	Limit         uint32              `json:"limit"`
}
