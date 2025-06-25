package entity

import (
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
	Status             string
	ApproveInvestor    bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
	CreatedBy          uint32
	UpdatedBy          uint32
}
