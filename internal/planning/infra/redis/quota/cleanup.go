package redisquota

import (
	"context"
	"time"
	"tqd/internal/usecase/quota"

	"github.com/redis/go-redis/v9"
)

/**
 * FindExpiredReservations lấy các reservation đã quá thời gian giữ quota.
 *
 * Redis key vẫn được giữ thêm một khoảng để cleanup có đủ dữ liệu CancelQuota.
 */
func (s *Store) FindExpiredReservations(
	ctx context.Context,
	before time.Time,
	limit int,
) ([]quota.Reservation, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 100
	}

	keys, err := s.client.ZRangeByScore(
		ctx,
		expiredSetKey,
		&redis.ZRangeBy{
			Min:   "-inf",
			Max:   formatUnix(before.Unix()),
			Count: int64(limit),
		},
	).Result()
	if err != nil {
		return nil, err
	}

	reservations := make([]quota.Reservation, 0, len(keys))
	for _, key := range keys {
		id := reservationID(key)
		if id == "" {
			_ = s.client.ZRem(ctx, expiredSetKey, key).Err()
			continue
		}

		exists, existsErr := s.client.Exists(ctx, key).Result()
		if existsErr != nil {
			return nil, existsErr
		}
		if exists == 0 {
			_ = s.client.ZRem(ctx, expiredSetKey, key).Err()
			continue
		}

		reservation, loadErr := s.loadReservation(ctx, id)
		if loadErr != nil {
			if loadErr == quota.ErrReservationNotFound {
				_ = s.client.ZRem(ctx, expiredSetKey, key).Err()
				continue
			}
			return nil, loadErr
		}
		if reservation.State != quota.StateReserved {
			_ = s.client.ZRem(ctx, expiredSetKey, key).Err()
			continue
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}
