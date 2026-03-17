package domain

import (
	_entity "common/domain/entity"
	"fmt"
)

// POI represents a Point of Interest
type POI struct {
	_entity.BaseEntity
	Code          string  `json:"code" gorm:"column:code;type:varchar(50);uniqueIndex"`
	Name          string  `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Description   string  `json:"description" gorm:"column:description;type:text"`
	Address       string  `json:"address" gorm:"column:address;type:text"`
	Phone         string  `json:"phone" gorm:"column:phone;type:varchar(50)"`
	Email         string  `json:"email" gorm:"column:email;type:varchar(255)"`
	Website       string  `json:"website" gorm:"column:website;type:varchar(255)"`
	Latitude      float64 `json:"latitude" gorm:"column:latitude;type:decimal(10,8)"`
	Longitude     float64 `json:"longitude" gorm:"column:longitude;type:decimal(11,8)"`
	CategoryID    uint64  `json:"categoryId" gorm:"column:category_id;index"`
	Rating        float64 `json:"rating" gorm:"column:rating;type:decimal(3,2);default:0"`
	ReviewCount   uint32  `json:"reviewCount" gorm:"column:review_count;default:0"`
	IsVerified    bool    `json:"isVerified" gorm:"column:is_verified;default:false"`
	IsActive      bool    `json:"isActive" gorm:"column:is_active;default:true"`
	IsFeatured    bool    `json:"isFeatured" gorm:"column:is_featured;default:false"`
	CoverImage    string  `json:"coverImage" gorm:"column:cover_image;type:varchar(500)"`
	ImagesJSON    string  `json:"-" gorm:"column:images;type:text"`
	TagsJSON      string  `json:"-" gorm:"column:tags;type:text"`
	AmenitiesJSON string  `json:"-" gorm:"column:amenities;type:text"`
	ViewCount     int64   `json:"viewCount" gorm:"column:view_count;default:0"`
	LikeCount     int64   `json:"likeCount" gorm:"column:like_count;default:0"`
	CreatedBy     uint64  `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy     uint64  `json:"updatedBy" gorm:"column:updated_by"`

	// Relations
	Category  *PoiCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID"`
	OpenHours []OpenHour   `json:"openHours,omitempty" gorm:"foreignKey:POIID;references:ID"`
	Amenities []Amenity    `json:"amenities,omitempty" gorm:"many2many:poi_amenities;"`
}

func (POI) TableName() string {
	return "pois"
}

func (p *POI) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.CategoryID == 0 {
		return fmt.Errorf("category ID is required")
	}
	if p.Latitude < -90 || p.Latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}
	if p.Longitude < -180 || p.Longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}
	return nil
}

// BeforeCreate gorm hook
func (p *POI) BeforeCreate() error {
	return p.Validate()
}

// BeforeUpdate gorm hook
func (p *POI) BeforeUpdate() error {
	return p.Validate()
}
