package entity

import (
	"organization/internal/enums"
	"time"
)

const (
	OrganizationStatusActive   = "active"
	OrganizationStatusInactive = "inactive"
	OrganizationStatusPending  = "pending"
)

type Organization struct {
	ID                 uint32
	Name               string
	TaxCode            string
	BusinessLicenseUrl *string
	Address            *string
	Phone              *string
	Email              *string
	Website            *string
	LogoUrl            *string
	Status             enums.EMemberStatus
	FoundedAt          *time.Time
	Description        *string
	ApproveInvestor    bool
	BusinessDomains    []*BusinessDomain
	BusinessDomainIds  []uint32
	CreatedAt          time.Time
	UpdatedAt          time.Time
	CreatedBy          uint32
	UpdatedBy          uint32
}
