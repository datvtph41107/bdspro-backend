package dto

import "encoding/json"

type SeoMeasurementImportRequest struct {
	Source      string                   `json:"source"`
	RequestedBy string                   `json:"requestedBy"`
	RangeStart  string                   `json:"rangeStart"`
	RangeEnd    string                   `json:"rangeEnd"`
	Checksum    string                   `json:"checksum"`
	Rows        []map[string]interface{} `json:"rows"`
	Metadata    map[string]interface{}   `json:"metadata"`
}

type SeoMeasurementImportResponse struct {
	RunID         uint64   `json:"runId"`
	Source        string   `json:"source"`
	Status        string   `json:"status"`
	RowCount      int64    `json:"rowCount"`
	AcceptedCount int64    `json:"acceptedCount"`
	RejectedCount int64    `json:"rejectedCount"`
	Errors        []string `json:"errors"`
	Idempotent    bool     `json:"idempotent"`
}

type SeoMeasurementSourceHealth struct {
	Source              string `json:"source"`
	Status              string `json:"status"`
	LastRunID           uint64 `json:"lastRunId"`
	LastAttemptAt       string `json:"lastAttemptAt"`
	LastSuccessAt       string `json:"lastSuccessAt"`
	LastRowCount        int64  `json:"lastRowCount"`
	StaleAfterHours     int    `json:"staleAfterHours"`
	AgeHours            int64  `json:"ageHours"`
	ConsecutiveFailures int    `json:"consecutiveFailures"`
	LastError           string `json:"lastError"`
	Fresh               bool   `json:"fresh"`
}

type SeoMeasurementImportRunRow struct {
	ID            uint64 `json:"id"`
	Source        string `json:"source"`
	Status        string `json:"status"`
	RequestedBy   string `json:"requestedBy"`
	RowCount      int64  `json:"rowCount"`
	AcceptedCount int64  `json:"acceptedCount"`
	RejectedCount int64  `json:"rejectedCount"`
	Checksum      string `json:"checksum"`
	ErrorMessage  string `json:"errorMessage"`
	StartedAt     string `json:"startedAt"`
	FinishedAt    string `json:"finishedAt"`
}

type SeoMeasurementImportRunListRequest struct {
	Page   uint32 `json:"page"`
	Size   uint32 `json:"size"`
	Source string `json:"source"`
	Status string `json:"status"`
}

type SeoMeasurementScoreRequest struct {
	Range string `json:"range"`
	Limit uint32 `json:"limit"`
}

type SeoMeasurementScoreResponse struct {
	Scanned   uint32 `json:"scanned"`
	Persisted uint32 `json:"persisted"`
	A         uint32 `json:"a"`
	B         uint32 `json:"b"`
	C         uint32 `json:"c"`
	D         uint32 `json:"d"`
	F         uint32 `json:"f"`
}

type SeoPageScoreListRequest struct {
	Page         uint32 `json:"page"`
	Size         uint32 `json:"size"`
	CanonicalURL string `json:"canonicalUrl"`
	Grade        string `json:"grade"`
}

type SeoPageScoreRow struct {
	SeoDomainID     uint64          `json:"seoDomainId"`
	CanonicalURL    string          `json:"canonicalUrl"`
	RangeKey        string          `json:"rangeKey"`
	MeasuredAt      string          `json:"measuredAt"`
	IndexScore      float64         `json:"indexScore"`
	VisibilityScore float64         `json:"visibilityScore"`
	EngagementScore float64         `json:"engagementScore"`
	ConversionScore float64         `json:"conversionScore"`
	TechnicalScore  float64         `json:"technicalScore"`
	IssuePenalty    float64         `json:"issuePenalty"`
	OverallScore    float64         `json:"overallScore"`
	Grade           string          `json:"grade"`
	Blockers        json.RawMessage `json:"blockers"`
	Warnings        json.RawMessage `json:"warnings"`
	SourceFreshness json.RawMessage `json:"sourceFreshness"`
}
