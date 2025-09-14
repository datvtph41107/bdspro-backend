package dto

// AssignRoleToUserRequest DTO cho request gán role cho user
type AssignRoleToUserRequest struct {
	UserID         uint64   `json:"userId" validate:"required"`
	RoleKeys       []uint32 `json:"roleKeys" validate:"required,min=1"`
	OrganizationID uint64   `json:"organizationId" validate:"required"`
}

// AssignRoleToUserResponse DTO cho response gán role cho user
type AssignRoleToUserResponse struct {
	AssignedRoles []RoleDTO `json:"assignedRoles"`
	Message       string    `json:"message"`
}

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
}

// PermissionDTO DTO cho permission
type PermissionDTO struct {
	ID             uint64          `json:"id"`
	Name           string          `json:"name"`
	Key            string          `json:"key"`
	Description    string          `json:"description"`
	PermissionType string          `json:"permissionType"`
	Module         string          `json:"module"`
	ParentID       *uint64         `json:"parentId,omitempty"`
	Children       []PermissionDTO `json:"children,omitempty"`
	CreatedAt      string          `json:"createdAt"`
	UpdatedAt      string          `json:"updatedAt"`
	Priority       uint32          `json:"priority"`
}

// ColorDTO DTO cho color
type ColorDTO struct {
	ID              uint64  `json:"id"`
	Name            string  `json:"name"`
	ContentColor    string  `json:"contentColor"`
	BackgroundColor string  `json:"backgroundColor"`
	ColorKey        string  `json:"colorKey"`
	HexCode         string  `json:"hexCode"`
	Description     *string `json:"description,omitempty"`
	IsActive        bool    `json:"isActive"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
	CreatedBy       uint32  `json:"createdBy"`
	UpdatedBy       uint32  `json:"updatedBy"`
}
