package filequeue

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunnerProcessesJobsAndStopsOnCancellation(t *testing.T) {
	queue := make(chan func(), 4)
	runner := New(queue, 2)
	ctx, cancel := context.WithCancel(context.Background())

	var count atomic.Int64
	jobDone := make(chan struct{})
	queue <- func() {
		if count.Add(1) == 2 {
			close(jobDone)
		}
	}
	queue <- func() {
		if count.Add(1) == 2 {
			close(jobDone)
		}
	}

	runDone := make(chan struct{})
	go func() {
		runner.Run(ctx)
		close(runDone)
	}()

	select {
	case <-jobDone:
	case <-time.After(time.Second):
		t.Fatal("jobs were not processed")
	}
	cancel()
	select {
	case <-runDone:
	case <-time.After(time.Second):
		t.Fatal("runner did not stop after cancellation")
	}
	if got := count.Load(); got != 2 {
		t.Fatalf("processed jobs = %d, want 2", got)
	}
}

func TestRunnerWaitsForInflightJobBeforeReturn(t *testing.T) {
	queue := make(chan func(), 1)
	runner := New(queue, 1)
	ctx, cancel := context.WithCancel(context.Background())

	jobStarted := make(chan struct{})
	releaseJob := make(chan struct{})
	queue <- func() {
		close(jobStarted)
		<-releaseJob
	}

	runDone := make(chan struct{})
	go func() {
		runner.Run(ctx)
		close(runDone)
	}()
	<-jobStarted
	cancel()
	select {
	case <-runDone:
		t.Fatal("runner returned before in-flight job completed")
	default:
	}
	close(releaseJob)
	select {
	case <-runDone:
	case <-time.After(time.Second):
		t.Fatal("runner did not return after in-flight job completed")
	}
}

func TestRunnerDrainsAcceptedBufferedJobsAfterCancellation(t *testing.T) {
	queue := make(chan func(), 4)
	runner := New(queue, 1)
	ctx, cancel := context.WithCancel(context.Background())

	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})
	var processed atomic.Int64
	queue <- func() {
		close(firstStarted)
		<-releaseFirst
		processed.Add(1)
	}
	queue <- func() { processed.Add(1) }
	queue <- func() { processed.Add(1) }

	done := make(chan struct{})
	go func() {
		runner.Run(ctx)
		close(done)
	}()

	<-firstStarted
	cancel()
	close(releaseFirst)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runner did not drain accepted jobs")
	}
	if got := processed.Load(); got != 3 {
		t.Fatalf("processed jobs = %d, want all 3 accepted jobs", got)
	}
}
