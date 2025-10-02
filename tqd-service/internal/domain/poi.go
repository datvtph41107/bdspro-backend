package domain

import (
	_models "common/domain/entity"
	"time"
)

// POI represents a Point of Interest
type POI struct {
	_models.BaseEntity
	Code           string     `json:"code" gorm:"size:50;uniqueIndex"`
	CategoryID     *uint64    `json:"categoryId" gorm:"column:category_id"`
	Name           string     `json:"name" gorm:"size:255;not null"`
	Address        string     `json:"address" gorm:"type:text"`
	Lat            float64    `json:"lat" gorm:"type:decimal(10,8)"`
	Lng            float64    `json:"lng" gorm:"type:decimal(11,8)"`
	GeoJSON        string     `json:"geoJson" gorm:"type:text"`
	Phone          string     `json:"phone" gorm:"size:20"`
	Website        string     `json:"website" gorm:"size:255"`
	RatingPoint    float64    `json:"ratingPoint" gorm:"type:decimal(3,2);default:0"`
	RatingCount    int        `json:"ratingCount" gorm:"default:0"`
	Source         string     `json:"source" gorm:"size:100"`
	Confidence     float64    `json:"confidence" gorm:"type:decimal(3,2);default:0"`
	LastVerifiedAt *time.Time `json:"lastVerifiedAt"`
	NextReviewAt   *time.Time `json:"nextReviewAt"`
	IsActive       bool       `json:"isActive" gorm:"default:true"`

	// Relations
	Category *PoiCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID"`
}

func (POI) TableName() string {
	return "pois"
}
