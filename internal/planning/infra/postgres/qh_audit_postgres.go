package postgres

import (
	"context"
	"time"

	"tqd/internal/usecase/resolver/layer_resolver/rule"
	"tqd/internal/usecase/resolver/layer_resolver/types"

	"gorm.io/gorm"
)

// QHAuditPostgres implements rule.AuditRepository
type QHAuditPostgres struct {
	DB *gorm.DB
}

func NewQHAuditPostgres(db *gorm.DB) rule.AuditRepository {
	return &QHAuditPostgres{DB: db}
}

func (r *QHAuditPostgres) Insert(ctx context.Context, entry *types.AuditEntry) error {
	return r.DB.WithContext(ctx).Table("qh_audit_entries").Create(entry).Error
}

func (r *QHAuditPostgres) Query(ctx context.Context, q *types.AuditQuery) ([]*types.AuditEntry, int64, error) {
	db := r.DB.WithContext(ctx).Table("qh_audit_entries")

	if q.EntityType != "" {
		db = db.Where("entity_type = ?", q.EntityType)
	}
	if q.EntityID > 0 {
		db = db.Where("entity_id = ?", q.EntityID)
	}
	if q.Action != "" {
		db = db.Where("action = ?", q.Action)
	}
	if q.Actor != "" {
		db = db.Where("actor = ?", q.Actor)
	}
	if !q.From.IsZero() {
		db = db.Where("timestamp >= ?", q.From)
	}
	if !q.To.IsZero() {
		db = db.Where("timestamp <= ?", q.To)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	offset := (q.Page - 1) * q.PageSize

	var entries []*types.AuditEntry
	if err := db.Order("timestamp DESC").Limit(q.PageSize).Offset(offset).Find(&entries).Error; err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

func (r *QHAuditPostgres) QueryByParcel(ctx context.Context, parcelID uint64) ([]*types.AuditEntry, error) {
	var entries []*types.AuditEntry
	err := r.DB.WithContext(ctx).Table("qh_audit_entries").
		Where("entity_type = ? AND entity_id = ?", "parcel", parcelID).
		Order("timestamp DESC").
		Limit(100).
		Find(&entries).Error
	return entries, err
}

var _ = time.Now
