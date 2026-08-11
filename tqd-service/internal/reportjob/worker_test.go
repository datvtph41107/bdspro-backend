package reportjob_test

import (
	"common/businessoperation"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"
	"tqd/infra/memoryreportjob"
	"tqd/internal/reportjob"
)

type fakeGenerator struct {
	output reportjob.Output
	err    error
}

func (f fakeGenerator) GenerateReport(context.Context, reportjob.Job) (reportjob.Output, error) {
	return f.output, f.err
}

func TestWorkerCompletesJob(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC)
	store := memoryreportjob.NewStore()
	job := reportjob.Job{
		ID:          "report_job_1",
		ReportID:    10,
		UserID:      42,
		Operation:   businessoperation.WorkspaceReportGenerate,
		OperationID: "op_report_worker_1",
		CommandKey:  "idem_report_worker_1",
		Status:      reportjob.StatusPending,
		AvailableAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if _, err := store.SaveJob(context.Background(), job); err != nil {
		t.Fatal(err)
	}

	worker := reportjob.NewWorker(
		store,
		fakeGenerator{output: reportjob.Output{PDFURL: "https://files.example/report.pdf", FileSize: 100}},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		"worker-1",
	)

	run, err := worker.RunOnce(context.Background())
	if err != nil || !run {
		t.Fatalf("RunOnce() = %v, %v", run, err)
	}

	saved, ok := store.Get(job.ID)
	if !ok || saved.Status != reportjob.StatusCompleted || saved.Attempts != 1 {
		t.Fatalf("saved job = %+v, %v", saved, ok)
	}
}

func TestWorkerSchedulesRetry(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC)
	store := memoryreportjob.NewStore()
	job := reportjob.Job{
		ID:          "report_job_2",
		ReportID:    11,
		UserID:      42,
		Operation:   businessoperation.WorkspaceReportGenerate,
		OperationID: "op_report_worker_2",
		CommandKey:  "idem_report_worker_2",
		Status:      reportjob.StatusPending,
		AvailableAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, _ = store.SaveJob(context.Background(), job)

	worker := reportjob.NewWorker(
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
	if saved.Status != reportjob.StatusPending || saved.Attempts != 1 || saved.LastError == "" {
		t.Fatalf("retry state = %+v", saved)
	}
}

func TestWorkerReleasesStaleJob(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 11, 5, 10, 0, 0, time.UTC)
	store := memoryreportjob.NewStore()
	job := reportjob.Job{
		ID:          "report_job_stale",
		ReportID:    12,
		UserID:      42,
		Operation:   businessoperation.WorkspaceReportGenerate,
		OperationID: "op_report_stale",
		CommandKey:  "idem_report_stale",
		Status:      reportjob.StatusPending,
		AvailableAt: now.Add(-time.Hour),
		CreatedAt:   now.Add(-time.Hour),
		UpdatedAt:   now.Add(-time.Hour),
	}
	_, _ = store.SaveJob(context.Background(), job)
	if _, found, err := store.ClaimNextJob(context.Background(), "dead-worker", now.Add(-10*time.Minute)); err != nil || !found {
		t.Fatalf("ClaimNextJob() found=%v err=%v", found, err)
	}

	worker := reportjob.NewWorker(
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
	if saved.Status != reportjob.StatusPending || saved.LockedAt != nil || saved.LockedBy != "" {
		t.Fatalf("released job = %+v", saved)
	}
}
