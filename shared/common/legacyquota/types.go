// common/quota/types.go
package legacyquota

import (
	"context"
	"time"
)

// QuotaInfo thông tin quota hiện tại của user trong một pool
type QuotaInfo struct {
	UserID    string    `json:"user_id"`
	Pool      string    `json:"pool"`
	Used      int64     `json:"used"`
	Limit     int64     `json:"limit"`
	Remaining int64     `json:"remaining"`
	ResetAt   time.Time `json:"reset_at"`
}

// ConsumeResult kết quả sau khi consume quota
type ConsumeResult struct {
	Success   bool       `json:"success"`
	QuotaInfo *QuotaInfo `json:"quota_info,omitempty"`
	Error     error      `json:"error,omitempty"`
}

// Reservation đại diện cho một reservation (2-phase commit)
type Reservation struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Pool      string    `json:"pool"`
	Amount    int64     `json:"amount"`
	ExpiresAt time.Time `json:"expires_at"`
}

// QuotaStorage interface cho các thao tác lưu trữ quota
type QuotaStorage interface {
	// GetUsage trả về số quota đã dùng của user trong pool
	GetUsage(ctx context.Context, userID, pool string) (int64, error)

	// GetLimit trả về giới hạn quota của pool cho user
	GetLimit(ctx context.Context, userID, pool string) (int64, error)

	// IncrementUsage tăng usage lên amount, trả về usage mới. Phải atomic.
	IncrementUsage(ctx context.Context, userID, pool string, amount int64) (int64, error)

	// SetUsage (chỉ dùng cho reset hoặc test)
	SetUsage(ctx context.Context, userID, pool string, usage int64, ttl time.Duration) error

	// CreateReservation tạo reservation với TTL
	CreateReservation(ctx context.Context, userID, pool string, amount int64, ttl time.Duration) (string, error)

	// ConfirmReservation xác nhận reservation và consume quota
	ConfirmReservation(ctx context.Context, reservationID string) error

	// CancelReservation huỷ reservation (giải phóng)
	CancelReservation(ctx context.Context, reservationID string) error
}

// PolicyEngine định nghĩa cách xác định pool và cost cho một action
type PolicyEngine interface {
	// Resolve trả về pool name và cost cho service+action
	Resolve(service, action string) (pool string, cost int64, err error)

	// GetUserLimit trả về limit của user cho một pool (có thể override theo plan)
	GetUserLimit(ctx context.Context, userID, pool string) (int64, error)
}

// QuotaManager là interface chính mà các service sẽ dùng
type QuotaManager interface {
	// CheckOnly kiểm tra quota còn đủ không (không trừ)
	CheckOnly(ctx context.Context, userID, service, action string) (*QuotaInfo, error)

	// Consume trừ quota ngay lập tức (chỉ gọi khi đã chắc chắn thành công)
	Consume(ctx context.Context, userID, service, action string) (*ConsumeResult, error)

	// Reserve tạo reservation, dùng cho 2-phase commit
	Reserve(ctx context.Context, userID, service, action string, ttl time.Duration) (*Reservation, error)

	// Confirm xác nhận reservation thành công
	Confirm(ctx context.Context, reservationID string) error

	// Cancel huỷ reservation
	Cancel(ctx context.Context, reservationID string) error
}
