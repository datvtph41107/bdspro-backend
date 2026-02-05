package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type StreamProducer struct {
	rdb *redis.Client
}

func NewStreamProducer(rdb *redis.Client) *StreamProducer {
	return &StreamProducer{rdb: rdb}
}

// Add pushes message into Redis Stream
func (p *StreamProducer) Add(
	ctx context.Context,
	stream string,
	values map[string]interface{},
) error {
	if stream == "" {
		return fmt.Errorf("stream name is empty")
	}
	if len(values) == 0 {
		return fmt.Errorf("stream payload is empty")
	}

	args := &redis.XAddArgs{
		Stream: stream,
		Values: values,
	}

	_, err := p.rdb.XAdd(ctx, args).Result()
	if err != nil {
		return fmt.Errorf("redis stream xadd failed: %w", err)
	}

	return nil
}

// AddWithTTL (optional) – useful for transient streams
func (p *StreamProducer) AddWithTTL(
	ctx context.Context,
	stream string,
	values map[string]interface{},
	ttl time.Duration,
) error {

	if err := p.Add(ctx, stream, values); err != nil {
		return err
	}
	if ttl > 0 {
		_ = p.rdb.Expire(ctx, stream, ttl).Err()
	}
	return nil
}

func (p *StreamProducer) AddDelayed(
	ctx context.Context,
	stream string,
	job map[string]interface{},
	delay time.Duration,
) error {
	if delay <= 0 {
		return fmt.Errorf("delay must be > 0")
	}

	b, err := json.Marshal(job)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("delay:%s", stream)
	score := float64(time.Now().Add(delay).UnixMilli())

	return p.rdb.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: b,
	}).Err()
}
