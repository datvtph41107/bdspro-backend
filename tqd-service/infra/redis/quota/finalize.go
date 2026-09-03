package redisquota

import (
	"context"
	"fmt"
	"tqd/internal/usecase/quota"

	"github.com/redis/go-redis/v9"
)

var commitScript = redis.NewScript(`
local state = redis.call('HGET', KEYS[1], 'state')
if not state then
    return {'not_found'}
end

local usage_key = redis.call('HGET', KEYS[1], 'usage_key')
local amount = tonumber(redis.call('HGET', KEYS[1], 'amount') or '0')
local limit = tonumber(redis.call('HGET', KEYS[1], 'limit') or '0')

if state == 'canceled' then
    return {'canceled'}
end

if state == 'reserved' then
    local current_reserved = tonumber(redis.call('HGET', usage_key, 'reserved') or '0')
    if current_reserved < amount then
        return {'invalid_state'}
    end
    redis.call('HINCRBY', usage_key, 'reserved', -amount)
    redis.call('HINCRBY', usage_key, 'used', amount)
    redis.call('HSET', KEYS[1], 'state', 'committed')
    redis.call('ZREM', KEYS[2], KEYS[1])
    state = 'committed'
end

local used = tonumber(redis.call('HGET', usage_key, 'used') or '0')
local reserved = tonumber(redis.call('HGET', usage_key, 'reserved') or '0')
local remaining = limit - used - reserved
if remaining < 0 then remaining = 0 end
return {'ok', state, tostring(used), tostring(reserved), tostring(remaining)}
`)

var cancelScript = redis.NewScript(`
local state = redis.call('HGET', KEYS[1], 'state')
if not state then
    return {'not_found'}
end

local usage_key = redis.call('HGET', KEYS[1], 'usage_key')
local amount = tonumber(redis.call('HGET', KEYS[1], 'amount') or '0')
local limit = tonumber(redis.call('HGET', KEYS[1], 'limit') or '0')

if state == 'committed' then
    return {'committed'}
end

if state == 'reserved' then
    local current_reserved = tonumber(redis.call('HGET', usage_key, 'reserved') or '0')
    if current_reserved < amount then
        return {'invalid_state'}
    end
    redis.call('HINCRBY', usage_key, 'reserved', -amount)
    redis.call('HSET', KEYS[1], 'state', 'canceled')
    redis.call('ZREM', KEYS[2], KEYS[1])
    state = 'canceled'
end

local used = tonumber(redis.call('HGET', usage_key, 'used') or '0')
local reserved = tonumber(redis.call('HGET', usage_key, 'reserved') or '0')
local remaining = limit - used - reserved
if remaining < 0 then remaining = 0 end
return {'ok', state, tostring(used), tostring(reserved), tostring(remaining)}
`)

func (s *Store) CommitQuota(ctx context.Context, id string) (quota.Reservation, error) {
	return s.finalize(ctx, id, commitScript, true)
}

func (s *Store) CancelQuota(ctx context.Context, id string) (quota.Reservation, error) {
	return s.finalize(ctx, id, cancelScript, false)
}

func (s *Store) finalize(
	ctx context.Context,
	id string,
	script *redis.Script,
	commit bool,
) (quota.Reservation, error) {
	if err := s.validate(); err != nil {
		return quota.Reservation{}, err
	}

	rKey := reservationKey(id)
	result, err := script.Run(ctx, s.client, []string{rKey, expiredSetKey}).Result()
	if err != nil {
		return quota.Reservation{}, err
	}

	values, ok := result.([]interface{})
	if !ok || len(values) == 0 {
		return quota.Reservation{}, fmt.Errorf("unexpected redis finalize result")
	}

	switch valueString(values[0]) {
	case "not_found":
		return quota.Reservation{}, quota.ErrReservationNotFound
	case "canceled":
		return quota.Reservation{}, quota.ErrReservationCanceled
	case "committed":
		return quota.Reservation{}, quota.ErrReservationCommitted
	case "invalid_state":
		return quota.Reservation{}, fmt.Errorf("runtime quota state is inconsistent")
	case "ok":
	default:
		return quota.Reservation{}, fmt.Errorf("unexpected redis finalize status %q", valueString(values[0]))
	}

	loaded, err := s.loadReservation(ctx, id)
	if err != nil {
		return quota.Reservation{}, err
	}
	if len(values) == 5 {
		loaded.State = quota.State(valueString(values[1]))
		loaded.Used, _ = valueInt64(values[2])
		loaded.Reserved, _ = valueInt64(values[3])
		loaded.Remaining, _ = valueInt64(values[4])
	}

	_ = commit // keeps call site explicit when reading the flow
	return loaded, nil
}
