package quota

import (
	commonoperation "common/operation"
	"common/request"
	"context"
	"errors"
	"fmt"
	"time"
	"tqd/internal/usecase/usage"
)

var (
	ErrActionUnavailable = errors.New("commercial action accounting is unavailable")
	ErrActionIdentity    = errors.New("Idempotency-Key is required for a metered action")
)

// ActionService surrounds a successful, read-only TQD action with commercial
// admission and durable usage accounting. The caller supplies the actual
// action, so quota is spent only when useful data was produced.
type ActionService struct {
	access ProfileAccessReader
	quota  *Service
	usage  usage.Store
	now    func() time.Time
}

func NewActionService(accessReader ProfileAccessReader, quotaService *Service, usageStore usage.Store) *ActionService {
	return &ActionService{
		access: accessReader,
		quota:  quotaService,
		usage:  usageStore,
		now:    func() time.Time { return time.Now().UTC() },
	}
}

// Run authorizes and accounts one logical user intent. A repeated request with
// the same Idempotency-Key can re-read the projection but cannot create a
// second durable usage event or spend quota twice.
func (s *ActionService) Run(ctx context.Context, operationCode commonoperation.Code, action func(context.Context) error) error {
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
	if !decision.Allowed {
		return &AccessDeniedError{Subject: decision.Subject, Operation: operationCode}
	}

	// An enabled, unmetered feature needs authorization but no accounting
	// identity or Redis reservation.
	if !decision.UsesQuota() {
		return action(ctx)
	}

	operationID, operationOK := request.OperationIDFromContext(ctx)
	idempotencyKey, idempotencyOK := request.IdempotencyKeyFromContext(ctx)
	if !operationOK || !idempotencyOK || idempotencyKey == "" {
		return ErrActionIdentity
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

	if err := action(ctx); err != nil {
		_, _ = s.quota.CancelQuota(ctx, reservation)
		return err
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
	if _, err := s.usage.SaveUsageEvent(ctx, event); err != nil {
		_, _ = s.quota.CancelQuota(ctx, reservation)
		return errors.Join(ErrActionUnavailable, err)
	}
	if _, err := s.quota.CommitQuota(ctx, reservation); err != nil {
		// The durable event is now authoritative. Reconciliation will repair the
		// Redis projection, and a retry with the same command remains idempotent.
		return errors.Join(ErrActionUnavailable, err)
	}
	return nil
}
