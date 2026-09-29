package admin_usecases

import (
	"bdspro/internal/domain"
	admin_repo "bdspro/internal/repo/admin"
	"common/case/crud"
)

type AdminPropertyTypeUsecase struct {
	crud.BaseUsecase[domain.PropertyType, admin_repo.AdminPropertyTypeRepo]
}

func NewAdminPropertyTypeUsecase(repo admin_repo.AdminPropertyTypeRepo) *AdminPropertyTypeUsecase {
	return &AdminPropertyTypeUsecase{
		crud.BaseUsecase[domain.PropertyType, admin_repo.AdminPropertyTypeRepo]{
			Repo: repo,
		},
	}
}
