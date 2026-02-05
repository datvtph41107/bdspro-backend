package outbox

import (
	"context"
	"core/kernel/queue"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type OutboxWorker struct {
	Repo     *Repo
	RDB      *redis.Client
	Producer *queue.StreamProducer
	Batch    int
	Interval time.Duration
	Ctx      context.Context
	Cancel   context.CancelFunc
}

func NewOutboxWorker(repo *Repo, rdb *redis.Client) *OutboxWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &OutboxWorker{
		Repo:     repo,
		RDB:      rdb,
		Producer: queue.NewStreamProducer(rdb),
		Batch:    100,
		Interval: 500 * time.Millisecond,
		Ctx:      ctx,
		Cancel:   cancel,
	}
}

func (w *OutboxWorker) Run(ctx context.Context, c interface{}) {
	for {
		select {
		case <-w.Ctx.Done():
			return
		default:
		}
		events, err := w.Repo.FetchPending(context.Background(), w.Batch)
		if err != nil {
			log.Printf("outbox fetch err: %v", err)
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
				"version":   ev.Version,
				"outbox_id": ev.ID,
			}
			if err := w.Producer.Add(context.Background(), "brain.events", values); err != nil {
				log.Printf("outbox publish err id=%d: %v", ev.ID, err)
				_ = w.Repo.IncrementAttempts(context.Background(), ev.ID)
				continue
			}
			_ = w.Repo.MarkProcessed(context.Background(), ev.ID)
		}
	}
}

func (w *OutboxWorker) Stop() {
	if w.Cancel != nil {
		w.Cancel()
	}
}
