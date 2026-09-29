package postgre

import (
	_dto "common/domain/dto"
	"context"

	"social/internal/domain"
	"social/internal/enums"

	"gorm.io/gorm"
)

// @bind: social/internal/repo.ReportRepo
type PostgreReport struct {
	db *gorm.DB
}

func NewPostgreReport(db *gorm.DB) *PostgreReport {
	return &PostgreReport{db: db}
}

func (r *PostgreReport) Create(ctx context.Context, report *domain.Report) (*domain.Report, error) {
	if err := r.db.WithContext(ctx).Create(report).Error; err != nil {
		return nil, err
	}
	return report, nil
}

func (r *PostgreReport) ExistByTargetIdAndUserIdAndTargetType(ctx context.Context, targetId uint64, userId uint64, targetType enums.TargetType) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Report{}).
		Where("target_id = ? AND user_id = ? AND target_type = ? AND deleted_at IS NULL", targetId, userId, targetType).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgreReport) GetByCreatedBy(ctx context.Context, createdBy uint64, dto *_dto.Pagable) ([]domain.Report, int64, error) {
	var reports []domain.Report
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.Report{}).
		Where("created_by = ? AND deleted_at IS NULL", createdBy).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Model(&domain.Report{}).
		Where("created_by = ? AND deleted_at IS NULL", createdBy).
		Offset(dto.GetOffset()).
		Limit(int(dto.GetSize())).
		Find(&reports).Error; err != nil {
		return nil, 0, err
	}
	return reports, total, nil
}

func (r *PostgreReport) UpdateStatus(ctx context.Context, id uint64, status enums.ReportStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Report{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("status", status).Error
}
