package postgres

import (
	"context"

	"gorm.io/gorm"
	"user/internal/domain/auth"
)

// @bind: user/internal/interface/repo.StatusRepository
type StatusPostgres struct {
	db *gorm.DB
}

// NewStatusRepository tạo mới StatusRepository
func NewStatusRepository(db *gorm.DB) *StatusPostgres {
	return &StatusPostgres{db: db}
}

// GetByID lấy status theo authID
func (r *StatusPostgres) GetByID(ctx context.Context, authID uint64) (*auth.UserStatusEntity, error) {
	var status auth.UserStatusEntity
	err := r.db.Where("auth_id = ?", authID).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// CreateStatus tạo mới status
func (r *StatusPostgres) CreateStatus(ctx context.Context, status *auth.UserStatusEntity) error {
	return r.db.Create(status).Error
}

// UpdateStatus cập nhật status
func (r *StatusPostgres) UpdateStatus(ctx context.Context, status *auth.UserStatusEntity) error {
	return r.db.Save(status).Error
}
