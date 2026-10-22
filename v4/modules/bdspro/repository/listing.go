package repository

import (
	"context"

	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/entity"

	"gorm.io/gorm"
)

// ListingRepository persists product listings on the v4 DSN.
type ListingRepository interface {
	CreateForProduct(ctx context.Context, listing *entity.Listing) error
	List(ctx context.Context, filter dto.ListingListRequest, ownerID uint64) ([]entity.Listing, int64, error)
}

type listingPostgres struct {
	db *gorm.DB
}

func NewListingPostgres(db *gorm.DB) ListingRepository {
	return &listingPostgres{db: db}
}

func (r *listingPostgres) CreateForProduct(ctx context.Context, listing *entity.Listing) error {
	if r.db == nil || listing == nil || listing.ProductID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := createSerialEntity(tx, "listings", listing); err != nil {
			return err
		}
		return tx.Model(&entity.Product{}).Where("id = ?", listing.ProductID).
			UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).Error
	})
}

func (r *listingPostgres) List(ctx context.Context, filter dto.ListingListRequest, ownerID uint64) ([]entity.Listing, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.Listing{}).Where("deleted_at IS NULL")
	if ownerID > 0 {
		q = q.Where("owner_id = ?", ownerID)
	}
	if filter.ProductID > 0 {
		q = q.Where("product_id = ?", filter.ProductID)
	}
	if filter.TransactionType > 0 {
		q = q.Where("transaction_type = ?", filter.TransactionType)
	}
	if filter.Text != "" {
		q = q.Where("(LOWER(title) LIKE '%' || LOWER(?) || '%' OR LOWER(content) LIKE '%' || LOWER(?) || '%')", filter.Text, filter.Text)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []entity.Listing
	err := q.Order("updated_at DESC, id DESC").
		Offset(filter.GetOffset()).
		Limit(filter.GetLimit()).
		Find(&items).Error
	return items, total, err
}
