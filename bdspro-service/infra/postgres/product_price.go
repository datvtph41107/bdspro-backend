package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud3"
	_db "common/db"
	_utils "common/utils"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.ProductPriceRepo
type GormProductPriceRepo struct {
	*_db.TransactionRepo
	crud3.BaseRepo[domain.ProductPrice]
}

func NewGormProductPriceRepo(gormDb *gorm.DB, db *_db.TransactionRepo) *GormProductPriceRepo {
	return &GormProductPriceRepo{
		TransactionRepo: db,
		BaseRepo:        crud3.BaseRepo[domain.ProductPrice]{DB: gormDb},
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
		// Kiểm tra sự khác biệt về các field giá
		diff := !_utils.CompareEqual(&oldPrice.Currency, &entity.Currency) ||
			!_utils.CompareEqual(&oldPrice.SalePrice, &entity.SalePrice) ||
			!_utils.CompareEqual(&oldPrice.SaleCommission, &entity.SaleCommission) ||
			!_utils.CompareEqual(&oldPrice.RentPrice, &entity.RentPrice) ||
			!_utils.CompareEqual(&oldPrice.RentCommission, &entity.RentCommission) ||
			!_utils.CompareEqual(&oldPrice.Deposite, &entity.Deposite) ||
			!_utils.CompareEqual(&oldPrice.RentPaymentCycle, &entity.RentPaymentCycle) ||
			!_utils.CompareEqual(&oldPrice.RentCommissionType, &entity.RentCommissionType) ||
			!_utils.CompareEqual(&oldPrice.SaleCommissionType, &entity.SaleCommissionType) ||
			oldPrice.ChannelPrice != entity.ChannelPrice

		// Kiểm tra sự khác biệt về ChangeNote
		diffNote := oldPrice.ChangeNote != entity.ChangeNote

		// Nếu không có sự khác biệt nào, không làm gì
		if !diff {
			entity.ID = oldPrice.ID

			// Nếu chỉ có ChangeNote khác, chỉ cập nhật ChangeNote của price hiện tại
			if diffNote {
				err := GetDB(ctx, r.DB).
					Model(&domain.ProductPrice{}).
					Where("id = ? AND deleted_at IS NULL", *priceId).
					Updates(map[string]interface{}{
						"change_note": entity.ChangeNote,
					}).Error
				if err != nil {
					return err
				}
			}
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

func (r *GormProductPriceRepo) GetHistoryProductPrice(
	ctx context.Context,
	req *dto.PriceHistorySearch,
) ([]*domain.ProductPrice, int64, error) {

	var total int64
	prices := make([]*domain.ProductPrice, 0)

	db := r.GetDB(ctx).
		Model(&domain.ProductPrice{}).
		Joins("JOIN products ON products.id = product_price.product_id").
		Where("products.id = ?", req.ProductID).
		Where("product_price.deleted_at IS NULL")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Order("product_price.created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&prices).Error

	if err != nil {
		return nil, 0, err
	}

	return prices, total, nil
}

func (r *GormProductPriceRepo) GetDistributePrice(ctx context.Context, productId *uint64) (*domain.ProductPrice, error) {
	var price *domain.ProductPrice

	originId := _utils.GetOriginIdFromContext(ctx)
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Query xử lý 3 trường hợp:
	// 1. Nếu product_user có priceId thì join theo field này
	// 2. Nếu product_user không có priceId thì vào distribution join theo priceId của distribution
	// 3. Nếu product_user đó là isOwner=true thì join trực tiếp theo priceId trong product (LastPriceID)
	query := `
		SELECT DISTINCT pp.*
		FROM product_user pu
		LEFT JOIN distributions d ON d.id = pu.distribute_id
		LEFT JOIN products p ON p.id = pu.product_id AND p.deleted_at IS NULL
		INNER JOIN product_price pp ON pp.deleted_at IS NULL
			AND pp.id = CASE
				-- Trường hợp 3: Nếu isOwner=true thì ưu tiên dùng LastPriceID từ product
				WHEN pu.is_owner = true AND p.last_price_id IS NOT NULL THEN p.last_price_id
				-- Trường hợp 1: Nếu product_user có priceId thì dùng priceId này
				WHEN pu.price_id IS NOT NULL THEN pu.price_id
				-- Trường hợp 2: Nếu không có priceId trong product_user thì dùng priceId từ distribution
				WHEN d.price_id IS NOT NULL THEN d.price_id
				ELSE NULL
			END
		WHERE pu.product_id = ?
			AND (pu.origin_profile_id = ? OR pu.profile_id = ?)
			AND pu.deleted_at IS NULL
			AND pp.id IS NOT NULL
		ORDER BY pp.created_at DESC
		LIMIT 1
	`

	err := r.GetDB(ctx).Raw(query, *productId, originId, profileId).Scan(&price).Error
	if err != nil {
		return nil, err
	}
	return price, nil
}
