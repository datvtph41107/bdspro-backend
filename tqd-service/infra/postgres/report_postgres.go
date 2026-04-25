package postgres

import (
	"context"
	"errors"
	"tqd/internal/domain"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type reportRepo struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) repo.ReportRepository {
	return &reportRepo{db: db}
}

func (r *reportRepo) Create(ctx context.Context, report *domain.Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *reportRepo) Update(ctx context.Context, report *domain.Report) error {
	return r.db.WithContext(ctx).Save(report).Error
}

func (r *reportRepo) UpdateStatus(ctx context.Context, id uint64, status uint32, errorMsg *string) error {
	updates := map[string]interface{}{"status": status}
	if errorMsg != nil {
		updates["error_message"] = *errorMsg
	}
	if status == 30 {
		updates["completed_at"] = gorm.Expr("NOW()")
	}
	return r.db.WithContext(ctx).Model(&domain.Report{}).Where("id = ?", id).Updates(updates).Error
}

func (r *reportRepo) UpdateFileInfo(ctx context.Context, id uint64, fileURL *string, fileSize *int64, fileHash *string, completedAt *string) error {
	updates := map[string]interface{}{
		"file_url":  fileURL,
		"file_size": fileSize,
		"file_hash": fileHash,
		"status":    30,
	}
	if completedAt != nil {
		updates["completed_at"] = *completedAt
	} else {
		updates["completed_at"] = gorm.Expr("NOW()")
	}
	return r.db.WithContext(ctx).Model(&domain.Report{}).Where("id = ?", id).Updates(updates).Error
}

func (r *reportRepo) GetByID(ctx context.Context, id uint64) (*domain.Report, error) {
	var report domain.Report
	err := r.db.WithContext(ctx).First(&report, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &report, err
}

func (r *reportRepo) ListByUser(ctx context.Context, userID uint64, reportType *uint32, status *uint32, page, limit int) ([]domain.Report, int64, error) {
	var reports []domain.Report
	var total int64
	query := r.db.WithContext(ctx).Model(&domain.Report{}).Where("user_id = ? AND deleted_at IS NULL", userID)
	if reportType != nil {
		query = query.Where("report_type = ?", *reportType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&reports).Error
	return reports, total, err
}

func (r *reportRepo) AdminList(ctx context.Context, userID *uint64, reportType *uint32, status *uint32, page, limit int) ([]domain.Report, int64, error) {
	var reports []domain.Report
	var total int64
	query := r.db.WithContext(ctx).Model(&domain.Report{}).Where("deleted_at IS NULL")
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if reportType != nil {
		query = query.Where("report_type = ?", *reportType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&reports).Error
	return reports, total, err
}

func (r *reportRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.Report{}, id).Error
}

func (r *reportRepo) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at IS NOT NULL AND expires_at < NOW()").Delete(&domain.Report{}).Error
}
