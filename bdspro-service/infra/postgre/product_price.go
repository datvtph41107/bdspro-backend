package postgre

import (
	"bdspro/internal/domain"
	"common/case/crud3"
	_utils "common/utils"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.ProductPriceRepo
type GormProductPriceRepo struct {
	crud3.BaseRepo[domain.ProductPrice]
}

func NewGormProductPriceRepo(db *gorm.DB) *GormProductPriceRepo {
	return &GormProductPriceRepo{
		BaseRepo: crud3.BaseRepo[domain.ProductPrice]{
			DB: db,
		},
	}
}

func (r *GormProductPriceRepo) FindByProductId(ctx context.Context, productId *uint64) *domain.ProductPrice {
	var price *domain.ProductPrice

	GetDB(ctx, r.DB).
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

// Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
func (r *GormProductPriceRepo) GetByID(c context.Context, id uint64) (*domain.ProductPrice, error) {
	var entity domain.ProductPrice
	err := GetDB(c, r.DB).Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormProductPriceRepo) UpdateProductID(c context.Context, productId *uint64, priceId *uint64) error {
	price, _ := r.GetByID(c, *priceId)
	if price == nil {
		return nil
	}
	price.ProductID = productId
	return GetDB(c, r.DB).Save(price).Error
}
func (r *GormProductPriceRepo) UpdatePrice(ctx context.Context, priceId *uint64, entity *domain.ProductPrice) error {
	if entity == nil {
		return nil
	}
	if priceId == nil {
		return GetDB(ctx, r.DB).Create(entity).Error
	}
	// oldPrice := r.FindByProductId(productId)
	oldPrice, _ := r.GetByID(ctx, *priceId)
	if oldPrice != nil {
		diff := !_utils.CompareEqual(&oldPrice.Currency, &entity.Currency) ||
			// !_utils.CompareEqual(&oldPrice.PriceOwner, &entity.PriceOwner) ||
			!_utils.CompareEqual(&oldPrice.SalePrice, &entity.SalePrice) ||
			!_utils.CompareEqual(&oldPrice.SaleCommission, &entity.SaleCommission) ||
			!_utils.CompareEqual(&oldPrice.RentPrice, &entity.RentPrice) ||
			!_utils.CompareEqual(&oldPrice.RentCommission, &entity.RentCommission) ||
			!_utils.CompareEqual(&oldPrice.Deposite, &entity.Deposite) ||
			!_utils.CompareEqual(&oldPrice.RentPaymentCycle, &entity.RentPaymentCycle)
			// oldPrice.SaleCommissionType != entity.Price.SaleCommissionType ||
			// oldPrice.RentCommissionType != entity.Price.RentCommissionType ||

		if !diff {
			entity.ID = oldPrice.ID
			return nil
		}
	}

	return GetDB(ctx, r.DB).Save(entity).Error
	// return _db.SaveWithAudit(c, &domain.ProductPrice{
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

func (r *GormProductPriceRepo) Create(c context.Context, entity *domain.ProductPrice) error {
	return GetDB(c, r.DB).Create(entity).Error
}
