package evaluate

import (
	"time"
	"user/internal/domain/entitlement"
)

/**
 * Subscription là dữ liệu tối thiểu cần để kiểm tra quyền sử dụng hiện tại.
 */
type Subscription struct {
	ID            uint64
	Subject       access.Subject
	PlanVersionID uint64
	PlanCode      string
	PlanVersion   string
	Status        string

	StartedAt   time.Time
	PeriodStart time.Time
	PeriodEnd   time.Time
	AccessUntil *time.Time
}

/**
 * AvailableAt kiểm tra subscription còn cho phép sử dụng tại thời điểm now.
 *
 * active dùng current period.
 * past_due/canceled chỉ tiếp tục nếu access_until còn hiệu lực.
 */
func (s Subscription) AvailableAt(now time.Time) bool {
	if s.ID == 0 || !s.Subject.IsValid() || s.PlanVersionID == 0 || s.StartedAt.IsZero() {
		return false
	}
	if now.Before(s.StartedAt) {
		return false
	}

	switch s.Status {
	case "active":
		if s.AccessUntil != nil {
			return now.Before(*s.AccessUntil)
		}
		return s.PeriodEnd.IsZero() || now.Before(s.PeriodEnd)

	case "past_due", "canceled":
		return s.AccessUntil != nil && now.Before(*s.AccessUntil)

	default:
		return false
	}
}
