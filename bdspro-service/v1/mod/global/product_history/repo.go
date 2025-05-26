package product_accecss

import (
	"bdspro/v1/models"
	"common/case/crud2"

	"gorm.io/gorm"
)

type ProductHistoryRepo struct {
	crud2.BaseRepo[models.ProductAccess]
	// PrivateFieldRepo *product_private.ProductPrivateRepo
}

func NewProductHistoryRepo(db *gorm.DB) *ProductHistoryRepo {
	return &ProductHistoryRepo{
		BaseRepo: crud2.BaseRepo[models.ProductAccess]{DB: db},
		// PrivateFieldRepo: PrivateFieldRepo,
	}
}
