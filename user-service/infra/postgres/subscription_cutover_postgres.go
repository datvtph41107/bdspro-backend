package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"user/internal/models"
	"user/internal/modules/subscription/cutover"
)

var _ cutover.Repository = (*CatalogRepo)(nil)

/** EnsureCutoverEvent inserts one idempotent activation proof. */
func (r *CatalogRepo) EnsureCutoverEvent(
	ctx context.Context,
	event cutover.Event,
) error {
	return r.ensureSubscriptionOperationalEvent(ctx, event)
}

/** EnsureRollbackEvent inserts one idempotent rollback proof. */
func (r *CatalogRepo) EnsureRollbackEvent(
	ctx context.Context,
	event cutover.Event,
) error {
	return r.ensureSubscriptionOperationalEvent(ctx, event)
}

func (r *CatalogRepo) ensureSubscriptionOperationalEvent(
	ctx context.Context,
	event cutover.Event,
) error {
	if r == nil || r.transactionRepo == nil {
		return fmt.Errorf("catalog repository is not configured")
	}
	if event.SubscriptionID == 0 || event.PlanVersionID == 0 ||
		event.EffectiveAt.IsZero() {
		return fmt.Errorf("cutover event identity is incomplete")
	}
	entity := models.CatalogSubscriptionEvent{
		SubscriptionID:  event.SubscriptionID,
		EventType:       event.EventType,
		ToPlanVersionID: copyUint64(&event.PlanVersionID),
		EffectiveAt:     event.EffectiveAt,
		Reason:          event.Run.Reason,
		ActorKind:       "migration",
		ActorID:         event.Run.ActorID,
		RequestID:       event.Run.RequestID,
		OperationID:     event.Run.OperationID,
		Metadata:        append(json.RawMessage(nil), event.Metadata...),
	}
	db := r.transactionRepo.GetDB(ctx)
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&entity).Error; err != nil {
		return err
	}
	var stored models.CatalogSubscriptionEvent
	if err := db.Where(
		"subscription_id = ? AND event_type = ?",
		event.SubscriptionID,
		event.EventType,
	).First(&stored).Error; err != nil {
		return err
	}
	if !sameOperationalEvent(stored, entity) {
		return cutover.ErrConcurrentChange
	}
	return nil
}

/** MarkCutover activates catalog reads for one manifest row. */
func (r *CatalogRepo) MarkCutover(
	ctx context.Context,
	entry cutover.Entry,
	evidenceChecksum string,
	run cutover.Run,
) error {
	if r == nil || r.transactionRepo == nil {
		return fmt.Errorf("catalog repository is not configured")
	}
	result := r.transactionRepo.GetDB(ctx).
		Model(&models.LegacyMembershipSubscriptionMapping{}).
		Where(
			"legacy_membership_id = ? AND subscription_id = ? AND source_checksum = ? "+
				"AND status = ? AND cutover_evidence_checksum = ''",
			entry.LegacyMembershipID,
			entry.SubscriptionID,
			entry.SourceChecksum,
			models.CatalogMigrationShadowRead,
		).
		Updates(map[string]any{
			"status":                    models.CatalogMigrationCutover,
			"cutover_evidence_checksum": evidenceChecksum,
			"cutover_at":                run.ExecutedAt,
			"cutover_by":                run.ActorID,
			"last_run_id":               run.RunID,
			"last_error":                "",
			"updated_at":                gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return cutover.ErrConcurrentChange
	}
	return nil
}

/** MarkRolledBack makes a cutover mapping terminal after rollback. */
func (r *CatalogRepo) MarkRolledBack(
	ctx context.Context,
	entry cutover.Entry,
	evidenceChecksum string,
	run cutover.Run,
) error {
	if r == nil || r.transactionRepo == nil {
		return fmt.Errorf("catalog repository is not configured")
	}
	result := r.transactionRepo.GetDB(ctx).
		Model(&models.LegacyMembershipSubscriptionMapping{}).
		Where(
			"legacy_membership_id = ? AND subscription_id = ? AND source_checksum = ? "+
				"AND status = ? AND rollback_evidence_checksum = ''",
			entry.LegacyMembershipID,
			entry.SubscriptionID,
			entry.SourceChecksum,
			models.CatalogMigrationCutover,
		).
		Updates(map[string]any{
			"status":                     models.CatalogMigrationRolledBack,
			"rollback_evidence_checksum": evidenceChecksum,
			"rolled_back_at":             run.ExecutedAt,
			"rolled_back_by":             run.ActorID,
			"last_run_id":                run.RunID,
			"last_error":                 "",
			"updated_at":                 gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return cutover.ErrConcurrentChange
	}
	return nil
}

func sameOperationalEvent(
	stored models.CatalogSubscriptionEvent,
	expected models.CatalogSubscriptionEvent,
) bool {
	if stored.SubscriptionID != expected.SubscriptionID ||
		stored.EventType != expected.EventType ||
		stored.ToPlanVersionID == nil || expected.ToPlanVersionID == nil ||
		*stored.ToPlanVersionID != *expected.ToPlanVersionID {
		return false
	}
	var storedMetadata map[string]any
	var expectedMetadata map[string]any
	if json.Unmarshal(stored.Metadata, &storedMetadata) != nil ||
		json.Unmarshal(expected.Metadata, &expectedMetadata) != nil {
		return false
	}
	storedJSON, storedErr := json.Marshal(storedMetadata)
	expectedJSON, expectedErr := json.Marshal(expectedMetadata)
	return storedErr == nil && expectedErr == nil && string(storedJSON) == string(expectedJSON)
}

/** GetCutoverEvent returns event evidence for operational verification. */
func (r *CatalogRepo) GetCutoverEvent(
	ctx context.Context,
	subscriptionID uint64,
	eventType string,
) (models.CatalogSubscriptionEvent, error) {
	if r == nil || r.transactionRepo == nil {
		return models.CatalogSubscriptionEvent{}, fmt.Errorf("catalog repository is not configured")
	}
	var entity models.CatalogSubscriptionEvent
	err := r.transactionRepo.GetDB(ctx).Where(
		"subscription_id = ? AND event_type = ?",
		subscriptionID,
		eventType,
	).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.CatalogSubscriptionEvent{}, cutover.ErrEvidenceInvalid
	}
	return entity, err
}
