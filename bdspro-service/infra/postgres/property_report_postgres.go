package postgres

import (
	"bdspro/internal/domain"
	property_repo "bdspro/internal/repo/property"
	"context"
	"errors"

	"gorm.io/gorm"
)

type propertyReportRepo struct {
	db *gorm.DB
}

func NewPropertyReportRepository(db *gorm.DB) property_repo.PropertyReportRepository {
	return &propertyReportRepo{db: db}
}

func (r *propertyReportRepo) Create(ctx context.Context, report *domain.PropertyReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *propertyReportRepo) GetByID(ctx context.Context, id uint64) (*domain.PropertyReport, error) {
	var report domain.PropertyReport
	err := r.db.WithContext(ctx).First(&report, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &report, err
}

func (r *propertyReportRepo) GetByLineageAndReporter(
	ctx context.Context,
	lineageID, reporterOriginID uint64,
) (*domain.PropertyReport, error) {
	var report domain.PropertyReport
	err := r.db.WithContext(ctx).
		Where("lineage_id = ? AND reporter_origin_id = ? AND status IN (10, 20)", lineageID, reporterOriginID).
		First(&report).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &report, err
}

func (r *propertyReportRepo) UpdateFields(ctx context.Context, id uint64, fields map[string]any) error {
	return r.db.WithContext(ctx).
		Model(&domain.PropertyReport{}).
		Where("id = ?", id).
		Updates(fields).Error
}

func (r *propertyReportRepo) CountPendingByLineage(ctx context.Context, lineageID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.PropertyReport{}).
		Where("lineage_id = ? AND status = 10", lineageID). // 10 = pending
		Count(&count).Error
	return count, err
}
