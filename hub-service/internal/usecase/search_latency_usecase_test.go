package usecase

import "testing"

func TestPercentileMs(t *testing.T) {
	sorted := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if got := percentileMs(sorted, 50); got != 5 {
		t.Fatalf("p50 got %d want 5", got)
	}
	if got := percentileMs(sorted, 95); got != 10 {
		t.Fatalf("p95 got %d want 10", got)
	}
	if got := percentileMs(sorted, 99); got != 10 {
		t.Fatalf("p99 got %d want 10", got)
	}
}

func TestSearchLatencyAggPercentiles(t *testing.T) {
	agg := &searchLatencyAgg{}
	for _, v := range []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100} {
		agg.addSample(v)
	}
	p50, p95, p99 := agg.percentiles()
	if p50 != 50 {
		t.Fatalf("p50 got %d want 50", p50)
	}
	if p95 != 100 {
		t.Fatalf("p95 got %d want 100", p95)
	}
	if p99 != 100 {
		t.Fatalf("p99 got %d want 100", p99)
	}
}
