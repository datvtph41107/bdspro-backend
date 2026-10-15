package evaluate

import (
	commonmetering "common/metering"
	"common/operation"
	"context"
	"errors"
	"time"
	"userv4/modules/user/domain/entitlement"
)

var (
	ErrInvalidSubject        = errors.New("subject is invalid")
	ErrInvalidOperation      = errors.New("operation is invalid")
	ErrAmbiguousSubscription = errors.New("multiple current subscriptions grant operation")
	ErrUsagePeriod           = errors.New("usage period boundary is not defined")
)

// Service giải quyết quyền sử dụng tại một thời điểm cụ thể từ subscription
// và PlanVersion. Đây là owner của entitlement: service chỉ trả lời được dùng
// hay không, meter nào và giới hạn bao nhiêu; nó không đọc số đã dùng, không
// giữ quota và không thực hiện thanh toán.
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

// GetAccess trả về bằng chứng quyền của operation tại thời điểm now. Kết quả
// được TQD lưu kèm usage event để sau này vẫn truy được người dùng đã được
// cấp quyền theo gói và phiên bản nào, kể cả khi catalog đã thay đổi.
func (s *Service) GetAccess(
	ctx context.Context,
	subject access.Subject,
	operationCode operation.Code,
	now time.Time,
) (access.Result, error) {
	if !subject.IsValid() {
		return access.Result{}, ErrInvalidSubject
	}
	if !operationCode.IsValid() {
		return access.Result{}, ErrInvalidOperation
	}

	result := access.Result{
		Subject:   subject,
		Operation: operationCode,
		Allowed:   false,
	}

	subscription, found, err := s.subscriptions.FindSubscriptionForOperation(
		ctx,
		string(subject.Type),
		subject.ID,
		operationCode,
		now,
	)
	if err != nil {
		return access.Result{}, err
	}
	if !found || !subscription.AvailableAt(now) {
		return result, nil
	}

	planAccess, found, err := s.plans.FindPlanAccess(
		ctx,
		subscription.PlanVersionID,
		operationCode,
	)
	if err != nil {
		return access.Result{}, err
	}
	if !found || !planAccess.Allowed {
		return result, nil
	}

	periodStart, periodEnd, err := usagePeriod(planAccess.Period, subscription)
	if err != nil {
		return access.Result{}, err
	}

	result.Allowed = true
	result.Metering = access.Metering{
		FeatureCode:    planAccess.FeatureCode,
		MeterCode:      commonmetering.Code(planAccess.MeterCode),
		UnitsPerAction: planAccess.UnitsPerAction,
		PolicyVersion:  subscription.PlanVersion,
	}
	result.Unlimited = planAccess.Unlimited
	result.Limit = planAccess.Limit
	result.Period = planAccess.Period
	result.PeriodStart = periodStart
	result.PeriodEnd = periodEnd
	result.SubscriptionID = subscription.ID
	result.PlanCode = subscription.PlanCode
	result.PlanVersion = subscription.PlanVersion
	if !result.IsValid() {
		return access.Result{}, errors.New("resolved access evidence is invalid")
	}

	return result, nil
}

// usagePeriod chọn khoảng thời gian mà quota áp dụng. Hiện source chỉ đủ dữ
// liệu chắc chắn cho subscription_cycle. Day/calendar_month cần quy ước
// timezone, còn lifetime cần mốc bắt đầu/kết thúc; khi chưa có quy ước, hàm
// fail-closed thay vì tự đoán và tính sai quyền khách hàng.
func usagePeriod(period access.Period, subscription Subscription) (time.Time, time.Time, error) {
	switch period {
	case access.PeriodNone:
		return time.Time{}, time.Time{}, nil
	case access.PeriodSubscriptionCycle:
		if subscription.PeriodStart.IsZero() || subscription.PeriodEnd.IsZero() ||
			!subscription.PeriodEnd.After(subscription.PeriodStart) {
			return time.Time{}, time.Time{}, ErrUsagePeriod
		}
		return subscription.PeriodStart, subscription.PeriodEnd, nil
	case access.PeriodDay, access.PeriodCalendarMonth, access.PeriodLifetime:
		return time.Time{}, time.Time{}, ErrUsagePeriod
	default:
		return time.Time{}, time.Time{}, ErrUsagePeriod
	}
}
