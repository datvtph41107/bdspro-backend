package repo

import (
	"auth/internal/domain"
	"context"
)

// IAuthConfigRepo interface cho repository auth_config
type IAuthConfigRepo interface {
	// GetByKey lấy config theo key
	GetByKey(ctx context.Context, key string) (*domain.AuthConfig, error)

	// GetByType lấy tất cả config theo type
	GetByType(ctx context.Context, configType string) ([]*domain.AuthConfig, error)

	// GetAll lấy tất cả config
	GetAll(ctx context.Context) ([]*domain.AuthConfig, error)

	// Create tạo config mới
	Create(ctx context.Context, config *domain.AuthConfig) error

	// Update cập nhật config
	Update(ctx context.Context, config *domain.AuthConfig) error

	// UpdateValue cập nhật giá trị config theo key
	UpdateValue(ctx context.Context, key string, value string) error

	// Delete xóa config
	Delete(ctx context.Context, id uint64) error
}
