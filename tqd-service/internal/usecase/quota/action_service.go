package quota

import (
	commonoperation "common/operation"
	"common/request"
	"common/requestlog"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"tqd/internal/usecase/usage"
)

var (
	ErrActionUnavailable = errors.New("commercial action accounting is unavailable")
	ErrActionIdentity    = errors.New("Idempotency-Key is required for a metered action")
)

// ReadActionUsageStore chỉ cung cấp bằng chứng cần để ghi nhận lượt đọc.
// Binding và acceptance thuộc cùng owner PostgreSQL. Acceptance trả identity
// của reservation thắng để caller không commit nhầm phần giữ dư.
type ReadActionUsageStore interface {
	BindReadCommand(context.Context, usage.ReadCommand) (bool, error)
	AcceptReadCommand(context.Context, usage.ReadCommand, *usage.Event) (string, error)
}

// ActionService ghi nhận quota cho hành động đọc đồng bộ, không có side effect
// nghiệp vụ trong callback. Thành công nghĩa là handler tạo được response và
// database xác nhận usage; không thể bảo đảm client đã nhận byte qua mạng.
// Report bất đồng bộ dùng transaction Report + Usage + Job riêng.
type ActionService struct {
	access ProfileAccessReader
	quota  *Service
	usage  ReadActionUsageStore
	now    func() time.Time
}

func NewActionService(accessReader ProfileAccessReader, quotaService *Service, usageStore ReadActionUsageStore) *ActionService {
	return &ActionService{
		access: accessReader,
		quota:  quotaService,
		usage:  usageStore,
		now:    func() time.Time { return time.Now().UTC() },
	}
}

// Run xử lý một ý định nghiệp vụ duy nhất của người dùng. Cùng một
// Idempotency-Key có thể được gửi lại khi mạng lỗi nhưng không được tạo thêm
// usage event hoặc trừ quota lần thứ hai.
func (s *ActionService) Run(ctx context.Context, operationCode commonoperation.Code, requestHash string, action func(context.Context) error) error {
	if s == nil || s.access == nil || s.quota == nil || s.usage == nil || action == nil {
		return ErrActionUnavailable
	}
	if !operationCode.IsValid() {
		return ErrInvalidInput
	}

	decision, err := s.access.GetAccess(ctx, operationCode)
	if err != nil {
		return errors.Join(ErrActionUnavailable, err)
	}
	if decision.Operation != operationCode {
		return fmt.Errorf("%w: access returned operation %q", ErrActionUnavailable, decision.Operation)
	}
	if !decision.IsValid() {
		return fmt.Errorf("%w: invalid access evidence", ErrActionUnavailable)
	}
	if !decision.Allowed {
		return &AccessDeniedError{Subject: decision.Subject, Operation: operationCode}
	}

	operationID, operationOK := request.OperationIDFromContext(ctx)
	idempotencyKey, idempotencyOK := request.IdempotencyKeyFromContext(ctx)
	if !operationOK || !idempotencyOK || idempotencyKey == "" {
		return ErrActionIdentity
	}

	command := usage.ReadCommand{Subject: decision.Subject, Operation: operationCode, CommandKey: idempotencyKey, RequestHash: requestHash}
	if !command.IsValid() {
		return ErrActionIdentity
	}
	accepted, err := s.usage.BindReadCommand(ctx, command)
	if err != nil {
		return errors.Join(ErrActionUnavailable, err)
	}
	// Read replay vẫn cần quyền hiện tại nhưng không reserve/tính lượt mới,
	// kể cả kỳ cũ đã hết hoặc meter của gói vừa thay đổi.
	if accepted {
		return action(ctx)
	}
	if !decision.UsesQuota() {
		if err := action(ctx); err != nil {
			return err
		}
		_, err := s.usage.AcceptReadCommand(ctx, command, nil)
		return err
	}

	reservation, err := s.quota.ReserveQuota(ctx, ReserveInput{
		Access:         decision,
		OperationID:    operationID,
		IdempotencyKey: idempotencyKey,
		Amount:         decision.Metering.UnitsPerAction,
		Now:            s.now(),
	})
	if err != nil {
		return err
	}

	logger := requestlog.FromContext(ctx, slog.Default()).With("reservation_id", reservation.ID, "operation", operationCode)
	// Lease thuộc command, có thể dùng chung bởi nhiều request retry. Không
	// hủy khi một request lỗi vì request khác có thể đang ghi usage. Cleanup
	// hết lease sẽ quyết định commit/cancel dựa trên database.
	acceptanceCtx, cancel, err := AcceptanceContext(ctx, reservation, s.now())
	if err != nil {
		return errors.Join(ErrActionUnavailable, err)
	}
	defer cancel()
	if err := action(acceptanceCtx); err != nil {
		logger.Warn("read failed; quota lease retained for reconciliation", "error", err)
		return err
	}
	// Callback phải tôn trọng deadline. Kiểm tra lại để không ghi usage sau
	// khi lease đã hết ngay cả khi callback bỏ qua context.
	if err := acceptanceCtx.Err(); err != nil {
		return errors.Join(ErrActionUnavailable, err)
	}

	event := usage.NewEvent(usage.NewEventInput{
		Subject:        reservation.Subject,
		Operation:      reservation.Operation,
		MeterCode:      reservation.MeterCode,
		OperationID:    reservation.OperationID,
		IdempotencyKey: reservation.IdempotencyKey,
		CommandKey:     reservation.CommandKey,
		ReservationID:  reservation.ID,
		Amount:         reservation.Amount,
		PeriodStart:    decision.PeriodStart,
		PeriodEnd:      decision.PeriodEnd,
		CreatedAt:      s.now(),
		Commercial: usage.CommercialEvidence{
			SubscriptionID: decision.SubscriptionID,
			PlanCode:       decision.PlanCode,
			PlanVersion:    decision.PlanVersion,
			PolicyVersion:  decision.Metering.PolicyVersion,
			LimitSnapshot:  decision.Limit,
		},
	})
	acceptedReservationID, err := s.usage.AcceptReadCommand(acceptanceCtx, command, &event)
	if err != nil {
		// Timeout/mất phản hồi không chứng minh insert đã rollback. Giữ lease
		// để cleanup kiểm tra bằng chứng thay vì mở lại quota có thể đã dùng.
		logger.Warn("read usage outcome unknown; quota lease retained for reconciliation", "error", err)
		return errors.Join(ErrActionUnavailable, err)
	}
	if acceptedReservationID != reservation.ID {
		// Transaction khác đã chấp nhận command. Chỉ reservation có usage
		// được commit; phần giữ dư có thể hủy vì winner đã là bằng chứng bền vững.
		if _, err := s.quota.CancelQuota(ctx, reservation); err != nil {
			logger.Warn("redundant read reservation pending cleanup", "error", err)
		}
		return nil
	}
	if _, err := s.quota.CommitQuota(ctx, reservation); err != nil {
		// Usage event trong database đã là nguồn sự thật. Nếu Redis lỗi, tiến
		// trình reconciliation sẽ sửa projection; retry cùng command vẫn an toàn.
		logger.Error("read accepted; quota projection pending repair", "error", err)
	}
	return nil
}
