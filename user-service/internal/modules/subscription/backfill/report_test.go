package backfill

import (
	"context"
	"testing"
	"time"

	"user/internal/models"
	"user/internal/modules/subscription/shadow"
)

type fakeMembershipSource struct {
	memberships []models.Membership
}

func (s fakeMembershipSource) ListLegacyMemberships(
	_ context.Context,
	afterID uint64,
	limit int,
) ([]models.Membership, error) {
	rows := make([]models.Membership, 0, limit)
	for _, membership := range s.memberships {
		if membership.ID > afterID && len(rows) < limit {
			rows = append(rows, membership)
		}
	}
	return rows, nil
}

type fakeComparator struct {
	entries map[uint64]shadow.Entry
}

func (c fakeComparator) CompareMembership(
	_ context.Context,
	membership models.Membership,
	_ time.Time,
) shadow.Entry {
	return c.entries[membership.ID]
}

type fakeEventEvidence struct {
	events map[uint64]uint64
}

func (e fakeEventEvidence) HasLegacyBackfillEvent(
	_ context.Context,
	subscriptionID uint64,
	legacyMembershipID uint64,
	sourceChecksum string,
) (bool, error) {
	return e.events[subscriptionID] == legacyMembershipID && sourceChecksum != "", nil
}

func membershipWithID(id uint64) models.Membership {
	paidAt := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	finishedAt := paidAt.AddDate(0, 1, 0)
	membership := models.Membership{
		ProfileID:    id + 100,
		PlanID:       7,
		FinishedDate: &finishedAt,
		OrderID:      id + 200,
		PaidAt:       &paidAt,
	}
	membership.ID = id
	return membership
}

func TestAuditorBuildsCoverageAndRollbackEvidence(t *testing.T) {
	auditor := NewAuditor(
		fakeMembershipSource{memberships: []models.Membership{
			membershipWithID(1), membershipWithID(2), membershipWithID(3),
		}},
		fakeComparator{entries: map[uint64]shadow.Entry{
			1: {
				LegacyMembershipID: 1,
				SubscriptionID:     11,
				MappingStatus:      models.CatalogMigrationShadowRead,
				Status:             shadow.StatusMatched,
			},
			2: {
				LegacyMembershipID: 2,
				SubscriptionID:     12,
				MappingStatus:      models.CatalogMigrationShadowRead,
				Status:             shadow.StatusPeriodMismatch,
			},
			3: {
				LegacyMembershipID: 3,
				Status:             shadow.StatusLegacyIncomplete,
			},
		}},
		fakeEventEvidence{events: map[uint64]uint64{11: 1, 12: 2}},
	)
	report, err := auditor.Build(
		context.Background(),
		time.Date(2026, 8, 4, 4, 0, 0, 0, time.UTC),
		2,
		true,
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.TotalLegacy != 3 || report.EligibleLegacy != 2 ||
		report.Observed != 2 || report.Matched != 1 || report.Mismatched != 1 {
		t.Fatalf("report = %+v", report)
	}
	if report.RollbackEvidence != 2 || report.Checksum == "" {
		t.Fatalf("rollback=%d checksum=%q", report.RollbackEvidence, report.Checksum)
	}
}

func TestAuditorRejectsEventForDifferentLegacyMembership(t *testing.T) {
	auditor := NewAuditor(
		fakeMembershipSource{memberships: []models.Membership{membershipWithID(1)}},
		fakeComparator{entries: map[uint64]shadow.Entry{
			1: {
				LegacyMembershipID: 1,
				SubscriptionID:     11,
				MappingStatus:      models.CatalogMigrationShadowRead,
				Status:             shadow.StatusMatched,
			},
		}},
		fakeEventEvidence{events: map[uint64]uint64{11: 99}},
	)
	report, err := auditor.Build(
		context.Background(),
		time.Date(2026, 8, 4, 4, 0, 0, 0, time.UTC),
		10,
		false,
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.RollbackEvidence != 0 || report.MissingRollbackProof != 1 {
		t.Fatalf("report = %+v", report)
	}
}

func TestStrictReadinessRequiresCompleteMatchedEvidence(t *testing.T) {
	report := CoverageReport{
		AsOf:             time.Date(2026, 8, 4, 4, 0, 0, 0, time.UTC),
		Checksum:         "checksum",
		EligibleLegacy:   10,
		Observed:         10,
		Matched:          10,
		RollbackEvidence: 10,
		Statuses:         map[string]int{"matched": 10},
	}
	result, err := EvaluateReadiness(report, StrictReadinessPolicy())
	if err != nil {
		t.Fatal(err)
	}
	if !result.Ready {
		t.Fatalf("reasons = %v", result.Reasons)
	}

	report.Mismatched = 1
	report.Matched = 9
	result, err = EvaluateReadiness(report, StrictReadinessPolicy())
	if err != nil {
		t.Fatal(err)
	}
	if result.Ready || len(result.Reasons) == 0 {
		t.Fatal("mismatch unexpectedly passed readiness")
	}
}

func TestReadinessRejectsExistingCutover(t *testing.T) {
	report := CoverageReport{
		AsOf:             time.Date(2026, 8, 4, 4, 0, 0, 0, time.UTC),
		Checksum:         "checksum",
		EligibleLegacy:   1,
		Observed:         1,
		Matched:          1,
		RollbackEvidence: 1,
		CutoverMappings:  1,
		Statuses:         map[string]int{"matched": 1},
	}
	result, err := EvaluateReadiness(report, StrictReadinessPolicy())
	if err != nil {
		t.Fatal(err)
	}
	if result.Ready {
		t.Fatal("existing cutover unexpectedly passed readiness")
	}
}
