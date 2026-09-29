package entity

import "time"

type OrganizationBranchMember struct {
	Id                   uint32
	OrganizationBranchId uint32
	UserId               uint32
	RoleId               uint64
	CreatedAt            time.Time
	UpdatedAt            time.Time
	CreatedBy            uint32
}
