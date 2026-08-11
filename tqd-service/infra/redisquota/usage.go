package redisquota

import (
	commonaccess "common/access"
	"common/businessoperation"
	"context"
	"fmt"
	"strconv"
	"time"
	"tqd/internal/quota"

	"github.com/redis/go-redis/v9"
)

func (s *Store) GetUsage(
	ctx context.Context,
	subject commonaccess.Subject,
	operation businessoperation.Code,
	periodStart time.Time,
	periodEnd time.Time,
) (quota.Usage, error) {
	if err := s.validate(); err != nil {
		return quota.Usage{}, err
	}

	values, err := s.client.HMGet(
		ctx,
		usageKey(subject, operation, periodStart, periodEnd),
		"used",
		"reserved",
	).Result()
	if err != nil {
		return quota.Usage{}, err
	}

	used, err := redisInt64(values, 0)
	if err != nil {
		return quota.Usage{}, err
	}
	reserved, err := redisInt64(values, 1)
	if err != nil {
		return quota.Usage{}, err
	}

	return quota.Usage{Used: used, Reserved: reserved}, nil
}

func (s *Store) SetUsed(
	ctx context.Context,
	subject commonaccess.Subject,
	operation businessoperation.Code,
	periodStart time.Time,
	periodEnd time.Time,
	used int64,
) error {
	if err := s.validate(); err != nil {
		return err
	}
	if used < 0 {
		return quota.ErrInvalidInput
	}

	key := usageKey(subject, operation, periodStart, periodEnd)
	if err := s.client.HSet(ctx, key, "used", used).Err(); err != nil {
		return err
	}

	keep := keepTTL(periodEnd, time.Now().UTC())
	if keep > 0 {
		return s.client.Expire(ctx, key, keep).Err()
	}
	return nil
}

var compareAndSetUsedScript = redis.NewScript(`
local current = tonumber(redis.call("HGET", KEYS[1], "used") or "0")
local expected = tonumber(ARGV[1])
if current ~= expected then
  return 0
end
redis.call("HSET", KEYS[1], "used", ARGV[2])
local keep_milliseconds = tonumber(ARGV[3])
if keep_milliseconds > 0 then
  redis.call("PEXPIRE", KEYS[1], keep_milliseconds)
end
return 1
`)

func (s *Store) CompareAndSetUsed(
	ctx context.Context,
	subject commonaccess.Subject,
	operation businessoperation.Code,
	periodStart time.Time,
	periodEnd time.Time,
	expected int64,
	used int64,
) (bool, error) {
	if err := s.validate(); err != nil {
		return false, err
	}
	if expected < 0 || used < 0 {
		return false, quota.ErrInvalidInput
	}

	keepMilliseconds := keepTTL(periodEnd, time.Now().UTC()).Milliseconds()
	result, err := compareAndSetUsedScript.Run(
		ctx,
		s.client,
		[]string{usageKey(subject, operation, periodStart, periodEnd)},
		expected,
		used,
		keepMilliseconds,
	).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func redisInt64(values []interface{}, index int) (int64, error) {
	if index >= len(values) || values[index] == nil {
		return 0, nil
	}
	value := fmt.Sprint(values[index])
	return strconv.ParseInt(value, 10, 64)
}
