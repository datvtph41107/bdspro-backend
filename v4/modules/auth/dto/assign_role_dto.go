package dto

// AssignRoleToUserRequest DTO cho request gán role cho user
type AssignRoleToUserRequest struct {
	UserID uint64 `json:"userId" validate:"required"`
	RoleID uint64 `json:"roleId" validate:"required"`
}

// AssignRoleToUserResponse DTO cho response gán role cho user
type AssignRoleToUserResponse struct {
	AssignedRole RoleDTO `json:"assignedRole"`
	Message      string  `json:"message"`
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
