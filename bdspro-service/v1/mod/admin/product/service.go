package a_product

import (
	"bdspro/v1/models"
	"common/case/crud2"
)

type AdminProductService struct {
	crud2.BaseService[models.Product, *AdminProductRepo]
}

func NewAdminProductService(repo *AdminProductRepo) *AdminProductService {
	return &AdminProductService{
		BaseService: crud2.BaseService[models.Product, *AdminProductRepo]{
			Repo: repo,
		},
	}
}
