package dto

// RoleDTO DTO cho role với thông tin màu sắc
type RoleDTO struct {
	ID              uint64          `json:"id"`
	RoleName        string          `json:"roleName"`
	RoleDescription string          `json:"roleDescription"`
	Key             string          `json:"key"`
	PermissionIDs   []uint64        `json:"permissionIds"`
	OrganizationID  uint64          `json:"organizationId"`
	Permissions     []PermissionDTO `json:"permissions"`
	IsDefault       bool            `json:"isDefault"`
	DomainType      string          `json:"domainType"`
	Color           *ColorDTO       `json:"color,omitempty"`
	RoleGroupID     *uint64         `json:"roleGroupId,omitempty"`
	RoleGroupName   *string         `json:"roleGroupName,omitempty"`
	AllowAssign     bool            `json:"allowAssign"`
	RoleKey         uint32          `json:"roleKey"`
	CreatedAt       string          `json:"createdAt"`
	UpdatedAt       string          `json:"updatedAt"`

	RoleColor   string `json:"roleColor"`
	RoleBgColor string `json:"roleBgColor"`
}
