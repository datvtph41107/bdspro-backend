package catalogmigration

import (
	"testing"
	"time"
)

const testPlanTermsChecksum = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func TestBuildSubscriptionDraftPreservesLegacyPeriod(t *testing.T) {
	paidAt := time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)
	finishedAt := paidAt.AddDate(1, 0, 0)
	createdAt := paidAt.Add(-24 * time.Hour)
	updatedAt := paidAt
	draft, err := BuildSubscriptionDraft(
		LegacyMembershipSnapshot{
			ID:           9,
			ProfileID:    42,
			PlanID:       7,
			FinishedDate: &finishedAt,
			OrderID:      88,
			PaidAt:       &paidAt,
			CreatedAt:    &createdAt,
			UpdatedAt:    &updatedAt,
		},
		MembershipDecision{
			ProductCode:       "qhpro",
			PlanCode:          "qhpro.legacy_pro",
			PlanVersion:       "1.0.0",
			PlanTermsChecksum: testPlanTermsChecksum,
			PeriodHandling:    "paid_at_to_finished_date",
		},
		paidAt.AddDate(0, 1, 0),
	)
	if err != nil {
		t.Fatalf("BuildSubscriptionDraft() error = %v", err)
	}
	if draft.Status != "active" {
		t.Fatalf("status = %q", draft.Status)
	}
	if draft.SubjectKind != "profile" || draft.SubjectID != "42" {
		t.Fatalf("subject = %s:%s", draft.SubjectKind, draft.SubjectID)
	}
	if draft.CurrentPeriodStart != paidAt || draft.CurrentPeriodEnd != finishedAt {
		t.Fatal("legacy period changed during migration")
	}
	if draft.SubscriptionKey != "sub_legacy_membership_9" {
		t.Fatalf("subscription key = %q", draft.SubscriptionKey)
	}
}

func TestBuildSubscriptionDraftMarksExpiredAtExclusiveEnd(t *testing.T) {
	paidAt := time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)
	finishedAt := paidAt.AddDate(0, 1, 0)
	createdAt := paidAt.Add(-24 * time.Hour)
	updatedAt := paidAt
	draft, err := BuildSubscriptionDraft(
		LegacyMembershipSnapshot{
			ID:           9,
			ProfileID:    42,
			PlanID:       7,
			FinishedDate: &finishedAt,
			PaidAt:       &paidAt,
			CreatedAt:    &createdAt,
			UpdatedAt:    &updatedAt,
		},
		MembershipDecision{
			ProductCode:       "qhpro",
			PlanCode:          "qhpro.legacy_pro",
			PlanVersion:       "1.0.0",
			PlanTermsChecksum: testPlanTermsChecksum,
			PeriodHandling:    "paid_at_to_finished_date",
		},
		finishedAt,
	)
	if err != nil {
		t.Fatal(err)
	}
	if draft.Status != "expired" {
		t.Fatalf("status = %q, want expired", draft.Status)
	}
}

func TestBuildSubscriptionDraftRejectsIncompleteHistory(t *testing.T) {
	if _, err := BuildSubscriptionDraft(
		LegacyMembershipSnapshot{ID: 1, ProfileID: 2, PlanID: 3},
		MembershipDecision{
			ProductCode:       "qhpro",
			PlanCode:          "qhpro.free",
			PlanVersion:       "1.0.0",
			PlanTermsChecksum: testPlanTermsChecksum,
			PeriodHandling:    "paid_at_to_finished_date",
		},
		time.Now(),
	); err == nil {
		t.Fatal("BuildSubscriptionDraft() accepted missing period evidence")
	}
}

func TestBuildSubscriptionDraftMarksFutureMembershipPending(t *testing.T) {
	paidAt := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	finishedAt := paidAt.AddDate(0, 1, 0)
	createdAt := paidAt.Add(-24 * time.Hour)
	updatedAt := paidAt
	draft, err := BuildSubscriptionDraft(
		LegacyMembershipSnapshot{
			ID:           9,
			ProfileID:    42,
			PlanID:       7,
			FinishedDate: &finishedAt,
			PaidAt:       &paidAt,
			CreatedAt:    &createdAt,
			UpdatedAt:    &updatedAt,
		},
		MembershipDecision{
			ProductCode:       "qhpro",
			PlanCode:          "qhpro.legacy_pro",
			PlanVersion:       "1.0.0",
			PlanTermsChecksum: testPlanTermsChecksum,
			PeriodHandling:    "paid_at_to_finished_date",
		},
		paidAt.Add(-time.Hour),
	)
	if err != nil {
		t.Fatal(err)
	}
	if draft.Status != "pending" {
		t.Fatalf("status = %q, want pending", draft.Status)
	}
	if draft.OrderReference != "" {
		t.Fatalf("order reference = %q, want empty", draft.OrderReference)
	}
}

func TestMembershipSourceChecksumChangesWithPeriod(t *testing.T) {
	paidAt := time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)
	finishedAt := paidAt.AddDate(0, 1, 0)
	createdAt := paidAt.Add(-24 * time.Hour)
	updatedAt := paidAt
	snapshot := LegacyMembershipSnapshot{
		ID: 1, ProfileID: 2, PlanID: 3,
		PaidAt: &paidAt, FinishedDate: &finishedAt,
		CreatedAt: &createdAt, UpdatedAt: &updatedAt,
	}
	first, err := MembershipSourceChecksum(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	changed := finishedAt.Add(time.Hour)
	snapshot.FinishedDate = &changed
	second, err := MembershipSourceChecksum(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("membership checksum ignored period change")
	}
}
