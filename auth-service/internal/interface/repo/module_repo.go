package repo

import (
	"context"

	"auth/internal/domain"

	_crud "common/domain/crud"
)

// @bind: auth/internal/interface/repo
type IModuleRepo interface {
	_crud.ICrudRepo[domain.Module]

	// GetAllActiveModules lấy tất cả module đang active
	GetAllActiveModules(ctx context.Context) ([]*domain.Module, error)

	// GetModuleByCode lấy module theo code
	GetModuleByCode(ctx context.Context, code string) (*domain.Module, error)
}
