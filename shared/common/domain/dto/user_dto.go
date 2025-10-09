package _dto

type ColorDTO struct {
	ID              uint64 `json:"id"`
	Code            string `json:"code"`
	Color           string `json:"color"`
	ContentColor    string `json:"contentColor"`
	BackgroundColor string `json:"backgroundColor"`
	ColorKey        string `json:"colorKey"`
	HexCode         string `json:"hexCode"`
}

type RoleDTO struct {
	ID         uint64    `json:"id"`
	RoleName   string    `json:"roleName"`
	RoleKey    uint32    `json:"roleKey"`
	ProfileID  uint64    `json:"profileId"`
	Color      *ColorDTO `json:"color"`
	Status     string    `json:"status"`
	StatusText string    `json:"statusText"`
}

type AuthUserProfileDTO struct {
	ProfileID uint64   `json:"profileId"`
	FullName  string   `json:"fullName"`
	Avatar    string   `json:"avatar"`
	Role      *RoleDTO `json:"role"`
	RoleID    uint64   `json:"roleId"`
	IsLocked  bool     `json:"isLocked"`
	LockType  string   `json:"lockType"`
}
