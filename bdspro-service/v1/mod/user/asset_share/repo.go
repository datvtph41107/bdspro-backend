package u_asset_share

import (
	"bdspro/v1/models"
	_owner "common/case/owner"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetShareRepo struct {
	_owner.OwnerRepo[models.AssetShare]
}

func NewAssetShareRepo(DB *gorm.DB) *AssetShareRepo {
	return &AssetShareRepo{
		OwnerRepo: _owner.OwnerRepo[models.AssetShare]{
			DB: DB,
		},
	}
}

func (r AssetShareRepo) ExistByTarget(c *gin.Context, ownerID, targetID, assetID uint64) (bool, error) {
	var count int64
	err := r.DB.
		Model(&models.AssetShare{}).
		// Select("1").
		Where("deleted_at IS NULL and target_id = ? and asset_id = ? and owner_id = ?", targetID, assetID, ownerID).
		Count(&count).
		Error
	return count > 0, err
}

// Cập nhật bản ghi
func (r AssetShareRepo) Update(c *gin.Context, id uint64, entity *_owner.OwnerModel[models.AssetShare]) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Update("permissions", entity.M.Permissions).Error
}
