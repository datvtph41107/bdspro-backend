package postgre

import (
	"context"
	"social/internal/domain"
	"social/internal/dto"
	"time"

	"gorm.io/gorm"
)

// @bind: social/internal/repo.ReportReasonRepo
type PostgreReportReason struct {
	db *gorm.DB
}

func NewPostgreReportReason(db *gorm.DB) *PostgreReportReason {
	return &PostgreReportReason{db: db}
}

func (r *PostgreReportReason) GetList(ctx context.Context, dto *dto.ReportReasonRequest) ([]domain.ReportReason, int64, error) {
	var reportReasons []domain.ReportReason
	var total int64

	query := r.db.WithContext(ctx).
		Model(&domain.ReportReason{}).
		Where("deleted_at IS NULL")

	if dto.ReasonName != "" {
		query = query.Where("name LIKE ?", "%"+dto.ReasonName+"%")
	}

	err := query.
		Order("created_at DESC").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&reportReasons).Error
	if err != nil {
		return nil, 0, err
	}

	query.Count(&total)

	return reportReasons, total, nil
}

func (r *PostgreReportReason) Create(ctx context.Context, reportReason *domain.ReportReason) (*domain.ReportReason, error) {
	err := r.db.WithContext(ctx).Create(reportReason).Error
	if err != nil {
		return nil, err
	}
	return reportReason, nil
}

func (r *PostgreReportReason) Update(ctx context.Context, reportReason *domain.ReportReason) (*domain.ReportReason, error) {
	err := r.db.WithContext(ctx).
		Model(&domain.ReportReason{}).
		Where("id = ?", reportReason.ID).
		Updates(reportReason).Error
	if err != nil {
		return nil, err
	}
	return reportReason, nil
}

func (r *PostgreReportReason) Delete(ctx context.Context, id uint64) error {
	err := r.db.WithContext(ctx).
		Model(&domain.ReportReason{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgreReportReason) GetPublic(ctx context.Context) ([]domain.ReportReason, error) {
	var reportReasons []domain.ReportReason
	err := r.db.WithContext(ctx).
		Model(&domain.ReportReason{}).
		Where("deleted_at IS NULL").
		Find(&reportReasons).Error
	if err != nil {
		return nil, err
	}
	return reportReasons, nil
}

func (r *PostgreReportReason) GetByID(ctx context.Context, id uint64) (*domain.ReportReason, error) {
	var reportReason domain.ReportReason
	err := r.db.WithContext(ctx).
		Model(&domain.ReportReason{}).
		Where("id = ?", id).
		First(&reportReason).Error
	if err != nil {
		return nil, err
	}
	return &reportReason, nil
}

func (r *PostgreReportReason) ExistedByID(ctx context.Context, id uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.ReportReason{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
