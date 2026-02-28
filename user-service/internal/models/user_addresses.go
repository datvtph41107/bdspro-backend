package models

import (
	"time"

	"github.com/lib/pq"
)

type UserAddress struct {
	ID          string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      string `gorm:"column:user_id;type:uuid;not null;index"`
	AddressType string `gorm:"column:address_type;size:20;not null"`

	FullAddress string  `gorm:"column:full_address;type:text;not null"`
	Street      *string `gorm:"size:255"`
	Ward        *string `gorm:"size:100"`
	District    *string `gorm:"size:100"`
	City        string  `gorm:"size:100;not null"`
	Region      *string `gorm:"size:50"`
	Country     string  `gorm:"size:50;default:'Vietnam'"`
	PostalCode  *string `gorm:"size:20"`

	Latitude  float64 `gorm:"column:latitude;type:double precision;not null"`
	Longitude float64 `gorm:"column:longitude;type:double precision;not null"`
	Geohash5  string  `gorm:"column:geohash_5;size:5;not null;index"`
	Geohash7  string  `gorm:"column:geohash_7;size:7;not null;index"`
	Geohash9  *string `gorm:"column:geohash_9;size:9"`

	PlaceID          *string `gorm:"column:place_id;size:255"`
	FormattedAddress *string `gorm:"column:formatted_address;type:text"`
	LocationType     *string `gorm:"column:location_type;size:50"`

	IsVerified         bool       `gorm:"column:is_verified;default:false;index"`
	VerifiedAt         *time.Time `gorm:"column:verified_at"`
	VerificationMethod *string    `gorm:"column:verification_method;size:50"`

	MovedInAt  *time.Time `gorm:"column:moved_in_at;type:date"`
	MovedOutAt *time.Time `gorm:"column:moved_out_at;type:date"`
	IsCurrent  bool       `gorm:"->;type:boolean GENERATED ALWAYS AS (moved_out_at IS NULL) STORED"`

	PropertyType      *string  `gorm:"column:property_type;size:50"`
	LivingAreaSqm     *float64 `gorm:"column:living_area_sqm;type:float"`
	NumberOfResidents *int     `gorm:"column:number_of_residents"`

	TotalCheckins int        `gorm:"column:total_checkins;default:0"`
	LastCheckinAt *time.Time `gorm:"column:last_checkin_at"`

	Tags      pq.StringArray `gorm:"column:tags;type:jsonb;default:'[]'::jsonb"`
	IsDefault bool           `gorm:"column:is_default;default:false;index:,composite:user_default,where:is_default=true"`
	Notes     *string        `gorm:"column:notes;type:text"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserAddress) TableName() string {
	return "user_addresses"
}
