package dto

// SwitchOrganizationRequest đại diện cho yêu cầu chuyển đổi tổ chức
type SwitchOrganizationRequest struct {
	OrganizationID uint64 `json:"organizationId" binding:"required" validate:"required"`
}

type SwitchOrganizationResponse struct {
	AuthId         uint64 `json:"authId"`
	ProfileId      uint64 `json:"profileId"`
	OrganizationId uint64 `json:"organizationId"`
	Role           string `json:"role"`
	FullName       string `json:"fullName"`
	Avatar         string `json:"avatar"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	AccessToken    string `json:"accessToken"`
}

type GroupMember struct {
	ID      uint64  `json:"id"`
	GroupId uint64  `json:"groupId"`
	UserId  uint64  `json:"userId"`
	Role    string  `json:"role"`
	RoleId  *uint64 `json:"roleId"`
}

type DealMember struct {
	DealId uint64  `json:"dealId"`
	UserId uint64  `json:"userId"`
	RoleId *uint64 `json:"roleId"`
}
