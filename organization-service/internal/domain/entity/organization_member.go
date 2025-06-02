package entity

import (
	"time"
)

type OrganizationMemberRole string

const (
	OrganizationMemberRoleAdmin  OrganizationMemberRole = "admin"
	OrganizationMemberRoleMember OrganizationMemberRole = "member"
)

type OrganizationMemberStatus string

const (
	OrganizationMemberStatusActive   OrganizationMemberStatus = "active"
	OrganizationMemberStatusInactive OrganizationMemberStatus = "inactive"
)

type OrganizationMember struct {
	ID             uint32
	OrganizationID uint32
	UserID         uint32
	Role           OrganizationMemberRole
	Status         OrganizationMemberStatus
	JoinedAt       time.Time
	RemovedAt      *time.Time

	Organization *Organization
}
