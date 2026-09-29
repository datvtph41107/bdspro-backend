package usageprojection

import (
	commonmetering "common/metering"
	"context"
	"io"
	"log/slog"
	"testing"
	"time"
	"tqd/internal/access"
	"tqd/internal/usecase/quota/memory"
	"tqd/internal/usecase/usage/memory"
)

type pagingScopeStore struct {
	offsets []int
	pages   map[int][]Scope
}

func (s *pagingScopeStore) ListUsageScopes(_ context.Context, offset int, _ int) ([]Scope, error) {
	s.offsets = append(s.offsets, offset)
	return s.pages[offset], nil
}

func TestRunnerRotatesThroughUsageScopePages(t *testing.T) {
	periodStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	scope := func(id string) Scope {
		return Scope{
			Subject:     access.Subject{Type: access.SubjectProfile, ID: id},
			MeterCode:   commonmetering.Code("workspace.report_generation.accepted"),
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
		}
	}
	store := &pagingScopeStore{pages: map[int][]Scope{
		0: {scope("1"), scope("2")},
		2: {scope("3")},
	}}
	runner := NewRunner(
		NewService(memoryquota.NewStore(), memoryusage.NewStore()),
		store,
		time.Minute,
		2,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	if err := runner.ReconcileAll(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(store.offsets) != 2 || store.offsets[0] != 0 || store.offsets[1] != 2 {
		t.Fatalf("page offsets = %v, want [0 2]", store.offsets)
	}
	if runner.offset != 0 {
		t.Fatalf("runner offset = %d, want reset after final short page", runner.offset)
	}
}
