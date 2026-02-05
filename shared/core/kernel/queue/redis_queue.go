package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	rdb       *redis.Client
	streamKey string // default stream key prefix
}

func NewRedisQueue(rdb *redis.Client, streamPrefix string) *RedisQueue {
	if streamPrefix == "" {
		streamPrefix = "app.queue"
	}
	return &RedisQueue{rdb: rdb, streamKey: streamPrefix}
}

type serializedJob struct {
	ID      string            `json:"id"`
	Payload string            `json:"payload"`
	Meta    map[string]string `json:"meta"`
}

func (q *RedisQueue) Push(ctx context.Context, queueName string, job JobPayload) error {
	if queueName == "" {
		queueName = "default"
	}
	sj := serializedJob{ID: job.ID, Payload: string(job.Payload), Meta: job.Meta}
	b, _ := json.Marshal(sj)
	key := q.streamKey + ":" + queueName
	_, err := q.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: key,
		Values: map[string]interface{}{"payload": string(b)},
	}).Result()
	return err
}

func (q *RedisQueue) Later(ctx context.Context, queueName string, job JobPayload, delay time.Duration) error {
	// delayed: push into zset delay:<queue>
	key := fmt.Sprintf("delay:%s:%s", q.streamKey, queueName)
	b, _ := json.Marshal(job)
	score := float64(time.Now().Add(delay).UnixMilli())
	return q.rdb.ZAdd(ctx, key, redis.Z{Score: score, Member: b}).Err()
}

func (q *RedisQueue) PopSerialized(ctx context.Context, queueName string) (id string, payload string, err error) {
	key := q.streamKey + ":" + queueName
	// XREADCOUNT 1 BLOCK 1000
	res, err := q.rdb.XRead(ctx, &redis.XReadArgs{
		Streams: []string{key, "$"},
		Count:   1,
		Block:   1000 * time.Millisecond,
	}).Result()
	if err != nil {
		return "", "", err
	}
	if len(res) == 0 { // no message
		return "", "", redis.Nil
	}
	// assume res[0].Messages[0]
	msg := res[0].Messages[0]
	if p, ok := msg.Values["payload"].(string); ok {
		return msg.ID, p, nil
	}
	return msg.ID, "", errors.New("payload missing")
}

func (q *RedisQueue) Ack(id string) error {
	// for streams we can do XACK if consumer groups used. For simple worker we'll XDel
	// Not implemented here; depends on consumer group usage
	return nil
}

func (q *RedisQueue) StartWorker(ctx context.Context, queueName string, handler func(ctx context.Context, payload string) error) {
	// Basic worker: move delayed items due, then XREAD
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			// move due delayed
			q.moveDelayed(ctx, queueName)
			// read 1 message
			id, payload, err := q.PopSerialized(ctx, queueName)
			if err != nil {
				if err == redis.Nil {
					time.Sleep(200 * time.Millisecond)
					continue
				}
				log.Printf("queue pop err: %v", err)
				time.Sleep(time.Second)
				continue
			}
			if err := handler(ctx, payload); err != nil {
				log.Printf("job handler error: %v", err)
				// leave it (for simplicity). Real impl should move to DLQ
			}
			// delete processed message
			_ = q.rdb.XDel(ctx, q.streamKey+":"+queueName, id)
		}
	}()
}

func (q *RedisQueue) moveDelayed(ctx context.Context, queueName string) {
	key := fmt.Sprintf("delay:%s:%s", q.streamKey, queueName)
	now := float64(time.Now().UnixMilli())
	items, _ := q.rdb.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{Min: "-inf", Max: fmt.Sprintf("%f", now), Offset: 0, Count: 10}).Result()
	if len(items) == 0 {
		return
	}
	for _, it := range items {
		// move to stream
		var job JobPayload
		if err := json.Unmarshal([]byte(it.Member.(string)), &job); err != nil {
			// fallback: remove
			_ = q.rdb.ZRem(ctx, key, it.Member).Err()
			continue
		}
		b, _ := json.Marshal(job)
		_, _ = q.rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: q.streamKey + ":" + queueName,
			Values: map[string]interface{}{"payload": string(b)},
		}).Result()
		_ = q.rdb.ZRem(ctx, key, it.Member).Err()
	}
}
