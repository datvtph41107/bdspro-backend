package bus

import (
	"context"
	"core/kernel/container"
	core_ctx "core/kernel/ctx"
	bus_type "core/kernel/event/bus/type"
	"core/kernel/event/outbox"
	"core/kernel/middleware"
	"encoding/json"
	"fmt"
	"reflect"
)

// Bus main structure
type Bus struct {
	container  *container.Container
	registry   *Registry
	outbox     *outbox.Repo
	middleware []middleware.Middleware
}

func NewBus(c *container.Container, repo *outbox.Repo) *Bus {
	b := &Bus{
		container:  c,
		registry:   NewRegistry(),
		outbox:     repo,
		middleware: []middleware.Middleware{},
	}
	// install default middleware order
	b.Use(middleware.TransactionMiddleware(c))
	b.Use(middleware.OutboxMiddleware(repo))
	b.Use(middleware.TraceMiddleware())
	b.Use(middleware.LoggingMiddleware())
	return b
}

func (b *Bus) Use(m middleware.Middleware) {
	b.middleware = append(b.middleware, m)
}

// RegisterHandler convenience
func (b *Bus) RegisterHandler(payloadPrototype interface{}, handler interface{}) error {
	return b.registry.Register(payloadPrototype, handler)
}

func (b *Bus) DispatchIfNeed(ctx context.Context, payload interface{}) error {
	return nil
}

// Dispatch entrypoint
func (b *Bus) Dispatch(ctx context.Context, payload interface{}) error {
	meta := b.introspect(payload)
	// ensure trace id exists
	trace := core_ctx.Trace(ctx)
	ctx = core_ctx.WithTrace(ctx, trace)

	// build pipeline
	var exec func(ctx context.Context) error
	exec = func(ctx context.Context) error {
		// decide: if async -> persist outbox inside tx (middleware does that)
		// else invoke handler(s)
		if meta.Async {
			// OutboxMiddleware will persist under tx; here do nothing
			return nil
		}
		return b.invokeHandlers(ctx, payload)
	}

	// apply middleware chain in reverse order
	for i := len(b.middleware) - 1; i >= 0; i-- {
		next := exec
		m := b.middleware[i]
		exec = func(m middleware.Middleware, nextFn func(context.Context) error) func(ctx context.Context) error {
			return func(ctx context.Context) error {
				return m.Handle(ctx, payload, nextFn)
			}
		}(m, next)
	}

	return exec(ctx)
}

// Introspect reads interfaces on payload to produce Meta
func (b *Bus) introspect(payload interface{}) Meta {
	meta := Meta{}
	t := reflect.TypeOf(payload)
	meta.TypeName = t.String()

	// Async if implements ShouldQueue
	if _, ok := payload.(bus_type.ShouldQueue); ok {
		meta.Async = true
		if q, ok := payload.(bus_type.ShouldQueue); ok {
			meta.QueueName = q.QueueName()
		}
	}
	if r, ok := payload.(bus_type.HasRetry); ok {
		meta.Retries = r.Retries()
	}
	if d, ok := payload.(bus_type.HasDelay); ok {
		meta.DelayMs = int64(d.Delay().Milliseconds())
	}
	// fallback: Events are async by default? We'll keep explicit only.
	return meta
}

// invokeHandlers: call all handlers registered for payload type
func (b *Bus) invokeHandlers(ctx context.Context, payload interface{}) error {
	handlers := b.registry.GetHandlersFor(payload)
	if len(handlers) == 0 {
		return nil
	}

	pv := reflect.ValueOf(payload)

	for _, he := range handlers {
		fn := he.fn
		fnType := fn.Type()
		argType := fnType.In(1)

		var argVal reflect.Value

		// ---- Case 1: exact match ----
		if pv.Type().AssignableTo(argType) {
			argVal = pv

			// ---- Case 2: handler wants pointer, payload is value ----
		} else if pv.Kind() != reflect.Ptr && pv.CanAddr() && pv.Addr().Type().AssignableTo(argType) {
			argVal = pv.Addr()

			// ---- Case 3: handler wants value, payload is pointer ----
		} else if pv.Kind() == reflect.Ptr && pv.Elem().Type().AssignableTo(argType) {
			argVal = pv.Elem()

			// ---- Case 4: fallback – decode by JSON ----
		} else {
			argPtr := reflect.New(argType)
			bb, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshal payload failed: %w", err)
			}
			if err := json.Unmarshal(bb, argPtr.Interface()); err != nil {
				return fmt.Errorf("cannot bind payload to handler arg %s: %w", argType, err)
			}
			argVal = argPtr.Elem()
		}

		// ---- Call handler ----
		out := fn.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			argVal,
		})

		if err, ok := out[0].Interface().(error); ok && err != nil {
			return err
		}
	}

	return nil
}
