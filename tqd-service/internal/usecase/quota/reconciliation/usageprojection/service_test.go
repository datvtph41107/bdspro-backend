package usageprojection

import (
	commonmetering "common/metering"
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

type racingQuotaStore struct {
	used     int64
	reserved int64
	casCalls int
}

func (s *racingQuotaStore) GetUsage(
	context.Context,
	access.Subject,
	commonmetering.Code,
	time.Time,
	time.Time,
) (quota.Usage, error) {
	return quota.Usage{Used: s.used, Reserved: s.reserved}, nil
}

func (s *racingQuotaStore) CompareAndSetUsed(
	context.Context,
	access.Subject,
	commonmetering.Code,
	time.Time,
	time.Time,
	int64,
	int64,
) (bool, error) {
	s.casCalls++
	// Simulate a quota commit winning immediately before the atomic CAS.
	s.used++
	return false, nil
}

func TestCheckAndFixRestoresSavedUsage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	periodStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}
	operation := operation.Code("workspace.report.generate")
	meterCode := commonmetering.Code("workspace.report_generation.accepted")

	quotaStore := memoryquota.NewStore()
	usageStore := memoryusage.NewStore()
	_, err := usageStore.SaveUsageEvent(ctx, usage.NewEvent(usage.NewEventInput{
		Subject: subject, Operation: operation, MeterCode: meterCode,
		OperationID: "op_report_900", IdempotencyKey: "idem_report_900", CommandKey: "idem_report_900",
		ReservationID: "reservation_900", Amount: 1, PeriodStart: periodStart, PeriodEnd: periodEnd,
		CreatedAt: time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC),
	}))
	if err != nil {
		t.Fatal(err)
	}

	service := NewService(quotaStore, usageStore)
	result, err := service.CheckAndFix(ctx, subject, meterCode, periodStart, periodEnd)
	if err != nil {
		t.Fatalf("CheckAndFix() error = %v", err)
	}
	if !result.Fixed || result.SavedUsed != 1 || result.RuntimeUsed != 0 {
		t.Fatalf("result = %+v, want saved=1 runtime=0 fixed=true", result)
	}

	runtime, err := quotaStore.GetUsage(ctx, subject, meterCode, periodStart, periodEnd)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.Used != 1 || runtime.Reserved != 0 {
		t.Fatalf("runtime usage = %+v, want used=1 reserved=0", runtime)
	}
}

func TestCheckAndFixDoesNotOverwriteConcurrentCommit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	periodStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}
	operation := operation.Code("workspace.report.generate")
	meterCode := commonmetering.Code("workspace.report_generation.accepted")
	usageStore := memoryusage.NewStore()
	_, err := usageStore.SaveUsageEvent(ctx, usage.NewEvent(usage.NewEventInput{
		Subject: subject, Operation: operation, MeterCode: meterCode,
		OperationID: "op_report_race", IdempotencyKey: "idem_report_race", CommandKey: "idem_report_race",
		ReservationID: "reservation_race", Amount: 1, PeriodStart: periodStart, PeriodEnd: periodEnd,
		CreatedAt: time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC),
	}))
	if err != nil {
		t.Fatal(err)
	}

	quotaStore := &racingQuotaStore{}
	result, err := NewService(quotaStore, usageStore).CheckAndFix(
		ctx,
		subject,
		meterCode,
		periodStart,
		periodEnd,
	)
	if err != nil {
		t.Fatal(err)
	}
	if result.Fixed {
		t.Fatalf("concurrent value must not be reported as overwritten: %+v", result)
	}
	if quotaStore.used != 1 || quotaStore.casCalls != 1 {
		t.Fatalf("quota store state used=%d cas_calls=%d", quotaStore.used, quotaStore.casCalls)
	}
}
