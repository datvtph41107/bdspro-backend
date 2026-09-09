package settlement

import (
	"time"

	subscriptiondomain "user/internal/domain/subscription"
	"user/internal/models"
)

func Decide(effect Effect, target Target, current Current, now time.Time) (Decision, error) {
	if err := effect.Validate(); err != nil || now.IsZero() {
		return Decision{}, ErrInvalidEffect
	}
	if !target.Matches(effect) {
		return Decision{}, ErrCatalogTermsConflict
	}
	if !current.Found {
		end := effect.OccurredAt.UTC().AddDate(0, 0, int(effect.SubscriptionTermDays))
		aggregate := subscriptiondomain.Aggregate{
			SubscriptionKey:    effect.SubscriptionKey(),
			SubjectKind:        effect.SubjectKind,
			SubjectID:          effect.SubjectID,
			ProductID:          target.ProductID,
			ProductCode:        effect.ProductCode,
			PlanVersionID:      effect.PlanVersionID,
			PlanCode:           effect.PlanCode,
			PlanVersion:        effect.PlanVersion,
			PlanTermsChecksum:  effect.TermsChecksum,
			PlanStatus:         target.PlanStatus,
			Status:             models.SubscriptionActive,
			StartedAt:          effect.OccurredAt.UTC(),
			CurrentPeriodStart: effect.OccurredAt.UTC(),
			CurrentPeriodEnd:   end,
			AccessUntil:        &end,
			OrderReference:     effect.EffectKey,
			CreatedAt:          now.UTC(),
			UpdatedAt:          now.UTC(),
		}
		return Decision{Action: ActionActivated, Subscription: aggregate, ToPlanVersionID: effect.PlanVersionID}, nil
	}
	if current.HasPendingChange || current.Aggregate.Status != models.SubscriptionActive {
		return Decision{}, ErrCurrentState
	}
	if current.TierRank <= 0 {
		return Decision{}, ErrCurrentState
	}
	from := current.Aggregate.PlanVersionID
	switch {
	case effect.TierRank == current.TierRank:
		return Decision{Action: ActionRejectedSameTier, Subscription: current.Aggregate.Copy(), FromPlanVersionID: &from, ToPlanVersionID: effect.PlanVersionID}, nil
	case effect.TierRank < current.TierRank:
		return Decision{Action: ActionRejectedDowngrade, Subscription: current.Aggregate.Copy(), FromPlanVersionID: &from, ToPlanVersionID: effect.PlanVersionID}, nil
	default:
		upgraded := current.Aggregate.Copy()
		upgraded.PlanVersionID = effect.PlanVersionID
		upgraded.PlanCode = effect.PlanCode
		upgraded.PlanVersion = effect.PlanVersion
		upgraded.PlanTermsChecksum = effect.TermsChecksum
		upgraded.PlanStatus = target.PlanStatus
		upgraded.OrderReference = effect.EffectKey
		upgraded.UpdatedAt = now.UTC()
		// Upgrade changes commercial authority immediately but preserves the
		// already-running subscription period. Usage/quota therefore do not reset.
		return Decision{Action: ActionUpgraded, Subscription: upgraded, FromPlanVersionID: &from, ToPlanVersionID: effect.PlanVersionID}, nil
	}
}
