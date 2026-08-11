package access

import (
	commonaccess "common/access"
	"common/businessoperation"
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidSubject   = errors.New("subject is invalid")
	ErrInvalidOperation = errors.New("business operation is invalid")
	ErrUsagePeriod      = errors.New("usage period boundary is not defined")
)

/**
 * Service lấy access hiện tại từ subscription và plan.
 *
 * Service chỉ trả lời "được dùng hay không" và "giới hạn bao nhiêu".
 * Nó không đọc used/reserved và không giữ quota.
 */
type Service struct {
	subscriptions SubscriptionStore
	plans         PlanAccessStore
}

func NewService(subscriptions SubscriptionStore, plans PlanAccessStore) *Service {
	return &Service{
		subscriptions: subscriptions,
		plans:         plans,
	}
}

/**
 * GetAccess trả về quyền sử dụng của operation tại thời điểm now.
 */
func (s *Service) GetAccess(
	ctx context.Context,
	subject commonaccess.Subject,
	operation businessoperation.Code,
	now time.Time,
) (commonaccess.Result, error) {
	if !subject.IsValid() {
		return commonaccess.Result{}, ErrInvalidSubject
	}
	if !operation.IsValid() {
		return commonaccess.Result{}, ErrInvalidOperation
	}

	result := commonaccess.Result{
		Subject:   subject,
		Operation: operation,
		Allowed:   false,
	}

	subscription, found, err := s.subscriptions.FindActiveSubscription(
		ctx,
		string(subject.Type),
		subject.ID,
		now,
	)
	if err != nil {
		return commonaccess.Result{}, err
	}
	if !found || !subscription.AvailableAt(now) {
		return result, nil
	}

	planAccess, found, err := s.plans.FindPlanAccess(
		ctx,
		subscription.PlanVersionID,
		operation,
	)
	if err != nil {
		return commonaccess.Result{}, err
	}
	if !found || !planAccess.Allowed {
		return result, nil
	}

	periodStart, periodEnd, err := usagePeriod(planAccess.Period, subscription)
	if err != nil {
		return commonaccess.Result{}, err
	}

	result.Allowed = true
	result.Unlimited = planAccess.Unlimited
	result.Limit = planAccess.Limit
	result.Period = planAccess.Period
	result.PeriodStart = periodStart
	result.PeriodEnd = periodEnd
	result.SubscriptionID = subscription.ID
	result.PlanCode = subscription.PlanCode
	result.PlanVersion = subscription.PlanVersion

	return result, nil
}

/**
 * usagePeriod chọn đúng khoảng quota mà source hiện tại đã có đủ dữ liệu.
 *
 * Report target đang dùng subscription_cycle.
 * Day/calendar_month cần chốt timezone, lifetime cần chốt mốc bắt đầu/kết thúc trước migration; hiện tại fail thay vì tự đoán.
 */
func usagePeriod(period commonaccess.Period, subscription Subscription) (time.Time, time.Time, error) {
	switch period {
	case commonaccess.PeriodNone:
		return time.Time{}, time.Time{}, nil
	case commonaccess.PeriodSubscriptionCycle:
		if subscription.PeriodStart.IsZero() || subscription.PeriodEnd.IsZero() ||
			!subscription.PeriodEnd.After(subscription.PeriodStart) {
			return time.Time{}, time.Time{}, ErrUsagePeriod
		}
		return subscription.PeriodStart, subscription.PeriodEnd, nil
	case commonaccess.PeriodDay, commonaccess.PeriodCalendarMonth, commonaccess.PeriodLifetime:
		return time.Time{}, time.Time{}, ErrUsagePeriod
	default:
		return time.Time{}, time.Time{}, ErrUsagePeriod
	}
}
