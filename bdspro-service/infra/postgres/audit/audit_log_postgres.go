package postgres_audit

import (
	"bdspro/internal/domain"
	property_repo "bdspro/internal/repo/property"
	"context"

	"gorm.io/gorm"
)

type auditLogRepoImpl struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) property_repo.AuditLogRepository {
	return &auditLogRepoImpl{db: db}
}

func (r *auditLogRepoImpl) Create(ctx context.Context, log *domain.PropertyAuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *auditLogRepoImpl) GetByProperty(ctx context.Context, propertyID uint64, limit int) ([]*domain.PropertyAuditLog, error) {
	var logs []*domain.PropertyAuditLog
	query := r.db.WithContext(ctx).Where("property_id = ?", propertyID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&logs).Error
	return logs, err
}
