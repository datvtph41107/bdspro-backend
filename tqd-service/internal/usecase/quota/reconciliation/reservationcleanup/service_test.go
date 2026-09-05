package reservationcleanup

import (
	"common/operation"
	"context"
	"testing"
	"time"
	"tqd/internal/access"
	"tqd/internal/usecase/quota"
	"tqd/internal/usecase/quota/memory"
	"tqd/internal/usecase/usage"
	"tqd/internal/usecase/usage/memory"
)

func TestRunOnceCommitsExpiredReservationWhenUsageExists(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC)
	periodStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}
	quotaStore := memoryquota.NewStore()
	usageStore := memoryusage.NewStore()

	reservation, err := quotaStore.ReserveQuota(context.Background(), quota.StoreReserveInput{
		ReservationID: "r1",
		Subject:       subject,
		Operation:     operation.Code("workspace.report.generate"),
		MeterCode:     "workspace.report_generation.accepted",
		OperationID:   "op_cleanup_1",
		CommandKey:    "idem_cleanup_1",
		Amount:        1,
		Limit:         100,
		PeriodStart:   periodStart,
		PeriodEnd:     periodEnd,
		ExpiresAt:     now.Add(-time.Minute),
	})
	if err != nil {
		t.Fatal(err)
	}

	event := usage.NewEvent(usage.NewEventInput{
		Subject: subject, Operation: operation.Code("workspace.report.generate"), MeterCode: "workspace.report_generation.accepted",
		OperationID: reservation.OperationID, CommandKey: reservation.CommandKey, ReservationID: reservation.ID,
		Amount: 1, PeriodStart: periodStart, PeriodEnd: periodEnd, CreatedAt: now,
	})
	if _, err := usageStore.SaveUsageEvent(context.Background(), event); err != nil {
		t.Fatal(err)
	}

	result, err := NewService(quotaStore, usageStore).RunOnce(context.Background(), now, 10)
	if err != nil {
		t.Fatal(err)
	}
	if result.Committed != 1 || result.Canceled != 0 {
		t.Fatalf("result = %+v", result)
	}

	current, _ := quotaStore.GetUsage(context.Background(), subject, "workspace.report_generation.accepted", periodStart, periodEnd)
	if current.Used != 1 || current.Reserved != 0 {
		t.Fatalf("usage = %+v", current)
	}
}

func TestRunOnceCancelsExpiredReservationWithoutUsage(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC)
	periodStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}
	quotaStore := memoryquota.NewStore()

	_, err := quotaStore.ReserveQuota(context.Background(), quota.StoreReserveInput{
		ReservationID: "r2",
		Subject:       subject,
		Operation:     operation.Code("workspace.report.generate"),
		MeterCode:     "workspace.report_generation.accepted",
		OperationID:   "op_cleanup_2",
		CommandKey:    "idem_cleanup_2",
		Amount:        1,
		Limit:         100,
		PeriodStart:   periodStart,
		PeriodEnd:     periodEnd,
		ExpiresAt:     now.Add(-time.Minute),
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := NewService(quotaStore, memoryusage.NewStore()).RunOnce(context.Background(), now, 10)
	if err != nil {
		t.Fatal(err)
	}
	if result.Canceled != 1 || result.Committed != 0 {
		t.Fatalf("result = %+v", result)
	}
}

func TestRunOnceDoesNotReuseCommandUsageAcrossQuotaPeriods(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 9, 2, 5, 0, 0, 0, time.UTC)
	firstStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	firstEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	secondStart := firstEnd
	secondEnd := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}
	quotaStore := memoryquota.NewStore()
	usageStore := memoryusage.NewStore()

	first, err := quotaStore.ReserveQuota(context.Background(), quota.StoreReserveInput{
		ReservationID: "period-reservation-1",
		Subject:       subject,
		Operation:     operation.Code("workspace.report.generate"),
		MeterCode:     "workspace.report_generation.accepted",
		OperationID:   "op_period_1",
		CommandKey:    "idem_same_command",
		Amount:        1,
		Limit:         100,
		PeriodStart:   firstStart,
		PeriodEnd:     firstEnd,
		ExpiresAt:     now.Add(-time.Hour),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = quotaStore.ReserveQuota(context.Background(), quota.StoreReserveInput{
		ReservationID: "period-reservation-2",
		Subject:       subject,
		Operation:     operation.Code("workspace.report.generate"),
		MeterCode:     "workspace.report_generation.accepted",
		OperationID:   "op_period_2",
		CommandKey:    "idem_same_command",
		Amount:        1,
		Limit:         100,
		PeriodStart:   secondStart,
		PeriodEnd:     secondEnd,
		ExpiresAt:     now.Add(-time.Hour),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = usageStore.SaveUsageEvent(context.Background(), usage.NewEvent(usage.NewEventInput{
		Subject: subject, Operation: operation.Code("workspace.report.generate"), MeterCode: "workspace.report_generation.accepted",
		OperationID: first.OperationID, CommandKey: first.CommandKey, ReservationID: first.ID,
		Amount: 1, PeriodStart: firstStart, PeriodEnd: firstEnd, CreatedAt: now.Add(-time.Hour),
	}))
	if err != nil {
		t.Fatal(err)
	}

	result, err := NewService(quotaStore, usageStore).RunOnce(context.Background(), now, 10)
	if err != nil {
		t.Fatal(err)
	}
	if result.Committed != 1 || result.Canceled != 1 {
		t.Fatalf("result = %+v, want one commit and one cancel", result)
	}

	secondUsage, err := quotaStore.GetUsage(
		context.Background(),
		subject,
		"workspace.report_generation.accepted",
		secondStart,
		secondEnd,
	)
	if err != nil {
		t.Fatal(err)
	}
	if secondUsage.Used != 0 || secondUsage.Reserved != 0 {
		t.Fatalf("second period usage = %+v, want no charge", secondUsage)
	}
}
