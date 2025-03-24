package models

type ProductHouserInfo struct {
	ProductID   uint64  `gorm:"primaryKey;not null" json:"product_id"`
	Bedrooms    int     `gorm:"check:bedrooms >= 0" json:"bedrooms"`
	Bathrooms   int     `gorm:"check:bathrooms >= 0" json:"bathrooms"`
	Floors      *int    `gorm:"check:floors >= 0" json:"floors,omitempty"`
	Furniture   *string `gorm:"type:VARCHAR(50)" json:"furniture,omitempty"`
	Orientation *string `gorm:"type:VARCHAR(50)" json:"orientation,omitempty"`
}

// GORM table name override
func (ProductHouserInfo) TableName() string {
	return "product_houser_info"
}
