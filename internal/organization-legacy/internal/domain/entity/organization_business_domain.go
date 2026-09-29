package entity

import "time"

type OrganizationBusinessDomain struct {
	ID               uint32
	OrganizationID   uint32
	BusinessDomainID uint32
	CreatedAt        time.Time
	CreatedBy        uint32
}
