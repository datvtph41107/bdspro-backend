package processing_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"
	"tqd/internal/usecase/generatedreport/memoryjob"
	"tqd/internal/usecase/generatedreport/processing"
)

type fakeGenerator struct {
	output processing.Output
	err    error
}

func (f fakeGenerator) GenerateReport(context.Context, processing.Job) (processing.Output, error) {
	return f.output, f.err
}

func TestWorkerCompletesJob(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC)
	store := memoryjob.NewStore()
	job := processing.Job{
		ID:          "report_job_1",
		ReportID:    10,
		UserID:      42,
		Operation:   "workspace.report.generate",
		OperationID: "op_report_worker_1",
		CommandKey:  "idem_report_worker_1",
		Status:      processing.StatusPending,
		AvailableAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if _, err := store.SaveJob(context.Background(), job); err != nil {
		t.Fatal(err)
	}

	worker := processing.NewWorker(
		store,
		fakeGenerator{output: processing.Output{PDFURL: "https://files.example/report.pdf", FileSize: 100}},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		"worker-1",
	)

	run, err := worker.RunOnce(context.Background())
	if err != nil || !run {
		t.Fatalf("RunOnce() = %v, %v", run, err)
	}

	saved, ok := store.Get(job.ID)
	if !ok || saved.Status != processing.StatusCompleted || saved.Attempts != 1 {
		t.Fatalf("saved job = %+v, %v", saved, ok)
	}
}

func TestWorkerSchedulesRetry(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC)
	store := memoryjob.NewStore()
	job := processing.Job{
		ID:          "report_job_2",
		ReportID:    11,
		UserID:      42,
		Operation:   "workspace.report.generate",
		OperationID: "op_report_worker_2",
		CommandKey:  "idem_report_worker_2",
		Status:      processing.StatusPending,
		AvailableAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, _ = store.SaveJob(context.Background(), job)

	worker := processing.NewWorker(
		store,
		fakeGenerator{err: errors.New("renderer unavailable")},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		"worker-1",
	)

	run, err := worker.RunOnce(context.Background())
	if err != nil || !run {
		t.Fatalf("RunOnce() = %v, %v", run, err)
	}

	saved, _ := store.Get(job.ID)
	if saved.Status != processing.StatusPending || saved.Attempts != 1 || saved.LastError == "" {
		t.Fatalf("retry state = %+v", saved)
	}
}

func TestWorkerReleasesStaleJob(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 11, 5, 10, 0, 0, time.UTC)
	store := memoryjob.NewStore()
	job := processing.Job{
		ID:          "report_job_stale",
		ReportID:    12,
		UserID:      42,
		Operation:   "workspace.report.generate",
		OperationID: "op_report_stale",
		CommandKey:  "idem_report_stale",
		Status:      processing.StatusPending,
		AvailableAt: now.Add(-time.Hour),
		CreatedAt:   now.Add(-time.Hour),
		UpdatedAt:   now.Add(-time.Hour),
	}
	_, _ = store.SaveJob(context.Background(), job)
	if _, found, err := store.ClaimNextJob(context.Background(), "dead-worker", now.Add(-10*time.Minute)); err != nil || !found {
		t.Fatalf("ClaimNextJob() found=%v err=%v", found, err)
	}

	worker := processing.NewWorker(
		store,
		fakeGenerator{},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		"worker-1",
	)
	released, err := worker.ReleaseStale(context.Background(), 5*time.Minute, 100)
	if err != nil || released != 1 {
		t.Fatalf("ReleaseStale() = %d, %v", released, err)
	}

	saved, _ := store.Get(job.ID)
	if saved.Status != processing.StatusPending || saved.LockedAt != nil || saved.LockedBy != "" {
		t.Fatalf("released job = %+v", saved)
	}
}

func TestStaleClaimCannotFinalizeAfterReclaim(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	t0 := time.Date(2026, 8, 19, 8, 0, 0, 0, time.UTC)
	store := memoryjob.NewStore()
	job := processing.Job{
		ID:          "report_job_fencing",
		ReportID:    99,
		UserID:      42,
		Operation:   "workspace.report.generate",
		OperationID: "op_report_fencing",
		CommandKey:  "idem_report_fencing",
		Status:      processing.StatusPending,
		AvailableAt: t0,
		CreatedAt:   t0,
		UpdatedAt:   t0,
	}
	if _, err := store.SaveJob(ctx, job); err != nil {
		t.Fatal(err)
	}

	claimA, found, err := store.ClaimNextJob(ctx, "worker-a", t0)
	if err != nil || !found {
		t.Fatalf("claim A found=%v err=%v", found, err)
	}
	if claimA.ClaimVersion != 1 {
		t.Fatalf("claim A version=%d, want 1", claimA.ClaimVersion)
	}

	released, err := store.ReleaseStaleJobs(ctx, t0.Add(time.Second), 1)
	if err != nil || released != 1 {
		t.Fatalf("release stale=%d err=%v", released, err)
	}

	claimB, found, err := store.ClaimNextJob(ctx, "worker-b", t0.Add(2*time.Second))
	if err != nil || !found {
		t.Fatalf("claim B found=%v err=%v", found, err)
	}
	if claimB.ClaimVersion != 2 {
		t.Fatalf("claim B version=%d, want 2", claimB.ClaimVersion)
	}

	if err := store.CompleteJob(ctx, claimA, processing.Output{}, t0.Add(3*time.Second)); !errors.Is(err, processing.ErrClaimLost) {
		t.Fatalf("stale CompleteJob error=%v, want ErrClaimLost", err)
	}
	if err := store.RetryJob(ctx, claimA, "stale", t0.Add(4*time.Second)); !errors.Is(err, processing.ErrClaimLost) {
		t.Fatalf("stale RetryJob error=%v, want ErrClaimLost", err)
	}
	if err := store.FailJob(ctx, claimA, "stale", t0.Add(5*time.Second)); !errors.Is(err, processing.ErrClaimLost) {
		t.Fatalf("stale FailJob error=%v, want ErrClaimLost", err)
	}

	if err := store.CompleteJob(ctx, claimB, processing.Output{PDFURL: "https://files.example/final.pdf"}, t0.Add(6*time.Second)); err != nil {
		t.Fatalf("current claimant CompleteJob error=%v", err)
	}
}
