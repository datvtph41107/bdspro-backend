package entity

import (
	"time"

	"gorm.io/gorm"
)

// HouseInfo maps to house_info (1:1 with products).
type HouseInfo struct {
	ProductID        uint64         `gorm:"column:product_id;primaryKey"`
	CreatedAt        *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy        *uint64        `gorm:"column:created_by"`
	UpdatedBy        *uint64        `gorm:"column:updated_by"`
	NumBedroom       *int32         `gorm:"column:num_bedroom"`
	NumBathroom      *int32         `gorm:"column:num_bathroom"`
	NumFloor         *int32         `gorm:"column:num_floor"`
	NumFront         *int32         `gorm:"column:num_front"`
	NumCarPark       *int32         `gorm:"column:num_car_park"`
	NumToilet        *int32         `gorm:"column:num_toilet"`
	Furniture        *string        `gorm:"column:furniture;size:50"`
	Orientation      *string        `gorm:"column:orientation;size:50"`
	OrientationHouse uint32         `gorm:"column:orientation_house"`
	CertificateHouse uint32         `gorm:"column:certificate_house"`
	RoadWidth        *float64       `gorm:"column:road_width;type:decimal(10,2)"`
	FrontWidth       *float64       `gorm:"column:front_width"`
	BackWidth        *float64       `gorm:"column:back_width"`
	Width            *float64       `gorm:"column:width"`
	Height           *float64       `gorm:"column:height"`
}

func (HouseInfo) TableName() string { return "house_info" }
