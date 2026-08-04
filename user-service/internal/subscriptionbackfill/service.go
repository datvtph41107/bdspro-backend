package subscriptionbackfill

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"common/productcatalog"
	"user/internal/catalogmigration"
	"user/internal/models"
	"user/internal/subscriptiondomain"
)

/** Service reviews and backfills subscriptions without cutting over reads. */
type Service struct {
	repository Repository
}

/** NewService creates the explicit migration service. */
func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

/** Review stores the approved interpretation of one legacy row. */
func (s *Service) Review(ctx context.Context, request ReviewRequest) error {
	if s == nil || s.repository == nil {
		return fmt.Errorf("subscription backfill repository is not configured")
	}
	if request.LegacyMembershipID == 0 || request.ReviewedBy == 0 || request.ReviewedAt.IsZero() {
		return fmt.Errorf("review identity is incomplete")
	}
	if err := request.Decision.Validate(); err != nil {
		return err
	}
	decision, err := json.Marshal(request.Decision)
	if err != nil {
		return err
	}
	return s.repository.WithTransaction(ctx, func(txCtx context.Context) error {
		candidate, err := s.repository.LoadCandidate(txCtx, request.LegacyMembershipID, true)
		if err != nil {
			return err
		}
		if candidate.Mapping.Status == models.CatalogMigrationShadowRead ||
			candidate.Mapping.Status == models.CatalogMigrationCutover {
			return fmt.Errorf("membership has already been migrated")
		}
		if err := validateTarget(candidate.Target, candidate.Membership); err != nil {
			return err
		}
		checksum, err := catalogmigration.MembershipSourceChecksum(
			toLegacySnapshot(candidate.Membership),
		)
		if err != nil {
			return err
		}
		if _, err := catalogmigration.BuildSubscriptionDraft(
			toLegacySnapshot(candidate.Membership),
			catalogmigration.MembershipDecision{
				ProductCode:       candidate.Target.ProductCode,
				PlanCode:          candidate.Target.PlanCode,
				PlanVersion:       candidate.Target.PlanVersion,
				PlanTermsChecksum: candidate.Target.PlanTermsChecksum,
				PeriodHandling:    request.Decision.PeriodHandling,
			},
			request.ReviewedAt,
		); err != nil {
			return fmt.Errorf("membership cannot be reviewed: %w", err)
		}
		return s.repository.SaveReview(txCtx, request, checksum, decision)
	})
}

/** Plan validates one reviewed candidate without writing. */
func (s *Service) Plan(ctx context.Context, legacyMembershipID uint64, asOf time.Time) (Result, error) {
	if s == nil || s.repository == nil {
		return Result{}, fmt.Errorf("subscription backfill repository is not configured")
	}
	candidate, err := s.repository.LoadCandidate(ctx, legacyMembershipID, false)
	if err != nil {
		if errors.Is(err, ErrCandidateBlocked) {
			return Result{
				LegacyMembershipID: legacyMembershipID,
				Status:             ResultBlocked,
				Detail:             err.Error(),
			}, nil
		}
		return Result{}, err
	}
	draft, _, err := prepare(
		candidate,
		asOf,
		models.CatalogMigrationReviewed,
	)
	if err != nil {
		return Result{
			LegacyMembershipID: legacyMembershipID,
			Status:             ResultBlocked,
			Detail:             err.Error(),
		}, nil
	}
	return Result{
		LegacyMembershipID: legacyMembershipID,
		Status:             ResultPlanned,
		Detail: fmt.Sprintf(
			"subscription_key=%s status=%s plan_version_id=%d",
			draft.SubscriptionKey,
			draft.Status,
			candidate.Target.PlanVersionID,
		),
	}, nil
}

/** Execute migrates one reviewed membership transactionally. */
func (s *Service) Execute(ctx context.Context, legacyMembershipID uint64, run Run) (Result, error) {
	if s == nil || s.repository == nil {
		return Result{}, fmt.Errorf("subscription backfill repository is not configured")
	}
	if legacyMembershipID == 0 {
		return Result{}, fmt.Errorf("legacy membership id is required")
	}
	if err := run.Validate(); err != nil {
		return Result{}, err
	}

	result := Result{LegacyMembershipID: legacyMembershipID}
	err := s.repository.WithTransaction(ctx, func(txCtx context.Context) error {
		candidate, err := s.repository.LoadCandidate(txCtx, legacyMembershipID, true)
		if err != nil {
			return err
		}
		if candidate.Mapping.Status == models.CatalogMigrationCutover {
			return fmt.Errorf("cutover mapping is read-only")
		}
		if candidate.Mapping.Status == models.CatalogMigrationShadowRead {
			if candidate.Mapping.SubscriptionID == nil {
				return blocking("migrated mapping has no subscription")
			}
			draft, decision, err := prepare(
				candidate,
				run.AsOf,
				models.CatalogMigrationShadowRead,
			)
			if err != nil {
				return err
			}
			subscription, err := s.repository.FindSubscriptionByKey(
				txCtx,
				draft.SubscriptionKey,
			)
			if err != nil {
				return err
			}
			if subscription.ID != *candidate.Mapping.SubscriptionID {
				return blocking("mapping points to a different subscription")
			}
			if err := matchSubscription(subscription, draft, candidate.Target, decision); err != nil {
				return blocking(err.Error())
			}
			result.SubscriptionID = subscription.ID
			result.Status = ResultIdempotent
			return nil
		}

		draft, decision, err := prepare(
			candidate,
			run.AsOf,
			models.CatalogMigrationReviewed,
		)
		if err != nil {
			return err
		}

		subscriptionID := uint64(0)
		existing, err := s.repository.FindSubscriptionByKey(txCtx, draft.SubscriptionKey)
		switch {
		case err == nil:
			if err := matchSubscription(existing, draft, candidate.Target, decision); err != nil {
				return blocking(err.Error())
			}
			subscriptionID = existing.ID
		case errors.Is(err, ErrNotFound):
			live, liveErr := s.repository.FindLiveSubscription(
				txCtx,
				models.SubscriptionSubjectProfile,
				draft.SubjectID,
				candidate.Target.ProductID,
			)
			switch {
			case liveErr == nil:
				return blocking(fmt.Sprintf(
					"subject already has live subscription %s",
					live.SubscriptionKey,
				))
			case !errors.Is(liveErr, ErrNotFound):
				return liveErr
			}
			subscriptionID, err = s.repository.CreateSubscription(
				txCtx,
				draft,
				candidate.Target,
				decision,
				run,
			)
			if err != nil {
				return err
			}
		default:
			return err
		}

		metadata, err := json.Marshal(map[string]any{
			"legacyMembershipId": legacyMembershipID,
			"legacyPlanId":       candidate.Membership.PlanID,
			"sourceChecksum":     draft.SourceChecksum,
			"planTermsChecksum":  draft.PlanTermsChecksum,
			"runId":              run.RunID,
		})
		if err != nil {
			return err
		}
		if err := s.repository.EnsureMigrationEvent(txCtx, Event{
			SubscriptionID: subscriptionID,
			PlanVersionID:  candidate.Target.PlanVersionID,
			EffectiveAt:    run.ExecutedAt,
			Run:            run,
			Metadata:       metadata,
		}); err != nil {
			return err
		}
		if err := s.repository.MarkShadowRead(
			txCtx,
			legacyMembershipID,
			subscriptionID,
			run,
		); err != nil {
			return err
		}
		result.SubscriptionID = subscriptionID
		result.Status = ResultMigrated
		return nil
	})
	if err == nil {
		return result, nil
	}
	if errors.Is(err, ErrCandidateBlocked) {
		reason := strings.TrimSpace(strings.TrimPrefix(err.Error(), ErrCandidateBlocked.Error()+": "))
		if reason == "" {
			reason = "candidate validation failed"
		}
		if blockErr := s.repository.MarkBlocked(ctx, legacyMembershipID, run, reason); blockErr != nil {
			return Result{}, fmt.Errorf("%w; record blocked failed: %v", err, blockErr)
		}
		return Result{
			LegacyMembershipID: legacyMembershipID,
			Status:             ResultBlocked,
			Detail:             reason,
		}, nil
	}
	return Result{
		LegacyMembershipID: legacyMembershipID,
		Status:             ResultFailed,
		Detail:             err.Error(),
	}, err
}

func prepare(
	candidate Candidate,
	asOf time.Time,
	allowedStatuses ...models.CatalogMigrationStatus,
) (catalogmigration.SubscriptionDraft, ReviewDecision, error) {
	statusAllowed := false
	for _, status := range allowedStatuses {
		if candidate.Mapping.Status == status {
			statusAllowed = true
			break
		}
	}
	if !statusAllowed {
		return catalogmigration.SubscriptionDraft{}, ReviewDecision{}, blocking(
			"membership mapping is not in an executable state",
		)
	}
	if candidate.Mapping.ReviewedBy == nil || *candidate.Mapping.ReviewedBy == 0 ||
		candidate.Mapping.ReviewedAt == nil {
		return catalogmigration.SubscriptionDraft{}, ReviewDecision{}, blocking(
			"review evidence is incomplete",
		)
	}
	var decision ReviewDecision
	if err := json.Unmarshal(candidate.Mapping.Decision, &decision); err != nil {
		return catalogmigration.SubscriptionDraft{}, ReviewDecision{}, blocking(
			"review decision is invalid",
		)
	}
	if err := decision.Validate(); err != nil {
		return catalogmigration.SubscriptionDraft{}, ReviewDecision{}, blocking(err.Error())
	}
	if err := validateTarget(candidate.Target, candidate.Membership); err != nil {
		return catalogmigration.SubscriptionDraft{}, ReviewDecision{}, blocking(err.Error())
	}
	checksum, err := catalogmigration.MembershipSourceChecksum(
		toLegacySnapshot(candidate.Membership),
	)
	if err != nil {
		return catalogmigration.SubscriptionDraft{}, ReviewDecision{}, blocking(err.Error())
	}
	if checksum != candidate.Mapping.SourceChecksum {
		return catalogmigration.SubscriptionDraft{}, ReviewDecision{}, blocking(
			"legacy membership changed after review",
		)
	}
	draft, err := catalogmigration.BuildSubscriptionDraft(
		toLegacySnapshot(candidate.Membership),
		catalogmigration.MembershipDecision{
			ProductCode:       candidate.Target.ProductCode,
			PlanCode:          candidate.Target.PlanCode,
			PlanVersion:       candidate.Target.PlanVersion,
			PlanTermsChecksum: candidate.Target.PlanTermsChecksum,
			PeriodHandling:    decision.PeriodHandling,
		},
		asOf,
	)
	if err != nil {
		return catalogmigration.SubscriptionDraft{}, ReviewDecision{}, blocking(err.Error())
	}
	return draft, decision, nil
}

func validateTarget(target PlanTarget, membership models.Membership) error {
	if target.LegacyPlanID == 0 || target.LegacyPlanID != membership.PlanID {
		return fmt.Errorf("approved plan mapping does not match membership")
	}
	if target.ProductID == 0 || target.PlanVersionID == 0 ||
		strings.TrimSpace(target.ProductCode) == "" ||
		strings.TrimSpace(target.PlanCode) == "" ||
		strings.TrimSpace(target.PlanVersion) == "" {
		return fmt.Errorf("approved plan target is incomplete")
	}
	if target.MappingStatus != models.CatalogMigrationShadowRead &&
		target.MappingStatus != models.CatalogMigrationCutover {
		return fmt.Errorf("approved plan target is not shadow-ready")
	}
	if target.PlanStatus != productcatalog.StatusActive || target.PlanPublishedAt == nil {
		return fmt.Errorf("approved plan version is not published and active")
	}
	if target.SubjectScope != productcatalog.SubjectScopeProfile &&
		target.SubjectScope != productcatalog.SubjectScopeAny {
		return fmt.Errorf("approved plan version does not support profile subscriptions")
	}
	if membership.PaidAt != nil && target.EffectiveFrom != nil &&
		membership.PaidAt.Before(*target.EffectiveFrom) {
		return fmt.Errorf("membership started before approved plan version")
	}
	if membership.PaidAt != nil && target.EffectiveUntil != nil &&
		!membership.PaidAt.Before(*target.EffectiveUntil) {
		return fmt.Errorf("membership started after approved plan version closed")
	}
	return nil
}

func matchSubscription(
	existing subscriptiondomain.Aggregate,
	draft catalogmigration.SubscriptionDraft,
	target PlanTarget,
	decision ReviewDecision,
) error {
	if err := existing.Validate(); err != nil {
		return fmt.Errorf("existing subscription is invalid: %w", err)
	}
	accessUntilMatches := existing.AccessUntil != nil &&
		existing.AccessUntil.Equal(draft.AccessUntil)
	if existing.SubscriptionKey != draft.SubscriptionKey ||
		existing.SubjectKind != models.SubscriptionSubjectProfile ||
		existing.SubjectID != draft.SubjectID ||
		existing.ProductID != target.ProductID ||
		existing.PlanVersionID != target.PlanVersionID ||
		string(existing.Status) != draft.Status ||
		!existing.StartedAt.Equal(draft.StartedAt) ||
		!existing.CurrentPeriodStart.Equal(draft.CurrentPeriodStart) ||
		!existing.CurrentPeriodEnd.Equal(draft.CurrentPeriodEnd) ||
		!accessUntilMatches ||
		existing.AutoRenew != decision.AutoRenew ||
		existing.OrderReference != draft.OrderReference ||
		!existing.CreatedAt.UTC().Equal(draft.CreatedAt.UTC()) ||
		!existing.UpdatedAt.UTC().Equal(draft.UpdatedAt.UTC()) {
		return ErrConflictingTarget
	}
	return nil
}

func toLegacySnapshot(membership models.Membership) catalogmigration.LegacyMembershipSnapshot {
	return catalogmigration.LegacyMembershipSnapshot{
		ID:           membership.ID,
		ProfileID:    membership.ProfileID,
		PlanID:       membership.PlanID,
		FinishedDate: copyTime(membership.FinishedDate),
		OrderID:      membership.OrderID,
		PaidAt:       copyTime(membership.PaidAt),
		CreatedAt:    copyTime(membership.CreatedAt),
		UpdatedAt:    copyTime(membership.UpdatedAt),
	}
}

func copyTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}

func blocking(reason string) error {
	return fmt.Errorf("%w: %s", ErrCandidateBlocked, strings.TrimSpace(reason))
}
