package usecase

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const spamThreshold = 5 * time.Second

// ErrBuildInProgress is returned when a build is already running for the same layer
// and the spam threshold has not been exceeded.
type ErrBuildInProgress struct {
	LayerID  uint64
	Progress BuildProgress
}

func (e *ErrBuildInProgress) Error() string {
	return fmt.Sprintf("build already in progress for layer %d: %s", e.LayerID, e.Progress.Message)
}

// ErrFamilyBuildInProgress is returned when a family PMTiles build is already running.
type ErrFamilyBuildInProgress struct {
	FamilyID uint64
	Progress BuildProgress
}

func (e *ErrFamilyBuildInProgress) Error() string {
	return fmt.Sprintf("build already in progress for family %d: %s", e.FamilyID, e.Progress.Message)
}

type buildJob struct {
	startTime time.Time
	progress  atomic.Value // BuildProgress
	cancel    context.CancelFunc
	gen       int64
}

type buildTracker struct {
	mu     sync.Mutex
	builds map[uint64]*buildJob
}

func newBuildTracker() *buildTracker {
	return &buildTracker{
		builds: make(map[uint64]*buildJob),
	}
}

// tryStart attempts to register a build job for the given layer.
// Returns a context for the build, its generation, and an error if blocked.
func (t *buildTracker) tryStart(layerID uint64) (context.Context, int64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	job, exists := t.builds[layerID]
	if !exists {
		ctx, gen := t.createJob(layerID)
		return ctx, gen, nil
	}

	if elapsed := time.Since(job.startTime); elapsed < spamThreshold {
		p, _ := job.progress.Load().(BuildProgress)
		return nil, 0, &ErrBuildInProgress{LayerID: layerID, Progress: p}
	}

	// Cancel the old build before replacing it.
	if job.cancel != nil {
		job.cancel()
	}

	ctx, gen := t.replaceJob(job)
	return ctx, gen, nil
}

func (t *buildTracker) createJob(layerID uint64) (context.Context, int64) {
	bgCtx, cancel := context.WithCancel(context.Background())
	gen := int64(1)
	job := &buildJob{
		startTime: time.Now(),
		cancel:    cancel,
		gen:       gen,
	}
	job.progress.Store(BuildProgress{Status: "processing", Message: "Starting..."})
	t.builds[layerID] = job
	return bgCtx, gen
}

func (t *buildTracker) replaceJob(job *buildJob) (context.Context, int64) {
	bgCtx, cancel := context.WithCancel(context.Background())
	gen := atomic.AddInt64(&job.gen, 1)
	job.startTime = time.Now()
	job.cancel = cancel
	job.progress.Store(BuildProgress{Status: "processing", Message: "Starting..."})
	return bgCtx, gen
}

// updateProgress only applies if the generation matches — prevents stale builds
// from overwriting the progress of a newer build.
func (t *buildTracker) updateProgress(layerID uint64, gen int64, p BuildProgress) {
	t.mu.Lock()
	defer t.mu.Unlock()

	job, exists := t.builds[layerID]
	if !exists || atomic.LoadInt64(&job.gen) != gen {
		return
	}
	job.progress.Store(p)
}

// finish only removes the entry if the generation matches — prevents stale builds
// from deleting the entry of a newer build.
func (t *buildTracker) finish(layerID uint64, gen int64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	job, exists := t.builds[layerID]
	if !exists || atomic.LoadInt64(&job.gen) != gen {
		return
	}
	job.cancel = nil
	delete(t.builds, layerID)
}

// Ensure ErrBuildInProgress implements the error interface.
var _ error = (*ErrBuildInProgress)(nil)
