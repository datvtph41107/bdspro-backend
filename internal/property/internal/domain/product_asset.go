package domain

import (
	_models "common/models"
)

// ProductAsset đại diện cho bảng liên kết many-to-many giữa sản phẩm và tài sản
type ProductAsset struct {
	_models.BaseEntity
	ProductID uint64 `gorm:"not null;index" json:"productId"`
	AssetID   uint64 `gorm:"not null;index" json:"assetId"`
}

func (ProductAsset) TableName() string {
	return "product_asset"
}
