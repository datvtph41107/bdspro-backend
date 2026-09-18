package quota

import (
	"context"
	"errors"
	"time"
)

var ErrAcceptanceLease = errors.New("quota acceptance lease is unavailable")

const acceptanceLeaseSafetyGap = 5 * time.Second

// AcceptanceContext giới hạn thời gian ghi nhận nghiệp vụ trước khi cleanup
// được phép thu hồi reservation. Cả lượt đọc và report dùng cùng quy ước này.
// Reservation đã commit chỉ còn là bằng chứng replay, không còn lease chờ thu
// hồi. Deadline của request gốc vẫn được giữ nguyên trong mọi trường hợp.
func AcceptanceContext(ctx context.Context, reservation Reservation, now time.Time) (context.Context, context.CancelFunc, error) {
	if !reservation.Required || reservation.State == StateCommitted {
		return ctx, func() {}, nil
	}
	deadline := reservation.ExpiresAt.Add(-acceptanceLeaseSafetyGap)
	if reservation.State == StateCanceled || reservation.ExpiresAt.IsZero() || !now.Before(deadline) {
		return nil, nil, ErrAcceptanceLease
	}
	bounded, cancel := context.WithDeadline(ctx, deadline)
	return bounded, cancel, nil
}
