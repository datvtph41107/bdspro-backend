package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Broadcaster interface {
	Serve(ctx context.Context) error
	JoinChannel(ctx context.Context, channel string, connID string) error
	LeaveChannel(ctx context.Context, channel string, connID string) error
	Broadcast(ctx context.Context, channel string, event string, payload interface{}) error
	RegisterLocalConnection(connID string, send func(msg []byte))
	UnregisterLocalConnection(connID string)
}

type RedisBroadcaster struct {
	rdb    *redis.Client
	inst   string   // instance id
	locals sync.Map // connID -> send func
}

func NewRedisBroadcaster(rdb *redis.Client, instanceID string) *RedisBroadcaster {
	return &RedisBroadcaster{rdb: rdb, inst: instanceID}
}

func (b *RedisBroadcaster) Serve(ctx context.Context) error {
	// subscribe to global pubsub for ws.fanout and pattern ws.channel.*
	ps := b.rdb.PSubscribe(ctx, "ws.channel.*", "ws.fanout")
	ch := ps.Channel()
	for {
		select {
		case <-ctx.Done():
			_ = ps.Close()
			return nil
		case msg := <-ch:
			// msg.Channel e.g. ws.channel.room.123 or ws.fanout
			var m struct {
				Event   string            `json:"event"`
				Payload interface{}       `json:"payload"`
				Meta    map[string]string `json:"meta,omitempty"`
			}
			if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
				log.Printf("ws: invalid payload %v", err)
				continue
			}
			// deliver to local sockets (simple approach: broadcast to all locals)
			b.locals.Range(func(k, v interface{}) bool {
				if send, ok := v.(func([]byte)); ok {
					data, _ := json.Marshal(m)
					go send(data)
				}
				return true
			})
		}
	}
}

func (b *RedisBroadcaster) JoinChannel(ctx context.Context, channel string, connID string) error {
	key := "channel:" + channel
	return b.rdb.SAdd(ctx, key, connID).Err()
}

func (b *RedisBroadcaster) LeaveChannel(ctx context.Context, channel string, connID string) error {
	key := "channel:" + channel
	return b.rdb.SRem(ctx, key, connID).Err()
}

func (b *RedisBroadcaster) Broadcast(ctx context.Context, channel string, event string, payload interface{}) error {
	m := map[string]interface{}{"event": event, "payload": payload}
	bs, _ := json.Marshal(m)
	// publish to pubsub
	return b.rdb.Publish(ctx, "ws.channel."+channel, bs).Err()
}

func (b *RedisBroadcaster) RegisterLocalConnection(connID string, send func(msg []byte)) {
	b.locals.Store(connID, send)
}

func (b *RedisBroadcaster) UnregisterLocalConnection(connID string) {
	b.locals.Delete(connID)
}
