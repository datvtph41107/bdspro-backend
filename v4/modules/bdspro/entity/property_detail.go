package entity

import "time"

// PropertyNamedTag is a lightweight id+name row (amenity, area region).
type PropertyNamedTag struct {
	ID   uint64 `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

// PropertyDeveloper is the developer linked via project (detail response).
type PropertyDeveloper struct {
	ID      uint64 `gorm:"column:id"`
	Name    string `gorm:"column:name"`
	Slug    string `gorm:"column:slug"`
	LogoURL string `gorm:"column:logo_url"`
	Phone   string `gorm:"column:phone"`
	Email   string `gorm:"column:email"`
	Website string `gorm:"column:website"`
}

// PropertyLineageDetail augments PropertyLineage for GET detail (not persisted as one row).
type PropertyLineageDetail struct {
	*PropertyLineage
	VerifiedNationalAt *time.Time
	Personalization    *PropertyUser
	Developer          *PropertyDeveloper
	Amenities          []PropertyNamedTag
	AreaRegions        []PropertyNamedTag
}
