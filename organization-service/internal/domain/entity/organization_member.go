package entity

import (
	"time"
)

// OrganizationMember represents the organization member entity
type OrganizationMember struct {
	ID             uint32     `json:"id" gorm:"primary_key;auto_increment"`
	OrganizationID uint32     `json:"organization_id" gorm:"not null"`
	UserID         uint32     `json:"user_id" gorm:"not null"`
	Role           string     `json:"role" gorm:"not null"`
	Status         string     `json:"status" gorm:"default:active"`
	JoinedAt       time.Time  `json:"joined_at" gorm:"not null"`
	RemovedAt      *time.Time `json:"removed_at"`

	// Relationships
	Organization *Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
}

// TableName specifies the table name for the OrganizationMember entity
func (OrganizationMember) TableName() string {
	return "organization_members"
} 