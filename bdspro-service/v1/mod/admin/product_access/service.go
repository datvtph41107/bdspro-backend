package a_product_access

import (
	product_accecss "bdspro/v1/mod/global/product_access"
	"bdspro/v1/models"
	"common/case/crud2"
)

type ProductAccessService struct {
	crud2.BaseService[models.ProductAccess, *product_accecss.ProductAccessRepo]
}

func NewProductAccessService(repo *product_accecss.ProductAccessRepo) *ProductAccessService {
	return &ProductAccessService{
		BaseService: crud2.BaseService[models.ProductAccess, *product_accecss.ProductAccessRepo]{Repo: repo},
	}
}
