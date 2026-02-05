package facade

import (
	"context"
	"core/kernel/container"
	"time"
)

func Queue(c *container.Container) (QueueI, error) {
	return container.ResolveTyped[QueueI](c, "queue")
}

// func Event(c *container.Container) (EventDispatcher, error) {
// 	return container.ResolveTyped[EventDispatcher](c, "event")
// }

// func WS(c *container.Container) (Broadcaster, error) {
// 	return container.ResolveTyped[Broadcaster](c, "ws")
// }

type QueueI interface {
	PushSerialized(payload string) error
	PopSerialized() (id string, payload string, err error)
	Ack(id string) error
	StartWorker()
}

type Job interface {
	Name() string
	Handle(ctx context.Context) error
	MaxRetry() int
	Backoff(attempt int) time.Duration
}
