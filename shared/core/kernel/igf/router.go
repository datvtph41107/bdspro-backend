package igf

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"core/kernel/event"

	"github.com/redis/go-redis/v9"
)

type Router struct {
	rdb      *redis.Client
	registry *event.SchemaRegistry
	routes   map[string][]string
	group    string
	consumer string
}

func NewRouter(rdb *redis.Client, registry *event.SchemaRegistry, routes map[string][]string) *Router {
	return &Router{rdb: rdb, registry: registry, routes: routes, group: "igf-router", consumer: "igf-1"}
}

func (r *Router) ensureGroup(ctx context.Context, stream string) {
	_ = r.rdb.XGroupCreateMkStream(ctx, stream, r.group, "0").Err()
}

func (r *Router) Serve(ctx context.Context) {
	stream := "brain.events"
	r.ensureGroup(ctx, stream)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		args := &redis.XReadGroupArgs{
			Group:    r.group,
			Consumer: r.consumer,
			Streams:  []string{stream, ">"},
			Count:    10,
			Block:    5 * time.Second,
		}
		res, err := r.rdb.XReadGroup(ctx, args).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			log.Printf("igf xread err: %v", err)
			time.Sleep(time.Second)
			continue
		}
		for _, st := range res {
			for _, msg := range st.Messages {
				raw, _ := msg.Values["payload"].(string)
				trace, _ := msg.Values["trace_id"].(string)
				typ, _ := msg.Values["type"].(string)
				if typ == "" {
					// try parse from payload
					var tmp map[string]interface{}
					_ = json.Unmarshal([]byte(raw), &tmp)
					if t, ok := tmp["type"].(string); ok {
						typ = t
					}
				}
				if typ == "" {
					// unknown type -> ack and continue
					_ = r.rdb.XAck(ctx, stream, r.group, msg.ID).Err()
					continue
				}
				// validate (if registry has schema)
				// if registry exists we attempt validation
				if r.registry != nil {
					// attempt parse event type->version if possible: default ver=1
					_ = r.registry.Validate(typ, 1, []byte(raw))
				}
				targets := r.routes[typ]
				if len(targets) == 0 {
					// default to ws.fanout
					targets = []string{"ws.fanout"}
				}
				for _, t := range targets {
					_ = r.rdb.XAdd(ctx, &redis.XAddArgs{
						Stream: t,
						Values: map[string]interface{}{
							"payload":  raw,
							"trace_id": trace,
							"type":     typ,
						},
					}).Err()
				}
				// ack after successful fanout
				_ = r.rdb.XAck(ctx, stream, r.group, msg.ID).Err()
			}
		}
	}
}

func (r *Router) ServeBackground(ctx context.Context) {
	go r.Serve(ctx)
}
