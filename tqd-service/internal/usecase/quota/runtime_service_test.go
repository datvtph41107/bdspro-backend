package quota_test

import (
	"common/operation"
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"tqd/internal/access"
	"tqd/internal/usecase/quota"
	"tqd/internal/usecase/quota/memory"
)

func quotaAccess(limit int64) access.Result {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	return access.Result{
		Subject:   access.Subject{Type: access.SubjectProfile, ID: "42"},
		Operation: operation.Code("workspace.report.generate"),
		Allowed:   true,
		Metering: access.Metering{
			FeatureCode:    "workspace.report.generate",
			MeterCode:      "workspace.report_generation.accepted",
			UnitsPerAction: 1,
			PolicyVersion:  "1.0.0",
		},
		Limit:       limit,
		Period:      access.PeriodSubscriptionCycle,
		PeriodStart: start,
		PeriodEnd:   end,
	}
}

func TestReserveCommitAndRetry(t *testing.T) {
	store := memoryquota.NewStore()
	service := quota.NewService(store)
	ctx := context.Background()

	input := quota.ReserveInput{
		Access:         quotaAccess(100),
		OperationID:    "op_report_500",
		IdempotencyKey: "idem_report_700",
		Amount:         1,
		Now:            time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC),
	}

	first, err := service.ReserveQuota(ctx, input)
	if err != nil {
		t.Fatalf("ReserveQuota() error = %v", err)
	}
	if first.State != quota.StateReserved || first.Reserved != 1 {
		t.Fatalf("unexpected first reservation: %+v", first)
	}

	committed, err := service.CommitQuota(ctx, first)
	if err != nil {
		t.Fatalf("CommitQuota() error = %v", err)
	}
	if committed.State != quota.StateCommitted || committed.Used != 1 || committed.Reserved != 0 {
		t.Fatalf("unexpected committed reservation: %+v", committed)
	}

	retry, err := service.ReserveQuota(ctx, input)
	if err != nil {
		t.Fatalf("retry ReserveQuota() error = %v", err)
	}
	if retry.ID != first.ID || retry.State != quota.StateCommitted || retry.Used != 1 {
		t.Fatalf("retry created another charge: first=%+v retry=%+v", first, retry)
	}
}

func TestConcurrentLastQuotaOnlyAllowsOneReservation(t *testing.T) {
	store := memoryquota.NewStore()
	service := quota.NewService(store)
	ctx := context.Background()
	access := quotaAccess(100)

	if err := store.SetUsed(
		ctx,
		access.Subject,
		access.Metering.MeterCode,
		access.PeriodStart,
		access.PeriodEnd,
		99,
	); err != nil {
		t.Fatal(err)
	}

	var allowed int32
	var exceeded int32
	var wg sync.WaitGroup

	for i, id := range []string{"op_A", "op_B"} {
		wg.Add(1)
		go func(index int, operationID string) {
			defer wg.Done()
			_, err := service.ReserveQuota(ctx, quota.ReserveInput{
				Access:         access,
				OperationID:    operationID,
				IdempotencyKey: "idem_" + operationID,
				Amount:         1,
			})
			switch {
			case err == nil:
				atomic.AddInt32(&allowed, 1)
			case errors.Is(err, quota.ErrQuotaExceeded):
				atomic.AddInt32(&exceeded, 1)
			default:
				t.Errorf("request %d error = %v", index, err)
			}
		}(i, id)
	}
	wg.Wait()

	if allowed != 1 || exceeded != 1 {
		t.Fatalf("allowed=%d exceeded=%d, want 1/1", allowed, exceeded)
	}
}

func TestOperationsSharingMeterCompeteForOneAllowance(t *testing.T) {
	t.Parallel()

	store := memoryquota.NewStore()
	service := quota.NewService(store)
	firstAccess := quotaAccess(1)
	secondAccess := quotaAccess(1)
	secondAccess.Operation = operation.Code("parcel.analysis.run")
	secondAccess.Metering.FeatureCode = "parcel.analysis.run"

	first, err := service.ReserveQuota(context.Background(), quota.ReserveInput{
		Access:      firstAccess,
		OperationID: "op_report_shared_meter",
		Amount:      1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if first.MeterCode != firstAccess.Metering.MeterCode {
		t.Fatalf("reservation meter = %q, want %q", first.MeterCode, firstAccess.Metering.MeterCode)
	}
	if _, err := service.CommitQuota(context.Background(), first); err != nil {
		t.Fatal(err)
	}

	_, err = service.ReserveQuota(context.Background(), quota.ReserveInput{
		Access:      secondAccess,
		OperationID: "op_analysis_shared_meter",
		Amount:      1,
	})
	if !errors.Is(err, quota.ErrQuotaExceeded) {
		t.Fatalf("second operation error = %v, want ErrQuotaExceeded", err)
	}
}

func TestCancelReturnsReservedQuota(t *testing.T) {
	store := memoryquota.NewStore()
	service := quota.NewService(store)
	ctx := context.Background()

	reservation, err := service.ReserveQuota(ctx, quota.ReserveInput{
		Access:      quotaAccess(10),
		OperationID: "op_cancel_1",
		Amount:      1,
	})
	if err != nil {
		t.Fatal(err)
	}

	canceled, err := service.CancelQuota(ctx, reservation)
	if err != nil {
		t.Fatal(err)
	}
	if canceled.State != quota.StateCanceled || canceled.Reserved != 0 || canceled.Used != 0 {
		t.Fatalf("unexpected cancel result: %+v", canceled)
	}
}

func TestReserveQuotaCanRetryAfterCanceledReservation(t *testing.T) {
	t.Parallel()

	store := memoryquota.NewStore()
	service := quota.NewService(store)
	accessResult := quotaAccess(100)

	first, err := service.ReserveQuota(context.Background(), quota.ReserveInput{
		Access:         accessResult,
		OperationID:    "op_retry_after_cancel",
		IdempotencyKey: "idem_retry_after_cancel",
		Amount:         1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.CancelQuota(context.Background(), first); err != nil {
		t.Fatal(err)
	}

	second, err := service.ReserveQuota(context.Background(), quota.ReserveInput{
		Access:         accessResult,
		OperationID:    "op_retry_after_cancel",
		IdempotencyKey: "idem_retry_after_cancel",
		Amount:         1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if second.ID == first.ID {
		t.Fatalf("retry reused canceled reservation %q", first.ID)
	}
	if second.State != quota.StateReserved {
		t.Fatalf("second state = %q", second.State)
	}
}

func TestAllowedOperationWithoutQuotaDoesNotCreateReservation(t *testing.T) {
	t.Parallel()

	service := quota.NewService(memoryquota.NewStore())
	result, err := service.ReserveQuota(context.Background(), quota.ReserveInput{
		Access: access.Result{
			Subject:   access.Subject{Type: access.SubjectProfile, ID: "42"},
			Operation: operation.Code("parcel.view"),
			Allowed:   true,
			Metering: access.Metering{
				FeatureCode:   "parcel.view",
				PolicyVersion: "1.0.0",
			},
			Period: access.PeriodNone,
		},
		OperationID: "op_parcel_view_1",
	})
	if err != nil {
		t.Fatalf("ReserveQuota() error = %v", err)
	}
	if result.Required || result.ID != "" {
		t.Fatalf("reservation = %+v, want no quota reservation", result)
	}
}
