package entity

import (
	"time"
)

type OrganizationMemberStatus uint32

const (
	OrganizationMemberStatusActive   OrganizationMemberStatus = 10
	OrganizationMemberStatusInactive OrganizationMemberStatus = 20
	OrganizationMemberStatusPending  OrganizationMemberStatus = 30
)

type OrganizationMember struct {
	ID             uint32
	OrganizationID uint32
	UserID         uint64
	RoleId         uint64
	RoleKey        uint32
	Status         OrganizationMemberStatus
	JoinedAt       time.Time
	RemovedAt      *time.Time

	Organization *Organization
	RoleName     string
	DealMember   *DealMember
	// Role         *OrganizationRole
}

func (m *OrganizationMember) GetRoleKey() uint32 {
	if m.RoleKey != 0 {
		return m.RoleKey
	}
	// if m.Role != nil {
	// 	return m.Role.RoleKey
	// }
	return 0
}

func (m *OrganizationMember) GetPermissions() []string {
	// if m.Role != nil {
	// 	permissions := make([]string, len(m.Role.Permissions))
	// 	for i, p := range m.Role.Permissions {
	// 		permissions[i] = p.Key
	// 	}
	// 	return permissions
	// }
	return nil
}
