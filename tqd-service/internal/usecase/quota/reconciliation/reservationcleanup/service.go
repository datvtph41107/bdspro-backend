package reservationcleanup

import (
	"context"
	"errors"
	"time"
	"tqd/internal/usecase/quota"
)

/**
 * QuotaStore là phần runtime state mà cleanup cần dùng.
 */
type QuotaStore interface {
	FindExpiredReservations(ctx context.Context, before time.Time, limit int) ([]quota.Reservation, error)
	CommitQuota(ctx context.Context, reservationID string) (quota.Reservation, error)
	CancelQuota(ctx context.Context, reservationID string) (quota.Reservation, error)
}

/**
 * UsageStore trả lời business đã accepted reservation này hay chưa.
 */
type UsageStore interface {
	HasUsageForReservation(ctx context.Context, reservationID string) (bool, error)
}

/**
 * Result cho biết batch vừa xử lý bao nhiêu reservation.
 */
type Action struct {
	ReservationID     string
	Outcome           string
	DurableUsageFound bool
}

type Result struct {
	Committed int
	Canceled  int
	Skipped   int
	Actions   []Action
}

/**
 * Service xử lý reservation đã quá thời gian giữ quota.
 *
 * Có usage event đúng reservation -> business đã accepted -> commit quota.
 * Không có evidence đúng reservation -> cancel quota, kể cả cùng command từng
 * được accepted ở một kỳ quota khác.
 */
type Service struct {
	quotaStore QuotaStore
	usageStore UsageStore
}

func NewService(quotaStore QuotaStore, usageStore UsageStore) *Service {
	return &Service{
		quotaStore: quotaStore,
		usageStore: usageStore,
	}
}

/**
 * RunOnce xử lý một batch reservation hết hạn.
 */
func (s *Service) RunOnce(
	ctx context.Context,
	now time.Time,
	limit int,
) (Result, error) {
	if s == nil || s.quotaStore == nil || s.usageStore == nil {
		return Result{}, errors.New("quota cleanup is not configured")
	}

	reservations, err := s.quotaStore.FindExpiredReservations(ctx, now, limit)
	if err != nil {
		return Result{}, err
	}

	result := Result{}
	for _, reservation := range reservations {
		hasUsage, err := s.usageStore.HasUsageForReservation(ctx, reservation.ID)
		if err != nil {
			return result, err
		}

		if hasUsage {
			_, err = s.quotaStore.CommitQuota(ctx, reservation.ID)
			switch {
			case err == nil:
				result.Committed++
				result.Actions = append(result.Actions, Action{
					ReservationID:     reservation.ID,
					Outcome:           "committed",
					DurableUsageFound: true,
				})
			case errors.Is(err, quota.ErrReservationNotFound),
				errors.Is(err, quota.ErrReservationCanceled):
				result.Skipped++
			default:
				return result, err
			}
			continue
		}

		_, err = s.quotaStore.CancelQuota(ctx, reservation.ID)
		switch {
		case err == nil:
			result.Canceled++
			result.Actions = append(result.Actions, Action{
				ReservationID:     reservation.ID,
				Outcome:           "canceled",
				DurableUsageFound: false,
			})
		case errors.Is(err, quota.ErrReservationNotFound),
			errors.Is(err, quota.ErrReservationCommitted):
			result.Skipped++
		default:
			return result, err
		}
	}

	return result, nil
}
