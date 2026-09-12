package dashboard

import (
	"context"
	"errors"
	"testing"
)

func TestDashboardStatsJobRunStopsWhenMetricReadFails(t *testing.T) {
	readErr := errors.New("count unavailable")
	job := NewDashboardStatsJobWithError(nil, "user_stats", func(context.Context) (int64, error) {
		return 0, readErr
	})

	err := job.Run(context.Background())
	if !errors.Is(err, readErr) {
		t.Fatalf("Run() error = %v, want %v", err, readErr)
	}
}
