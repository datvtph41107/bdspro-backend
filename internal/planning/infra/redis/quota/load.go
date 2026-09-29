package redisquota

import (
	commonmetering "common/metering"
	"common/operation"
	"context"
	"fmt"
	"strconv"
	"time"
	"tqd/internal/access"
	"tqd/internal/usecase/quota"
)

func (s *Store) loadReservation(ctx context.Context, id string) (quota.Reservation, error) {
	values, err := s.client.HGetAll(ctx, reservationKey(id)).Result()
	if err != nil {
		return quota.Reservation{}, err
	}
	if len(values) == 0 {
		return quota.Reservation{}, quota.ErrReservationNotFound
	}

	amount, err := strconv.ParseInt(values["amount"], 10, 64)
	if err != nil {
		return quota.Reservation{}, fmt.Errorf("invalid reservation amount: %w", err)
	}
	limit, err := strconv.ParseInt(values["limit"], 10, 64)
	if err != nil {
		return quota.Reservation{}, fmt.Errorf("invalid reservation limit: %w", err)
	}
	expiresUnix, err := strconv.ParseInt(values["expires_at"], 10, 64)
	if err != nil {
		return quota.Reservation{}, fmt.Errorf("invalid reservation expires_at: %w", err)
	}
	periodStartUnix, err := strconv.ParseInt(values["period_start"], 10, 64)
	if err != nil {
		return quota.Reservation{}, fmt.Errorf("invalid reservation period_start: %w", err)
	}
	periodEndUnix, err := strconv.ParseInt(values["period_end"], 10, 64)
	if err != nil {
		return quota.Reservation{}, fmt.Errorf("invalid reservation period_end: %w", err)
	}

	return quota.Reservation{
		Required: true,
		ID:       values["id"],
		Subject: access.Subject{
			Type: access.SubjectType(values["subject_type"]),
			ID:   values["subject_id"],
		},
		Operation:      operation.Code(values["operation"]),
		MeterCode:      commonmetering.Code(values["meter_code"]),
		OperationID:    values["operation_id"],
		IdempotencyKey: values["idempotency_key"],
		CommandKey:     values["command_key"],
		Amount:         amount,
		Limit:          limit,
		PeriodStart:    time.Unix(periodStartUnix, 0).UTC(),
		PeriodEnd:      time.Unix(periodEndUnix, 0).UTC(),
		State:          quota.State(values["state"]),
		ExpiresAt:      time.Unix(expiresUnix, 0).UTC(),
	}, nil
}

func valueString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func valueInt64(value interface{}) (int64, error) {
	return strconv.ParseInt(valueString(value), 10, 64)
}

func formatUnix(value int64) string {
	return strconv.FormatInt(value, 10)
}
