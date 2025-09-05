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
	UserID         uint32
	RoleId         uint64
	Status         OrganizationMemberStatus
	JoinedAt       time.Time
	RemovedAt      *time.Time

	Organization *Organization
	Role         *OrganizationRole
	RoleName     string
	DealMember   *DealInvitation
}

func (m *OrganizationMember) GetRoleKey() string {
	if m.Role != nil {
		return m.Role.Key
	}
	return ""
}

func (m *OrganizationMember) GetPermissions() []string {
	if m.Role != nil {
		permissions := make([]string, len(m.Role.Permissions))
		for i, p := range m.Role.Permissions {
			permissions[i] = p.Key
		}
		return permissions
	}
	return nil
}
