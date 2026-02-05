package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	corectx "core/kernel/ctx"
	bus_type "core/kernel/event/bus/type"
	"core/kernel/queue"
)

// RetryEnvelope is the stable envelope we push to retry queue/stream
type RetryEnvelope struct {
	Event     string                 `json:"event"`
	Payload   json.RawMessage        `json:"payload"`
	RetryLeft int                    `json:"retry_left"`
	TraceID   string                 `json:"trace_id"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
}

func RetryMiddleware(producer *queue.StreamProducer, retryStream string) Middleware {
	if retryStream == "" {
		retryStream = "brain.retry"
	}
	return MiddlewareFunc(func(ctx context.Context, payload interface{}, next func(context.Context) error) error {
		err := next(ctx)
		if err == nil {
			return nil
		}
		// check retries on payload
		retries := 0
		if r, ok := payload.(bus_type.HasRetry); ok {
			retries = r.Retries()
		}
		if retries <= 0 {
			return err
		}
		delay := time.Second
		if d, ok := payload.(bus_type.HasDelay); ok {
			delay = d.Delay()
		}
		// marshall payload into stable JSON
		raw, jerr := json.Marshal(payload)
		if jerr != nil {
			// cannot JSON payload — fail fast
			return err
		}
		env := RetryEnvelope{
			Event:     fmt.Sprintf("%T", payload), // consumer can map type name to handler
			Payload:   raw,
			RetryLeft: retries - 1,
			TraceID:   corectx.Trace(ctx),
		}
		data, _ := json.Marshal(env)
		job := map[string]interface{}{
			"payload": string(data),
			"trace":   env.TraceID,
		}
		if addErr := producer.AddDelayed(ctx, retryStream, job, delay); addErr != nil {
			// swallow add error but still return original handler error
			corectx.Logger(ctx).Error("RETRY_PUSH_ERROR", "failed to push retry", addErr)
		}
		return err
	})
}
