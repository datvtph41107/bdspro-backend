package usecase

import (
	"context"

	"user/internal/domain"
	"user/internal/interface/repo"

	_crud "common/domain/crud"
)

type ModuleUsecase struct {
	_crud.BaseUsecase[domain.Module, repo.IModuleRepo]
}

func NewModuleUsecase(moduleRepo repo.IModuleRepo) *ModuleUsecase {
	return &ModuleUsecase{
		BaseUsecase: _crud.BaseUsecase[domain.Module, repo.IModuleRepo]{
			Repo: moduleRepo,
		},
	}
}

// GetAllActiveModules lấy tất cả module đang active
func (u *ModuleUsecase) GetAllActiveModules(ctx context.Context) ([]*domain.Module, error) {
	modules, err := u.Repo.GetAllActiveModules(ctx)
	if err != nil {
		return nil, err
	}
	return modules, nil
}

// GetModuleByCode lấy module theo code
func (u *ModuleUsecase) GetModuleByCode(ctx context.Context, code string) (*domain.Module, error) {
	module, err := u.Repo.GetModuleByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return module, nil
}
