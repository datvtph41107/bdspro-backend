package dto

type RoleDTO struct {
	ID          uint64 `json:"id"`
	RoleName    string `json:"roleName"`
	RoleKey     string `json:"roleKey"`
	RoleColor   string `json:"roleColor"`
	RoleBgColor string `json:"roleBgColor"`
}
