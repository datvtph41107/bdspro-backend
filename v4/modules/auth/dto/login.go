package dto

// HTTP-only login DTOs (path/query binding). Shared login types live in auth.go.

type RefreshTokenHTTPRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type PhoneCheckResponse struct {
	Existed bool    `json:"existed"`
	UserId  *uint64 `json:"userId,omitempty"`
}

// PhoneCheckPath is GET {prefix}/auth/phone/check/{phone}.
type PhoneCheckPath struct {
	Phone string `path:"phone"`
}
