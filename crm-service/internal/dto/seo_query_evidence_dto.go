package dto

import (
	"encoding/json"
	"time"
)

type UpsertSeoQueryEvidenceRequest struct {
	UserNeedID uint64

	RawQuery string

	SourceSystem    string
	SourceKind      string
	ObservationType string
	SourceRecordKey string
	SourceReference string

	Locale         string
	CountryCode    string
	DeviceCategory string

	ObservedFrom *time.Time
	ObservedTo   *time.Time
	ValidUntil   *time.Time

	Impressions  *int64
	Clicks       *int64
	CTR          *float64
	AvgPosition  *float64
	SearchVolume *int64

	Confidence string
	Metadata   json.RawMessage
}

type SeoQueryEvidenceListRequest struct {
	ObjectiveID uint64
	UserNeedID  uint64

	Keyword      string
	SourceSystem string

	SourceKind      string
	ObservationType string
	DeviceCategory  string
	Confidence      string
	Status          string

	CurrentOnly bool

	Page uint32
	Size uint32
}

type SeoQueryEvidenceStatusCommandRequest struct {
	ID             uint64
	ExpectedStatus string
}
