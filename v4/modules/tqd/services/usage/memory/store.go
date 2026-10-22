package memoryusage

import (
	commonmetering "common/metering"
	"context"
	"errors"
	"mdn4/modules/tqd/access"
	"mdn4/modules/tqd/services/usage"
	"sync"
	"time"
)

/**
 * Store là in-memory usage store dùng cho test và đọc flow.
 */
type Store struct {
	mu       sync.Mutex
	events   map[string]usage.Event
	commands map[readCommandKey]readCommandState
}

func NewStore() *Store {
	return &Store{events: make(map[string]usage.Event), commands: make(map[readCommandKey]readCommandState)}
}

func (s *Store) SaveUsageEvent(
	ctx context.Context,
	event usage.Event,
) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if !event.IsValid() {
		// Giữ cùng ngữ nghĩa với PostgreSQL: dữ liệu sai là lỗi, không phải
		// replay thành công. Test double không được che lỗi production.
		return false, errors.New("usage event is invalid")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[event.UsageKey]; exists {
		return false, nil
	}
	s.events[event.UsageKey] = event
	return true, nil
}

func (s *Store) HasUsageForReservation(
	ctx context.Context,
	reservationID string,
) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if reservationID == "" {
		return false, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, event := range s.events {
		if event.ReservationID == reservationID {
			return true, nil
		}
	}
	return false, nil
}

func (s *Store) SumUsage(
	ctx context.Context,
	subject access.Subject,
	meterCode commonmetering.Code,
	periodStart time.Time,
	periodEnd time.Time,
) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var total int64
	for _, event := range s.events {
		if event.Subject == subject &&
			event.MeterCode == meterCode &&
			event.PeriodStart.Equal(periodStart) &&
			event.PeriodEnd.Equal(periodEnd) {
			total += event.Amount
		}
	}
	return total, nil
}
