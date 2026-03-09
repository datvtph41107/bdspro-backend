package postgres

import (
	_db "common/db"
	"context"
	"errors"

	"hub/helpers"
	"hub/internal/domain"
	"hub/internal/repo"

	"gorm.io/gorm"
)

type ApplinkPostgre struct {
	db *_db.TransactionRepo
}

func NewApplinkPostgre(db *_db.TransactionRepo) repo.IApplinkRepo {
	return &ApplinkPostgre{db: db}
}

// ExistsByCode kiểm tra code đã tồn tại chưa
func (r *ApplinkPostgre) ExistsByCode(ctx context.Context, code uint64) (bool, error) {
	var count int64
	err := r.db.GetDB(ctx).
		Model(&domain.Applink{}).
		Where("code = ?", code).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create insert applink mới — GORM populate ID + timestamps sau insert
func (r *ApplinkPostgre) Create(ctx context.Context, a *domain.Applink) error {
	return r.db.GetDB(ctx).Create(a).Error
}

// GetByCode tìm theo code — hit uniqueIndex idx_applink_code
func (r *ApplinkPostgre) GetByCode(ctx context.Context, code uint64) (*domain.Applink, error) {
	var a domain.Applink
	err := r.db.GetDB(ctx).
		Where("code = ?", code).
		First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrApplinkNotFound
		}
		return nil, err
	}
	return &a, nil
}
