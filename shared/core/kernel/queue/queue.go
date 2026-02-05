package queue

import (
	"context"
	"time"
)

type JobPayload struct {
	ID      string
	Payload []byte // serialized envelope JSON
	Meta    map[string]string
}

type Queue interface {
	Push(ctx context.Context, queueName string, job JobPayload) error
	Later(ctx context.Context, queueName string, job JobPayload, delay time.Duration) error
	// Optionally Later: Pull/Consume with ack semantics, but we'll prefer Redis Streams consumer for that
	StartWorker(
		ctx context.Context,
		queue string,
		handler func(ctx context.Context, payload string) error,
	)
}
