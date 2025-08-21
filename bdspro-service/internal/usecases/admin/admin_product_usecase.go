package admin_usecases

import (
	"bdspro/internal/domain"
	admin_repo "bdspro/internal/repo/admin"
	"common/case/crud"
)

type AdminProductUsecase struct {
	crud.BaseUsecase[domain.Product, admin_repo.AdminProductRepo]
}

func NewAdminProductUsecase(repo admin_repo.AdminProductRepo) *AdminProductUsecase {
	return &AdminProductUsecase{
		crud.BaseUsecase[domain.Product, admin_repo.AdminProductRepo]{
			Repo: repo,
		},
	}
}
