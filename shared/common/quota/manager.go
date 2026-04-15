package quota

import (
	"context"
	"time"
)

type quotaManager struct {
	storage QuotaStorage
	policy  PolicyEngine
}

func NewQuotaManager(storage QuotaStorage, policy PolicyEngine) QuotaManager {
	return &quotaManager{
		storage: storage,
		policy:  policy,
	}
}

func (m *quotaManager) CheckOnly(ctx context.Context, userID, service, action string) (*QuotaInfo, error) {
	start := time.Now()

	pool, _, err := m.policy.Resolve(service, action)
	if err != nil {
		return nil, err
	}
	limit, err := m.policy.GetUserLimit(ctx, userID, pool)
	if err != nil {
		return nil, err
	}
	used, err := m.storage.GetUsage(ctx, userID, pool)
	if err != nil {
		return nil, err
	}
	remaining := limit - used
	latencyMs := time.Since(start).Milliseconds()

	// Ghi log check quota
	LogCheck(ctx, userID, service, action, pool, used, limit, remaining, endOfMonth(), latencyMs)

	// Cảnh báo nếu sắp hết (<10%)
	if remaining > 0 && remaining <= limit/10 {
		LogWarn(ctx, userID, service, action, pool, remaining, limit,
			"quota almost exhausted")
	}

	return &QuotaInfo{
		UserID:    userID,
		Pool:      pool,
		Used:      used,
		Limit:     limit,
		Remaining: remaining,
		ResetAt:   endOfMonth(),
	}, nil
}

func (m *quotaManager) Consume(ctx context.Context, userID, service, action string) (*ConsumeResult, error) {
	start := time.Now()

	pool, cost, err := m.policy.Resolve(service, action)
	if err != nil {
		return nil, err
	}
	limit, err := m.policy.GetUserLimit(ctx, userID, pool)
	if err != nil {
		return nil, err
	}
	usedBefore, _ := m.storage.GetUsage(ctx, userID, pool)

	newUsed, err := m.storage.IncrementUsage(ctx, userID, pool, cost)
	if err != nil {
		// Ghi log lỗi consume
		LogError(ctx, userID, service, action, err, map[string]interface{}{
			"pool": pool,
			"cost": cost,
		})
		return &ConsumeResult{Success: false, Error: err}, err
	}

	success := newUsed <= limit
	remaining := limit - newUsed
	if remaining < 0 {
		remaining = 0
	}
	latencyMs := time.Since(start).Milliseconds()

	// Ghi log consume
	LogConsume(ctx, userID, service, action, pool, cost, usedBefore, newUsed, limit, remaining, endOfMonth(), latencyMs)

	if !success {
		return &ConsumeResult{
			Success: false,
			QuotaInfo: &QuotaInfo{
				UserID:    userID,
				Pool:      pool,
				Used:      newUsed,
				Limit:     limit,
				Remaining: remaining,
				ResetAt:   endOfMonth(),
			},
			Error: ErrQuotaExceeded,
		}, nil
	}

	return &ConsumeResult{
		Success: true,
		QuotaInfo: &QuotaInfo{
			UserID:    userID,
			Pool:      pool,
			Used:      newUsed,
			Limit:     limit,
			Remaining: remaining,
			ResetAt:   endOfMonth(),
		},
	}, nil
}

func (m *quotaManager) Reserve(ctx context.Context, userID, service, action string, ttl time.Duration) (*Reservation, error) {
	pool, cost, err := m.policy.Resolve(service, action)
	if err != nil {
		return nil, err
	}
	limit, err := m.policy.GetUserLimit(ctx, userID, pool)
	if err != nil {
		return nil, err
	}
	used, err := m.storage.GetUsage(ctx, userID, pool)
	if err != nil {
		return nil, err
	}
	if used+cost > limit {
		return nil, ErrQuotaExceeded
	}
	id, err := m.storage.CreateReservation(ctx, userID, pool, cost, ttl)
	if err != nil {
		return nil, err
	}
	return &Reservation{
		ID:        id,
		UserID:    userID,
		Pool:      pool,
		Amount:    cost,
		ExpiresAt: time.Now().Add(ttl),
	}, nil
}

func (m *quotaManager) Confirm(ctx context.Context, reservationID string) error {
	return m.storage.ConfirmReservation(ctx, reservationID)
}

func (m *quotaManager) Cancel(ctx context.Context, reservationID string) error {
	return m.storage.CancelReservation(ctx, reservationID)
}

func endOfMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
}
