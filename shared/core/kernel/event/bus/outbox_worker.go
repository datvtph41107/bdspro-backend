package bus

import (
	"context"
	"core/kernel/event/outbox"
	"core/kernel/queue"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type OutboxWorker struct {
	Repo     *outbox.Repo
	RDB      *redis.Client
	Producer *queue.StreamProducer
	Batch    int
	Interval time.Duration
	cancel   context.CancelFunc
}

func NewOutboxWorker(repo *outbox.Repo, rdb *redis.Client) *OutboxWorker {
	return &OutboxWorker{
		Repo:     repo,
		RDB:      rdb,
		Producer: queue.NewStreamProducer(rdb),
		Batch:    100,
		Interval: 300 * time.Millisecond,
	}
}

func (w *OutboxWorker) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	w.cancel = cancel
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		events, err := w.Repo.FetchPending(ctx, w.Batch)
		if err != nil {
			log.Printf("outbox worker fetch err: %v", err)
			time.Sleep(time.Second)
			continue
		}
		if len(events) == 0 {
			time.Sleep(w.Interval)
			continue
		}
		for _, ev := range events {
			payloadBytes, _ := json.Marshal(ev.Payload)
			values := map[string]interface{}{
				"payload":   string(payloadBytes),
				"trace_id":  ev.TraceID,
				"type":      ev.EventType,
				"outbox_id": ev.ID,
			}
			if err := w.Producer.Add(ctx, "brain.events", values); err != nil {
				log.Printf("outbox publish err id=%d: %v", ev.ID, err)
				_ = w.Repo.IncrementAttempts(ctx, ev.ID)
				continue
			}
			if err := w.Repo.MarkProcessed(ctx, ev.ID); err != nil {
				log.Printf("outbox mark processed err id=%d: %v", ev.ID, err)
			}
		}
	}
}

func (w *OutboxWorker) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
}
