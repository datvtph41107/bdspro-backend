package memoryreportjob

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"
	"tqd/internal/reportjob"
)

/**
 * Store là in-memory report job queue dùng cho test worker flow.
 */
type Store struct {
	mu sync.Mutex

	jobs    map[string]reportjob.Job
	outputs map[string]reportjob.Output
}

func NewStore() *Store {
	return &Store{
		jobs:    make(map[string]reportjob.Job),
		outputs: make(map[string]reportjob.Output),
	}
}

func (s *Store) SaveJob(ctx context.Context, job reportjob.Job) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.jobs[job.ID]; exists {
		return false, nil
	}
	s.jobs[job.ID] = job
	return true, nil
}

func (s *Store) ClaimNextJob(
	ctx context.Context,
	workerID string,
	now time.Time,
) (reportjob.Job, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ids := make([]string, 0, len(s.jobs))
	for id := range s.jobs {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		job := s.jobs[id]
		if job.Status != reportjob.StatusPending || job.AvailableAt.After(now) {
			continue
		}
		job.Status = reportjob.StatusRunning
		job.Attempts++
		job.LockedAt = timePtr(now)
		job.LockedBy = workerID
		job.UpdatedAt = now
		s.jobs[id] = job
		return job, true, nil
	}
	return reportjob.Job{}, false, nil
}

func (s *Store) CompleteJob(
	ctx context.Context,
	job reportjob.Job,
	output reportjob.Output,
	now time.Time,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.jobs[job.ID]
	if !ok || current.Status != reportjob.StatusRunning {
		return errors.New("report job is no longer running")
	}
	current.Status = reportjob.StatusCompleted
	current.LockedAt = nil
	current.LockedBy = ""
	current.LastError = ""
	current.UpdatedAt = now
	s.jobs[job.ID] = current
	s.outputs[job.ID] = output
	return nil
}

func (s *Store) RetryJob(
	ctx context.Context,
	job reportjob.Job,
	lastError string,
	availableAt time.Time,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.jobs[job.ID]
	if !ok || current.Status != reportjob.StatusRunning {
		return errors.New("report job is no longer running")
	}
	current.Status = reportjob.StatusPending
	current.AvailableAt = availableAt
	current.LockedAt = nil
	current.LockedBy = ""
	current.LastError = lastError
	current.UpdatedAt = availableAt
	s.jobs[job.ID] = current
	return nil
}

func (s *Store) FailJob(
	ctx context.Context,
	job reportjob.Job,
	lastError string,
	now time.Time,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.jobs[job.ID]
	if !ok || current.Status != reportjob.StatusRunning {
		return errors.New("report job is no longer running")
	}
	current.Status = reportjob.StatusFailed
	current.LockedAt = nil
	current.LockedBy = ""
	current.LastError = lastError
	current.UpdatedAt = now
	s.jobs[job.ID] = current
	return nil
}

func (s *Store) ReleaseStaleJobs(
	ctx context.Context,
	lockedBefore time.Time,
	limit int,
) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 {
		limit = 100
	}

	count := 0
	for id, job := range s.jobs {
		if count >= limit {
			break
		}
		if job.Status != reportjob.StatusRunning || job.LockedAt == nil || !job.LockedAt.Before(lockedBefore) {
			continue
		}
		job.Status = reportjob.StatusPending
		job.AvailableAt = lockedBefore
		job.LockedAt = nil
		job.LockedBy = ""
		job.UpdatedAt = lockedBefore
		s.jobs[id] = job
		count++
	}
	return count, nil
}

func (s *Store) Get(id string) (reportjob.Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, ok := s.jobs[id]
	return job, ok
}

func timePtr(value time.Time) *time.Time {
	return &value
}
