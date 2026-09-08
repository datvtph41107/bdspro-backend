package entity

import (
	_models "common/models"
	"fmt"
	"time"
)

const (
	OrganizationStatusActive   = "active"
	OrganizationStatusInactive = "inactive"
	OrganizationStatusPending  = "pending"
)

type Organization struct {
	_models.BaseEntity
	Name                string `gorm:"not null;index"`
	TaxCode             string `gorm:"not null;index:idx_organizations_tax_code"`
	BusinessLicenseUrl  *string
	Address             *string
	DetailAddress       *string
	Phone               *string
	Email               *string
	Website             *string
	LogoUrl             *string
	Status              uint32 `gorm:"index:idx_organizations_status_created_at"`
	ApproveInvestor     bool   `gorm:"default:false"`
	AccountBankInvestor *uint64
	FoundedAt           *time.Time
	Description         *string
	OwnerId             uint64            `gorm:"index:idx_organizations_owner_id"`
	BusinessDomainIds   []uint32          `gorm:"-"`
	BusinessDomains     []*BusinessDomain `gorm:"many2many:organization_business_domains;joinForeignKey:OrganizationID;joinReferences:BusinessDomainID"`
}

func (Organization) TableName() string {
	return "organizations"
}

// GenerateBusinessCode generates business code in format TC + 6 digits (padded with zeros)
func (o *Organization) GenerateBusinessCode() string {
	return fmt.Sprintf("TC%06d", o.ID)
}

type OrganizationWithMemberCount struct {
	Organization
	TotalMembers uint32 `gorm:"column:total_members"`
}
