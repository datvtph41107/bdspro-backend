package postgres

import (
	"bdspro/infra/shared"
	"bdspro/internal/domain"
	"bdspro/internal/dto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.AssetShareRepo
type AssetShareRepo struct {
	shared.PostgreCrud[domain.AssetShare, *dto.AssetShareSearchDTO]
}

func NewAssetShareRepo(DB *gorm.DB) *AssetShareRepo {
	return &AssetShareRepo{
		PostgreCrud: shared.PostgreCrud[domain.AssetShare, *dto.AssetShareSearchDTO]{
			DB: DB,
		},
	}
}

func (r AssetShareRepo) ExistByTarget(c *gin.Context, ownerID, targetID, assetID uint64) (bool, error) {
	var count int64
	err := r.DB.
		Model(&domain.AssetShare{}).
		// Select("1").
		Where("deleted_at IS NULL and target_id = ? and asset_id = ? and owner_id = ?", targetID, assetID, ownerID).
		Count(&count).
		Error
	return count > 0, err
}

// Cập nhật bản ghi
// func (r AssetShareRepo) Update(c *gin.Context, id uint64, entity *domain.AssetShare) error {
// 	return r.DB.WithContext(c).
// 		Model(entity).
// 		Where("id = ? and deleted_at is null", id).
// 		Update("permissions", entity.Permissions).Error
// }
