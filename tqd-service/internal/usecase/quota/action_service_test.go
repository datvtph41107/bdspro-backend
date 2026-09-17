package quota_test

import (
	"common/operation"
	"common/request"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
	"tqd/internal/access"
	"tqd/internal/usecase/quota"
	memoryquota "tqd/internal/usecase/quota/memory"
	"tqd/internal/usecase/quota/reconciliation/reservationcleanup"
	"tqd/internal/usecase/usage"
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
		if err := service.Run(ctx, decision.Operation, strings.Repeat("a", 64), func(context.Context) error {
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

func TestActionServiceReconcilesFailedReadWithoutUsage(t *testing.T) {
	decision := meteredAccess()
	runtimeStore := memoryquota.NewStore()
	durableStore := memoryusage.NewStore()
	service := quota.NewActionService(accessReader{result: decision}, quota.NewService(runtimeStore), durableStore)
	wantErr := errors.New("parcel unavailable")

	err := service.Run(actionContext(t, "idem_failed_read"), decision.Operation, strings.Repeat("a", 64), func(context.Context) error { return wantErr })
	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}
	usage, err := runtimeStore.GetUsage(context.Background(), decision.Subject, decision.Metering.MeterCode, decision.PeriodStart, decision.PeriodEnd)
	if err != nil {
		t.Fatal(err)
	}
	if usage.Used != 0 || usage.Reserved != 1 {
		t.Fatalf("runtime usage = %+v, want lease retained without spending", usage)
	}
	result, err := reservationcleanup.NewService(runtimeStore, durableStore).RunOnce(context.Background(), time.Now().Add(time.Hour), 10)
	if err != nil || result.Canceled != 1 || result.Committed != 0 {
		t.Fatalf("cleanup = %+v, err=%v", result, err)
	}
}

type failingCommitStore struct{ *memoryquota.Store }

func (s failingCommitStore) CommitQuota(context.Context, string) (quota.Reservation, error) {
	return quota.Reservation{}, errors.New("redis response lost")
}

func TestActionServiceReturnsResultAfterDurableUsageDespiteRedisFailure(t *testing.T) {
	decision := meteredAccess()
	runtimeStore := memoryquota.NewStore()
	durable := memoryusage.NewStore()
	service := quota.NewActionService(accessReader{decision}, quota.NewService(failingCommitStore{runtimeStore}), durable)
	err := service.Run(actionContext(t, "idem_commit_failure"), decision.Operation, strings.Repeat("a", 64), func(context.Context) error { return nil })
	if err != nil {
		t.Fatalf("accepted read must retain its successful response: %v", err)
	}
	result, err := reservationcleanup.NewService(runtimeStore, durable).RunOnce(context.Background(), time.Now().Add(time.Hour), 10)
	if err != nil || result.Committed != 1 || result.Canceled != 0 {
		t.Fatalf("cleanup = %+v, err=%v", result, err)
	}
}

type lostSaveResponse struct{ *memoryusage.Store }

func (s lostSaveResponse) AcceptReadCommand(ctx context.Context, command usage.ReadCommand, event *usage.Event) (string, error) {
	if _, err := s.Store.AcceptReadCommand(ctx, command, event); err != nil {
		return "", err
	}
	return "", errors.New("postgres response lost after commit")
}

func TestActionServiceRetainsLeaseOnUnknownDatabaseOutcome(t *testing.T) {
	decision := meteredAccess()
	runtimeStore := memoryquota.NewStore()
	durable := memoryusage.NewStore()
	service := quota.NewActionService(accessReader{decision}, quota.NewService(runtimeStore), lostSaveResponse{durable})
	err := service.Run(actionContext(t, "idem_unknown_save"), decision.Operation, strings.Repeat("a", 64), func(context.Context) error { return nil })
	if !errors.Is(err, quota.ErrActionUnavailable) {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := reservationcleanup.NewService(runtimeStore, durable).RunOnce(context.Background(), time.Now().Add(time.Hour), 10)
	if err != nil || result.Committed != 1 || result.Canceled != 0 {
		t.Fatalf("committed database usage must keep its quota: cleanup=%+v err=%v", result, err)
	}
}

// Hai request cùng command chia sẻ lease. Request lỗi không được giải phóng
// quota trong khi request còn lại đang chuẩn bị ghi usage.
func TestActionServiceFailedConcurrentRetryDoesNotCancelWinningRead(t *testing.T) {
	decision := meteredAccess()
	runtimeStore := memoryquota.NewStore()
	durable := memoryusage.NewStore()
	service := quota.NewActionService(accessReader{decision}, quota.NewService(runtimeStore), durable)
	ctx := actionContext(t, "idem_concurrent_read")
	entered := make(chan struct{})
	release := make(chan struct{})
	finished := make(chan error, 1)
	go func() {
		finished <- service.Run(ctx, decision.Operation, strings.Repeat("a", 64), func(context.Context) error {
			close(entered)
			<-release
			return nil
		})
	}()
	<-entered
	want := errors.New("second request failed")
	err := service.Run(ctx, decision.Operation, strings.Repeat("a", 64), func(context.Context) error { return want })
	close(release)
	winningErr := <-finished
	if !errors.Is(err, want) || winningErr != nil {
		t.Fatalf("failed retry=%v, winning read=%v", err, winningErr)
	}
	state, err := runtimeStore.GetUsage(ctx, decision.Subject, decision.Metering.MeterCode, decision.PeriodStart, decision.PeriodEnd)
	if err != nil || state.Used != 1 || state.Reserved != 0 {
		t.Fatalf("runtime usage=%+v, err=%v", state, err)
	}
}

func TestActionServiceDoesNotCommitNewPeriodReservationForOldUsage(t *testing.T) {
	decision := meteredAccess()
	runtimeStore := memoryquota.NewStore()
	durable := memoryusage.NewStore()
	ctx := actionContext(t, "idem_across_periods")
	for range 2 {
		service := quota.NewActionService(accessReader{decision}, quota.NewService(runtimeStore), durable)
		if err := service.Run(ctx, decision.Operation, strings.Repeat("a", 64), func(context.Context) error { return nil }); err != nil {
			t.Fatal(err)
		}
		decision.PeriodStart = decision.PeriodEnd
		decision.PeriodEnd = decision.PeriodStart.Add(30 * 24 * time.Hour)
	}
	result, err := reservationcleanup.NewService(runtimeStore, durable).RunOnce(context.Background(), time.Now().Add(time.Hour), 10)
	if err != nil || result.Canceled != 0 || result.Committed != 0 {
		t.Fatalf("old command must not consume new period: cleanup=%+v err=%v", result, err)
	}
}

func TestActionServiceRejectsMalformedAccessBeforeUnmeteredBypass(t *testing.T) {
	decision := meteredAccess()
	decision.Metering.MeterCode = ""
	service := quota.NewActionService(accessReader{decision}, quota.NewService(memoryquota.NewStore()), memoryusage.NewStore())
	err := service.Run(actionContext(t, "idem_invalid_access"), decision.Operation, strings.Repeat("a", 64), func(context.Context) error {
		t.Fatal("invalid entitlement must not execute business handler")
		return nil
	})
	if !errors.Is(err, quota.ErrActionUnavailable) {
		t.Fatalf("unexpected error: %v", err)
	}
}
