package services

import (
	"context"
	"testing"
)

func TestParseSyncTimestamp(t *testing.T) {
	cases := []struct {
		in   any
		want int64
	}{
		{nil, 0},
		{"1790325256127", 1790325256127},
		{[]byte("1790325256127"), 1790325256127},
		{int64(1790325256127), 1790325256127},
		{int(42), 42},
		{float64(99), 99},
		{"", 0},
	}
	for _, tc := range cases {
		if got := parseSyncTimestamp(tc.in); got != tc.want {
			t.Fatalf("parseSyncTimestamp(%v)=%d want %d", tc.in, got, tc.want)
		}
	}
}

func TestLocalSyncService_NilSyncReturnsEmpty(t *testing.T) {
	svc := NewLocalSyncService(nil)
	got, err := svc.GetAppTimestamps(context.Background())
	if err != nil {
		t.Fatalf("GetAppTimestamps: %v", err)
	}
	if got == nil || got.Timestamps == nil {
		t.Fatalf("expected empty timestamps map, got %#v", got)
	}
	if len(got.Timestamps) != 0 {
		t.Fatalf("expected 0 timestamps, got %d", len(got.Timestamps))
	}
}
