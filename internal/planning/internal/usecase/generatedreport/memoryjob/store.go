package memoryjob

import (
	"context"
	"sort"
	"sync"
	"time"
	"tqd/internal/usecase/generatedreport/processing"
)

/**
 * Store là in-memory report job queue dùng cho test worker flow.
 */
type Store struct {
	mu sync.Mutex

	jobs    map[string]processing.Job
	outputs map[string]processing.Output
}

func NewStore() *Store {
	return &Store{
		jobs:    make(map[string]processing.Job),
		outputs: make(map[string]processing.Output),
	}
}

func (s *Store) SaveJob(ctx context.Context, job processing.Job) (bool, error) {
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
) (processing.Job, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ids := make([]string, 0, len(s.jobs))
	for id := range s.jobs {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		job := s.jobs[id]
		if job.Status != processing.StatusPending || job.AvailableAt.After(now) {
			continue
		}
		job.Status = processing.StatusRunning
		job.Attempts++
		job.ClaimVersion++
		job.LockedAt = timePtr(now)
		job.LockedBy = workerID
		job.UpdatedAt = now
		s.jobs[id] = job
		return job, true, nil
	}
	return processing.Job{}, false, nil
}

func (s *Store) CompleteJob(
	ctx context.Context,
	job processing.Job,
	output processing.Output,
	now time.Time,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.jobs[job.ID]
	if !ok ||
		current.Status != processing.StatusRunning ||
		current.LockedBy != job.LockedBy ||
		current.ClaimVersion != job.ClaimVersion {
		return processing.ErrClaimLost
	}
	current.Status = processing.StatusCompleted
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
	job processing.Job,
	lastError string,
	availableAt time.Time,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.jobs[job.ID]
	if !ok ||
		current.Status != processing.StatusRunning ||
		current.LockedBy != job.LockedBy ||
		current.ClaimVersion != job.ClaimVersion {
		return processing.ErrClaimLost
	}
	current.Status = processing.StatusPending
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
	job processing.Job,
	lastError string,
	now time.Time,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.jobs[job.ID]
	if !ok ||
		current.Status != processing.StatusRunning ||
		current.LockedBy != job.LockedBy ||
		current.ClaimVersion != job.ClaimVersion {
		return processing.ErrClaimLost
	}
	current.Status = processing.StatusFailed
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
		if job.Status != processing.StatusRunning || job.LockedAt == nil || !job.LockedAt.Before(lockedBefore) {
			continue
		}
		job.Status = processing.StatusPending
		job.AvailableAt = lockedBefore
		job.LockedAt = nil
		job.LockedBy = ""
		job.UpdatedAt = lockedBefore
		s.jobs[id] = job
		count++
	}
	return count, nil
}

func (s *Store) Get(id string) (processing.Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, ok := s.jobs[id]
	return job, ok
}

func timePtr(value time.Time) *time.Time {
	return &value
}
