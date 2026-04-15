package quota

import (
	"context"
	"fmt"
	"strconv"
	"time"

	_redis "common/redis"

	"github.com/google/uuid"
)

type RedisQuotaStorage struct {
	redis *_redis.RedisService
}

func NewRedisQuotaStorage(redisService *_redis.RedisService) *RedisQuotaStorage {
	return &RedisQuotaStorage{redis: redisService}
}

func (s *RedisQuotaStorage) usageKey(userID, pool string) string {
	return fmt.Sprintf("quota:usage:%s:%s", userID, pool)
}

func (s *RedisQuotaStorage) limitKey(userID, pool string) string {
	return fmt.Sprintf("quota:limit:%s:%s", userID, pool)
}

func (s *RedisQuotaStorage) reservationKey(reservationID string) string {
	return fmt.Sprintf("quota:reserve:%s", reservationID)
}

func (s *RedisQuotaStorage) GetUsage(ctx context.Context, userID, pool string) (int64, error) {
	key := s.usageKey(userID, pool)
	val, err := s.redis.Get(key)
	if err != nil {
		// Key không tồn tại -> 0
		return 0, nil
	}
	return strconv.ParseInt(val, 10, 64)
}

func (s *RedisQuotaStorage) GetLimit(ctx context.Context, userID, pool string) (int64, error) {
	key := s.limitKey(userID, pool)
	val, err := s.redis.Get(key)
	if err != nil {
		return 0, nil
	}
	return strconv.ParseInt(val, 10, 64)
}

func (s *RedisQuotaStorage) IncrementUsage(ctx context.Context, userID, pool string, amount int64) (int64, error) {
	key := s.usageKey(userID, pool)
	// Dùng IncrementBy (cần thêm method trong RedisService)
	if err := s.redis.IncrementBy(key, amount); err != nil {
		return 0, err
	}
	newValStr, err := s.redis.Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(newValStr, 10, 64)
}

func (s *RedisQuotaStorage) SetUsage(ctx context.Context, userID, pool string, usage int64, ttl time.Duration) error {
	key := s.usageKey(userID, pool)
	return s.redis.SetWithTime(key, usage, ttl)
}

func (s *RedisQuotaStorage) CreateReservation(ctx context.Context, userID, pool string, amount int64, ttl time.Duration) (string, error) {
	id := uuid.New().String()
	key := s.reservationKey(id)
	value := fmt.Sprintf("%s|%s|%d", userID, pool, amount)
	err := s.redis.Client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *RedisQuotaStorage) ConfirmReservation(ctx context.Context, reservationID string) error {
	key := s.reservationKey(reservationID)
	val, err := s.redis.Get(key)
	if err != nil {
		return ErrReservationNotFound
	}
	var userID, pool string
	var amount int64
	fmt.Sscanf(val, "%s|%s|%d", &userID, &pool, &amount)
	_, err = s.IncrementUsage(ctx, userID, pool, amount)
	if err != nil {
		return err
	}
	return s.redis.Delete(key)
}

func (s *RedisQuotaStorage) CancelReservation(ctx context.Context, reservationID string) error {
	key := s.reservationKey(reservationID)
	return s.redis.Delete(key)
}
