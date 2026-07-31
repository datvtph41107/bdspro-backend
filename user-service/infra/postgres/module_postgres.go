package postgres

import (
	"context"

	"user/internal/domain"
	"user/internal/interface/repo"

	_db "common/db"
	_provider "common/provider"
)

type ModuleRepo struct {
	_provider.CrudRepo[domain.Module]
}

func NewModuleRepo(db *_db.TransactionRepo) repo.IModuleRepo {
	repo := &ModuleRepo{}
	repo.Init(repo, db)
	return repo
}

// GetAllActiveModules lấy tất cả module đang active
func (r *ModuleRepo) GetAllActiveModules(ctx context.Context) ([]*domain.Module, error) {
	var modules []*domain.Module
	err := r.GetDB(ctx).Where("is_active = ?", true).Order("priority ASC, created_at ASC").Find(&modules).Error
	if err != nil {
		return nil, err
	}
	return modules, nil
}

// GetModuleByCode lấy module theo code
func (r *ModuleRepo) GetModuleByCode(ctx context.Context, code string) (*domain.Module, error) {
	var module domain.Module
	err := r.GetDB(ctx).Where("code = ? AND is_active = ?", code, true).First(&module).Error
	if err != nil {
		return nil, err
	}
	return &module, nil
}
