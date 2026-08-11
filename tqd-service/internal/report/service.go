package report

import (
	"common/businessoperation"
	"common/identitycontext"
	"common/requestcontext"
	"common/requestlog"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"tqd/internal/quota"
	"tqd/internal/reportjob"
	"tqd/internal/usage"
)

/**
 * Service giữ flow tạo generated report ở một chỗ dễ đọc.
 *
 * Flow chính:
 * 1. Lấy Actor + Operation-ID + BusinessOperation từ context.
 * 2. Dựng report data từ parcel/region.
 * 3. Lấy access từ user-service.
 * 4. Reserve quota nếu operation có quota.
 * 5. Lưu report + usage event trong cùng DB transaction.
 * 6. Commit reservation sau khi DB đã accepted report.
 */
type Service struct {
	source       SourceStore
	reports      ReportStore
	accessClient AccessClient
	quotaService *quota.Service
	logger       *slog.Logger
	now          func() time.Time
}

func NewService(
	source SourceStore,
	reports ReportStore,
	accessClient AccessClient,
	quotaService *quota.Service,
	logger *slog.Logger,
) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		source:       source,
		reports:      reports,
		accessClient: accessClient,
		quotaService: quotaService,
		logger:       logger,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

/**
 * CreateGeneratedReport tạo hoặc trả lại report đã được accepted cho cùng command.
 *
 * Nếu DB transaction thất bại, reservation được cancel vì business state chưa commit.
 * Nếu DB đã commit nhưng Redis commit quota thất bại, report vẫn là accepted;
 * request trả error để retry cùng command có thể hoàn tất quota mà không tạo report mới.
 */
func (s *Service) CreateGeneratedReport(
	ctx context.Context,
	input Input,
) (Result, error) {
	if s == nil || s.source == nil || s.reports == nil || s.accessClient == nil || s.quotaService == nil {
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

	operation, ok := businessoperation.FromContext(ctx)
	if !ok {
		return Result{}, ErrOperationMissing
	}
	if operation != businessoperation.WorkspaceReportGenerate {
		return Result{}, ErrOperationMismatch
	}

	operationID, ok := requestcontext.OperationIDFromContext(ctx)
	if !ok {
		return Result{}, ErrOperationIDMissing
	}
	idempotencyKey, _ := requestcontext.IdempotencyKeyFromContext(ctx)
	commandKey := operationID
	if idempotencyKey != "" {
		commandKey = idempotencyKey
	}

	requestHash, err := BuildRequestHash(input)
	if err != nil {
		return Result{}, err
	}

	existing, found, err := s.reports.FindReportByCommand(ctx, actor.ProfileID, commandKey)
	if err != nil {
		return Result{}, err
	}
	if found {
		if existing.RequestHash != "" && existing.RequestHash != requestHash {
			return Result{}, ErrCommandConflict
		}
		logger.Info(
			"generated report already accepted",
			"report_id", existing.ID,
			"outcome", OutcomeExisting,
		)
		return Result{Report: existing, Outcome: OutcomeExisting}, nil
	}

	reportData, err := s.buildReport(ctx, actor.ProfileID, commandKey, input)
	if err != nil {
		logger.Warn("create report input rejected", "error", err)
		return Result{}, err
	}
	reportData.RequestHash = requestHash
	reportData.Operation = operation
	reportData.OperationID = operationID
	reportData.JobID = reportjob.BuildID(actor.ProfileID, commandKey)

	accessResult, err := s.accessClient.GetReportAccess(ctx)
	if err != nil {
		logger.Error("get report access failed", "error", err)
		return Result{}, err
	}
	if accessResult.Operation != operation {
		return Result{}, fmt.Errorf("report access returned operation %q, want %q", accessResult.Operation, operation)
	}

	amount, hasAmount := quota.AmountForOperation(operation)
	if accessResult.UsesQuota() && !hasAmount {
		return Result{}, fmt.Errorf("quota amount is missing for operation %q", operation)
	}

	reservation, err := s.quotaService.ReserveQuota(ctx, quota.ReserveInput{
		Access:         accessResult,
		OperationID:    operationID,
		IdempotencyKey: idempotencyKey,
		Amount:         amount,
		Now:            s.now(),
	})
	if err != nil {
		logger.Warn("reserve report quota failed", "error", err)
		return Result{}, err
	}

	logger = logger.With(
		"subject_type", string(accessResult.Subject.Type),
		"subject_id", accessResult.Subject.ID,
		"quota_required", reservation.Required,
	)
	if reservation.Required {
		logger = logger.With(
			"reservation_id", reservation.ID,
			"quota_limit", reservation.Limit,
			"quota_used", reservation.Used,
			"quota_reserved", reservation.Reserved,
			"quota_remaining", reservation.Remaining,
		)
	}

	var usageEvent *usage.Event
	if reservation.Required {
		event := usage.NewEvent(
			reservation.Subject,
			reservation.Operation,
			reservation.OperationID,
			reservation.IdempotencyKey,
			reservation.CommandKey,
			reservation.ID,
			reservation.Amount,
			accessResult.PeriodStart,
			accessResult.PeriodEnd,
			s.now(),
		)
		usageEvent = &event
	}

	saved, created, finalizeReservation, err := s.reports.AcceptReport(ctx, reportData, usageEvent)
	if err != nil {
		cancelErr := s.cancelReservation(ctx, reservation, logger)
		if cancelErr != nil {
			return Result{}, errors.Join(err, cancelErr)
		}
		return Result{}, err
	}

	// Từ đây DB đã commit. Không rollback report chỉ vì bước trả response hay Redis gặp lỗi.
	if reservation.Required && !finalizeReservation {
		if err := s.cancelReservation(ctx, reservation, logger); err != nil {
			return Result{}, err
		}
		logger.Info(
			"generated report race reused existing durable command",
			"report_id", saved.ID,
			"outcome", OutcomeExisting,
		)
		return Result{Report: saved, Outcome: OutcomeExisting}, nil
	}

	committed, err := s.quotaService.CommitQuota(ctx, reservation)
	if err != nil {
		logger.Error(
			"report accepted but quota commit failed",
			"report_id", saved.ID,
			"usage_key", usageKey(usageEvent),
			"error", err,
		)
		return Result{}, fmt.Errorf("%w: %v", ErrQuotaCommitFailed, err)
	}

	outcome := OutcomeExisting
	if created {
		outcome = OutcomeAccepted
	}

	logger.Info(
		"generated report accepted",
		"report_id", saved.ID,
		"outcome", outcome,
		"reservation_state", committed.State,
		"usage_key", usageKey(usageEvent),
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

	_, err := s.quotaService.CancelQuota(ctx, reservation)
	if err != nil {
		logger.Error(
			"cancel report quota failed",
			"reservation_id", reservation.ID,
			"error", err,
		)
		return fmt.Errorf("cancel quota reservation: %w", err)
	}

	logger.Info(
		"report quota canceled",
		"reservation_id", reservation.ID,
	)
	return nil
}

func usageKey(event *usage.Event) string {
	if event == nil {
		return ""
	}
	return event.UsageKey
}
