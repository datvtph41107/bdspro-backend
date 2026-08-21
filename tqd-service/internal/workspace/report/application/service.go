package application

import (
	"common/identitycontext"
	commonoperation "common/operation"
	"common/requestcontext"
	"common/requestlog"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"tqd/internal/quota"
	"tqd/internal/usage"
	"tqd/internal/workspace/report/processing"
)

// Service sở hữu create-report flow.
//
// Service điều phối business capability, nhưng không sở hữu:
//   - Access policy evaluation
//   - Redis implementation
//   - PostgreSQL mechanics
//   - Worker execution
//
// Nó chịu trách nhiệm xác định thứ tự:
//
//	Access → Reserve → Acceptance → Quota finalization.
type Service struct {
	source              SourceStore
	acceptance          AcceptanceStore
	accessClient        AccessClient
	quotaService        Quota
	processingAvailable bool
	logger              *slog.Logger
	now                 func() time.Time
}

func NewService(
	source SourceStore,
	acceptance AcceptanceStore,
	accessClient AccessClient,
	quotaService Quota,
	processingAvailable bool,
	logger *slog.Logger,
) *Service {
	if logger == nil {
		logger = slog.Default()
	}

	return &Service{
		source:              source,
		acceptance:          acceptance,
		accessClient:        accessClient,
		quotaService:        quotaService,
		processingAvailable: processingAvailable,
		logger:              logger,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

// CreateGeneratedReport tạo mới hoặc trả lại durable report
// đã được accepted cho cùng command.
//
// Failure semantics:
//
//	BEFORE Acceptance commit:
//	    command có thể fail và reservation phải được release.
//
//	AFTER Acceptance commit:
//	    Report + UsageEvent + Job đã là durable truth.
//	    Redis quota chỉ là runtime projection.
//	    Commit/Cancel failure sau checkpoint này không được
//	    đảo ngược accepted business outcome.
func (s *Service) CreateGeneratedReport(
	ctx context.Context,
	input Input,
) (Result, error) {
	if s == nil ||
		s.source == nil ||
		s.acceptance == nil ||
		s.accessClient == nil ||
		s.quotaService == nil {
		return Result{}, errors.New("report service is not configured")
	}

	logger := requestlog.FromContext(ctx, s.logger)

	actor, ok := identitycontext.FromContext(ctx)
	if !ok {
		return Result{}, ErrActorMissing
	}
	if actor.ProfileID == 0 {
		return Result{}, ErrProfileMissing
	}

	currentOperation, ok := commonoperation.FromContext(ctx)
	if !ok {
		return Result{}, ErrOperationMissing
	}

	operationID, ok := requestcontext.OperationIDFromContext(ctx)
	if !ok {
		return Result{}, ErrOperationIDMissing
	}

	idempotencyKey, ok :=
		requestcontext.IdempotencyKeyFromContext(ctx)
	if !ok || idempotencyKey == "" {
		return Result{}, ErrCommandKeyMissing
	}
	commandKey := idempotencyKey

	requestHash, err := BuildRequestHash(input)
	if err != nil {
		return Result{}, err
	}

	// Current Access is evaluated before replay. Replaying an old durable
	// result must not bypass a permission/entitlement decision that changed
	// after the original command was accepted.
	accessResult, err :=
		s.accessClient.GetAccess(
			ctx,
			currentOperation,
		)
	if err != nil {
		logger.Error(
			"get report access failed",
			"error", err,
		)
		return Result{}, err
	}
	if accessResult.Operation != currentOperation {
		return Result{}, fmt.Errorf(
			"report access returned operation %q, want %q",
			accessResult.Operation,
			currentOperation,
		)
	}
	if !accessResult.Allowed {
		return Result{}, quota.ErrAccessDenied
	}

	// Durable command pre-check happens after current Access. A replay skips
	// Reserve/charge, but it never bypasses today's permission decision.
	existing, found, err :=
		s.acceptance.FindReportByCommand(
			ctx,
			actor.ProfileID,
			commandKey,
		)
	if err != nil {
		return Result{}, err
	}

	if found {
		if existing.RequestHash != "" &&
			existing.RequestHash != requestHash {
			return Result{}, ErrCommandConflict
		}

		logger.Info(
			"generated report already accepted",
			"report_id", existing.ID,
			"outcome", OutcomeExisting,
		)

		return Result{
			Report:  existing,
			Outcome: OutcomeExisting,
		}, nil
	}

	// A durable replay remains readable when processing is degraded, but a
	// new command must not be accepted when no runner can make progress.
	if !s.processingAvailable {
		return Result{}, ErrProcessingUnavailable
	}

	reportData, err :=
		s.buildReport(
			ctx,
			actor.ProfileID,
			commandKey,
			input,
		)
	if err != nil {
		logger.Warn(
			"create report input rejected",
			"error", err,
		)
		return Result{}, err
	}

	reportData.RequestHash = requestHash
	reportData.Operation = currentOperation
	reportData.OperationID = operationID
	reportData.JobID =
		processing.BuildID(
			actor.ProfileID,
			commandKey,
		)

	// Reserve là online quota state.
	//
	// Chưa có business acceptance ở checkpoint này.
	reservation, err :=
		s.quotaService.ReserveQuota(
			ctx,
			quota.ReserveInput{
				Access:         accessResult,
				OperationID:    operationID,
				IdempotencyKey: idempotencyKey,
				Amount: accessResult.
					Metering.
					UnitsPerAction,
				Now: s.now(),
			},
		)
	if err != nil {
		logger.Warn(
			"reserve report quota failed",
			"error", err,
		)
		return Result{}, err
	}

	logger = logger.With(
		"subject_type",
		string(accessResult.Subject.Type),
		"subject_id",
		accessResult.Subject.ID,
		"quota_required",
		reservation.Required,
	)

	if reservation.Required {
		logger = logger.With(
			"reservation_id",
			reservation.ID,
			"quota_limit",
			reservation.Limit,
			"quota_used",
			reservation.Used,
			"quota_reserved",
			reservation.Reserved,
			"quota_remaining",
			reservation.Remaining,
		)
	}

	var usageEvent *usage.Event

	if reservation.Required {
		event := usage.NewEvent(usage.NewEventInput{
			Subject:        reservation.Subject,
			Operation:      reservation.Operation,
			MeterCode:      reservation.MeterCode,
			OperationID:    reservation.OperationID,
			IdempotencyKey: reservation.IdempotencyKey,
			CommandKey:     reservation.CommandKey,
			ReservationID:  reservation.ID,
			Amount:         reservation.Amount,
			PeriodStart:    accessResult.PeriodStart,
			PeriodEnd:      accessResult.PeriodEnd,
			CreatedAt:      s.now(),
			Commercial: usage.CommercialEvidence{
				SubscriptionID: accessResult.SubscriptionID,
				PlanCode:       accessResult.PlanCode,
				PlanVersion:    accessResult.PlanVersion,
				PolicyVersion:  accessResult.Metering.PolicyVersion,
				LimitSnapshot:  accessResult.Limit,
			},
		})

		usageEvent = &event
	}

	acceptedAt := s.now()
	acceptance := Acceptance{
		Report: reportData,
		Usage:  usageEvent,
		Job: processing.Job{
			ID:          reportData.JobID,
			UserID:      reportData.UserID,
			Operation:   reportData.Operation,
			OperationID: reportData.OperationID,
			CommandKey:  reportData.CommandKey,
			Status:      processing.StatusPending,
			AvailableAt: acceptedAt,
			CreatedAt:   acceptedAt,
			UpdatedAt:   acceptedAt,
		},
	}

	// -------------------------------
	// BUSINESS ACCEPTANCE CHECKPOINT
	// -------------------------------
	//
	// Acceptance nói rõ transaction phải commit:
	//
	//   Report
	//   optional UsageEvent
	//   Processing Job
	//
	// PostgreSQL adapter arbitrate command và persist các facts này.
	// Nó không quyết Redis reservation lifecycle.
	accepted, err := s.acceptance.AcceptReport(ctx, acceptance)

	if err != nil {
		// FAILURE BEFORE ACCEPTANCE.
		//
		// Business state chưa commit.
		// Reservation phải được release.
		cancelErr :=
			s.cancelReservation(
				ctx,
				reservation,
				logger,
			)

		if cancelErr != nil {
			return Result{},
				errors.Join(
					err,
					cancelErr,
				)
		}

		return Result{}, err
	}

	saved := accepted.Report
	created := accepted.Created

	// ---------------------------------
	// FROM HERE BUSINESS IS ACCEPTED
	// ---------------------------------
	//
	// Report + UsageEvent + Job đã là durable truth.
	//
	// Redis chỉ là runtime quota projection.
	// Failure sau điểm này không được đổi business outcome.
	outcome := OutcomeExisting
	if created {
		outcome = OutcomeAccepted
	}

	// Race case:
	//
	// Transaction phát hiện durable command khác đã thắng.
	// Reservation hiện tại trở thành redundant reservation.
	if reservation.Required &&
		!reservationMatchesDurableUsage(
			reservation,
			accepted.DurableUsage,
		) {

		if err :=
			s.cancelReservation(
				ctx,
				reservation,
				logger,
			); err != nil {

			// Operational debt.
			//
			// Không return business failure.
			logger.Error(
				"redundant quota reservation pending repair",
				"report_id",
				saved.ID,
				"reservation_id",
				reservation.ID,
				"error",
				err,
			)
		}

		logger.Info(
			"generated report race reused existing durable command",
			"report_id",
			saved.ID,
			"outcome",
			OutcomeExisting,
		)

		return Result{
			Report:  saved,
			Outcome: OutcomeExisting,
		}, nil
	}

	// Happy path:
	// finalize online quota projection.
	committed, err :=
		s.quotaService.CommitQuota(
			ctx,
			reservation,
		)

	if err != nil {
		// IMPORTANT:
		//
		// DB đã accepted.
		// Đây không còn là application/business failure.
		//
		// Durable UsageEvent cho phép Repair hoàn tất
		// Redis projection sau đó.
		logger.Error(
			"accepted report quota projection pending repair",
			"report_id",
			saved.ID,
			"reservation_id",
			reservation.ID,
			"usage_key",
			usageKey(usageEvent),
			"error",
			err,
		)

		return Result{
			Report:  saved,
			Outcome: outcome,
		}, nil
	}

	logger.Info(
		"generated report accepted",
		"report_id",
		saved.ID,
		"outcome",
		outcome,
		"reservation_state",
		committed.State,
		"usage_key",
		usageKey(usageEvent),
	)

	return Result{
		Report:  saved,
		Outcome: outcome,
	}, nil
}

func (s *Service) cancelReservation(
	ctx context.Context,
	reservation quota.Reservation,
	logger *slog.Logger,
) error {
	if !reservation.Required {
		return nil
	}

	_, err :=
		s.quotaService.CancelQuota(
			ctx,
			reservation,
		)

	if err != nil {
		logger.Error(
			"cancel report quota failed",
			"reservation_id",
			reservation.ID,
			"error",
			err,
		)

		return fmt.Errorf(
			"cancel quota reservation: %w",
			err,
		)
	}

	logger.Info(
		"report quota canceled",
		"reservation_id",
		reservation.ID,
	)

	return nil
}

// reservationMatchesDurableUsage trả lời bằng durable evidence,
// không bằng một boolean do persistence adapter suy diễn.
//
// Với current Redis-D bridge, chỉ reservation đã tạo usage event
// của command thắng mới được finalize. Reservation của race loser phải cancel.
func reservationMatchesDurableUsage(
	reservation quota.Reservation,
	event *usage.Event,
) bool {
	if !reservation.Required {
		return true
	}
	return event != nil &&
		event.ReservationID == reservation.ID
}

func usageKey(
	event *usage.Event,
) string {
	if event == nil {
		return ""
	}

	return event.UsageKey
}
