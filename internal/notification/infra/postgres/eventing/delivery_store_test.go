package eventing

import (
	"context"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDeliveryClaimQueryIncludesOnlyExpiredRunningClaims(t *testing.T) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=test dbname=test sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		DryRun:                 true,
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("open dry-run database: %v", err)
	}

	now := time.Date(2026, time.September, 12, 3, 0, 0, 0, time.UTC)
	store := NewDeliveryStore(db)
	stmt := store.claimCandidateQuery(context.Background(), db, now).Find(&deliveryRow{}).Statement
	query := strings.Join(strings.Fields(stmt.SQL.String()), " ")

	if !strings.Contains(query, "status IN") {
		t.Fatalf("claim query = %q, want pending/retry/unknown status predicate", query)
	}
	if !strings.Contains(query, "status =") || !strings.Contains(query, "lease_until IS NOT NULL") {
		t.Fatalf("claim query = %q, want running claims fenced by an explicit lease", query)
	}
	if strings.Count(query, "lease_until <=") < 2 {
		t.Fatalf("claim query = %q, want both normal and running claims fenced by lease expiry", query)
	}

	wantStatus := map[string]bool{"pending": false, "retry": false, "unknown": false, "running": false}
	for _, value := range stmt.Vars {
		status, ok := value.(string)
		if ok {
			if _, tracked := wantStatus[status]; tracked {
				wantStatus[status] = true
			}
		}
	}
	for status, found := range wantStatus {
		if !found {
			t.Fatalf("claim query vars = %#v, missing status %q", stmt.Vars, status)
		}
	}
}
