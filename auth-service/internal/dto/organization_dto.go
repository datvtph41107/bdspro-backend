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
