package domain

import _models "common/models"

// PropertyProduct liên kết nhiều products với một property (để link product con với property)
type PropertyProduct struct {
	_models.BaseEntity
	PropertyID uint64           `gorm:"type:int8;not null;index" json:"propertyId"` // Liên kết với Property
	Property   *PropertyLineage `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`
	ProductID  uint64           `gorm:"type:int8;not null;index" json:"productId"` // Liên kết với Product
	Product    *Product         `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
}

// TableName override tên bảng
func (PropertyProduct) TableName() string {
	return "property_product"
}
