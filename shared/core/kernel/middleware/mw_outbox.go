package middleware

import (
	"context"
	"encoding/json"
	"fmt"

	corectx "core/kernel/ctx"
	"core/kernel/event/outbox"
)

// DomainEvent: events intended to be persisted in outbox must implement this
type DomainEvent interface {
	EventName() string
	Version() int
}

// OutboxMiddleware persists DomainEvent into outbox using tx from ctx.
func OutboxMiddleware(repo *outbox.Repo) Middleware {
	return MiddlewareFunc(func(ctx context.Context, payload interface{}, next func(context.Context) error) error {
		// only handle DomainEvent (explicit)
		ev, ok := payload.(DomainEvent)
		if !ok {
			return next(ctx)
		}
		tx := corectx.Tx(ctx)
		if tx == nil {
			return fmt.Errorf("outbox requires DB transaction in context")
		}
		// serialize payload to JSON map
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		var m map[string]interface{}
		if err := json.Unmarshal(b, &m); err != nil {
			return err
		}
		out := &outbox.OutboxEvent{
			EventType: ev.EventName(),
			Version:   ev.Version(),
			Payload:   outbox.JSONMap(m),
			TraceID:   corectx.Trace(ctx),
		}
		if err := repo.SaveTx(ctx, tx, out); err != nil {
			return err
		}
		// Do not call handler now — outbox worker will publish
		return nil
	})
}
