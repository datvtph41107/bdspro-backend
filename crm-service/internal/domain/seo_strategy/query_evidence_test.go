package seo_strategy

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func queryEvidenceInt64(value int64) *int64       { return &value }
func queryEvidenceFloat64(value float64) *float64 { return &value }

func validQueryEvidenceFixture(t *testing.T) *SeoQueryEvidence {
	t.Helper()

	now := time.Date(2026, time.July, 14, 10, 0, 0, 0, time.UTC)
	evidence, err := CreateSeoQueryEvidence(
		1,
		2,
		"Tra cứu quy hoạch số tờ số thửa",
		"google_search_console",
		QueryEvidenceSourceFirstParty,
		QueryObservationObserved,
		"gsc|2026-06|query-001",
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	return evidence
}

func TestNormalizeSearchQuery(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"  Tra   CỨU Quy Hoạch  ":      "tra cứu quy hoạch",
		"TRA CỨU\nTHỬA ĐẤT":            "tra cứu thửa đất",
		"xem\tquy hoạch":               "xem quy hoạch",
		"quy hoạch số tờ số thửa":      "quy hoạch số tờ số thửa",
		"  bản đồ   quy hoạch\nTP.HCM": "bản đồ quy hoạch tp.hcm",
	}

	for input, want := range tests {
		input, want := input, want
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if got := NormalizeSearchQuery(input); got != want {
				t.Fatalf("NormalizeSearchQuery(%q)=%q, want %q", input, got, want)
			}
		})
	}
}

func TestBuildQueryEvidenceHashIsStable(t *testing.T) {
	t.Parallel()

	first := BuildQueryEvidenceHash("tra cứu quy hoạch", "vi-VN")
	second := BuildQueryEvidenceHash("  TRA CỨU   QUY HOẠCH ", "VI-vn")
	if first != second {
		t.Fatalf("hash phải ổn định: %s != %s", first, second)
	}
	if len(first) != 64 {
		t.Fatalf("hash length=%d, want 64", len(first))
	}
	if first == BuildQueryEvidenceHash("tra cứu quy hoạch", "en-US") {
		t.Fatal("locale khác nhau phải tạo hash khác nhau")
	}
}

func TestCreateSeoQueryEvidenceDefaults(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.July, 14, 10, 0, 0, 0, time.UTC)
	evidence, err := CreateSeoQueryEvidence(
		10,
		20,
		"  Tra   CỨU Quy Hoạch  ",
		"internal_search",
		QueryEvidenceSourceFirstParty,
		QueryObservationObserved,
		"internal|20260714|abc",
		now,
	)
	if err != nil {
		t.Fatal(err)
	}

	if evidence.RawQuery != "Tra   CỨU Quy Hoạch" {
		t.Fatalf("raw query=%q", evidence.RawQuery)
	}
	if evidence.NormalizedQuery != "tra cứu quy hoạch" {
		t.Fatalf("normalized query=%q", evidence.NormalizedQuery)
	}
	if evidence.Locale != "vi-VN" || evidence.CountryCode != "VN" {
		t.Fatalf("locale/country=%q/%q", evidence.Locale, evidence.CountryCode)
	}
	if evidence.DeviceCategory != QueryDeviceUnknown {
		t.Fatalf("device=%q", evidence.DeviceCategory)
	}
	if evidence.Confidence != ObjectiveConfidenceExploratory {
		t.Fatalf("confidence=%q", evidence.Confidence)
	}
	if evidence.Status != QueryEvidenceStatusActive {
		t.Fatalf("status=%q", evidence.Status)
	}
	if string(evidence.Metadata) != "{}" {
		t.Fatalf("metadata=%s", evidence.Metadata)
	}
	if evidence.CreatedAt != now || evidence.UpdatedAt != now {
		t.Fatal("createdAt/updatedAt không đúng")
	}
}

func TestValidateQueryEvidenceValid(t *testing.T) {
	t.Parallel()

	evidence := validQueryEvidenceFixture(t)
	impressions := int64(100)
	clicks := int64(15)
	ctr := 0.15
	position := 4.8
	volume := int64(1200)

	evidence.Impressions = &impressions
	evidence.Clicks = &clicks
	evidence.CTR = &ctr
	evidence.AvgPosition = &position
	evidence.SearchVolume = &volume
	evidence.Metadata = json.RawMessage(`{"property":"sc-domain:qhpro.vn"}`)

	if err := ValidateQueryEvidence(evidence); err != nil {
		t.Fatal(err)
	}
}

func TestValidateQueryEvidenceRejectsInvalidIdentity(t *testing.T) {
	t.Parallel()

	tests := map[string]func(*SeoQueryEvidence){
		"objective id": func(v *SeoQueryEvidence) { v.ObjectiveID = 0 },
		"user need id": func(v *SeoQueryEvidence) { v.UserNeedID = 0 },
		"raw query":    func(v *SeoQueryEvidence) { v.RawQuery = "   " },
		"raw query too long": func(v *SeoQueryEvidence) {
			v.RawQuery = strings.Repeat("x", 1001)
		},
		"source system": func(v *SeoQueryEvidence) { v.SourceSystem = "Google Search Console" },
		"source kind": func(v *SeoQueryEvidence) {
			v.SourceKind = QueryEvidenceSourceKind("invalid")
		},
		"observation type": func(v *SeoQueryEvidence) {
			v.ObservationType = QueryObservationType("invalid")
		},
		"source record key": func(v *SeoQueryEvidence) { v.SourceRecordKey = " " },
		"locale":            func(v *SeoQueryEvidence) { v.Locale = " " },
		"country code":      func(v *SeoQueryEvidence) { v.CountryCode = "VNM" },
		"device": func(v *SeoQueryEvidence) {
			v.DeviceCategory = QueryDeviceCategory("watch")
		},
		"confidence": func(v *SeoQueryEvidence) {
			v.Confidence = ObjectiveConfidence("invalid")
		},
		"status": func(v *SeoQueryEvidence) {
			v.Status = QueryEvidenceStatus("invalid")
		},
		"metadata": func(v *SeoQueryEvidence) { v.Metadata = json.RawMessage(`{`) },
	}

	for name, mutate := range tests {
		name, mutate := name, mutate
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			evidence := validQueryEvidenceFixture(t)
			mutate(evidence)
			if err := ValidateQueryEvidence(evidence); err == nil {
				t.Fatalf("%s sai phải bị từ chối", name)
			}
		})
	}
}

func TestValidateQueryEvidenceRejectsInvalidPeriod(t *testing.T) {
	t.Parallel()

	evidence := validQueryEvidenceFixture(t)
	from := time.Date(2026, time.July, 14, 10, 0, 0, 0, time.UTC)
	to := from.Add(-time.Hour)
	evidence.ObservedFrom = &from
	evidence.ObservedTo = &to

	if err := ValidateQueryEvidence(evidence); err == nil {
		t.Fatal("observedTo trước observedFrom phải bị từ chối")
	}
}

func TestValidateQueryEvidenceMetrics(t *testing.T) {
	t.Parallel()

	tests := map[string]func(*SeoQueryEvidence){
		"negative impressions": func(v *SeoQueryEvidence) { v.Impressions = queryEvidenceInt64(-1) },
		"negative clicks":      func(v *SeoQueryEvidence) { v.Clicks = queryEvidenceInt64(-1) },
		"clicks exceed impressions": func(v *SeoQueryEvidence) {
			v.Impressions = queryEvidenceInt64(5)
			v.Clicks = queryEvidenceInt64(6)
		},
		"negative ctr":    func(v *SeoQueryEvidence) { v.CTR = queryEvidenceFloat64(-0.1) },
		"ctr over one":    func(v *SeoQueryEvidence) { v.CTR = queryEvidenceFloat64(1.1) },
		"zero position":   func(v *SeoQueryEvidence) { v.AvgPosition = queryEvidenceFloat64(0) },
		"negative volume": func(v *SeoQueryEvidence) { v.SearchVolume = queryEvidenceInt64(-1) },
	}

	for name, mutate := range tests {
		name, mutate := name, mutate
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			evidence := validQueryEvidenceFixture(t)
			mutate(evidence)
			if err := ValidateQueryEvidence(evidence); err == nil {
				t.Fatalf("%s phải bị từ chối", name)
			}
		})
	}
}

func TestValidateQueryEvidenceRejectsHashMismatch(t *testing.T) {
	t.Parallel()

	evidence := validQueryEvidenceFixture(t)
	evidence.QueryHash = strings.Repeat("a", 64)
	if err := ValidateQueryEvidence(evidence); err == nil {
		t.Fatal("hash không khớp phải bị từ chối")
	}
}

func TestParseQueryEvidenceEnums(t *testing.T) {
	t.Parallel()

	if got, err := ParseQueryEvidenceSourceKind("first_party"); err != nil || got != QueryEvidenceSourceFirstParty {
		t.Fatalf("source kind got=%q err=%v", got, err)
	}
	if _, err := ParseQueryEvidenceSourceKind("firstparty"); err == nil {
		t.Fatal("source kind sai phải bị từ chối")
	}
	if got, err := ParseQueryObservationType("observed"); err != nil || got != QueryObservationObserved {
		t.Fatalf("observation got=%q err=%v", got, err)
	}
	if got, err := ParseQueryDeviceCategory(""); err != nil || got != QueryDeviceUnknown {
		t.Fatalf("blank device got=%q err=%v", got, err)
	}
	if got, err := ParseQueryEvidenceStatus(""); err != nil || got != QueryEvidenceStatusActive {
		t.Fatalf("blank status got=%q err=%v", got, err)
	}
}

func TestIsQueryEvidenceCurrent(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.July, 14, 10, 0, 0, 0, time.UTC)
	evidence := validQueryEvidenceFixture(t)
	if !IsQueryEvidenceCurrent(*evidence, now) {
		t.Fatal("active evidence phải current")
	}

	expiredAt := now.Add(-time.Second)
	evidence.ValidUntil = &expiredAt
	if IsQueryEvidenceCurrent(*evidence, now) {
		t.Fatal("expired evidence không được current")
	}

	evidence.ValidUntil = nil
	evidence.Status = QueryEvidenceStatusRejected
	if IsQueryEvidenceCurrent(*evidence, now) {
		t.Fatal("rejected evidence không được current")
	}
}

func TestIsQueryEvidenceObserved(t *testing.T) {
	t.Parallel()

	evidence := validQueryEvidenceFixture(t)
	for _, observationType := range []QueryObservationType{
		QueryObservationObserved,
		QueryObservationReported,
	} {
		evidence.ObservationType = observationType
		if !IsQueryEvidenceObserved(*evidence) {
			t.Fatalf("%s phải được coi là observation", observationType)
		}
	}
	for _, observationType := range []QueryObservationType{
		QueryObservationSuggested,
		QueryObservationAssumed,
	} {
		evidence.ObservationType = observationType
		if IsQueryEvidenceObserved(*evidence) {
			t.Fatalf("%s chưa phải observation thực", observationType)
		}
	}
}
