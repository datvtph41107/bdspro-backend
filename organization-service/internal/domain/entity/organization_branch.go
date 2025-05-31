package entity

import (
	"time"
)

// OrganizationBranch represents a branch of an organization
type OrganizationBranch struct {
	ID             uint32    `json:"id" gorm:"primary_key;auto_increment"`
	OrganizationID uint32    `json:"organization_id" gorm:"not null"`
	Name           string    `json:"name" gorm:"not null"`
	Address        string    `json:"address" gorm:"not null"`
	Phone          *string   `json:"phone"`
	Email          *string   `json:"email"`
	Manager        *string   `json:"manager"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	Organization *Organization `json:"organization" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for the OrganizationBranch entity
func (OrganizationBranch) TableName() string {
	return "organization_branches"
} 