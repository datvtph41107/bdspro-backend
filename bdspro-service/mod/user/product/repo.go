package u_product

import (
	"bdspro/data"
	"bdspro/enums"
	"bdspro/models"
	"common/case/crud"
	_db "common/db"
	_utils "common/utils"
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserProductRepo struct {
	crud.BaseRepo[models.Product]
}

func NewUserProductRepo(db *gorm.DB) *UserProductRepo {
	return &UserProductRepo{
		BaseRepo: crud.BaseRepo[models.Product]{DB: db},
	}
}

func (r *UserProductRepo) UpdateStatus(ctx context.Context, id uint64, status enums.ProductStatus) error {
	return r.DB.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ? and deleted_at is null", id).
		Update("sale_status", status).Error
}

func (r *UserProductRepo) UpdateProduct(c *gin.Context, product *models.Product) error {
	return _db.SaveWithAudit(c, product)
}

func (r *UserProductRepo) UpdateAsset(asset *models.Asset) error {
	return r.DB.Save(asset).Error
}

func (r *UserProductRepo) GetByID(ctx context.Context, id uint64) (*models.Product, error) {
	var product models.Product
	err := r.DB.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *UserProductRepo) SearchProductWithOwner(profileId uint64, dto data.ProductSearchRequest, pageable _utils.Page) ([]models.Product, error) {
	var results []models.Product

	query := `
		SELECT p.*
		FROM product_access pa
		JOIN products p ON pa.product_id = p.id
		LIMIT ? OFFSET ?
	`

	err := r.DB.Raw(query, pageable.Limit, pageable.Skip).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
