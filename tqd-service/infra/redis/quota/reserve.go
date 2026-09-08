package redisquota

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"tqd/internal/usecase/quota"

	"github.com/redis/go-redis/v9"
)

var reserveScript = redis.NewScript(`
local existing_key = redis.call('GET', KEYS[2])
if existing_key then
    local state = redis.call('HGET', existing_key, 'state')
    if not state or state == 'canceled' then
        -- Reservation cũ không còn giữ business state; retry được tạo reservation mới.
        redis.call('DEL', KEYS[2])
    else
        local id = redis.call('HGET', existing_key, 'id') or ''
        local used = tonumber(redis.call('HGET', KEYS[1], 'used') or '0')
        local reserved = tonumber(redis.call('HGET', KEYS[1], 'reserved') or '0')
        local limit = tonumber(redis.call('HGET', existing_key, 'limit') or ARGV[9])
        local remaining = limit - used - reserved
        if remaining < 0 then remaining = 0 end

        return {'existing', id, state, tostring(used), tostring(reserved), tostring(remaining)}
    end
end

local used = tonumber(redis.call('HGET', KEYS[1], 'used') or '0')
local reserved = tonumber(redis.call('HGET', KEYS[1], 'reserved') or '0')
local amount = tonumber(ARGV[8])
local limit = tonumber(ARGV[9])

if used + reserved + amount > limit then
    local remaining = limit - used - reserved
    if remaining < 0 then remaining = 0 end
    return {'exceeded', '', '', tostring(used), tostring(reserved), tostring(remaining)}
end

local new_reserved = redis.call('HINCRBY', KEYS[1], 'reserved', amount)
redis.call('HSET', KEYS[1], 'used', used)

redis.call('HSET', KEYS[3],
    'id', ARGV[1],
    'subject_type', ARGV[2],
    'subject_id', ARGV[3],
    'operation', ARGV[4],
	'meter_code', ARGV[15],
    'operation_id', ARGV[5],
    'idempotency_key', ARGV[6],
    'command_key', ARGV[7],
    'amount', ARGV[8],
    'limit', ARGV[9],
    'state', 'reserved',
    'usage_key', KEYS[1],
    'expires_at', ARGV[10],
    'period_start', ARGV[13],
    'period_end', ARGV[14]
)

redis.call('SET', KEYS[2], KEYS[3])
if tonumber(ARGV[11]) > 0 then
    redis.call('PEXPIRE', KEYS[2], ARGV[11])
    redis.call('PEXPIRE', KEYS[3], ARGV[11])
end
if tonumber(ARGV[12]) > 0 then
    redis.call('PEXPIRE', KEYS[1], ARGV[12])
end
redis.call('ZADD', KEYS[4], ARGV[10], KEYS[3])

local remaining = limit - used - new_reserved
if remaining < 0 then remaining = 0 end
return {'reserved', ARGV[1], 'reserved', tostring(used), tostring(new_reserved), tostring(remaining)}
`)

func (s *Store) ReserveQuota(ctx context.Context, input quota.StoreReserveInput) (quota.Reservation, error) {
	if err := s.validate(); err != nil {
		return quota.Reservation{}, err
	}
	if !input.MeterCode.IsValid() {
		return quota.Reservation{}, quota.ErrInvalidInput
	}

	now := time.Now().UTC()
	keep := keepTTL(input.PeriodEnd, now)
	usageTTL := keep

	uKey := usageKey(input.Subject, input.MeterCode, input.PeriodStart, input.PeriodEnd)
	cKey := commandKey(uKey, input.CommandKey)
	rKey := reservationKey(input.ReservationID)

	result, err := reserveScript.Run(
		ctx,
		s.client,
		[]string{uKey, cKey, rKey, expiredSetKey},
		input.ReservationID,
		string(input.Subject.Type),
		input.Subject.ID,
		string(input.Operation),
		input.OperationID,
		input.IdempotencyKey,
		input.CommandKey,
		strconv.FormatInt(input.Amount, 10),
		strconv.FormatInt(input.Limit, 10),
		strconv.FormatInt(input.ExpiresAt.Unix(), 10),
		strconv.FormatInt(keep.Milliseconds(), 10),
		strconv.FormatInt(usageTTL.Milliseconds(), 10),
		strconv.FormatInt(input.PeriodStart.Unix(), 10),
		strconv.FormatInt(input.PeriodEnd.Unix(), 10),
		string(input.MeterCode),
	).Result()
	if err != nil {
		return quota.Reservation{}, err
	}

	values, ok := result.([]interface{})
	if !ok || len(values) < 1 {
		return quota.Reservation{}, fmt.Errorf("unexpected redis reserve result")
	}

	status := valueString(values[0])
	switch status {
	case "exceeded":
		if len(values) != 6 {
			return quota.Reservation{}, fmt.Errorf("unexpected redis reserve result length")
		}
	case "reserved", "existing":
		if len(values) != 6 {
			return quota.Reservation{}, fmt.Errorf("unexpected redis reserve result length")
		}
	default:
		return quota.Reservation{}, fmt.Errorf("unexpected redis reserve status %q", status)
	}

	used, err := valueInt64(values[3])
	if err != nil {
		return quota.Reservation{}, err
	}
	reserved, err := valueInt64(values[4])
	if err != nil {
		return quota.Reservation{}, err
	}
	remaining, err := valueInt64(values[5])
	if err != nil {
		return quota.Reservation{}, err
	}
	if status == "exceeded" {
		return quota.Reservation{}, &quota.ExhaustedError{
			Subject:     input.Subject,
			Operation:   input.Operation,
			MeterCode:   input.MeterCode,
			Limit:       input.Limit,
			Used:        used,
			Reserved:    reserved,
			Remaining:   remaining,
			PeriodStart: input.PeriodStart,
			PeriodEnd:   input.PeriodEnd,
		}
	}

	reservation := quota.Reservation{
		Required:       true,
		ID:             valueString(values[1]),
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
		Used:           used,
		Reserved:       reserved,
		Remaining:      remaining,
		State:          quota.State(valueString(values[2])),
		ExpiresAt:      input.ExpiresAt,
	}

	if status == "existing" {
		loaded, loadErr := s.loadReservation(ctx, reservation.ID)
		if loadErr != nil {
			return quota.Reservation{}, loadErr
		}
		loaded.Used = used
		loaded.Reserved = reserved
		loaded.Remaining = remaining
		return loaded, nil
	}

	return reservation, nil
}
