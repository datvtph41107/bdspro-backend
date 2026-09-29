package publish

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	catalogdomain "user/internal/domain/plan"
)

var (
	ErrPlanVersionNotFound = errors.New("catalog plan version not found")
	ErrPlanVersionNotDraft = errors.New("catalog plan version is not draft")
	ErrConcurrentPublish   = errors.New("catalog plan version changed during publication")
	ErrParentRetired       = errors.New("catalog product or plan is retired")
)

/** Transaction executes one publication atomically. */
type Transaction interface {
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}

/** Repository persists catalog publication state. */
type Repository interface {
	LoadPlanVersionForUpdate(
		ctx context.Context,
		planVersionID uint64,
	) (catalogdomain.PlanVersionAggregate, error)
	MarkPlanVersionPublished(
		ctx context.Context,
		input PublishRecord,
	) error
}

/** Command publishes one reviewed draft. */
type Command struct {
	PlanVersionID  uint64
	EffectiveFrom  time.Time
	EffectiveUntil *time.Time
	ActorID        uint64
	PublishedAt    time.Time
}

/** PublishRecord is the guarded persistence update. */
type PublishRecord struct {
	PlanVersionID  uint64
	ProductID      uint64
	PlanID         uint64
	EffectiveFrom  time.Time
	EffectiveUntil *time.Time
	PublishedAt    time.Time
	PublishedBy    uint64
	TermsChecksum  string
}

/** Result describes the immutable published contract. */
type Result struct {
	PlanVersionID uint64
	ProductCode   string
	PlanCode      string
	Version       string
	TermsChecksum string
	PublishedAt   time.Time
}

/** Service validates and publishes plan versions. */
type Service struct {
	transaction Transaction
	repository  Repository
	now         func() time.Time
}

/** NewService creates the publication service. */
func NewService(transaction Transaction, repository Repository) *Service {
	return &Service{
		transaction: transaction,
		repository:  repository,
		now:         time.Now,
	}
}

/** Publish validates and commits one draft in one transaction. */
func (s *Service) Publish(ctx context.Context, command Command) (Result, error) {
	if s == nil || s.transaction == nil || s.repository == nil {
		return Result{}, fmt.Errorf("catalog publication service is not configured")
	}
	if command.PlanVersionID == 0 {
		return Result{}, fmt.Errorf("plan version id is required")
	}
	if command.ActorID == 0 {
		return Result{}, fmt.Errorf("publisher actor id is required")
	}
	if command.EffectiveFrom.IsZero() {
		return Result{}, fmt.Errorf("effective_from is required")
	}
	if command.EffectiveUntil != nil &&
		!command.EffectiveUntil.After(command.EffectiveFrom) {
		return Result{}, fmt.Errorf("effective window is invalid")
	}

	publishedAt := command.PublishedAt
	if publishedAt.IsZero() {
		publishedAt = s.now().UTC()
	}
	if publishedAt.After(command.EffectiveFrom) {
		return Result{}, fmt.Errorf("publication cannot occur after effective_from")
	}

	var result Result
	err := s.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		aggregate, err := s.repository.LoadPlanVersionForUpdate(
			txCtx,
			command.PlanVersionID,
		)
		if err != nil {
			return err
		}
		if aggregate.Status != catalogdomain.StatusDraft {
			return ErrPlanVersionNotDraft
		}
		if aggregate.ProductStatus == catalogdomain.StatusRetired ||
			aggregate.PlanStatus == catalogdomain.StatusRetired {
			return ErrParentRetired
		}
		if aggregate.ProductStatus != catalogdomain.StatusDraft &&
			aggregate.ProductStatus != catalogdomain.StatusActive {
			return fmt.Errorf("invalid product status %q", aggregate.ProductStatus)
		}
		if aggregate.PlanStatus != catalogdomain.StatusDraft &&
			aggregate.PlanStatus != catalogdomain.StatusActive {
			return fmt.Errorf("invalid plan status %q", aggregate.PlanStatus)
		}
		if strings.TrimSpace(aggregate.ProductDisplayName) == "" {
			return fmt.Errorf("product display name is required")
		}

		contract := aggregate.Contract(
			catalogdomain.StatusActive,
			command.EffectiveFrom,
			command.EffectiveUntil,
		)
		_, err = catalogdomain.NewRegistry(catalogdomain.Snapshot{
			Version: "1.0.0",
			Products: []catalogdomain.Product{
				{
					Code:        aggregate.ProductCode,
					DisplayName: aggregate.ProductDisplayName,
					Plans:       []catalogdomain.PlanVersion{contract},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("invalid publication contract: %w", err)
		}

		checksum, err := catalogdomain.TermsChecksum(contract)
		if err != nil {
			return fmt.Errorf("calculate terms checksum: %w", err)
		}
		record := PublishRecord{
			PlanVersionID:  command.PlanVersionID,
			ProductID:      aggregate.ProductID,
			PlanID:         aggregate.PlanID,
			EffectiveFrom:  command.EffectiveFrom,
			EffectiveUntil: copyTime(command.EffectiveUntil),
			PublishedAt:    publishedAt,
			PublishedBy:    command.ActorID,
			TermsChecksum:  checksum,
		}
		if err := s.repository.MarkPlanVersionPublished(txCtx, record); err != nil {
			return err
		}

		result = Result{
			PlanVersionID: command.PlanVersionID,
			ProductCode:   aggregate.ProductCode,
			PlanCode:      aggregate.PlanCode,
			Version:       aggregate.Version,
			TermsChecksum: checksum,
			PublishedAt:   publishedAt,
		}
		return nil
	})
	if err != nil {
		return Result{}, err
	}
	return result, nil
}

func copyTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}
