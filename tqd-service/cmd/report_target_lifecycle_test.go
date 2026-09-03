package cmd

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"tqd/internal/usecase/generatedreport/memoryjob"
	"tqd/internal/usecase/generatedreport/processing"
)

type checkpointGenerator struct{}

func (checkpointGenerator) GenerateReport(
	context.Context,
	processing.Job,
) (processing.Output, error) {
	return processing.Output{}, nil
}

func TestReportTargetWaitTracksStartedRunner(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	worker := processing.NewWorker(
		memoryjob.NewStore(),
		checkpointGenerator{},
		logger,
		"checkpoint-worker",
	)
	target := &reportTarget{
		jobRunner: processing.NewRunner(worker, time.Hour, logger),
	}

	ctx, cancel := context.WithCancel(context.Background())
	target.Start(ctx)
	cancel()

	done := make(chan struct{})
	go func() {
		target.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Report target Wait did not observe runner cancellation")
	}
}
