package processing

import (
	"context"
	"time"
)

/**
 * Store giữ durable job queue trong PostgreSQL.
 */
type Store interface {
	SaveJob(ctx context.Context, job Job) (created bool, err error)

	ClaimNextJob(
		ctx context.Context,
		workerID string,
		now time.Time,
	) (Job, bool, error)

	CompleteJob(
		ctx context.Context,
		job Job,
		output Output,
		now time.Time,
	) error

	RetryJob(
		ctx context.Context,
		job Job,
		lastError string,
		availableAt time.Time,
	) error

	FailJob(
		ctx context.Context,
		job Job,
		lastError string,
		now time.Time,
	) error

	ReleaseStaleJobs(
		ctx context.Context,
		lockedBefore time.Time,
		limit int,
	) (int, error)
}
