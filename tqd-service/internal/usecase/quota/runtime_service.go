package quota

import (
	"common/request"
	"context"
	"errors"
	"strings"
	"time"
	"tqd/internal/access"

	"github.com/google/uuid"
)

const defaultReservationTTL = 5 * time.Minute

/**
 * ReserveInput gom access đã được xác nhận với identity của command hiện tại.
 */
type ReserveInput struct {
	Access access.Result

	OperationID    string
	IdempotencyKey string

	Amount int64
	Now    time.Time
	TTL    time.Duration
}

/**
 * Service kiểm tra quota runtime và giữ quota trước khi business state thay đổi.
 */
type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

/**
 * ReserveQuota giữ trước quota cho command hiện tại.
 *
 * Idempotency-Key được ưu tiên làm command key.
 * Nếu không có thì dùng Operation-ID.
 */
func (s *Service) ReserveQuota(ctx context.Context, input ReserveInput) (Reservation, error) {
	if s == nil || s.store == nil {
		return Reservation{}, errors.New("quota store is not configured")
	}
	if !input.Access.IsValid() {
		return Reservation{}, ErrInvalidInput
	}
	if !input.Access.Allowed {
		return Reservation{}, ErrAccessDenied
	}
	if !input.Access.UsesQuota() {
		return Reservation{
			Required:  false,
			Subject:   input.Access.Subject,
			Operation: input.Access.Operation,
			MeterCode: input.Access.Metering.MeterCode,
		}, nil
	}
	if input.Amount <= 0 {
		return Reservation{}, ErrInvalidInput
	}

	operationID := strings.TrimSpace(input.OperationID)
	if !request.IsValidOperationID(operationID) {
		return Reservation{}, ErrInvalidInput
	}

	idempotencyKey := strings.TrimSpace(input.IdempotencyKey)
	if idempotencyKey != "" && !request.IsValidIdempotencyKey(idempotencyKey) {
		return Reservation{}, ErrInvalidInput
	}

	commandKey := operationID
	if idempotencyKey != "" {
		commandKey = idempotencyKey
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}
	ttl := input.TTL
	if ttl <= 0 {
		ttl = defaultReservationTTL
	}

	return s.store.ReserveQuota(ctx, StoreReserveInput{
		ReservationID:  uuid.NewString(),
		Subject:        input.Access.Subject,
		Operation:      input.Access.Operation,
		MeterCode:      input.Access.Metering.MeterCode,
		OperationID:    operationID,
		IdempotencyKey: idempotencyKey,
		CommandKey:     commandKey,
		Amount:         input.Amount,
		Limit:          input.Access.Limit,
		PeriodStart:    input.Access.PeriodStart,
		PeriodEnd:      input.Access.PeriodEnd,
		ExpiresAt:      now.Add(ttl),
	})
}

/**
 * CommitQuota xác nhận phần quota đã được giữ trước đó.
 *
 * Gọi lại với reservation đã commit không được trừ quota thêm lần nữa.
 */
func (s *Service) CommitQuota(ctx context.Context, reservation Reservation) (Reservation, error) {
	if !reservation.Required {
		return reservation, nil
	}
	if reservation.ID == "" {
		return Reservation{}, ErrInvalidInput
	}
	return s.store.CommitQuota(ctx, reservation.ID)
}

/**
 * CancelQuota trả lại quota đang giữ khi business action chưa được accepted.
 */
func (s *Service) CancelQuota(ctx context.Context, reservation Reservation) (Reservation, error) {
	if !reservation.Required {
		return reservation, nil
	}
	if reservation.ID == "" {
		return Reservation{}, ErrInvalidInput
	}
	return s.store.CancelQuota(ctx, reservation.ID)
}
