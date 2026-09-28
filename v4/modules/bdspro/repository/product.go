package repository

import (
	"context"

	"userv4/modules/bdspro/entity"

	"gorm.io/gorm"
)

// ProductRepository is the product-facing store on the v4 DSN.
type ProductRepository interface {
	ListMine(ctx context.Context, userID uint64, page, size uint32, archived *bool) ([]*entity.Product, int64, error)
	GetByID(ctx context.Context, id uint64) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) error
	SoftDelete(ctx context.Context, id, profileID uint64) error
}

type productPostgres struct {
	db *gorm.DB
}

func NewProductPostgres(db *gorm.DB) ProductRepository {
	return &productPostgres{db: db}
}

func (r *productPostgres) ListMine(ctx context.Context, userID uint64, page, size uint32, archived *bool) ([]*entity.Product, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.Product{}).
		Joins("JOIN product_user pu ON pu.product_id = products.id AND pu.deleted_at IS NULL").
		Where("pu.profile_id = ?", userID).
		Where("products.deleted_at IS NULL")
	if archived != nil {
		q = q.Where("products.archived = ?", *archived)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := int(size)
	if limit <= 0 {
		limit = 50
	}
	offset := int(page) * limit

	var list []*entity.Product
	err := q.Preload("Price").
		Order("products.updated_at DESC NULLS LAST").
		Limit(limit).Offset(offset).
		Find(&list).Error
	return list, total, err
}

func (r *productPostgres) GetByID(ctx context.Context, id uint64) (*entity.Product, error) {
	if r.db == nil {
		return nil, nil
	}
	var p entity.Product
	err := r.db.WithContext(ctx).Preload("Price").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&p).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productPostgres) Create(ctx context.Context, product *entity.Product) error {
	if r.db == nil || product == nil {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		price := product.Price
		product.Price = nil
		if price != nil {
			if err := tx.Create(price).Error; err != nil {
				return err
			}
			product.LastPriceID = &price.ID
		}
		if err := tx.Create(product).Error; err != nil {
			return err
		}
		if price != nil {
			pid := product.ID
			if err := tx.Model(price).Update("product_id", pid).Error; err != nil {
				return err
			}
			product.Price = price
		}
		pu := &entity.ProductUser{
			ProfileID:       product.OwnerID,
			ProductID:       product.ID,
			IsOwner:         true,
			OwnerOf:         product.OwnerOf,
			OriginProfileID: &product.OwnerID,
		}
		return tx.Create(pu).Error
	})
}

func (r *productPostgres) SoftDelete(ctx context.Context, id, profileID uint64) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("id = ? AND owner_id = ?", id, profileID).Delete(&entity.Product{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			var count int64
			if err := tx.Model(&entity.ProductUser{}).
				Where("product_id = ? AND profile_id = ? AND deleted_at IS NULL", id, profileID).
				Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return gorm.ErrRecordNotFound
			}
			if err := tx.Delete(&entity.Product{}, id).Error; err != nil {
				return err
			}
		}
		return tx.Where("product_id = ?", id).Delete(&entity.ProductUser{}).Error
	})
}
