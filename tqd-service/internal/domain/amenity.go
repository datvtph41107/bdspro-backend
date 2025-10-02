package domain

import (
	_models "common/domain/entity"
)

// Amenity represents facilities or services available at a location
type Amenity struct {
	_models.BaseEntity
	Name        string `json:"name" gorm:"size:255;not null"`
	Description string `json:"description" gorm:"type:text"`
	Icon        string `json:"icon" gorm:"size:255"` // Icon URL or class
}

func (Amenity) TableName() string {
	return "amenities"
}
