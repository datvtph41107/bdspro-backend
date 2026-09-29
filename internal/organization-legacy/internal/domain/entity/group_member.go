package entity

import "time"

type GroupMemberRole string

const (
	GroupMemberRoleAdmin  GroupMemberRole = "admin"
	GroupMemberRoleMember GroupMemberRole = "member"
)

type GroupMemberStatus string

const (
	GroupMemberStatusActive   GroupMemberStatus = "active"
	GroupMemberStatusInactive GroupMemberStatus = "inactive"
)

type GroupMember struct {
	Id        uint32
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uint32
	UpdatedBy uint32

	GroupId uint32
	UserId  uint32
	Role    GroupMemberRole
	RoleId  *uint64
	Status  GroupMemberStatus
}
