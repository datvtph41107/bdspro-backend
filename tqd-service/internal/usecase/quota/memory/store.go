package memoryquota

import (
	commonmetering "common/metering"
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
	"tqd/internal/access"
	"tqd/internal/usecase/quota"
)

type usageEntry struct {
	used     int64
	reserved int64
}

type savedReservation struct {
	quota.Reservation
	periodStart time.Time
	periodEnd   time.Time
}

/**
 * Store là in-memory implementation dùng cho test và đọc flow.
 *
 * mutex giúp Reserve/Commit/Cancel có cùng semantics atomic như Redis target.
 */
type Store struct {
	mu sync.Mutex

	usage        map[string]*usageEntry
	reservations map[string]*savedReservation
	commands     map[string]string
}

func NewStore() *Store {
	return &Store{
		usage:        make(map[string]*usageEntry),
		reservations: make(map[string]*savedReservation),
		commands:     make(map[string]string),
	}
}

func (s *Store) ReserveQuota(ctx context.Context, input quota.StoreReserveInput) (quota.Reservation, error) {
	if !input.MeterCode.IsValid() {
		return quota.Reservation{}, quota.ErrInvalidInput
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	usageKey := buildUsageKey(input.Subject, input.MeterCode, input.PeriodStart, input.PeriodEnd)
	commandKey := buildCommandKey(usageKey, input.CommandKey)

	if reservationID, ok := s.commands[commandKey]; ok {
		saved, exists := s.reservations[reservationID]
		if !exists {
			delete(s.commands, commandKey)
		} else if saved.State == quota.StateCanceled {
			// Business chưa accepted nên retry cùng command được phép giữ quota lại.
			delete(s.commands, commandKey)
		} else {
			return s.withCurrentUsage(saved), nil
		}
	}

	usage := s.usage[usageKey]
	if usage == nil {
		usage = &usageEntry{}
		s.usage[usageKey] = usage
	}

	if usage.used+usage.reserved+input.Amount > input.Limit {
		return quota.Reservation{}, &quota.ExhaustedError{
			Subject:     input.Subject,
			Operation:   input.Operation,
			MeterCode:   input.MeterCode,
			Limit:       input.Limit,
			Used:        usage.used,
			Reserved:    usage.reserved,
			Remaining:   max64(input.Limit-usage.used-usage.reserved, 0),
			PeriodStart: input.PeriodStart,
			PeriodEnd:   input.PeriodEnd,
		}
	}

	usage.reserved += input.Amount
	reservation := quota.Reservation{
		Required:       true,
		ID:             input.ReservationID,
		Subject:        input.Subject,
		Operation:      input.Operation,
		MeterCode:      input.MeterCode,
		OperationID:    input.OperationID,
		IdempotencyKey: input.IdempotencyKey,
		CommandKey:     input.CommandKey,
		Amount:         input.Amount,
		Limit:          input.Limit,
		PeriodStart:    input.PeriodStart,
		PeriodEnd:      input.PeriodEnd,
		Used:           usage.used,
		Reserved:       usage.reserved,
		Remaining:      max64(input.Limit-usage.used-usage.reserved, 0),
		State:          quota.StateReserved,
		ExpiresAt:      input.ExpiresAt,
	}

	s.reservations[input.ReservationID] = &savedReservation{
		Reservation: reservation,
		periodStart: input.PeriodStart,
		periodEnd:   input.PeriodEnd,
	}
	s.commands[commandKey] = input.ReservationID

	return reservation, nil
}

func (s *Store) CommitQuota(ctx context.Context, reservationID string) (quota.Reservation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	saved, ok := s.reservations[reservationID]
	if !ok {
		return quota.Reservation{}, quota.ErrReservationNotFound
	}
	if saved.State == quota.StateCommitted {
		return s.withCurrentUsage(saved), nil
	}
	if saved.State == quota.StateCanceled {
		return quota.Reservation{}, quota.ErrReservationCanceled
	}

	usageKey := buildUsageKey(saved.Subject, saved.MeterCode, saved.periodStart, saved.periodEnd)
	usage := s.usage[usageKey]
	if usage == nil || usage.reserved < saved.Amount {
		return quota.Reservation{}, fmt.Errorf("runtime quota state is inconsistent")
	}

	usage.reserved -= saved.Amount
	usage.used += saved.Amount
	saved.State = quota.StateCommitted

	return s.withCurrentUsage(saved), nil
}

func (s *Store) CancelQuota(ctx context.Context, reservationID string) (quota.Reservation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	saved, ok := s.reservations[reservationID]
	if !ok {
		return quota.Reservation{}, quota.ErrReservationNotFound
	}
	if saved.State == quota.StateCanceled {
		return s.withCurrentUsage(saved), nil
	}
	if saved.State == quota.StateCommitted {
		return quota.Reservation{}, quota.ErrReservationCommitted
	}

	usageKey := buildUsageKey(saved.Subject, saved.MeterCode, saved.periodStart, saved.periodEnd)
	usage := s.usage[usageKey]
	if usage == nil || usage.reserved < saved.Amount {
		return quota.Reservation{}, fmt.Errorf("runtime quota state is inconsistent")
	}

	usage.reserved -= saved.Amount
	saved.State = quota.StateCanceled

	return s.withCurrentUsage(saved), nil
}

func (s *Store) GetUsage(
	ctx context.Context,
	subject access.Subject,
	meterCode commonmetering.Code,
	periodStart time.Time,
	periodEnd time.Time,
) (quota.Usage, error) {
	if !meterCode.IsValid() {
		return quota.Usage{}, quota.ErrInvalidInput
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	usage := s.usage[buildUsageKey(subject, meterCode, periodStart, periodEnd)]
	if usage == nil {
		return quota.Usage{}, nil
	}
	return quota.Usage{Used: usage.used, Reserved: usage.reserved}, nil
}

func (s *Store) SetUsed(
	ctx context.Context,
	subject access.Subject,
	meterCode commonmetering.Code,
	periodStart time.Time,
	periodEnd time.Time,
	used int64,
) error {
	if used < 0 {
		return quota.ErrInvalidInput
	}
	if !meterCode.IsValid() {
		return quota.ErrInvalidInput
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := buildUsageKey(subject, meterCode, periodStart, periodEnd)
	entry := s.usage[key]
	if entry == nil {
		entry = &usageEntry{}
		s.usage[key] = entry
	}
	entry.used = used
	return nil
}

func (s *Store) CompareAndSetUsed(
	ctx context.Context,
	subject access.Subject,
	meterCode commonmetering.Code,
	periodStart time.Time,
	periodEnd time.Time,
	expected int64,
	used int64,
) (bool, error) {
	if expected < 0 || used < 0 {
		return false, quota.ErrInvalidInput
	}
	if !meterCode.IsValid() {
		return false, quota.ErrInvalidInput
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := buildUsageKey(subject, meterCode, periodStart, periodEnd)
	entry := s.usage[key]
	current := int64(0)
	if entry != nil {
		current = entry.used
	}
	if current != expected {
		return false, nil
	}
	if entry == nil {
		entry = &usageEntry{}
		s.usage[key] = entry
	}
	entry.used = used
	return true, nil
}

func (s *Store) FindExpiredReservations(ctx context.Context, before time.Time, limit int) ([]quota.Reservation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 100
	}

	ids := make([]string, 0)
	for id, reservation := range s.reservations {
		if reservation.State != quota.StateReserved || reservation.ExpiresAt.After(before) {
			continue
		}
		ids = append(ids, id)
	}
	sort.Strings(ids)
	if len(ids) > limit {
		ids = ids[:limit]
	}

	result := make([]quota.Reservation, 0, len(ids))
	for _, id := range ids {
		result = append(result, s.withCurrentUsage(s.reservations[id]))
	}
	return result, nil
}

func (s *Store) withCurrentUsage(saved *savedReservation) quota.Reservation {
	result := saved.Reservation
	usageKey := buildUsageKey(saved.Subject, saved.MeterCode, saved.periodStart, saved.periodEnd)
	usage := s.usage[usageKey]
	if usage != nil {
		result.Used = usage.used
		result.Reserved = usage.reserved
		result.Remaining = max64(result.Limit-usage.used-usage.reserved, 0)
	}
	return result
}

func buildUsageKey(
	subject access.Subject,
	meterCode commonmetering.Code,
	periodStart time.Time,
	periodEnd time.Time,
) string {
	return fmt.Sprintf(
		"%s:%s:%s:%d:%d",
		subject.Type,
		subject.ID,
		meterCode,
		periodStart.Unix(),
		periodEnd.Unix(),
	)
}

func buildCommandKey(usageKey, commandKey string) string {
	return usageKey + ":" + commandKey
}

func max64(left, right int64) int64 {
	if left > right {
		return left
	}
	return right
}
