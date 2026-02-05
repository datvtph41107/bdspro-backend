package kernel

import (
	"context"
	"core/kernel/container"
	"core/kernel/queue"
	"time"
)

var (
	appC *container.Container
)

func Queue() queue.Queue {
	rsv, err := container.ResolveTyped[queue.Queue](appC, "queue")
	if err != nil {
		panic(err)
	}
	return rsv
}

func Push(
	ctx context.Context,
	queueName string,
	job queue.JobPayload,
) error {
	return Queue().Push(ctx, queueName, job)
}

func Later(
	ctx context.Context,
	queueName string,
	job queue.JobPayload,
	delay time.Duration,
) error {
	return Queue().Later(ctx, queueName, job, delay)
}
