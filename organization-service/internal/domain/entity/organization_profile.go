package entity

import (
	"time"
)

// OrganizationProfile represents the detailed profile of an organization
type OrganizationProfile struct {
	OrganizationID     uint32     `json:"organization_id" gorm:"primary_key"`
	Description        *string    `json:"description"`
	Industry          *string    `json:"industry"`
	NumberOfEmployees *int32     `json:"number_of_employees"`
	FoundingDate      *time.Time `json:"founding_date"`

	// Relationships
	Organization *Organization `json:"organization" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for the OrganizationProfile entity
func (OrganizationProfile) TableName() string {
	return "organization_profiles"
} 