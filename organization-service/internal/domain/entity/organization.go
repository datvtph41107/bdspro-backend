package entity

import (
	_models "common/models"
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
	BusinessDomainIds   []uint32          `gorm:"-"`
	BusinessDomains     []*BusinessDomain `gorm:"many2many:organization_business_domains;joinForeignKey:OrganizationID;joinReferences:BusinessDomainID"`
}

func (Organization) TableName() string {
	return "organizations"
}
