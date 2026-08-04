package postgres

import (
	"encoding/json"
	"testing"
	"time"

	"user/internal/models"
)

func TestToBackfillMappingCopiesReviewEvidence(t *testing.T) {
	reviewer := uint64(7)
	reviewedAt := time.Date(2026, 8, 4, 4, 0, 0, 0, time.UTC)
	decision := json.RawMessage(`{"periodHandling":"paid_at_to_finished_date"}`)
	mapping := toBackfillMapping(models.LegacyMembershipSubscriptionMapping{
		LegacyMembershipID: 9,
		SourceChecksum:     "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Status:             models.CatalogMigrationReviewed,
		Decision:           decision,
		ReviewedBy:         &reviewer,
		ReviewedAt:         &reviewedAt,
		BlockedReason:      "",
		LastRunID:          "run_1",
	})
	if mapping.LegacyMembershipID != 9 || mapping.ReviewedBy == nil ||
		*mapping.ReviewedBy != reviewer || mapping.ReviewedAt == nil {
		t.Fatalf("mapping = %+v", mapping)
	}
	mapping.Decision[0] = '['
	*mapping.ReviewedBy = 99
	if decision[0] == '[' || reviewer != 7 {
		t.Fatal("mapping retained mutable model fields")
	}
}

func TestSameLegacyBackfillEventRequiresMatchingEvidence(t *testing.T) {
	planVersionID := uint64(11)
	expected := models.CatalogSubscriptionEvent{
		SubscriptionID:  19,
		EventType:       "migration.legacy_membership_backfilled",
		ToPlanVersionID: &planVersionID,
		Metadata: json.RawMessage(`{
			"legacyMembershipId":9,
			"legacyPlanId":7,
			"sourceChecksum":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			"planTermsChecksum":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
		}`),
	}
	stored := expected
	stored.Metadata = append(json.RawMessage(nil), expected.Metadata...)
	if !sameLegacyBackfillEvent(stored, expected) {
		t.Fatal("matching event evidence was rejected")
	}
	stored.Metadata = json.RawMessage(`{
		"legacyMembershipId":10,
		"legacyPlanId":7,
		"sourceChecksum":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"planTermsChecksum":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	}`)
	if sameLegacyBackfillEvent(stored, expected) {
		t.Fatal("event for another legacy membership was accepted")
	}
}
