package usagecheck

import (
	commonaccess "common/access"
	"common/businessoperation"
	"context"
	"testing"
	"time"
	"tqd/infra/memoryquota"
	"tqd/infra/memoryusage"
	"tqd/internal/quota"
	"tqd/internal/usage"
)

type racingQuotaStore struct {
	used     int64
	reserved int64
	casCalls int
}

func (s *racingQuotaStore) GetUsage(
	context.Context,
	commonaccess.Subject,
	businessoperation.Code,
	time.Time,
	time.Time,
) (quota.Usage, error) {
	return quota.Usage{Used: s.used, Reserved: s.reserved}, nil
}

func (s *racingQuotaStore) CompareAndSetUsed(
	context.Context,
	commonaccess.Subject,
	businessoperation.Code,
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
	subject := commonaccess.Subject{Type: commonaccess.SubjectProfile, ID: "42"}
	operation := businessoperation.WorkspaceReportGenerate

	quotaStore := memoryquota.NewStore()
	usageStore := memoryusage.NewStore()
	_, err := usageStore.SaveUsageEvent(ctx, usage.NewEvent(
		subject,
		operation,
		"op_report_900",
		"idem_report_900",
		"idem_report_900",
		"reservation_900",
		1,
		periodStart,
		periodEnd,
		time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC),
	))
	if err != nil {
		t.Fatal(err)
	}

	service := NewService(quotaStore, usageStore)
	result, err := service.CheckAndFix(ctx, subject, operation, periodStart, periodEnd)
	if err != nil {
		t.Fatalf("CheckAndFix() error = %v", err)
	}
	if !result.Fixed || result.SavedUsed != 1 || result.RuntimeUsed != 0 {
		t.Fatalf("result = %+v, want saved=1 runtime=0 fixed=true", result)
	}

	runtime, err := quotaStore.GetUsage(ctx, subject, operation, periodStart, periodEnd)
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
	subject := commonaccess.Subject{Type: commonaccess.SubjectProfile, ID: "42"}
	operation := businessoperation.WorkspaceReportGenerate
	usageStore := memoryusage.NewStore()
	_, err := usageStore.SaveUsageEvent(ctx, usage.NewEvent(
		subject,
		operation,
		"op_report_race",
		"idem_report_race",
		"idem_report_race",
		"reservation_race",
		1,
		periodStart,
		periodEnd,
		time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC),
	))
	if err != nil {
		t.Fatal(err)
	}

	quotaStore := &racingQuotaStore{}
	result, err := NewService(quotaStore, usageStore).CheckAndFix(
		ctx,
		subject,
		operation,
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
