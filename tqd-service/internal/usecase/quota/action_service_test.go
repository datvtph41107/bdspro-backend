package quota_test

import (
	"common/operation"
	"common/request"
	"context"
	"errors"
	"testing"
	"time"
	"tqd/internal/access"
	"tqd/internal/usecase/quota"
	memoryquota "tqd/internal/usecase/quota/memory"
	memoryusage "tqd/internal/usecase/usage/memory"
)

type accessReader struct{ result access.Result }

func (r accessReader) GetAccess(context.Context, operation.Code) (access.Result, error) {
	return r.result, nil
}

func actionContext(t *testing.T, idempotencyKey string) context.Context {
	t.Helper()
	ctx, err := request.BindOperationID(context.Background(), "op_action_test")
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = request.BindIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}

func meteredAccess() access.Result {
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	return access.Result{
		Subject:        access.Subject{Type: access.SubjectProfile, ID: "42"},
		Operation:      quota.ParcelDetailViewOperation,
		Allowed:        true,
		Limit:          2,
		Period:         access.PeriodSubscriptionCycle,
		PeriodStart:    start,
		PeriodEnd:      start.Add(30 * 24 * time.Hour),
		SubscriptionID: 7,
		PlanCode:       "qhpro.standard",
		PlanVersion:    "1.0.0",
		Metering: access.Metering{
			FeatureCode:    string(quota.ParcelDetailViewOperation),
			MeterCode:      "workspace.parcel_detail.viewed",
			UnitsPerAction: 1,
			PolicyVersion:  "1.0.0",
		},
	}
}

func TestActionServiceAccountsSuccessfulIntentOnce(t *testing.T) {
	decision := meteredAccess()
	runtimeStore := memoryquota.NewStore()
	durableStore := memoryusage.NewStore()
	service := quota.NewActionService(accessReader{result: decision}, quota.NewService(runtimeStore), durableStore)
	ctx := actionContext(t, "idem_parcel_42")

	called := 0
	for range 2 {
		if err := service.Run(ctx, decision.Operation, func(context.Context) error {
			called++
			return nil
		}); err != nil {
			t.Fatalf("Run() error = %v", err)
		}
	}

	used, err := durableStore.SumUsage(ctx, decision.Subject, decision.Metering.MeterCode, decision.PeriodStart, decision.PeriodEnd)
	if err != nil {
		t.Fatal(err)
	}
	if called != 2 || used != 1 {
		t.Fatalf("action calls=%d durable used=%d, want 2 and 1", called, used)
	}
}

func TestActionServiceCancelsFailedRead(t *testing.T) {
	decision := meteredAccess()
	runtimeStore := memoryquota.NewStore()
	durableStore := memoryusage.NewStore()
	service := quota.NewActionService(accessReader{result: decision}, quota.NewService(runtimeStore), durableStore)
	wantErr := errors.New("parcel unavailable")

	err := service.Run(actionContext(t, "idem_failed_read"), decision.Operation, func(context.Context) error { return wantErr })
	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}
	usage, err := runtimeStore.GetUsage(context.Background(), decision.Subject, decision.Metering.MeterCode, decision.PeriodStart, decision.PeriodEnd)
	if err != nil {
		t.Fatal(err)
	}
	if usage.Used != 0 || usage.Reserved != 0 {
		t.Fatalf("runtime usage = %+v, want empty", usage)
	}
}
