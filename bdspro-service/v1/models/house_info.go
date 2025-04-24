package models

import _models "common/models"

type HouseInfo struct {
	_models.BaseEntityNotId
	ProductID   uint64  `gorm:"primaryKey;not null" json:"product_id"`
	NumBedroom  *int    `json:"numBedroom,omitempty"`                          // Số phòng ngủ (optional)
	NumBathroom *int    `json:"numBathroom,omitempty"`                         // Số WC (optional)
	NumFloor    *int    `json:"numFloor,omitempty"`                            // Số tầng (optional)
	NumFront    *int    `json:"numFront,omitempty"`                            // Số mặt tiền
	NumCarPark  *int    `json:"numCarPark,omitempty"`                          // Số chỗ để ôtô
	Furniture   *string `gorm:"type:varchar(50)" json:"furniture,omitempty"`   // Nội thất: "full", "partial", "none"
	Orientation *string `gorm:"type:varchar(50)" json:"orientation,omitempty"` // Hướng nhà (Đông, Tây, Nam,...)
	NumToilet   *int    `json:"numToilet,omitempty"`
}

// GORM table name override
func (HouseInfo) TableName() string {
	return "house_info"
}
