package usecase

import (
	"context"
	"testing"
	"time"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
)

type runtimeEventRepoStub struct {
	created bool
	event   *seo_domain.SeoRuntimeEvent
}

func (r *runtimeEventRepoStub) CreateIdempotent(_ context.Context, event *seo_domain.SeoRuntimeEvent) (bool, error) {
	r.event = event
	return r.created, nil
}

func TestSeoRuntimeEventUsecasePersistsSanitizedEvent(t *testing.T) {
	now := time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)
	repository := &runtimeEventRepoStub{created: true}
	uc := NewSeoRuntimeEventUsecase(repository)
	uc.now = func() time.Time { return now }
	indexable := true

	response, err := uc.Ingest(context.Background(), &dto.SeoRuntimeEventIngestRequest{
		SchemaVersion: 1,
		Source:        "nextjs",
		Event:         "seo_page_view",
		SessionID:     "session-1",
		CanonicalURL:  "https://qhpro.vn/dia-ban/phuong-ben-nghe-26734",
		CanonicalPath: "/dia-ban/phuong-ben-nghe-26734",
		ModuleID:      "administrative-unit",
		Environment:   "production",
		Indexable:     &indexable,
		Timestamp:     now.Add(-time.Second).Format(time.RFC3339Nano),
		Metrics:       map[string]float64{"lcp": 1234.5},
		Metadata:      map[string]interface{}{"surface": "seo-detail"},
		UserAgent:     "test-agent",
	})
	if err != nil {
		t.Fatalf("Ingest returned error: %v", err)
	}
	if response == nil || !response.Accepted || !response.Persisted || response.Idempotent || len(response.EventKey) != 64 {
		t.Fatalf("unexpected response: %#v", response)
	}
	if repository.event == nil || repository.event.ModuleID != "administrative-unit" || repository.event.CanonicalPath != "/dia-ban/phuong-ben-nghe-26734" {
		t.Fatalf("unexpected persisted event: %#v", repository.event)
	}
}

func TestSeoRuntimeEventUsecaseRejectsInvalidContract(t *testing.T) {
	now := time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)
	uc := NewSeoRuntimeEventUsecase(&runtimeEventRepoStub{created: true})
	uc.now = func() time.Time { return now }

	cases := []*dto.SeoRuntimeEventIngestRequest{
		{SchemaVersion: 2, Source: "nextjs", Event: "seo_page_view"},
		{SchemaVersion: 1, Source: "browser", Event: "seo_page_view"},
		{SchemaVersion: 1, Source: "nextjs", Event: "unknown_event"},
		{SchemaVersion: 1, Source: "nextjs", Event: "seo_page_view", SessionID: "s", ModuleID: "administrative-unit", CanonicalPath: "../x", CanonicalURL: "https://qhpro.vn/x", Timestamp: now.Format(time.RFC3339Nano)},
		{SchemaVersion: 1, Source: "nextjs", Event: "seo_page_view", SessionID: "s", ModuleID: "administrative-unit", CanonicalPath: "/dia-ban/a", CanonicalURL: "https://qhpro.vn/dia-ban/b", Timestamp: now.Format(time.RFC3339Nano)},
	}

	for index, req := range cases {
		if _, err := uc.Ingest(context.Background(), req); err == nil {
			t.Fatalf("case %d expected validation error", index)
		}
	}
}

func TestSeoRuntimeEventUsecaseReportsIdempotentDuplicate(t *testing.T) {
	now := time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)
	uc := NewSeoRuntimeEventUsecase(&runtimeEventRepoStub{created: false})
	uc.now = func() time.Time { return now }

	response, err := uc.Ingest(context.Background(), &dto.SeoRuntimeEventIngestRequest{
		SchemaVersion: 1,
		Source:        "react-native",
		Event:         "rn_activation",
		SessionID:     "session-1",
		CanonicalURL:  "https://qhpro.vn/dia-ban/a",
		CanonicalPath: "/dia-ban/a",
		ModuleID:      "administrative-unit",
		Timestamp:     now.Format(time.RFC3339Nano),
	})
	if err != nil {
		t.Fatalf("Ingest returned error: %v", err)
	}
	if response == nil || !response.Accepted || response.Persisted || !response.Idempotent {
		t.Fatalf("unexpected duplicate response: %#v", response)
	}
}
