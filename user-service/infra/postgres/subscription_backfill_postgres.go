package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"user/internal/catalogmigration"
	"user/internal/models"
	"user/internal/subscriptionbackfill"
	"user/internal/subscriptiondomain"
)

var (
	_ subscriptionbackfill.Repository             = (*CatalogRepo)(nil)
	_ subscriptionbackfill.LegacyMembershipSource = (*CatalogRepo)(nil)
	_ subscriptionbackfill.EventEvidence          = (*CatalogRepo)(nil)
)

/** WithTransaction executes one atomic subscription migration. */
func (r *CatalogRepo) WithTransaction(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	if r == nil || r.transactionRepo == nil {
		return fmt.Errorf("catalog repository is not configured")
	}
	return r.transactionRepo.WithTransaction(ctx, fn)
}

/** LoadCandidate loads and optionally locks one migration aggregate. */
func (r *CatalogRepo) LoadCandidate(
	ctx context.Context,
	legacyMembershipID uint64,
	forUpdate bool,
) (subscriptionbackfill.Candidate, error) {
	if r == nil || r.transactionRepo == nil {
		return subscriptionbackfill.Candidate{}, fmt.Errorf(
			"catalog repository is not configured",
		)
	}
	if legacyMembershipID == 0 {
		return subscriptionbackfill.Candidate{}, fmt.Errorf(
			"legacy membership id is required",
		)
	}
	db := r.transactionRepo.GetDB(ctx)
	membershipQuery := db.Model(&models.Membership{}).
		Where("id = ? AND deleted_at IS NULL", legacyMembershipID)
	if forUpdate {
		membershipQuery = membershipQuery.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	var membership models.Membership
	if err := membershipQuery.First(&membership).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return subscriptionbackfill.Candidate{}, subscriptionbackfill.ErrNotFound
		}
		return subscriptionbackfill.Candidate{}, err
	}

	mapping := models.LegacyMembershipSubscriptionMapping{
		LegacyMembershipID: membership.ID,
		Status:             models.CatalogMigrationPendingReview,
		Decision:           json.RawMessage(`{}`),
	}
	mappingQuery := db.Model(&models.LegacyMembershipSubscriptionMapping{}).
		Where("legacy_membership_id = ?", membership.ID)
	if forUpdate {
		mappingQuery = mappingQuery.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := mappingQuery.First(&mapping).Error; err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {
		return subscriptionbackfill.Candidate{}, err
	}

	planMappingQuery := db.Model(&models.LegacyPlanCatalogMapping{}).
		Where("legacy_plan_id = ?", membership.PlanID)
	if forUpdate {
		planMappingQuery = planMappingQuery.Clauses(clause.Locking{Strength: "SHARE"})
	}
	var planMapping models.LegacyPlanCatalogMapping
	if err := planMappingQuery.First(&planMapping).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return subscriptionbackfill.Candidate{}, fmt.Errorf(
				"%w: approved legacy plan mapping not found",
				subscriptionbackfill.ErrCandidateBlocked,
			)
		}
		return subscriptionbackfill.Candidate{}, err
	}
	if planMapping.PlanVersionID == nil || *planMapping.PlanVersionID == 0 {
		return subscriptionbackfill.Candidate{}, fmt.Errorf(
			"%w: approved legacy plan mapping is incomplete",
			subscriptionbackfill.ErrCandidateBlocked,
		)
	}

	var version models.CatalogPlanVersion
	if err := db.First(&version, *planMapping.PlanVersionID).Error; err != nil {
		return subscriptionbackfill.Candidate{}, err
	}
	var plan models.CatalogPlan
	if err := db.First(&plan, version.PlanID).Error; err != nil {
		return subscriptionbackfill.Candidate{}, err
	}
	var product models.CatalogProduct
	if err := db.First(&product, plan.ProductID).Error; err != nil {
		return subscriptionbackfill.Candidate{}, err
	}

	candidate := subscriptionbackfill.Candidate{
		Membership: membership,
		Mapping:    toBackfillMapping(mapping),
		Target: subscriptionbackfill.PlanTarget{
			LegacyPlanID:      membership.PlanID,
			ProductID:         product.ID,
			ProductCode:       product.Code,
			PlanVersionID:     version.ID,
			PlanCode:          plan.Code,
			PlanVersion:       version.Version,
			PlanTermsChecksum: version.TermsChecksum,
			PlanStatus:        version.Status,
			PlanPublishedAt:   copyTime(version.PublishedAt),
			SubjectScope:      version.SubjectScope,
			EffectiveFrom:     copyTime(version.EffectiveFrom),
			EffectiveUntil:    copyTime(version.EffectiveUntil),
			MappingStatus:     planMapping.Status,
		},
	}
	if mapping.SubscriptionID != nil && *mapping.SubscriptionID > 0 {
		subscription, err := r.GetSubscription(ctx, *mapping.SubscriptionID)
		if err != nil {
			return subscriptionbackfill.Candidate{}, err
		}
		candidate.ExistingSubscription = &subscription
	}
	return candidate, nil
}

/** SaveReview persists an explicit membership interpretation. */
func (r *CatalogRepo) SaveReview(
	ctx context.Context,
	request subscriptionbackfill.ReviewRequest,
	checksum string,
	decision json.RawMessage,
) error {
	if r == nil || r.transactionRepo == nil {
		return fmt.Errorf("catalog repository is not configured")
	}
	db := r.transactionRepo.GetDB(ctx)
	var existing models.LegacyMembershipSubscriptionMapping
	err := db.Where("legacy_membership_id = ?", request.LegacyMembershipID).
		First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		entity := models.LegacyMembershipSubscriptionMapping{
			LegacyMembershipID: request.LegacyMembershipID,
			SourceChecksum:     checksum,
			Status:             models.CatalogMigrationReviewed,
			Decision:           append(json.RawMessage(nil), decision...),
			ReviewedBy:         copyUint64(&request.ReviewedBy),
			ReviewedAt:         copyTime(&request.ReviewedAt),
			BlockedReason:      "",
			LastRunID:          "",
			LastError:          "",
		}
		return db.Create(&entity).Error
	}
	if err != nil {
		return err
	}
	if existing.SubscriptionID != nil {
		return fmt.Errorf("membership mapping already points to a subscription")
	}
	result := db.Model(&models.LegacyMembershipSubscriptionMapping{}).
		Where(
			"legacy_membership_id = ? AND status IN ? AND subscription_id IS NULL",
			request.LegacyMembershipID,
			[]models.CatalogMigrationStatus{
				models.CatalogMigrationPendingReview,
				models.CatalogMigrationBlocked,
			},
		).
		Updates(map[string]any{
			"source_checksum": checksum,
			"status":          models.CatalogMigrationReviewed,
			"decision":        decision,
			"reviewed_by":     request.ReviewedBy,
			"reviewed_at":     request.ReviewedAt,
			"blocked_reason":  "",
			"last_error":      "",
			"updated_at":      gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return subscriptionbackfill.ErrConcurrentChange
	}
	return nil
}

/** ListReviewedMembershipIDs returns a stable execution batch. */
func (r *CatalogRepo) ListReviewedMembershipIDs(
	ctx context.Context,
	afterID uint64,
	limit int,
) ([]uint64, error) {
	if r == nil || r.transactionRepo == nil {
		return nil, fmt.Errorf("catalog repository is not configured")
	}
	if limit <= 0 || limit > 1000 {
		return nil, fmt.Errorf("batch size must be between 1 and 1000")
	}
	var ids []uint64
	err := r.transactionRepo.GetDB(ctx).
		Model(&models.LegacyMembershipSubscriptionMapping{}).
		Where(
			"legacy_membership_id > ? AND status = ? AND subscription_id IS NULL",
			afterID,
			models.CatalogMigrationReviewed,
		).
		Order("legacy_membership_id ASC").
		Limit(limit).
		Pluck("legacy_membership_id", &ids).Error
	return ids, err
}

/** FindSubscriptionByKey loads one stable subscription identity. */
func (r *CatalogRepo) FindSubscriptionByKey(
	ctx context.Context,
	key string,
) (subscriptiondomain.Aggregate, error) {
	if r == nil || r.transactionRepo == nil {
		return subscriptiondomain.Aggregate{}, fmt.Errorf(
			"catalog repository is not configured",
		)
	}
	var entity models.CatalogSubscription
	err := r.transactionRepo.GetDB(ctx).
		Where("subscription_key = ?", key).
		First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return subscriptiondomain.Aggregate{}, subscriptionbackfill.ErrNotFound
	}
	if err != nil {
		return subscriptiondomain.Aggregate{}, err
	}
	return r.GetSubscription(ctx, entity.ID)
}

/** FindLiveSubscription detects conflicting ownership before insert. */
func (r *CatalogRepo) FindLiveSubscription(
	ctx context.Context,
	subjectKind models.SubscriptionSubjectKind,
	subjectID string,
	productID uint64,
) (subscriptiondomain.Aggregate, error) {
	if r == nil || r.transactionRepo == nil {
		return subscriptiondomain.Aggregate{}, fmt.Errorf(
			"catalog repository is not configured",
		)
	}
	var entity models.CatalogSubscription
	err := r.transactionRepo.GetDB(ctx).
		Where(
			"subject_kind = ? AND subject_id = ? AND product_id = ? AND status IN ?",
			subjectKind,
			subjectID,
			productID,
			[]models.SubscriptionStatus{
				models.SubscriptionPending,
				models.SubscriptionActive,
				models.SubscriptionPastDue,
			},
		).
		First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return subscriptiondomain.Aggregate{}, subscriptionbackfill.ErrNotFound
	}
	if err != nil {
		return subscriptiondomain.Aggregate{}, err
	}
	return r.GetSubscription(ctx, entity.ID)
}

/** CreateSubscription inserts one reviewed legacy subscription. */
func (r *CatalogRepo) CreateSubscription(
	ctx context.Context,
	draft catalogmigration.SubscriptionDraft,
	target subscriptionbackfill.PlanTarget,
	decision subscriptionbackfill.ReviewDecision,
	run subscriptionbackfill.Run,
) (uint64, error) {
	if r == nil || r.transactionRepo == nil {
		return 0, fmt.Errorf("catalog repository is not configured")
	}
	status := models.SubscriptionStatus(draft.Status)
	accessUntil := draft.AccessUntil
	entity := models.CatalogSubscription{
		SubscriptionKey:    draft.SubscriptionKey,
		SubjectKind:        models.SubscriptionSubjectKind(draft.SubjectKind),
		SubjectID:          draft.SubjectID,
		ProductID:          target.ProductID,
		PlanVersionID:      target.PlanVersionID,
		Status:             status,
		StartedAt:          draft.StartedAt,
		CurrentPeriodStart: draft.CurrentPeriodStart,
		CurrentPeriodEnd:   draft.CurrentPeriodEnd,
		AccessUntil:        &accessUntil,
		AutoRenew:          decision.AutoRenew,
		OrderReference:     draft.OrderReference,
		CreatedAt:          draft.CreatedAt,
		UpdatedAt:          draft.UpdatedAt,
	}
	if status == models.SubscriptionExpired {
		endedAt := draft.CurrentPeriodEnd
		entity.EndedAt = &endedAt
	}
	if err := r.transactionRepo.GetDB(ctx).Create(&entity).Error; err != nil {
		return 0, err
	}
	return entity.ID, nil
}

/** EnsureMigrationEvent inserts one idempotent append-only proof. */
func (r *CatalogRepo) EnsureMigrationEvent(
	ctx context.Context,
	event subscriptionbackfill.Event,
) error {
	if r == nil || r.transactionRepo == nil {
		return fmt.Errorf("catalog repository is not configured")
	}
	entity := models.CatalogSubscriptionEvent{
		SubscriptionID:  event.SubscriptionID,
		EventType:       subscriptionbackfill.EventTypeLegacyMembershipBackfilled,
		ToPlanVersionID: copyUint64(&event.PlanVersionID),
		EffectiveAt:     event.EffectiveAt,
		Reason:          "Legacy membership backfilled after explicit review",
		ActorKind:       "migration",
		ActorID:         event.Run.ActorID,
		RequestID:       event.Run.RequestID,
		OperationID:     event.Run.OperationID,
		Metadata:        append(json.RawMessage(nil), event.Metadata...),
	}
	db := r.transactionRepo.GetDB(ctx)
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).
		Create(&entity).Error; err != nil {
		return err
	}
	var stored models.CatalogSubscriptionEvent
	if err := db.Where(
		"subscription_id = ? AND event_type = ?",
		event.SubscriptionID,
		subscriptionbackfill.EventTypeLegacyMembershipBackfilled,
	).First(&stored).Error; err != nil {
		return err
	}
	if !sameLegacyBackfillEvent(stored, entity) {
		return subscriptionbackfill.ErrConflictingTarget
	}
	return nil
}

/** MarkShadowRead activates non-authoritative comparison only. */
func (r *CatalogRepo) MarkShadowRead(
	ctx context.Context,
	legacyMembershipID uint64,
	subscriptionID uint64,
	run subscriptionbackfill.Run,
) error {
	if r == nil || r.transactionRepo == nil {
		return fmt.Errorf("catalog repository is not configured")
	}
	result := r.transactionRepo.GetDB(ctx).
		Model(&models.LegacyMembershipSubscriptionMapping{}).
		Where(
			"legacy_membership_id = ? AND status = ? AND subscription_id IS NULL",
			legacyMembershipID,
			models.CatalogMigrationReviewed,
		).
		Updates(map[string]any{
			"subscription_id": subscriptionID,
			"status":          models.CatalogMigrationShadowRead,
			"migrated_at":     run.ExecutedAt,
			"last_run_id":     run.RunID,
			"blocked_reason":  "",
			"last_error":      "",
			"updated_at":      gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return subscriptionbackfill.ErrConcurrentChange
	}
	return nil
}

/** MarkBlocked preserves review failure without mutating legacy data. */
func (r *CatalogRepo) MarkBlocked(
	ctx context.Context,
	legacyMembershipID uint64,
	run subscriptionbackfill.Run,
	reason string,
) error {
	if r == nil || r.transactionRepo == nil {
		return fmt.Errorf("catalog repository is not configured")
	}
	result := r.transactionRepo.GetDB(ctx).
		Model(&models.LegacyMembershipSubscriptionMapping{}).
		Where(
			"legacy_membership_id = ? AND status IN ?",
			legacyMembershipID,
			[]models.CatalogMigrationStatus{
				models.CatalogMigrationPendingReview,
				models.CatalogMigrationReviewed,
				models.CatalogMigrationShadowRead,
				models.CatalogMigrationBlocked,
			},
		).
		Updates(map[string]any{
			"status":         models.CatalogMigrationBlocked,
			"blocked_reason": reason,
			"last_run_id":    run.RunID,
			"last_error":     reason,
			"updated_at":     gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return subscriptionbackfill.ErrConcurrentChange
	}
	return nil
}

/** ListLegacyMemberships returns a bounded audit page. */
func (r *CatalogRepo) ListLegacyMemberships(
	ctx context.Context,
	afterID uint64,
	limit int,
) ([]models.Membership, error) {
	if r == nil || r.transactionRepo == nil {
		return nil, fmt.Errorf("catalog repository is not configured")
	}
	if limit <= 0 || limit > 1000 {
		return nil, fmt.Errorf("batch size must be between 1 and 1000")
	}
	var memberships []models.Membership
	err := r.transactionRepo.GetDB(ctx).
		Where("id > ? AND deleted_at IS NULL", afterID).
		Order("id ASC").
		Limit(limit).
		Find(&memberships).Error
	return memberships, err
}

/** HasLegacyBackfillEvent checks append-only rollback evidence. */
func (r *CatalogRepo) HasLegacyBackfillEvent(
	ctx context.Context,
	subscriptionID uint64,
	legacyMembershipID uint64,
	sourceChecksum string,
) (bool, error) {
	if r == nil || r.transactionRepo == nil {
		return false, fmt.Errorf("catalog repository is not configured")
	}
	var count int64
	err := r.transactionRepo.GetDB(ctx).
		Model(&models.CatalogSubscriptionEvent{}).
		Where(
			"subscription_id = ? AND event_type = ? "+
				"AND metadata ->> 'legacyMembershipId' = ? "+
				"AND metadata ->> 'sourceChecksum' = ?",
			subscriptionID,
			subscriptionbackfill.EventTypeLegacyMembershipBackfilled,
			fmt.Sprintf("%d", legacyMembershipID),
			sourceChecksum,
		).
		Count(&count).Error
	return count > 0, err
}

type legacyBackfillEventMetadata struct {
	LegacyMembershipID uint64 `json:"legacyMembershipId"`
	LegacyPlanID       uint64 `json:"legacyPlanId"`
	SourceChecksum     string `json:"sourceChecksum"`
	PlanTermsChecksum  string `json:"planTermsChecksum"`
}

func sameLegacyBackfillEvent(
	stored models.CatalogSubscriptionEvent,
	expected models.CatalogSubscriptionEvent,
) bool {
	if stored.SubscriptionID != expected.SubscriptionID ||
		stored.EventType != expected.EventType ||
		stored.ToPlanVersionID == nil || expected.ToPlanVersionID == nil ||
		*stored.ToPlanVersionID != *expected.ToPlanVersionID {
		return false
	}
	var storedMetadata legacyBackfillEventMetadata
	var expectedMetadata legacyBackfillEventMetadata
	if json.Unmarshal(stored.Metadata, &storedMetadata) != nil ||
		json.Unmarshal(expected.Metadata, &expectedMetadata) != nil {
		return false
	}
	return storedMetadata == expectedMetadata &&
		storedMetadata.LegacyMembershipID > 0 &&
		storedMetadata.LegacyPlanID > 0 &&
		storedMetadata.SourceChecksum != "" &&
		storedMetadata.PlanTermsChecksum != ""
}

func toBackfillMapping(
	entity models.LegacyMembershipSubscriptionMapping,
) subscriptionbackfill.Mapping {
	return subscriptionbackfill.Mapping{
		LegacyMembershipID:       entity.LegacyMembershipID,
		SubscriptionID:           copyUint64(entity.SubscriptionID),
		SourceChecksum:           entity.SourceChecksum,
		Status:                   entity.Status,
		Decision:                 append(json.RawMessage(nil), entity.Decision...),
		ReviewedBy:               copyUint64(entity.ReviewedBy),
		ReviewedAt:               copyTime(entity.ReviewedAt),
		MigratedAt:               copyTime(entity.MigratedAt),
		BlockedReason:            entity.BlockedReason,
		LastRunID:                entity.LastRunID,
		LastError:                entity.LastError,
		CutoverEvidenceChecksum:  entity.CutoverEvidenceChecksum,
		CutoverAt:                copyTime(entity.CutoverAt),
		CutoverBy:                entity.CutoverBy,
		RollbackEvidenceChecksum: entity.RollbackEvidenceChecksum,
		RolledBackAt:             copyTime(entity.RolledBackAt),
		RolledBackBy:             entity.RolledBackBy,
	}
}
