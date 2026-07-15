package mapper

import (
	"encoding/json"
	"testing"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
)

func TestSeoQueryEvidenceMapperToPb(t *testing.T) {
	t.Parallel()

	zone := time.FixedZone("GMT+7", 7*60*60)
	observedFrom := time.Date(2026, time.July, 1, 0, 0, 0, 0, zone)
	observedTo := time.Date(2026, time.July, 14, 23, 59, 59, 0, zone)
	validUntil := time.Date(2026, time.October, 15, 0, 0, 0, 0, zone)
	impressions := int64(100)
	clicks := int64(12)
	ctr := 0.12
	position := 4.5
	volume := int64(1200)

	item := &seo_strategy.SeoQueryEvidence{
		ID:              1,
		ObjectiveID:     10,
		UserNeedID:      20,
		RawQuery:        "Tra cứu quy hoạch số tờ số thửa",
		NormalizedQuery: "tra cứu quy hoạch số tờ số thửa",
		QueryHash:       "0123456789012345678901234567890123456789012345678901234567890123",
		SourceSystem:    "google_search_console",
		SourceKind:      seo_strategy.QueryEvidenceSourceFirstParty,
		ObservationType: seo_strategy.QueryObservationObserved,
		SourceRecordKey: "gsc|2026-07|query-1",
		SourceReference: "sc-domain:qhpro.vn",
		Locale:          "vi-VN",
		CountryCode:     "VN",
		DeviceCategory:  seo_strategy.QueryDeviceMobile,
		ObservedFrom:    &observedFrom,
		ObservedTo:      &observedTo,
		ValidUntil:      &validUntil,
		Impressions:     &impressions,
		Clicks:          &clicks,
		CTR:             &ctr,
		AvgPosition:     &position,
		SearchVolume:    &volume,
		Confidence:      seo_strategy.ObjectiveConfidenceSupported,
		Status:          seo_strategy.QueryEvidenceStatusActive,
		Metadata:        json.RawMessage(`{"property":"qhpro.vn"}`),
		CreatedAt:       observedFrom,
		UpdatedAt:       observedTo,
	}

	result := NewSeoQueryEvidenceMapper().ToPb(item)
	if result == nil {
		t.Fatal("response không được nil")
	}
	if result.GetId() != 1 || result.GetObjectiveId() != 10 || result.GetUserNeedId() != 20 {
		t.Fatalf("identity=%+v", result)
	}
	if result.GetStatus() != "active" {
		t.Fatalf("status=%q", result.GetStatus())
	}
	if result.Impressions == nil || result.GetImpressions() != 100 {
		t.Fatalf("impressions=%v", result.Impressions)
	}
	if result.ObservedFrom == nil || result.GetObservedFrom() != "2026-06-30T17:00:00Z" {
		t.Fatalf("observedFrom=%q", result.GetObservedFrom())
	}
	if result.GetMetadataJson() != `{"property":"qhpro.vn"}` {
		t.Fatalf("metadata=%q", result.GetMetadataJson())
	}
}
