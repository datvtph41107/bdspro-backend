package postgres

import (
	"context"
	"time"

	"tqd/internal/usecase/resolver/layer_resolver/rule"

	"gorm.io/gorm"
)

// QHLayerLifecyclePostgres implements rule.LayerRepository
type QHLayerLifecyclePostgres struct {
	DB *gorm.DB
}

func NewQHLayerLifecyclePostgres(db *gorm.DB) rule.LayerRepository {
	return &QHLayerLifecyclePostgres{DB: db}
}

func (r *QHLayerLifecyclePostgres) GetByID(ctx context.Context, layerID uint64) (*rule.LayerRecord, error) {
	var rec rule.LayerRecord
	err := r.DB.WithContext(ctx).Table("qh_layers").
		Select("id, name, family_id, parent_id, replaced_by_id, version, state, created_at, created_by").
		Where("id = ?", layerID).
		First(&rec).Error
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *QHLayerLifecyclePostgres) GetByFamilyID(ctx context.Context, familyID uint64) ([]*rule.LayerRecord, error) {
	var recs []*rule.LayerRecord
	err := r.DB.WithContext(ctx).Table("qh_layers").
		Select("id, name, family_id, parent_id, replaced_by_id, version, state, created_at, created_by").
		Where("family_id = ?", familyID).
		Order("version ASC").
		Find(&recs).Error
	return recs, err
}

func (r *QHLayerLifecyclePostgres) GetStateTransitions(ctx context.Context, layerID uint64) ([]*rule.StateTransitionRecord, error) {
	var recs []*rule.StateTransitionRecord
	err := r.DB.WithContext(ctx).Table("qh_lifecycle_transitions").
		Select("id, layer_id, from_state, to_state, actor, reason, legal_doc, created_at").
		Where("layer_id = ?", layerID).
		Order("created_at ASC").
		Find(&recs).Error
	return recs, err
}

var _ = time.Now
