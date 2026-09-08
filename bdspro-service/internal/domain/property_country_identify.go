// internal/domain/identifier.go
package domain

import (
	"time"

	"gorm.io/gorm"
)

// CountryIdentifier represents the core identifier entity
type CountryIdentifier struct {
	ID             string         `gorm:"column:id;primaryKey;size:40"`
	LandParcelCode string         `gorm:"column:land_parcel_code;size:30;index"`
	ProjectCode    string         `gorm:"column:project_code;size:30;index"`
	ProvinceID     string         `gorm:"column:province_id;size:5;not null;index:idx_location"`
	DistrictID     string         `gorm:"column:district_id;size:7;not null;index:idx_location"`
	WardID         string         `gorm:"column:ward_id;size:11;not null;index:idx_location"`
	Latitude       float64        `gorm:"column:latitude;type:decimal(10,8)"`
	Longitude      float64        `gorm:"column:longitude;type:decimal(11,8)"`
	Type           string         `gorm:"column:type;size:3;not null;index"`
	LegalStatus    string         `gorm:"column:legal_status;size:3;index"`
	CurrentOwnerID string         `gorm:"column:current_owner_id;size:20;index"`
	CreatedAt      time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName specifies the table name
func (CountryIdentifier) TableName() string {
	return "property_identifiers"
}

// BeforeCreate GORM hook để sinh ID
func (i *CountryIdentifier) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		// Gọi function PostgreSQL để sinh mã định danh
		return tx.Raw("SELECT generate_identifier_id(?, ?) as id",
			i.ProvinceID, i.Type).Scan(&i.ID).Error
	}
	return nil
}
