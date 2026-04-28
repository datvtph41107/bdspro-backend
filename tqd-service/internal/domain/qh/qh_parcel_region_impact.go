package qh_domain

import "time"

type QHParcelRegionImpact struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// Foreign keys
	ParcelID uint64 `gorm:"not null;index" json:"parcelId"`
	RegionID uint64 `gorm:"not null;index" json:"regionId"`
	LayerID  uint64 `gorm:"not null;index" json:"layerId"`

	// Impact metrics
	IntersectionAreaSqm float64 `gorm:"default:0" json:"intersectionAreaSqm"`
	OverlapPercentage   float64 `gorm:"default:0" json:"overlapPercentage"`

	// Cache data
	RegionLabel         string `gorm:"type:text" json:"regionLabel"`
	LandUseCodeAtImpact string `gorm:"type:varchar(20)" json:"landUseCodeAtImpact"`

	// Temporal validity
	EffectiveDate *time.Time `gorm:"type:date" json:"effectiveDate"`
	ExpiryDate    *time.Time `gorm:"type:date" json:"expiryDate"`

	// Audit
	CalculatedAt time.Time `gorm:"autoCreateTime" json:"calculatedAt"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (QHParcelRegionImpact) TableName() string {
	return "qh_parcel_region_impact"
}

// HasImpact - Kiểm tra có ảnh hưởng không
func (i *QHParcelRegionImpact) HasImpact() bool {
	return i.IntersectionAreaSqm > 0
}

// IsActive - Kiểm tra có đang hiệu lực không
func (i *QHParcelRegionImpact) IsActive() bool {
	now := time.Now()
	if i.EffectiveDate != nil && i.EffectiveDate.After(now) {
		return false
	}
	if i.ExpiryDate != nil && i.ExpiryDate.Before(now) {
		return false
	}
	return true
}
