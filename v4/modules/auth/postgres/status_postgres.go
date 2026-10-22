package postgres

import (
	"context"
	"errors"

	"mdn4/modules/auth/domain/auth"

	"gorm.io/gorm"
)

// @bind: mdn4/internal/interface/repo.AuthStatusRepository
type StatusPostgres struct {
	db *gorm.DB
}

// NewStatusRepository tạo mới AuthStatusRepository
func NewStatusRepository(db *gorm.DB) *StatusPostgres {
	return &StatusPostgres{db: db}
}

// GetByID lấy auth_status theo authID. Nếu bản ghi bị miss thì tạo skeleton đăng ký.
func (r *StatusPostgres) GetByID(ctx context.Context, authID uint64) (*auth.AuthStatus, error) {
	var status auth.AuthStatus
	err := r.db.WithContext(ctx).Where("auth_id = ?", authID).First(&status).Error
	if err == nil {
		return &status, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	status = auth.AuthStatus{
		IdDomain: auth.IdDomain{AuthID: authID},
		Active:   true,
		Verified: false,
	}
	if err := r.db.WithContext(ctx).Create(&status).Error; err != nil {
		if readErr := r.db.WithContext(ctx).Where("auth_id = ?", authID).First(&status).Error; readErr != nil {
			return nil, err
		}
	}
	return &status, nil
}

func (r *StatusPostgres) CreateStatus(ctx context.Context, status *auth.AuthStatus) error {
	return r.db.WithContext(ctx).Create(status).Error
}

func (r *StatusPostgres) UpdateStatus(ctx context.Context, status *auth.AuthStatus) error {
	return r.db.WithContext(ctx).Save(status).Error
}
