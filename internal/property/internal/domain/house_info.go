package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

type HouseInfo struct {
	_models.BaseEntityNotId
	ProductID   uint64  `gorm:"primaryKey;not null" json:"product_id"`
	NumBedroom  *int32  `json:"numBedroom,omitempty"`                          // Số phòng ngủ (optional)
	NumBathroom *int32  `json:"numBathroom,omitempty"`                         // Số WC (optional)
	NumFloor    *int32  `json:"numFloor,omitempty"`                            // Số tầng (optional)
	NumFront    *int32  `json:"numFront,omitempty"`                            // Số mặt tiền
	NumCarPark  *int32  `json:"numCarPark,omitempty"`                          // Số chỗ để ôtô
	Furniture   *string `gorm:"type:varchar(50)" json:"furniture,omitempty"`   // Nội thất: "full", "partial", "none"
	Orientation *string `gorm:"type:varchar(50)" json:"orientation,omitempty"` // Hướng nhà (Đông, Tây, Nam,...)
	NumToilet   *int32  `json:"numToilet,omitempty"`

	OrientationHouse enums.EHouseOrient      `gorm:"type:int8" json:"orientationHouse,omitempty"`   // Hướng nhà
	OrientationName  string                  `gorm:"-" json:"orientationName,omitempty"`            // Tên hướng nhà
	CertificateHouse enums.EHouseCertificate `gorm:"type:int8" json:"certificateHouse,omitempty"`   // Chứng nhận
	CertificateName  string                  `gorm:"-" json:"certificateName,omitempty"`            // Tên giấy chứng nhận
	RoadWidth        *float64                `gorm:"type:decimal(10,2)" json:"roadWidth,omitempty"` // Đường vào (m)

	FrontWidth *float64 `json:"frontWidth,omitempty"` // mặt tiền ? m
	BackWidth  *float64 `json:"backWidth,omitempty"`  // hậu ? m
	Width      *float64 `json:"width,omitempty"`      // rộng ? m
	Height     *float64 `json:"height,omitempty"`     // dài ? m
}

// GORM table name override
func (HouseInfo) TableName() string {
	return "house_info"
}
