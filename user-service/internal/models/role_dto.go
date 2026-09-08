package models

// RoleDTO chứa thông tin role từ auth service
type RoleDTO struct {
	ID              uint64 `json:"id"`
	RoleName        string `json:"roleName"`
	RoleKey         uint32 `json:"roleKey"`
	RoleDescription string `json:"roleDescription"`
}
