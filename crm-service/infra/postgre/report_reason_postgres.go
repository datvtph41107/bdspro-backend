package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/repo"
	"errors"

	"gorm.io/gorm"
)

type ReportReasonPostgres struct {
	db *gorm.DB
}

// @bind: crm/internal/repo.ReportReasonRepo
func NewReportReasonPostgres(db *gorm.DB) repo.ReportReasonRepo {
	return &ReportReasonPostgres{db: db}
}

// GetAll gets all report reasons
func (r *ReportReasonPostgres) GetAll(ctx context.Context) ([]domain.ReportReason, error) {
	var reasons []domain.ReportReason
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL and is_active = true").
		Find(&reasons).Error
	if err != nil {
		return nil, err
	}
	return reasons, nil
}

func (r *ReportReasonPostgres) Create(ctx context.Context, reason *domain.ReportReason) (*domain.ReportReason, error) {
	if err := r.db.WithContext(ctx).Create(reason).Error; err != nil {
		return nil, err
	}
	return reason, nil
}

func (r *ReportReasonPostgres) GetByID(ctx context.Context, id uint64) (*domain.ReportReason, error) {
	var reason domain.ReportReason
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&reason).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reason, err
}

func (r *ReportReasonPostgres) GetList(ctx context.Context, req *domain.ReportReasonListRequest) ([]domain.ReportReason, int64, error) {
	var reasons []domain.ReportReason
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.ReportReason{})

	// Apply filters
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	} else {
		query = query.Where("is_active = true")
	}

	query = query.Where("deleted_at IS NULL")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if req.Page > 0 && req.Size > 0 {
		offset := (req.Page - 1) * req.Size
		query = query.Offset(offset).Limit(req.Size)
	}

	// Execute query
	err := query.Order("name ASC").Find(&reasons).Error
	if err != nil {
		return nil, 0, err
	}

	return reasons, total, nil
}

func (r *ReportReasonPostgres) Update(ctx context.Context, reason *domain.ReportReason) error {
	return r.db.WithContext(ctx).
		Model(&domain.ReportReason{}).
		Where("id = ? AND deleted_at IS NULL", reason.ID).
		Updates(reason).Error
}

func (r *ReportReasonPostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&domain.ReportReason{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}