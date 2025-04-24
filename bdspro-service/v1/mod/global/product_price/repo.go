package product_price

import (
	"bdspro/v1/models"
	"common/case/crud2"
	_db "common/db"
	_utils "common/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductPriceRepo struct {
	crud2.BaseRepo[models.ProductPrice]
}

func NewProductPriceRepo(db *gorm.DB) *ProductPriceRepo {
	return &ProductPriceRepo{
		BaseRepo: crud2.BaseRepo[models.ProductPrice]{
			DB: db,
		},
	}
}

func (r *ProductPriceRepo) FindByProductId(productId *uint64) *models.ProductPrice {
	var price *models.ProductPrice

	_db.DB.
		Where("product_id = ? AND deleted_at IS NULL", productId).
		Order("created_at DESC").
		Limit(1).
		First(&price)

	return price
}

// func Compare[T1, T2 any](t1 T1, t2 T2) bool {
// 	if reflect.TypeOf(t1) == reflect.TypeOf(t2) {
// 		return reflect.DeepEqual(t1, t2)
// 	}
// 	return false
// }

func (r *ProductPriceRepo) UpdateProductID(c *gin.Context, productId *uint64, priceId *uint64) error {
	price, _ := r.GetByID(c, *priceId)
	price.ProductID = productId
	return _db.SaveWithAudit(c, price)
}
func (r *ProductPriceRepo) UpdatePrice(c *gin.Context, priceId uint64, entity *models.ProductPrice) error {
	if entity == nil {
		return nil
	}
	// oldPrice := r.FindByProductId(productId)
	oldPrice, _ := r.GetByID(c, priceId)
	if oldPrice != nil {
		diff := !_utils.CompareEqual(&oldPrice.Currency, &entity.Currency) ||
			!_utils.CompareEqual(&oldPrice.PriceOwner, &entity.PriceOwner) ||
			!_utils.CompareEqual(&oldPrice.SalePrice, &entity.SalePrice) ||
			!_utils.CompareEqual(&oldPrice.SaleCommission, &entity.SaleCommission) ||
			!_utils.CompareEqual(&oldPrice.RentPrice, &entity.RentPrice) ||
			!_utils.CompareEqual(&oldPrice.RentCommission, &entity.RentCommission) ||
			!_utils.CompareEqual(&oldPrice.RentPaymentCycle, &entity.RentPaymentCycle)
			// oldPrice.SaleCommissionType != entity.Price.SaleCommissionType ||
			// oldPrice.RentCommissionType != entity.Price.RentCommissionType ||

		if !diff {
			return nil
		}
	}

	return _db.SaveWithAudit(c, entity)
	// return _db.SaveWithAudit(c, &models.ProductPrice{
	// 	ProductID:      productId,
	// 	Currency:       entity.Currency,
	// 	PriceOwner:     entity.PriceOwner,
	// 	SalePrice:      entity.SalePrice,
	// 	SaleCommission: entity.SaleCommission,
	// 	// SaleCommissionType: entity.SaleCommissionType,
	// 	RentPrice:      entity.RentPrice,
	// 	RentCommission: entity.RentCommission,
	// 	// RentCommissionType: entity.RentCommissionType,
	// 	RentPaymentCycle: entity.RentPaymentCycle,
	// })
}
