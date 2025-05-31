package entity

import (
	"time"
)

// OrganizationCertification represents a certification or license of an organization
type OrganizationCertification struct {
	ID               uint32     `json:"id" gorm:"primary_key;auto_increment"`
	OrganizationID   uint32     `json:"organization_id" gorm:"not null"`
	CertificationType string     `json:"certification_type" gorm:"not null"`
	DocumentURL      string     `json:"document_url" gorm:"not null"`
	IssuedDate       *time.Time `json:"issued_date"`
	ExpiryDate       *time.Time `json:"expiry_date"`
	Status           string     `json:"status" gorm:"default:valid"`

	// Relationships
	Organization *Organization `json:"organization" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for the OrganizationCertification entity
func (OrganizationCertification) TableName() string {
	return "organization_certifications"
}

// BeforeCreate hook to validate expiry date and set status
func (oc *OrganizationCertification) BeforeCreate() error {
	if oc.ExpiryDate != nil && oc.ExpiryDate.Before(time.Now()) {
		oc.Status = "expired"
	}
	return nil
}

// BeforeUpdate hook to update status based on expiry date
func (oc *OrganizationCertification) BeforeUpdate() error {
	if oc.ExpiryDate != nil {
		if oc.ExpiryDate.Before(time.Now()) {
			oc.Status = "expired"
		} else {
			oc.Status = "valid"
		}
	}
	return nil
} 