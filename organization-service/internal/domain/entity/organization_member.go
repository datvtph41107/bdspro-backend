package entity

import (
	"time"
)



type OrganizationMemberStatus uint32

const (
	OrganizationMemberStatusActive   OrganizationMemberStatus = 1
	OrganizationMemberStatusInactive OrganizationMemberStatus = 2
	OrganizationMemberStatusPending  OrganizationMemberStatus = 3
)

type OrganizationMember struct {
	ID             uint32
	OrganizationID uint32
	UserID         uint32
	RoleId           uint32
	Status         OrganizationMemberStatus
	JoinedAt       time.Time
	RemovedAt      *time.Time

	Organization *Organization
}
