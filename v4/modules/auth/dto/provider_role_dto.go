package dto

// ProviderRoleDTO chứa thông tin role từ auth service (client contract).
type ProviderRoleDTO struct {
	ID              uint64 `json:"id"`
	RoleName        string `json:"roleName"`
	RoleKey         uint32 `json:"roleKey"`
	RoleDescription string `json:"roleDescription"`
}
